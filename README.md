# myBlogs
没有什么顺序，纯粹是给自己做的一个笔记，大部分是在巨人肩膀上汇总。
### 参考文章
参考了很多大牛的文章或者书，需要放在最显眼的位置。

另外由于收集得太多，很多都忘记出处，如果谁发现有自己的或者哪位大神的内容被我参考了，请直接在评论处贴链接，我把它们放到这里：
1. 《MySQL DBA 修炼之道》
2. [MySQL索引优化全攻略（菜鸟）](www.runoob.com/w3cnote/mysql-index.html)
3. [mysql性能优化之索引优化](https://www.cnblogs.com/yyjie/p/7486975.html)，这博客写得挺好。
4. [数据库--视图的基本概念以及作用](https://blog.csdn.net/buhuikanjian/article/details/53105416)
5. 《Java网络编程》Elliotte Rusty Harold著

# MySQL（数据库相关）

## MySQL的优点
* 开源免费又高效
* 良好的安全连接，自带查询解析、sql语句优化
* 使用读写锁（细化到行）
* 事物隔离和多版本并发控制提高并发
* 完备的事务日志记录
* 强大的存储引擎提供高效查询（表记录可达百万级）
* 如果是InnoDB，还可在崩溃后进行完整的恢复

## 事务处理
**事务**：恢复和并发控制的**基本单位**，四个特性（ACID酸性）：
1. **原子性**（Atomicity）：不可拆分，要么做（提交）要么不做（回滚）；MySQL通过**redo log**重做日志实现原子性。
> 在执行SQL语句前，会先写入redo log **buffer**，再执行SQL语句；如果语句执行出错就会根据redo log buffer中的记录执行回滚。


2. **一致性**（Consistency）：事务应确保数据库的状态从一个**一致性状态**转变为另一个一致性状态。一致性状态的含义是：数据库中的数据**应满足约束**。通过**undo log**实现一致性。
> 在写入redo log buffer之前会写入undo log。undo log是**逻辑日志**，会根据之前的SQL语句进行相应回滚，比如之前是insert，回滚就是delete。 除了回滚，undo还有一个作用是MVCC，通过undo实现**非锁定读取**。 并且undo log也会产生redo log,因为undolog也需要持久性保护。


3. **隔离性**（Isolation）：多个事务并行执行时，一个事务的执行不应影响其他事务的执行。四种隔离级别，隔离级别越高，越能保证数据的完整性和一致性，性能也越差。解决四种问题：
> 1. 丢失更新 -> Read Uncommitted 事务可以看到其他事务更改但未提交的数据，依旧存在脏读问题，几乎没什么用。
> 2. 脏读 -> Read Committed 事务可以看到它执行的时候，其他事务已经提交的数据，解决了（不允许）脏读，允许不可重复读，没办法实现可重复读。
> 3. 不可重复读 -> Repeatable Read （解决包括前面所有1,2,3,）
>同一个事务内，同一个查询请求多次执行，则获取的记录集是相同的（也就是实现了可重复读，这个需要保留旧的行版本），但不能杜绝幻读。
> 4. 幻读 -> Serializable (解决包括前面所有)，它将锁施加在所有访问的数据上。

4. **持久性**（Durability）：一旦事务提交，则所做的修改就会永久保存到数据库。 即便系统崩溃，修改的数据也不会丢失。实现原理：
> 在事务commit之前会将，redo log buffer中的数据持久化到硬盘中的redo log file，这样在commit的时候，硬盘中就已经有修改后的数据了

> Innodb下也解决幻读 ?  
MySQL InnoDB的可重复读并不保证避免幻读，需要应用使用加锁读来保证。而这个加锁度使用到的机制就是next-key locks。

**纯undo log-》 redo/undo log**

1. 单纯的undo保证了原子性和持久性，需要事务提交之前将undo buffer数据写入磁盘undo（磁盘），**浪费大量I/O**
2. 引入redo log记录数据修改后的值，可以避免数据在事务提交之前必须写入磁盘的要求，**减少I/O**

这两者就是为了保证原子性和持久性。单纯的undo log需要两次I/O.

**MyISAM引擎不支持事务，只有InnoDB（默认隔离级别：repetable read）或者Falcon**

### **MVCC**
 全称：多版本并发控制。InnoDB有。 实现查询一些**正在被另外事务更新的行**，可以看到它们被更新之前的值。 这样查询就不用等待另一个事务释放锁。
 
 1. 给每行增加两个隐藏字段来实现MVCC，一个用来记录**数据行创建时间**，另一个用来记录**行的过期(删除)时间**.（实际操作中，存储的不是时间，而是事务的版本号）
 2. 快照读：读取的是历史版本，普通的SELECT操作就是快照读。 当前读：读取的是最新版本。

### 锁
1. Record Locks（记录锁）：在索引记录上加锁。
2. Gap Locks（间隙锁）：在索引记录之间加锁，或者在第一个索引记录之前，或者在最后一个索引记录之后。
3. Next-Key Locks：上面两者结合，都加锁。

* 利用MVCC实现一致性非锁定读，保证了可重复读

### 悲观锁和乐观锁
**悲观锁**：老认为别人会修改它所要操作的数据，实现依靠数据库底层。在数据进行提交更新的时候，才会正式对数据的冲突与否进行检测。

**乐观锁**：相反。只有在提交更新的时候去检查数据的状态。通常是给数据增加一个字段来标识数据的版本。


## 索引优化

### 索引
为特定的mysql字段进行一些特定的算法排序，比如二叉树的算法和哈希算法，优化查询速度。 

* Explain优化查询检测 可自动分析语句，提前结束搜索。
* 

## 创建索引的技巧
1. 维度高的列创建索引
2. 数据列中不重复的值出现的个数越高，维度越高。（比如性别这种就不适合建立索引）
3. 对where，on,group by,order by中出现的列建立索引
4. 对较长的字符串使用前缀索引
5. 不要过多创建索引，这会增加额外的磁盘空间，对于DML操作的速度影响很大，因为每增删一次都要从新建立索引。
6. 使用组合索引，可以减少文件索引大小，在使用速度要由于多个单列索引。


## 组合索引和前缀索引 （索引技巧不是类型）


## 避免写出一些不走索引的sql

* ` where 'age'+10=30;`--所有索引列参与了计算
* ` where left('data'.4) < 1990;`--同上
* ` like "%你好%"` -- 不走索引  `like "你好%"` --走索引
* 正则表达式不走索引
* 字符串和数字比较不走索引
* 如果条件中有`or`，需要全部条件都建立索引，所以尽量避免用`or`
* 如果mysql估计是用全表扫描要比使用索引快，则不使用索引。
* 不在索引上做任何操作（计算/函数/类型转换），会导致索引失效
* 存储引擎不能使用索引中范围条件右边的列
* 不等于（<> 或！=）会变成全表扫描
* is null ,not null也无法使用索引
* 字符串不走单引号索引失效

## 索引的弊端
* 查询操作频繁的列创建索引。因为索引会降低增加、删除、更新操作的速度，因为这些操作后需要对索引文件进行重新排序或更新。 
* 不过互联网应用，查询语句远多于DML的语句，所以，一般只要在大批数据要导入的时候，先删除索引，再批量插入数据，最后添加索引。


## 数据库视图

**视图（子查询）**：是从一个或多个表导出的虚拟的表，其内容由查询定义。具有普通表的结构，但是不实现数据存储。


**对视图的修改**：单表视图一般用于查询和修改，会改变基本表的数据，

**多表视图**一般用于**查询**，不会改变基本表的数据。

``` SQl
--创建视图--
create or replace view v_student as select * from student;
--从视图中检索数据--
select * from v_student;
--删除视图--
drop view v_student;
```
**作用：**

* 简化了操作，把经常使用的数据定义为视图。

* 安全性，用户只能查询和修改能看到的数据。

* 逻辑上的独立性，屏蔽了真实表的结构带来的影响。
* 


# JVM(内存、)

## 内存

### Java内存的分配策略

Java 程序运行时的内存分配策略有三种,分别是**静态分配**,**栈式分配**,和**堆式分配**，对应的，三种存储策略使用的内存空间主要分别是静态存储区（也称方法区）、栈区和堆区。

**静态存储区（方法区）**：主要存放**静态数据**、**全局static数据**和**常量**。这块内存在程序编译时就已经分配好，并且在程序整个运行期间都存在。

**栈区** ：**线程私有**。当方法被执行时，**方法体内的局部变量**（其中包括基础数据类型、**对象的引用**）都在栈上创建，并在方法执行结束时这些局部变量所持有的内存将会**自动被释放**。因为栈内存分配运算内置于处理器的**指令**集中，效率很高，但是分配的内存**容量有限**。

**堆区**： 又称**动态内存分配**，通常就是指在程序运行时直接 new 出来的内存，也就是**对象的实例**。这部分内存在不使用时将会由 Java 垃圾回收器来负责回收。

### 栈
#### Java虚拟机栈
**线程私有的**，它的生命周期跟线程相同。
* 每个方法的调用到完成，其实就是对应一个**栈帧**在虚拟机栈中出栈和入栈的过程。
* 虚拟机栈就是**执行Java方法的内存模型服务**：每个方法的执行时都会创建一个栈帧，用于存储：
    >  局部变量表、操作数栈、动态链接、方法出口等信息
    > 1. 局部变量表：存放编译期可知的各种**数据类型、对象应用、returnAddress类型**
    > 2. 操作数栈：大多数指令都要从这里：弹出数据，执行运算，把结果压回。
    > 3. 动态连接：每个栈帧都包含一个指向**运行时常量池**（方法区的一部分）中该帧**所属方法**的**引用**。
    > 4. 方法出口：返回方法被调用的位置，恢复上层方法的**局部变量**和**操作数栈**。类似递归时...
    > 
    > 局部表量表所需的内存空间在**编译期间**完成分配。 在方法运行期间不会改变局部变量表的大小。
* 虚拟机栈定义了两种异常：
    1. StackOverfloError：线程请求的**深度大于虚拟机所允许的深度**。
    2. OutOfMemoryError：虚拟机栈动态扩展时无法申请到足够的内存。

**与本地方法栈比较**
1. 本地方法栈为虚拟机使用到**的Native方法**服务。 HotSpot虚拟机直接把本地方法栈和虚拟机栈合二为一。
2. 虚拟机栈为虚拟机**执行Java方法（也就是字节码）**服务。

### 方法区
**线程共享**，用于存储**已被虚拟机加载的类信息、常量、静态变量、运行时常量池**，即编译后的代码，方法区也叫：**持久代**（Permanent Generation），Non-Heap（非堆）。因为它存放的信息与垃圾回收关系不打，可以**选择不实现垃圾回收**。

方法区的内存回收主要针对：**常量池的回收**和**类的卸载**。

**运行时常量池**：
* jdk1.7之后，字符串常量池已经从方法区挪到堆中了。
* 对比下，常量池：常量池数据编译期被确定，是Class文件中的一部分，存储了类、方法、接口等中的**常量**，也包括字符串常量。
* 运行时常量池：方法区的一部分，所有线程共享。虚拟机加载Class后把常量池中的数据放入**运行时常量池**。



### 内存泄露

**什么是内存泄露？**
> 指程序中己动态分配的堆内存由于某种原因程序未释放或无法释放，造成系统内存的浪费，导致程序运行速度减慢甚至系统崩溃等严重后果。 出现**可达、无用的对象**。

**如何导致的？**

长生命周期的对象持有短生命周期对象的引用就很可能发生内存泄漏，尽管短生命周期对象已经不再需要，但是因为长生命周期持有它的引用而导致不能被回收，这就是Java中内存泄漏的发生场景

1. 单例造成的内存泄漏
2. 非静态内部类创建静态实例造成的内存泄漏
> 解决办法：将该内部类设为静态内部类或将该内部类抽取出来封装成一个单例，如果需要使用Context，就使用Application的Context。
3. Handler造成的内存泄漏
> 解决办法：将Handler类独立出来或者使用静态内部类，这样便可以避免内存泄漏。
4. 线程造成的内存泄漏
> 解决办法： 将AsyncTask和Runnable类独立出来或者使用静态内部类，这样便可以避免内存泄漏。
5. 资源未关闭，监听器未关闭。

**怎么解决？**

1.  * 将内部类改为静态内部类
    * 静态内部类中使用弱引用来引用外部类的成员变量

2. 尽量避免使用 static 成员变量

## GC



# 集合

集合类的基本接口:**Collection** 、**Map**。

## 集合和数组的转换
1. 数组——>集合： `Arrays.asList`包装器可以实现
```java
    String[] values = ...;
    HashSet<String> staff = new HashSet<>(Arrays.asList(values));
```
2. 集合——>数组： 困难一些，通过`toArray()`方法得到的是对象数组，需要处理一下
``` java
   String[] values = (String[]) staff.toArrays(); //Error! 不能强制转换 
   String[] values = staff.toArrays(new String[0]); //0可以换成指定大小 staff.size()
```

## Collection接口

实现有;List(LinkedList、ArrayList)、Set(HashSet、EnumSet、TreeSet)、Queue（PriorityQueue）。

* 两个基本方法：
``` java 
boolean add(E element);
Iterator<E> iterator()
```
* 迭代器：Iterator接口包含四个方法：
``` java
E next() ;  
boolean hasNext(); 
void remove() ; 
deafult void forEachRemaining(Consumer<? super E> action);
```
* 其实next()取下一个，之前需要判断有没有下一个；最重要的是可以直接使用` for each`循环操作；
```java 
for (String element : c){
    do something with element
}
```
* Collection接口扩展了Iterator接口，所以标准类库中的**任何集合**都可以使用`for each`。顺序是从0开始，而HashSet这种，每个元素将会按照某种随机的次序出现，每次出现不一定一致，对于求总和的倒是无碍。

* next()读取元素的的同时，迭代位置会向前移动。 所以可以认为Iterator迭代器是位于*两个元素之间*。
* remove() 会删除上次调用` next()` 方法时返回的元素。删除两个元素，需要配合着调用两次`next()`
* **Collection提供的其他方法**：
``` java
//返回一个用于访问集合中每个元素的迭代器/
Iterator<E> iterator(); /它本身也有三个函数：hashNext();next();remove()

int size();
boolean isEmpty();
boolean add(Object element)；
boolean contains(Object obj);
boolean equals(Object obj);
boolean containsAll(Collection<?> other)； //这个集合包含other集合中的所有元素，返回true
boolean retainAll(Collection<?> other)； //这个集合删除other集合中的所有元素，返回true
还有remove()、removeAll()、clear()、
```
* Collection一些很nice的功能函数，用好了可以省很多事：
``` java
    min();max(); //返回集合中最小或者最大的元素
    copy(List<? super T> to, LIst<T> from);//将源列表的所有元素复制到目标列表的相应位置，目标列表长度大于等于源列表
    fill(List<? super T> l, T value); //将列表中所有位置设置为相同的值！！！！很关键。
    addALL();//将所有的值添加到集合中，如果集合改变了，返回true
    a.removeAll(b);//删除a中所有在b中出现的元素；
    a.retainAll(b);//保留...同上
    replaceAll(List<T> l, T oldValue, T newValue); //将集合中所有oldValue替换成newValue
    swap(list,int i,int j); //交换他们的位置
    reverse(list); //把列表元素逆序 ，时间复杂度是 O（n）
    frequency(Collection<?> c,Object o); //返回c中与对象o相同的元素**个数**
    boolean disjoint(Collection<?> c,Collection<?> c2); //如果两个集合没有相同元素，返回true
    removeIf(Predicate<? super E> filter); //删除所有匹配的元素
    replaceALl(UnaryOperator<E> op);//对列表中所有元素应用这个操作
```

### 链表的一些骚操作
Java中的链表都是双向链表（有pre的），删除以前需要绕来绕去，Java改进了：
``` java
//添加三个元素，后删除第二个元素操作如下：
List<String> staff = new LinkedList<>();
staff.add("A");
staff.add("B");
staff.add("C");
//重点来了！！！-------
Iterator iter = staff.iterator();
//-------------------
String first = iter.next();
String second = iter.next(); //越过第二个元素啦，迭代器在二、三之间了
iter.remove()； //remove删除迭代器上一个越过的元素
```
* ！！！add()操作也一样。会在迭代器之前位置添加。不过只有ListIterator接口有提供add方法。 无顺序的Set是没有这种的。
* ！！双向链表，有个hasPrevious(),previous().这个时候调用`reomove()`会删除刚刚越过的元素，一样的；
* ! `set(newValue)`操作也一样，会用新的值替代迭代器刚刚越过的元素
* **很有意思**：LinkedList提供了一个访问某个特定元素的`get`方法：listDemo.get(i); 但其实还是遍历，只是如果索引大于size()/2就从链表尾部开始搜索。

如果列表只有少量元素，那么就用ArrayList吧！ 避免使用以**整数**索引表示列表位置的所有方法。

### ArrayList
ArrayList实现了List接口,它是以**数组的方式来**实现的,数组的特性是快速的随机读取值.

循环List列表的三种方式: 
* 使用普通for循环,用get方法获取; 
* 使用Iterator迭代器,使用next方法遍历;
* 使用增强for循环,直接输出!  

> 由此可见第三种方法是最方便，最简洁的!

**与Vector比较**：ArrayList方法不是同步的，线程不安全，而Vector是同步的，线程安全，不过会在同步操作上耗费大量时间。

### LinkedList
1. 采用链表的方式来实现List接口的,它本身有自己特定的方法，如: addFirst(),addLast(),getFirst(),removeFirst()等.
2. 由于是采用链表实现的,因此在进行insert和remove动作时在效率上要比ArrayList要好得多!
3. 适合用来实现Stack(堆栈)与Queue(队列),前者先进后出，后者是先进先出.
4. 

### ArrayList与LinkedList的区别

ArrayList与LinkedList都是List接口的实现类,List接口继承了Collection接口,Collection接口又继承了Iterable接口,因此可以看出List同时拥有了**Collection与Iterable接口的特性**.

### HashSet
* hashCode，自己实现hashCode方法应该与equal方法一起，如果`a.equals(b)`成立，则a与b的散列码必须相同。
* 散列表用**链表数组实现**,找位置：
    1. 计算散列码；
    2. 用散列码与桶的总数取余；
    3. 无其他元素，直接放进去；
    4. 有的话，比较hashCode是否相同？相同，不管：不相同，在该桶的链上插入，桶满了（链长度到达阈值），变为平衡二叉树（后改进为红黑树）
* 通常把桶数设定为预计插入元素的75%~150%；
* 散列表太满了，需要增加桶树，就需要将原来的所有元素从新插入新表；满的判定一般用 **装填因子**（一般0.75），超过75%的位置都已经填入元素，就增加双倍桶数目；

**TreeSet**：树集。变成有序集合。可以按照顺序打印，使用红黑树去实现，插入操作会慢一点。需要实现Comparable接口


## Map接口
也叫映射。是对**键**进行散列、比较。值不能散列或比较。

实现有：HashMap、TreeMap、EnumMap、

**树映射还是散列映射？**

看需求。速度还是顺序性？

add() -> put() ; next() -> get();


### HashMap

**底层实现**：HashMap是基于**哈希表**的M**ap接口**的**非同步**实现，Java最基本数据结构就是两种，一种是数组，一种是引用。**所有的数据结构都可以用这两个基本结构来构造的**，HashMap也不例外。HashMap实际上是一个“**链表散列**”的数据结构，即数组和链表的结合体。**HashMap底层**就是一**个数组结构**，**数组中的每一项又是一个链表**。当新建一个HashMap的时候，就会初始化一个数组。Entry就是数组中的元素，每个 Map.Entry 其实就是一个key-value对，它持有一个指向下一个元素的引用，这就构成了链表。

**存取实现**：
* **put**元素的时候:
1. 先根据**key的hashCode重新计算hash值**，根据hash值得到这个元素在**数组中的位置**（即下标）。
2. 如果数组该位置上已经**存放有**其他元素了，那么在这个位置上的元素将**以链表的形式存放**，**新加入的放在链头，最先加入的放在尾**。
3. 如果数组该位置上没有元素，就直接将该元素放到此数组中的该位置上。

* **get**元素的时候

### TreeMap

### 常用方法

**Map**：
``` java
    V get(Object key); //获取与键对应的值，无返回null
    default V getOrDefault(Object key,V defaultValue) // 无返回默认值defaultValue
    V put(K key,V value); //插入数据，如果键已存在，则覆盖原数据，返回对应的旧值；如果键不存在，返回null。
    boolean containsKey(Object key); //是否已有这个键
    boolean containsValu(Object value); //是否已有这个值
    default void forEach();
```

**HashMap**:
``` java
    HashMap();
    HashMap(int initialCapacity);
    HashMap(int initialCapacity,float loadFactor); //给定容量和装填因子构造一个空散列映射，默认装填因子0.75.
```
**TreeMap**：
``` java
    TreeMap(); //为实现Comparable接口的键构造一个空的树映射
    TreeMap(Comparator<? super K> c);//使用一个指定的比较器对键排序
```
**SortedMap**：
``` java
    Comparator<? super K> comparator(); //返回对键排序的比较器，如果是用Comparable接口的CompareTo方法进行比较的，返回null
    K firstKey();
    K lastKey();//返回映射中的最小元素和最大元素
```

## 队列 Queue

### 一些方法
Queue：
* `boolean add(E element)` 和 `boolean offer(E element)`：尾部添加元素，如果满了，第一个方法是抛出异常，第二个是返回false。
* E remove() ; E poll(); **删除并返回**队头元素，为空的话，第一个抛出异常，第二个返回nul。
* E element(); E peek();**返回但不删除**头部元素，如果为空，第一个抛出异常，第二个返回null。

Deque：
* `void addFirst(E element); void addLast(E element)`和`boolean offerFirst(E element); boolean offerLast(E element)`:前两个添加失败是抛出异常，后两个是返回false。
* 下同。都是加一个First，Last的区别
* 
### 优先队列 PriorityQueue
小顶堆的感觉，用二叉树实现，每次remove都是删除最小的。 典型示例：任务调度，有优先级的。

其add()自动按照构建小顶堆的方式添加元素。

# 进程、多线程、线程池

## 一些常见问题
1. 如果对`Thread`派生子类，就应当**只**覆盖`run()`，而不要覆盖其他方法！ Thread类的其他方法（star(),interrupt(),join(),sleep()等标准方法）都有非常特定的语义。可以根据需要**提供额外的构造函数和其他方法**。
2. 线程sleep和wait的区别？哪个涉及锁的释放？
> wait释放了锁，使得其他线程可以使用同步控制块或者方法。
>wait，notify和notifyAll只能在同步控制方法或者同步控制块里面使用，而sleep可以在任何地方使用（使用范围）
sleep必须捕获异常，而wait，notify和notifyAll不需要捕获异常

## 进程
1. 进程和线程的区别？
> * 进程：进程是程序的一次执行过程；进程是是正在运行程序的抽象；系统资源（如内存、文件）以进程为单位分配；操作系统为每个进程分配了独立的地址空间；操作系统通过“调度”把控制权交给进程。
> * 线程：有标识符ID；有状态及状态转换；不运行时需要保存上下文环境（需要程序计数器等寄存器）；有自己的栈和栈指针
；共享所在进程的地址空间和其它资源。
> * 区别： 
        定义方面：进程是程序在某个数据集合上的一次运行活动；线程是进程中的一个执行路径。（进程可以创建多个线程）
        角色方面：在支持线程机制的系统中，进程是系统资源分配的单位，线程是CPU调度的单位。
        资源共享方面：进程之间不能共享资源，而线程共享所在进程的地址空间和其它资源。同时线程还有自己的栈和栈指针，程序计数器等寄存器。
        独立性方面：进程有自己独立的地址空间，而线程没有，线程必须依赖于进程而存在。
        开销方面。进程切换的开销较大。线程相对较小。（前面也提到过，引入线程也出于了开销的考虑。）
2. 为什么要引入线程？ 
>进程有利于资源的管理和保护,但是： 
>1. 进程切换的代价、开销比较大； 在进程内创建、终止线程比创建、终止进程要快。
>2. 在一个进程内也需要并行执行多个程序，实现不同的功能。 性能也快很多
>3. 进程有时候性能比较低。

## 锁-多线程

### Synchronized 和 ReentrantLock的区别？


### CAS
比较并交换。属于乐观锁技术。

优点：确保对内存的读-改-写操作都是原子操作执行。

缺点：**ABA问题**；循环时间长开销大；**只能保证一个共享变量**的原子操作。

总结：**线程冲突较少**的情况使用。

### AQS  （AbstractQueuedSynchronizer）
是一个用于构建锁和同步容器的框架。 concurrent包许多类都是基于AQS构建的。例如：
> ReentrantLock、Semaphore、FutureTask.

解决了 ：在实现同步容器时设计的大量细节问题。

AQS使用一个FIFO的双向队列表示排队等待锁的线程。队列头节点称为**哨兵节点**。其他节点都维护一个等待状态`waitStatus`。

AQS还有一个表示状态的字段state。

### 什么是原子操作？

i++不是一个原子操作：读，加，写。
解决：jdk1.5后，concurrent.atomic包提供了int和long类型的装类，可以保证操作的原子性，而不需要使用同步。



# IO

## 几种IO比较
**BIO**：同步阻塞IO，阻塞整个步骤，如果连接少，他的延迟是最低的，因为一个线程只处理一个连接，适用于少连接且延迟低的场景，比如说数据库连接。

**NIO**：同步非阻塞IO，阻塞业务处理但不阻塞数据接收，适用于高并发且处理简单的场景，比如聊天软件。

**AIO**：异步IO，他的数据请求和数据处理都是异步的，数据请求一次返回一次，适用于长连接的业务场景。

## NIO
JDK1.4开始，增加了新的io模式**new IO**。 Socket也属于IO的一种，nio为它提供了：
> ServerSocketChannel 和 SocketChannel

**三个重要概念**：`Buffer`(所送货物)、` Channel`（送货车）、` Selector`（分拣员）。
1. **Buffer**
> **四个属性**：mark <= postion <= limit <= capacity
>* capacity:容量，最多可以保存多少元素，创建初设定后无法改变；
>* limit：可以使用的上限，当前有20个元素，就只能操作20个。这个值需要 <= capacity；
>* position:当前所操作元素所在的索引位置，从0开始，随着get和put方法自动更新；
>* mark：暂存postion，可以通过reset方法，将postion恢复到mark位置。
>
> **两个方法**：
> * clear():初始化limit = capacity、position = 0、mark= -1三个属性。
> * flip():保存数据后让position加1，读数据需要将position位置设置为limit。

2. **Channel**
3. **Selector**

**NioSocket**中服务端的处理过程：
1. 创建`ServerSocketChannel`并设置相应参数；
2. 创建`Selector`并**注册**到`ServerSocketChanel`上；
3. 调用`Selector`的`select`方法**等待请求**；
4. `Selector`接受到请求后使用`selectedKeyr`返回`SelectionKey`集合；
5. 使用`SelectionKey`获取到`Channel`、`Selector`和操作类型并进行具体操作。

## Netty

### 优势
1. Netty为什么传输快？ 零拷贝。
2. 为什么说Netty封装好？
    * Channel:表示一个连接
    > ChannelHandler，用于处理业务请求;ChannelHandlerContext，用于传输业务数据;ChannelPipeline，用于保存处理过程需要用到的ChannelHandler和ChannelHandlerContext。
    * ByteBuf:使用方便
    > Heap Buffer 堆缓冲区;Direct Buffer 直接缓冲区;Composite Buffer 复合缓冲区
    * Codec : Netty中的编码/解码器

# 网络

# TCP/IP


# Spring MVC

## 相关术语
* DispatcherServlet
* HandlerAdapter
* HandlerMapping
* ViewResolver
* 

## 常用注解
* Controller  指示Spring类的实例是一个控制器
* RequestMapping  指示一个请求处理方法
* GetMapping
* PostMapping
* RequestParam  参数绑定
* 

## 标签库
* form
* input
* password
* hidden
* 

## ORM和MyBatis
### ORM
对象/关系数据映射

### MyBatis
* SqlSession 类似JDBC中的Connection
* 


# Spring
## 注解那些事

* 
1. @Resource默认按照名称方式进行bean匹配   J2EE的注解
2. @Autowired默认按照类型方式进行bean匹配   Spring的注解
