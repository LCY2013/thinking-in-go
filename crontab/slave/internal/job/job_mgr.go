package service

import (
	"context"
	entity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
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
			kv:      client.KV,
			lease:   client.Lease,
			watcher: client.Watcher,
		}
	})

	// 启动任务监听
	err = G_MGR.WatchJobs(context.Background())

	// 启动强杀kill任务通知
	G_MGR.watchKiller()

	return
}

// WatchJobs 监听任务变化
func (mgr *Mgr) WatchJobs(ctx context.Context) (err error) {
	var (
		getRespn           *clientv3.GetResponse
		kvPair             *mvccpb.KeyValue
		job                *entity.JobEntity
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *entity.JobEvent
	)
	// get /cron/jobs/目录下所有的任务，并且获知当前集群的revision
	if getRespn, err = mgr.kv.Get(ctx, constants.JobDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前有那些任务
	for _, kvPair = range getRespn.Kvs {
		// 反序列化json得到job
		if job, err = entity.UnpackJobEntity(kvPair.Value); err != nil {
			log.WithFields(log.Fields{
				"WatchJobs": "err",
			}).Error(err)
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
		watchStartRevision = getRespn.Header.Revision + 1
		// 监听/cron/jobs/目录的后续变化
		watchChan = mgr.watcher.Watch(context.Background(), constants.JobDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					//反序列化Job，
					if job, err = entity.UnpackJobEntity(watchEvent.Kv.Value); err != nil {
						continue
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
					jobName = entity.ExtractJobName(string(watchEvent.Kv.Key))

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
