package main

import (
	"aurora-v1.0/aiface"
	"aurora-v1.0/anet"
	"fmt"
)

type PingRouter struct {
	anet.BaseRouter //一定要先继承基础 BaseRouter
}

// Handle 主业务
func (pr *PingRouter) Handle(req aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))

	//回写数据
	err := req.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloAuroraRouter struct {
	anet.BaseRouter
}

func (ha *HelloAuroraRouter) Handle(req aiface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))

	err := req.GetConnection().SendMsg(1, []byte("Hello Aurora Router v0.9"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnBegin
// @Description: 创建连接的时候执行
func DoConnBegin(conn aiface.IConnection) {
	fmt.Println("DoConnBegin is called ...")

	//========在创建链接之后，设置两个链接属性===========
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Kevin")
	conn.SetProperty("age", 17)
	//=============================================
	err := conn.SendMsg(2, []byte("doConn begin..."))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnLost
// @Description: 连接断开的时候执行
func DoConnLost(conn aiface.IConnection) {
	//========在链接销毁之前，查询conn的Name，Home属性===
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	//==============================================
	fmt.Println("doConnLost is called")
}
func main() {
	s := anet.NewServer() //创建一个server句柄

	//注册连接hook回调函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	// 绑定一个路由处理业务
	s.AddRouter(0, &PingRouter{})        //添加路由1
	s.AddRouter(1, &HelloAuroraRouter{}) //添加路由2

	s.Serve() //启动server
}