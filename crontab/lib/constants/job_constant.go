package constants

const (
	// JobDir 任务根目录
	JobDir = "/cron/jobs/"
	// JobKillDir 任务取消目录
	JobKillDir = "/cron/killer/"
	// JobLockDir 任务锁目录
	JobLockDir = "/cron/lock/"

	// JobEventSave 保存任务事件
	JobEventSave = 1

	// JobEventDelete 删除任务事件
	JobEventDelete = 2

	// JobEventKill 强杀任务事件
	JobEventKill = 3
)
