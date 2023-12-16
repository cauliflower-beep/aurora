package test

import (
	"aurora-v0.5/anet"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	fmt.Printf("pack res. sendData1|%v\n", sendData1)

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
	fmt.Printf("pack res. sendData2|%v\n", sendData2)

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)
	fmt.Printf("TCP stream. sendData|%v\n", sendData1)

	//向服务器端写数据
	_, _ = conn.Write(sendData1)

	//客户端阻塞
	select {}
}

// TestBinaryRead go test -run TestBinaryRead -v
func TestBinaryRead(t *testing.T) {
	type Msg struct {
		DataLen uint32
		Id      uint32
		Data    uint32
	}
	msg := &Msg{}
	binaryData := []byte{7, 0, 0, 0, 1, 0, 0, 0}
	dataBuff := bytes.NewReader(binaryData)

	// 读取dataLen
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("dataLen|%d\n", msg.DataLen)

	// dataLen读完之后，下次读取的位置被移动到了前4个字节之后，所以再读就是从第5个字节开始的
	_ = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	fmt.Printf("Id|%d\n", msg.Id)

	// 此时8个字节都已经读完了，再读就是空值
	_ = binary.Read(dataBuff, binary.LittleEndian, &msg.Data)
	fmt.Printf("Data|%d\n", msg.Data)

	// 如果想从头开始读，可以用 Seek 方法将位置重置
	_, _ = dataBuff.Seek(0, io.SeekStart)
	_ = binary.Read(dataBuff, binary.LittleEndian, &msg.Data)
	fmt.Printf("Data|%d\n", msg.Data)
}
