package main

import (
	"aurora/aiface"
	"aurora/anet"
	"fmt"
	"time"
)

// 用户自定义的心跳检测消息处理方法
func myClientHeartBeatMsg(conn aiface.IConnection) []byte {
	return []byte("heartbeat, I am Client, I am alive")
}

// 用户自定义的远程连接不存活时的处理方法
func myClientOnRemoteNotAlive(conn aiface.IConnection) {
	fmt.Println("myClientOnRemoteNotAlive is Called, connID=", conn.GetConnID(), "remoteAddr = ", conn.RemoteAddr())
	//关闭链接
	conn.Stop()
}

// 用户自定义的心跳检测消息处理方法
type myClientHeartBeatRouter struct {
	anet.BaseRouter
}

func (r *myClientHeartBeatRouter) Handle(request aiface.IRequest) {
	// 业务处理
	fmt.Println("in myClientHeartBeatRouter Handle, recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := anet.NewClient("127.0.0.1", 8999)

	myHeartBeatMsgID := 88888

	//启动心跳检测
	client.StartHeartBeatWithOption(3*time.Second, &aiface.HeartBeatOption{
		MakeMsg:          myClientHeartBeatMsg,
		OnRemoteNotAlive: myClientOnRemoteNotAlive,
		Router:           &myClientHeartBeatRouter{},
		HeadBeatMsgID:    uint32(myHeartBeatMsgID),
	})

	//启动客户端client
	client.Start()

	select {}
}
