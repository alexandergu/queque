package queue

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type ExecutionRegistry struct {
	mu       sync.Mutex
	registry map[uuid.UUID]context.CancelFunc
}

func newExecutionRegistry() *ExecutionRegistry {
	return &ExecutionRegistry{
		registry: make(map[uuid.UUID]context.CancelFunc),
	}
}

func (r *ExecutionRegistry) Add(id uuid.UUID, cancelFunc context.CancelFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registry[id] = cancelFunc
}

func (r *ExecutionRegistry) Remove(id uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.registry, id)
}

func (r *ExecutionRegistry) Cancel(id uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cancelFunc, ok := r.registry[id]; ok {
		cancelFunc()
	}
}
