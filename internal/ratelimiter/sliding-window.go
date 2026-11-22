package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowRateLimiter struct {
	sync.RWMutex
	clients map[string][]time.Time
	limit   int
	window  int64
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  int64(window.Seconds()),
	}
}

func (rl *SlidingWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	rl.clients[ip] = append(rl.clients[ip], time.Now())

	// Remove timestamps that are outside the window
	for len(rl.clients[ip]) > 0 && time.Since(rl.clients[ip][0]).Seconds() > float64(rl.window) {
		rl.clients[ip] = rl.clients[ip][1:]
	}

	// Check if the number of requests in the current window exceeds the limit
	if len(rl.clients[ip]) > rl.limit {
		// Calculate the time until the next request is allowed
		retryAfter := time.Duration(rl.window)*time.Second - time.Since(rl.clients[ip][0])
		return false, retryAfter
	}

	return true, 0
}
