package aiface

// IMsgHandle
//  @Description: 提供worker启动、处理消息业务调用等接口
//  @Description: 消息管理抽象层
type IMsgHandle interface {
	DoMsgHandler(request IRequest) //马上以非阻塞方式处理消息
	/*
		为消息添加具体的处理逻辑, msgID，支持整型，字符串
	*/
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()                    //启动worker工作池
	SendMsgToTaskQueue(request IRequest) //将消息交给TaskQueue,由worker进行处理

	Decode(request IRequest)
	AddInterceptor(interceptor Interceptor) //注册责任链任务入口，每个拦截器处理完后，数据都会传递至下一个拦截器，使得消息可以层层处理层层传递，顺序取决于注册顺序
}
