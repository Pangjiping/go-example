package concurrency_limiter

import (
	"fmt"
	"runtime"
)

// 使用buffered channel来限制goroutine

func busi1(ch chan bool, i int) {
	fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
	<-ch
}

func ExecBufferedChan() {
	// 模拟用户需求业务的数量
	task_cnt := 100

	ch := make(chan bool, 5)

	for i := 0; i < task_cnt; i++ {
		ch <- true
		go busi1(ch, i)
	}
}
