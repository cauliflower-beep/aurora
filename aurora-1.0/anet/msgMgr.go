package anet

import (
	"aurora-v1.0/aiface"
	"aurora-v1.0/utils"
	"fmt"
	"strconv"
)

type MsgMgr struct {
	Apis map[uint32]aiface.IRouter //存放每个MsgId对应的处理方法的map属性

	WorkPoolSize uint32                 //业务工作Worker池的数量
	TaskQueue    []chan aiface.IRequest //消息队列 worker池中的goroutine从这个队列中取任务
}

func NewMsgMgr() *MsgMgr {
	return &MsgMgr{
		Apis: make(map[uint32]aiface.IRouter),

		// worker会从对应的队列中获取客户端的请求数据并且处理掉
		WorkPoolSize: utils.GlobalObject.WorkerPoolSize,
		// 一个worker对应一个queue 用来缓冲request请求信息，后续提供给worker池调用
		TaskQueue: make([]chan aiface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// StartOneWorker
// @Description: 启动一个worker工作流程
func (mm *MsgMgr) StartOneWorker(workerID int, taskQueue chan aiface.IRequest) {
	fmt.Println("worker Id = ", workerID, " is started.")

	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的req，并绑定执行的业务方法
		case req := <-taskQueue:
			mm.DoMsgHandler(req)
		}
	}
}

// StartWorkerPool
// @Description: 启动worker工作池
func (mm *MsgMgr) StartWorkerPool() {
	//遍历需要启动worker的数量，依次启动
	for i := 0; i < int(mm.WorkPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mm.TaskQueue[i] = make(chan aiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前worker，阻塞等待对应的任务队列是否有消息传递进来
		go mm.StartOneWorker(i, mm.TaskQueue[i])
	}
}

// SendMsgToTaskQueue
// @Description: 将消息交给taskQueue 由worker进行处理
func (mm *MsgMgr) SendMsgToTaskQueue(req aiface.IRequest) {
	//根据ConnId来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条链接的workerId
	workerId := req.GetConnection().GetConnID() % mm.WorkPoolSize
	fmt.Println("add connId = ", req.GetConnection().GetConnID(), " req msgId = ", req.GetMsgId(), "to workerId = ", workerId)
	//将请求消息发送给任务队列
	mm.TaskQueue[workerId] <- req
}
