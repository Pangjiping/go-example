package limiter

import (
	"fmt"
	"time"
)

// Task 每个请求到来，需要把执行的业务逻辑封装成Task，放入漏桶，等待worker取出执行
type Task struct {
	handler func() Result // worker从漏桶取出请求对象后要执行的业务逻辑函数
	resChan chan Result   // 等待worker执行并返回结果的channel
	taskID  int
}

// Result 封装业务逻辑的执行结果
type Result struct{}

// handler 模拟封装业务逻辑的函数
func handler() Result {
	time.Sleep(300 * time.Millisecond)
	return Result{}
}

func NewTask(id int) Task {
	return Task{
		handler: handler,
		resChan: make(chan Result),
		taskID:  id,
	}
}

// 漏桶的具体实现
type LeakyBucket struct {
	BucketSize int       // 漏桶大小
	NumWorker  int       // 同时从漏桶中获取任务执行的worker数量
	bucket     chan Task // 存放任务的漏桶
}

func NewLeakyBucket(bucketSize int, numWorker int) *LeakyBucket {
	return &LeakyBucket{
		BucketSize: bucketSize,
		NumWorker:  numWorker,
		bucket:     make(chan Task, bucketSize),
	}
}

func (b *LeakyBucket) validate(task Task) bool {
	// 如果漏桶容量达到上限，返回false
	select {
	case b.bucket <- task:
	default:
		fmt.Printf("request[id=%d] is refused!\n", task.taskID)
		return false
	}

	// 等待worker执行
	<-task.resChan
	fmt.Printf("request[id=%d] is running!\n", task.taskID)
	return true
}

func (b *LeakyBucket) Start() {
	// 开启worker从漏桶中获取任务并执行
	go func() {
		for i := 0; i < b.NumWorker; i++ {
			go func() {
				for {
					task := <-b.bucket
					result := task.handler()
					task.resChan <- result
				}
			}()
		}
	}()
}
