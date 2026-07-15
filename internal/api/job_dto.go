package api

import (
	"encoding/json"
	"fmt"

	"github.com/alexandergu/queque/internal/queue"
)

type CreateJobData struct {
	Type     string
	Payload  json.RawMessage
	Priority int
}

func (dto CreateJobData) Validate() error {
	if dto.Priority <= 0 {
		return fmt.Errorf("priority must be 1 or more")
	}

	if dto.Type == "" {
		return fmt.Errorf("type field is mandatory")
	}

	if len(dto.Payload) == 0 {
		return fmt.Errorf("payload is mandatory")
	}

	return nil
}

func (dto CreateJobData) ToJobDto() queue.JobDto {
	return queue.JobDto{
		Type:     dto.Type,
		Payload:  dto.Payload,
		Priority: dto.Priority,
	}
}
