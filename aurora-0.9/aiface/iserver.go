package aiface

// IServer
//
//	@Description: 定义服务接口
type IServer interface {
	Start() //启动服务器方法
	Stop()  //停止服务器方法
	Serve() //开启业务服务方法

	AddRouter(msgId uint32, router IRouter) //路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
	GetConnMgr() IConnMgr                   //得到链接管理

	/*
		有时候，在创建链接之前、断开连接之后，
		我们希望执行一些用户自定义的业务
		那就需要给aurora增添两个连接创建后和断开前时机的回调函数，一般也称作Hook(钩子)函数
	*/
	SetOnConnStart(func(conn IConnection)) //设置该server连接创建时的Hook函数
	SetOnConnStop(func(conn IConnection))  //设置该server连接断开时的Hook函数
	CallOnConnStart(conn IConnection)      //调用链接创建时的 Hook函数
	CallOnConnStop(conn IConnection)       //调用链接断开时的 Hook函数
}
