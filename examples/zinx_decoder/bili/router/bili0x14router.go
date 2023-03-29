package router

import (
	"aurora/aiface"
	"aurora/anet"
	"aurora/examples/zinx_decoder/bili/utils"
	"aurora/examples/zinx_decoder/decode"
	"bytes"
	"fmt"
)

type Data0x14Router struct {
	anet.BaseRouter
}

func (this *Data0x14Router) Handle(request aiface.IRequest) {
	fmt.Println("Data0x14Router Handle", request.GetMessage().GetData())
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case decode.HtlvCrcData:
			_data := _response.(decode.HtlvCrcData)
			fmt.Println("Data0x14Router", _data)
			buffer := pack14(_data)
			request.GetConnection().Send(buffer)
		}
	}
}

// 头码   功能码 数据长度      Body                         CRC
// A2      10     0E        0102030405060708091011121314 050B
func pack14(_data decode.HtlvCrcData) []byte {
	_data.Data[0] = 0xA1
	buffer := bytes.NewBuffer(_data.Data[:len(_data.Data)-2])
	crc := utils.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
