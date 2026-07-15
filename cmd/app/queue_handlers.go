package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/alexandergu/queque/internal/queue"
)

var QueueHandlers = map[string]queue.Handler{
	"type1": func(bytes []byte) ([]byte, error) {
		payload := Payload{}

		if err := json.Unmarshal(bytes, &payload); err != nil {
			return nil, fmt.Errorf("payload error")
		}

		if err := payload.Validate(); err != nil {
			return nil, err
		}

		duration := time.Second * 1
		if payload.Duration != nil {
			duration = time.Duration(*payload.Duration) * time.Second
		}

		var failCh <-chan time.Time
		work := time.After(duration)

		if payload.FailChance != nil && rand.Float32() < *payload.FailChance {
			failAfter := time.Duration(rand.Int63n(int64(duration)))
			failCh = time.After(failAfter)
		}

		select {
		case <-failCh:
			return nil, fmt.Errorf("simulated failure")
		case <-work:
			return nil, nil
		}
	},
}
