package queue

import "fmt"

type WorkerNotFoundError struct {
	ID string
}

func (err *WorkerNotFoundError) Error() string {
	return fmt.Sprintf("worker id: %s not found", err.ID)
}
