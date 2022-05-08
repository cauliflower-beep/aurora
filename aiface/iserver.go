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
}
