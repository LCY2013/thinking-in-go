package service

import (
	"context"
	"fmt"
	entity "github.com/LCY2013/thinking-in-go/crontab/domain"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/lib/errors"
	"github.com/LCY2013/thinking-in-go/crontab/slave/configs"
	"github.com/LCY2013/thinking-in-go/crontab/tools"
	consistent "github.com/LCY2013/thinking-in-go/crontab/tools/consistenthash"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"hash/fnv"
	"sync"
	"time"
)

// HashNode 一致性hash节点名称
type HashNode struct {
	Name string `json:"name"`
}

func (h *HashNode) String() string {
	return h.Name
}

type hasher struct{}

func (hs hasher) Sum64(data []byte) uint64 {
	h := fnv.New64()
	_, err := h.Write(data)
	if err != nil {
		log.Error(err)
		return 0
	}
	return h.Sum64()
}

// WorkerNode 工作节点信息
type WorkerNode struct {
	consistent.Consistent

	client     *clientv3.Client
	kv         clientv3.KV
	lease      clientv3.Lease
	watcher    clientv3.Watcher
	cancelFunc context.CancelFunc // 取消上下文
	leaseId    clientv3.LeaseID   // 租约id

	isMaster bool
	masterIp string
	ip       string
}

var (
	// GWorkerNode 工作节点信息
	GWorkerNode *WorkerNode

	// MgrOnce 控制并发
	once = sync.Once{}
)

// InitWorkerNode 初始化一致性hash node信息
func InitWorkerNode() (err error) {
	once.Do(func() {
		var (
			config clientv3.Config
			client *clientv3.Client
			ipV4   string
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

		// 获取本地ip
		if ipV4, err = tools.GetLocalIP(); err != nil {
			return
		}

		if configs.Conf().Serves == nil || len(configs.Conf().Serves) == 0 {
			err = errors.ErrServesConfigNotFound
			return
		}

		if configs.Conf().Consistent.Hash.Type == constants.IpPort {
			ipV4 = fmt.Sprintf("%s:%d", ipV4, configs.Conf().Serves[0].ServePort)
		}

		// 得到KV和Lease的API子集
		GWorkerNode = &WorkerNode{
			client:  client,
			kv:      clientv3.NewKV(client),
			lease:   clientv3.NewLease(client),
			watcher: clientv3.NewWatcher(client),

			Consistent: *consistent.New(nil, consistent.Config{
				PartitionCount:    23,
				ReplicationFactor: 20,
				Load:              1.25,
				Hasher:            hasher{},
			}),

			ip: ipV4,
		}
	})

	// 选主
	GWorkerNode.confirmMaster()

	async.GO(func() {
		// 监听所有工作节点
		GWorkerNode.WatchWorkers(context.TODO())
	})

	// 监听所有工作节点
	//GWorkerNode.WatchWorkers(context.TODO())

	return
}

// WatchWorkers 监听工作节点变化信息，用于重新分配任务
func (w *WorkerNode) WatchWorkers(ctx context.Context) {
	var (
		getResp    *clientv3.GetResponse
		kvPair     *mvccpb.KeyValue
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		member     consistent.Member
		memberName string
		workerInfo string

		err error
	)

	for w.masterIp == "" {
		time.Sleep(50 * time.Millisecond)
	}

	// get /cron/worker/register/目录下所有的任务，并且获知当前集群的revision
	if getResp, err = w.kv.Get(ctx, constants.JobWorkerRegisterDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前有那些节点
	for _, kvPair = range getResp.Kvs {
		// 去掉工作节点前缀
		memberName = entity.ExtractWorkerIP(string(kvPair.Key))
		// 一致性hash节点添加
		w.Consistent.Add(&HashNode{
			Name: memberName,
		})
	}

	// 从该revision向后监听变化事件
	async.GO(func() {
		// 监听协程
		// 监听/cron/worker/register/目录的后续变化
		watchChan = w.watcher.Watch(context.Background(), constants.JobWorkerRegisterDir, clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 工作节点新增事件
					// 加入任务到该节点，一致性hash重新分配，任务重新分配，分配给新加入的节点 /cron/worker/jobs/新节点 任务
					memberName = string(watchEvent.Kv.Key)
					// 去掉工作节点前缀
					memberName = entity.ExtractWorkerIP(memberName)

					// 先判断该节点是否存在
					for _, member = range w.Consistent.GetMembers() {
						if memberName == member.String() {
							goto RETRY
						}
					}
					w.Consistent.Add(&HashNode{
						Name: memberName,
					})

					if w.isMaster {
						// 通知job分配调度重新实现，修改对应全部任务信息里面的nodeIP
						G_MGR.PushWorkerChangeEvent(&entity.WorkerChangeEvent{
							WorkerName: memberName,
							ChangeType: constants.WorkerEventAdd,
						})
					}
				RETRY:
					continue
				case mvccpb.DELETE: // 工作节点删除事件
					// 先移除本地成员，再分配该节点（/cron/worker/jobs/老节点）的任务，然后删除该节点数据
					memberName = string(watchEvent.Kv.Key)
					// 去掉工作节点前缀
					memberName = entity.ExtractWorkerIP(memberName)
					workerInfo = fmt.Sprintf("%s%s", constants.WorkerJobs, memberName)

					// 查询是否移除该节点
					if getResp, err = w.kv.Get(context.TODO(), workerInfo, clientv3.WithCountOnly()); err != nil {
						log.WithFields(log.Fields{
							"WatchWorkers":     err,
							"WatchWorkersInfo": workerInfo,
						}).Warn("移除工作节点")
						continue
					}
					if getResp.Count > 0 {
						log.WithFields(log.Fields{
							"WatchWorkers":     "",
							"WatchWorkersInfo": workerInfo,
						}).Warn("移除已存在线上的工作节点")
						continue
					}
					log.WithFields(log.Fields{
						"WatchWorkers":     "success",
						"WatchWorkersInfo": workerInfo,
					}).Info("移除工作节点")
					w.Consistent.Remove(memberName)

					if w.isMaster {
						// 通知任务重新分配，修改对应全部任务信息里面的nodeIP
						G_MGR.PushWorkerChangeEvent(&entity.WorkerChangeEvent{
							WorkerName: memberName,
							ChangeType: constants.WorkerEventDelete,
						})
					}
				}
			}
		}
	})
}

func (w *WorkerNode) LocateKey(key []byte) (member consistent.Member) {
	for w.Consistent.GetMembers() == nil || len(w.Consistent.GetMembers()) == 0 {
		log.WithFields(log.Fields{
			"LocateKey": "...",
		}).Info()
		time.Sleep(time.Millisecond * 100)
	}
	member = w.Consistent.LocateKey(key)
	log.WithFields(log.Fields{
		"LocateKey": "LocateKey",
	}).Info(member.String())
	return
}

// confirmMaster 选择主工作节点
func (w *WorkerNode) confirmMaster() {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		err            error
	)

	// 处理续租应答
	async.GO(func() {
		for {
			// 锁路径
			lockKey = constants.WorkerMasterDir

			if getResp, err = w.kv.Get(context.TODO(), lockKey); err != nil {
				return
			}
			if getResp.Kvs != nil && len(getResp.Kvs) > 0 && string(getResp.Kvs[0].Value) == w.ip {
				w.isMaster = true
				w.masterIp = string(getResp.Kvs[0].Value)
				return
			}
			if getResp.Kvs != nil && len(getResp.Kvs) > 0 && string(getResp.Kvs[0].Value) != w.ip {
				w.masterIp = string(getResp.Kvs[0].Value)
				goto RETRY
			}

			// 创建租约
			if leaseGrantResp, err = w.lease.Grant(context.TODO(), 5); err != nil {
				return
			}

			// context用于取消自动续租
			cancelCtx, cancelFunc = context.WithCancel(context.TODO())

			// 租约ID
			leaseId = leaseGrantResp.ID

			// 自动续租
			if keepRespChan, err = w.lease.KeepAlive(cancelCtx, leaseId); err != nil {
				goto RETRY
			}

			// 创建事物
			txn = w.kv.Txn(context.TODO())

			// 事物抢锁
			txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
				Then(clientv3.OpPut(lockKey, w.ip, clientv3.WithLease(leaseId))).
				Else(clientv3.OpGet(lockKey))

			// 提交事物
			if txnResp, err = txn.Commit(); err != nil {
				goto RETRY
			}

			// 成功返回，失败释放租约
			if !txnResp.Succeeded { // 锁被占用
				err = errors.ErrLockAlreadyRequired
				goto RETRY
			}

			// 锁成功
			w.leaseId = leaseId
			w.cancelFunc = cancelFunc
			w.isMaster = true
			w.masterIp = w.ip

			log.WithFields(log.Fields{
				"work-master": w.ip,
			}).Logf(log.InfoLevel, "confirmMaster")

			// job master 执行
			G_MGR.StartMaster()

			for {
				select {
				case keepResp = <-keepRespChan: // 自动续租的应答
					if keepResp == nil {
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
	})
}

/*// watchMaster 监听master变化事件
func (w *WorkerNode) watchMaster(ctx context.Context) {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
	)

	// 从该revision向后监听变化事件
	async.GO(func() {
		// 监听协程
		// 监听/cron/worker/register/目录的后续变化
		watchChan = w.watcher.Watch(context.Background(), constants.WorkerMasterDir)
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.DELETE: // master工作节点删除事件
					w.confirmMaster()
				}
			}
		}
	})
}*/
