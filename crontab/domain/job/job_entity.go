package entity

import (
	"encoding/json"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// JobEntity 定时任务
type JobEntity struct {
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // shell 命令
	CronExpr string `json:"cronExpr"` // cron 表达式
}

// UnpackJobEntity 反序列化job信息
func UnpackJobEntity(value []byte) (ret *JobEntity, err error) {
	var (
		job *JobEntity
	)

	job = &JobEntity{}

	if err = json.Unmarshal(value, &job); err != nil {
		return
	}

	ret = job
	return
}

// ExtractJobName 从ETCD的key中获取对应的任务名称
// /cron/jobs/job0 -> job0
func ExtractJobName(jobName string) string {
	return strings.TrimPrefix(jobName, constants.JobDir)
}

// JobEvent job event
type JobEvent struct {
	EventType int // SAVE、DELETE
	Job       *JobEntity
}

// BuildJobEvent 任务变化事件有两种：1）更新任务 2）删除任务
func BuildJobEvent(eventType int, job *JobEntity) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// JobSchedulePlan 任务调度计划
type JobSchedulePlan struct {
	Job      *JobEntity           // 要调度的任务
	Expr     *cronexpr.Expression // 解析好的cronexpr表达式
	NextTime time.Time            // 下次调度时间
}

// BuildJobSchedulePlan 构建任务调度计划
func BuildJobSchedulePlan(job *JobEntity) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	// 解析JOB中的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		log.WithFields(log.Fields{
			"BuildJobSchedulePlan-JobEventSave": err,
		}).Log(log.InfoLevel)
		return
	}

	// 生成任务调度计划对象
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}
