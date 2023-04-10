package test

import (
	"aurora-v0.8/anet"
	"fmt"
	"net"
	"testing"
)

// TestUnpack go test -run TestUnpack -v
func TestUnpack(t *testing.T) {
	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	//创建一个封包对象 dp
	dp := anet.NewDataPack()

	//封装一个msg1包
	msg1 := &anet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &anet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	_, _ = conn.Write(sendData1)

	//客户端阻塞
	select {}
}
