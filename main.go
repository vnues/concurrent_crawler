package main

import (
	"./engine"
	"./scheduler"
	"./zhenai/parse"
	"fmt"
)

func main(){
	fmt.Println("开始启动项目")
	e:= engine.ConcurrentEngine{
		//给它传QueuedScheduler 意思是QueuedScheduler有实现这个方法
		//我感觉跟实例化对象一样
		Scheduler:&scheduler.QueuedScheduler{},
		WorkerCount:100,
	}
	e.Run(engine.Request{
		Url:"http://www.zhenai.com/zhenghun",
		ParseFunc:parse.ParseCityList,
	})
}

// channel可以说是全局变量的--你可以这样理解 引用传递可以这样用
//我们不想让多个worker抢一个in 管道
//而是多个worker都有自己的in管道

//s.workerChan <- w
//s.workerChan是一个管道类型 然后存入的变量又是管道类型
//这个w 就是in

//其实worker队列就是in队列   --好的告辞

//管道是放值的  是我们要的结果 传递的变量

//这里的worker队列意思是把worker放进去队列执行（我们通过in管道去控制worker达到worker队列
//如何把worker放进去队列
// 意思我这种把worker放进去队列是通过in放进去队列执行的 因为多个worker对应多个in管道类型 -- 相当于映射

//最终的结果我们是要把request放入in