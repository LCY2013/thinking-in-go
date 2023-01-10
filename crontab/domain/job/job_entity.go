package entity

import (
	"context"
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

// JobExecuteInfo 任务执行状态
type JobExecuteInfo struct {
	Job        *JobEntity         // 任务信息
	PlanTime   time.Time          // 计划调度时间
	RealTime   time.Time          // 实际的调度时间
	CancelCtx  context.Context    `json:"-"` // command命令任务取消的context
	CancelFunc context.CancelFunc `json:"-"` // command命令任务取消的func
}

// BuildJobExecuteInfo 构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, // 计划调度时间
		RealTime: time.Now(),               // 真实调度时间
	}

	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

// JobExecuteResult 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          // 脚本输出
	Err         error           // 脚本错误原因
	StartTime   time.Time       // 启动时间
	EndTime     time.Time       // 结束时间
}

// ExtractKillerName 从ETCD的key中获取对应的强杀任务名称
// /cron/killer/job0 -> job0
func ExtractKillerName(killerName string) string {
	return strings.TrimPrefix(killerName, constants.JobKillDir)
}

// JobLog 任务执行日志
type JobLog struct {
	JobName      string `bson:"jobName"`      // 任务名字
	Command      string `bson:"command"`      // 脚本命令
	Err          string `bson:"err"`          // 错误原因
	Output       string `bson:"output"`       // 输出信息
	PlanTime     int64  `bson:"planTime"`     // 计划开始时间
	ScheduleTime int64  `bson:"scheduleTime"` // 实际调度时间
	StartTime    int64  `bson:"startTime"`    // 任务执行开始时间
	EndTime      int64  `bson:"endTime"`      // 任务执行结束时间
}
