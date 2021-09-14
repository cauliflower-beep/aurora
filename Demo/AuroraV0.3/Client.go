package main

import (
	"fmt"
	"net"
	"time"
)

/*
模拟客户端
 */
func main(){
	fmt.Println("client starting...")
	time.Sleep(1 *time.Second)
	//连接远程服务器，得到一个conn连接
	conn ,err := net.Dial("tcp","127.0.0.1:8999")
	if err != nil{
		fmt.Println("client start err,exit!")
		return
	}

	//连接调用write写数据
	for {
		_,err := conn.Write([]byte("hello aurora V0.2 !"))
		if err != nil{
			fmt.Println("write conn err",err)
		}
		buf := make([]byte,512)
		cnt,err := conn.Read(buf)
		if err != nil{
			fmt.Println("read buf error")
			return
		}
		//为了能够原样显示，还需要把byte转换为string
		fmt.Println("server call back:%s,cnt = %d\n",string(buf[:cnt]),cnt)
		//这里要加一个cou阻塞，以便于cpu去处理别的事情，否则会一直卡在这个循环中，过分消耗资源
		time.Sleep(1 *time.Second)
	}
}