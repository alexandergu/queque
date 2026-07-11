package api

import (
	"fmt"
	"net/http"

	"github.com/alexandergu/queque/internal/queue"
	"github.com/google/uuid"
)

func RenderJob(job queue.JobSnapshot) any {
	return job
}

func GetJobFromRequest(registry *queue.Engine, r *http.Request) (queue.JobSnapshot, error) {
	id := r.PathValue("id")

	if id == "" {
		return queue.JobSnapshot{}, fmt.Errorf("job not found")
	}

	ID, err := uuid.Parse(id)
	if err != nil {
		return queue.JobSnapshot{}, fmt.Errorf("job not found")
	}

	job, err := registry.GetJob(ID)
	if err != nil {
		return queue.JobSnapshot{}, fmt.Errorf("job not found")
	}

	return job, nil
}
