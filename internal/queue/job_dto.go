package queue

import "encoding/json"

type JobDto struct {
	Type       string
	Payload    json.RawMessage
	Priority   int
	MaxAttempt int
}
