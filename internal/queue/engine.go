package queue

import "time"

type Engine struct {
	queue    *JobQueue
	pool     *WorkerPool
	handlers *HandlerRegistry

	jobs chan *Job
	wake chan struct{}
}

func NewEngine() *Engine {
	jobsChannel := make(chan *Job)

	engine := &Engine{
		queue:    NewJobQueue(),
		handlers: NewHandlerRegistry(),
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

func (e *Engine) Run(data JobDto) (*Job, error) {
	job := NewJob(data)
	e.queue.Push(job)

	select {
	case e.wake <- struct{}{}:
	default:
	}

	return job, nil
}

func (e *Engine) RegisterHandler(id string, handler Handler) {
	e.handlers.Register(id, handler)
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
	if !job.State.CanTransition(JobStateRunning) {
		return
	}

	handler, err := e.handlers.GetById(job.Type)

	if err != nil {
		_ = job.transitionTo(JobStateFailed)

		return
	}

	_ = job.transitionTo(JobStateRunning)
	job.StartedAt = time.Now()

	result, err := handler(job.Payload)

	if err != nil {
		_ = job.transitionTo(JobStateFailed)

		return
	}

	_ = job.transitionTo(JobStateCompleted)
	job.Result = result
	job.FinishedAt = time.Now()
}
