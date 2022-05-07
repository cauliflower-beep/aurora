package aiface

import "net"

//定义连接模块的抽象层
type IConnection interface {
	//启动链接 让当前的连接准备开始工作
	Start()

	//停止链接 结束当前连接的工作
	Stop()

	//获取当前连接的绑定socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前连接模块的连接ID
	GetConnID() uint32

	//获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr

	//发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error
}

//定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
