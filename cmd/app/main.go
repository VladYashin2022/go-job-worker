package main

import (
	"go-job-worker/internal/model"
	"go-job-worker/internal/queue"
	"go-job-worker/internal/worker"
	"sync"
)

func main() {
	qSize := 5
	q := queue.NewQueue(qSize)

	var jobsWg sync.WaitGroup
	var workersWg sync.WaitGroup

	workerCount := 3

	for i := 1; i <= workerCount; i++ {
		workersWg.Add(1)
		go worker.StartWorker(i, q.JobsCh, &jobsWg, &workersWg)
	}

	for i := 1; i <= 5; i++ {

		job := model.Job{
			ID:      i,
			Type:    "email",
			Payload: "send welcome email",
		}

		jobsWg.Add(1)
		q.JobsCh <- job
	}

	close(q.JobsCh)

	jobsWg.Wait()
	workersWg.Wait()

}
