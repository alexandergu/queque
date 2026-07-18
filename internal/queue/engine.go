package queue

import (
	"context"
	"errors"
	"log"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type Engine struct {
	registry   *JobRegistry
	queue      *JobQueue
	pool       *WorkerPool
	handlers   *HandlerRegistry
	eventBus   *EventBus
	executions *ExecutionRegistry

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
		registry:   newJobRegistry(),
		queue:      newJobQueue(),
		handlers:   newHandlerRegistry(),
		eventBus:   newEventBus(),
		executions: newExecutionRegistry(),
		jobs:       jobsChannel,
		wake:       make(chan struct{}, 1),
		ctx:        ctx,
		cancel:     cancel,
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
		return JobSnapshot{}, &EngineRunError{"engine has already stopped"}
	}

	job := NewJob(data)

	e.registry.Add(job)
	e.queue.Push(job)
	e.eventBus.Publish(Event{Type: EventTypeJobScheduled, Job: job.toSnapshot()})

	select {
	case e.wake <- struct{}{}:
	default:
	}

	return job.toSnapshot(), nil
}

func (e *Engine) Cancel(id uuid.UUID) (JobSnapshot, error) {
	job, err := e.registry.Get(id)
	if err != nil {
		return JobSnapshot{}, err
	}

	err = job.cancel()
	if err != nil {
		return JobSnapshot{}, err
	}

	e.executions.Cancel(id)
	e.eventBus.Publish(Event{EventTypeJobCancelled, job.toSnapshot()})

	return job.toSnapshot(), nil
}

func (e *Engine) Stop() error {
	e.cancel()
	e.wg.Wait()
	e.pool.Stop()
	e.eventBus.Close()

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

func (e *Engine) ResizeWorkersCount(count int) error {
	e.pool.Resize(count)

	return nil
}

func (e *Engine) WorkersCount() int {
	return e.pool.Len()
}

func (e *Engine) Subscribe() (<-chan Event, func()) {
	return e.eventBus.Subscribe()
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

func (e *Engine) process(workerCtx context.Context, job *Job) {
	handler, reason := e.handlers.GetById(job.Type)
	if reason != nil {
		if err := job.fail(reason); err != nil {
			log.Printf("job %s failed to mark as failed", job.ID)

			return
		}

		e.eventBus.Publish(Event{EventTypeJobFailed, job.toSnapshot()})

		return
	}

	ctx, cancel := context.WithCancel(workerCtx)
	defer cancel()

	e.executions.Add(job.ID, cancel)
	defer e.executions.Remove(job.ID)

	reason = job.run()
	if reason != nil {
		if job.toSnapshot().State == JobStateCancelled {
			return
		}

		if err := job.fail(reason); err != nil {
			log.Printf("job %s failed to mark as failed", job.ID)

			return
		}

		e.eventBus.Publish(Event{EventTypeJobFailed, job.toSnapshot()})

		return
	}

	e.eventBus.Publish(Event{EventTypeJobRunning, job.toSnapshot()})
	result, reason := handler(ctx, job.Payload)

	if reason != nil {
		if errors.Is(reason, context.Canceled) {
			if err := job.cancel(); err == nil {
				e.eventBus.Publish(Event{EventTypeJobCancelled, job.toSnapshot()})
			}

			return
		}

		if err := job.fail(reason); err != nil {
			log.Printf("job %s failed to mark as failed", job.ID)

			return
		}

		e.eventBus.Publish(Event{EventTypeJobFailed, job.toSnapshot()})

		return
	}

	reason = job.complete(result)
	if reason != nil {
		log.Printf("job %s failed to mark as completed", job.ID)

		return
	}

	e.eventBus.Publish(Event{EventTypeJobCompleted, job.toSnapshot()})
}
