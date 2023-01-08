package job

import (
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
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
	// G_MGR 单例
	G_MGR *Mgr
	// MgrOnce 控制并发
	mgrOnce = sync.Once{}
)

// InitMgr 初始化管理器
func InitMgr() (err error) {
	mgrOnce.Do(func() {
		var (
			config clientv3.Config
			client *clientv3.Client
			lease  clientv3.Lease
			kv     clientv3.KV
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
		G_MGR = &Mgr{
			client: client,
			kv:     kv,
			lease:  lease,
		}
	})
	return
}
