package main

import (
	"aurora/aiface"
	"aurora/alog"
	"aurora/anet"
	"fmt"
	"time"
)

type TestRouter struct {
	anet.BaseRouter
}

// PreHandle -
func (t *TestRouter) PreHandle(req aiface.IRequest) {
	//使用场景模拟  完整路由计时
	start := time.Now()

	fmt.Println("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Println("该路由组执行完成耗时：", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req aiface.IRequest) {
	fmt.Println("--> Call Handle")

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		fmt.Println(err)
	}
}

// PostHandle -
func (t *TestRouter) PostHandle(req aiface.IRequest) {
	fmt.Println("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := anet.NewServer()
	s.AddRouter(1, &TestRouter{})
	alog.SetLogger(new(MyLogger))
	s.Serve()
}
