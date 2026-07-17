package api

import (
	"fmt"
	"net/http"

	"github.com/alexandergu/queque/internal/httpx"
	"github.com/alexandergu/queque/internal/queue"
	"github.com/google/uuid"
)

func RenderJob(job queue.JobSnapshot) any {
	return job
}

func GetJobFromRequest(registry *queue.Engine, r *http.Request) (queue.JobSnapshot, error) {
	id := r.PathValue("id")

	if id == "" {
		return queue.JobSnapshot{}, &httpx.ConvertError{Message: fmt.Sprintf("job id is empty")}
	}

	ID, err := uuid.Parse(id)
	if err != nil {
		return queue.JobSnapshot{}, &httpx.ConvertError{Message: fmt.Sprintf("invalid job id")}
	}

	job, err := registry.GetJob(ID)
	if err != nil {
		return queue.JobSnapshot{}, err
	}

	return job, nil
}
