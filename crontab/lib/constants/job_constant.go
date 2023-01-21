package constants

const (
	// JobDir 任务根目录
	JobDir = "/cron/jobs/"
	// JobKillDir 任务取消目录
	JobKillDir = "/cron/killer/"
	// JobLockDir 任务锁目录
	JobLockDir = "/cron/lock/"
	// JobWorkerRegisterDir 工作节点注册目录
	JobWorkerRegisterDir = "/cron/worker/register/"
	// WorkerJobs 每个工作节点的所有任务信息
	WorkerJobs = "/cron/worker/jobs/"
	// WorkerMasterDir 工作节点master目录
	WorkerMasterDir = "/cron/worker/master"

	// JobEventSave 保存任务事件
	JobEventSave = 1

	// JobEventDelete 删除任务事件
	JobEventDelete = 2

	// JobEventKill 强杀任务事件
	JobEventKill = 3

	// IP 一致性hash类型
	IP = "IP"

	// IpPort 一致性hash类型
	IpPort = "IP:PORT"

	// WorkerEventAdd 工作节点新增事件
	WorkerEventAdd = 1

	// WorkerEventDelete 工作节点删除事件
	WorkerEventDelete = 2
)
