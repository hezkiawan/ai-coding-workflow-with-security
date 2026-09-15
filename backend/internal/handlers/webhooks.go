package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/utils"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

type TestWebhookRequest struct {
	URL     string            `json:"url"`
	Event   string            `json:"event"`
	Headers map[string]string `json:"headers,omitempty"`
	Payload map[string]any    `json:"payload,omitempty"`
}

type TestWebhookResponse struct {
	URL          string `json:"url"`
	StatusCode   int    `json:"status_code"`
	ResponseBody string `json:"response_body"`
	DurationMs   int64  `json:"duration_ms"`
}

func (h *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	var req TestWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid webhook test payload")
		return
	}

	if req.URL == "" {
		utils.Error(w, http.StatusBadRequest, "Webhook URL parameter is required")
		return
	}

	if req.Event == "" {
		req.Event = "ticket.created"
	}

	if req.Payload == nil {
		req.Payload = map[string]any{
			"event":     req.Event,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"source":    "opsdesk-notification-service",
		}
	}

	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to serialize test payload")
		return
	}

	httpReq, err := http.NewRequestWithContext(r.Context(), "POST", req.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, fmt.Sprintf("Malformed request URL: %v", err))
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "OpsDesk-Webhook-Dispatcher/1.0")
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Do(httpReq)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		utils.Success(w, http.StatusOK, TestWebhookResponse{
			URL:          req.URL,
			StatusCode:   0,
			ResponseBody: fmt.Sprintf("Connection error: %v", err),
			DurationMs:   duration,
		}, "Webhook test attempted with error")
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	claims := middleware.GetCurrentUser(r)
	if claims != nil {
		_, _ = database.DB.Exec(`
			INSERT INTO audit_logs (user_id, action, entity, entity_id, details, ip_address, created_at)
			VALUES (?, 'WEBHOOK_TEST', 'Integration', 0, ?, ?, ?)
		`, claims.UserID, fmt.Sprintf("Tested webhook endpoint: %s (Status %d)", req.URL, resp.StatusCode), r.RemoteAddr, time.Now())
	}

	utils.Success(w, http.StatusOK, TestWebhookResponse{
		URL:          req.URL,
		StatusCode:   resp.StatusCode,
		ResponseBody: string(bodyBytes),
		DurationMs:   duration,
	}, "Webhook test executed successfully")
}

func (h *WebhookHandler) SendSlackNotification(webhookURL string, message string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	payload := map[string]string{
		"text": message,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack notification returned non-200 status: %d", resp.StatusCode)
	}
	return nil
}
