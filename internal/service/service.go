package service

import (
	"go-job-worker/internal/model"
	"go-job-worker/internal/storage"
)

type JobsService struct {
	storage *storage.JobsStorage
	ch      chan model.Job
}

func NewJobsService(s *storage.JobsStorage, jobsCh chan model.Job) *JobsService {
	return &JobsService{
		storage: s,
		ch:      jobsCh,
	}
}

func (s *JobsService) GetJob(id int) (model.Job, error) {
	return s.storage.GetByID(id)
}

func (s *JobsService) GetAllJobs() []model.Job {
	return s.storage.GetAll()
}

func (s *JobsService) CreateJob(job model.Job) {
	s.storage.Create(job)
	s.ch <- job
}

type requestJob struct {
}

type patchJob struct {
}
