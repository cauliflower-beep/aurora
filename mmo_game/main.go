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

	//========同步周边玩家的上线信息
	player.SyncSurrounding()

	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

// OnConnectionRemove
// @Description: 连接丢失之后的hook  todo
func OnConnectionRemove(conn aiface.IConnection) {
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("get pid failed")
		return
	}
	//========得到当前掉线的玩家
	player := core.WorldMgrObj.Players[pid.(int32)]

	//========将当前离线的玩家添从 WorldMgrObj 中删除
	core.WorldMgrObj.RemovePlayer(player)

	fmt.Println("=====> Player pidId = ", 0, " offline ====")
}
func main() {
	//创建服务器句柄
	s := anet.NewServer()

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionRemove)

	//注册路由
	s.AddRouter(2, &api.WorldChatApi{}) //世界聊天
	s.AddRouter(3, &api.MoveApi{})

	//启动服务
	s.Serve()
}
