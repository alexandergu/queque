package queue

import (
	"sync"

	"github.com/google/uuid"
)

type WorkerPool struct {
	mu      sync.Mutex
	workers []*Worker

	jobs    <-chan *Job
	process func(*Job)
}

func NewWorkerPool(jobs <-chan *Job, process func(*Job)) *WorkerPool {
	return &WorkerPool{
		jobs:    jobs,
		process: process,
	}
}

func (pool *WorkerPool) Start(count int) {
	pool.Resize(count)
}

func (pool *WorkerPool) Stop() {}

func (pool *WorkerPool) Resize(count int) {
	if count < 0 {
		count = 0
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	currentLength := len(pool.workers)

	if currentLength < count {
		for i := currentLength; i < count; i++ {
			w := &Worker{
				ID:   uuid.NewString(),
				quit: make(chan struct{}),
			}

			pool.workers = append(pool.workers, w)
			go pool.runWorker(w)
		}
	} else if currentLength > count {
		for i := count; i < currentLength; i++ {
			close(pool.workers[i].quit)
		}

		pool.workers = pool.workers[:count]
	}
}

func (pool *WorkerPool) runWorker(worker *Worker) {
	for {
		select {
		case <-worker.quit:
			return

		case job, ok := <-pool.jobs:
			if !ok {
				return
			}

			pool.process(job)
		}
	}
}
