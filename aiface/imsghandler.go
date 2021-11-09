/*
消息管理模块，主要实现两个功能：
1.根据消息id匹配不同的路由；
2.将多路由与对应的消息id保存在map中
 */

package aiface

type IMsgHandle interface{
	DoMsgHandler(request IRequest)			//马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter)	//为消息添加具体的处理逻辑
}