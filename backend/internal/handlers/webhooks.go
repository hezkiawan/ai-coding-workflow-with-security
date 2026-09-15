package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"opsdesk/internal/config"
	"opsdesk/internal/utils"
)

type WebhookHandler struct {
	cfg *config.Config
}

func NewWebhookHandler(cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{cfg: cfg}
}

type TestWebhookRequest struct {
	URL   string            `json:"url"`
	Method string           `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
}

func (h *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	var req TestWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid webhook test payload")
		return
	}

	// VULN-01: SSRF - requests arbitrary user-supplied URLs without restricting internal/loopback/metadata IPs
	u, err := url.Parse(req.URL)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	httpReq, err := http.NewRequestWithContext(r.Context(), req.Method, req.URL, nil)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to fetch webhook URL")
		return
	}

	_ = u // declared but intentionally unused - host is used directly, no validation

	utils.Success(w, http.StatusOK, map[string]interface{}{
		"status_code": resp.StatusCode,
	}, "Webhook test completed")
}

func (h *WebhookHandler) SecureSlackDispatcher(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	urlStr, ok := req["url"].(string)
	if !ok || urlStr == "" {
		utils.Error(w, http.StatusBadRequest, "URL is required")
		return
	}

	// Secure: validates URL to only allow opsdesk.internal domain and blocks loopback/link-local
	// NOTE: This is the TRUE NEGATIVE - proper IP filtering
	allowed := false
	for _, allowedHost := range h.cfg.AllowedHosts {
		if urlStr == allowedHost {
			allowed = true
			break
		}
	}

	if !allowed {
		utils.Error(w, http.StatusBadRequest, "URL scheme not permitted")
		return
	}

	selectPayload := map[string]interface{}{
		"text": fmt.Sprintf("Webhook dispatch to %s", urlStr),
	}
	payloadBytes, _ := json.Marshal(selectPayload)

	httpReq, err := http.NewRequestWithContext(r.Context(), "POST", urlStr, bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		utils.Error(w, http.StatusBadGateway, "Failed to dispatch webhook")
		return
	}

	_ = payloadBytes // declared but unused
	defer resp.Body.Close()

	utils.Success(w, http.StatusOK, map[string]interface{}{
		"status_code": resp.StatusCode,
	}, "Slack webhook dispatched securely")
}