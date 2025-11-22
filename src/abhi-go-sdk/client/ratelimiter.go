package client

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified configuration
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	if config == nil || !config.Enabled {
		return nil
	}

	return &RateLimiter{
		tokens:     float64(config.BurstSize),
		maxTokens:  float64(config.BurstSize),
		refillRate: config.RequestsPerSecond,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request should be allowed based on rate limiting rules
func (rl *RateLimiter) Allow() bool {
	if rl == nil {
		return true // No rate limiting
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	
	// Refill tokens based on elapsed time
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.maxTokens {
		rl.tokens = rl.maxTokens
	}
	
	rl.lastRefill = now

	// Check if we have tokens available
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}

// Wait blocks until a token is available or the context is canceled
func (rl *RateLimiter) Wait(ctx context.Context) error {
	if rl == nil {
		return nil // No rate limiting
	}

	for {
		if rl.Allow() {
			return nil
		}

		// Calculate wait time until next token is available
		rl.mutex.Lock()
		waitTime := time.Duration((1.0 - rl.tokens) / rl.refillRate * float64(time.Second))
		rl.mutex.Unlock()

		if waitTime > time.Millisecond {
			waitTime = time.Millisecond * 10 // Minimum wait time
		}

		select {
		case <-time.After(waitTime):
			// Continue to next iteration
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// GetAvailableTokens returns the current number of available tokens
func (rl *RateLimiter) GetAvailableTokens() float64 {
	if rl == nil {
		return -1 // No rate limiting
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	
	tokens := rl.tokens + elapsed*rl.refillRate
	if tokens > rl.maxTokens {
		tokens = rl.maxTokens
	}

	return tokens
}

// rateLimitTransport wraps an HTTP transport with rate limiting
type rateLimitTransport struct {
	transport   http.RoundTripper
	rateLimiter *RateLimiter
}

func (rt *rateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter approval
	if err := rt.rateLimiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return rt.transport.RoundTrip(req)
}