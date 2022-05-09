package aiface

/*
增加连接管理模块的目的是对连接进行更好的管理
例如超过一定数量的连接过来，为了保证服务器能够及时响应，需要拒绝这些连接
一台服务器最多能开启多少链接，是看linux内核的
*/

type IConnMgr interface {
	// 添加链接
	AddConn(conn IConnection)
	//删除链接
	RemoveConn(conn IConnection)
	//根据connID获取连接
	GetConn(connID uint32) (IConnection, error)
	//得到当前连接总数
	Len() int
	//清除并终止所有连接
	ClearConn()
}
