package concurrency_limiter

import (
	"fmt"
	"runtime"
	"sync"
)

// 只使用sync同步机制
// 没用,WaitGroup是一般用来阻塞线程的

var wg = sync.WaitGroup{}

func busi2(i int) {
	fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
	wg.Done()
}

func ExecWaitGroup() {
	task_cnt := 100
	for i := 0; i < task_cnt; i++ {
		wg.Add(1)
		go busi2(i)
	}

	wg.Wait()
}
