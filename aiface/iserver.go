/// Package aiface 主要提供aurora全部抽象层接口定义.
package aiface

/*
需要暴露给包外的方法，就放到接口里面去
*/
//定义服务接口
type IServer interface {
	Start() //启动服务器方法
	Stop()  //停止服务器方法
	Serve() //开启业务服务方法
	//AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId uint32, router IRouter)
	GetConnMgr() IConnMgr //获取该server的连接管理器
	/*
		注册和调用钩子函数也是为了系统的扩展性考虑。
		之前的架构中，连接建立之后就直接是客户端与服务器之间的request交互了，
		添加钩子函数之后可以做一些额外的事情，比如客户端连接进来，可以广播给其他玩家
	*/
	//注册OnConnStart 钩子函数的方法
	RegOnConnStart(func(conn IConnection))
	//注册OnConnStop 钩子函数的方法
	RegOnConnStop(func(conn IConnection))
	//调用OnConnStart 钩子函数的方法
	CallOnStart(conn IConnection)
	//调用OnConnStop 钩子函数的方法
	CallOnStop(conn IConnection)
}
