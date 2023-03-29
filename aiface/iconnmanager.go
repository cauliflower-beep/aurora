package aiface

// IConnManager
//  @Description: 连接管理抽象层
//  @Description: 连接管理相关,包括添加、删除、通过一个连接ID获得连接对象，当前连接数量、清空全部连接等方法
type IConnManager interface {
	Add(IConnection)                                                       //添加链接
	Remove(IConnection)                                                    //删除连接
	Get(uint64) (IConnection, error)                                       //利用ConnID获取链接
	Len() int                                                              //获取当前连接
	ClearConn()                                                            //删除并停止所有链接
	GetAllConnID() []uint64                                                //获取所有连接ID
	Range(func(uint64, IConnection, interface{}) error, interface{}) error //遍历所有连接
}
