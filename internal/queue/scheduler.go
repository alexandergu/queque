package queue

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	mu      sync.Mutex
	heap    JobDelayedHeap
	promote func(*Job)

	wake   chan struct{}
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func newScheduler(promote func(*Job)) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		wake:    make(chan struct{}, 1),
		promote: promote,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
}

func (s *Scheduler) Schedule(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	heap.Push(&s.heap, job)

	select {
	case s.wake <- struct{}{}:
	default:
	}
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	for {
		s.mu.Lock()

		if s.heap.Len() == 0 {
			s.mu.Unlock()

			select {
			case <-s.ctx.Done():
				return
			case <-s.wake:
				continue
			}
		}

		now := time.Now()
		next := s.heap[0]

		if next.AvailableAt.Before(now) {
			heap.Pop(&s.heap)
			s.mu.Unlock()
			s.promote(next)

			continue
		}

		s.mu.Unlock()
		delay := next.AvailableAt.Sub(now)
		timer := time.NewTimer(delay)

		select {
		case <-s.ctx.Done():
			timer.Stop()

			return
		case <-timer.C:
			continue
		case <-s.wake:
			timer.Stop()

			continue
		}

	}
}
