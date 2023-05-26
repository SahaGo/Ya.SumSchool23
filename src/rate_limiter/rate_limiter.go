package rate_limiter

import (
	"math"
	"sync"
	"time"
)

type RateLimiter struct {
	rate                int64 //RPS
	maxTokens           int64
	currentTokens       float64
	lastRefillTimestamp time.Time
	mutex               sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		rate:                10, //let it be fixed
		maxTokens:           10,
		currentTokens:       10,
		lastRefillTimestamp: time.Now(),
	}
}

func (r *RateLimiter) refill() {
	now := time.Now()
	end := time.Since(r.lastRefillTimestamp)
	tokensToBeAdded := end.Seconds() * float64(r.rate)
	r.currentTokens = math.Min(r.currentTokens+tokensToBeAdded, float64(r.maxTokens))
	r.lastRefillTimestamp = now
}

func (r *RateLimiter) RegisterCall() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.refill()
	if r.currentTokens >= 1 {
		r.currentTokens = r.currentTokens - 1
		return true
	}
	return false
}
