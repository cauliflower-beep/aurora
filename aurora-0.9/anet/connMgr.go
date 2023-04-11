package anet

import (
	"aurora-v0.9/aiface"
	"errors"
	"fmt"
	"sync"
)

// ConnMgr
// @Description: 链接管理模块
type ConnMgr struct {
	connections map[uint32]aiface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

// NewConnMgr
// @Description: 创建一个链接管理
func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connections: make(map[uint32]aiface.IConnection),
	}
}

// Add
// @Description: 链接管理模块新增一个链接
func (cm *ConnMgr) Add(conn aiface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	connId := conn.GetConnID()
	cm.connections[connId] = conn
	fmt.Println("conn add to connMgr succ: conn num =  ", cm.Len())
}

// Remove
// @Description: 链接管理模块中删除一个链接
// @Description: 该方法只是单纯的将conn从map中摘掉
func (cm *ConnMgr) Remove(conn aiface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connMgr remove conn succ: connId = ", conn.GetConnID(), " conn num = ", cm.Len())
}

// GetConnById
// @Description: 根据connId获取链接
func (cm *ConnMgr) GetConnById(connID uint32) (aiface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn Not Found")
	}
}

// Len
// @Description: 获取当前链接数量
func (cm *ConnMgr) Len() int {
	return len(cm.connections)
}

// ClearConn
// @Description: 清除并停止所有链接
func (cm *ConnMgr) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connId, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connId)
	}
	fmt.Println("clear all conn succ: conn num = ", cm.Len())
}
