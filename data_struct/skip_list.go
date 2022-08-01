package data_struct

import (
	"math/rand"
	"sync"
)

type SkipNodeInt struct {
	key   int64
	value interface{}
	next  []*SkipNodeInt
}

type SkipListInt struct {
	SkipNodeInt
	mutex  sync.RWMutex
	update []*SkipNodeInt
	maxl   int
	skip   int
	level  int
	length int32
}

func NewSkipListInt(skip ...int) *SkipListInt {
	list := &SkipListInt{}
	list.maxl = 32
	list.skip = 4
	list.level = 0
	list.length = 0
	list.SkipNodeInt.next = make([]*SkipNodeInt, list.maxl)
	list.update = make([]*SkipNodeInt, list.maxl)

	if len(skip) == 1 && skip[0] > 1 {
		list.skip = skip[0]
	}
	return list
}

func (l *SkipListInt) Get(key int64) interface{} {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var prev = &l.SkipNodeInt
	var next *SkipNodeInt
	for i := l.level - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = next
			next = prev.next[i]
		}
	}

	if next != nil && next.key == key {
		return next.value
	} else {
		return nil
	}
}

func (l *SkipListInt) Set(key int64, val interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 获取每层的前驱节点 => l.update
	var prev = &l.SkipNodeInt
	var next *SkipNodeInt
	for i := l.level - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = next
			next = prev.next[i]
		}
		l.update[i] = prev
	}

	// 如果key已经存在
	if next != nil && next.key == key {
		next.value = val
		return
	}

	// 随机生成新节点的层数
	level := l.randomLevel()
	if level > l.level {
		level = l.level + 1
		l.level = level
		l.update[l.level-1] = &l.SkipNodeInt
	}

	// 申请新的节点
	node := &SkipNodeInt{}
	node.key = key
	node.value = val
	node.next = make([]*SkipNodeInt, level)

	// 调整next指向
	for i := 0; i < level; i++ {
		node.next[i] = l.update[i].next[i]
		l.update[i].next[i] = node
	}
	l.length++
}

func (l *SkipListInt) Delete(key int64) interface{} {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	//获取每层的前驱节点=>list.update
	var prev = &l.SkipNodeInt
	var next *SkipNodeInt
	for i := l.level - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = next
			next = prev.next[i]
		}
		l.update[i] = prev
	}

	//结点不存在
	node := next
	if next == nil || next.key != key {
		return nil
	}

	//调整next指向
	for i, v := range node.next {
		if l.update[i].next[i] == node {
			l.update[i].next[i] = v
			if l.SkipNodeInt.next[i] == nil {
				l.level -= 1
			}
		}
		l.update[i] = nil
	}

	l.length--
	return node.value
}

func (l *SkipListInt) Len() int32 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.length
}

func (l *SkipListInt) randomLevel() int {
	i := 1
	for ; i < l.maxl; i++ {
		if rand.Int()%l.skip != 0 {
			break
		}
	}
	return i
}
