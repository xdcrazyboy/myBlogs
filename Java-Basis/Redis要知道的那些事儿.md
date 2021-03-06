## 老规矩，显示参考资料
1. [Redis是单线程的，但Redis为什么这么快？](https://zhuanlan.zhihu.com/p/42272979)
2. [Redis容灾策略](https://blog.csdn.net/irean_lau/article/details/51360277)


## Redis是什么

* Redis是一个开源的内存中的数据结构存储系统，它可以用作：**数据库、缓存和消息中间件**。

* Redis 内置了**复制**（Replication），LUA脚本（Lua scripting）， LRU驱动事件（LRU eviction），**事务**（Transactions） 和**不同级别的磁盘持久化**（Persistence），并通过 Redis**哨兵**（Sentinel）和**自动分区**（Cluster）**提供高可用性**（High Availability）。

* 提供了持久化的选项,根据实际情况，可以每隔一定时间将数据集导出到磁盘（快照），或者追加到命令日志中（AOF只追加文件），他会在执行写命令时，将被执行的写命令复制到硬盘里面。 也可以关闭持久化功能，只使用高效网络数据缓存功能。

* 不用表，不需要对数据进行关联

* 数据库的工作模式按存储方式可分为：**硬盘数据库和内存数据库**。Redis 将数据储存在内存里面，读写数据的时候都不会受到硬盘 I/O 速度的限制，**所以速度极快**。

**比较**：速度不比采用**单进程多线程**的同样**基于内存的 KV 数据库** Memcached 差！

## 为什么快？
1. 基于**内存的**，大部分操作都是在内存中。
2. **数据结构简单**，对数据操作也简单
3. 采用**单线程**，避免了不必要的**上下文切换和竞争条件**；无多线程，避免了**线程切换**而消耗 CPU；无锁，不存在**加锁释放锁操作**，**无死锁**消耗。
4. 多路I/O复用模型，非阻塞IO；
   > 多路I/O复用模型是利用 **select、poll、epoll 可以同时监察多个流的 I/O 事件**的能力，在**空闲**的时候，会把**当前线程阻塞掉**，当有一个或多个流有 I/O 事件时，就从阻塞态中唤醒，于是程序就会轮询一遍所有的流（**epoll 是只轮询那些真正发出了事件的流**），并且只依次顺序的处理就绪的流，这种做法就**避免了大量的无用操作**。
5. 使用**底层模型**不同，它们之间底层实现方式以及与客户端之间通**信的应用协议**不一样，Redis直接自己构建了**VM机制**，因为一般的系统**调用系统函数**的话，**会浪费一定的时间去移动和请求**；

**多路复用**：“多路”指的是**多个网络连接**，“复用”指的是**复用同一个线程**。采用多路 I/O 复用技术可以让单个线程高效的处理多个连接请求（尽量**减少网络 IO 的时间消耗**），且 Redis 在内存中操作数据的速度非常快，也就是说内存内的操作不会成为影响Redis性能的瓶颈，主要由以上几点造就了 Redis 具有很高的吞吐量。

**为什么用单线程**：因为Redis是基于内存的操作，**CPU不是Redis的瓶颈**，Redis的瓶颈最有可能是机器内存的大小或者网络带宽。既然单线程容易实现，而且CPU不会成为瓶颈。

多核CPU不是浪费了？ 可以开多个redis配合使用，新版本也有写情况可以使用多线程。 避免耗时操作，影响redis的并发能力。

**扩展**：
1、单进程多线程模型：MySQL、Memcached、Oracle（Windows版本）；

2、多进程模型：Oracle（Linux版本）；

3、Nginx有两类进程，一类称为Master进程(相当于管理进程)，另一类称为Worker进程（实际工作进程）。启动方式有两种：

（1）单进程启动：仅有一个进程，充当Master进程，也充当Worker进程的角色。

（2）多进程启动：有且仅有一个Master进程，至少有一个Worker进程工作。

Master进程主要进行一些全局性的初始化工作和管理Worker的工作；事件处理是在Worker中进行的。



## 哪些数据结构

1. **字符串（String）**
2. **散列（Hash）**
3. **列表（List）**
4. **集合（Set）**
5. **有序集合**（Sorted Set或者是**ZSet**）与范围查询
6. Bitmaps
7. Hyperloglogs
8. 地理空间（Geospatial）索引半径查询

## 持久化

**RDB持久化**——可以在指定的**时间间隔**内生成数据集的时间点**快照**（point-in-time snapshot）

* RDB优点：
    1. 文件紧凑，保存了某个时间点上的数据仅，**适合用于备份**。
    2. 回合灾难恢复，它只有一个文件，可以加密后放到别的数据中心。
    3. 可以最大化redis的性能，在保存RDB文件时只需要fork一个子进程，让它负责处理就行，父进程无需执行任何I、O操作。
    4. 在恢复**大数据集**速度比AOF快。
* RDB缺点：
    1. 出故障停机，会丢几分钟数据（需要保存整个数据集状态，消耗大，肯定会有间隔去进行保存）
    2. fork子进程在数据集庞大，cpu时间紧张时，可能会比较耗时，这期间会停止处理客户端。AOF也会，但是持久性更好。

**AOF持久化**——**记录**服务器执行的所有**写操作**命令，并在服务器启动时，通过**重新执行这些命令**来还原数据集。 以 Redis协议的格式来保存，新命令会被追加到文件的末尾。 
> Redis 还可以在**后台对AOF文件**进行**重写**（rewrite），使得 AOF 文件的**体积不会超出**保存数据集状态所需的实际大小。

* AOF优点
  1. 让redis更耐久。fsync设置为一秒，也能保持良好性能，出故障也就丢失1秒数据，还可以设置无fsync。
  2. AOF文件是只进行追加操作的日志文件，中途停机，没写完整，也可以用redis-check-aof轻易修复
  3. 可以重写AOF，缩小体积，优化不错。
  4. 有序保存了对数据库执行的所有写入操作，很容易读懂，对文件分析很方便。导出也方便。
* AOF缺点
  1. AOF体积大于RDB。
  2. 采用fsync时，速度可能慢于AOF，处理巨大数据载入时，RDB可以提供更有保证的最大延迟时间。
  3. 有数据没办法恢复的bug？以前？

**两种都用**： 当 Redis 重启时， 它会**优先使用AOF**文件来还原数据集， 因为 AOF 文件保存的数据集通常比 RDB 文件所保存的数据集**更完整**.

## 缓存

## Redis容灾策略

**主从备份+哨兵**： 
1. 采用master-slave方式 
2. 为了得到好的读写性能，master不做任何的持久化 
3. slave同时开启Snapshot和AOF来进行持久化，保证数据的安全性 
4. 当master挂掉后，修改slave为master 
5. 恢复原master数据，修改原先master为slave，启动slave 
6. 若master与slave都挂掉后，调用命令通过aof和snapshot进行恢复恢复时要先确保恢复文件都正确了，才能启动主库；
>也可以先启动slave，将master与slave对调开源方案[codis](http://navyaijm.blog.51cto.com/4647068/1637688)

**哨兵的作用**

1. 监控：监控主从是否正常
2. 通知：出现问题时，可以通知相关人员
3. 故障迁移：自动主从切换
4. 统一的配置管理：连接者询问sentinel取得主从的地址 

Raft算法核心: 可视图
[Raft Visualization (算法演示)](http://thesecretlivesofdata.com/raft/)

[使用主从结构+哨兵（sentinel）来进行容灾](http://blog.csdn.net/liuwei063608/article/details/50520163)

