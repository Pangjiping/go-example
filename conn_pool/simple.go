package conn_pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

// 实现一个简单的线程安全的连接池

var (
	ErrPoolClosed = errors.New("Connection Pool has Closed")
)

type SimplePool struct {
	m       sync.Mutex                // 保证线程安全
	res     chan io.Closer            // 连接存储的channel
	factory func() (io.Closer, error) // 新建连接的工厂方法
	closed  bool                      // 连接池是否关闭
}

func NewSimplePool(fn func() (io.Closer, error), size uint) (*SimplePool, error) {
	if size <= 0 {
		return nil, errors.New("Invalid.Size")
	}

	return &SimplePool{
		factory: fn,
		res:     make(chan io.Closer, size),
	}, nil
}

// Acquire 获取连接
func (p *SimplePool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.res:
		log.Println("Acquire: 共享资源")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquire: 新生成资源")
		return p.factory()
	}
}

// Release 释放连接
func (p *SimplePool) Release(r io.Closer) {
	// 保证线程安全
	p.m.Lock()
	defer p.m.Unlock()

	// 如果连接池关闭，直接释放即可
	if p.closed {
		r.Close()
		return
	}

	select {
	case p.res <- r:
		log.Println("连接正常释放")
	default:
		log.Println("连接池满了，直接关闭连接")
		r.Close()
	}
}

// Close 关闭连接池
func (p *SimplePool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	// 关闭通道
	close(p.res)

	// 关闭通道里的资源
	for r := range p.res {
		r.Close()
	}
}
