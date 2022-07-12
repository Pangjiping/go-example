package concurrency_limiter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	concurrency_num = int32(5)
)

func Test_InitLimiter(t *testing.T) {
	limiter := NewConcurrencyLimiter(concurrency_num)
	assert.NotNil(t, limiter)
	assert.NotNil(t, limiter.limit)
	assert.Equal(t, int32(5), limiter.limit)
	assert.NotNil(t, limiter.cond)
	assert.NotNil(t, limiter.mu)
	assert.Equal(t, int32(0), limiter.blockingNum)
	assert.Equal(t, int32(0), limiter.runningNum)
}

func Test_Limiter(t *testing.T) {
	limiter := NewConcurrencyLimiter(concurrency_num)

	// goroutine
	for i := 0; i < 10; i++ {
		go func(i int) {
			limiter.Get()
			t.Logf("[INFO] goroutine %d start", i)
			defer t.Logf("[INFO] goroutine %d finished", i)
			defer limiter.Release()

			time.Sleep(time.Second)
		}(i)
	}

	time.Sleep(10 * time.Second)

}
