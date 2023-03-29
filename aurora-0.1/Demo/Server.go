package main

import "aurora/anet"

func main() {
	s := anet.NewServer("[aurora v0.1]") //创建一个server句柄
	s.Serve()                            //启动server
}
