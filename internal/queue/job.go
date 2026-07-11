package queue

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID       uuid.UUID
	Type     string
	State    State
	Priority int

	Payload []byte
	Result  []byte

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

func (j *Job) transitionTo(s State) error {
	if !j.State.CanTransition(s) {
		return fmt.Errorf("state transition error")
	}

	j.State = s

	return nil
}
