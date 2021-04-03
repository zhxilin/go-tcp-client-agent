package service

import (
	"encoding/binary"
	"errors"
	"go-tcp-client-agent/core/intf"
	"go-tcp-client-agent/core/model"
	"io"
	"net"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	eConnStatus_None = iota
	eConnStatus_Connected
	eConnStatus_Disconnected
)

const (
	EEvType_None = iota
	EEvType_Connected
	EEvType_Disconnected
	EEvType_Data
)

type SendTask struct {
	Data []byte
}

type RecvTask struct {
	MsgID uint16
	Data  []byte
	Len   int
}

type GtaConnection struct {
	conn       net.Conn
	status     int32
	sendQueue  chan *SendTask
	eventQueue intf.IEventQueue
	parser     intf.IParser
	logger     intf.ILogger
}

func NewConnection(c net.Conn) (*GtaConnection, error) {
	return &GtaConnection{
		conn:      c,
		status:    eConnStatus_None,
		sendQueue: make(chan *SendTask, 128),
	}, nil
}

func newEvent(et int, data interface{}) *model.EventItem {
	return &model.EventItem{
		Type:     et,
		UserData: data,
	}
}

func (c *GtaConnection) SetEventQueue(q intf.IEventQueue) {
	c.eventQueue = q
}

func (c *GtaConnection) SetParser(p intf.IParser) {
	c.parser = p
}

func (c *GtaConnection) SetLogger(l intf.ILogger) {
	c.logger = l
}

func (c *GtaConnection) Run() {
	if c.logger == nil {
		panic("No logger")
	}

	if c.parser == nil {
		panic("No parser")
	}

	if c.eventQueue == nil {
		panic("No event queue")
	}

	go func() {
		defer func() {
			c.close()
			c.pushEvent(eConnStatus_Disconnected, nil)
		}()

		c.logger.LogInfo("Connected")
		atomic.StoreInt32(&c.status, eConnStatus_Connected)
		c.pushEvent(EEvType_Connected, nil)

		go c.processSend()
		c.processRecv()
	}()
}

func (c *GtaConnection) Shutdown() {
	if atomic.LoadInt32(&c.status) != eConnStatus_Connected {
		return
	}

	select {
	case c.sendQueue <- nil:
		{

		}
	case <-time.After(5 * time.Second):
		{
			c.close()
		}
	}
}

func (c *GtaConnection) Send(msgID uint16, m proto.Message) error {
	msg := c.parser.Pack(msgID, m)
	return c.send(msg)
}

func (c *GtaConnection) close() {
	if atomic.CompareAndSwapInt32(&c.status, eConnStatus_Connected, eConnStatus_Disconnected) {
		c.conn.Close()
	}
}

func (c *GtaConnection) pushEvent(et int, data interface{}) {
	if c.eventQueue == nil {
		return
	}

	c.eventQueue.Push(newEvent(et, data))
}

func (c *GtaConnection) processSend() error {
	defer func() {
		c.logger.LogInfo("Send Process End")
	}()

	if atomic.LoadInt32(&c.status) != eConnStatus_Connected {
		return errors.New("connection closed")
	}

	for {
		evt, ok := <-c.sendQueue
		if !ok {
			return nil
		}

		if evt == nil {
			c.close()
			return nil
		}

		_, err := c.conn.Write(evt.Data)
		if err != nil {
			return err
		}
	}
}

func (c *GtaConnection) processRecv() error {
	defer func() {
		c.logger.LogInfo("Recv Process End")
	}()

	buffer := make([]byte, 1024)

	for {
		bufLen, err := c.conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				c.logger.LogError("Conn Read Error: %s", err)
			}
			break
		}

		offset := 0
		leftLen := bufLen
		for {

			id, payload, packSize, err := c.parser.UnPack(buffer[offset:], leftLen)
			if err != nil {
				continue
			}

			task := &RecvTask{
				MsgID: id,
				Data:  payload,
				Len:   binary.Size(payload),
			}
			c.pushEvent(EEvType_Data, task)

			offset += packSize
			leftLen -= packSize

			if leftLen <= 0 {
				break
			}
		}
	}
	return nil
}

func (c *GtaConnection) send(msg []byte) error {
	task := &SendTask{
		Data: msg,
	}

	select {
	case c.sendQueue <- task:
		{

		}
	case <-time.After(5 * time.Second):
		{
			c.close()
			return errors.New("send timeout")
		}
	}

	return nil
}

func (c *GtaConnection) IsConnected() bool {
	return atomic.LoadInt32(&c.status) == eConnStatus_Connected
}
