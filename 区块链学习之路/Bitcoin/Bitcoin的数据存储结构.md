# Bitcoin Core 0.11 (ch 2): Data Storage
原文链接：[https://en.bitcoin.it/wiki/Bitcoin_Core_0.11_(ch_2):_Data_Storage](https://en.bitcoin.it/wiki/Bitcoin_Core_0.11_(ch_2):_Data_Storage)

本文描述了比特币核心如何以及在何处存储区块链数据。

# 概观
There are basically four pieces of data that are maintained:

blocks/blk*.dat: the actual Bitcoin blocks, in network format, dumped in raw on disk. They are only needed for rescanning missing transactions in a wallet, reorganizing to a different part of the chain, and serving the block data to other nodes that are synchronizing.
blocks/index/*: this is a LevelDB database that contains metadata about all known blocks, and where to find them on disk. Without this, finding a block would be very slow.
chainstate/*: this is a LevelDB database with a compact representation of all currently unspent transaction outputs and some metadata about the transactions they are from. The data here is necessary for validating new incoming blocks and transactions. It can theoretically be rebuilt from the block data (see the -reindex command line option), but this takes a rather long time. Without it, you could still theoretically do validation indeed, but it would mean a full scan through the blocks (7 GB as of may 2013) for every output being spent.
blocks/rev*.dat: these contain "undo" data. You can see blocks as 'patches' to the chain state (they consume some unspent outputs, and produce new ones), and see the undo data as reverse patches. They are necessary for rolling back the chainstate, which is necessary in case of reorganisations.
Note that the LevelDB's are redundant in the sense that they can be rebuilt from the block data. But validation and other operations would become intolerably slow without them.

See here: StackExchange post by Pieter Wuille (2013)


# Raw Block data (blk*.dat)
Block files store the raw blocks as they were received over the network.

Block files are about 128 MB, allocated in 16 MB chunks to prevent excessive fragmentation. As of October 2015, the block chain is stored in about 365 block files, for a total of about 45 GB.

Each block file (blk1234.dat) has a corresponding undo file (rev1234.dat) which contains the data necessary to remove blocks from the blockchain in the event of a reorganization (fork).

Info about the block files is stored in the block index (the LevelDB) in two places:

General info about the files themselves is held in the "f" records in the block index LevelDB (meaning keys "fxxxx", where "xxxx" is the 4 digit file number), including:
Number of blocks stored in the file
File size (and the corresponding undo file size)
Lowest and highest block in the file
Timestamps - earlier and latest blocks in the file
Info about where to find a particular block on disk is in the "b" ("b" = block) record:
Each block contains a pointer to the block is on disk (a file number and an offset)

Accessing the block data files from the code

The block files are accessed through:

1) DiskBlockPos: a struct that is simply a pointer to a block's location on disk (a file number and an offset.)

2) vInfoBlockFiles: a vector of BlockFileInfo objects. This variable is used to perform such tasks as:

Determine whether new blocks can fit into the current file or a new file needs to be created
Calculate the total disk usage by block & undo files
Iterate through the block files and find ones that can be pruned
Blocks are written to disk as soon as they are received, in AcceptBlock. (The actual disk write operation is in WriteBlockToDisk [main.cpp:1164]). Note that there is some overlap of the code that accesses block files with the code that accesses and writes to the coins database (/chainstate). There is a complex system of when to flush state to disk. None of this code affects block files, which are simply written to disk when received. Once they have been received and stored, the block files are only needed for serving blocks to other nodes.


More info about block files

See here: the commit that puts multiple blocks in a block file (2012)

# Block index (leveldb)
The block index holds metadata about all known blocks, including where the block is stored on disk.

Note that the set of "known blocks" is a superset of the longest chain, because it includes blocks that were received and processed but are not part of the active chain - for example, orphaned blocks that were detached from the active chain in a small reorganization.

Terminology

The terminology can be a little confusing here, because while people normally think of the "blockchain" as being synonymous with the active chain (an uninterrupted, linear chain of X blocks starting with the genesis block and continuing to the current tip), there are some places in the code where "blockchain" refers to the active chain plus the numerous, mostly short forks off the chain that our node happens to know about.

a) Block Tree

A better term for the set of known blocks stored on disk is "block tree," as this term contemplates a tree structure with numerous branches (albeit small ones) from the main chain. Indeed, the block index LevelDB is accessed through the "CBlockTreeDB" wrapper class, defined in src/txdb.h. Note that it's perfectly fine, indeed it is expected, that different nodes would have slightly different block trees; what matters is that they agree on the active chain.

Key-value pairs

Inside the actual LevelDB, the used key/value pairs are:

   'b' + 32-byte block hash -> block index record. Each record stores:
       * The block header
       * The height.
       * The number of transactions.
       * To what extent this block is validated.
       * In which file, and where in that file, the block data is stored.
       * In which file, and where in that file, the undo data is stored.
  'f' + 4-byte file number -> file information record. Each record stores:
       * The number of blocks stored in the block file with that number.
       * The size of the block file with that number ($DATADIR/blocks/blkNNNNN.dat).
       * The size of the undo file with that number ($DATADIR/blocks/revNNNNN.dat).
       * The lowest and highest height of blocks stored in the block file with that number.
       * The lowest and highest timestamp of blocks stored in the block file with that number.
   'l' -> 4-byte file number: the last block file number used.
   'R' -> 1-byte boolean ('1' if true): whether we're in the process of reindexing.
   'F' + 1-byte flag name length + flag name string -> 1 byte boolean ('1' if true, '0' if false): various flags that can be on or off. Currently defined flags include:
        * 'txindex': Whether the transaction index is enabled.
   't' + 32-byte transaction hash -> transaction index record. These are optional and only exist if 'txindex' is enabled (see above). Each record stores:
       * Which block file number the transaction is stored in.
       * Which offset into that file the block the transaction is part of is stored at.
       * The offset from the start of that block to the position where that transaction itself is stored.

See here: StackExchange post by Pieter Wuille (2014)


Data Access Layer

The database is accessed through CBlockTreeDB wrapper class. See txdb.h.

The wrapper is instantiated in a global variable called pblocktree, defined in main.cpp.

CBlockIndex

Blocks stored in the database are represented in memory as CBlockIndex objects. An object of this type is first created after the header is received; the code does not wait to receive the full block. When headers are received over the network, they are streamed into a vector of CBlockHeaders, which are then checked. Each header that checks out causes a new CBlockIndex to be created, which is stored to the database.

CBlock / CBlockHeader

Note that these objects have little to do with the /blocks LevelDB. A CBlock holds the full set of transactions in the block, the data for which is stored in two places - in full, in raw format, in the blk???.dat files, and in pruned format in the UTXO database. The block index database cares not for such details, since it holds only the metadata for the block.

Loading the block database into memory

The entire database is loaded into memory on startup. See LoadBlockIndexGuts (txdb.cpp). This only takes a few seconds.

The blocks ('b' keys) are loaded into the global "mapBlockIndex" variable. "mapBlockIndex" is an unordered_map that holds CBlockIndex for each block in the entire block tree; not just the active chain.

mapBlockIndex is described in more detail in Chapter 6 - The Blockchain.

The block file metadata ('f' keys) is loaded into vInfoBlockFiles.

# The UTXO set (chainstate leveldb)
The UTXO database was introduced in 2012 in pull request #1677 - "Ultraprune."

The idea behind "Ultraprune" is to reduce the size of (prune) the set of past transactions, keeping only those parts of past transactions that are necessary to validate later transactions.

Say you have a transaction T1 which takes two inputs and sends to 3 outputs: O1,O2,O3. Two of those outputs (O1, O2) have been used as inputs in a later transaction, T2. Once T2 has been mined, T1 only has one item of interest (O3). There's no reason to keep T1 around in its entirety. Instead, a slimmed-down version of T1 will suffice, consisting only of O3 (locking script and amount) and certain basic information about T1 (height, whether it is a coinbase, etc.)

The description of ultraprune is on the specific "ultraprune" commit within the pull:

-------------
This switches bitcoin's transaction/block verification logic to use a "coin database", which contains all unredeemed transaction output scripts, amounts and heights.
The name ultraprune comes from the fact that instead of a full transaction index, we only (need to) keep an index with unspent outputs. For now, the blocks themselves are kept as usual, although they are only necessary for serving, rescanning and reorganizing.
The basic data structures are CCoins (representing the coins of a single transaction), and CCoinsView (representing a state of the coins database). There are several implementations for CCoinsView. A dummy, one backed by the coins database (coins.dat), one backed by the memory pool, and one that adds a cache on top of it. FetchInputs, ConnectInputs, ConnectBlock, DisconnectBlock, ... now operate on a generic CCoinsView.
The block switching logic now builds a single cached CCoinsView with changes to be committed to the database before any changes are made. This means no uncommitted changes are ever read from the database, and should ease the transition to another database layer which does not support transactions (but does support atomic writes), like LevelDB.
For the getrawtransaction() RPC call, access to a txid-to-disk index would be preferable. As this index is not necessary or even useful for any other part of the implementation, it is not provided. Instead, getrawtransaction() uses the coin database to find the block height, and then scans that block to find the requested transaction. This is slow, but should suffice for debug purposes.
-----------------

See: Ultraprune - July 2012

## 一些术语

- "UTXO (Unspent Transaction Out):" An output from a transaction. This is colloquially referred to as a "coin." For this reason, the UTXO db is sometimes referred to as the "coins database."

"UTXO set / coins database / chainstate database:" These terms are more or less synonymous and are used interchangeably.

- "Provably Unspendable:" A coin is provably unspendable if its scriptPubKey cannot be satisfied - for example, an OP_RETURN. A provably unspendable coin can be eliminated from the utxo database regardless of its amount.


## Key-value pairs

The records in the chainstate levelDB are:

  > 'c' + 32-byte transaction hash -> unspent transaction output record for that transaction. These records are only present for transactions that have at least one unspent output left. Each record stores:
    > - The version of the transaction.
    > - Whether the transaction was a coinbase or not.
    > - Which height block contains the transaction.
    > - Which outputs of that transaction are unspent.
    > - The scriptPubKey and amount for those unspent outputs.


   'B' -> 32-byte block hash: the block hash up to which the database represents the unspent transaction outputs.

See here: [StackExchange post by Pieter Wuille (2014)](http://bitcoin.stackexchange.com/questions/28168/what-are-the-keys-used-in-the-blockchain-leveldb-ie-what-are-the-keyvalue-pair)

## 数据访问层和缓存

访问UTXO数据库比块索引复杂得多。这是因为它的性能对比特币系统的整体性能至关重要。块索引对性能不是那么关键，因为只有几十万个块，并且运行在不错的硬件上的节点可以在几秒钟内检索并滚动它们（并且不需要经常这样做。）另一个在UTXO数据库中有数百万个硬币，必须对每个进入mempool或包含在块中的事务的每个输入进行检查和修改。

正如sipa在ultraprune提交中所说：

> 基本数据结构是CCoins（代表单个交易的硬币）和CCoinsView（代表硬币数据库的状态）。CCoinsView有几种实现方式。一个虚拟的，一个由硬币数据库（coins.dat）支持，一个由内存池支持，另一个在其上添加缓存。

This is not stated as clearly as it might have been, however; at least, not for the current state of the code.

在0.11中，CoinsView的实例化是：

- dummy
- database
- pCoinsTip (a cache backed by the database)
- "validation cache" (u在由pCoinsTip支持时使用，在连接块时使用)\
  
Separate from that chain of caches is the memory pool's CoinsView, which is backed by the database.


The class diagram (data types) for the views is:

      CCoinsView (abstract class)
             /            \
         ViewDB          ViewBacked 
      (database)          /      \
                   ViewMempool   ViewCache

Each class has one key characteristic:

>- View is the base class, declaring methods for verifying that coins exist (HaveCoins), retrieving coins (GetCoins), etc.
>- ViewDB has code to interact with the LevelDB.
>- ViewBacked has a pointer to another View; thus it is "backed" by another view (version) of the UTXO set.
>- ViewCache has a cache (a map of CCoins).
>- ViewMempool associates a mempool with a view.

那些是定义的类;而对象图是：:

            Database       
           /       \
       MemPool     Blockchain cache (pcoinsTip) 
     View/Cache            \
                         Validation cache

这是一个总结视图实例的表：:

Object	Type	Backed By?	Description / Purpose
DB view	ViewDB	n/a	Represents the UTXO set according to the /chainstate LevelDB. Retrieves coins and flushes changes to the LevelDB.
Creation in code (instantiation): see init.cpp:1131
pCoinsTip
(blockchain cache)	ViewCache	DB view	Holds the UTXO set corresponding to the active chain's tip. Retrieves/flushes to the database view.
Creation in code: see init.cpp:1133
Validation cache	ViewCache	pCoinsTip	This cache's lifetime is within ConnectTip (or DisconnectTip).
Its purpose is to keep track of modifications to the UTXO set while processing a block.
If the block validates, the cache is flushed to pcoinsTip.
If the block fails, the cache is discarded. 
Creation in code: see main.cpp:2231: CCoinsViewCache view(pcoinsTip);
Mempool view	ViewMemPool	pCoinsTip	This object brings the mempool into view, meaning it can see both a UTXO set and the mempool.
Its purpose is to enable validation of chains of transactions, a.k.a. "zero-confirmation" transactions. (If chains of transactions weren't permitted, the mempool could simply validate against pcoinsTip.)
Thus, when queried, it can check if a given input can be found either in the mempool (i.e., "zero-conf") or in the blockchain's utxo set ("confirmed.")
Note that this object is not a cache; rather, it is a view that is used by the object below, which does contain a cache. 
Creation in code: Its lifetime is that of AcceptToMemoryPool in main.cpp.
Mempool cache	ViewCache	Mempool view	The cache for the mempool. It contains a cache and sets its backend to be the mempool view.
Creation in code: Its lifetime is also that of AcceptToMemoryPool in main.cpp.


### *Loading the UTXO set*

Access to the coins database is initialized in init.cpp: 1131-1133:
``` c++
pcoinsdbview = new CCoinsViewDB(nCoinDBCache, false, fReindex);
pcoinscatcher = new CCoinsViewErrorCatcher(pcoinsdbview);
pcoinsTip = new CCoinsViewCache(pcoinscatcher);
```
The code starts by initializing a CoinsViewDB, which is equipped with methods to load coins from the LevelDB. 

The error catcher is a little hack that can be ignored. 

Next, the code initalizes pCoinsTip, which is the cache representing the state of the active chain, and is backed by the database view. 

### *Cache vs. Database*

*coins.cpp*中的`FetchCoins`函数演示了代码如何使用缓存与数据库：

``` c++
//首先，代码在缓存中搜索给定交易ID的币。
1   CCoinsMap::iterator it = cacheCoins.find(txid);
//如果找到，则返回“获取”的币。
2   if (it != cacheCoins.end())
3     return it;
4   CCoins tmp;
//如果没有，它将搜索数据库。
5   if (!base->GetCoins(txid, tmp))
6     return cacheCoins.end();
//如果在数据库中找到，它将更新缓存。
7   CCoinsMap::iterator ret = cacheCoins.insert(std::make_pair(txid,CCoinsCacheEntry())).first;
```

Note: if the cache's backend is another cache, then the term "database" really means "parent cache." 

### *将验证缓存刷新到区块链缓存*

在连接块之后，在它超出范围之前，验证缓存将刷新到区块链缓存。范围在ConnectTip中捕获，具体地，在代码块main.cpp：2231-2243中捕获。

在该代码块中，有一个ConnectBlock调用，在此期间代码将新硬币存储在验证缓存中。（具体来说，请参阅main.cpp中的UpdateCoins（）。）

在代码块的末尾，刷新验证缓存。由于其“父视图”也是一个缓存（pcoinsTip，即“区块链缓存”），代码将调用父级的ViewCache :: BatchWrite，后者将更新的硬币条目交换到自己的缓存中。（运行中的多态性：稍后，当区块链缓存刷新到数据库视图时，代码将运行CoinsViewDB :: BatchWrite，其最后一行写入LevelDB。）

总之，验证缓存的使用很简单：它在上述代码块中被实例化、使用、刷新，并超出范围。



### *将区块链缓存刷新到数据库*

Flushing the validate cache was simple because the code only shuffled items between two caches in memory (of which no one is aware outside of the caching code.) Flushing the blockchain cache to the database is a bit more complicated. 

At the lowest level, the mechanics of flushing the blockchain cache (pcoinsTip) is the same as the validation cache: the Flush() method calls BatchWrite on its backend (the "base" pointer), and in this case that means BatchWrite on the database view.

 Up a level, Flush() is called from FlushStateToDisk (FSTD) - main.cpp:2098. 

FlushStateToDisk使用给定模式在几个不同的点调用：:


Flush Mode	Description	When called
IF_NEEDED	Flush only if the cache is over its size limit.	Right after connecting (or disconnecting) a block and flushing the validation cache.
See ConnectTip / DisconnectTip.
ALWAYS	Flush cache.	During initialization only.
PERIODIC	Here, the code considers other data points to decide whether to flush.
Is the code almost over its size limit?
Has it been a long time since the cache was flushed?
If so, then proceed.	At end of ActivateBestChain()
(Code comment: "write changes periodically to disk, after relay").

The idea is to flush the block cache frequently (to avoid having to download a large number of blocks if the program crashes), but the coins cache infrequently (in order to maximize the benefit from the coins cache.)

Specifically, the block cache is guaranteed to be flushed once an hour, whereas the coins cache once per day. (See here: Sipa comment on PR 6102)

The FlushStateToDisk code is well-commented so for more info, the curious reader can check main.cpp.

# Raw undo data (rev*.dat)

The undo data contains the information that is necessary to disconnect or "roll back" a block: specifically, the coins that were spent by the block in question.

So, the data being written is essentially a set of CTxOut objects. (A CTxOut is simply an amount and a script - see primitives/transaction.h:107-108).

The matter is complicated slightly by the fact that if the coin is the last one being spent by its transaction, the undo data needs to store the transaction's metadata (the txn's block height, whether it's a coinbase, and its version.) So, if you have a transaction T with outputs O1,O2,O3 spent in that order, for O1 and O2 all that will be written to the undo file is the amount and the script. For 03, the undo file will have the amount, the script, plus T's height and version, and whether T is a coinbase.

The undo data is written to the raw file with the following code:

fileout << blockundo; (main.cpp:1567 [UndoWriteToDisk])
This line of code calls the serialization function on the CBlockUndo - which is basically just a vector of coins (CTxOuts.) Finally, a checksum is written to the undo file. The checksum is used during initialization to verify that any undo data being checked is intact. See Pull 2145

The undo data is used when disconnecting a block. The DisconnectBlock() code is discussed further down this wiki page in The Blockchain: Reorganizations.

# Use of LevelDB
LevelDB is a key-value store that was introduced to store the block index and UTXO set (chainstate) in 2012 as part of the complex "Ultraprune" pull (PR 1677). See here: the 27 commits on Ultraprune.

On the subject of why LevelDB is used, core developer Greg Maxwell stated the following to the bitcoin-dev mailing list in October 2015:

I think people are falling into a trap of thinking "It's a <database>, I know a <black box> for that!"; but the application and needs are very specialized here. . . It just so happens that on the back of the very bitcoin specific cryptographic consensus algorithim there was a slot where a pre-existing high performance key-value store fit; and so we're using one and saving ourselves some effort...
One might ask whether different nodes could use different databases - as long as they retrieve the same data, what's the difference? The issue here is "bug-for-bug compatibility" - if one database has a bug that causes records to not be returned under certain circumstances, then all other nodes bst have the same bug, else the network could fork as a result.

Greg Maxwell stated the following in the same thread referenced above (in response to a proposal to switch to using sqlite):

...[D]atabases sometimes have errors which cause them to fail to return records, or to return stale data. And if those exist consistency must be maintained; and "fixing" the bug can cause a divergence in consensus state that could open users up to theft.
Case in point, prior to leveldb's use in Bitcoin Core it had a bug that, under rare conditions, could cause it to consistently return not found on records that were really there. . . Leveldb fixed this serious bug in a minor update. But deploying a fix like this in an uncontrolled manner in the bitcoin network would potentially cause a fork in the consensus state; so any such fix would need to be rolled out in an orderly manner.