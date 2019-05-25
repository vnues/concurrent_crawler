## go并发版爬虫

### Scheduler实现一、架构图：

![](./1.png)



### worker: Fether和Parser，把worker并发
```go
func worker(r Request) (ParseResult,error){
	body,err :=fetch.Fetcher(r.Url)
	fmt.Printf("我正在爬取这个地址Fetching %s\n",r.Url)
	if err!=nil{
		return ParseResult{},err
	}
	ParseResult :=r.ParseFunc(body)
	return ParseResult,nil
}


//开启一个协程
func createWorker(in chan Request,out chan ParseResult){
	go func(){
		 for{
			 request := <-in
			 parseResult,err:=worker(request)
			 if err!=nil{
                 continue
			 }
			 out<-parseResult
		 }
	}()

}
```

- Scheduler: work并发之后，会面临多对多的并发任务的分配，有很对的request,很多worker在等着做它们，Scheduler去分配这些任务。

- Scheduler实现一：Scheduler收到一个个request,所有worker公用一个输入(in)，所有worker在**同一个channel里面去抢下一个request**。谁抢到谁做，这种不行，存在等待问题，解决如Scheduler实现二


### Scheduler实现二、架构图：

![](./2.png)
```go
//engine
//这里的goroutine听过是main函数里，因为main函数也是一个goroutine,而run函数不是goroutine
	 func (e *ConcurrentEngine) Run(seeds ...Request){
            ...     	
     	   for{
     		   result :=<-out
     		   for _,request:=range result.Requests{
     			   e.Scheduler.Submit(request)
     		   }
     		   //打印返回来的parseResult
     		   //打印返回来的parseResult
     		   for _,item :=range result.Items{
     		   	fmt.Printf("got item: %d %v\n",itemCount,item)
     			   itemCount++
     		   }
     
     	   }
     }
```
> result :=<-out 和  e.Scheduler.Submit(request) 在一个协程里，你要等到in发送成功out才能执行不然就会死锁，存在等待


### Submit(将种子输入到in)
```go
package scheduler

import (
	"../engine"
)



//实现interface接口

type  SimpleScheduler struct {
	//其实就是in管道
   workChan chan engine.Request
   //也就是当需要改变SimpleScheduler的属性的时候就需要指针传递
}
//要求引用者是指针类型并且符合SimpleScheduler类型
//为什么传过去要指针类型 因为我们的workChan需要被操作 我们不希望拷贝的workChan来操作
func (s *SimpleScheduler) Submit(r engine.Request){
         //这里就是将Request种子放进in管道
         //但是你有没有看到in管道变量是engine下的 不是这个包的
         //所以我们得把in给拿过来
	go func(){
		s.workChan <-r
	}()
}

func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request){
	s.workChan=c
}
```

- channel的通信起码需要两个goroutine才能进行，如果两个channel放在一个goroutine那么就存在等待问题，一个死锁或者卡死另外一个很明显不能继续运行了

- 假如你要输入数据到channel，那么最好就是提前通知接收者我准备发送数据
```go
   //开启多个createWorker--嗷嗷待哺要从in管道拿到种子
	   for i :=0;i<e.WorkerCount;i++{
		//createWorker需要从in管道拿到种子然后解析以后把结果放进去out管道
		createWorker(in,out)
	    }
	    //要把种子送机skeduler队列的方法--种子是通过in channel去拿的
	    //channel:这段放后面是提前跟别人说你要在这里等着拿数据
	   for _,r:=range seeds{
		e.Scheduler.Submit(r)
	    }
```