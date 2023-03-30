package anet

import (
	"aurora-v0.2/aiface"
	"fmt"
	"net"
)

type Connection struct {
	Conn     *net.TCPConn //当前连接的socket TCP套接字
	ConnID   uint32       //当前连接的ID 也可以称作为SessionID，ID全局唯一
	isClosed bool         //当前连接的关闭状态

	handleAPI aiface.HandFunc //该连接的处理方法api

	ExitBuffChan chan bool //告知该连接已经退出|停止的channel
}

// NewConnection
//  @Description: 新建连接
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi aiface.HandFunc) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callbackApi,
		ExitBuffChan: make(chan bool, 1),
	}
}

// StartReader
//  @Description: 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running!")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err|", err)
			c.ExitBuffChan <- true
			continue
		}

		//调用当前连接业务，这里执行的是当前conn绑定的 HandFunc
		if err = c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("connID:%d handle is error! err:%v\n", c.ConnID, err)
			c.ExitBuffChan <- true
			return
		}
	}
}

// Start
//  @Description: 启动连接
func (c *Connection) Start() {
	go c.StartReader() //处理该连接读取到客户端数据之后的请求业务
	for {
		select {
		case <-c.ExitBuffChan:
			return //得到退出消息，不再阻塞
		}
	}
}

// Stop
//  @Description: 停止连接，结束当前连接状态
//  @receiver c
func (c *Connection) Stop() {
	// 如果当前连接已经关闭，可以直接返回
	if c.isClosed {
		return
	}
	c.isClosed = true

	// todo Connection Stop()如果用户注册了该链接的关闭回调业务，在此刻应该显示调用

	_ = c.Conn.Close() //关闭tcp连接

	c.ExitBuffChan <- true //通知从缓冲队列读取数据的业务，该连接已经关闭
	close(c.ExitBuffChan)  //关闭该连接全部管道
}

// GetConnID
//  @Description: 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetTCPConnection
//  @Description: 获取当前连接原始的socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// RemoteAddr
//  @Description:获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
