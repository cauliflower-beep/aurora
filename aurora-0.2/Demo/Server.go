package main

import "aurora-v0.2/anet"

func main() {
	s := anet.NewServer("[aurora v0.2]") //创建一个server句柄
	s.Serve()                            //启动server
}
