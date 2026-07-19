package queue

import (
	"context"
	"errors"
	"testing"
	"time"

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
