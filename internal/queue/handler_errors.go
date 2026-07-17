package queue

import "fmt"

type HandlerNotFoundError struct {
	ID string
}

func (err *HandlerNotFoundError) Error() string {
	return fmt.Sprintf("handler id: %s not found", err.ID)
}
