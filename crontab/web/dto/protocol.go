package dto

// Job 定时任务
type Job struct {
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // shell 命令
	CronExpr string `json:"cronExpr"` // cron 表达式
}
