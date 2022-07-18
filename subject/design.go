package subject

import "sync"

type BasicInterface interface {
	NotifyAll()
	Register(next BasicInterface)
	Deregister(next BasicInterface)
}

// FirstPublisher 第一级发布者
// 从本地clusterInfo中获取信息
type FirstPublisher struct {
	Name            string
	Region          string
	Observers       map[string]BasicInterface // 下一级的观察者列表
	ClusterInfoChan chan string
	rwLock          sync.RWMutex
}

func NewFirstPublisher(name, region string, maxBuffer int) BasicInterface {
	return nil
}

// SecondPublisher 第二级发布者
// 从tag service中获取tags信息，对属性进行改写 [时间]
type SecondPublisher struct {
	name      string
	tagChan   chan string
	Observers map[string]BasicInterface
	rwLock    sync.RWMutex
}

func NewSecondPublisher(region string, maxBuffer int) BasicInterface {
	return nil
}
