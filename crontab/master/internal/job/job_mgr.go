package service

import (
	"context"
	"encoding/json"
	"fmt"
	jobEntity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
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
			kv:     client.KV,
			lease:  client.Lease,
		}
	})
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
	)

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
