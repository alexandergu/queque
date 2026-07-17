package queue

import "fmt"

type EngineRunError struct {
	Reason string
}

func (err *EngineRunError) Error() string {
	return fmt.Sprintf("engine can not run due to: %s", err.Reason)
}
