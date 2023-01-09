package service

import (
	entity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"math/rand"
	"os/exec"
	"time"
)

// Executor 执行器
type Executor struct {
}

var (
	GExecutor *Executor
)

// ExecuteJob 执行一个任务
func (e *Executor) ExecuteJob(info *entity.JobExecuteInfo) {
	async.GO(func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *entity.JobExecuteResult
			jobLock *JobLock
		)

		// 任务结果
		result = &entity.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		// 初始化分布式锁
		jobLock = G_MGR.CreateJobLock(info.Job.Name)

		// 记录任务开始时间
		result.StartTime = time.Now()

		// 随机休眠（0-1）s
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		if err = jobLock.TryLock(); err != nil {
			result.EndTime = time.Now()
			result.Err = err
			// 任务执行结束，把执行的结果返回给Schedule， Scheduler会从executingTable中删除掉执行记录
			GScheduler.PushJobResult(result)
			return
		}
		// 释放锁
		defer jobLock.UnLock()

		// 上锁成功后重置任务启动时间
		result.StartTime = time.Now()

		// 执行shell命令
		cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)

		// 执行并捕获输出
		output, err = cmd.CombinedOutput()

		// 记录任务结束时间
		result.EndTime = time.Now()

		// 执行结果
		result.Output = output
		result.Err = err

		// 任务执行结束，把执行的结果返回给Schedule， Scheduler会从executingTable中删除掉执行记录
		GScheduler.PushJobResult(result)
	})
}

// InitExecutor 初始化执行器
func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
