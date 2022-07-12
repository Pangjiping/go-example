package etcd_go

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// putToEtcd set key and value
func putToEtcd() error {
	// 用于写etcd的键值对
	kv := clientv3.NewKV(etcdClient)

	// PUT请求，clientv3.WithPrevKV()表示获取上一个版本的kv
	putResp, err := kv.Put(context.Background(), "/cron/jobs/job1", "hello", clientv3.WithPrevKV())
	if err != nil {
		return err
	}

	// 获取版本号
	fmt.Println("Revision: ", putResp.Header.Revision)

	// 如果有上一个key，输出kv的值
	if putResp.PrevKv != nil {
		fmt.Println("PrevValue: ", string(putResp.PrevKv.Value))
	}

	return nil
}

// getFromEtcd
func getFromEtcd() error {
	kv := clientv3.NewKV(etcdClient)

	// 简单的get操作
	getResp, err := kv.Get(context.Background(), "/cron/jobs/job1", clientv3.WithCountOnly())
	if err != nil {
		return err
	}
	fmt.Println(getResp.Count)

	// get with prefix
	getResp, err = kv.Get(context.Background(), "/cron/jobs", clientv3.WithPrefix())
	if err != nil {
		return err
	}
	fmt.Println(getResp.Kvs)

	return nil
}

// deleteFromEtcd
func deleteFromEtcd() error {
	kv := clientv3.NewKV(etcdClient)

	// 删除指定kv
	delResp, err := kv.Delete(context.Background(), "/cron/jobs/job1", clientv3.WithPrevKV())
	if err != nil {
		return err
	}

	// 被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair := range delResp.PrevKvs {
			fmt.Println("delete:", string(kvpair.Key), string(kvpair.Value))
		}
	}

	// 删除目录下的所有key
	delResp, err = kv.Delete(context.Background(), "/cron/jobs/", clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 删除从这个key开始的后面的两个key
	delResp, err = kv.Delete(context.Background(), "/cron/jobs/job1", clientv3.WithFromKey(), clientv3.WithLimit(2))
	if err != nil {
		return err
	}

	return nil
}

func watchEtcd() error {
	// 创建一个用于读写的kv
	kv := clientv3.NewKV(etcdClient)

	// 模拟etcd中kv的变化，每隔1s执行一次put-del操作
	go func() {
		for {
			kv.Put(context.Background(), "/cron/jobs/job7", "i am job7")
			kv.Delete(context.Background(), "/cron/jobs/job7")
			time.Sleep(time.Second * 1)
		}
	}()

	// 先get到当前的值，并监听后续变化
	getResp, err := kv.Get(context.Background(), "/cron/jobs/job7")
	if err != nil {
		return err
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：", string(getResp.Kvs[0].Value))
	}

	// 监听的revision起点
	watchStartRevision := getResp.Header.Revision + 1

	// 创建一个watcher
	watcher := clientv3.NewWatcher(etcdClient)

	// 启动监听
	fmt.Println("从这个版本开始监听：", watchStartRevision)

	// 设置5s的watch时间
	ctx, cancelFunc := context.WithCancel(context.Background())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	watchRespChan := watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	// 得到kv的变化事件，从chan中取值
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events { //.Events是一个切片
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value),
					"revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了：", "revision:", event.Kv.ModRevision)
			}
		}
	}

	return nil
}

func leaseEtcd() error {
	// 用于申请租约
	lease := clientv3.NewLease(etcdClient)

	// 申请一个10s的租约
	leaseGrantResp, err := lease.Grant(context.Background(), 10) //10s
	if err != nil {
		return err
	}

	// 拿到租约的id
	leaseID := leaseGrantResp.ID

	// 自动续租
	keepRespChan, err := lease.KeepAlive(context.Background(), leaseID)
	if err != nil {
		return err
	}

	// 处理续租应答的协程
	go func() {
		select {
		case keepResp := <-keepRespChan:
			if keepRespChan == nil {
				fmt.Println("lease has expired")
				goto END
			} else {
				// 每秒会续租一次
				fmt.Println("收到自动续租应答", keepResp.ID)
			}
		}
	END:
	}()

	// 用于读写etcd的键值对
	kv := clientv3.NewKV(etcdClient)

	// put一个key-value，关联租约，实现10s后过期
	// 防止程序宕机
	putResp, err := kv.Put(context.Background(), "/cron/lock/job1", "",
		clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}
	fmt.Println("put success", putResp.Header.Revision)

	for {
		getResp, err := kv.Get(context.Background(), "/cron/lock/job1")
		if err != nil {
			return err
		}
		if getResp.Count == 0 {
			return err
		} else {
			fmt.Println(getResp.Kvs)
			time.Sleep(2 * time.Second)
		}
	}
}

func operatorEtcd() error {
	kv := clientv3.NewKV(etcdClient)

	// 创建putop
	putOp := clientv3.OpPut("/cron/jobs/job7", "")

	// 执行op
	opResp, err := kv.Do(context.Background(), putOp)
	if err != nil {
		return err
	}

	fmt.Println("写入的revision：", opResp.Put().Header.Revision)

	// 创建getOp
	getOp := clientv3.OpGet("/cron/jobs/job7")

	// 执行op
	getResp, err := kv.Do(context.Background(), getOp)
	if err != nil {
		return err
	}
	fmt.Println("revision：", getResp.Get().Kvs[0].ModRevision)
	fmt.Println("取到的值为：", getResp.Get().Kvs[0].Value)

	return nil
}
