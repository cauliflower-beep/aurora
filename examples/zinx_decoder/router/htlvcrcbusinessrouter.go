package router

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"aurora/examples/zinx_decoder/decode"
)

type HtlvCrcBusinessRouter struct {
	anet.BaseRouter
}

func (this *HtlvCrcBusinessRouter) Handle(request aiface.IRequest) {
	//alog.Ins().DebugF("Call HtlvCrcBusinessRouter Handle %d %s\n", request.GetMessage().GetMsgID(), hex.EncodeToString(request.GetMessage().GetData()))
	msgID := request.GetMessage().GetMsgID()
	if msgID == 0x10 {
		_response := request.GetResponse()
		if _response != nil {
			switch _response.(type) {
			case decode.HtlvCrcData:
				tlvData := _response.(decode.HtlvCrcData)
				alog.Ins().DebugF("do msgid=0x10 data business %+v\n", tlvData)
			}
		}
	}
}
