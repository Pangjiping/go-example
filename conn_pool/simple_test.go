package conn_pool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"sync"
	"testing"
	"time"
)

type mockDB struct {
	name string
}

func (m *mockDB) Close() error {
	fmt.Printf("mockDB: %s, 连接正常关闭\n", m.name)
	return nil
}

var fn = func() (io.Closer, error) {
	return &mockDB{name: "new db"}, nil
}

func Test_SimplePool(t *testing.T) {
	pool, err := NewSimplePool(fn, 5)
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.Acquire()
			time.Sleep(time.Minute)
		}()
	}

	wg.Wait()
}
