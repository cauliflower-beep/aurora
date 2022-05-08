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
	fmt.Println("recv from client:msgId = ",request.GetMsgID(),
		",data=",string(request.GetData()))
	err := request.GetConnection().SendMsg(200,[]byte("ping..ping..ping"))
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
	fmt.Println("recv from client:msgId = ",request.GetMsgID(),
		",data=",string(request.GetData()))
	err := request.GetConnection().SendMsg(1,[]byte("hahaha...you are so funny"))
	if err != nil {
		fmt.Println(err)
	}
}


func main() {
	//创建一个server句柄
	s := anet.NewServer("[zinx V0.6]")
	//给当前服务端添加一个自定义的router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&Zxx{})

	//2 开启服务
	s.Serve()
}
