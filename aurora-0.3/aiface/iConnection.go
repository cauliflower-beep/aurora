package aiface

import "net"

// IConnection
//  @Description: 连接接口
type IConnection interface {
	Start()                         //启动连接，让当前连接开始工作
	Stop()                          //停止连接，结束当前连接状态
	GetConnID() uint32              //获取当前连接ID
	RemoteAddr() net.Addr           //获取远程客户端地址信息
	GetTCPConnection() *net.TCPConn //从当前连接获取原始的socket TCPConn
}

// HandFunc
//  @Description: 统一处理连接业务的接口
//  @Description: 想要指定一个conn处理业务，只需要定义一个HandFunc，然后和该链接绑定即可
//  @param *net.TCPConn socket原生链接
//  @param []byte 客户端请求的数据
//  @param int 客户端请求的数据长度
type HandFunc func(*net.TCPConn, []byte, int) error