package handlers

import (
	"net/http"
	"runtime"
	"time"

	"opsdesk/internal/config"
	"opsdesk/internal/database"
	"opsdesk/internal/utils"
)

type SystemHandler struct {
	cfg       *config.Config
	startTime time.Time
}

func NewSystemHandler(cfg *config.Config) *SystemHandler {
	return &SystemHandler{
		cfg:       cfg,
		startTime: time.Now(),
	}
}

func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbStatus := "healthy"
	if err := database.DB.Ping(); err != nil {
		dbStatus = "unreachable"
	}

	utils.Success(w, http.StatusOK, map[string]interface{}{
		"status":      "online",
		"database":    dbStatus,
		"uptime":      time.Since(h.startTime).String(),
		"environment": h.cfg.Environment,
		"timestamp":   time.Now().UTC(),
	})
}

func (h *SystemHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var ticketCount, userCount, assetCount int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM tickets").Scan(&ticketCount)
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM assets").Scan(&assetCount)

	utils.Success(w, http.StatusOK, map[string]interface{}{
		"memory_alloc_mb":   m.Alloc / 1024 / 1024,
		"memory_total_mb":   m.TotalAlloc / 1024 / 1024,
		"goroutines":        runtime.NumGoroutine(),
		"num_cpu":           runtime.NumCPU(),
		"active_tickets":    ticketCount,
		"registered_users":  userCount,
		"registered_assets": assetCount,
	})
}
