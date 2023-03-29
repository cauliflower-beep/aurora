package main

import "aurora/anet"

/*
基于aurora框架开发的服务器端应用程序
 */

func main(){
	//创建一个server句柄，使用框架中的api
	s := anet.NewServer("[aurorav0.1]")
	//启动server
	s.Serve()
}
