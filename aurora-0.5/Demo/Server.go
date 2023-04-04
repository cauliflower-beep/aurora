package main

import (
	"aurora-v0.5/aiface"
	"aurora-v0.5/anet"
	"fmt"
)

func main() {
	s := anet.NewServer() //创建一个server句柄

	// 绑定一个路由处理业务
	s.AddRouter(&PingRouter{})

	s.Serve() //启动server
}

type PingRouter struct {
	anet.BaseRouter //一定要先继承基础 BaseRouter
}

// Handle 主业务
func (pr *PingRouter) Handle(req aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))

	//回写数据
	err := req.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}
