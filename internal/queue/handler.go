package queue

import "context"

type Handler func(context.Context, []byte) ([]byte, error)
