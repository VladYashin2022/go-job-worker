package queue

import "go-job-worker/internal/model"

type JobQueue struct {
	JobsCh chan model.Job
}

func NewQueue(size int) *JobQueue {
	return &JobQueue{
		JobsCh: make(chan model.Job, size),
	}
}
