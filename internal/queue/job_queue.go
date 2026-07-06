package queue

import (
	"container/heap"
	"sync"
)

type JobQueue struct {
	mu   sync.RWMutex
	heap JobHeap
}

func NewJobQueue() *JobQueue {
	return &JobQueue{}
}

func (queue *JobQueue) Push(job *Job) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	heap.Push(&queue.heap, job)
}

func (queue *JobQueue) Pop() (*Job, bool) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.heap.Len() == 0 {
		return nil, false
	}

	return heap.Pop(&queue.heap).(*Job), true
}

func (queue *JobQueue) Len() int {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	return queue.heap.Len()
}
