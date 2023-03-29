/**
 * @author uuxia
 * @date 15:57 2023/3/10
 * @description 拦截器管理
 **/

package acode

import "aurora/aiface"

// InterceptorChain
// HTLV+CRC，H头码，T功能码，L数据长度，V数据内容
// +------+-------+---------+--------+--------+
// | 头码  | 功能码 | 数据长度 | 数据内容 | CRC校验 |
// | 1字节 | 1字节  | 1字节   | N字节   |  2字节  |
// +------+-------+---------+--------+--------+
type InterceptorChain struct {
	body       []aiface.Interceptor
	head, tail aiface.Interceptor
	request    aiface.Request
}

func NewInterceptorBuilder() aiface.InterceptorBuilder {
	return &InterceptorChain{
		body: make([]aiface.Interceptor, 0),
	}
}

func (this *InterceptorChain) Head(interceptor aiface.Interceptor) {
	this.head = interceptor
}

func (this *InterceptorChain) Tail(interceptor aiface.Interceptor) {
	this.tail = interceptor
}

func (this *InterceptorChain) AddInterceptor(interceptor aiface.Interceptor) {
	this.body = append(this.body, interceptor)
}

func (this *InterceptorChain) Execute(request aiface.Request) aiface.Response {
	this.request = request
	var interceptors []aiface.Interceptor
	if this.head != nil {
		interceptors = append(interceptors, this.head)
	}
	if len(this.body) > 0 {
		interceptors = append(interceptors, this.body...)
	}
	if this.tail != nil {
		interceptors = append(interceptors, this.tail)
	}
	chain := NewRealInterceptorChain(interceptors, 0, request)
	return chain.Proceed(this.request)
}
