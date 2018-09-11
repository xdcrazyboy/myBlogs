
## HTTP/1.1  VS HTTP/2.0
1. HTTP/2.0 对消息头采用了HPACK压缩，**提升了传输效率**。
2. **基于帧和流的多路复用**，真正实现了基于一个链接的多请求并发处理。
3. 支持服务器推送


## SPDY/2  VS HTTP/2.0

HTTP/2.0是托体于SPDY/2.0（Google整的），但还是有些区别：
1. HTTP/2.0 **支持明文传输**，SPDY强制使用HTTPS。
2. HTTP/2.0消息头采用HPACK压缩，而SPDY采用DELEFT。

