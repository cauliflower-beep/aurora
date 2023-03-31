package main

import (
	"aurora-v0.3/aiface"
	"aurora-v0.3/anet"
	"fmt"
)

func main() {
	s := anet.NewServer("[aurora v0.3]") //创建一个server句柄

	// 绑定一个路由处理业务
	s.AddRouter(&PingRouter{})

	s.Serve() //启动server
}

type PingRouter struct {
	anet.BaseRouter //一定要先继承基础 BaseRouter
}

// PreHandle 前置业务
func (pr *PingRouter) PreHandle(req aiface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Handle 主业务
func (pr *PingRouter) Handle(req aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	conn := req.GetConnection().GetTCPConnection()
	data := req.GetData()
	if _, err := conn.Write(data); err != nil {
		fmt.Println("write back buf err|", err)
	}
}

// PostHandle 后置业务
func (pr *PingRouter) PostHandle(req aiface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}
