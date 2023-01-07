package container

import "github.com/LCY2013/thinking-in-go/crontab/master/internal/web/controller"

type WebContainer struct {
	JobController *controller.JobController
}

func NewContainer(JobController *controller.JobController) *WebContainer {
	return &WebContainer{
		JobController: JobController,
	}
}
