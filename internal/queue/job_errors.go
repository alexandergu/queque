package queue

import (
	"fmt"

	"github.com/google/uuid"
)

type JobTransitionError struct {
	From State
	To   State
}

func (err *JobTransitionError) Error() string {
	return fmt.Sprintf("job transition error from %s to %s", err.From, err.To)
}

type JobNotFoundError struct {
	id uuid.UUID
}

func (err *JobNotFoundError) Error() string {
	return fmt.Sprintf("job ID: %s not found", err.id)
}
