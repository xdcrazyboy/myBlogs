 看个视频教程总结的，参考了课程笔记。

## 1. 什么是maven
Maven是基于POM（工程对象模型），通过一小段描述来对项目的代码、报告、文件进管理的工具。

Maven是一个跨平台的项目管理工具，它是使用java开发的，它要依赖于jdk1.6及以上.

Maven主要有两大功能：管理依赖、项目构建。依赖指的就是jar包。

 **核心概念**
1. 	坐标
2.	依赖管理
3.	生命周期
4.	插件
5.	继承
6.	聚合

## 2. 与其他工具的比较


### Eclipse
使用eclipse进行项目**构建**，相对来说，步骤**比较零散，不好操作**

###	Ant
它是一个专门的项目**构建**工具，它可以通过一些配置来完成项目构建，这些配置要明确的告诉ant，源码包在哪？目标class文件应该存放在哪？资源文件应该在哪.**麻烦**

###	Maven
它是一个项目**管理**工具，他也是一个项目构建工具，通过使用maven，可以对项目进行快速简单的构建，它不需要告诉maven很多信息，但是需要安装maven去的规范去进行代码的开发。也就是说maven是有约束的。项目构建：编译、打包、部署.


## 3. 安装Maven
1. 官网下载：http://maven.apache.org
2. 直接解压，然后配置path路径。
3. cmd-》执行 mvn -v。显示正常则安装成功。

### 配置

**全局配置**：在maven安装目录的conf里面有一个settings.xml文件，这个文件就是maven的全局配置文件。默认在系统的用户目录下的m2/repository中，该目录是本地仓库的目录。

**用户配置**：用户配置文件的地址：~/.m2/settings.xml，该文件默认是没有，需要将全局配置文件拷贝一份到该目录下。 建议：改成自己的项目专有仓库。

## 4. 创建Maven工程 
Project: Maven工程结构
* src
    * |-main
       > -java        —— 存放项目的.java文件  
       > -resources   —— 存放项目资源文件，如spring, hibernate配置文件
    * -test
       > -java        ——存放所有测试.java文件，如JUnit测试类  
       > -resources   —— 测试资源文件
* target              —— 目标文件输出位置例如.class、.jar、.war文件
* pom.xml             ——maven项目核心配置文件

手动创建maven工程结构个个目录，Hello.java  TestHello.java; pom.xml文件自动构建，也可以复制。

### Maven命令的使用

Maven的命令要在pom.xml所在目录中去执行

* mvn compile  -- 编译
* mvn clean  -- 清除命令，清除已经编译好的class文件，具体说清除的是target整个目录
* mvn test -- 测试命令，该命令会将test目录中的源码进行编译
* mvn package -- 打包
* mvn install --安装命令，会将打好的包，安装到本地仓库

还可以有组合命令：

* mvn clean compile
* mvn clean install
* ....


## 5. M2Eclipse （插件，如果用idea可以忽略）

如何查看是否安装成功？ -> window-> preference->maven-> Install ->改成本地的maven

可以修改用户配置文件的默认位置。 maven-> User Setting

## 6. Maven的核心概念
### 坐标 

``` java
<dependencies>
	<dependency>
		<groupId>junit</groupId>
		<artifactId>junit</artifactId>
		<version>4.10</version>
		<scope>test</scope>
	</dependency>	
</dependencies>

```

> 目的： 为了定位一个唯一确定的jar包。
> 主要组成： 
> * groupId  定义当前Maven组织名称
> * artifactId  定义实际项目名称
> * version  定义当前项目的当前版本

### 依赖管理
**依赖范围**： scope 用于控制依赖和编译、测试、运行的classpath的关系。主要有三种：
* compile: 默认。 编译、测试、运行时，classpath都有效。
* test: 测试有效
* provide: 编译、测试。 
* runtime: 运行时提供。 例如： jdbc 驱动

**依赖传递**

*范围区别* ：
* 当第二直接依赖范围是compile时，可正常传递； 
* 第二是test时，不传递； 
* ...provided时，只传递为第一直接依赖范围为provided的依赖，且传递性范围同样为provided； 
* ...runtime，传递范围与第一直接依赖一致，compile例外，此时传递范围为runtime。 

**依赖冲突**
* 跨pom文件的冲突：  就近原则-依赖最直接最近的。
* 同一个pom文件的冲突：  就近原则-越下面的越近。

**可选依赖**

Optional标签表示：该依赖是否可选。
``` java
<optinal>true</optional>    //默认为false：该依赖会传递下去；  true： 该依赖不会传递下去
```

**排除依赖**
Exclusions标签：可以排除依赖。 在引入first.jar 时，可添加此标签，标注 排除掉first.jar传递过来的某个依赖。

### 生命周期

项目构建的步骤集合： clean生命周期，default生命周期，site生命周期。

生命周期由多个阶段（Phase）组成。

**Clean 生命周期**

 * pre-clean  -  执行一些需要在clean之前完成的工作
 * clean  -  移出所有上一次构建生成的文件
 * post-clean  -  执行一些需要在clean之后立刻完成的工作

mvn post-clean 该命令会自动执行对应什么周期 前面的命令 pre-clean clean .简化命令输入

**Default 生命周期（重点）**

* validate 
* generate-sources 
* process-sources 
* generate-resources 
* process-resources 复制并处理资源文件，至目标目录，准备打包。 
* **compile** 编译项目的源代码。 
* process-classes 
* generate-test-sources 
* process-test-sources 
* generate-test-resources 
* process-test-resources 复制并处理资源文件，至目标测试目录。
* test-compile 编译测试源代码。 
* process-test-classes 
* **test** 使用合适的单元测试框架运行测试。这些测试代码不会被打包或部署。
* prepare-package 
* **package** 接受编译好的代码，打包成可发布的格式，如 JAR 。 
* pre-integration-test 
* integration-test 
* post-integration-test 
* verify 
* **install** 将包安装至本地仓库，以让其它项目依赖。 
* **deploy** 将最终的包复制到远程的仓库，以让其它开发人员与项目共享。

**运行任何一个阶段的时候，它前面的所有阶段都会被运行**，这也就是为什么我们运行mvn install 的时候，代码会被编译，测试，打包。此外，Maven的插件机制是完全依赖Maven的生命周期的，因此理解生命周期至关重要。

**Site生命周期**-生成项目站点
* pre-site  -  执行一些需要在生成站点文档之前完成的工作 
* site  -  生成项目的站点文档 
* post-site  -  执行一些需要在生成站点文档之后完成的工作，并且为部署做准备 
* site-deploy  -  将生成的站点文档部署到特定的服务器上 

用以生成和发布Maven站点，文档及统计数据自动生成

### 插件
Maven的核心仅仅定义了抽象的生命周期，具体的任务都是交由插件完成的。每个插件都能实现一个功能，每个功能就是一个插件目标。Maven的生命周期与插件目标相互绑定，以完成某个具体的构建任务。
>例如compile就是插件maven-compiler-plugin的一个插件目标

**编译插件**
``` java
<build>
    <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-compiler-plugin</artifactId>
        <configuration>
            <source>1.7</source>
            <target>1.7</target>
            <encoding>UTF-8</encoding>
        </configuration>
    </plugin>
</build>
```
修改配置文件后，在工程上点击右键选择maven→update project configration

**Tomcat插件**
如果使用maven的tomcat插件的话，那么本地则不需要安装tomcat。

1. 创建一个maven工程
2. 选择打包方式为 war 才能创建web工程
3. web-app下创建WEB-INF（及它里面的web.xml）和index.jsp

使用tomcat插件运行web工程
1. run as -> build -> tomcat:run 
2. 可以换默认的tomcat版本，通过pom.xml添加build-plugin。
``` java
<plugin>
	<!-- 配置插件 -->
	<groupId>org.apache.tomcat.maven</groupId>
	<artifactId>tomcat7-maven-plugin</artifactId>
	<configuration>
		<port>8080</port>
		<path>/</path>
	</configuration>
</plugin>
```

### 继承

指pom文件的继承
1. 创建父工程：不用模板，选择打包方式必须为pom方式。
2. 创建子工程：创新工程、修改老工程。 选择打包方式界面-下方-Parent Project处添加父工程。

结果就是：
> 子工程pom文件多了\<parent>\</parent>标签。（修改老工程其实就是直接修改其pom文件：增加此标签，指定老工程的父工程）

作用：
>父工程统一依赖jar包：在父工程中对jar包进行依赖，在子工程中都会继承此依赖。

**父工程统一管理版本号**

jar包依赖太多了，每个要写一个dependency很冗余，于是：
1. Maven使用dependencyManagement管理依赖的版本号。
注意：此处只是定义依赖jar包的版本号，并不实际依赖。如果子工程中需要依赖jar包还需要添加dependency节点。
> 父工程：\<dependencyManagement> \<dependencyManagement/>
> 子工程：需要依赖的话自己定义\<dependency>，但是直接使用父类的版本号，自己不用指定version。

2. 父工程抽取版本号
父工程统一管理依赖后，会有很多，所以对版本号抽象出来，单独定义，方便修改管理。
``` java
<properties>
	<junit.version>4.9</junit.version>
</properties>

<dependencyManagement>
	<dependencies>
		<dependency>
			<groupId>junit</groupId>
			<artifactId>junit</artifactId>
			<version>${junit.version}</version>
		</dependency>
	</dependencies>
</dependencyManagement>
```

### 聚合
聚合工程的多个模块。例如：电商项目中，包括商品模块、订单模块、用户模块等。就可以对不同的模块单独创建工程，最终在打包时，将不同的模块聚合到一起。

例如同一个项目中的表现层、业务层、持久层，也可以分层创建不同的工程，最后打包运行时，再聚合到一起。
1. 创建一个聚合工程，聚合工程的**打包方式也必须为pom**。 
2. 创建持久层，在聚合工程上右键，创建Maven Module，然后就取名字，其他比如打包方式默认jar。会自动把聚合工程当做父类。
3. 创建业务层。同持久层一样，jar。
4. 创建表现层，**打包方式选择war**.
5. 然后修改父工程的pom文件，对父工程进行 tomcat7:run 运行。

## 7. Maven仓库管理
**Maven仓库**：用来统一存储所有Maven共享构建的位置就是仓库。根据Maven坐标定义每个构建在仓库中唯一存储路径大致为：groupId/artifactId/version/artifactId-version.packaging。

**仓库分类**
1. 本地仓库
>默认在~/.m2/repository. 若用户自己配置了，就以配置的地址为准
2. 远程仓库
  + [中央仓库](http://repo1.maven.org/maven2)（不包含有版权的jar包）
  + 私服:是一种特殊的远程仓库，它是架设在局域网内的仓库.

**Maven私服**
1. **安装Nexus**

> 为所有来自中央仓库的构建安装提供本地缓存。
下载网站：http://nexus.sonatype.org/ （war包最近的似乎没了，自己百度下）
>+ 第一步：将下载的nexus的war包复制到tomcat下的webapps目录。
>+ 第二步：启动tomcat。nexus将在c盘创建sonatype-work目录【C:\Users\当前用户\sonatype-work\nexus】。


2. **Nexus的本地目录**
目录不复杂，暂略。

3. **访问Nexus**
>访问URL: http://localhost:8080/nexus-版本号/
默认账号:
用户名： admin
密码： admin123



**仓库有四种类型**
*	group(仓库组)：一组仓库的集合
*	hosted(宿主)：配置第三方仓库 （包括公司内部私服 ） 
*	proxy(代理)：私服会对中央仓库进行代理，用户连接私服，私服自动去中央仓库下载jar包或者插件 
*	virtual(虚拟)：兼容Maven1 版本的jar或者插件

**Nexus的仓库和仓库组介绍**
*	3rd party: 一个策略为Release的宿主类型仓库，用来部署无法从公共仓库获得的第三方发布版本构建
*	Apache Snapshots: 一个策略为Snapshot的代理仓库，用来代理Apache Maven仓库的快照版本构建
*	Central: 代理Maven中央仓库
*	Central M1 shadow: 代理Maven1 版本 中央仓库
*	Codehaus Snapshots: 一个策略为Snapshot的代理仓库，用来代理Codehaus Maven仓库的快照版本构件
*	Releases: 一个策略为Release的宿主类型仓库，用来部署组织内部的发布版本构件
*	Snapshots: 一个策略为Snapshot的宿主类型仓库，用来部署组织内部的快照版本构件
*	Public Repositories:该仓库组将上述所有策略为Release的仓库聚合并通过一致的地址提供服务

5. ** 配置所有构建均从私服下载**

在本地仓库的setting.xml中配置如下：
``` java
<mirrors>
	 <mirror>
		 <!--此处配置所有的构建均从私有仓库中下载 *代表所有，也可以写central -->
		 <id>nexus</id>
		 <mirrorOf>*</mirrorOf>
		 <url>http://localhost:8080/nexus-2.7.0-06/content/groups/public/</url>
	 </mirror>
 </mirrors>
```

6. **部署构建到Nexus**
- 第一步：Nexus的访问权限控制.在本地仓库的setting.xml中配置如下：
     
	 ``` java
    <server>
 		<id>releases</id>
		<username>admin</username>
		<password>admin123</password>
	</server>
	<server>
		<id>snapshots</id>
		<username>admin</username>
		<password>admin123</password>
	</server>
     ```
- 第二步：配置pom文件.在需要构建的项目中修改pom文件:
``` java
<distributionManagement>
        <!-- 正式版本发布位置地址 -->
		<repository>
			<id>releases</id>
			<name>Internal Releases</name>
			<url>http://localhost:8080/nexus-2.7.0-06/content/repositories/releases/</url>
		</repository>
        <!-- 快照版本发布位置地址 -->
		<snapshotRepository>
			<id>snapshots</id>
			<name>Internal Snapshots</name>
			<url>http://localhost:8080/nexus-2.7.0-06/content/repositories/snapshots/</url>
		</snapshotRepository>
	</distributionManagement>
```
- 第三步：执行maven的deploy命令. mvn deploy
