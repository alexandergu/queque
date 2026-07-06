package queue

import "fmt"

type Engine struct {
	queue *JobQueue
	pool  *WorkerPool

	jobs chan *Job
	wake chan struct{}
}

func NewEngine() *Engine {
	jobsChannel := make(chan *Job)
	workerPool := NewWorkerPool(jobsChannel, func(job *Job) {
		fmt.Println("process callback", job.ID)
	})

	return &Engine{
		queue: NewJobQueue(),
		pool:  workerPool,
		jobs:  jobsChannel,
		wake:  make(chan struct{}, 1),
	}
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
