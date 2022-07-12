package etcd_go

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	etcdClient *clientv3.Client
)

// 初始化etcd客户端链接
func init() {
	var err error
	config := clientv3.Config{
		Endpoints:   []string{"xxx.xxx.xxx.xxx:2379"},
		DialTimeout: 5 * time.Second,
	}

	etcdClient, err = clientv3.New(config)
	if err != nil {
		panic(err.Error())
	}
}
