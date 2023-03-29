package s_router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
)

type HelloZinxRouter struct {
	anet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request aiface.IRequest) {
	alog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	alog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(3, []byte("Hello Zinx Router[FromServer]"))
	if err != nil {
		alog.Error(err)
	}
}
