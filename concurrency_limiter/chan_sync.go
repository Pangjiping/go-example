package concurrency_limiter

import (
	"fmt"
	"runtime"
)

// channel和sync组合使用

func busi3(ch chan bool, i int) {
	fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
	<-ch
	wg.Done()
}

func ExecChanWithSync() {
	//模拟用户需求go业务的数量
	task_cnt := 100

	ch := make(chan bool, 3)

	for i := 0; i < task_cnt; i++ {
		wg.Add(1)
		ch <- true
		go busi3(ch, i)
	}

	wg.Wait()
}
