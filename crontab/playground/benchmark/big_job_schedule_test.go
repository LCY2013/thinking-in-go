package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var (
	crontabMaster = "http://127.0.0.1:8080/job/save"
)

// JobEntity 定时任务
type JobEntity struct {
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // shell 命令
	CronExpr string `json:"cronExpr"` // cron 表达式
}

func TestBigJobs(t *testing.T) {
	var (
		buffer   bytes.Buffer
		resp     *http.Response
		job      *JobEntity
		jobValue []byte
		err      error
	)

	for i := 0; i < 100000; i++ {
		job = &JobEntity{
			Name:     fmt.Sprintf("job%d", i),
			Command:  fmt.Sprintf("echo hello%d", i),
			CronExpr: "*/5 * * * * * *",
		}

		if jobValue, err = json.Marshal(job); err != nil {
			t.Error(err)
		}
		buffer.Write(jobValue)
		if resp, err = http.Post(crontabMaster, "application/json", &buffer); err != nil {
			t.Error(err)
		}
		t.Log(resp)
		buffer.Reset()
	}
}
