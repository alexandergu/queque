package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Errors  []Violation `json:"errors,omitempty"`
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Ok(w http.ResponseWriter) {
	writeJson(w, http.StatusOK, struct{}{})
}

func Resources[T any](w http.ResponseWriter, items []T, renderer func(T) any) {
	result := make([]any, len(items))

	for i, item := range items {
		result[i] = renderer(item)
	}

	writeJson(w, http.StatusOK, result)
}

func Resource[T any](w http.ResponseWriter, item T, renderer func(T) any) {
	result := renderer(item)

	writeJson(w, http.StatusOK, result)
}

func Error(w http.ResponseWriter, err error) {
	if _, ok := errors.AsType[*NotFoundError](err); ok {
		writeJson(w, http.StatusNotFound, errorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})

		return
	}

	if validationError, ok := errors.AsType[*ValidationError](err); ok {
		writeJson(w, http.StatusBadRequest, errorResponse{
			Error:   "validation_error",
			Message: validationError.Message,
			Errors:  validationError.Errors,
		})

		return
	}

	if _, ok := errors.AsType[*BadRequestError](err); ok {
		writeJson(w, http.StatusBadRequest, errorResponse{
			Error:   "bad_request",
			Message: err.Error(),
		})

		return
	}

	if _, ok := errors.AsType[*ConflictError](err); ok {
		writeJson(w, http.StatusConflict, errorResponse{
			Error:   "conflict",
			Message: err.Error(),
		})

		return
	}

	writeJson(w, http.StatusInternalServerError, errorResponse{
		Error:   "internal_error",
		Message: err.Error(),
	})
}
