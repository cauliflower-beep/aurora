package aiface

// IDataPack
//  @Description: 消息的打包和解包方法
//  @Description: 直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
type IDataPack interface {
	GetHeadLen() uint32                //获取包头长度方法
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (IMessage, error)   //拆包方法
}

const (
	// ZinxDataPack Zinx标准封包和拆包方式
	ZinxDataPack string = "zinx_pack"

	//...(+)
	//自定义封包方式在此添加
)

const (
	// ZinxMessage Zinx默认标准报文协议格式
	ZinxMessage string = "zinx_message"
)
