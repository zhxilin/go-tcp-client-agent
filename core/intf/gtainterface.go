package intf

import (
	"go-tcp-client-agent/core/model"

	"google.golang.org/protobuf/proto"
)

type IParser interface {
	Pack(uint16, proto.Message) []byte
	UnPack([]byte, int) (uint16, []byte, int, error)
}

type IEventQueue interface {
	Push(*model.EventItem)
	Pop() <-chan *model.EventItem
}

type ILogger interface {
	LogTrace(f string, v ...interface{})
	LogInfo(f string, v ...interface{})
	LogWarn(f string, v ...interface{})
	LogError(f string, v ...interface{})
	LogFatal(f string, v ...interface{})
	Close()
}
