package anet

import (
	"aurora/aiface"
	"errors"
	"fmt"
	"sync"
)

/*
连接管理模块
*/

type ConnMgr struct {
	connections map[uint32]aiface.IConnection //管理的连接集合
	/*
		这里需要加锁，因为涉及并发时对map的增删改查
		路由模块也有个map，但那里只是添加多路由，不加锁也ok
	*/
	connLock sync.RWMutex // 保护连接集合的读写锁
}

//创建连接管理模块
func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connections: make(map[uint32]aiface.IConnection),
		/*
			这里不需要对读写锁进行初始化，golang的读写锁是一个开箱即用的工具，
			只需要对它简单声明就可以使用
		*/
	}
}

// 添加链接
func (connMgr *ConnMgr) AddConn(conn aiface.IConnection) {
	// 保护共享资源map，加写锁
	/*
		对于同一个互斥锁的锁定操作和解锁操作总是应该成对出现的，这是一个惯用法，
		应该养成良好的习惯
	*/
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到connMgr中
	connMgr.connections[conn.GetConnID()] = conn
	//养成随手打日志的习惯，前期可能会麻烦一点，但是出了问题看的就比较清楚
	fmt.Println("connection add to ConnManager succ:conn num = ", connMgr.Len())
}

//删除链接
func (connMgr *ConnMgr) RemoveConn(conn aiface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("remove connection from ConnManager succ:conn num = ", connMgr.Len())

}

//根据connID获取连接
func (connMgr *ConnMgr) GetConn(connID uint32) (aiface.IConnection, error) {
	// 保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn not found")
	}

}

//得到当前连接总数
func (connMgr *ConnMgr) Len() int {
	return len(connMgr.connections)
}

//清除并终止所有连接
func (connMgr *ConnMgr) ClearConn() {
	// 保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("clear All conn succ! conn num = ", connMgr.Len())
}
