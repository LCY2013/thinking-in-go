package etcd_usage

import (
	"context"
	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		t.Error(err)
		return
	}

	_ = client
}

// TestEtcdPut etcd put
func (s *etcdTestSuite) TestEtcdPut() {
	var (
		kv          clientv3.KV
		putResponse *clientv3.PutResponse
		err         error
	)

	// 用于读写etcd健值对
	kv = clientv3.NewKV(s.client)

	if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hello"); err != nil {
		//if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/job2", "world"); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("Revision: %v", putResponse.Header.Revision)
	if putResponse.PrevKv != nil {
		s.T().Logf("PrevKv: %s", putResponse.PrevKv.Value)
	}
}

// TestEtcdGet etcd get
func (s *etcdTestSuite) TestEtcdGet() {
	var (
		kv          clientv3.KV
		getResponse *clientv3.GetResponse
		err         error
	)

	// 用于读写etcd健值对
	kv = clientv3.NewKV(s.client)

	if getResponse, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("Kvs: %v", getResponse.Kvs)
	s.T().Logf("Revision: %v \n", getResponse.Header.Revision)

	s.T().Log("----------------------------")
	// 只返回数量
	if getResponse, err = kv.Get(context.TODO(), "/cron/jobs/job1", clientv3.WithCountOnly()); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("With Count: %v \n", getResponse.Count)

	s.T().Log("----------------------------")
	// 返回某个前缀的数据
	if getResponse, err = kv.Get(context.TODO(), "/cron/jobs", clientv3.WithPrefix()); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("With Prefix Kvs: %v \n", getResponse.Kvs)

	s.T().Log("----------------------------")
	// 返回某个前缀的数据
	if getResponse, err = kv.Get(context.TODO(), "/cron/jobs", clientv3.WithFromKey()); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("With From Kvs: %v \n", getResponse.Kvs)
}

// TestEtcdGetPrefix etcd get prefix
func (s *etcdTestSuite) TestEtcdGetPrefix() {
	var (
		kv          clientv3.KV
		getResponse *clientv3.GetResponse
		err         error
	)

	// 用于读写etcd健值对
	kv = clientv3.NewKV(s.client)

	// 返回某个前缀的数据
	if getResponse, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("With Prefix Kvs: %v \n", getResponse.Kvs)
	s.T().Log("----------------------------")
}

// TestEtcdDelete etcd删除机制
func (s *etcdTestSuite) TestEtcdDelete() {
	var (
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		err     error
		idx     int
		kvpair  *mvccpb.KeyValue
	)

	// 用于读写etcd健值对
	kv = clientv3.NewKV(s.client)

	// 删除指定某个KV
	/*if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV()); err != nil {
		s.T().Error(err)
		return
	}
	if len(delResp.PrevKvs) > 0 {
		for idx, kvpair = range delResp.PrevKvs {
			s.T().Logf("Delete [%d] kv: key-%s , value-%s \n", idx, kvpair.Key, kvpair.Value)
		}
	}*/

	// 删除从某个key开始的若干个key
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithFromKey(), clientv3.WithPrevKV()); err != nil {
		s.T().Error(err)
		return
	}
	if len(delResp.PrevKvs) > 0 {
		for idx, kvpair = range delResp.PrevKvs {
			s.T().Logf("Delete [%d] kv: key-%s , value-%s \n", idx, kvpair.Key, kvpair.Value)
		}
	}
}

// TestEtcdLease etcd 租约机制测试
func (s *etcdTestSuite) TestEtcdLease() {
	var (
		kv                     clientv3.KV
		err                    error
		lease                  clientv3.Lease
		leaseId                clientv3.LeaseID
		leaseGrantResp         *clientv3.LeaseGrantResponse
		putResp                *clientv3.PutResponse
		getResp                *clientv3.GetResponse
		keepRespChan           <-chan *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveResponse *clientv3.LeaseKeepAliveResponse
	)

	// 申请一个租约(lease)
	lease = clientv3.NewLease(s.client)

	// 申请一个10s的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		s.T().Error(err)
		return
	}

	// 拿到租约ID
	leaseId = leaseGrantResp.ID

	// 5秒后取消自动续租
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// 总共续租有15s的生命周期

	// 自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		s.T().Error(err)
		return
	}

	// 处理续租应答
	go func() {
		for {
			select {
			case leaseKeepAliveResponse = <-keepRespChan:
				if leaseKeepAliveResponse == nil {
					s.T().Log("租约失效")
					goto END
				}
				// 每秒会续租一次，所以就会收到一次应答
				s.T().Logf("收到自动应答: %d", leaseKeepAliveResponse.ID)
			}
		}
	END:
	}()

	// 获取kv API子集
	kv = clientv3.NewKV(s.client)

	// PUT 一个kv， 与租约关联，然后才能实现10s后自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		s.T().Error(err)
		return
	}

	s.T().Log("lock successful: ", putResp.Header.Revision)

	// 定时查看key是否过期
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			s.T().Error(err)
			return
		}
		if getResp.Count == 0 {
			s.T().Logf("%s 已过期", "/cron/lock/job1")
			break
		}
		time.Sleep(2 * time.Second)
		s.T().Logf("未过期: %v", getResp.Kvs)
	}
}

// TestEtcdLease etcd 租约机制测试
func (s *etcdTestSuite) TestEtcdWatch() {
	var (
		kv                 clientv3.KV
		putResp            *clientv3.PutResponse
		delResp            *clientv3.DeleteResponse
		getResp            *clientv3.GetResponse
		err                error
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchChan          clientv3.WatchChan
		watchResponse      clientv3.WatchResponse
		event              *clientv3.Event
	)

	// 获取kv API子集
	kv = clientv3.NewKV(s.client)

	// 模拟etcd中kv的变化
	go func() {
		for {
			if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job8", "i am job8"); err != nil {
				s.T().Error(err)
				return
			}
			_ = putResp

			time.Sleep(1 * time.Second)

			if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job8"); err != nil {
				s.T().Error(err)
				return
			}
			_ = delResp

			time.Sleep(1 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job8"); err != nil {
		s.T().Error(err)
		return
	}

	// 现在key是否存在
	if len(getResp.Kvs) != 0 {
		s.T().Logf("当前值: %s", getResp.Kvs[0].Value)
	}

	// 当前etcd集群事物ID，单调递增
	watchStartRevision = getResp.Header.Revision + 1

	// 创建一个watcher
	watcher = clientv3.NewWatcher(s.client)

	// 启动监听
	s.T().Logf("从该版本向后监听: %d", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, cancelFunc)
	watchChan = watcher.Watch(ctx, "/cron/jobs/job8", clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResponse = range watchChan {
		for _, event = range watchResponse.Events {
			switch event.Type {
			case mvccpb.PUT:
				s.T().Logf("创建KEY: %s, 修改前: %s, 修改后: %s, 创建Revision: %d, 修改Revision: %d", event.Kv.Key, event.PrevKv, event.Kv.Value, event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				s.T().Logf("删除KEY: %s, 创建Revision: %d, 修改Revision: %d", event.Kv.Key, event.Kv.CreateRevision, event.Kv.ModRevision)
			}
		}
	}
}

// TestEtcdLease etcd operation
func (s *etcdTestSuite) TestEtcdOperation() {
	var (
		kv     clientv3.KV
		opResp clientv3.OpResponse
		err    error
		putOp  clientv3.Op
		getOp  clientv3.Op
	)

	// 获取kv API子集
	kv = clientv3.NewKV(s.client)

	// 创建Op operation
	putOp = clientv3.OpPut("/cron/jobs/job8", "")

	// 执行OP
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("put Revision: %v", opResp.Put().Header.Revision)

	// 创建GET OP
	getOp = clientv3.OpGet("/cron/jobs/job8")

	// 执行OP
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("数据Revision: %v", opResp.Get().Header.Revision) // 未删除时，第一次 执行 create rev == mod rev
	s.T().Logf("数据value: %v", opResp.Get().Kvs)

	// kv.Do(op)

	// kv.Put
	// kv.Get
	// kv.Delete
}

// TestEtcdLease etcd optimistic lock(乐观锁)
func (s *etcdTestSuite) TestEtcdOptimisticLock() {
	var (
		kv                     clientv3.KV
		lease                  clientv3.Lease
		leaseGrantResp         *clientv3.LeaseGrantResponse
		leaseId                clientv3.LeaseID
		ctx                    context.Context
		cancelFunc             context.CancelFunc
		keepRespChan           <-chan *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveResponse *clientv3.LeaseKeepAliveResponse
		txn                    clientv3.Txn
		txnResp                *clientv3.TxnResponse
		genUuid                uuid.UUID

		err error
	)

	// lease 实现锁自动过期
	// op操作
	// txn事物：if else then

	// 1、上锁（创建租约，自动续租，拿着租约去抢占一个key）
	// 申请一个租约(lease)
	lease = clientv3.NewLease(s.client)

	// 申请一个5s的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		s.T().Error(err)
		return
	}

	// 拿到租约ID
	leaseId = leaseGrantResp.ID

	// 准备一个取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.Background())

	// 确保函数退出后，自动续租会停止
	defer cancelFunc()
	defer func(lease clientv3.Lease, ctx context.Context, id clientv3.LeaseID) {
		_, err = lease.Revoke(ctx, id)
		if err != nil {
			s.T().Error(err)
		}
	}(lease, context.TODO(), leaseId)

	// 自动续租，直到主动取消或者5s后过期
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		s.T().Error(err)
		return
	}

	// 处理续租应答
	go func() {
		for {
			select {
			case leaseKeepAliveResponse = <-keepRespChan:
				if leaseKeepAliveResponse == nil {
					goto END
				}
				// 每秒会续租一次，所以就会收到一次应答
				s.T().Logf("收到自动续租应答: %d", leaseKeepAliveResponse.ID)
			}
		}
	END:
	}()

	// if 不存在key， then设置它 else 抢锁失败
	// 获取kv API子集
	kv = clientv3.NewKV(s.client)

	// 创建事物
	txn = kv.Txn(context.TODO())

	// 定义事物
	genUuid = uuid.New()

	// if key 不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		// then try lock
		Then(clientv3.OpPut("/cron/lock/job9", genUuid.URN(), clientv3.WithLease(leaseId))).
		// else lock failed
		Else(clientv3.OpGet("/cron/lock/job9"))

	// 事物提交
	if txnResp, err = txn.Commit(); err != nil {
		s.T().Error(err)
		return
	}

	// 是否抢到了锁
	if !txnResp.Succeeded {
		s.T().Logf("failed to lock, already occupied by: %s", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}

	// 2、处理业务
	s.T().Log("handle something start")
	time.Sleep(5 * time.Second)
	s.T().Log("handle something end")

	// 3、释放锁（取消自动续租，释放租约）
	// defer 会把租约释放掉，关联的KV就会被删除了
}
