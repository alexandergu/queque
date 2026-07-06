package queue

type JobDto struct {
	Type     string
	Payload  []byte
	Priority int
}
