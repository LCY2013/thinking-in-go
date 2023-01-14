package etcd

import (
	"context"
	"fmt"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/lib/errors"
	"github.com/LCY2013/thinking-in-go/crontab/slave/configs"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"sync"
	"time"
)

// Register 注册信息
type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string // local IP
	port    int    // port
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
		if localIP, err = getLocalIP(); err != nil {
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
	return
}

// getLocalIP 获取本地网卡IP
func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)

	// 获取所有网卡信息
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	// 获取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址：ipv4，ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPv6
			if ipNet.IP.To4() == nil {
				continue
			}
			ipv4 = ipNet.IP.String() // xxxx.xxxx.xxxx.xxxx
			return
		}
	}

	err = errors.ERR_NO_LOCAL_IP_FOUND

	return
}

// KeepOnline worker 节点注册
func (r *Register) KeepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepAlive      <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		err            error
	)

	for {
		cancelFunc = nil

		// 注册路径
		regKey = fmt.Sprintf("%s%s", constants.JobWorkerRegisterDir, r.localIP)

		// 创建租约
		if leaseGrantResp, err = r.lease.Grant(context.TODO(), 10); err != nil {
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
