package api

import (
	"fmt"

	"github.com/alexandergu/queque/internal/httpx"
)

type WorkersCount struct {
	Count int `json:"count"`
}

type ResizeWorkersData struct {
	Count int `json:"count"`
}

func (dto ResizeWorkersData) Validate() error {
	if dto.Count < 0 || dto.Count > 10 {
		return &httpx.ValidationError{
			Message: "validation failed",
			Errors: []httpx.Violation{
				{
					Path:    "count",
					Message: fmt.Sprintf("count must be between 0 and 10, got %d", dto.Count),
				},
			},
		}
	}

	return nil
}
