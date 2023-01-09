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

type DelJobRequest struct {
	Name string `json:"name"` // 任务名称
}

type DelJobResponse struct {
	jobentity.JobEntity
}

func (c *JobController) DelJob(ctx context.Context, delJobRequest *DelJobRequest) (*DelJobResponse, error) {
	var (
		err    error
		oldJob *jobentity.JobEntity
	)

	log.WithFields(log.Fields{
		"DelJob": "DelJob",
	}).Logf(log.InfoLevel, "%+v", *delJobRequest)
	if delJobRequest == nil {
		return &DelJobResponse{}, nil
	}

	// 从etcd中删除
	oldJob, err = jobservice.G_MGR.DeleteJob(ctx, delJobRequest.Name)

	resp := &DelJobResponse{}

	if oldJob != nil {
		resp.JobEntity = *oldJob
	}

	return resp, err
}

type ListJobRequest struct {
}

type ListJobResponse struct {
	JobList []*jobentity.JobEntity `json:"list"`
}

func (c *JobController) ListJob(ctx context.Context, listJobRequest *ListJobRequest) (*ListJobResponse, error) {
	var (
		err     error
		listJob []*jobentity.JobEntity
	)

	log.WithFields(log.Fields{
		"ListJob": "ListJob",
	}).Logf(log.InfoLevel, "%+v", *listJobRequest)
	if listJobRequest == nil {
		return &ListJobResponse{}, nil
	}

	// 从etcd中删除
	listJob, err = jobservice.G_MGR.ListJob(ctx)

	resp := &ListJobResponse{}

	if listJob != nil {
		resp.JobList = listJob
	}

	return resp, err
}

type KillJobRequest struct {
	Name string `json:"name"` // 任务名称
}

type KillJobResponse struct {
}

func (c *JobController) KillJob(ctx context.Context, killJobRequest *KillJobRequest) error {
	var (
		err error
	)

	log.WithFields(log.Fields{
		"KillJob": "KillJob",
	}).Logf(log.InfoLevel, "%+v", *killJobRequest)
	if killJobRequest == nil {
		return nil
	}

	// 从etcd中删除
	err = jobservice.G_MGR.KillJob(ctx, killJobRequest.Name)

	return err
}