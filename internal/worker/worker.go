package worker

import (
	"context"
	"fmt"
	"go-job-worker/internal/model"
	"math/rand"
	"sync"
	"time"
)

func StartWorker(
	id int,
	jobs <-chan model.Job,
	jobsWg *sync.WaitGroup,
	workersWg *sync.WaitGroup,
	sem chan struct{},
	ctx context.Context,
) {
	defer workersWg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("worker %d exiting (channel closed)\n", id)
				return
			}

			sem <- struct{}{} // add semafor
			fmt.Printf("worker %d, sem +1\n", id)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			fmt.Printf(
				"worker %d processing job %d: %s\n",
				id,
				job.ID,
				job.Payload,
			)
			jobsWg.Done()

			<-sem
			fmt.Printf("worker %d, sem -1\n", id)

		case <-ctx.Done():
			fmt.Println("worker stopped")
			return
		}
	}
}
