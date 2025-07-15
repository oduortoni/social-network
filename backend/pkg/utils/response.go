package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message,omitempty"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
