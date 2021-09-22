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
	err := request.GetConnection().SendMsg(1,[]byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}


func main() {
	//创建一个server句柄
	s := anet.NewServer("[zinx V0.5]")
	//给当前服务端添加一个自定义的router
	s.AddRouter(&PingRouter{})

	//2 开启服务
	s.Serve()
}
