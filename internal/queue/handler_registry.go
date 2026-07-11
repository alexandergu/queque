package queue

import "fmt"

type HandlerRegistry struct {
	registry map[string]Handler
}

func newHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		registry: make(map[string]Handler),
	}
}

func (r *HandlerRegistry) Register(id string, h Handler) {
	r.registry[id] = h
}

func (r *HandlerRegistry) GetById(id string) (Handler, error) {
	handler, ok := r.registry[id]

	if !ok {
		return nil, fmt.Errorf("not exist handler")
	}

	return handler, nil
}
