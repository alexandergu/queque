package api

import "fmt"

type ResizeWorkersDto struct {
	Count int
}

func (dto ResizeWorkersDto) Validate() error {
	if dto.Count < 0 || dto.Count > 10 {
		return fmt.Errorf("count must be more than 0 and less then 10")
	}

	return nil
}
