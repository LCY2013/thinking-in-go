package etcd_usage

import (
	"context"
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

func (s *etcdTestSuite) TestEtcdCURD() {
	var (
		kv          clientv3.KV
		putResponse *clientv3.PutResponse
		err         error
	)

	// 用于读写etcd健值对
	kv = clientv3.NewKV(s.client)

	if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hello"); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("Revision: %v", putResponse.Header.Revision)
	if putResponse.PrevKv != nil {
		s.T().Logf("PrevKv: %s", putResponse.PrevKv.Value)
	}
}
