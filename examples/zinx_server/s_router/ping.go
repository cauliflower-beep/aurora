package s_router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
)

//ping test 自定义路由
type PingRouter struct {
	anet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request aiface.IRequest) {

	alog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	alog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(2, []byte("pong...pong...pong[FromServer]"))
	if err != nil {
		alog.Error(err)
	}
}
