package queue

import "sync"

type EventBus struct {
	mu   sync.Mutex
	subs map[chan Event]struct{}
}

func newEventBus() *EventBus {
	return &EventBus{
		subs: make(map[chan Event]struct{}),
	}
}

func (bus *EventBus) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 16)

	bus.mu.Lock()
	bus.subs[ch] = struct{}{}
	bus.mu.Unlock()

	unsubscribe := func() {
		bus.mu.Lock()
		defer bus.mu.Unlock()

		if _, ok := bus.subs[ch]; ok {
			delete(bus.subs, ch)
			close(ch)
		}
	}

	return ch, unsubscribe
}

func (bus *EventBus) Publish(event Event) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	for ch := range bus.subs {
		select {
		case ch <- event:
		default:
		}
	}
}

func (bus *EventBus) Close() {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	for ch := range bus.subs {
		delete(bus.subs, ch)
		close(ch)
	}
}
