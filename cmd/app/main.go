package main

import (
	"context"
	"fmt"
	"go-job-worker/internal/model"
	"go-job-worker/internal/producer"
	"go-job-worker/internal/queue"
	"go-job-worker/internal/worker"
	"os"
	"os/signal"
	"sync"
)

func main() {
	qSize := 5
	q := queue.NewQueue(qSize)

	var jobsWg sync.WaitGroup
	var workersWg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	sem := make(chan struct{}, 5) //semafor

	//producer записывает job в канал jobs и отправляет сигнал, когда закончит
	producerDone := make(chan struct{})
	jobsCount := 50
	go producer.Producer(jobsCount, q.JobsCh, &jobsWg, producerDone, ctx)

	//worker обрабатывает job из канала jobs
	workerCount := 10

	for i := 1; i <= workerCount; i++ {
		workersWg.Add(1)
		go worker.StartWorker(
			i,
			q.JobsCh,
			&jobsWg,
			&workersWg,
			sem,
			ctx,
		)
	}

	select {
	case <-producerDone:
		fmt.Println("producer finished")
	case <-sigCh:
		fmt.Println("\nshutdown signal received")
		cancel()
	}

	shutdown(q.JobsCh, &jobsWg, &workersWg)

	fmt.Println("program finished")
}

func shutdown(
	ch chan model.Job,
	jobsWg *sync.WaitGroup,
	workersWg *sync.WaitGroup,
) {
	fmt.Println("channel closed")
	close(ch)

	fmt.Println("waiting jobs...")
	jobsWg.Wait()

	fmt.Println("waiting workers...")
	workersWg.Wait()
}
