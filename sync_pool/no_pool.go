package sync_pool

import (
	"fmt"
	"sync"
	"time"
)

// TestClient 定义一个简单的客户端对象,使用用户uid初始化
type TestClient struct {
	Uid    string
	Action string
}

func (c *TestClient) DoAction() {
	if c.Action != "" {
		time.Sleep(200 * time.Millisecond) // 模拟业务场景,sleep 10ms
	}
}

func (c *TestClient) reset() {
	c.Uid = ""
	c.Action = ""
}

func (c *TestClient) register(uid string, action ...string) {
	c.Uid = uid
	if len(action) > 0 {
		c.Action = action[0]
	}
}

func NewTestClient(uid string, action ...string) *TestClient {
	client := &TestClient{
		Uid: uid,
	}
	if len(action) > 0 {
		client.Action = action[0]
	}
	return client
}

// MultiClientWithoutPool 在goroutine中不断的初始化client
func MultiClientWithoutPool(uid string) {
	start := time.Now().UnixMilli()
	action := "DoSomething"
	defer func(s int64) {
		fmt.Printf("MultiClientWithoutPool 用时%d\n", time.Now().UnixMilli()-s)
	}(start)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := NewTestClient(uid, action)
			client.DoAction()
		}()
	}
	wg.Wait()
}
