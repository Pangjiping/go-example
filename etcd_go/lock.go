package etcd_go

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func lockByEtcd() {
	// 1. 上锁（创建租约，自动续租，拿着租约去抢占一个key ）
	// 用于申请租约
	lease := clientv3.NewLease(etcdClient)

	// 申请一个10s的租约
	leaseGrantResp, err := lease.Grant(context.Background(), 10) //10s
	if err != nil {
		fmt.Println(err)
		return
	}

	// 拿到租约的id
	leaseID := leaseGrantResp.ID

	// 准备一个用于取消续租的context
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 确保函数退出后，自动续租会停止
	defer cancelFunc()
	// 确保函数退出后，租约会失效
	defer lease.Revoke(context.Background(), leaseID)

	// 自动续租
	keepRespChan, err := lease.KeepAlive(ctx, leaseID)
	if err != nil {
		fmt.Println(err)
		return
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

	// if key 不存在，then设置它，else抢锁失败
	kv := clientv3.NewKV(etcdClient)
	// 创建事务
	txn := kv.Txn(context.Background())
	// 如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job7"), "=", 0)).
		Then(clientv3.OpPut("/cron/jobs/job7", "", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet("/cron/jobs/job7")) //如果key存在

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用了：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2. 处理业务（锁内，很安全）

	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)

	// 3. 释放锁（取消自动续租，释放租约）
	// defer会取消续租，释放锁
}
