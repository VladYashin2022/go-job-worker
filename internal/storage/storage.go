package storage

import (
	"errors"
	"go-job-worker/internal/model"
)

type JobsStorage struct {
	jobs map[int]model.Job
}

func NewJobsStorage() *JobsStorage {
	return &JobsStorage{
		jobs: make(map[int]model.Job),
	}
}

func (s *JobsStorage) Create(job model.Job) {
	s.jobs[job.ID] = job
}

func (s *JobsStorage) GetByID(id int) (model.Job, error) {
	job, ok := s.jobs[id]
	if ok != true {
		return model.Job{}, errors.New("Job not found")
	}

	return job, nil
}

func (s *JobsStorage) GetAll() []model.Job {
	jobs := make([]model.Job, 0, len(s.jobs))

	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}
