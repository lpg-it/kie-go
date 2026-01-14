package kie

import (
	"net/http"
	"time"
)

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithBaseURL sets a custom API base URL.
// Use this for testing or private deployments.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
// The provided client should be configured for production use.
// If nil, a high-performance default client is used.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithTimeout sets the request timeout.
// This applies to individual HTTP requests, not to WaitForTask operations.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithRetry configures retry behavior.
// Use DefaultRetryConfig() as a starting point.
func WithRetry(config *RetryConfig) Option {
	return func(c *Client) {
		if config != nil {
			c.retryConfig = config
		}
	}
}

// WithMaxRetries sets the maximum number of retries.
// Default is 3 retries.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		c.retryConfig.MaxRetries = n
	}
}

// WithDebug enables debug mode.
// Debug mode logs request/response details but NEVER logs sensitive data
// such as API keys or authentication tokens.
func WithDebug(enabled bool) Option {
	return func(c *Client) {
		c.debug = enabled
	}
}

// WithTransport configures the HTTP transport for high-concurrency scenarios.
//
// Use HighConcurrencyTransportConfig() for 10k+ QPS workloads:
//
//	client := kie.NewClient("api-key",
//	    kie.WithTransport(kie.HighConcurrencyTransportConfig()),
//	)
func WithTransport(config *TransportConfig) Option {
	return func(c *Client) {
		if config != nil {
			c.httpClient.Transport = createTransport(config)
		}
	}
}

// RetryConfig configures retry behavior.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts.
	// Set to 0 to disable retries.
	MaxRetries int

	// InitialBackoff is the initial backoff duration.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration.
	MaxBackoff time.Duration

	// Multiplier is the backoff multiplier.
	Multiplier float64

	// Jitter adds randomness to prevent thundering herd.
	Jitter float64
}

// DefaultRetryConfig returns the default retry configuration.
// Suitable for most production use cases.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
		Jitter:         0.1,
	}
}

// NoRetry returns a configuration that disables retries.
func NoRetry() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 0,
	}
}

// calculateBackoff calculates the backoff duration for a given attempt.
func (r *RetryConfig) calculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	backoff := float64(r.InitialBackoff)
	for i := 1; i < attempt; i++ {
		backoff *= r.Multiplier
	}

	if backoff > float64(r.MaxBackoff) {
		backoff = float64(r.MaxBackoff)
	}

	// Add jitter: backoff ± jitter%
	// Using simple deterministic jitter based on attempt number
	jitterFactor := 1.0 + (float64(attempt%2)*2-1)*r.Jitter
	backoff *= jitterFactor

	return time.Duration(backoff)
}
