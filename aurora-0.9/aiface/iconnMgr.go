package aiface

/*
	链接管理模块
	限定链接个数，如果超过一定量的客户端个数，aurora为了保证后端的及时响应，拒绝链接请求
*/

type IConnMgr interface {
	Add(conn IConnection)                           //添加链接
	Remove(conn IConnection)                        //删除链接
	GetConnById(connID uint32) (IConnection, error) //利用ConnID获取链接
	Len() int                                       //获取当前连接数量
	ClearConn()                                     //删除并停止所有连接
}
