package queue

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobSnapshot struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	State    State     `json:"state"`
	Priority int       `json:"priority"`

	Payload json.RawMessage `json:"payload"`
	Result  json.RawMessage `json:"result"`
	Error   string          `json:"error"`

	CreatedAt  time.Time `json:"createdAt"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}
