package queue

type State string

const (
	JobStateScheduled State = "scheduled"
	JobStateRetrying  State = "retrying"
	JobStateRunning   State = "running"
	JobStateCompleted State = "completed"
	JobStateFailed    State = "failed"
	JobStateCancelled State = "cancelled"
)

var fsm = map[State]map[State]bool{
	JobStateScheduled: {
		JobStateRunning:   true,
		JobStateCancelled: true,
		JobStateFailed:    true,
	},
	JobStateRunning: {
		JobStateCompleted: true,
		JobStateFailed:    true,
		JobStateCancelled: true,
		JobStateRetrying:  true,
	},
	JobStateRetrying: {
		JobStateRunning:   true,
		JobStateScheduled: true,
		JobStateCancelled: true,
	},
	JobStateCompleted: {},
	JobStateCancelled: {},
}

func (s State) CanTransition(to State) bool {
	return fsm[s][to]
}
