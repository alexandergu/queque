package queue

import "context"

type Worker struct {
	ID     string
	cancel context.CancelFunc
}
