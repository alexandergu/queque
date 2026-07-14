package queue

type EventType string

type Event struct {
	Type EventType
	Job  JobSnapshot
}

const (
	EventTypeJobScheduled EventType = "job.scheduled"
	EventTypeJobRunning   EventType = "job.running"
	EventTypeJobFailed    EventType = "job.failed"
	EventTypeJobCompleted EventType = "job.completed"
	EventTypeJobCancelled EventType = "job.cancelled"
)
