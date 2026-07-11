package queue

import (
	"slices"

	"github.com/google/uuid"
)

type Engine struct {
	registry *JobRegistry
	queue    *JobQueue
	pool     *WorkerPool
	handlers *HandlerRegistry

	jobs chan *Job
	wake chan struct{}
}

func NewEngine() *Engine {
	jobsChannel := make(chan *Job)

	engine := &Engine{
		registry: newJobRegistry(),
		queue:    newJobQueue(),
		handlers: newHandlerRegistry(),
		jobs:     jobsChannel,
		wake:     make(chan struct{}, 1),
	}
	engine.pool = NewWorkerPool(jobsChannel, engine.process)

	return engine
}

func (e *Engine) Start() {
	go e.dispatch()
	e.pool.Start(1)
}

func (e *Engine) Run(data JobDto) (JobSnapshot, error) {
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
	for {
		job, ok := e.queue.Pop()

		if !ok {
			<-e.wake

			continue
		}

		e.jobs <- job
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
