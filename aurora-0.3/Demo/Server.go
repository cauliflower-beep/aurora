package main

import "aurora-v0.3/anet"

func main() {
	s := anet.NewServer("[aurora v0.3]") //创建一个server句柄
	s.Serve()                            //启动server
}
