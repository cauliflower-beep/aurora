package router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
)

type HelloRouter struct {
	anet.BaseRouter
}

func (hr *HelloRouter) Handle(request aiface.IRequest) {
	alog.Ins().InfoF(string(request.GetData()))
}
