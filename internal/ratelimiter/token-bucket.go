package ratelimiter

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type TokenBucketRateLimiter struct {
	sync.RWMutex
	clients map[string]*Client
	config  Config
}

func NewTokenBucketLimiter(cfg Config) *TokenBucketRateLimiter {
	rl := &TokenBucketRateLimiter{
		clients: make(map[string]*Client),
		config:  cfg,
	}

	go rl.cleanupClients()

	return rl
}

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

func (rl *TokenBucketRateLimiter) getClient(ip string) *Client {
	rl.Lock()
	_, exists := rl.clients[ip]
	defer rl.Unlock()

	if !exists {
		limiter := rate.NewLimiter(rate.Every(rl.config.TimeFrame), rl.config.RequestsPerTimeFrame)
		rl.Lock()

		rl.clients[ip] = &Client{
			Limiter:  limiter,
			LastSeen: time.Now(),
		}

		rl.Unlock()
	}

	v := rl.clients[ip]
	v.LastSeen = time.Now()

	return v
}

func (rl *TokenBucketRateLimiter) Allow(ip string) (bool, float64) {
	client := rl.getClient(ip)

	allowed := client.Limiter.Allow()
	tokens := client.Limiter.Tokens()

	return allowed, tokens
}

// Every minute check the map for client entries that haven't been seen for
// more than 3 minutes and delete the entries.
func (rl *TokenBucketRateLimiter) cleanupClients() {
	for {
		time.Sleep(time.Minute)
		log.Println("Cleaning up clients...", rl.clients)

		rl.Lock()
		for ip, v := range rl.clients {
			if time.Since(v.LastSeen) > 3*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.Unlock()
	}
}
