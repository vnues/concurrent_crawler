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
		Scheduler:&scheduler.SimpleScheduler{},
		WorkerCount:100,
	}
	e.Run(engine.Request{
		Url:"http://www.zhenai.com/zhenghun",
		ParseFunc:parse.ParseCityList,
	})
}
