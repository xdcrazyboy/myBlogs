# 安装
## Linux
下载稳定版本：
> 没有windows版本
``` shell
//下载
wget http://download.redis.io/redis-stable.tar.gz
//解压
tar xzf redis-stable.tar.gz
//编译
cd redis-stable
make
//将可执行程序复制到/usr/local/bin目录中，方便后面执行程序不需要输入完整路径
make install
//推荐运行之前做好测试
make test

```