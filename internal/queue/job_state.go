package queue

type State string

const (
	JobStateScheduled State = "scheduled"
	JobStateRunning   State = "running"
	JobStateCompleted State = "completed"
	JobStateFailed    State = "failed"
	JobStateCancelled State = "cancelled"
)

var fsm = map[State]map[State]bool{
	JobStateScheduled: {
		JobStateRunning:   true,
		JobStateCancelled: true,
	},
	JobStateRunning: {
		JobStateCompleted: true,
		JobStateFailed:    true,
		JobStateCancelled: true,
	},
	JobStateCompleted: {},
	JobStateFailed:    {},
	JobStateCancelled: {},
}

func (s State) CanTransition(to State) bool {
	return fsm[s][to]
}
