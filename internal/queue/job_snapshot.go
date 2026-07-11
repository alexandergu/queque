package queue

import (
	"time"

	"github.com/google/uuid"
)

type JobSnapshot struct {
	ID       uuid.UUID
	Type     string
	State    State
	Priority int

	Payload []byte
	Result  []byte

	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
}
