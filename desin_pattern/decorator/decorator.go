package main

import (
	"fmt"
	"log"
	"time"
)

// 定义一个类型为函数的装饰器
type FunctionDecorator func(func() error) func() error

// 记录执行时间的装饰器
func LogExecutionTime(f func() error) func() error {
	return func() error {
		start := time.Now()
		err := f()
		log.Printf("Function took %s", time.Since(start))
		return err
	}
}

// 模拟业务函数
func BusinessLogic() error {
	time.Sleep(2 * time.Second)
	fmt.Println("Business logic executed")
	return nil
}

//func main() {
//	// 装饰 BusinessLogic
//	decoratedFunc := LogExecutionTime(BusinessLogic)
//	decoratedFunc()
//}
