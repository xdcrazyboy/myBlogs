
# 来源书籍

[《Docker — 从入门到实践》](https://github.com/yeasy/docker_practice) ，大家看详细版本，可以直接去这里看。
- 《Docker In Practise》
- 《Docker IN Action》


# Docker是什么
## 定义

## Docker的好处
  1. 替代虚拟机
   >只关心应用程序的话。 启动速度比虚拟机快，迁移更轻量，得益于它的分层文件系统，与其他人分享变更也更简单便捷。非常适合脚本化。
  2. 软件原型
   >快速体验软件，而又避免干扰目前的配置或配备一个虚拟机的麻烦。
   1. 打包软件
   >对于Linux用户而言，Docker镜像实际上没有依赖，非常适用于打包软件。
   1. 让微服务架构成为可能
   >有助于将一个复杂系统分解成一系列可组合的部分，让用户可以用更加离散的方式划分服务，更加易于管理和可插拔。
   1. 网络建模
   >因为一台机器上启动数百个隔离的容器，可以很好的建网络模型做测试。
   1. 离线时启用全栈生产力
   2. 降低调试支出
   3. 文档化软件依赖及接触点
   4.  启动持续交付
   >**更具有可重现性和可复制性**。

## Docker 镜像
Docker 镜像是一个特殊的文件系统，除了提供容器运行时所需的程序、库、资源、配置等文件外，还包含了一些为运行时准备的一些配置参数（如匿名卷、环境变量、用户等）。镜像不包含任何动态数据，其内容在构建之后也不会被改变。

镜像不是由一个文件组成，而是一组文件系统组成，或者说由多层文件系统联合组成。
> 层层构建，前一层是后一层的基础。每一层尽量只包括该层需要添加的东西，任何额外的东西应该在该层构建结束前清理掉。
> 
> 分层
存储优势：使得镜像的服用、定制变得更加容易，可以在原来基础上，添加新层，定制新内容。


## Docker容器
镜像（Image）之余容器（Container） = 类 之余 实例

镜像是静态的定义，容器是镜像运行的**实体**，容器可以被创建、启动、停止、删除、暂停。

容器的实质是**进程**，但它运行于属于自己的独立的命名空间。 容器拥有自己的root文件系统、自己的网络配置、自己的进程空间、甚至自己的用户ID空间。容器是独立运行的一个或一组应用，以及它们的运行态环境。

因为容器有这种隔离的特性，很多初学者把它跟虚拟机混淆。

容器也是分层的，以镜像为基础层，上面创建当前容器的存储层。

>容器存储层生存周期跟容器一样，容器不应该向其存储层写入任何数据，容器存储要保持无状态化，所有文件的写入操作，都因该使用**数据卷**、或者绑定宿主的目录。 数据卷的生存周期独立于容器，因此使用它后，容器删除，数据不会丢失。

## Docker仓库
Docker Registry：提供集中的存储、分发镜像的服务。

一个Docker Registry可以包含多个仓库（Repository）；每个仓库可以包含多个标签（Tag)；每个标签对应一个镜像。
> <仓库名>:<标签>   ——> ubuntu:16.04 (不写标签，则默认latest)

Docker Registry公开服务

1. 官方的：[Docker Hub](https://hub.docker.com/)
2. 其他：[CoreOS的Quay.io](https://quay.io/repository/)
3. 谷歌的(Kubernetes用的就是这个):[Google Container Registry](https://cloud.google.com/container-registry/)
4. 国内的：[阿里云镜像库](https://cr.console.aliyun.com/)、[时速云镜像仓库](https://hub.tenxcloud.com/)、[网易云镜像服务](https://c.163.com/hub#/m/library/)
5. *加速器*：[阿里加速器](https://cr.console.aliyun.com/#/accelerator)、[DaoCloud加速器](https://www.daocloud.io/mirror#accelerator-doc)
   

私有Docker Registry：官方提供了这样的镜像，你可以直接使用，但是免费版功能有限。



# 安装 Docker

## 版本
CE（社区免费版）和EE（企业版）。我们安CE的stable版本（六个月发布一个）。

官网有各种环境下的[安装指南](https://docs.docker.com/engine/installation/).我只说Ubuntu下的安装方式。

## Ubuntu

支持版本：14.04、16.04、17.10、18.04。
建议使用16.04或者18.04.这里用16.04

1. 卸载旧版本。不管有没有，直接运行：
``` shell
$ sudo apt-get remove docker \
               docker-engine \ 
               docker.io
```
2. 使用Apt安装

``` shell
$ sudo apt-get update
//添加CA证书
$ sudo aptsudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
```
这里使用国内源。

``` shell
//为了确认所下载软件包的合法性，需要添加软件源的 GPG 密钥。
$ curl -fsSL https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu/gpg | sudo apt-key add -

//向 source.list 中添加 Docker 软件源(stable版本的)
$ sudo add-apt-repository \
    "deb [arch=amd64] https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu \
    $(lsb_release -cs) \
    stable"
```

3. 安装Docker CE

``` shell
$ sudo apt-get update

$ sudo apt-get install docker-ce
```

4. 启动Docker CE
```shel
$ sudo systemctl enable docker
$ sudo systemctl start docker
```   

5. 建立docker用户组
   
   出于安全考虑一般 Linux 系统上不会直接使用 root 用户。因此，更好地做法是将需要使用 docker 的用户加入 docker 用户组。

```shell
//建立docker组：
$ sudo groupadd docker

//将当前用户加入docker组
$ sudo usermod -aG docker $USER
```

注意：**退出当前终端并重新登录**，进行如下测试。
6. 测试Docker是否安装正确

``` shell
//若能正常输出以下信息，则说明安装成功。
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
ca4f61b1923c: Pull complete
Digest: sha256:be0cd392e45be79ffeffa6b05338b98ebb16c87b255f48e297ec7f98e123905c
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://cloud.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/engine/userguide/
```


## 删除镜像
1. 如果有实例的容器，无法删除该镜像。
2. 删除行为分为两类，一类是 Untagged ，另一类是 Deleted 。
3. 我们之前介绍过，镜像的唯一标识是其ID 和摘要，而一个镜像可以有多个标签。
4. 因此当我们使用上面命令删除镜像的时候，实际上是在要求删除某个标签的镜像，只有该镜像下的标签全部Untagged后才会出发Deleted操作。

## Docker容器的特点和优势
1. 容器中仅运行了指定的 bash 应用。这种特点使得 Docker 对资源的利用率极高，是货真价实的轻量级虚拟化。


## 容器创建过程
当利用 docker run 来创建容器时，Docker 在后台运行的标准操作包括：
1. 检查本地是否存在指定的镜像，不存在就从公有仓库下载启动
2. 利用镜像创建并启动一个容器
3. 分配一个文件系统，并在只读的镜像层外面挂载一层可读写层
4. 从宿主主机配置的网桥接口中桥接一个虚拟接口到容器中去
5. 从地址池配置一个 ip 地址给容器
6. 执行用户指定的应用程序
7. 执行完毕后容器被终止

## 进入容器操作
某些时候需要进入容器进行操作，包括使用 `docker attach` 命令或 `docker exec` 命令，推荐大家使用 docker exec 命令。

因为：当通过exec进入容器交互界面，exit后不会停止容器，前者会。


## 使用仓库，寻找好的镜像
1. docker login //登录账号
2. docker search centos //搜索镜像，以centos为例，--filter=stars=N 找收藏数达到多少以上的。
3. docker logout //登出账号
4. docker pull centos //下载容器
5. docker push username/ubuntu:18.04 //推送自己的容器到仓库

## 自己构建镜像
### 利用 commit 理解镜像构成
>不要使用 docker commit 定制镜像，定制镜像应该使用 Dockerfile 来完成。
现在让我们以定制一个 Web 服务器为例子，来讲解镜像是如何构建的。
1. `>docker run --name webserver -d -p 80:80 nginx` 
   > 这条命令会用 nginx 镜像启动一个容器，命名为 webserver ，并且映射了 80 端口，这样我们可以用浏览器去访问这个 nginx 服务器。在浏览器输入`localhost`登录界面，可以看到Nginx的欢迎界面。
2. 我们修改下默认欢迎界面：
   ``` shell
   $ docker exec -it webserver bash
   root@3729b97e8226:/# echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
   root@3729b97e8226:/# exit
      exit
   ```
3. 刷新下浏览器，发现界面变了。我们可以用下面的命令将容器保存为镜像：
   > docker commit [选项] <容器ID或容器名> [<仓库名>[:<标签>]]
   ``` shell
   $ docker commit \
      --author "Tao Wang <twang2218@gmail.com>" \
      --message "修改了默认网页" \
      webserver \
      nginx:v2
   ```
4. 最后通过`docker image ls nginx`就可以看到新定制的nginx的镜像。
   
> 这种方法会修改很多层的东西，积累着不会删除，让镜像很臃肿。

### 使用 Dockerfile 定制镜像
>从刚才的 docker commit 的学习中，我们了解到，镜像的定制实际上就是定制每一层所添加的配置、文件。如果我们可以把每一层修改、安装、构建、操作的命令都写入一个脚本，用这个脚本来构建、定制镜像，那么之前提及的无法重复的问题、镜像构建透明性的问题、体积的问题就都会解决。这个脚本就是 Dockerfile。

**Dockerfile**:一个文本文件，其内包含了一条条的 指令(Instruction)，每一条指令构建一层，因此每一条指令的内容，就是描述该层应当如何构建。文件内容：
```docker
FROM nginx
RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
``` 
其中，**FROM** 指定基础镜像。在 [Docker Hub](https://hub.docker.com/search?q=&type=image&image_filter=official) 上有非常多的高质量的官方镜像，有可以直接拿来使用的服务类的镜像，如 nginx、redis、mongo、mysql、httpd、php、tomcat 等；也有一些方便开发、构建、运行各种语言应用的镜像，如 node、openjdk、python、ruby、golang 等。可以在其中寻找一个最符合我们最终目标的镜像为基础镜像进行定制。

`FROM scratch`不使用任何镜像，直接从第一层开始创建，比较少用。

 **RUN执行命令**，由于命令行的强大能力，RUN 指令在定制镜像时是最常用的指令之一。其格式有两种：
-  _shell_ 格式：RUN <命令>，就像直接在命令行中输入的命令一样。刚才写的 Dockerfile 中的 RUN 指令就是这种格式。每个RUN会创建一层，所以如果是多条执行语句，可以使用`buildDeps='gcc libc6-dev make wget' && apt-get update && ...`
   ```docker
   RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
   ```
- _exec_ 格式：RUN ["可执行文件", "参数1", "参数2"]，这更像是函数调用中的格式。
   ```docker

   ```

**构建镜像**（以nginx为例子）
1. 在Dockerfile所在目录运行以下命令
   ```
   > docker build -t nginx:v3 .
   ```
   针对上个命令最后的 . 的一些说明： 那是上下文路径而不是Dockerfile的路径。

   - 理解构建上下文对于镜像构建是很重要的，避免犯一些不应该的错误。比如有些初学者在发现 COPY /opt/xxxx /app 不工作后，于是干脆将 Dockerfile 放到了硬盘根目录去构建，结果发现 docker build 执行后，在发送一个几十 GB 的东西，极为缓慢而且很容易构建失败。那是因为这种做法是在让 docker build 打包整个硬盘，这显然是使用错误。

   - 一般来说，应该会将 Dockerfile 置于一个空目录下，或者项目根目录下。如果该目录下没有所需文件，那么应该把所需文件复制一份过来。如果目录下有些东西确实不希望构建时传给 Docker 引擎，那么可以用 .gitignore 一样的语法写一个 .dockerignore，该文件是用于剔除不需要作为上下文传递给 Docker 引擎的。

   - 那么为什么会有人误以为 . 是指定 Dockerfile 所在目录呢？这是因为在默认情况下，如果不额外指定 Dockerfile 的话，会将上下文目录下的名为 Dockerfile 的文件作为 Dockerfile。

   - 这只是默认行为，实际上 Dockerfile 的文件名并不要求必须为 Dockerfile，而且并不要求必须位于上下文目录中，比如可以用 -f ../Dockerfile.php 参数指定某个文件作为 Dockerfile。

2. docker build 还支持从 URL 构建，比如可以直接从 Git repo 中构建：（一般默认找Dockerfile文件）
   ```
   $ docker build https://github.com/twang2218/gitlab-ce-zh.git#:11.1
   ```

   这行命令指定了构建所需的 Git repo，并且指定默认的 master 分支，构建目录为 /11.1/，然后 Docker 就会自己去 git clone 这个项目、切换到指定分支、并进入到指定目录后开始构建

3. 用给定的 tar 压缩包构建
   ```
   $ docker build http://server/context.tar.gz
   ```

### Dockerfile命令详解
已经介绍了FROM、RUN，还有其他一些常用命令：
* COPY 复制文件
* ADD 更高级的复制文件
* CMD 容器启动命令
* ENTRYPOINT 入口点
* ENV 设置环境变量
* ARG 构建参数
* VOLUME 定义匿名卷
* EXPOSE 暴露端口
* WORKDIR 指定工作目录
* USER 指定当前用户
* HEALTHCHECK 健康检查
* ONBUILD 为他人作嫁衣裳

参考文档：
* Dockerfie 官方文档：https://docs.docker.com/engine/reference/builder/

* Dockerfile 最佳实践文档：https://docs.docker.com/develop/develop-images/dockerfile_best-practices/

* Docker 官方镜像 Dockerfile：https://github.com/docker-library/docs


# 容器
## 启动容器 （docker run）
1. 一次性的：
   >$ docker run ubuntu:18.04 /bin/echo 'Hello world'
2. 启动一个bash终端，有交互：
   >$ docker run -t -i ubuntu:18.04 /bin/bash

当利用 docker run 来创建容器时，Docker 在后台运行的标准操作包括：
   * 检查本地是否存在指定的镜像，不存在就从公有仓库下载
   * 利用镜像创建并启动一个容器
   * 分配一个文件系统，并在只读的镜像层外面挂载一层可读写层
   * 从宿主主机配置的网桥接口中桥接一个虚拟接口到容器中去
   * 从地址池配置一个 ip 地址给容器
   * 执行用户指定的应用程序
   * 执行完毕后容器被终止

3. 启动已终止的容器
   ```
   $ docker container start
   ```
4. 后台运行 -d
   ```
   $ docker run -d ubuntu:18.04 /bin/sh -c "while true; do echo hello world; sleep 1; done"
   ```
   > 查看后台运行日志：
   ```
   $ docker container logs [container ID or NAMES]
   ```
## 终止容器 （docker container stop）
- `docker container stop` -> 终止容器
- `docker container ls -a ` -> 可查看到终止容器
- `docker container start` -> 重新启动终止的容器
- `docker container restart` -> 将容器终止，然后再启动
   
## 导出导入容器
- `docker export 7691a814370e > ubuntu.tar` -> 导出容器
- `cat ubuntu.tar | docker import - test/ubuntu:v1.0` -> 导入容器
- `docker import http://example.com/exampleimage.tgz example/imag
erepo` -> 从url导入容器
>用户既可以使用 docker load 来导入镜像存储文件到本地镜像库，也可以
使用 docker import 来导入一个容器快照到本地镜像库。这两者的区别在于容
器快照文件将丢弃所有的历史记录和元数据信息（即仅保存容器当时的快照状
态），而镜像存储文件将保存完整记录，体积也要大。此外，从容器快照文件导入
时可以重新指定标签等元数据信息。