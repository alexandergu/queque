package queue

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type JobRegistry struct {
	mu   sync.Mutex
	jobs map[uuid.UUID]*Job
}

func newJobRegistry() *JobRegistry {
	return &JobRegistry{jobs: make(map[uuid.UUID]*Job)}
}

func (registry *JobRegistry) Add(job *Job) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.jobs[job.ID] = job
}

func (registry *JobRegistry) Get(id uuid.UUID) (*Job, error) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	job, ok := registry.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job %s not found", id)
	}

	return job, nil
}

func (registry *JobRegistry) All() []*Job {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	list := make([]*Job, 0, len(registry.jobs))
	for _, job := range registry.jobs {
		list = append(list, job)
	}

	return list
}
