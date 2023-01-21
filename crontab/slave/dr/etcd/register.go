package etcd

import (
	"context"
	"fmt"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/slave/configs"
	"github.com/LCY2013/thinking-in-go/crontab/tools"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"sync"
	"time"
)

// Register 注册信息
type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string // local IP
	success bool   // success

}

var (
	GRegister *Register
	// MgrOnce 控制并发
	registerOnce = sync.Once{}
)

// InitRegister 初始化服务注册器
func InitRegister() (err error) {
	registerOnce.Do(func() {
		var (
			config  clientv3.Config
			client  *clientv3.Client
			localIP string
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

		// 获取本机IP
		if localIP, err = tools.GetLocalIP(); err != nil {
			return
		}

		// 得到KV和Lease的API子集
		GRegister = &Register{
			client:  client,
			kv:      clientv3.NewKV(client),
			lease:   clientv3.NewLease(client),
			localIP: localIP,
		}
	})

	// 注册工作节点
	async.GO(func() {
		GRegister.KeepOnline()
	})

	for !GRegister.success {
		time.Sleep(50 * time.Millisecond)
	}
	return
}

// KeepOnline worker 节点注册
func (r *Register) KeepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		getResp        *clientv3.GetResponse
		keepAlive      <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		err            error
	)

	for {
		cancelFunc = nil

		// 注册路径
		regKey = fmt.Sprintf("%s/%s", constants.JobWorkerRegisterDir, r.localIP)

		// 如果是IP:PORT类型
		if configs.Conf().Consistent.Hash.Type == constants.IpPort {
			// 注册工作节点信息
			if configs.Conf().Serves == nil || len(configs.Conf().Serves) == 0 {
				panic("serves config not found")
			}
			regKey = fmt.Sprintf("%s%s:%d", constants.JobWorkerRegisterDir, r.localIP, configs.Conf().Serves[0].ServePort)
		}

		// 先判断是否已经注册, 解决同样的key被使用后，一个移除导致都不可能
		if getResp, err = r.kv.Get(context.TODO(), regKey); err != nil {
			goto RETRY
		}
		if getResp.Count > 0 {
			log.WithFields(log.Fields{
				"register": "已经注册",
				"regKey":   regKey,
			}).Info("KeepOnline")
			os.Exit(-1)
		}

		// 创建租约，这里可以做成配置
		if leaseGrantResp, err = r.lease.Grant(context.TODO(), 3); err != nil {
			//if leaseGrantResp, err = r.lease.Grant(context.TODO(), 20); err != nil {
			goto RETRY
		}

		// 自动续租
		if keepAlive, err = r.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		// 注册到etcd
		if _, err = r.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}

		// 处理续租应答
		for {
			select {
			case keepAliveResp = <-keepAlive:
				// 成功就后续执行
				r.success = true
				if keepAliveResp == nil {
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}
}
