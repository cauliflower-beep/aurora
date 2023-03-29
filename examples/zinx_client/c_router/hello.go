package c_router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
)

type HelloRouter struct {
	anet.BaseRouter
}

//HelloZinxRouter Handle
func (this *HelloRouter) Handle(request aiface.IRequest) {
	alog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	alog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

}
