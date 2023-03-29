/**
 * @author uuxia
 * @date 15:56 2023/3/10
 * @description 责任链模式
 **/

package acode

import "aurora/aiface"

type RealInterceptorChain struct {
	request      aiface.Request
	position     int
	interceptors []aiface.Interceptor
}

func (this *RealInterceptorChain) Request() aiface.Request {
	return this.request
}

func (this *RealInterceptorChain) Proceed(request aiface.Request) aiface.Response {
	if this.position < len(this.interceptors) {
		chain := NewRealInterceptorChain(this.interceptors, this.position+1, request)
		interceptor := this.interceptors[this.position]
		response := interceptor.Intercept(chain)
		return response
	}
	return request
}

func NewRealInterceptorChain(list []aiface.Interceptor, pos int, request aiface.Request) aiface.Chain {
	return &RealInterceptorChain{
		request:      request,
		position:     pos,
		interceptors: list,
	}
}
