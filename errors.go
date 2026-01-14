package kie

import (
	"errors"
	"fmt"
	"time"
)

// Error codes from KIE API
const (
	ErrCodeSuccess          = 200
	ErrCodeBadRequest       = 400
	ErrCodeUnauthorized     = 401
	ErrCodeInsufficientFund = 402
	ErrCodeNotFound         = 404
	ErrCodeValidation       = 422
	ErrCodeRateLimit        = 429
	ErrCodeInternal         = 500
)

// Common errors
var (
	// ErrNilInput is returned when a required input is nil
	ErrNilInput = errors.New("kie: input cannot be nil")

	// ErrEmptyTaskID is returned when task ID is empty
	ErrEmptyTaskID = errors.New("kie: task ID cannot be empty")

	// ErrTaskFailed is returned when a task fails
	ErrTaskFailed = errors.New("kie: task failed")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("kie: operation timed out")

	// ErrMaxRetriesExceeded is returned when max retries are exceeded
	ErrMaxRetriesExceeded = errors.New("kie: max retries exceeded")
)

// APIError represents an error returned by the KIE API.
type APIError struct {
	// HTTPStatus is the HTTP status code
	HTTPStatus int

	// Code is the API error code
	Code int

	// Message is the error message from the API
	Message string

	// RequestID is the request ID for debugging (if available)
	RequestID string

	// RetryAfter is the suggested wait time before retrying (for 429 errors)
	// Zero value means no Retry-After header was present
	RetryAfter time.Duration

	// retryable indicates if this error can be retried
	retryable bool
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("kie: API error (code=%d, http=%d, request_id=%s): %s",
			e.Code, e.HTTPStatus, e.RequestID, e.Message)
	}
	return fmt.Sprintf("kie: API error (code=%d, http=%d): %s",
		e.Code, e.HTTPStatus, e.Message)
}

// IsRetryable returns true if this error can be retried.
func (e *APIError) IsRetryable() bool {
	return e.retryable
}

// newAPIError creates a new API error from code and message.
func newAPIError(httpStatus, code int, message string) *APIError {
	return &APIError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		retryable:  isRetryableCode(code) || isRetryableHTTPStatus(httpStatus),
	}
}

// isRetryableCode checks if an API error code is retryable.
func isRetryableCode(code int) bool {
	switch code {
	case ErrCodeRateLimit, ErrCodeInternal:
		return true
	default:
		return false
	}
}

// isRetryableHTTPStatus checks if an HTTP status is retryable.
func isRetryableHTTPStatus(status int) bool {
	switch status {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// IsRetryable returns true if the error can be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}

	return false
}

// IsAuthError returns true if the error is an authentication error.
func IsAuthError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == ErrCodeUnauthorized || apiErr.HTTPStatus == 401
	}
	return false
}

// IsRateLimitError returns true if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == ErrCodeRateLimit || apiErr.HTTPStatus == 429
	}
	return false
}

// IsInsufficientFundError returns true if the error is an insufficient funds error.
func IsInsufficientFundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == ErrCodeInsufficientFund || apiErr.HTTPStatus == 402
	}
	return false
}

// IsValidationError returns true if the error is a validation error.
func IsValidationError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == ErrCodeValidation || apiErr.HTTPStatus == 422
	}
	return false
}

// TaskFailedError represents a task that failed on the server side.
type TaskFailedError struct {
	TaskID   string
	FailCode string
	FailMsg  string
}

// Error implements the error interface.
func (e *TaskFailedError) Error() string {
	return fmt.Sprintf("kie: task %s failed (code=%s): %s", e.TaskID, e.FailCode, e.FailMsg)
}

// Is checks if target error is ErrTaskFailed.
func (e *TaskFailedError) Is(target error) bool {
	return target == ErrTaskFailed
}

// ResultValidationError indicates an anomaly in the response data.
// This is returned when the task state is "success" but the result data is missing or invalid.
// The TaskInfo is still returned, allowing the user to inspect the raw response.
type ResultValidationError struct {
	TaskID  string
	Message string
}

// Error implements the error interface.
func (e *ResultValidationError) Error() string {
	return fmt.Sprintf("kie: result validation error for task %s: %s", e.TaskID, e.Message)
}

// ErrResultValidation is a sentinel error for result validation failures.
var ErrResultValidation = errors.New("kie: result validation failed")

// Is checks if target error is ErrResultValidation.
func (e *ResultValidationError) Is(target error) bool {
	return target == ErrResultValidation
}
