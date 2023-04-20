package api

import (
	"aurora-v1.0/aiface"
	"aurora-v1.0/anet"
	"auroraTags/mmo_game/core"
	"auroraTags/mmo_game/pb/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	anet.BaseRouter
}

func (m *MoveApi) Handle(req aiface.IRequest) {
	//1.将客户端传来的proto协议解码
	msg := &pb.Position{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("position unmarshal error|", err)
		return
	}

	//2.得知当前的消息是从哪个玩家传递过来的，从连接属性pid中获取
	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("get pid error|", err)
		req.GetConnection().Stop()
		return
	}

	fmt.Printf("user pid = %d , move(%f,%f,%f,%f)", pid, msg.X, msg.Y, msg.Z, msg.V)

	//3.根据pid得到player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4.让player对象发起移动位置信息广播
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
