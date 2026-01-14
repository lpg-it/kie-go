package kie

import (
	"context"
	"sync"
	"time"
)

// rateLimiter implements a simple token bucket rate limiter.
type rateLimiter struct {
	mu sync.Mutex

	rate     float64   // tokens per second
	burst    int       // maximum tokens
	tokens   float64   // current tokens
	lastTime time.Time // last token update time
}

// newRateLimiter creates a new rate limiter.
func newRateLimiter(rps int) *rateLimiter {
	return &rateLimiter{
		rate:     float64(rps),
		burst:    rps,
		tokens:   float64(rps),
		lastTime: time.Now(),
	}
}

// Wait blocks until a token is available or the context is canceled.
func (l *rateLimiter) Wait(ctx context.Context) error {
	for {
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.lastTime).Seconds()
		l.tokens += elapsed * l.rate
		if l.tokens > float64(l.burst) {
			l.tokens = float64(l.burst)
		}
		l.lastTime = now

		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}

		// Calculate wait time
		waitTime := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Continue loop to try again
		}
	}
}

// TryAcquire attempts to acquire a token without blocking.
// Returns true if successful.
func (l *rateLimiter) TryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTime).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastTime = now

	if l.tokens >= 1 {
		l.tokens--
		return true
	}

	return false
}

// WithRateLimit enables rate limiting with the given requests per second.
func WithRateLimit(rps int) Option {
	return func(c *Client) {
		if rps > 0 {
			c.rateLimiter = newRateLimiter(rps)
		}
	}
}
