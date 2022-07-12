package concurrency_limiter

import "sync"

// 并发执行，但是要限制并发数量
type concurrencyLimiter struct {
	runningNum  int32
	limit       int32
	blockingNum int32
	cond        *sync.Cond
	mu          *sync.Mutex
}

// NewConcurrencyLimiter 创建一个并发限制器，limit为限制数量，可以通过Reset()调整limit
// 每次调用Get()来获取并发权限来创建一个goroutine，完成任务后通过Release()释放资源
func NewConcurrencyLimiter(limit int32) *concurrencyLimiter {
	m := new(sync.Mutex)
	return &concurrencyLimiter{
		limit: limit,
		cond:  sync.NewCond(m),
		mu:    m,
	}
}

func (c *concurrencyLimiter) Reset(limit int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	tmp := c.limit
	c.limit = limit
	blockingNum := c.blockingNum

	// 优先唤醒阻塞的任务
	if limit-tmp > 0 && blockingNum > 0 {
		for i := int32(0); i < limit-tmp && blockingNum > 0; i++ {
			c.cond.Signal()
			blockingNum--
		}
	}
}

// Get 没有资源时会阻塞
func (c *concurrencyLimiter) Get() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.runningNum < c.limit {
		c.runningNum++
		return
	}

	// block here
	c.blockingNum++
	for !(c.runningNum < c.limit) {
		c.cond.Wait()
	}

	c.runningNum++
	c.blockingNum--
}

func (c *concurrencyLimiter) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.blockingNum > 0 {
		c.runningNum--
		c.cond.Signal()
		return
	}
	c.runningNum--
}
