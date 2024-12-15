package main

import "fmt"

// Subject 定义了代理和真实对象需要实现的接口
type Subject interface {
	Request() string
}

// RealSubject 真实主题，执行实际操作
type RealSubject struct{}

func (r *RealSubject) Request() string {
	return "RealSubject: Handling request"
}

// Proxy 代理对象，封装了对 RealSubject 的访问
type Proxy struct {
	realSubject *RealSubject
}

func (p *Proxy) Request() string {
	// 在转发请求之前，可以执行一些额外的逻辑
	fmt.Println("Proxy: Before calling RealSubject")

	if p.realSubject == nil {
		p.realSubject = &RealSubject{}
	}

	// 调用真实对象的 Request 方法
	result := p.realSubject.Request()

	// 在转发请求之后，可以执行其他操作
	fmt.Println("Proxy: After calling RealSubject")

	return result
}
func main() {
	// 创建一个代理对象
	proxy := &Proxy{}

	// 通过代理来调用 Request 方法
	fmt.Println(proxy.Request())
}
