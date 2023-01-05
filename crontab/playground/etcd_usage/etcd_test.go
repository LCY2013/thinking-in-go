package etcd_usage

import (
	"context"
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
