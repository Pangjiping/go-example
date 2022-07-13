package concurrency_limiter

import (
	"fmt"
	"math"
	"runtime"
)

// 利用无缓冲channel与任务发送/执行分离方式

func busi4(ch chan int) {

	for t := range ch {
		fmt.Println("go task = ", t, ", goroutine count = ", runtime.NumGoroutine())
		wg.Done()
	}
}

func sendTask(task int, ch chan int) {
	wg.Add(1)
	ch <- task
}

func ExecUnbufferedChan() {
	ch := make(chan int) //无buffer channel
	goCnt := 3           //启动goroutine的数量

	for i := 0; i < goCnt; i++ {
		//启动go
		go busi4(ch)
	}

	taskCnt := math.MaxInt64 //模拟用户需求业务的数量
	for t := 0; t < taskCnt; t++ {
		//发送任务
		sendTask(t, ch)
	}

	wg.Wait()
}
