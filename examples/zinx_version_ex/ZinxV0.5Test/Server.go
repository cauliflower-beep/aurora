package main

import (
	"fmt"

	"aurora/aiface"
	"aurora/anet"
)

//ping test 自定义路由
type PingRouter struct {
	anet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回写数据
	/*
		_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
		if err != nil {
			fmt.Println("call back ping ping ping error")
		}
	*/
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server句柄
	s := anet.NewServer()

	//配置路由
	s.AddRouter(5, &PingRouter{})

	//开启服务
	s.Serve()
}
