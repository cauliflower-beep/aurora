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

//Test PreHandle
func (this *PingRouter) PreHandle(request aiface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test Handle
func (this *PingRouter) Handle(request aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
func (this *PingRouter) PostHandle(request aiface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {
	//创建一个server句柄
	s := anet.NewServer("[zinx V0.3]")
	//给当前服务端添加一个自定义的router
	s.AddRouter(&PingRouter{})

	//2 开启服务
	s.Serve()
}