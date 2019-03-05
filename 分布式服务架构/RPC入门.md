
# RPC是什么

RPC框架是构成微服务最重要的组成部分之一。

# RPC组件
## 服务发现

## 服务治理

## 远程调用

## 调用链分析

## 网关

# RPC调用过程
Client——》Proxy——》Protocol（通过协议序列化字节流）——》network（可通过netty网络框架）——————netWork——》Protocol——》Proxy——》Service

# Protocol  
序列化成byte[]数组，方便网络传输
序列化可选协议：
- jdk的序列化方法（不利于跨语言调用）
- json可读性强，但是速度慢，体积大。 jackson是json的解析框架之一。
- protobuf，kyro，Hessian等

# Server

# Client 
