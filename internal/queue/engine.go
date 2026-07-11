package queue

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type Engine struct {
	registry *JobRegistry
	queue    *JobQueue
	pool     *WorkerPool
	handlers *HandlerRegistry

	jobs   chan *Job
	wake   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewEngine() *Engine {
	jobsChannel := make(chan *Job)
	ctx, cancel := context.WithCancel(context.Background())

	engine := &Engine{
		registry: newJobRegistry(),
		queue:    newJobQueue(),
		handlers: newHandlerRegistry(),
		jobs:     jobsChannel,
		wake:     make(chan struct{}, 1),
		ctx:      ctx,
		cancel:   cancel,
	}
	engine.pool = NewWorkerPool(jobsChannel, engine.process)

	return engine
}

func (e *Engine) Start() {
	e.wg.Add(1)
	go e.dispatch()
	e.pool.Start(1)
}

func (e *Engine) Run(data JobDto) (JobSnapshot, error) {
	if e.ctx.Err() != nil {
		return JobSnapshot{}, fmt.Errorf("engine has already stopped")
	}

	job := NewJob(data)

	e.registry.Add(job)
	e.queue.Push(job)

	select {
	case e.wake <- struct{}{}:
	default:
	}

	return job.toSnapshot(), nil
}

func (e *Engine) Stop() error {
	e.cancel()
	e.wg.Wait()
	e.pool.Stop()

	return nil
}

func (e *Engine) RegisterHandler(id string, handler Handler) {
	e.handlers.Register(id, handler)
}

func (e *Engine) GetJobs() []JobSnapshot {
	jobs := e.registry.All()
	snapshots := make([]JobSnapshot, 0, len(jobs))

	for _, job := range jobs {
		snapshots = append(snapshots, job.toSnapshot())
	}

	slices.SortFunc(snapshots, func(a, b JobSnapshot) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	return snapshots
}

func (e *Engine) GetJob(id uuid.UUID) (JobSnapshot, error) {
	job, err := e.registry.Get(id)
	if err != nil {
		return JobSnapshot{}, err
	}

	return job.toSnapshot(), nil
}

func (e *Engine) dispatch() {
	defer e.wg.Done()

	for {
		job, ok := e.queue.Pop()

		if !ok {
			select {
			case <-e.wake:
				continue
			case <-e.ctx.Done():
				return
			}
		}

		select {
		case <-e.ctx.Done():
			e.queue.Push(job)

			return
		case e.jobs <- job:
			continue
		}
	}
}

func (e *Engine) process(job *Job) {
	handler, err := e.handlers.GetById(job.Type)
	if err != nil {
		_ = job.fail()

		return
	}

	err = job.run()
	if err != nil {
		_ = job.fail()

		return
	}

	result, err := handler(job.Payload)
	if err != nil {
		_ = job.fail()

		return
	}

	_ = job.complete(result)
}
