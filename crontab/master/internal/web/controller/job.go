package controller

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type JobController struct {
}

func NewJobController() *JobController {
	return &JobController{}
}

type CreateJobRequest struct {
	JobName string `json:"jobName"`
}

func (c *JobController) CreateJob(ctx context.Context, createJobRequest *CreateJobRequest) error {
	log.WithFields(log.Fields{
		"CreateJob": "CreateJob",
	}).Logf(log.InfoLevel, "%+v", *createJobRequest)
	return nil
}
