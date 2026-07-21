package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/goleak"
)

func TestEngine_RunCompleteJob(t *testing.T) {
	defer goleak.VerifyNone(t)

	engine := NewEngine()
	engine.RegisterHandler("testHandler", func(ctx context.Context, bytes []byte) ([]byte, error) {
		return nil, nil
	})

	engine.Start()
	defer engine.Stop()

	ch, unsubscribe := engine.Subscribe()
	defer unsubscribe()

	snapshot, err := engine.Run(JobDto{Type: "testHandler"})
	if err != nil {
		t.Fatalf("job cannot be run: %s", err.Error())
	}

	waitForEvent(t, ch, EventTypeJobCompleted)

	job, err := engine.GetJob(snapshot.ID)
	if err != nil {
		t.Fatalf("error during get the job by ID: %s", err.Error())
	}

	if job.State != JobStateCompleted {
		t.Errorf("incorrect job state %s, want %s", job.State, JobStateCompleted)
	}

	if job.StartedAt.IsZero() {
		t.Errorf("startedAt is zero")
	}

	if job.FinishedAt.IsZero() {
		t.Errorf("finishedAt is zero")
	}
}

func TestEngine_RunFailedJob(t *testing.T) {
	defer goleak.VerifyNone(t)

	errorMessage := "testError"
	engine := NewEngine()
	engine.RegisterHandler("testHandler", func(ctx context.Context, bytes []byte) ([]byte, error) {
		return nil, errors.New(errorMessage)
	})

	engine.Start()
	defer engine.Stop()

	ch, unsubscribe := engine.Subscribe()
	defer unsubscribe()

	snapshot, err := engine.Run(JobDto{Type: "testHandler"})
	if err != nil {
		t.Fatalf("job cannot be run: %s", err.Error())
	}

	waitForEvent(t, ch, EventTypeJobFailed)

	job, err := engine.GetJob(snapshot.ID)
	if err != nil {
		t.Fatalf("error during get the job by ID: %s", err.Error())
	}

	if job.State != JobStateFailed {
		t.Errorf("incorrect job state %s, want %s", job.State, JobStateFailed)
	}

	if job.Error != errorMessage {
		t.Errorf("incorrect error message %s, want %s", job.Error, errorMessage)
	}

	if job.StartedAt.IsZero() {
		t.Errorf("startedAt is zero")
	}

	if job.FinishedAt.IsZero() {
		t.Errorf("finishedAt is zero")
	}
}

func TestEngine_CancelJob(t *testing.T) {
	defer goleak.VerifyNone(t)

	engine := NewEngine()
	engine.RegisterHandler("testHandler", func(ctx context.Context, bytes []byte) ([]byte, error) {
		<-ctx.Done()

		return nil, ctx.Err()
	})

	engine.Start()
	defer engine.Stop()

	ch, unsubscribe := engine.Subscribe()
	defer unsubscribe()

	snapshot, err := engine.Run(JobDto{Type: "testHandler"})
	if err != nil {
		t.Fatalf("job cannot be run: %s", err.Error())
	}

	waitForEvent(t, ch, EventTypeJobRunning)
	_, err = engine.Cancel(snapshot.ID)
	if err != nil {
		t.Fatalf("job cannot be cancel: %s", err.Error())
	}
	waitForEvent(t, ch, EventTypeJobCancelled)

	job, err := engine.GetJob(snapshot.ID)
	if err != nil {
		t.Fatalf("error during get the job by ID: %s", err.Error())
	}

	if job.State != JobStateCancelled {
		t.Errorf("incorrect job state %s, want %s", job.State, JobStateCancelled)
	}

	if job.StartedAt.IsZero() {
		t.Errorf("startedAt is zero")
	}

	if job.FinishedAt.IsZero() {
		t.Errorf("finishedAt is zero")
	}
}

func TestEngine_RunUnknownHandler(t *testing.T) {
	defer goleak.VerifyNone(t)

	unknownHandlerID := "unknownHandler"
	engine := NewEngine()
	engine.Start()
	defer engine.Stop()

	ch, unsubscribe := engine.Subscribe()
	defer unsubscribe()

	snapshot, err := engine.Run(JobDto{Type: unknownHandlerID})
	if err != nil {
		t.Fatalf("job cannot be run: %s", err.Error())
	}

	waitForEvent(t, ch, EventTypeJobFailed)

	job, err := engine.GetJob(snapshot.ID)
	if err != nil {
		t.Fatalf("error during get the job by ID: %s", err.Error())
	}

	if job.State != JobStateFailed {
		t.Errorf("incorrect job state %s, want %s", job.State, JobStateFailed)
	}

	wantErr := &HandlerNotFoundError{ID: unknownHandlerID}
	if job.Error != wantErr.Error() {
		t.Errorf("unexpected error message: %s, want %s", job.Error, wantErr.Error())
	}

	if !job.StartedAt.IsZero() {
		t.Errorf("startedAt is not zero")
	}

	if job.FinishedAt.IsZero() {
		t.Errorf("finishedAt is zero")
	}
}

func TestEngine_GetJobs(t *testing.T) {
	defer goleak.VerifyNone(t)

	engine := NewEngine()
	engine.RegisterHandler("testHandler", func(ctx context.Context, payload []byte) ([]byte, error) {
		return nil, nil
	})

	engine.Start()
	defer engine.Stop()

	if jobs := engine.GetJobs(); len(jobs) != 0 {
		t.Fatalf("got %d jobs, want 0", len(jobs))
	}

	ch, unsubscribe := engine.Subscribe()
	defer unsubscribe()

	const count = 3
	ids := make([]uuid.UUID, 0, count)

	for i := range count {
		snapshot, err := engine.Run(JobDto{Type: "testHandler"})
		if err != nil {
			t.Fatalf("job %d cannot be run: %s", i, err.Error())
		}

		ids = append(ids, snapshot.ID)
	}

	for range count {
		waitForEvent(t, ch, EventTypeJobCompleted)
	}

	jobs := engine.GetJobs()
	if len(jobs) != count {
		t.Fatalf("got %d jobs, want %d", len(jobs), count)
	}

	for i, job := range jobs {
		if job.ID != ids[i] {
			t.Errorf("job at index %d has ID %s, want %s", i, job.ID, ids[i])
		}

		if job.State != JobStateCompleted {
			t.Errorf("job at index %d has state %s, want %s", i, job.State, JobStateCompleted)
		}

		if i > 0 && job.CreatedAt.Before(jobs[i-1].CreatedAt) {
			t.Errorf("jobs are not sorted by createdAt at index %d", i)
		}
	}
}

func TestEngine_GetJob(t *testing.T) {
	defer goleak.VerifyNone(t)

	payload := json.RawMessage(`{"value":42}`)
	result := json.RawMessage(`{"ok":true}`)

	engine := NewEngine()
	engine.RegisterHandler("testHandler", func(ctx context.Context, got []byte) ([]byte, error) {
		if !bytes.Equal(got, payload) {
			t.Errorf("got payload %s, want %s", got, payload)
		}

		return result, nil
	})

	engine.Start()
	defer engine.Stop()

	_, err := engine.GetJob(uuid.New())
	if err == nil {
		t.Fatal("expected an error for unknown job ID")
	}

	if _, ok := errors.AsType[*JobNotFoundError](err); !ok {
		t.Fatalf("unexpected error %T, want *JobNotFoundError", err)
	}

	ch, unsubscribe := engine.Subscribe()
	defer unsubscribe()

	snapshot, err := engine.Run(JobDto{
		Type:       "testHandler",
		Payload:    payload,
		Priority:   7,
		MaxAttempt: 3,
	})
	if err != nil {
		t.Fatalf("job cannot be run: %s", err.Error())
	}

	waitForEvent(t, ch, EventTypeJobCompleted)

	job, err := engine.GetJob(snapshot.ID)
	if err != nil {
		t.Fatalf("error during get the job by ID: %s", err.Error())
	}

	if job.ID != snapshot.ID {
		t.Errorf("incorrect job ID %s, want %s", job.ID, snapshot.ID)
	}

	if job.Type != "testHandler" {
		t.Errorf("incorrect job type %s, want testHandler", job.Type)
	}

	if job.Priority != 7 {
		t.Errorf("incorrect priority %d, want 7", job.Priority)
	}

	if job.MaxAttempt != 3 {
		t.Errorf("incorrect maxAttempt %d, want 3", job.MaxAttempt)
	}

	if job.Attempt != 1 {
		t.Errorf("incorrect attempt %d, want 1", job.Attempt)
	}

	if !bytes.Equal(job.Payload, payload) {
		t.Errorf("incorrect payload %s, want %s", job.Payload, payload)
	}

	if !bytes.Equal(job.Result, result) {
		t.Errorf("incorrect result %s, want %s", job.Result, result)
	}
}

func waitForEvent(t *testing.T, ch <-chan Event, wantEvent EventType) {
	t.Helper()

	for {
		select {
		case event := <-ch:
			if event.Type == wantEvent {
				return
			}

			if event.Type == EventTypeJobCancelled || event.Type == EventTypeJobCompleted || event.Type == EventTypeJobFailed {
				t.Fatalf("got terminal event type %s, want %s", event.Type, wantEvent)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout, waiting for event %s", wantEvent)
		}
	}
}
