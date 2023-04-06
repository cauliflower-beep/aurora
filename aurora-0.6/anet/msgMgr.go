package anet

import (
	"aurora-v0.6/aiface"
	"fmt"
	"strconv"
)

type MsgMgr struct {
	Apis map[uint32]aiface.IRouter //存放每个MsgId对应的处理方法的map属性
}

func NewMsgMgr() *MsgMgr {
	return &MsgMgr{
		Apis: make(map[uint32]aiface.IRouter),
	}
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (mm *MsgMgr) DoMsgHandler(req aiface.IRequest) {
	handler, ok := mm.Apis[req.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", req.GetMsgId(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

// AddRouter 为消息添加具体的处理逻辑
func (mm *MsgMgr) AddRouter(msgId uint32, router aiface.IRouter) {
	// 判断当前的msg绑定的API处理方法是否已存在
	if _, ok := mm.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	// 添加msg与api的绑定关系
	mm.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}
