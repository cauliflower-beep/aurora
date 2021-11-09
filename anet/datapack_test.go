package anet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只是负责测试datapack拆包、封包的单元测试
func TestDataPack(t *testing.T) {
	/*
	模拟服务器
	 */
	//1、创建socketTCP
	listener,err := net.Listen("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建一个go承载 负责从客户端处理业务
	go func() {
		//2、从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				 fmt.Println("server accept err:", err)
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				//----->拆包的过程<------
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for{
					//1 第一次从conn读，先读出流中的head部分
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)  //ReadFull 会把msg填充满为止
					if err != nil {
						fmt.Println("read head error")
					}
					//将headData字节流 拆包到msg中
					msgHead,err := dp.Unpack(headData)
					if err != nil{
						fmt.Println("server unpack err:", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						//msg 是有data数据的，需要再次读取data数据
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						//根据dataLen从io中读取字节流
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil{
		fmt.Println("client dial err:", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg包一起发送，看是否能分开获取
	//封装第一个msg1包
	msg1 := &Message{
		Id:0,
		DataLen:5,
		Data:[]byte{'h', 'e', 'l', 'l', 'o'},
	}

	//打包msg1，变成一个二进制文件
	sendData1, err := dp.Pack(msg1)
	if err!= nil{
		fmt.Println("client pack msg1 err:", err)
		return
	}

	//封装第二个msg包
	msg2 := &Message{
		Id:1,
		DataLen:7,
		Data:[]byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err!= nil{
		fmt.Println("client temp msg2 err:", err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	conn.Write(sendData1)

	//客户端阻塞
	select{}
	//测试的时候注意先把全局配置文件的读取过程注释掉
}
