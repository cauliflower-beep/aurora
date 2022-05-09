package main

import (
	"aurora/aiface"
	"aurora/anet"
	"fmt"
)

//ping test 自定义路由
type PingRouter struct {
	anet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端数据，再回写ping..ping..ping
	fmt.Println("recv from client:msgId = ", request.GetMsgID(),
		",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}

//ping test 自定义路由
type Zxx struct {
	anet.BaseRouter
}

//Test Handle
func (this *Zxx) Handle(request aiface.IRequest) {
	fmt.Println("im zxx!im listening!!")
	//先读取客户端数据，再回写hahaha
	fmt.Println("recv from client:msgId = ", request.GetMsgID(),
		",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("hahaha...you are so funny"))
	if err != nil {
		fmt.Println(err)
	}
}

//创建连接之后执行钩子函数
func DoConnBegin(conn aiface.IConnection) {
	fmt.Println("---->doConnBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnBegin")); err != nil {
		fmt.Println(err)
	}
}

//连接断开之前需要执行的函数
func DoConnStop(conn aiface.IConnection) {
	fmt.Println("---->doConnStop is Called...")
	fmt.Println("connID =", conn.GetConnID(), "is Lost...")
}

func main() {
	//创建一个server句柄
	s := anet.NewServer("[zinx V0.6]")

	//注册钩子函数
	s.RegOnConnStart(DoConnBegin)
	s.RegOnConnStop(DoConnStop)

	//给当前服务端添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &Zxx{})

	//2 开启服务
	s.Serve()
}
