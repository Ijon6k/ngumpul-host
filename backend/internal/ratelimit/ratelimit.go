package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"ngumpul-host/backend/internal/response"
)

// IPRateLimiter implements a thread-safe sliding window rate limiter per client IP.
type IPRateLimiter struct {
	mu          sync.RWMutex
	limit       int
	window      time.Duration
	records     map[string][]time.Time
	stopCleanup chan struct{}
}

// New creates an IPRateLimiter and starts a background cleanup worker.
func New(limit int, window time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limit:       limit,
		window:      window,
		records:     make(map[string][]time.Time),
		stopCleanup: make(chan struct{}),
	}

	// Periodically purge stale IP records every 2 * window (minimum 1 minute)
	cleanupInterval := 2 * window
	if cleanupInterval < time.Minute {
		cleanupInterval = time.Minute
	}

	go limiter.cleanupLoop(cleanupInterval)

	return limiter
}

// Allow checks if the given IP is within the rate limit.
func (rl *IPRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	timestamps, exists := rl.records[ip]
	if !exists {
		rl.records[ip] = []time.Time{now}
		return true
	}

	// Filter timestamps within current window
	valid := timestamps[:0]
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.records[ip] = valid
		return false
	}

	rl.records[ip] = append(valid, now)
	return true
}

// Stop terminates the background cleanup routine.
func (rl *IPRateLimiter) Stop() {
	close(rl.stopCleanup)
}

func (rl *IPRateLimiter) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopCleanup:
			return
		case now := <-ticker.C:
			rl.mu.Lock()
			cutoff := now.Add(-rl.window)
			for ip, times := range rl.records {
				validCount := 0
				for _, t := range times {
					if t.After(cutoff) {
						validCount++
					}
				}
				if validCount == 0 {
					delete(rl.records, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// GetClientIP extracts the real client IP from headers or RemoteAddr.
// X-Real-IP is prioritized because reverse proxies like Nginx set it unconditionally
// from the TCP socket remote address ($remote_addr), preventing client header spoofing.
func GetClientIP(r *http.Request) string {
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Fallback to X-Forwarded-For if X-Real-IP is absent
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

// Middleware creates a chi/http middleware for rate limiting using this limiter.
func (rl *IPRateLimiter) Middleware(customMessage ...string) func(http.Handler) http.Handler {
	msg := "Too many requests. Please try again later."
	if len(customMessage) > 0 && customMessage[0] != "" {
		msg = customMessage[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)
			if !rl.Allow(ip) {
				w.Header().Set("Retry-After", "60")
				response.Error(w, http.StatusTooManyRequests, msg)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
