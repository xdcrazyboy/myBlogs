# 参考资料
《Tomcat内核设计剖析》
《Tomcat架构解析》 刘光瑞

# Tomcat总体架构

## 总体设计

1.**静态设计**

### Server

一个`Server`—<包含>—>多个`Service`;Server表示整个Servlet容器。

每个`Service`—负责维护—多个`Connector`和**一个**`Container`（`Engine`）
    
- Connector:负责开启Socket并监听客户端请求、返回响应数据
- Container：负责具体的请求处理。下能够执行客户端请求并返回响应的一类对象

### Container
不同级别的容器：

- Engine：由Container—重命名，只负责请求处理，不管请求链接、协议。表示整个Servlet引擎。
- Context：一个Engine容器能够找到合适多个Web应用（用Context表示一个应用）;
    > 拥有start()/stop()方法。用于表示ServletContext。
- Host：将每个域名抽象为一个虚拟机，每个Host可以包含多个Context。而一个Engine包含多个Host。
- Wrapper：Web应用中定义的Servlet。

**Lifecycle**：针对所有拥有生命周期的管理特性的组件，抽象了一个Lifecycle通用接口，该接口定义了生命周期管理的核心方法。
    
- Init():初始化组件。
- start():启动组件。
- stop():停止组件。
- destroy():销毁组件。
- LifecycleListener：支持添加事件监听器，支持组件状态以及状态之间的转换。

**Pipeline和Value**：前者用于构造职责链，后者代表职责链上的每个处理器。

Tomcat每个层级的容器均通过PipeLine和Value进行请求处理。

### Connector
功能：
- 监听服务器端口，读取来自客户端的请求
- 将**请求**数据按照指定协议进行**解析**
- 根据你请求地址**匹配正确的容器**进行处理
- 将响应返回客户端
  
包含：
- `ProtocolHandler`
    - `Processor`：按照指定协议读取数据。然后按照请求地址映射到具体的容器进行处理，也就是**请求映射**。
    - `AbstractEndpoint` ：启动Socket监听

由`Service`维护：
   > **请求映射**除了考虑映射规则的实现外，还要考虑容器组件的注册和销毁。
    > 用`Mapper`和`MapperListener`实现上诉功能.
        > - Mapper:用于维护容器映射信息，同时按照映射规则查找容器
        > - MapperListener:实现了ContainerListener和LifecycleListener，用于容器组件状态发生变化时，注册、取消对应的容器映射信息。

### Executor 解决并发的问题。
Tomcat中的Executor由Service维护，因此**同一个Service中的组件可以共享一个线程池**。

Endpoint会启动一组线程来监听Socket端口，当接受到客户端请求后，会创建请求处理对象，并交给线程池处理，由此支持并发处理客户端请求。

### Bootstrap 和 Catania

* Bootstrap:作为应用服务器启动入口。负责创建`Catalina`实例，根据执行参数调用Catalina相关方法完成针对应用服务器的操作（启动、停止）。
* Catalina:通过这个类提供一个Shell程序，用于解析`server.xml`创建各个组件。
* Digester：解析XML文件。

**为什么不直接通过Catania启动，而是又提供了Bootstrap**

1. Bootstrap与Tomcat应用服务器完全松耦合，它可以直接依赖JRE运行并为Tomcat应用服务器创建共享类加载器，用于构造Catania实例以及整个Tomcat服务器。
2. 启动入口和核心环境的解耦，灵活的组织中间件产品的结构，尤其是类加载器的方案。


2.动态设计

### Tomcat启动
- 统一按照生命周期管理接口Lifecycle的定义进行启动。
- 首先，调用init()方法进行组件的逐级初始化，然后调用start()方法进行启动。
- Bootstrap——》Catania——》Server——》Service——》Executor、Engine（Host->Context）、Connector（ProtocolHander）

### 请求处理
- Endpoint——》Processor——》CoyoteAdapter

### 类加载器

1. JVM默认的3个类加载器：
   - **Bootstrap**Class Loader
    >加载JVM提供的基础运行类，jre/lib目录下的核心类库。
   - **Extension** Class Loader
    >JVM会自动加载，放到jre/lib/ext目录下，不推荐将应用程序依赖的类库放置到扩展目录下，该目录对所有基于该JVM运行的应用程序可见。
   - **System** Class Loader
    >用于加载环境变量CLASSPATH指定目录下的或者-classpath运行参数指定的jar包。一般用于加载应用程序Jar包及其启动入口类（Tomcat的Bootstrap类即由System类加载器加载。）

2. Tomcat加载器
   
   应用服务器一般会自行创建类加载器以实现更灵活的控制，一方面是对规范的实现（Servlet规范要求每个Web应用都有一个独立的类加载器实例），另一方面也有架构层面的考虑：

   - **隔离性**：Web应用类库相互隔离，避免依赖库或者应用包相互影响。
   - **灵活性**：可以只针对一个Web应用进行重新部署，不影响其他应用，更加灵活。
   - **性能**：只搜索自己的Jar包，效率会高一些。

    Tomcat的类加载方案：
    - Common Class Loader：以System为父加载器，路径为common.loader
        - Catalina Class Loader:用于加载Tomcat应用服务器的类加载器，路径为server.loader，默认为空，用父类加载应用服务器。
        - Shared Class loader:是所有Web应用的父加载器，默认为空。
            - Web Appl Class loader：以Shared为附加在其，加载`/WEB-INF/classers`目录下的未压缩的Class和资源文件，以及`/WEB-INF/lib`目录下的Jar包。只对当前web应用可见。

    架构分析补充：
    - **共享**：
        - Common类加载器：实现了 **Jar包在应用服务器和Web应用之间的共享**。
        - Shared类加载器：实现了 **Jar包在Web应用之间的共享**。
    - **隔离性**：
        - 服务器与Web应用的隔离，理论上，除了Servlet规范定义的接口外，Web应用不应该依赖服务器任何的实现类（通过 **JVM的安全策略许可** 实现），这样才有助于Web应用的可移植性。 正因如此，Tomcat通过Catalina类加载器加载服务器依赖的包，以便应用服务器与Web应用更好地隔离。


    Web应用类加载器
    - **Java默认的类加载机制**的委派过程：
        1. 从缓存中加载；
        2. 没有？从父类加载器中加载；（双亲委派0
        3. 没有？从当前类加载器加载；
        4. 如果还是没有，抛出异常。
    - Tomcat中Web应用类加载器过程：
        1. 从缓存中加载
        2. 没有？从JVM的Bootstrap类加载器加载
        3. 没有？从当前类加载器加载（先`WEB-INF/classes`,再`WEB-INF/lib`）；
        4. 没有？从父加载器加载，而父类加载器使用默认的委派模式（System、Common、Shared）
    - Tomcat提供`delegate`属性用于控制是否启用Java委派模式，默认为`false`。  
