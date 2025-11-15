package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"merge-queue/internal/config"
	"merge-queue/pkg/utils"
)

// RateLimitMiddleware implements basic rate limiting.
type RateLimitMiddleware struct {
	config        *config.Config
	logger        *utils.Logger
	response      *utils.ResponseHelper
	clients       map[string]*clientInfo
	mutex         sync.RWMutex
	cleanupTicker *time.Ticker
}

// clientInfo tracks request information for a client.
type clientInfo struct {
	requests []time.Time
	lastSeen time.Time
}

// NewRateLimitMiddleware creates a new rate limiting middleware.
func NewRateLimitMiddleware(cfg *config.Config, logger *utils.Logger) *RateLimitMiddleware {
	rlm := &RateLimitMiddleware{
		config:   cfg,
		logger:   logger,
		response: utils.NewResponseHelper(),
		clients:  make(map[string]*clientInfo),
	}

	// Start cleanup routine.
	rlm.cleanupTicker = time.NewTicker(5 * time.Minute)
	go rlm.cleanupOldClients()

	return rlm
}

// Handler returns the rate limiting middleware handler.
func (rlm *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rlm.config.Features.RateLimitPerMin <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := rlm.getClientIP(r)

		if rlm.isRateLimited(clientIP) {
			rlm.logger.Warn("Rate limit exceeded for client %s", clientIP)
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rlm.config.Features.RateLimitPerMin))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", "60")
			rlm.response.SendError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		rlm.recordRequest(clientIP)

		// Add rate limit headers.
		remaining := rlm.getRemainingRequests(clientIP)
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rlm.config.Features.RateLimitPerMin))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		next.ServeHTTP(w, r)
	})
}

// Stop stops the cleanup routine.
func (rlm *RateLimitMiddleware) Stop() {
	if rlm.cleanupTicker != nil {
		rlm.cleanupTicker.Stop()
	}
}

// Helper methods.

func (rlm *RateLimitMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header.
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to remote address.
	return r.RemoteAddr
}

func (rlm *RateLimitMiddleware) isRateLimited(clientIP string) bool {
	rlm.mutex.RLock()
	defer rlm.mutex.RUnlock()

	client, exists := rlm.clients[clientIP]
	if !exists {
		return false
	}

	// Count requests in the last minute.
	now := time.Now()
	cutoff := now.Add(-time.Minute)

	count := 0
	for _, reqTime := range client.requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	return count >= rlm.config.Features.RateLimitPerMin
}

func (rlm *RateLimitMiddleware) recordRequest(clientIP string) {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	now := time.Now()

	client, exists := rlm.clients[clientIP]
	if !exists {
		client = &clientInfo{
			requests: make([]time.Time, 0),
			lastSeen: now,
		}
		rlm.clients[clientIP] = client
	}

	// Add current request.
	client.requests = append(client.requests, now)
	client.lastSeen = now

	// Clean up old requests.
	cutoff := now.Add(-time.Minute)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range client.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	client.requests = validRequests
}

func (rlm *RateLimitMiddleware) getRemainingRequests(clientIP string) int {
	rlm.mutex.RLock()
	defer rlm.mutex.RUnlock()

	client, exists := rlm.clients[clientIP]
	if !exists {
		return rlm.config.Features.RateLimitPerMin
	}

	// Count requests in the last minute.
	now := time.Now()
	cutoff := now.Add(-time.Minute)

	count := 0
	for _, reqTime := range client.requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	remaining := rlm.config.Features.RateLimitPerMin - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

func (rlm *RateLimitMiddleware) cleanupOldClients() {
	for range rlm.cleanupTicker.C {
		rlm.mutex.Lock()

		cutoff := time.Now().Add(-10 * time.Minute)
		for clientIP, client := range rlm.clients {
			if client.lastSeen.Before(cutoff) {
				delete(rlm.clients, clientIP)
			}
		}

		rlm.mutex.Unlock()
	}
}
