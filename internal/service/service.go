package service

import (
	"go-job-worker/internal/model"
	"go-job-worker/internal/storage"
)

type JobsService struct {
	storage *storage.JobsStorage
	ch      chan model.Job
	nextInt int // для создания id
}

func NewJobsService(s *storage.JobsStorage, jobsCh chan model.Job) *JobsService {
	return &JobsService{
		storage: s,
		ch:      jobsCh,
		nextInt: 0,
	}
}

func (s *JobsService) GetJob(id int) (model.Job, error) {
	return s.storage.GetByID(id)
}

func (s *JobsService) GetAllJobs() []model.Job {
	return s.storage.GetAll()
}

func (s *JobsService) CreateJob(jobType, jobPayload string) (model.Job, error) {
	s.nextInt = s.nextInt + 1
	jobID := s.nextInt

	var job model.Job = model.Job{
		ID:      jobID,
		Type:    jobType,
		Payload: jobPayload,
	}
	s.storage.Create(job) // запись job в storage
	s.ch <- job
	return job, nil //заглушка, чтобы потом подставить сюда error DB
}
