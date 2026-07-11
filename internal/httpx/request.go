package httpx

import (
	"encoding/json"
	"net/http"
)

type Validatable interface {
	Validate() error
}

func Convert[T Validatable](r *http.Request) (T, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, &ConvertError{"Convert error"}
	}

	if err := data.Validate(); err != nil {
		return data, &ValidationError{"Validation error"}
	}

	return data, nil
}
