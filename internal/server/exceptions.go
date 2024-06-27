package server

import (
	"encoding/json"
	"net/http"
)

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func NotFoundError(w http.ResponseWriter, item string) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": item + " not found"})
	return
}

func BadRequestError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request payload"})
}

func UnauthorizedError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func ForbiddenError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
}
