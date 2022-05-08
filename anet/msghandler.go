package anet

import (
	"aurora/aiface"
	"aurora/utils"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32]aiface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	// 负责worker池取任务的消息队列，与worker池中的worker数量应该是一一对应的
	TaskQueue []chan aiface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

//创建msghandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]aiface.IRouter),
		/*
			这里需要注意，TaskQueue 是每个 Queue队列中对应的request最大数量，所以要设置最大值
		*/
		TaskQueue:      make([]chan aiface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取 工作池中的worker通常是跟cpu核心数一一对应的，虽然go中的goroutine更加轻量级，这里还是保持传统
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request aiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router aiface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId, "succ!")
}

/*
StartWorkerPool
启动一个worker工作池,开启工作池的动作只能发生一次（一个aurora框架只能有一个worker工作池）
这个方法应该暴露，因为aurora需要在某个地方启动工作池
*/
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1.当前的worker对应的channel消息队列 开辟空间 第0个worker就用第0个channel...
		/*
			这里需要注意，所有worker共用一个消息队列也可以，需要加锁，竞争等待
			有时候没必要一对一，会造成浪费；
			但每个worker独享一个队列也有好处，可以减少冲突
		*/
		mh.TaskQueue[i] = make(chan aiface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)
		//2.启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan aiface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request aiface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		" request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
