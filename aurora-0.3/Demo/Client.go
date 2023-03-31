package main

import (
	"fmt"
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
		_, err = conn.Write([]byte("hello aurora v0.3! "))
		if err != nil {
			fmt.Println("write conn err", err)
		}

		buf := make([]byte, 512)
		//var cnt int
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		//fmt.Println(fmt.Sprintf("server call back:%scnt = %d", buf, cnt))
		fmt.Println(string(buf))
		//这里要加一个cpu阻塞，以便于cpu去处理别的事情，否则会一直卡在这个循环中，过分消耗资源
		time.Sleep(1 * time.Second)
	}
}
