package scheduler

import (
	entity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/lib/constants"
	log "github.com/sirupsen/logrus"
	"time"
)

// Scheduler 任务调度
type Scheduler struct {
	jobEventChan      chan *entity.JobEvent              // etcd 任务事件队列
	jobPlanTable      map[string]*entity.JobSchedulePlan // 任务执行计划表
	jobExecutingTable map[string]*entity.JobExecuteInfo  // 任务执行表
	jobResultChan     chan *entity.JobExecuteResult      // 任务结果队列
}

var (
	GScheduler *Scheduler
)

// InitScheduler 初始化调度器
func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		jobEventChan:      make(chan *entity.JobEvent, 1000),
		jobPlanTable:      make(map[string]*entity.JobSchedulePlan, 100),
		jobExecutingTable: make(map[string]*entity.JobExecuteInfo, 100),
		jobResultChan:     make(chan *entity.JobExecuteResult, 1000),
	}

	// 启动调度协程
	async.GO(func() {
		GScheduler.scheduleLoop()
	})
	return
}

// scheduleLoop 调度协程
func (s *Scheduler) scheduleLoop() {
	var (
		jobEvent      *entity.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *entity.JobExecuteResult
	)

	// 执行一遍初始化
	scheduleAfter = s.TrySchedule()

	// 调度的延时定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	// 定时任务JobEntity
	for {
		select {
		case jobEvent = <-s.jobEventChan: // 监听任务变化事件
			// 对内存中的任务事件实时同步
			s.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: // 最近的任务到期执行
		case jobResult = <-s.jobResultChan: // 监听任务执行结果
			// 处理任务执行结果
			s.handleJobResult(jobResult)
		}

		// 执行一遍初始化
		scheduleAfter = s.TrySchedule()
		// 重置调度时间间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

// handleJobEvent 处理任务事件
func (s *Scheduler) handleJobEvent(jobEvent *entity.JobEvent) {
	var (
		err            error
		jobSchedulePan *entity.JobSchedulePlan
		jobExisted     bool
	)

	switch jobEvent.EventType {
	case constants.JobEventSave: // 保存任务事件
		// 构建一个任务调度表
		if jobSchedulePan, err = entity.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			log.WithFields(log.Fields{
				"BuildJobSchedulePlan-JobEventSave": err,
			}).Log(log.InfoLevel)
			return
		}

		// 更新信息到调度表里面
		s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePan
	case constants.JobEventDelete: // 删除任务事件
		if jobSchedulePan, jobExisted = s.jobPlanTable[jobEvent.Job.Name]; !jobExisted {
			return
		}
		delete(s.jobPlanTable, jobEvent.Job.Name)
	}
}

// PushJobEvent 推送任务变化事件
func (s *Scheduler) PushJobEvent(jobEvent *entity.JobEvent) {
	s.jobEventChan <- jobEvent
}

// TrySchedule 计算任务调度状态
func (s *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *entity.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)

	// 如果任务表没有任务，就随意休眠多久
	if len(s.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	// 获取当前时间
	now = time.Now()

	// 遍历当前节点的所有任务
	for _, jobPlan = range s.jobPlanTable {
		// 过期的任务立即执行
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			log.WithFields(log.Fields{
				"TrySchedule-DO": jobPlan,
			}).Log(log.InfoLevel)
			// 尝试执行任务
			s.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下一次执行的时间
		}

		// 统计最近要过期的任务的时间（N秒后过期 == scheduleAfter）
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	if nearTime != nil {
		// 下次调度时间(最近要执行的任务调度时间-当前时间)
		scheduleAfter = (*nearTime).Sub(now)
	}
	return
}

// TryStartJob 尝试执行任务
func (s *Scheduler) TryStartJob(jobPlan *entity.JobSchedulePlan) {
	var (
		jobExecuteInfo *entity.JobExecuteInfo
		jobExecuting   bool
	)
	// 调度和执行两件事

	// 执行的任务可能允许很久，1分钟会调度60次，但是只执行一次
	if jobExecuteInfo, jobExecuting = s.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		return
	}

	// 如果任务正在执行，跳过本次调度
	jobExecuteInfo = entity.BuildJobExecuteInfo(jobPlan)

	// 保存执行状态
	s.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	// 执行任务
	GExecutor.ExecuteJob(jobExecuteInfo)
}

// PushJobResult 回传执行结果
func (s *Scheduler) PushJobResult(jobResult *entity.JobExecuteResult) {
	s.jobResultChan <- jobResult
}

// handleJobResult 处理任务结果信息
func (s *Scheduler) handleJobResult(result *entity.JobExecuteResult) {
	// 删除任务执行状态
	delete(s.jobExecutingTable, result.ExecuteInfo.Job.Name)

	log.WithFields(log.Fields{
		"handleJobResult": result,
		"output":          string(result.Output),
	}).Log(log.InfoLevel)
}
