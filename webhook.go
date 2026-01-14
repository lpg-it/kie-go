package kie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// WebhookHandler handles KIE callback/webhook requests.
//
// Use this to receive task completion notifications instead of polling.
type WebhookHandler struct {
	// secretKey is used for signature verification (optional)
	secretKey string
}

// NewWebhookHandler creates a new webhook handler.
//
// If secretKey is provided, all incoming requests will be verified
// using HMAC-SHA256 signature.
func NewWebhookHandler(secretKey string) *WebhookHandler {
	return &WebhookHandler{
		secretKey: secretKey,
	}
}

// WebhookPayload represents the structure of a KIE webhook callback.
type WebhookPayload struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID       string `json:"taskId"`
		Model        string `json:"model"`
		State        string `json:"state"`
		Param        string `json:"param"`
		ResultJSON   string `json:"resultJson"`
		FailCode     string `json:"failCode"`
		FailMsg      string `json:"failMsg"`
		CostTime     int64  `json:"costTime"`
		CompleteTime int64  `json:"completeTime"`
		CreateTime   int64  `json:"createTime"`
	} `json:"data"`
}

// ParseRequest parses a webhook HTTP request into TaskInfo.
//
// This method reads the request body, validates the signature (if configured),
// and parses the payload into a TaskInfo struct.
//
// Example:
//
//	handler := kie.NewWebhookHandler("your-secret-key")
//
//	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
//	    info, err := handler.ParseRequest(r)
//	    if err != nil {
//	        http.Error(w, err.Error(), http.StatusBadRequest)
//	        return
//	    }
//	    log.Printf("Task %s completed: %s", info.TaskID, info.State)
//	    w.WriteHeader(http.StatusOK)
//	})
func (h *WebhookHandler) ParseRequest(r *http.Request) (*TaskInfo, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("webhook: expected POST, got %s", r.Method)
	}

	// Read body (limit to 10MB to prevent DoS)
	body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		return nil, fmt.Errorf("webhook: failed to read body: %w", err)
	}
	defer r.Body.Close()

	// Verify signature if configured
	if h.secretKey != "" {
		if err := h.verifySignature(r, body); err != nil {
			return nil, err
		}
	}

	// Parse payload
	return h.ParsePayload(body)
}

// ParsePayload parses a raw webhook payload into TaskInfo.
//
// Use this when you have the raw body bytes and have already
// handled signature verification.
func (h *WebhookHandler) ParsePayload(body []byte) (*TaskInfo, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("webhook: failed to parse payload: %w", err)
	}

	// Validate response code
	if payload.Code != 200 {
		return nil, fmt.Errorf("webhook: API error (code=%d): %s", payload.Code, payload.Msg)
	}

	// Convert to TaskInfo
	info := &TaskInfo{
		TaskID:       payload.Data.TaskID,
		Model:        payload.Data.Model,
		State:        TaskState(payload.Data.State),
		Param:        payload.Data.Param,
		ResultJSON:   payload.Data.ResultJSON,
		FailCode:     payload.Data.FailCode,
		FailMsg:      payload.Data.FailMsg,
		CostTime:     payload.Data.CostTime,
		CompleteTime: payload.Data.CompleteTime,
		CreateTime:   payload.Data.CreateTime,
	}

	return info, nil
}

// verifySignature verifies the HMAC-SHA256 signature of the request.
//
// Expected header: X-KIE-Signature: sha256=<hex-encoded-signature>
func (h *WebhookHandler) verifySignature(r *http.Request, body []byte) error {
	signature := r.Header.Get("X-KIE-Signature")
	if signature == "" {
		return fmt.Errorf("webhook: missing signature header")
	}

	// Parse signature header (format: "sha256=<hex>")
	parts := strings.SplitN(signature, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return fmt.Errorf("webhook: invalid signature format")
	}

	expectedSig, err := hex.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("webhook: invalid signature encoding: %w", err)
	}

	// Compute expected signature
	mac := hmac.New(sha256.New, []byte(h.secretKey))
	mac.Write(body)
	computedSig := mac.Sum(nil)

	// Constant-time comparison to prevent timing attacks
	if !hmac.Equal(expectedSig, computedSig) {
		return fmt.Errorf("webhook: signature verification failed")
	}

	return nil
}

// Handler returns an http.HandlerFunc that processes webhook callbacks.
//
// This is a convenience method for creating an HTTP handler that:
// 1. Parses and validates the webhook request
// 2. Calls your callback function with the TaskInfo
// 3. Responds with appropriate HTTP status codes
//
// Example:
//
//	handler := kie.NewWebhookHandler("secret")
//
//	http.HandleFunc("/callback", handler.Handler(func(info *kie.TaskInfo) {
//	    log.Printf("Task completed: %s", info.TaskID)
//	    // Process the result...
//	}))
func (h *WebhookHandler) Handler(callback func(*TaskInfo)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := h.ParseRequest(r)
		if err != nil {
			// Log error but don't expose details
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Call user callback
		callback(info)

		// Acknowledge receipt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

// HandlerWithError returns an http.HandlerFunc with error handling support.
//
// Unlike Handler, this version allows your callback to return an error,
// which will result in a 500 response (telling KIE to retry).
func (h *WebhookHandler) HandlerWithError(callback func(*TaskInfo) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := h.ParseRequest(r)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Call user callback
		if err := callback(info); err != nil {
			// Return 500 to trigger retry
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Acknowledge receipt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}
