package queue

import "time"

func generateQuadraticDelay(attempts int) time.Duration {
	return time.Duration(attempts*attempts) * time.Second
}
