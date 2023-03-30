package aiface

// IServer
//  @Description: 定义服务接口
type IServer interface {
	Start() //启动服务器方法
	Stop()  //停止服务器方法
	Serve() //开启业务服务方法
}
