package main

import (
	"aurora/aiface"
	"aurora/anet"
	"aurora/examples/zinx_decoder/bili/router"
	"aurora/examples/zinx_decoder/decode"
)

func DoConnectionBegin(conn aiface.IConnection) {
}

func DoConnectionLost(conn aiface.IConnection) {
}

func main() {
	server := anet.NewServer(func(s *anet.Server) {
		s.Port = 9090
		/*
			s.LengthField = aiface.LengthField{
				MaxFrameLength:      math.MaxUint8 + 4,
				LengthFieldOffset:   2,
				LengthFieldLength:   1,
				LengthAdjustment:    2,
				InitialBytesToStrip: 0,
			}
		*/
	})
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	server.AddInterceptor(decode.NewHTLVCRCDecoder())
	server.AddRouter(0x10, &router.Data0x10Router{})
	server.AddRouter(0x13, &router.Data0x13Router{})
	server.AddRouter(0x14, &router.Data0x14Router{})
	server.AddRouter(0x15, &router.Data0x15Router{})
	server.AddRouter(0x16, &router.Data0x16Router{})
	server.Serve()
}
