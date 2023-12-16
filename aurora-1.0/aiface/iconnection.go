package aiface

import "net"

// IConnection 连接接口
type IConnection interface {
	Start()                         //启动连接，让当前连接开始工作
	Stop()                          //停止连接，结束当前连接状态
	GetConnID() uint32              //获取当前连接ID
	RemoteAddr() net.Addr           //获取远程客户端地址信息
	GetTCPConnection() *net.TCPConn //从当前连接获取原始的socket TCPConn

	SendMsg(msgId uint32, data []byte) error //新增方法 直接将Message数据发送给Tcp客户端 希望给用户返回一个TLV格式的数据

	SendBuffMsg(msgId uint32, data []byte) error //直接将Message数据发送给远程的Tcp客户端（有缓冲）

	/*
		使用链接处理的时候，希望和链接绑定一些用户的数据，或者参数
		所以这里添加了一些给链接设定属性参数的方法
		是不是莫名熟悉呢？在优品做智能预警的时候，也需要在请求参数中加上一些客户端标识，来区分ios或者安卓等，方便调试
	*/
	SetProperty(key string, value interface{})   //设置链接属性
	GetProperty(key string) (interface{}, error) //获取链接属性
	RemoveProperty(key string)                   //移除链接属性
}

// HandFunc
//
//	@Description: 统一处理连接业务的接口
//	@Description: 想要指定一个conn处理业务，只需要定义一个HandFunc，然后和该链接绑定即可
//	@param *net.TCPConn socket原生链接
//	@param []byte 客户端请求的数据
//	@param int 客户端请求的数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
