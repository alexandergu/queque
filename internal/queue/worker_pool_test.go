package queue

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/goleak"
)

func TestWorkerPool_CreateWorkers(t *testing.T) {
	defer goleak.VerifyNone(t)

	jobs := make(chan *Job)
	process := func(ctx context.Context, job *Job) {}

	pool := NewWorkerPool(jobs, process)
	defer pool.Stop()

	pool.Start(3)

	if length := pool.Len(); length != 3 {
		t.Errorf("Length is %d, want %d", length, 3)
	}
}

func TestWorkerPool_ResizeWorkers(t *testing.T) {
	testCases := []struct {
		name   string
		start  int
		resize int
		want   int
	}{
		{"grow", 1, 2, 2},
		{"shrink", 5, 2, 2},
		{"to zero", 5, 0, 0},
		{"from zero", 0, 2, 2},
		{"to negative", 5, -2, 0},
		{"equal", 2, 2, 2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			jobs := make(chan *Job)
			process := func(ctx context.Context, job *Job) {}

			pool := NewWorkerPool(jobs, process)
			defer pool.Stop()

			pool.Start(testCase.start)
			pool.Resize(testCase.resize)

			if length := pool.Len(); length != testCase.want {
				t.Errorf("Length after resize is %d, want %d", length, testCase.want)
			}
		})
	}
}

func TestWorkerPool_Process(t *testing.T) {
	defer goleak.VerifyNone(t)

	count := 5
	wantIds := make(map[uuid.UUID]struct{}, count)
	processed := make(chan uuid.UUID, count)

	jobs := make(chan *Job)
	process := func(ctx context.Context, job *Job) {
		processed <- job.ID
	}

	pool := NewWorkerPool(jobs, process)
	defer pool.Stop()

	pool.Start(2)

	for range count {
		job := NewJob(JobDto{})
		wantIds[job.ID] = struct{}{}
		jobs <- job
	}

	for range count {
		select {
		case id := <-processed:
			if _, ok := wantIds[id]; !ok {
				t.Errorf("unexpected job ID %s", id)
			}

			delete(wantIds, id)

		case <-time.After(time.Second):
			t.Fatalf("timeout, %d jobs unprocessed", len(wantIds))
		}
	}
}
