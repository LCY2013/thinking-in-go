package service

import (
	"context"
	"fmt"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/lib/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// JobLock job分布式锁
type JobLock struct {
	kv    clientv3.KV
	lease clientv3.Lease

	jobName    string             // 任务名称
	cancelFunc context.CancelFunc // 取消上下文
	leaseId    clientv3.LeaseID   // 租约id

	isLocked bool
}

// InitJobLock 初始化分布式锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

// TryLock try lock
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)

	// 创建租约
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	// context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	// 租约ID
	leaseId = leaseGrantResp.ID

	// 自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	// 处理续租应答
	async.GO(func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan: // 自动续租的应答
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	})

	// 创建事物
	txn = jobLock.kv.Txn(context.TODO())
	// 锁路径
	lockKey = fmt.Sprintf("%s%s", constants.JobLockDir, jobLock.jobName)

	// 事物抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// 提交事物
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	// 成功返回，失败释放租约
	if !txnResp.Succeeded { // 锁被占用
		err = errors.ErrLockAlreadyRequired
		goto FAIL
	}

	// 锁成功
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	jobLock.isLocked = true
	return

FAIL:
	cancelFunc() // 取消自动续租
	jobLock.lease.Revoke(context.TODO(), leaseId)
	return
}

// UnLock un lock
func (jobLock *JobLock) UnLock() {
	if jobLock.isLocked {
		jobLock.cancelFunc() // 取消自动续租的协程
		_, _ = jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)
	}
}
