# 自己积累常用
- docker rm $(docker ps -a -q) //删除容器，-a是列出所有的容器，-q是列出ID属性。
- docker container prune //删除所有处于终止状态的容器
- docker container rm [容器ID] //删除容器
- docker image rm （可以是镜像短 ID 、 镜像长 ID 、 镜像名 或者 镜像摘要，短ID可以只指出前几位字符）
- docker pull [选项] [Docker Registry 地址[:端口号]/]仓库名[:标签]
  > docker pull ubuntu:18.04   //默认地址为官方镜像Docker Hub 用户library，因此将会获取官方镜像 library/ubuntu
仓库中标签为 18.04 的镜像
- docker image ls //查看镜像，加 -a 可以查看包括中间镜像。还可以有过滤参数，比如 -q 只列出ID。或者指定格式：
  ```shell
  docker image ls --format "{{.ID}}: {{.Repository}}"
  ```
- docker container ls  //查看容器
- docker container stop [容器ID] //关闭容器，start和restart可以打开
- docker run ubuntu:18.04 /bin/echo 'Hello world' //启动容器
  > docker run -t -i ubuntu:18.04 /bin/bash  //开个伪终端交互模式，使用`-d`容器会进入后台，可以查看日志获取容器运行信息，`docker container logs [ID]`
- docker





# 所有命令

- docker build -t friendlyhello .  # Create image using this directory's dockerfile
- docker run -p 4000:80 friendlyhello  # Run "friendlyhello" mapping port 4000 to 80
- docker run -d -p 4000:80 friendlyhello         # Same thing, but in detached mode
- docker container ls                                # List all running containers
- docker container ls -a             # List all containers, even those not running
- docker container stop <hash>           # Gracefully stop the specified container
- docker container kill <hash>         # Force shutdown of the specified container
- docker container rm <hash>        # Remove specified container from this machine
- docker container rm $(docker container ls -a -q)         # Remove all containers
- docker image ls -a                             # List all images on this machine
- docker image rm <image id>            # Remove specified image from this machine
- docker image rm $(docker image ls -a -q)   # Remove all images from this machine
- docker login             # Log in this CLI session using your docker credentials
- docker tag <image> username/repository:tag  # Tag <image> for upload to registry
- docker push username/repository:tag            # Upload tagged image to registry
- docker run username/repository:tag                   # Run image from a registry




