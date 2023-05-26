package rate_limiter

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHighRPS(t *testing.T) {
	limiter := NewRateLimiter()

	//10 requests made really fast and must be OK
	for i := 0; i < 10; i++ {
		require.True(t, limiter.RegisterCall(), "rate limiter should work with RPS < 10")
	}
	require.False(t, limiter.RegisterCall(), "rate limiter should fail with RPS > 10")
}

func TestHighRPSWithRefill(t *testing.T) {
	limiter := NewRateLimiter()

	//10 requests made really fast and must be OK
	for i := 0; i < 10; i++ {
		require.True(t, limiter.RegisterCall(), "rate limiter should work with RPS < 10")
	}
	//then let rate limiter refill
	time.Sleep(time.Second * 1)
	//10 requests made really fast and must be OK after refill
	for i := 0; i < 10; i++ {
		require.True(t, limiter.RegisterCall(), "rate limiter should work with RPS < 10")
	}
}
