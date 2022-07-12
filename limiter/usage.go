package limiter

import (
	"fmt"
	"time"
)

// 滑动窗口
// func use() {
// 	limiter := NewSliding(100*time.Millisecond, time.Second, 10)

// 	for i := 0; i < 5; i++ {
// 		fmt.Println(limiter.IsLimited())
// 	}

// 	time.Sleep(100 * time.Millisecond)
// 	for i := 0; i < 5; i++ {
// 		fmt.Println(limiter.IsLimited())
// 	}

// 	// 这个请求触发限流
// 	fmt.Println(limiter.IsLimited())

// 	for _, v := range limiter.windows[limiter.getUidOrIp()] {
// 		fmt.Println(v.timestamp, v.count)
// 	}

// 	fmt.Println("one thousand years later ...")
// 	time.Sleep(time.Second)
// 	for i := 0; i < 7; i++ {
// 		fmt.Println(limiter.IsLimited())
// 	}
// 	for _, v := range limiter.windows[limiter.getUidOrIp()] {
// 		fmt.Println(v.timestamp, v.count)
// 	}
// }

// 漏桶
// func use() {
// 	bucket := NewLeakyBucket(10, 4)
// 	bucket.Start()

// 	var wg sync.WaitGroup
// 	for i := 0; i < 20; i++ {
// 		wg.Add(1)
// 		go func(id int) {
// 			defer wg.Done()
// 			task := NewTask(id)
// 			bucket.validate(task)
// 		}(i)
// 	}
// 	wg.Wait()
// }

// 令牌桶
func use() {
	tokenBucket := NewTokenBucket(5, 100*time.Millisecond)
	for i := 0; i < 6; i++ {
		fmt.Println(tokenBucket.IsLimited())
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println(tokenBucket.IsLimited())
}
