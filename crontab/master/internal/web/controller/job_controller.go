package controller

import (
	"context"
	jobentity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	jobservice "github.com/LCY2013/thinking-in-go/crontab/master/internal/job"
	log "github.com/sirupsen/logrus"
)

type JobController struct {
}

func NewJobController() *JobController {
	return &JobController{}
}

type CreateJobRequest struct {
	jobentity.JobEntity
}

type CreateJobResponse struct {
	jobentity.JobEntity
}

func (c *JobController) CreateJob(ctx context.Context, createJobRequest *CreateJobRequest) (*CreateJobResponse, error) {
	var (
		err    error
		oldJob *jobentity.JobEntity
	)

	log.WithFields(log.Fields{
		"CreateJob": "CreateJob",
	}).Logf(log.InfoLevel, "%+v", *createJobRequest)
	if createJobRequest == nil {
		return &CreateJobResponse{}, nil
	}

	// 保存到etcd中
	oldJob, err = jobservice.G_MGR.SaveJob(ctx, &createJobRequest.JobEntity)

	resp := &CreateJobResponse{}

	if oldJob != nil {
		resp.JobEntity = *oldJob
	}

	return resp, err
}
