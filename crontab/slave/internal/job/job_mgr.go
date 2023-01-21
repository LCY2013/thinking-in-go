package service

import (
	"context"
	"encoding/json"
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
	"sync"
	"time"
)

// Mgr 任务管理器
type Mgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher

	workerNodeChangeEvent chan *entity.WorkerChangeEvent

	ip string
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
		G_MGR = &Mgr{
			client:                client,
			kv:                    client.KV,
			lease:                 client.Lease,
			watcher:               client.Watcher,
			ip:                    ipV4,
			workerNodeChangeEvent: make(chan *entity.WorkerChangeEvent, 5),
		}
	})

	// 启动任务监听
	err = G_MGR.WatchJobs(context.TODO())
	if err != nil {
		return
	}

	// 启动强杀kill任务通知
	G_MGR.watchKiller()

	// 处理新增节点，任务分配
	G_MGR.SchedulerAddWorkerNode()

	return
}

// WatchJobs 监听任务变化
func (mgr *Mgr) WatchJobs(ctx context.Context) (err error) {
	var (
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *entity.JobEntity
		//watchStartRevision int64
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobName    string
		jobEvent   *entity.JobEvent
		nodeJob    string
		nodeInfo   string
		localIP    string
		member     consistent.Member
		existent   bool
	)

	// 获取本机IP
	if localIP, err = tools.GetLocalIP(); err != nil {
		return
	}

	if configs.Conf().Consistent.Hash.Type == constants.IpPort {
		nodeJob = fmt.Sprintf("%s%s:%d/", constants.WorkerJobs, localIP, configs.Conf().Serves[0].ServePort)
	}
	if nodeJob == "" {
		nodeJob = fmt.Sprintf("%s%s/", constants.WorkerJobs, localIP)
	}

	// get /cron/jobs/目录下所有的任务，并且获知当前集群的revision
	if getResp, err = mgr.kv.Get(ctx, nodeJob, clientv3.WithPrefix()); err != nil {
		return
	}

	for GWorkerNode.GetMembers() == nil || len(GWorkerNode.GetMembers()) == 0 {
		time.Sleep(50 * time.Millisecond)
	}

	// TODO: 处理这里的时候需要master节点的时间，用于处理上次调度未被清理的任务数据，每次master启动都需要记录一个启动时间
	// 当前有那些任务
	for _, kvPair = range getResp.Kvs {
		// 反序列化json得到job
		if job, err = entity.UnpackJobEntity(kvPair.Value); err != nil {
			log.WithFields(log.Fields{
				"WatchJobs": "err",
			}).Error(err)
			continue
		}

		nodeInfo = string(kvPair.Key)
		nodeInfo = entity.ExtractNodeInfoName(nodeInfo, job.Name)
		// 获取key对应的node IP 与自己的ip节点比较如果不是自己节点，就删除自己的信息
		if nodeInfo != job.NodeIp {

			// 清除老数据
			for _, member = range GWorkerNode.GetMembers() {
				if member.String() == nodeJob {
					existent = true
					break
				}
			}

			if !existent {
				// 删除对应的job信息
				_, err = mgr.kv.Delete(context.TODO(), string(kvPair.Key))
				if err != nil {
					log.Error(err)
				}
			}

			continue
		}

		jobEvent = entity.BuildJobEvent(constants.JobEventSave, job)

		//job同步给调度协程(scheduler)
		log.WithFields(log.Fields{
			"jobEvent-All": jobEvent,
		}).Log(log.InfoLevel)

		GScheduler.PushJobEvent(jobEvent)
	}

	// 从该revision向后监听变化事件
	async.GO(func() {
		// 监听协程
		// 从GET时刻的后续版本开始监听
		//watchStartRevision = getResp.Header.Revision + 1
		// 监听/cron/jobs/目录的后续变化
		watchChan = mgr.watcher.Watch(context.Background(), nodeJob, clientv3.WithPrefix(), clientv3.WithPrevKV())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					//反序列化Job，
					if job, err = entity.UnpackJobEntity(watchEvent.Kv.Value); err != nil {
						continue
					}

					nodeInfo = string(watchEvent.Kv.Key)
					nodeInfo = entity.ExtractNodeInfoName(nodeInfo, job.Name)
					// 获取key对应的node IP 与自己的ip节点比较如果不是自己节点，就删除自己的信息
					if nodeInfo != job.NodeIp {
						for GWorkerNode.GetMembers() == nil || len(GWorkerNode.GetMembers()) == 0 {
							log.WithFields(log.Fields{
								"slave-WatchJobs":         nodeInfo,
								"slave-WatchJobs-writing": "...",
							}).Info()
							time.Sleep(100 * time.Millisecond)
						}

						// 清除老数据
						for _, member = range GWorkerNode.GetMembers() {
							if member.String() == nodeJob {
								existent = true
								break
							}
						}

						if !existent {
							// 删除对应的job信息
							_, err = mgr.kv.Delete(context.TODO(), string(watchEvent.Kv.Key))
							if err != nil {
								log.Error(err)
							}
							continue
						}
					}

					// 构建一个Event更新事件
					jobEvent = entity.BuildJobEvent(constants.JobEventSave, job)

					log.WithFields(log.Fields{
						"jobEvent-PUT": jobEvent,
					}).Log(log.InfoLevel)

					// 推一个更新事件给scheduler
					GScheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: // 任务删除事件
					// delete /cron/jobs/job0
					jobName = entity.ExtractNodeJobName(string(watchEvent.Kv.Key), nodeJob)

					// 构建一个删除的凭证
					job = &entity.JobEntity{
						Name: jobName,
					}

					// 构建一个删除Event
					jobEvent = entity.BuildJobEvent(constants.JobEventDelete, job)

					log.WithFields(log.Fields{
						"jobEvent-DELETE": jobEvent,
					}).Log(log.InfoLevel)

					// 推一个删除事件给scheduler
					GScheduler.PushJobEvent(jobEvent)
				}
			}
		}
	})

	return
}

// SchedulerJobsForAddWorkerNode 处理重新调度任务到新的工作节点上
func (mgr *Mgr) schedulerJobsForAddWorkerNode(changeEvent *entity.WorkerChangeEvent) {
	// 获取工作节点信息
	var (
		getResp  *clientv3.GetResponse
		kvPair   *mvccpb.KeyValue
		job      *entity.JobEntity
		nodeJob  string
		nodeInfo string
		localIP  string
		member   consistent.Member
		jobKey   string
		jobValue []byte
		jobEvent *entity.JobEvent

		err error
	)

	// 获取本机IP
	if localIP, err = tools.GetLocalIP(); err != nil {
		return
	}

	if configs.Conf().Consistent.Hash.Type == constants.IpPort {
		nodeJob = fmt.Sprintf("%s%s:%d/", constants.WorkerJobs, localIP, configs.Conf().Serves[0].ServePort)
	}
	if nodeJob == "" {
		nodeJob = fmt.Sprintf("%s%s/", constants.WorkerJobs, localIP)
	}

	// get /cron/jobs/目录下所有的任务，并且获知当前集群的revision
	if getResp, err = mgr.kv.Get(context.TODO(), nodeJob, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前有那些任务
	for _, kvPair = range getResp.Kvs {
		// 反序列化json得到job
		if job, err = entity.UnpackJobEntity(kvPair.Value); err != nil {
			log.WithFields(log.Fields{
				"WatchJobs": "err",
			}).Error(err)
			continue
		}

		nodeInfo = string(kvPair.Key)
		nodeInfo = entity.ExtractNodeInfoName(nodeInfo, job.Name)

		member = GWorkerNode.LocateKey([]byte(job.Name))

		// 该任务无变化继续在该任务节点执行
		if member.String() == nodeInfo {
			continue
		}

		// 如果不等于就要把该任务从这个工作节点删除，然后移入其他节点
		// 先移入其他节点
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
		jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, entity.ExtractWorkerIP(member.String()), job.Name)

		// 保存到etcd中，并且获取以前的值信息
		if _, err = mgr.kv.Put(context.TODO(), jobKey, string(jobValue)); err != nil {
			continue
		}

		// 构建一个删除Event
		jobEvent = entity.BuildJobEvent(constants.JobEventDelete, job)

		log.WithFields(log.Fields{
			"schedulerJobsForAddWorkerNode-DELETE": jobEvent,
		}).Log(log.InfoLevel)

		// 推一个删除事件给scheduler
		GScheduler.PushJobEvent(jobEvent)

		// 移除远程该工作节点记录，删除对应的job信息
		_, err = mgr.kv.Delete(context.TODO(), string(kvPair.Key))
		if err != nil {
			log.WithFields(log.Fields{
				"delete-from": member,
				"delete-to":   nodeJob,
			}).Error(err)
		}

		log.WithFields(log.Fields{
			"delete-from": member,
			"delete-to":   nodeJob,
		}).Info()
	}
}

// SchedulerAddWorkerNode 将自身节点的任务按一致性hash处理，如果是其他节点就移入其他节点，如果是自身的节点就不用移动
func (mgr *Mgr) SchedulerAddWorkerNode() {
	var (
		changeEvent *entity.WorkerChangeEvent
	)

	async.GO(func() {
		for changeEvent = range mgr.workerNodeChangeEvent {
			if changeEvent == nil {
				continue
			}

			switch changeEvent.ChangeType {
			case constants.WorkerEventAdd:
				mgr.schedulerJobsForAddWorkerNode(changeEvent)
				// 交由master处理，然后均匀分配给存活的节点
				log.WithFields(log.Fields{
					"WorkerEventDelete": changeEvent,
				}).Info("WorkerEventDelete")
			case constants.WorkerEventDelete:
				// 交由master处理，然后均匀分配给存活的节点
				log.WithFields(log.Fields{
					"WorkerEventDelete": changeEvent,
				}).Info("WorkerEventDelete")
			}
		}
	})
}

// PushWorkerNodeChangeEvent 处理新节点加入后的任务分配问题
func (mgr *Mgr) PushWorkerNodeChangeEvent(changeEvent *entity.WorkerChangeEvent) {
	mgr.workerNodeChangeEvent <- changeEvent
}

// DeleteWorkerDataCallback 关闭服务时删除工作节点分配任务数据
func (mgr *Mgr) DeleteWorkerDataCallback(ctx context.Context) {
	/*var (
		jobNodeKey string
		err        error
	)

	jobNodeKey = fmt.Sprintf("%s%s", constants.WorkerJobs, GWorkerNode.ip)
	if _, err = mgr.kv.Delete(context.TODO(), jobNodeKey, clientv3.WithPrefix()); err != nil {
		log.WithFields(log.Fields{
			"DeleteWorkerDataCallback": "err",
		}).Error(err)
	}*/
}

/*// SchedulerNewWorkerNodeJobs 调度计算后的任务到新节点上去
func (mgr *Mgr) SchedulerNewWorkerNodeJobs(changeEvent *entity.WorkerChangeEvent) {
	// 计算所有节点的任务信息，先删除现在在老节点但是本来分配在新节点的任务，然后分配该任务到新节点
	var (
		getResp  *clientv3.GetResponse
		kvPair   *mvccpb.KeyValue
		job      *entity.JobEntity
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
		if job, err = entity.UnpackJobEntity(kvPair.Value); err != nil {
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
func (mgr *Mgr) SchedulerOldWorkerNodeJobs(changeEvent *entity.WorkerChangeEvent) {
	//分配删除节点上所有的任务到其他工作节点上， 然后删除该节点
	var (
		workerJobsKey string
		getResp       *clientv3.GetResponse
		kvpair        *mvccpb.KeyValue
		job           *entity.JobEntity
		member        consistent.Member
		jobKey        string
		jobValue      []byte
		err           error
	)

	// 获取子任务路径
	workerJobsKey = fmt.Sprintf("%s%s", constants.WorkerJobs, changeEvent.WorkerName)

	if getResp, err = mgr.kv.Get(context.TODO(), workerJobsKey, clientv3.WithPrefix()); err != nil {
		log.WithFields(log.Fields{
			"SchedulerOldWorkerNodeJobs": "err",
		}).Error(err)
		return
	}

	// 遍历
	for _, kvpair = range getResp.Kvs {
		//反序列化Job，
		if job, err = entity.UnpackJobEntity(kvpair.Value); err != nil {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs": "err",
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

		jobKey = fmt.Sprintf("%s%s/%s", constants.WorkerJobs, job.NodeIp, job.Name)
		if _, err = mgr.kv.Put(context.TODO(), jobKey, string(jobValue)); err != nil {
			log.WithFields(log.Fields{
				"SchedulerOldWorkerNodeJobs": "err",
			}).Error(err)
			continue
		}
	}

	// 删除该节点
	if _, err = mgr.kv.Delete(context.TODO(), workerJobsKey); err != nil {
		log.WithFields(log.Fields{
			"SchedulerOldWorkerNodeJobs": "err",
		}).Error(err)
		return
	}
}*/

// CreateJobLock 构建一个job锁
func (mgr *Mgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	return InitJobLock(jobName, mgr.kv, mgr.lease)
}

// watchKiller 监听强杀任务通知
func (mgr *Mgr) watchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		job        *entity.JobEntity
		jobEvent   *entity.JobEvent
		jobName    string
	)

	// 监听/cron/killer目录
	async.GO(func() {
		// 监听协程
		// 监听/cron/killer目录变化
		watchChan = mgr.watcher.Watch(context.Background(), constants.JobKillDir, clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 杀死任务通知
					// /cron/killer/job0 -> job0
					jobName = entity.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &entity.JobEntity{
						Name: jobName,
					}

					jobEvent = entity.BuildJobEvent(constants.JobEventKill, job)
					// 推一个变化事件给scheduler
					GScheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: // killer标记过期，被自动删除

				}

			}
		}
	})
}
