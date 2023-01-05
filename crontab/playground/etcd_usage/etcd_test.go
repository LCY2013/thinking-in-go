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

func (s *etcdTestSuite) TestEtcdPut() {
	var (
		kv          clientv3.KV
		putResponse *clientv3.PutResponse
		err         error
	)

	// 用于读写etcd健值对
	kv = clientv3.NewKV(s.client)

	//if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hello"); err != nil {
	if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/job2", "world"); err != nil {
		s.T().Error(err)
		return
	}
	s.T().Logf("Revision: %v", putResponse.Header.Revision)
	if putResponse.PrevKv != nil {
		s.T().Logf("PrevKv: %s", putResponse.PrevKv.Value)
	}
}

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
