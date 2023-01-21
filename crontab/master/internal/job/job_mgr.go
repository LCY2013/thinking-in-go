package service

import (
	"context"
	"encoding/json"
	"fmt"
	jobEntity "github.com/LCY2013/thinking-in-go/crontab/domain"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
	consistent "github.com/LCY2013/thinking-in-go/crontab/tools/consistenthash"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

// Mgr 任务管理器
type Mgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher

	startMaster chan struct{}

	workerChangedChan chan *jobEntity.WorkerChangeEvent
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
			client:  client,
			kv:      clientv3.NewKV(client),
			lease:   clientv3.NewLease(client),
			watcher: clientv3.NewWatcher(client),

			startMaster:       make(chan struct{}, 1),
			workerChangedChan: make(chan *jobEntity.WorkerChangeEvent, 10),
		}
	})

	async.GO(func() {
		// 启动全量处理
		err = G_MGR.WatchWholeJobs(context.TODO())
		if err != nil {
			log.Error(err)
		}
	})

	// 启动工作节点变更事件监听处理
	G_MGR.handleWorkerChangeEvent()
	return
}

// SaveJob 保存job信息
func (mgr *Mgr) SaveJob(ctx context.Context, job *jobEntity.JobEntity) (oldJob *jobEntity.JobEntity, err error) {
	// 把任务保存到/cron/jobs/任务名称 -> json
	var (
		jobKey   string
		jobValue []byte
		putResp  *clientv3.PutResponse
		preJob   jobEntity.JobEntity
		member   consistent.Member
	)

	// 没有工作节点也不用一直等待。。。
	if GWorkerNode.GetMembers() != nil && len(GWorkerNode.GetMembers()) > 0 {
		// 选择对应的nodeIP
		member = GWorkerNode.LocateKey([]byte(job.Name))

		if job.NodeIp == "" {
			job.NodeIp = jobEntity.ExtractWorkerIP(member.String())
		}
		job.OldNodeIp = job.NodeIp
		job.NodeIp = jobEntity.ExtractWorkerIP(member.String())
	}

	// etcd 的保存key
	jobKey = fmt.Sprintf("%s%s", constants.JobDir, job.Name)

	// 序列化任务信息
	if jobValue, err = json.Marshal(job); err != nil {
		return nil, err
	}

	// 保存到etcd中，并且获取以前的值信息
	if putResp, err = mgr.kv.Put(ctx, jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return nil, err
	}

	// 如果时更新，那么返回新值
	if putResp.PrevKv == nil {
		return
	}

	// 反序列化到老值上面
	_ = json.Unmarshal(putResp.PrevKv.Value, &preJob)

	oldJob = &preJob

	return
}

// DeleteJob 删除job信息
func (mgr *Mgr) DeleteJob(ctx context.Context, jobName string) (oldJob *jobEntity.JobEntity, err error) {
	// 把任务保存到/cron/jobs/任务名称 -> json
	var (
		jobKey  string
		delResp *clientv3.DeleteResponse
		preJob  jobEntity.JobEntity
	)

	// etcd 的保存key
	jobKey = fmt.Sprintf("%s%s", constants.JobDir, jobName)

	// 保存到etcd中，并且获取以前的值信息
	if delResp, err = mgr.kv.Delete(ctx, jobKey, clientv3.WithPrevKV()); err != nil {
		return nil, err
	}

	// 如果时更新，那么返回新值
	if delResp.PrevKvs == nil || len(delResp.PrevKvs) == 0 {
		return
	}

	// 反序列化到老值上面
	_ = json.Unmarshal(delResp.PrevKvs[0].Value, &preJob)

	oldJob = &preJob

	return
}

// ListJob 列举所有job信息
func (mgr *Mgr) ListJob(ctx context.Context) (jobList []*jobEntity.JobEntity, err error) {
	var (
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		dirKey  string
		job     *jobEntity.JobEntity
	)

	// 任务根目录
	dirKey = constants.JobDir

	// 获取目录下所有任务信息
	if getResp, err = mgr.kv.Get(ctx, dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// 初始化数组空间
	jobList = make([]*jobEntity.JobEntity, 0)

	// 遍历所有任务进行反序列化
	for _, kvPair = range getResp.Kvs {
		job = &jobEntity.JobEntity{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}

	return
}

// KillJob kill job
// ./etcdctl watch "/cron/killer/" --prefix
func (mgr *Mgr) KillJob(ctx context.Context, jobName string) (err error) {
	// 更新一些key=/cron/killer/任务名
	var (
		killerKey      string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)

	// 通知worker节点kill 对应的任务
	killerKey = fmt.Sprintf("%s%s", constants.JobKillDir, jobName)

	// 让worker监听一次put操作，创建一个租约让其稍后自动过期
	if leaseGrantResp, err = mgr.lease.Grant(ctx, 1); err != nil {
		return
	}

	// 租约ID
	leaseId = leaseGrantResp.ID

	// 设置killer标记
	if _, err = mgr.kv.Put(ctx, killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}

	return
}

// WatchWholeJobs 监听全部的任务变化
func (mgr *Mgr) WatchWholeJobs(ctx context.Context) (err error) {
	// master开始干活
	<-mgr.startMaster

	var (
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *jobEntity.JobEntity
		//watchStartRevision int64
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		member     consistent.Member
		jobKey     string
		jobValue   []byte
	)

	// worker节点是否已经注册到指定数量
	for {
		if getResp, err = mgr.kv.Get(ctx, constants.JobWorkerRegisterDir, clientv3.WithPrefix(), clientv3.WithCountOnly()); err != nil {
			log.Error(err)
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if getResp.Count < int64(configs.Conf().Consistent.Hash.WorkerNodeNum) {
			log.WithFields(log.Fields{
				"WatchWholeJobs": "stating ...",
			}).Info()
			time.Sleep(500 * time.Millisecond)
			continue
		}

		break
	}

	// get /cron/jobs/目录下所有的任务，并且获知当前集群的revision
	if getResp, err = mgr.kv.Get(ctx, constants.JobDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前有那些任务
	for _, kvPair = range getResp.Kvs {
		// 分配任务到对应的任务节点

		// 反序列化json得到job
		if job, err = jobEntity.UnpackJobEntity(kvPair.Value); err != nil {
			log.WithFields(log.Fields{
				"WatchJobs": "err",
			}).Error(err)
			continue
		}

		// 重新计算工作节点
		if member = GWorkerNode.LocateKey([]byte(job.Name)); member == nil {
			continue
		}

		job.OldNodeIp = job.NodeIp
		job.NodeIp = member.String()

		// 序列化任务信息
		if jobValue, err = json.Marshal(job); err != nil {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs": "err",
			}).Error(err)
			continue
		}

		if _, err = mgr.kv.Put(context.TODO(), string(kvPair.Key), string(jobValue)); err != nil {
			log.WithFields(log.Fields{
				"WatchWholeJobs": "err",
			}).Error(err)
			continue
		}

		// 将任务加入到对应节点下
		// etcd 的保存分发任务到对应的节点的key
		jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, jobEntity.ExtractWorkerIP(member.String()), job.Name)

		// 保存到etcd中，并且获取以前的值信息
		if _, err = mgr.kv.Put(ctx, jobKey, string(jobValue)); err != nil {
			continue
		}
	}

	// 从该revision向后监听变化事件
	async.GO(func() {
		// 监听协程
		// 从GET时刻的后续版本开始监听
		//watchStartRevision = getResp.Header.Revision + 1
		// 监听/cron/jobs/目录的后续变化
		watchChan = mgr.watcher.Watch(context.Background(), constants.JobDir, clientv3.WithPrefix(), clientv3.WithPrevKV())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					//反序列化Job，
					if job, err = jobEntity.UnpackJobEntity(watchEvent.Kv.Value); err != nil {
						continue
					}

					// 重新计算工作节点
					if member = GWorkerNode.LocateKey([]byte(job.Name)); member == nil {
						continue
					}

					job.OldNodeIp = job.NodeIp
					job.NodeIp = member.String()

					// 序列化任务信息
					if jobValue, err = json.Marshal(job); err != nil {
						log.WithFields(log.Fields{
							"SchedulerOldWorkerNodeJobs": "err",
						}).Error(err)
						continue
					}

					// 将任务加入到对应节点下
					// etcd 的保存分发任务到对应的节点的key
					jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, jobEntity.ExtractWorkerIP(member.String()), job.Name)

					// 保存到etcd中，并且获取以前的值信息
					if _, err = mgr.kv.Put(ctx, jobKey, string(jobValue)); err != nil {
						continue
					}

				case mvccpb.DELETE: // 任务删除事件
					// delete /cron/jobs/job0

					//反序列化Job，
					if job, err = jobEntity.UnpackJobEntity(watchEvent.PrevKv.Value); err != nil {
						continue
					}

					// etcd 的保存分发任务到对应的节点的key
					jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, job.NodeIp, job.Name)

					// 删除某个node的节点任务信息
					if _, err = mgr.kv.Delete(ctx, jobKey); err != nil {
						continue
					}

				}
			}
		}
	})

	return
}

// handleWorkerChangeEvent 处理工作节点变更事件
func (mgr *Mgr) handleWorkerChangeEvent() {
	var (
		changeEvent *jobEntity.WorkerChangeEvent
	)

	async.GO(func() {
		for changeEvent = range mgr.workerChangedChan {
			if changeEvent == nil {
				continue
			}

			switch changeEvent.ChangeType {
			case constants.WorkerEventAdd:
				mgr.SchedulerNewWorkerNodeJobs(changeEvent)

			case constants.WorkerEventDelete:
				mgr.SchedulerOldWorkerNodeJobs(changeEvent)
			}
		}
	})
}

// SchedulerNewWorkerNodeJobs 调度计算后的任务到新节点上去
func (mgr *Mgr) SchedulerNewWorkerNodeJobs(changeEvent *jobEntity.WorkerChangeEvent) {
	// 计算所有节点的任务信息，先删除现在在老节点但是本来分配在新节点的任务，然后分配该任务到新节点
	var (
		getResp  *clientv3.GetResponse
		kvPair   *mvccpb.KeyValue
		job      *jobEntity.JobEntity
		member   consistent.Member
		jobValue []byte
		jobKey   string
		err      error
	)
	// get /cron/jobs/目录下所有的任务，并且获知当前集群的revision
	if getResp, err = mgr.kv.Get(context.TODO(), constants.JobDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前有那些任务
	for _, kvPair = range getResp.Kvs {
		// 分配任务到对应的任务节点

		// 反序列化json得到job
		if job, err = jobEntity.UnpackJobEntity(kvPair.Value); err != nil {
			log.WithFields(log.Fields{
				"SchedulerNewWorkerNodeJobs": "err",
			}).Error(err)
			continue
		}

		// 重新计算工作节点
		if member = GWorkerNode.LocateKey([]byte(job.Name)); member == nil {
			continue
		}

		job.OldNodeIp = job.NodeIp
		job.NodeIp = member.String()

		// 序列化任务信息
		if jobValue, err = json.Marshal(job); err != nil {
			log.WithFields(log.Fields{
				"SchedulerNewWorkerNodeJobs": "err",
			}).Error(err)
			continue
		}

		// 将任务加入到对应节点下
		// etcd 的保存分发任务到对应的节点的key
		jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, job.NodeIp, job.Name)

		// 保存到etcd中，并且获取以前的值信息
		if _, err = mgr.kv.Put(context.TODO(), jobKey, string(jobValue)); err != nil {
			log.WithFields(log.Fields{
				"SchedulerNewWorkerNodeJobs": "err",
			}).Error(err)
			continue
		}
	}
}

// SchedulerOldWorkerNodeJobs 调度老节点的任务到其他节点
func (mgr *Mgr) SchedulerOldWorkerNodeJobs(changeEvent *jobEntity.WorkerChangeEvent) {
	//分配删除节点上所有的任务到其他工作节点上， 然后删除该节点
	var (
		workerJobsKey string
		getResp       *clientv3.GetResponse
		kvpair        *mvccpb.KeyValue
		job           *jobEntity.JobEntity
		member        consistent.Member
		jobKey        string
		jobValue      []byte
		err           error
	)

	// 获取子任务路径
	workerJobsKey = fmt.Sprintf("%s%s", constants.WorkerJobs, changeEvent.WorkerName)

	if getResp, err = mgr.kv.Get(context.TODO(), workerJobsKey, clientv3.WithPrefix()); err != nil {
		log.WithFields(log.Fields{
			"SchedulerOldWorkerNodeJobs":           "err",
			"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
		}).Error(err)
		return
	}

	// 遍历
	for _, kvpair = range getResp.Kvs {
		//反序列化Job，
		if job, err = jobEntity.UnpackJobEntity(kvpair.Value); err != nil {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs":           "err",
				"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
			}).Error(err)
			continue
		}

		if GWorkerNode.GetMembers() == nil || len(GWorkerNode.GetMembers()) == 0 {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs":           "最终关闭的工作节点信息",
				"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
			}).Info()
			break
		}

		// 重新计算工作节点
		if member = GWorkerNode.LocateKey([]byte(job.Name)); member == nil {
			continue
		}
		job.OldNodeIp = job.NodeIp
		job.NodeIp = member.String()

		// 序列化任务信息
		if jobValue, err = json.Marshal(job); err != nil {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs":           "err",
				"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
			}).Error(err)
			continue
		}

		jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, job.NodeIp, job.Name)
		if _, err = mgr.kv.Put(context.TODO(), jobKey, string(jobValue)); err != nil {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs":           "err",
				"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
			}).Error(err)
			continue
		}
	}

	// 删除该节点
	if _, err = mgr.kv.Delete(context.TODO(), workerJobsKey, clientv3.WithPrefix()); err != nil {
		log.WithFields(log.Fields{
			"SchedulerOldWorkerNodeJobs":           "err",
			"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
		}).Error(err)
		return
	}
	log.WithFields(log.Fields{
		"SchedulerOldWorkerNodeJobs":           "success",
		"SchedulerOldWorkerNodeJobsWorkerNode": workerJobsKey,
	}).Info()
}

// PushWorkerChangeEvent 推送worker变更事件
func (mgr *Mgr) PushWorkerChangeEvent(event *jobEntity.WorkerChangeEvent) {
	mgr.workerChangedChan <- event
}

// StartMaster 启动master执行
func (mgr *Mgr) StartMaster() {
	if mgr.startMaster != nil {
		mgr.startMaster <- struct{}{}
		mgr.startMaster = nil
	}
}
