package kie

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("kie: circuit breaker is open")

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	// CircuitClosed means the circuit is closed and requests are allowed.
	CircuitClosed CircuitState = iota
	// CircuitOpen means the circuit is open and requests are blocked.
	CircuitOpen
	// CircuitHalfOpen means the circuit is testing if the service has recovered.
	CircuitHalfOpen
)

// CircuitBreakerConfig configures the circuit breaker behavior.
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of consecutive failures before opening.
	// Default: 5
	FailureThreshold int

	// SuccessThreshold is the number of consecutive successes to close from half-open.
	// Default: 2
	SuccessThreshold int

	// Timeout is how long to wait before testing recovery.
	// Default: 30 seconds
	Timeout time.Duration
}

// DefaultCircuitBreakerConfig returns the default circuit breaker configuration.
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
	}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu sync.RWMutex

	config *CircuitBreakerConfig
	state  CircuitState

	failures  int
	successes int
	lastFail  time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}
	return &CircuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// Allow returns true if a request should be allowed.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if timeout has passed
		if time.Since(cb.lastFail) > cb.config.Timeout {
			cb.state = CircuitHalfOpen
			cb.successes = 0
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return true
	}
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.state = CircuitClosed
			cb.failures = 0
			cb.successes = 0
		}
	case CircuitClosed:
		cb.failures = 0
	}
}

// RecordFailure records a failed request.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFail = time.Now()

	switch cb.state {
	case CircuitClosed:
		cb.failures++
		if cb.failures >= cb.config.FailureThreshold {
			cb.state = CircuitOpen
		}
	case CircuitHalfOpen:
		cb.state = CircuitOpen
		cb.failures = cb.config.FailureThreshold
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitClosed
	cb.failures = 0
	cb.successes = 0
}

// CircuitBreakerStats contains statistics about the circuit breaker.
type CircuitBreakerStats struct {
	State    CircuitState
	Failures int
}

// Stats returns the current statistics of the circuit breaker.
func (cb *CircuitBreaker) Stats() CircuitBreakerStats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return CircuitBreakerStats{
		State:    cb.state,
		Failures: cb.failures,
	}
}

// WithCircuitBreaker enables circuit breaker with the given configuration.
func WithCircuitBreaker(config *CircuitBreakerConfig) Option {
	return func(c *Client) {
		c.circuitBreaker = NewCircuitBreaker(config)
	}
}
