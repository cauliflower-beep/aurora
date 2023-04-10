package main

import (
	"aurora-v0.8/anet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client starting...")
	time.Sleep(1 * time.Second)
	//连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	//连接调用write写数据
	for {
		// 发送封包 message 消息
		dp := anet.NewDataPack()
		msg, _ := dp.Pack(anet.NewMsgPackage(1, []byte("aurora v0.7 client2 test message...")))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write conn err", err)
		}

		// 接收服务器返回
		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error|", err)
		}

		// 将headData字节流，拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err|", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// msg 是有data数据的，需要再次读取data数据
			serverMsg := msgHead.(*anet.Message)
			serverMsg.Data = make([]byte, serverMsg.GetDataLen())

			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, serverMsg.Data)
			if err != nil {
				fmt.Println("server unpack data err|", err)
				return
			}
			fmt.Println("==> recv msg:ID=", serverMsg.Id, ", len=", serverMsg.DataLen, ", data=", string(serverMsg.Data))
		}
		//这里要加一个cpu阻塞，以便于cpu去处理别的事情，否则会一直卡在这个循环中，过分消耗资源
		time.Sleep(1 * time.Second)
	}
}
