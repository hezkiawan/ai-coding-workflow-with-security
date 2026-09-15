package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	JSON(w, status, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, errMessage string) {
	JSON(w, status, APIResponse{
		Success: false,
		Error:   errMessage,
	})
}
