# go-tcp-client-agent
## Features
- 一个通用的TCP客户端分层框架
- 统一ILogger和IParser接口,支持自定义logger和parser中间件,支持多种消息协议
- 通过依赖注入的形式提供事件队列,解耦业务层和网络层.
- 使用channel分离消息收发通道
- 连接状态一致性

## How to
本实例结合proto协议说明使用步骤

1. 生成proto协议
```shell
$ ./genmsg.sh
```
将core/proto目录下的pb协议转成.pb.go

2. 生成并运行
```shell
$ make build
$ ./build/go-tcp-client-agent
```
或者直接运行
```shell
$ make run
```
