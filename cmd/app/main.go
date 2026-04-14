package main

import (
	"context"
	"fmt"
	"go-job-worker/internal/handler"
	"go-job-worker/internal/model"
	"go-job-worker/internal/queue"
	"go-job-worker/internal/service"
	"go-job-worker/internal/storage"
	"go-job-worker/internal/worker"
	"os"
	"os/signal"
	"sync"
)

func main() {
	//создаем структуру с каналом jobs
	qSize := 5
	q := queue.NewQueue(qSize)

	var jobsWg sync.WaitGroup
	var workersWg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	//канал для отлова сигналов о прекращении работы приложения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	//создаем хранилище (мапу), сервис для работы с jobs и сервер
	st := storage.NewJobsStorage()
	svc := service.NewJobsService(st, q.JobsCh)
	h := handler.NewHandler(svc)
	//запускаем сервер
	go h.Run()

	sem := make(chan struct{}, 5) //semafor
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

	<-ctx.Done()
	fmt.Println("\nshutdown signal received")
	cancel()

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
