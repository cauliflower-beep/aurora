package main

import (
	"aurora/aiface"
	"aurora/anet"
	"time"
)

func main() {
	// 客户端向服务端发送消息
	client := anet.NewClient("127.0.0.1", 8999)
	client.SetOnConnStart(func(connection aiface.IConnection) {
		_ = connection.SendMsg(1, []byte("hello zinx"))
	})
	client.Start()
	time.Sleep(time.Second)
}
