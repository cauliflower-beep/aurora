package aiface

type Inotify interface {
	HasIdConn(id uint64) bool                                        // 是否有这个id
	ConnNums() int                                                   // 存储的map长度
	SetNotifyID(Id uint64, conn IConnection)                         // 添加链接
	GetNotifyByID(Id uint64) (IConnection, error)                    // 得到某个链接
	DelNotifyByID(Id uint64)                                         // 删除某个链接
	NotifyToConnByID(Id uint64, MsgId uint32, data []byte) error     // 通知某个id的方法
	NotifyAll(MsgId uint32, data []byte) error                       // 通知所有人
	NotifyBuffToConnByID(Id uint64, MsgId uint32, data []byte) error // 通过缓冲队列通知某个id的方法
	NotifyBuffAll(MsgId uint32, data []byte) error                   // 缓冲队列通知所有人
}
