package api

import (
	"encoding/json"
	"fmt"

	"github.com/alexandergu/queque/internal/httpx"
	"github.com/alexandergu/queque/internal/queue"
)

type CreateJobData struct {
	Type       string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
	Priority   int             `json:"priority"`
	MaxAttempt int             `json:"maxAttempt"`
}

func (dto CreateJobData) Validate() error {
	var violations []httpx.Violation

	if dto.Priority <= 0 {
		violations = append(violations, httpx.Violation{
			Path:    "priority",
			Message: fmt.Sprintf("priority must be at least 1, got %d", dto.Priority),
		})
	}

	if dto.Type == "" {
		violations = append(violations, httpx.Violation{
			Path:    "type",
			Message: "type is required",
		})
	}

	if len(dto.Payload) == 0 {
		violations = append(violations, httpx.Violation{
			Path:    "payload",
			Message: "payload is required",
		})
	}

	if dto.MaxAttempt < 1 || dto.MaxAttempt > 10 {
		violations = append(violations, httpx.Violation{
			Path:    "maxAttempt",
			Message: fmt.Sprintf("maxAttempt must be between 1 and 10"),
		})
	}

	if len(violations) > 0 {
		return &httpx.ValidationError{
			Message: "validation failed",
			Errors:  violations,
		}
	}

	return nil
}

func (dto CreateJobData) ToJobDto() queue.JobDto {
	return queue.JobDto{
		Type:       dto.Type,
		Payload:    dto.Payload,
		Priority:   dto.Priority,
		MaxAttempt: dto.MaxAttempt,
	}
}
