package aiface

// IMsgMgr 消息管理抽象层
type IMsgMgr interface {
	DoMsgHandler(req IRequest)              //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑

	StartWorkerPool()                //启动worker工作池 go的调度算法已经做的很极致了 但大数量的goroutine依然会带来不必要的环境切换成本 这些成本应在服务器端节省掉 所以我们限定处理业务的goroutine数量
	SendMsgToTaskQueue(req IRequest) //将消息交给TaskQueue, 由worker进行处理
}
