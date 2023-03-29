package c_router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"fmt"
)

//ping test 自定义路由
type PingRouter struct {
	anet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request aiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	alog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	alog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	fmt.Println("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	if err := request.GetConnection().SendBuffMsg(1, []byte("Hello[from client]")); err != nil {
		alog.Error(err)
	}
}
