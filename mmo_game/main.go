package main

import (
	"aurora-v1.0/aiface"
	"aurora-v1.0/anet"
	"auroraTags/mmo_game/api"
	"auroraTags/mmo_game/core"
	"fmt"
)

// OnConnectionAdd
// @Description: 绑定hook，自动触发
// @Description: 在连接创建完毕之后服务器自动回复客户端玩家ID和坐标
func OnConnectionAdd(conn aiface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)

	player.SyncPid()

	player.BroadCastStartPosition()

	//========将当前新上线的玩家添加到 WorldMgrObj 中
	core.WorldMgrObj.AddPlayer(player)

	//========将当前新上线的玩家链接属性绑定pid
	player.Conn.SetProperty("pid", player.Pid)

	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}
func main() {
	//创建服务器句柄
	s := anet.NewServer()

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)

	//注册路由
	s.AddRouter(2, &api.WorldChatApi{}) //世界聊天

	//启动服务
	s.Serve()
}
