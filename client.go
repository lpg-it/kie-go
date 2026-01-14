package kie

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	// DefaultBaseURL is the default KIE API endpoint
	DefaultBaseURL = "https://api.kie.ai"

	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
)

// Client is the KIE API client.
// It is safe for concurrent use by multiple goroutines.
type Client struct {
	// httpClient is the underlying HTTP client with connection pooling
	httpClient *http.Client

	// baseURL is the API base URL
	baseURL string

	// apiKey is the API key for authentication (stored securely, never logged)
	apiKey string

	// bufferPool provides zero-allocation buffer reuse
	bufferPool *sync.Pool

	// retryConfig controls retry behavior
	retryConfig *RetryConfig

	// circuitBreaker provides stability protection
	circuitBreaker *CircuitBreaker

	// rateLimiter controls request rate (requests per second)
	rateLimiter *rateLimiter

	// metrics is the metrics collector (optional)
	metrics Metrics

	// tracer is the distributed tracer (optional)
	tracer Tracer

	// debug enables debug logging (never logs sensitive data)
	debug bool

	// mu protects concurrent access to client fields
	mu sync.RWMutex
}

// TransportConfig configures the HTTP transport for high-concurrency scenarios.
type TransportConfig struct {
	// MaxIdleConns is the maximum number of idle connections across all hosts.
	// Default: 100. For 10k+ QPS, set to 1000+.
	MaxIdleConns int

	// MaxIdleConnsPerHost is the maximum number of idle connections per host.
	// Default: 100. For 10k+ QPS, set to 1000+.
	MaxIdleConnsPerHost int

	// MaxConnsPerHost is the maximum number of connections per host.
	// Default: 0 (unlimited).
	MaxConnsPerHost int
}

// DefaultTransportConfig returns the default transport configuration.
// Suitable for moderate load (<1000 QPS).
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     0,
	}
}

// HighConcurrencyTransportConfig returns transport configuration for high load.
// Suitable for 10k+ QPS scenarios.
func HighConcurrencyTransportConfig() *TransportConfig {
	return &TransportConfig{
		MaxIdleConns:        10000,
		MaxIdleConnsPerHost: 10000,
		MaxConnsPerHost:     0,
	}
}

// highPerformanceTransport creates an optimized HTTP transport.
// Configuration based on production-grade Go clients (aws-sdk-go-v2, google-cloud-go).
func highPerformanceTransport() *http.Transport {
	return createTransport(DefaultTransportConfig())
}

// createTransport creates an HTTP transport with the given configuration.
func createTransport(config *TransportConfig) *http.Transport {
	return &http.Transport{
		// Connection pooling (configurable)
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,

		// Timeouts
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Performance optimizations
		DisableCompression: false,
		ForceAttemptHTTP2:  true,

		// Dial configuration
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		// TLS configuration
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
}

// NewClient creates a new KIE API client.
//
// The client is configured with:
//   - HTTP/2 support with connection pooling
//   - Automatic retry with exponential backoff
//   - Buffer pooling for zero-allocation request/response handling
//
// Example:
//
//	client := kie.NewClient("your-api-key")
//
//	// With options
//	client := kie.NewClient("your-api-key",
//	    kie.WithTimeout(60*time.Second),
//	    kie.WithRetry(3, kie.DefaultBackoff()),
//	)
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Transport: highPerformanceTransport(),
			Timeout:   DefaultTimeout,
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				// Pre-allocate 4KB buffer, typical request size
				return bytes.NewBuffer(make([]byte, 0, 4096))
			},
		},
		retryConfig: DefaultRetryConfig(),
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// getBuffer retrieves a buffer from the pool.
// The caller MUST call putBuffer when done.
func (c *Client) getBuffer() *bytes.Buffer {
	buf := c.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// putBuffer returns a buffer to the pool.
func (c *Client) putBuffer(buf *bytes.Buffer) {
	// Don't pool buffers that grew too large
	if buf.Cap() <= 1<<20 { // 1MB limit
		c.bufferPool.Put(buf)
	}
}

// doRequest executes an HTTP request with retry support.
// This is the core method that ensures reliability.
func (c *Client) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	// Start tracing span if tracer is configured
	var span Span
	if c.tracer != nil {
		ctx, span = c.tracer.Start(ctx, "kie.doRequest")
		defer span.End()
		span.SetAttribute("http.method", req.Method)
		span.SetAttribute("http.url", req.URL.String())
	}

	// Record metrics
	var startTime time.Time
	if c.metrics != nil {
		startTime = time.Now()
	}

	// Check circuit breaker first
	if c.circuitBreaker != nil && !c.circuitBreaker.Allow() {
		if span != nil {
			span.RecordError(ErrCircuitOpen)
		}
		if c.metrics != nil {
			c.metrics.IncCounter("kie_requests_total", "status", "circuit_open")
		}
		return nil, ErrCircuitOpen
	}

	// Apply rate limiting
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			if span != nil {
				span.RecordError(err)
			}
			return nil, err
		}
	}

	// Set authentication header
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kie-go/"+Version)

	var lastErr error
	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			// Record retry metric
			if c.metrics != nil {
				c.metrics.IncCounter("kie_retries_total")
			}

			// Calculate backoff duration
			backoff := c.retryConfig.calculateBackoff(attempt)

			// Wait with context cancellation support
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Clone request for retry (body needs to be re-readable)
		reqClone := req.Clone(ctx)
		if req.Body != nil && req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to get request body: %w", err)
			}
			reqClone.Body = body
		}

		resp, err := c.httpClient.Do(reqClone)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			// Record failure for circuit breaker
			if c.circuitBreaker != nil {
				c.circuitBreaker.RecordFailure()
			}
			if span != nil {
				span.RecordError(err)
			}
			if isRetryableNetworkError(err) {
				continue
			}
			if c.metrics != nil {
				c.metrics.IncCounter("kie_requests_total", "status", "network_error")
				c.metrics.ObserveHistogram("kie_request_duration_seconds", time.Since(startTime).Seconds(), "status", "error")
			}
			return nil, lastErr
		}

		// Read response body
		body, err := c.readResponseBody(resp)
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			if c.circuitBreaker != nil {
				c.circuitBreaker.RecordFailure()
			}
			if span != nil {
				span.RecordError(err)
			}
			continue
		}

		// Record response status in span
		if span != nil {
			span.SetAttribute("http.status_code", resp.StatusCode)
		}

		// Validate response
		if err := validateHTTPResponse(resp, body); err != nil {
			lastErr = err
			// Record failure for circuit breaker (only for server errors)
			if c.circuitBreaker != nil && isServerError(err) {
				c.circuitBreaker.RecordFailure()
			}
			if span != nil {
				span.RecordError(err)
			}
			if IsRetryable(err) {
				continue
			}
			if c.metrics != nil {
				c.metrics.IncCounter("kie_requests_total", "status", fmt.Sprintf("%d", resp.StatusCode))
				c.metrics.ObserveHistogram("kie_request_duration_seconds", time.Since(startTime).Seconds(), "status", "error")
			}
			return nil, err
		}

		// Success - record for circuit breaker
		if c.circuitBreaker != nil {
			c.circuitBreaker.RecordSuccess()
		}

		// Record success metrics
		if c.metrics != nil {
			c.metrics.IncCounter("kie_requests_total", "status", "success")
			c.metrics.ObserveHistogram("kie_request_duration_seconds", time.Since(startTime).Seconds(), "status", "success")
		}

		return body, nil
	}

	// Record max retries exceeded
	if c.metrics != nil {
		c.metrics.IncCounter("kie_requests_total", "status", "max_retries")
		c.metrics.ObserveHistogram("kie_request_duration_seconds", time.Since(startTime).Seconds(), "status", "error")
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isServerError checks if the error is a server-side error (5xx).
func isServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatus >= 500
	}
	return false
}

// readResponseBody reads the response body efficiently using buffer pool.
func (c *Client) readResponseBody(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, nil
	}

	// Use pooled buffer for reading
	buf := c.getBuffer()
	defer c.putBuffer(buf)

	// Limit response size to prevent memory exhaustion (10MB)
	limited := io.LimitReader(resp.Body, 10<<20)
	if _, err := buf.ReadFrom(limited); err != nil {
		return nil, err
	}

	// Make a copy since buffer will be reused
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// isRetryableNetworkError checks if a network error is retryable.
func isRetryableNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check for timeout
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	// Check for connection reset, refused, etc.
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Temporary()
	}

	return false
}

// Close releases resources held by the client.
// After Close, the client should not be used.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}
