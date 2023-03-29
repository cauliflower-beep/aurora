package main

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"aurora/examples/zinx_client/c_router"
	"fmt"
	"os"
	"os/signal"
	"time"
)

// 客户端自定义业务
func business(conn aiface.IConnection) {

	for {
		err := conn.SendMsg(1, []byte("Ping...[FromClient]"))
		if err != nil {
			fmt.Println(err)
			alog.Error(err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn aiface.IConnection) {
	alog.Debug("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

// 连接断开的时候执行
func DoClientConnectedLost(conn aiface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		alog.Error("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		alog.Error("Conn Property Home = ", home)
	}

	alog.Debug("DoClientConnectedLost is Called ... ")
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := anet.NewClient("127.0.0.1", 8999)

	//添加首次建立链接时的业务
	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	//注册收到服务器消息业务路由
	client.AddRouter(0, &c_router.PingRouter{})

	//启动客户端client
	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}
