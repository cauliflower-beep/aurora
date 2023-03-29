/**
* @Author: Aceld
* @Date: 2020/12/24 00:24
* @Mail: danbing.at@gmail.com
*    zinx server demo
 */
package main

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"aurora/examples/zinx_server/s_router"
)

// 创建连接的时候执行
func DoConnectionBegin(conn aiface.IConnection) {
	alog.Ins().InfoF("DoConnecionBegin is Called ...")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.kancloud.cn/@aceld")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		alog.Error(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn aiface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		alog.Ins().InfoF("Conn Property Name = %v", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		alog.Ins().InfoF("Conn Property Home = %v", home)
	}

	alog.Ins().InfoF("Conn is Lost")
}

func main() {
	//创建一个server句柄
	s := anet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	//开启服务
	s.Serve()
}
