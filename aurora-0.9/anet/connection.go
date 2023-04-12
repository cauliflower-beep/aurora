package anet

import (
	"aurora-v0.9/aiface"
	"aurora-v0.9/utils"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	// TCPServer 字段是为了配合server中的connMgr字段，建立conn <-> server 相互索引的关系
	TcpServer aiface.IServer // 当前conn属于哪个server，在conn初始化的时候添加
	Conn      *net.TCPConn   //当前连接的socket TCP套接字
	ConnID    uint32         //当前连接的ID 也可以称作为SessionID，ID全局唯一
	isClosed  bool           //当前连接的关闭状态

	MsgMgr aiface.IMsgMgr //消息id和对应处理方法的管理模块

	ExitBuffChan chan bool //告知该连接已经退出|停止的channel

	msgChan     chan []byte //无缓冲管道 用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte //有缓冲通道 用于读、写两个goroutine之间的消息通信
}

// NewConnection 新建连接
func NewConnection(server aiface.IServer, conn *net.TCPConn, connID uint32, msgMgr aiface.IMsgMgr) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgMgr:       msgMgr,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte), //msgChan初始化
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MasMsgChanLen),
	}
	c.TcpServer.GetConnMgr().Add(c) //将新创建的链接添加到链接管理模块中
	return c
}

// StartWriter
// @Description: 写消息goroutine 服务器将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[write goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err, " conn writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send data error:", err, " conn writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}

		case <-c.ExitBuffChan:
			//conn已经关闭
			return
		}
	}
}

// StartReader 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("[reader goroutine is running!]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn reader exit!]")
	defer c.Stop()

	for {
		// 创建拆包解包对象
		dp := NewDataPack()

		// 读取客户端的msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			c.ExitBuffChan <- true
			//continue
			return
		}

		// 拆包，得到msgId 和 dataLen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			c.ExitBuffChan <- true
			//continue
			return
		}

		// 根据 dataLen 读取data，放在 msg.Data 中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				c.ExitBuffChan <- true
				//continue
				return
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的request数据
		req := Request{
			iConn: c,
			msg:   msg, // 之前的buf 改成msg
		}

		// 这里并没有强制使用多任务worker机制
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经启动工作池，消息交给worker处理
			c.MsgMgr.SendMsgToTaskQueue(&req)
		} else {
			//调用当前连接业务，这里执行的是当前conn绑定的 Router
			go c.MsgMgr.DoMsgHandler(&req)
		}
	}
}

// Start 启动连接 让当前连接开始工作
func (c *Connection) Start() {
	// 开启用户从客户端读取数据流程的goroutine
	go c.StartReader()
	// 开启用于写回客户端数据流程的goroutine
	go c.StartWriter()

	//============================
	//按照用户传递进来的创建连接时需要处理的业务方法，执行hook函数
	//放在这里比较好，因为连接已经开始工作了
	c.TcpServer.CallOnConnStart(c)
	//============================

	for {
		select {
		case <-c.ExitBuffChan:
			return //得到退出消息，不再阻塞
		}
	}
}

// Stop 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	// 如果当前连接已经关闭，可以直接返回
	if c.isClosed {
		return
	}
	c.isClosed = true

	//=============================
	// 如果用户注册了该链接的关闭回调业务，在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)
	//=============================

	_ = c.Conn.Close() //关闭tcp连接

	c.ExitBuffChan <- true             //通知从缓冲队列读取数据的业务，该连接已经关闭
	c.TcpServer.GetConnMgr().Remove(c) // 将当前连接从连接管理器中删除
	close(c.ExitBuffChan)              //关闭该连接全部管道
	close(c.msgBuffChan)
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetTCPConnection 获取当前连接原始的socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将msg转发给Tcp客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection is closed when sendMsg")
	}
	//将data封包并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	//写回客户端
	c.msgChan <- msg //将之前直接回写给conn.write的方法改为发送给channel 供write读取

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed when sendBuffMsg")
	}
	//将data封包并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	//写回客户端
	c.msgBuffChan <- msg
	return nil
}