package httpx

import (
	"encoding/json"
	"net/http"
)

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

func NotFound(w http.ResponseWriter) {
	writeJson(w, http.StatusNotFound, struct{}{})
}

func Error(w http.ResponseWriter, err error) {
	writeJson(w, http.StatusInternalServerError, struct{}{})
}
