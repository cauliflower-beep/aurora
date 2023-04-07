package aiface

// IMsgMgr 消息管理抽象层
type IMsgMgr interface {
	DoMsgHandler(req IRequest)              //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
}
