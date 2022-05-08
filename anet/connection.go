package anet

import (
	"aurora/aiface"
	"errors"
	"fmt"
	"io"
	"net"
)

/*
连接层的抽象是把链接和数据以及处理接口封装在一起
*/

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool

	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler aiface.IMsgHandle

	// 无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte

	//给缓冲队列发送数据的channel，
	// 如果向缓冲队列发送数据，那么把数据发送到这个channel下
	//	SendBuffChan chan []byte

}

//初始化连接模块的方法
func NewConntion(conn *net.TCPConn, connID uint32, msgHandler aiface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
		//		SendBuffChan: make(chan []byte, 512),
	}

	return c
}

/*
	StartWriter 写消息Goroutine， 用户将数据发送给客户端
	读写分离的好处是，后续可以针对这两个go程扩展业务
	例如，涉及到数据库操作，读写分离还是很有必要的
*/
func (c *Connection) StartWriter() {

	defer fmt.Println(c.RemoteAddr().String(), " conn Writer exit!")
	defer c.Stop()

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			//conn已经关闭
			return
		}
	}
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}

		msg.SetData(data)
		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		//根据绑定好的msgid找到对应的api业务执行
		go c.MsgHandler.DoMsgHandler(&req)
	}

}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	//开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()

	// 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//for {
	//	select {
	//	case <- c.ExitBuffChan:
	//		//得到退出消息，不再阻塞
	//		return
	//	}
	//}

}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	// 告知writer关闭
	c.ExitBuffChan <- true

	//通知从缓冲队列读数据的业务，该链接已经关闭
	//c.ExitBuffChan <- true

	//关闭该链接全部管道 回收资源
	close(c.ExitBuffChan)
	close(c.msgChan)
	//close(c.SendBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//将数据发送给缓冲队列，通过专门从缓冲队列读数据的go写给客户端
func (c *Connection) SendBuff(data []byte) error {
	return nil
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	/*
		旧版本框架中，读写放在一起处理了，读到数据之后在同一协程中写回数据给客户端。
		这样的扩展性是较差的，例如想要加入消息队列，就很不方便了；
		为了实现解耦，可以采用读写分离的方案来操作。
	*/
	//写回客户端(旧版本)
	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Write msg id ", msgId, " error ")
	//	c.ExitBuffChan <- true
	//	return errors.New("conn Write error")
	//}

	//将数据写入管道，发送给写go程
	c.msgChan <- binaryMsg

	return nil
}
