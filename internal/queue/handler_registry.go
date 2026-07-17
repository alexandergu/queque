package queue

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
		return nil, &HandlerNotFoundError{id}
	}

	return handler, nil
}
