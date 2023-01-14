package worker

import (
	"context"
	"github.com/LCY2013/thinking-in-go/crontab/domain"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

// Mgr 任务管理器
type Mgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	// GWorkerMgr 单例
	GWorkerMgr *Mgr
	// MgrOnce 控制并发
	mgrOnce = sync.Once{}
)

// InitWorkerMgr 初始化工作节点管理器
func InitWorkerMgr() (err error) {
	mgrOnce.Do(func() {
		var (
			config clientv3.Config
			client *clientv3.Client
		)

		// 初始化配置
		config = clientv3.Config{
			Endpoints:   configs.Conf().Etcd.Server.Endpoints,                                     // 连接地址
			DialTimeout: time.Duration(configs.Conf().Etcd.Server.DialTimeout) * time.Millisecond, // 连接超时
		}

		// 建立连接
		if client, err = clientv3.New(config); err != nil {
			return
		}

		// 得到KV和Lease的API子集
		GWorkerMgr = &Mgr{
			client: client,
			kv:     clientv3.NewKV(client),
			lease:  clientv3.NewLease(client),
		}
	})

	return
}

// ListWorkers 获取在线worker列表
func (mgr *Mgr) ListWorkers() (workerArr []string, err error) {
	var (
		getResp *clientv3.GetResponse
		kv      *mvccpb.KeyValue
	)

	// 初始化结果集
	workerArr = make([]string, 0)

	// 获取目录下所有kv
	if getResp, err = mgr.kv.Get(context.TODO(), constants.JobWorkerRegisterDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 解析每隔节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key:/cron/worker/register/xxxx.xxxx.xxxx.xxxx
		workerArr = append(workerArr, domain.ExtractWorkerIP(string(kv.Key)))
	}

	return
}
