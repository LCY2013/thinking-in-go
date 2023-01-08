package entity

import (
	"encoding/json"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	"strings"
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
	job       *JobEntity
}

// BuildJobEvent 任务变化事件有两种：1）更新任务 2）删除任务
func BuildJobEvent(eventType int, job *JobEntity) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		job:       job,
	}
}
