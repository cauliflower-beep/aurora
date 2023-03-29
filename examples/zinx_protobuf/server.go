package main

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"aurora/aurora_app_demo/mmo_game/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type PositionServerRouter struct {
	anet.BaseRouter
}

//Ping Handle
func (this *PositionServerRouter) Handle(request aiface.IRequest) {

	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Position Unmarshal error ", err, " data = ", request.GetData())
		return
	}

	fmt.Printf("recv from client : msgId=%+v, data=%+v\n", request.GetMsgID(), msg)

	msg.X += 1
	msg.Y += 1
	msg.Z += 1
	msg.V += 1

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("proto Marshal error = ", err, " msg = ", msg)
		return
	}

	err = request.GetConnection().SendMsg(0, data)

	if err != nil {
		alog.Error(err)
	}
}

func main() {
	//创建一个server句柄
	s := anet.NewServer()

	//配置路由
	s.AddRouter(0, &PositionServerRouter{})

	//开启服务
	s.Serve()
}
