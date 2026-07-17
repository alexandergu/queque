package api

import (
	"fmt"

	"github.com/alexandergu/queque/internal/httpx"
)

type WorkersCount struct {
	Count int `json:"count"`
}

type ResizeWorkersDto struct {
	Count int
}

func (dto ResizeWorkersDto) Validate() *httpx.ValidationError {
	if dto.Count < 0 || dto.Count > 10 {
		return &httpx.ValidationError{
			Message: "resize workers payload validation error",
			Errors: []httpx.Violation{
				{
					Path:    "Count",
					Message: fmt.Sprintf("count must be between 0 and 10, got %d", dto.Count),
				},
			},
		}
	}

	return nil
}
