## 说明
自己正在找工作吧，看到挺多面经，感觉通过面经不是为了投机取巧去应付面试，而是知道重点，知道哪些知识是需要掌握的（我相信面试官不会为了面试而面试，出些过时的题目）。 

想着把大家的面经集合起来，原本想用爬虫，奈何不熟练，就先自己慢慢总结吧。

这只是告诉大家什么很重要，具体的肯定需要大家去学习，去琢磨，只有打通了它们的联系，才能融会贯通。

# 面向对象的一些基础

## 六大原则
共性：高内聚低耦合

1. 单一职责：一个类只负责一个功能领域中的相应职责
2. 接口隔离：使用多个专门的接口，而不是单一的总接口。
3. 开闭原则：对扩展开放，对修改关闭。 即软件实体尽量在不修改原有代码的情况下进行扩展。
4. 里氏替换原则：所有引用父类的地方必须能透明地引用其之类对象（反之不然）。
5. 依赖倒置原则：抽象不应该依赖于细节，细节应该依赖于抽象。 要针对接口编程，而不是针对实现编程。
6. 迪米特法则：一个软件实体应当 **尽可能少地与其他实体发生相互作用**。

## 接口和抽象类

## 多态

# 集合

## 集合接口和实现类直接的关系

## Hashmap

## ConcurrentHashMap

## List

## Set



# Java虚拟机

## 内存模型
1. Java的内存分区
2. 

## 垃圾回收GC
1. 回收方式、回收算法？
2. CMS和G1了解么？ CMS解决什么问题，说一下回收过程？CMS回收停顿几次？ 为什么停顿两次？


## 栈溢出、内存泄漏
1. Java栈为什么会发生内存溢出，Java堆呢？ 举个场景？（集合类持有对象——>那么集合类是怎么解决这个问题的？ 软引用和弱引用？ 那这些引用的区别？虚引用）
2. 

## 系统调优

# 多线程、高并发

## 进程和线程

## 线程池


# 锁

## 数据库中的锁

## 并发中的锁

## volatile 和 Synchronized 和 lock
1. 它们的使用方式和实现原理有什么区别？
   > syn用于方法和代码块，锁对象和类以及方法； lock一般锁代码块，可以和condition搭配使用；
   > 实现原理;syn使用底层的mutex锁，需要系统调用，而Lock则使用AQS实现。 
2. Synchronized锁升级的过程，偏向锁，轻量级锁，重量级锁，它们分别怎么实现，解决哪些问题？ 什么时候发生锁升级。

### Lock

## CAS

## AQS

## cycle那一堆



# 数据库（通用的、MySQL）

## SQL

## 优化

## 索引
### 聚集索引

1. 什么是聚集索引？
2. 聚集索引的作用是什么？（经常用于范围查询）
3. 为什么一个表只能有一个聚集索引？（代表了物理排列顺序，类似电话本，该索引可以有多个列-组合索引）

## 数据库连接池


# IO

## BIO、NIO、AIO

## 阻塞 vs 非阻塞

## Selector 和 poll 、epoll


# 网络

## TCP

## IP

IPv4的头部：
- 版本号
- 头部长度
- 服务类型（ToS）
- 总长度
- ID号
- 标志Flag
- 分片偏移量
- 生存时间 TTL 跳数
- 协议号 （TCP、UDP、ICMP）
- 头部校验和
- **发送方IP地址**
- **接收方IP地址**
- 可选字段

## HTTP

## 网络编程

## Socket


# 设计模式



## 单例模式

## 代理模式

## 工厂模式（简单工厂、抽象工厂）

## 


# Spring

## IOC

### IOC容器

### Bean初始化？ 生命周期？
Spring IOC容器对Bean的生命周期进行管理的过程如下：
1. 通过构造器或者工厂方法创建Bean实例
2. 为Bean的属性设置值和对其它Bean的引用
3. 将Bean的实例传递给bean后置处理器的`postProcessBeforeInitialization`方法
4. 调用Bean的初始化方法
5. 将Bean的实例传递给bean后置处理器的`postProcessAfterInitialization`方法
6. bean可以用了
7. 当容器关闭，调用Bean的销毁方法

Bean的作用域
1. singleton
2. prototype：容器在接到该类型对象的请求的时候，会每次都重新生成一个新的对象实例给请求方。
3. request、session和global session ：只适用于web应用程序

## AOP

### 动态代理？
### Cglib VS Java虚拟机实现

## 事务
1. 事务的ACID属性 
2. 事务的隔离级别
3. 事务的传播机制
4. 声明式事务 VS 编程式事务

### 声明式事务

# Spring MVC

# Struts2

# MyBatis 、Hibernate

# Tomcat

## Tomcat的类加载器结构。