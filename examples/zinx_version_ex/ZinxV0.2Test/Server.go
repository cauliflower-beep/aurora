package main

import (
	"aurora/anet"
)

//Server 模块的测试函数
func main() {

	/*
		服务端测试
	*/
	//1 创建一个server 句柄 s
	// s := anet.NewServer("[zinx V0.2]")

	s := anet.NewServer()

	//2 开启服务
	s.Serve()
}
