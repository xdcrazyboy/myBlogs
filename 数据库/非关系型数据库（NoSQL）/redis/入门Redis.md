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

# 入门

## 数据类型

### 字符串类型

- `INCR` 自增  （需要是原子操作）
  > 实践： 1. 文章访问量统计； 2.生成自增ID； 3. 存储文章数据
- `INCRBY key increment`  -- 可以指定一次增加的数值，increment参数表示该数值
- 对应的减少命令有`DECR`、`DECRBY`
- `INCRBYFLOAT` 增加指定浮点数
- `APPEND key value` 向尾部追加 值（字符串）
- `STRLEN key` 获取字符串长度
- `MGET/MSET key1 (v1) key2 (v2) ...`通知获取多个键值
- 位操作：`GETBIT`、`SETBIT`、`BITCOUNT`、`BITOP`