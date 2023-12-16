package aiface

/*
IRouter
路由接口，这里面路由是 框架用户 给该链接自定义的处理业务方法
路由里面的IRequest 则包含该链接的链接信息和客户端请求的数据信息
*/
type IRouter interface {
	PreHandle(req IRequest)  //处理conn业务之前的钩子方法(前置业务)
	Handle(req IRequest)     //处理conn业务的方法(主业务)
	PostHandle(req IRequest) //处理conn业务之后的钩子方法(后置业务)
}

/*
	router实际上的作用是，服务端应用可以给aurora框架配置当前链接的处理业务方法，
	之前的v0.2处理链接请求的方法是固定的，现在可以自定义，并且有3种接口可以重写。
*/
