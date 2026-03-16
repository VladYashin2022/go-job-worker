package worker

import (
	"fmt"
	"go-job-worker/internal/model"
	"sync"
)

func StartWorker(
	id int,
	jobs <-chan model.Job,
	jobsWg *sync.WaitGroup,
	workersWg *sync.WaitGroup,
) {
	defer workersWg.Done()

	for job := range jobs {

		fmt.Printf(
			"worker %d processing job %d: %s\n",
			id,
			job.ID,
			job.Payload,
		)

		jobsWg.Done()
	}
}
