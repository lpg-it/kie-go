package kie

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// apiResponse represents the standard API response structure.
type apiResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    T      `json:"data"`
}

// parseResponseData parses the API response and extracts the data.
func parseResponseData[T any](body []byte) (*T, error) {
	var resp apiResponse[T]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("kie: failed to parse response: %w", err)
	}

	// Check for API-level errors
	if resp.Code != ErrCodeSuccess {
		return nil, newAPIError(0, resp.Code, resp.Message)
	}

	return &resp.Data, nil
}

// validateHTTPResponse checks the HTTP response for errors.
func validateHTTPResponse(resp *http.Response, body []byte) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Try to parse error from body
	var apiResp struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}

	message := string(body)
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Message != "" {
		message = apiResp.Message
	}

	apiErr := newAPIError(resp.StatusCode, apiResp.Code, message)

	// Check for Retry-After header
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		// Parse as seconds (simplified - could also be a date)
		var seconds int
		if _, err := fmt.Sscanf(retryAfter, "%d", &seconds); err == nil {
			apiErr.RetryAfter = time.Duration(seconds) * time.Second
		}
	}

	apiErr.RequestID = resp.Header.Get("X-Request-ID")

	return apiErr
}
