
# 来源书籍

[《Docker — 从入门到实践》](https://github.com/yeasy/docker_practice) ，大家看详细版本，可以直接去这里看。

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

容器的实质是**进程**，但它运行于属于自己的独立的命名空间。 容器拥有自己的root文件系统、自己的网络配置、自己的进程空间、甚至自己的用户ID空间。

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