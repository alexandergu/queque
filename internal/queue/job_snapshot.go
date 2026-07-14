package queue

import (
	"time"

	"github.com/google/uuid"
)

type JobSnapshot struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	State    State     `json:"state"`
	Priority int       `json:"priority"`

	Payload []byte `json:"payload"`
	Result  []byte `json:"result"`
	Error   string `json:"error"`

	CreatedAt  time.Time `json:"createdAt"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}
