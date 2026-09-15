package handlers

import (
	"net/http"
	"time"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"version": "1.0.0",
	}, "Health check passed")
}

func (h *SystemHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, http.StatusOK, map[string]interface{}{
		"system_uptime": time.Since(time.Unix(0, 0)).String(),
	}, "System metrics")
}