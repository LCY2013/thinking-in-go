package cron_usage

import (
	"github.com/gorhill/cronexpr"
	"testing"
	"time"
)

func TestCronexpr(t *testing.T) {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	// linux crontab 只支持5个
	// 那一分钟(0-59)，那个小时(0-23)，那一天(1-31)，那个月(1-12)，星期几(0-6)
	// cronexpr 多支持了 秒、年配置(2018-2099)

	/*if expr, err = cronexpr.Parse("* * * * *"); err != nil {
		t.Error(err)
		return
	}*/

	// 每隔5分钟执行一次
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		t.Error(err)
		return
	}

	// 当前时间
	now = time.Now()

	// 下次调度时间
	nextTime = expr.Next(now)
	t.Logf("%s", nextTime)

	// 等待这个定时器超时
	time.AfterFunc(nextTime.Sub(now), func() {
		t.Logf("%s", time.Now())
	})

	time.Sleep(6 * time.Second)
}

// CronJob 任务抽象
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time // expr.Next(now
}

func TestCronexprScheduling(t *testing.T) {
	// 需要有一个调度协程， 它定时监测所有的Cron任务，谁过期就执行谁

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob // key 任务的名字
	)

	scheduleTable = make(map[string]*CronJob)

	// now
	now = time.Now()

	// 1、定义两个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册到调度表
	scheduleTable["job-1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册到调度表
	scheduleTable["job-2"] = cronJob

	// 启动调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)

		// 定时检查一下任务调度表
		for {
			now = time.Now()

			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程，执行这个任务
					go func(jobName string) {
						t.Logf("[%s]-执行 ", jobName)
					}(jobName)

					// 计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					t.Logf("[%s]-下次执行时间:%s", jobName, cronJob.nextTime)
				}
			}

			// 休眠100ms
			select {
			case <-time.NewTimer(100 * time.Millisecond).C:
			}
		}
	}()

	time.Sleep(100 * time.Second)
}
