package main

import (
	"go-tcp-client-agent/core"
	"go-tcp-client-agent/core/model"
	"os"
	"os/signal"
)

func main() {

	cfg := &model.Config{
		Id:   0,
		Host: "127.0.0.1",
		Port: 1922,
	}

	if agent := core.NewTcpClient(cfg); agent != nil {
		agent.Run()
		defer agent.Shutdown()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
