package main

import (
	"aurora/anet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
模拟客户端
 */
func main(){
	fmt.Println("client1 starting...")
	time.Sleep(1 *time.Second)
	//连接远程服务器，得到一个conn连接
	conn ,err := net.Dial("tcp","127.0.0.1:8999")
	if err != nil{
		fmt.Println("client1 start err,exit!")
		return
	}

	//连接调用write写数据
	for {
		//发封包message消息
		dp := anet.NewDataPack()
		msg, err := dp.Pack(anet.NewMsgPackage(1, []byte("Aurora V0.6 Client1 Test Message")))
		if err != nil {
			fmt.Println("pack error ", err)
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		//先读出流中的head部分,得到id和len
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*anet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}