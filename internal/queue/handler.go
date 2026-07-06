package queue

type Handler func([]byte) ([]byte, error)
