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
