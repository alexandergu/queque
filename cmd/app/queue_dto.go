package main

import "fmt"

type Payload struct {
	FailChance *float32 `json:"failChance"`
	Duration   *int64   `json:"duration"`
}

func (p Payload) Validate() error {
	if p.Duration != nil && (*p.Duration <= 0 || *p.Duration > 300) {
		return fmt.Errorf("payload error: duration must be more 0 and less than 5 minutes")
	}

	if p.FailChance != nil && (*p.FailChance < 0 || *p.FailChance > 1) {
		return fmt.Errorf("payload error: fail Chance must be between 0 and 1")
	}

	return nil
}
