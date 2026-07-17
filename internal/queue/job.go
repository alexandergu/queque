package queue

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	mu sync.Mutex

	ID       uuid.UUID
	Type     string
	State    State
	Priority int

	Payload json.RawMessage
	Result  json.RawMessage
	Error   string

	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time

	seq int64
}

func NewJob(dto JobDto) *Job {
	return &Job{
		ID:        uuid.New(),
		Type:      dto.Type,
		State:     JobStateScheduled,
		Priority:  dto.Priority,
		Payload:   dto.Payload,
		CreatedAt: time.Now(),
	}
}

func (j *Job) run() error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if !j.State.CanTransition(JobStateRunning) {
		return &JobTransitionError{j.State, JobStateRunning}
	}

	j.State = JobStateRunning
	j.StartedAt = time.Now()

	return nil
}

func (j *Job) complete(result []byte) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if !j.State.CanTransition(JobStateCompleted) {
		return &JobTransitionError{j.State, JobStateCompleted}
	}

	j.State = JobStateCompleted
	j.FinishedAt = time.Now()
	j.Result = result

	return nil
}

func (j *Job) fail(reason error) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if !j.State.CanTransition(JobStateFailed) {
		return &JobTransitionError{j.State, JobStateFailed}
	}

	j.State = JobStateFailed
	j.FinishedAt = time.Now()
	j.Error = reason.Error()

	return nil
}

func (j *Job) cancel() error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if !j.State.CanTransition(JobStateCancelled) {
		return &JobTransitionError{j.State, JobStateCancelled}
	}

	j.State = JobStateCancelled
	j.FinishedAt = time.Now()

	return nil
}

func (j *Job) toSnapshot() JobSnapshot {
	j.mu.Lock()
	defer j.mu.Unlock()

	return JobSnapshot{
		ID:         j.ID,
		Type:       j.Type,
		State:      j.State,
		Priority:   j.Priority,
		Payload:    j.Payload,
		Result:     j.Result,
		Error:      j.Error,
		CreatedAt:  j.CreatedAt,
		StartedAt:  j.StartedAt,
		FinishedAt: j.FinishedAt,
	}
}
