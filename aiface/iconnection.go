package aiface

import (
	"context"
	"net"
)

// IConnection
//  @Description: 全部连接相关接口声明
type IConnection interface {
	Start()                   //启动连接，让当前连接开始工作
	Stop()                    //停止连接，结束当前连接状态
	Context() context.Context //返回ctx，用于用户自定义的go程获取连接退出状态

	GetConnection() net.Conn //从当前连接获取原始的socket TCPConn
	GetConnID() uint64       //获取当前连接ID
	RemoteAddr() net.Addr    //获取链接远程地址信息
	LocalAddr() net.Addr     //获取链接本地地址信息

	Send(data []byte) error
	SendToQueue(data []byte) error
	SendMsg(msgID uint32, data []byte) error     //直接将Message数据发送给远程的TCP客户端(无缓冲)
	SendBuffMsg(msgID uint32, data []byte) error //直接将Message数据发送给远程的TCP客户端(有缓冲)

	SetProperty(key string, value interface{})   //设置链接属性
	GetProperty(key string) (interface{}, error) //获取链接属性
	RemoveProperty(key string)                   //移除链接属性
	IsAlive() bool                               //判断当前连接是否存活
	SetHeartBeat(checker IHeartbeatChecker)      //设置心跳检测器
}
