package anet

import (
	"aurora/aiface"
	"aurora/alog"
	"fmt"
	"time"
)

type HeartbeatChecker struct {
	interval time.Duration // 心跳检测时间间隔
	quitChan chan bool     // 退出信号

	makeMsg aiface.HeartBeatMsgFunc //用户自定义的心跳检测消息处理方法

	onRemoteNotAlive aiface.OnRemoteNotAlive //用户自定义的远程连接不存活时的处理方法

	msgID  uint32         // 心跳的消息ID
	router aiface.IRouter //用户自定义的心跳检测消息业务处理路由

	conn aiface.IConnection // 绑定的链接
}

// HeatBeatDefaultRouter
// @Description: 收到remote心跳消息的默认回调路由业务
type HeatBeatDefaultRouter struct {
	BaseRouter
}

// Handle -
func (r *HeatBeatDefaultRouter) Handle(req aiface.IRequest) {
	alog.Ins().InfoF("Recv Heartbeat from %s, MsgID = %+v, Data = %s",
		req.GetConnection().RemoteAddr(), req.GetMsgID(), string(req.GetData()))
}

// 默认的心跳消息生成函数
func makeDefaultMsg(conn aiface.IConnection) []byte {
	msg := fmt.Sprintf("heartbeat [%s->%s]", conn.LocalAddr(), conn.RemoteAddr())
	return []byte(msg)
}

// 默认的心跳检测函数
func notAliveDefaultFunc(conn aiface.IConnection) {
	alog.Ins().InfoF("Remote connection %s is not alive, stop it", conn.RemoteAddr())
	conn.Stop()
}

// NewHeartbeatChecker 创建心跳检测器
func NewHeartbeatChecker(interval time.Duration) aiface.IHeartbeatChecker {
	heartbeat := &HeartbeatChecker{
		interval: interval,
		quitChan: make(chan bool),

		//均使用默认的心跳消息生成函数和远程连接不存活时的处理方法
		makeMsg:          makeDefaultMsg,
		onRemoteNotAlive: notAliveDefaultFunc,
		msgID:            aiface.HeartBeatDefaultMsgID,
		router:           &HeatBeatDefaultRouter{},
	}

	return heartbeat
}

func (h *HeartbeatChecker) SetOnRemoteNotAlive(f aiface.OnRemoteNotAlive) {
	if f != nil {
		h.onRemoteNotAlive = f
	}
}

func (h *HeartbeatChecker) SetHeartbeatMsgFunc(f aiface.HeartBeatMsgFunc) {
	if f != nil {
		h.makeMsg = f
	}
}

func (h *HeartbeatChecker) BindRouter(msgID uint32, router aiface.IRouter) {
	if router != nil && msgID != aiface.HeartBeatDefaultMsgID {
		h.msgID = msgID
		h.router = router
	}
}

func (h *HeartbeatChecker) start() {
	ticker := time.NewTicker(h.interval)
	for {
		select {
		case <-ticker.C:
			h.check()
		case <-h.quitChan:
			ticker.Stop()
			return
		}
	}
}

// Start 启动心跳检测
func (h *HeartbeatChecker) Start() {
	go h.start()
}

// 停止心跳检测
func (h *HeartbeatChecker) Stop() {
	alog.Ins().InfoF("heartbeat checker stop, connID=%+v", h.conn.GetConnID())
	h.quitChan <- true
}

func (h *HeartbeatChecker) SendHeartBeatMsg() error {

	msg := h.makeMsg(h.conn)

	err := h.conn.SendMsg(h.msgID, msg)
	if err != nil {
		alog.Ins().ErrorF("send heartbeat msg error: %v, msgId=%+v msg=%+v", err, h.msgID, msg)
		return err
	}

	return nil
}

// 执行心跳检测
func (h *HeartbeatChecker) check() (err error) {

	if h.conn == nil {
		return nil
	}

	if !h.conn.IsAlive() {
		h.onRemoteNotAlive(h.conn)
	} else {
		err = h.SendHeartBeatMsg()
	}

	return err
}

// BindConn 绑定一个链接
func (h *HeartbeatChecker) BindConn(conn aiface.IConnection) {
	h.conn = conn
	conn.SetHeartBeat(h)
}

// CloneTo 克隆到一个指定的链接上
func (h *HeartbeatChecker) Clone() aiface.IHeartbeatChecker {

	heatbeat := &HeartbeatChecker{
		interval:         h.interval,
		quitChan:         make(chan bool),
		makeMsg:          h.makeMsg,
		onRemoteNotAlive: h.onRemoteNotAlive,
		msgID:            h.msgID,
		router:           h.router,
		conn:             nil, //绑定的链接需要重新赋值
	}

	return heatbeat
}

func (h *HeartbeatChecker) MsgID() uint32 {
	return h.msgID
}

func (h *HeartbeatChecker) Router() aiface.IRouter {
	return h.router
}