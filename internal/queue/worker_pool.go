package queue

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type WorkerPool struct {
	mu      sync.Mutex
	workers []*Worker

	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	stopped bool

	jobs    <-chan *Job
	process func(*Job)
}

func NewWorkerPool(jobs <-chan *Job, process func(*Job)) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		jobs:    jobs,
		process: process,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (pool *WorkerPool) Start(count int) {
	pool.Resize(count)
}

func (pool *WorkerPool) Stop() {
	pool.mu.Lock()
	pool.stopped = true
	pool.mu.Unlock()

	pool.cancel()
	pool.wg.Wait()
}

func (pool *WorkerPool) Len() int {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return len(pool.workers)
}

func (pool *WorkerPool) Resize(count int) {
	if count < 0 {
		count = 0
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.stopped {
		return
	}

	currentLength := len(pool.workers)

	if currentLength < count {
		for i := currentLength; i < count; i++ {
			w := &Worker{
				ID:   uuid.NewString(),
				quit: make(chan struct{}),
			}

			pool.workers = append(pool.workers, w)

			pool.wg.Add(1)
			go pool.runWorker(w)
		}
	} else if currentLength > count {
		for i := count; i < currentLength; i++ {
			close(pool.workers[i].quit)
		}

		pool.workers = slices.Delete(pool.workers, count, currentLength)
	}
}

func (pool *WorkerPool) runWorker(worker *Worker) {
	defer pool.wg.Done()

	for {
		select {
		case <-pool.ctx.Done():
			return
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

func (pool *WorkerPool) stopWorker(id string) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for i, worker := range pool.workers {
		if worker.ID != id {
			continue
		}

		close(worker.quit)
		pool.workers = slices.Delete(pool.workers, i, i+1)

		return nil
	}

	return &WorkerNotFoundError{id}
}
