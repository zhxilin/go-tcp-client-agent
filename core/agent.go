package core

import (
	"go-tcp-client-agent/core/intf"
	"go-tcp-client-agent/core/middleware"
	"go-tcp-client-agent/core/model"
	"go-tcp-client-agent/core/service"
	"net"
	"strconv"
	"time"
)

type gtaMsgHandler func(uint16, []byte, int) error

type GtaTcpClient struct {
	config     *model.Config
	conn       *service.GtaConnection
	eventQueue chan *model.EventItem
	handlers   map[uint16]gtaMsgHandler
	logger     intf.ILogger
}

func NewTcpClient(cfg *model.Config) *GtaTcpClient {
	logger := middleware.NewLogger(cfg.Id, false, true)

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	logger.LogInfo("Ready to connect to proxy server [%s]", addr)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		logger.LogError("Connect to proxy server [%s] failed", addr)
		return nil
	}

	connection, err := service.NewConnection(c)
	if err != nil {
		return nil
	}

	t := &GtaTcpClient{
		config:     cfg,
		conn:       connection,
		eventQueue: make(chan *model.EventItem, 1024),
		logger:     logger,
	}

	connection.SetLogger(logger)
	connection.SetParser(middleware.NewMsgParser())
	connection.SetEventQueue(t)

	return t
}

func (cli *GtaTcpClient) Run() {
	go cli.onConnEvent()
	cli.initHandlers()
	cli.conn.Run()
}

func (cli *GtaTcpClient) Shutdown() {
	cli.logger.LogInfo("Shutdown tcp client.")
	close(cli.eventQueue)
	cli.conn.Shutdown()

	cli.logger.Close()
}

func (cli *GtaTcpClient) Push(evt *model.EventItem) {
	if cli.eventQueue == nil {
		return
	}

	for {
		select {
		case cli.eventQueue <- evt:
			{
				return
			}
		case <-time.After(5 * time.Second):
			{
				return
			}
		}
	}
}

func (cli *GtaTcpClient) Pop() <-chan *model.EventItem {
	return cli.eventQueue
}

func (cli *GtaTcpClient) onConnEvent() {
	if cli.eventQueue == nil {
		return
	}

	for {
		evt := <-cli.Pop()
		if evt == nil {
			return
		}

		// cli.logger.LogTrace("On Socket Event: %d", evt.EventType)
		if evt.Type == service.EEvType_Data {
			recvTask := evt.UserData.(*service.RecvTask)
			cli.dispatchMessage(recvTask.MsgID, recvTask.Data, recvTask.Len)
		} else if evt.Type == service.EEvType_Connected {
			cli.register()
		}
	}
}

func (cli *GtaTcpClient) dispatchMessage(msgID uint16, payload []byte, len int) {
	cli.logger.LogTrace("Recv msg id: %d, data len: %d", msgID, len)

	handler, found := cli.handlers[msgID]
	if !found {
		return
	}

	handler(msgID, payload, len)
}

func (cli *GtaTcpClient) initHandlers() {
	cli.handlers = make(map[uint16]gtaMsgHandler)

	// Add your own msg handler here.
}

func (cli *GtaTcpClient) addHandler(msgID uint16, handler gtaMsgHandler) {
	cli.handlers[msgID] = handler
}

func (cli *GtaTcpClient) register() error {
	//Call cli.conn.Send() to send your own message to server.

	return nil
}
