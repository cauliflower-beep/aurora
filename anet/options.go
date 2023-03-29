package anet

import "aurora/aiface"

// Option Server的服务Option
type Option func(s *Server)

// WithPacket 只要实现Packet 接口可自由实现数据包解析格式，如果没有则使用默认解析格式
func WithPacket(pack aiface.IDataPack) Option {
	return func(s *Server) {
		s.SetPacket(pack)
	}
}

// ClientOption Client的客户端Option
type ClientOption func(c *Client)

// WithPacketClient Client的客户端Option
func WithPacketClient(pack aiface.IDataPack) ClientOption {
	return func(c *Client) {
		c.SetPacket(pack)
	}
}
