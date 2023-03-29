package router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"aurora/apack"
)

type TLVBusinessRouter struct {
	anet.BaseRouter
}

func (this *TLVBusinessRouter) Handle(request aiface.IRequest) {
	alog.Ins().DebugF("Call TLVRouter Handle %d %+v\n", request.GetMessage().GetMsgID(), request.GetMessage().GetData())
	msgID := request.GetMessage().GetMsgID()
	if msgID == 0x00000001 {
		_response := request.GetResponse()
		if _response != nil {
			switch _response.(type) {
			case apack.TLVDecoder:
				tlvData := _response.(apack.TLVDecoder)
				alog.Ins().DebugF("do msgid=0x00000001 data business %+v\n", tlvData)
			}
		}
	}

}
