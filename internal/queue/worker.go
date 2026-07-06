package queue

type Worker struct {
	ID   string
	quit chan struct{}
}
