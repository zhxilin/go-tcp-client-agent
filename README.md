# go-tcp-client-agent
## Features
- An universal multi-layer tcp client framework.
- Universal interface `ILogger` and `IParser`, support custom middleware for logger and parser.
- Support kinds of message protocol, such as proto, thrift, json, xml etc.
- Provide event queue injection, decouple business logic and network layer.
- Seperate send and receive using channel.
- Sync connection state in multiple goroutine.

## Requirements
- Golang v1.11+ Tested (go mod suppored)
- Protoc & Protoc-gen-go
```shell
$ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

## Quick start
The example instance using proto, follow the steps to run:

1. Download dependecies
```shell
$ go mod download
```

2. Generate proto message
```shell
$ ./genmsg.sh
```
It will auto convert `.pb` files under core/proto to `.pb.go`

3. Build and run
```shell
$ make build
$ ./build/go-tcp-client-agent
```

Or run directly

```shell
$ make run
```

## How to use

- Main code usage

```golang
//define config info to connect to remote server.
cfg := &model.Config{
    Id:   0,
    Host: "127.0.0.1",
    Port: 1922,
}

//create a tcp client then begin to run.
if agent := core.NewTcpClient(cfg); agent != nil {
    agent.Run()
    defer agent.Shutdown()
}
```

- Receive message from server
In the `core/agent.go`, focus on the `initHandlers` function. 
Add your own message handler for different message id there.

```golang
func (cli *GtaTcpClient) initHandlers() {
	cli.handlers = make(map[uint16]gtaMsgHandler)

	// Add your own msg handler here.
    //e.g.
    addHandler(100, MyFirstHandler)
}

func MyFirstHandler(msgID uint16, data []byte, len int) error {
    //Deserialize data in your own protocol.
    ack := &msg.MyFirstAck{}
    proto.Unmarshal(data, ack)
    //...
}
```

- Send message to server
In the `core/agent.go`, using `GtaTcpClient.conn.Send()` function to send message.

```golang
func (cli *GtaTcpClient) register() error {
	//Call cli.conn.Send() to send your own message to server.
	req := &msg.MyFirstReq{}
	err := cli.conn.Send(100, req)
	return err
}
```