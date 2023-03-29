package interceptors

import (
	"aurora/aiface"
	"aurora/alog"
)

// 自定义拦截器1

type MyInterceptor struct {
}

func (m *MyInterceptor) Intercept(chain aiface.Chain) aiface.Response {
	request := chain.Request()
	// 这一层是自定义拦截器处理逻辑，这里只是简单打印输入
	iRequest := request.(aiface.IRequest)
	alog.Ins().InfoF("自定义拦截器, 收到消息：%s", iRequest.GetData())
	return chain.Proceed(chain.Request())
}
