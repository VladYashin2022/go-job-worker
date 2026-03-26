package producer

import (
	"context"
	"fmt"
	"go-job-worker/internal/model"
	"sync"
)

func Producer(
	count int,
	jobCh chan<- model.Job,
	jobsWg *sync.WaitGroup,
	prodDone chan struct{},
	ctx context.Context,
) {
	defer close(prodDone)

	for i := 1; i <= count; i++ {

		job := model.Job{
			ID:      i,
			Type:    "email",
			Payload: "send welcome email",
		}

		select {
		case <-ctx.Done():
			return
		case jobCh <- job:
			jobsWg.Add(1)
			fmt.Println("sending jobs")
		}
	}

}
