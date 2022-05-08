/*
消息管理模块，主要实现两个功能：
1.根据消息id匹配不同的路由；
2.将多路由与对应的消息id保存在map中
*/

package aiface

type IMsgHandle interface {
	/*
		下面的非阻塞方式处理消息，是相对于 reader、writer两个go程来说的。
		V0.7版本的框架中，假设某一时刻有10w连接进来，那服务端就会启动30w个go程：
		10w reader 10w writer 10w DoMsgHandler
		读写go程是阻塞式的，有数据就读，没有数据就循环阻塞，并不会占用cpu资源；
		但处理消息的go程是非阻塞式的，会立即参与运算，占用cpu资源
		在10w个处理消息的go程之间切换，会浪费一定的系统资源
	*/
	DoMsgHandler(request IRequest)          //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)
}
