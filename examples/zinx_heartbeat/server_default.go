package main

import (
	"aurora/anet"
	"time"
)

func main() {
	s := anet.NewServer()

	//启动心跳检测
	s.StartHeartBeat(5 * time.Second)

	s.Serve()
}
