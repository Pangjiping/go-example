package sync_pool

import (
	"fmt"
	"sync"
	"time"
)

func NewTestClientForPool() interface{} {
	return &TestClient{}
}

type TestClientPool struct {
	pool *sync.Pool
}

func (p *TestClientPool) DoAction(uid string, action ...string) {
	client := p.pool.Get().(*TestClient)
	defer func() {
		client.reset()
		p.pool.Put(client)
	}()

	client.register(uid, action...)
	client.DoAction()
}

func NewTestClientPool() *TestClientPool {
	return &TestClientPool{
		pool: &sync.Pool{
			New: NewTestClientForPool,
		},
	}
}

// MultiClientWithPool 使用链接池的概念,在池内Get连接
func MultiClientWithPool(uid string) {
	start := time.Now().UnixMilli()
	defer func(s int64) {
		fmt.Printf("MultiClientWithoutPool 用时%d\n", time.Now().UnixMilli()-s)
	}(start)

	action := "DoSomething"
	pool := NewTestClientPool()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.DoAction(uid, action)
		}()
	}
	wg.Wait()
}
