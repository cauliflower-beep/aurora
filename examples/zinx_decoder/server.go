package main

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"aurora/examples/zinx_decoder/decode"
	"aurora/examples/zinx_decoder/router"
)

func DoConnectionBegin(conn aiface.IConnection) {
	alog.Ins().InfoF("DoConnectionBegin is Called ...")
}

func DoConnectionLost(conn aiface.IConnection) {
	alog.Ins().InfoF("Conn is Lost")
}

func main() {
	//创建一个server句柄
	s := anet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//s.AddRouter(0x00000001, &router.TLVBusinessRouter{}) //TLV协议对应业务功能
	//处理HTLVCRC协议数据
	s.SetDecoder(decode.NewHTLVCRCDecoder())
	s.AddRouter(0x10, &router.HtlvCrcBusinessRouter{}) //TLV协议对应业务功能，因为client.go中模拟数据funcode字段为0x10
	s.AddRouter(0x13, &router.HtlvCrcBusinessRouter{}) //TLV协议对应业务功能，因为client.go中模拟数据funcode字段为0x13

	//开启服务
	s.Serve()
}
