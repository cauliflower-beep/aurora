package main

import (
	"aurora/anet"
	"aurora/examples/zinx_interceptor/interceptors"
	"aurora/examples/zinx_interceptor/router"
)

func main() {
	// 创建server 对象
	server := anet.NewServer()
	// 添加路由映射
	server.AddRouter(1, &router.HelloRouter{})
	// 添加自定义拦截器
	server.AddInterceptor(&interceptors.MyInterceptor{})
	// 启动服务
	server.Serve()
}
