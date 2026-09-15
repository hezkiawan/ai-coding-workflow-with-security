package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type DiagnosticsHandler struct{}

func NewDiagnosticsHandler() *DiagnosticsHandler {
	return &DiagnosticsHandler{}
}

type PingRequest struct {
	Host  string `json:"host"`
	Count int    `json:"count"`
}

type PingResponse struct {
	Host     string   `json:"host"`
	Resolved []string `json:"resolved_ips"`
	Status   string   `json:"status"`
	Latency  int64    `json:"latency_ms"`
}

type PortCheckRequest struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout_seconds"`
}

type PortCheckResponse struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Open   bool   `json:"open"`
	Banner string `json:"banner,omitempty"`
}

func (h *DiagnosticsHandler) PingAsset(w http.ResponseWriter, r *http.Request) {
	assetIDStr := chi.URLParam(r, "id")
	var host string

	if assetIDStr != "" {
		assetID, err := strconv.Atoi(assetIDStr)
		if err == nil {
			var asset models.Asset
			err = database.DB.QueryRow(
				`SELECT id, name, location, ip_address FROM assets WHERE id = ?`,
				assetID,
			).Scan(&asset.ID, &asset.Name, &asset.Location, &asset.IPAddress)
			if err == nil {
				if asset.IPAddress != "" {
					host = asset.IPAddress
				} else {
					host = asset.Location
				}
			}
		}
	}

	var req PingRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Host != "" {
			host = req.Host
		}
	}

	if host == "" {
		utils.Error(w, http.StatusBadRequest, "Target host or valid asset ID is required")
		return
	}

	cleanHost := strings.TrimSpace(host)
	start := time.Now()
	ips, err := net.LookupHost(cleanHost)
	latency := time.Since(start).Milliseconds()

	status := "resolved"
	if err != nil {
		status = "unresolved"
		ips = []string{}
	}

	claims := middleware.GetCurrentUser(r)
	if claims != nil {
		_, _ = database.DB.Exec(`
			INSERT INTO audit_logs (user_id, action, entity, entity_id, details, ip_address, created_at)
			VALUES (?, 'DIAGNOSTIC_LOOKUP', 'Asset', 0, ?, ?, ?)
		`, claims.UserID, fmt.Sprintf("DNS lookup executed for host: %s", cleanHost), r.RemoteAddr, time.Now())
	}

	utils.Success(w, http.StatusOK, PingResponse{
		Host:     cleanHost,
		Resolved: ips,
		Status:   status,
		Latency:  latency,
	}, "Host lookup completed")
}

func (h *DiagnosticsHandler) PortCheck(w http.ResponseWriter, r *http.Request) {
	var req PortCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Host == "" || req.Port <= 0 || req.Port > 65535 {
		utils.Error(w, http.StatusBadRequest, "Valid host and port (1-65535) are required")
		return
	}

	timeoutSec := req.Timeout
	if timeoutSec <= 0 || timeoutSec > 10 {
		timeoutSec = 3
	}

	cleanHost := strings.TrimSpace(req.Host)
	targetAddr := net.JoinHostPort(cleanHost, strconv.Itoa(req.Port))

	d := net.Dialer{Timeout: time.Duration(timeoutSec) * time.Second}
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", targetAddr)
	if err != nil {
		utils.Success(w, http.StatusOK, PortCheckResponse{
			Host: cleanHost,
			Port: req.Port,
			Open: false,
		}, "Port check completed")
		return
	}
	defer conn.Close()

	var banner string
	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	if n > 0 {
		banner = strings.TrimSpace(string(buf[:n]))
	}

	utils.Success(w, http.StatusOK, PortCheckResponse{
		Host:   cleanHost,
		Port:   req.Port,
		Open:   true,
		Banner: banner,
	}, "Port check completed")
}
