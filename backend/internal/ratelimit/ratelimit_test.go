package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIPRateLimiter_Allow(t *testing.T) {
	limiter := New(3, 50*time.Millisecond)
	defer limiter.Stop()

	ip := "192.168.1.100"

	// First 3 requests should be allowed
	for i := 1; i <= 3; i++ {
		if !limiter.Allow(ip) {
			t.Fatalf("Request %d should have been allowed", i)
		}
	}

	// 4th request should be rejected
	if limiter.Allow(ip) {
		t.Fatalf("Request 4 should have been blocked by rate limit")
	}

	// Different IP should still be allowed
	otherIP := "192.168.1.101"
	if !limiter.Allow(otherIP) {
		t.Fatalf("Different IP should be allowed")
	}

	// Wait for window to expire
	time.Sleep(60 * time.Millisecond)

	// Should be allowed again
	if !limiter.Allow(ip) {
		t.Fatalf("Request after window expiration should be allowed")
	}
}

func TestIPRateLimiter_Middleware(t *testing.T) {
	limiter := New(2, 50*time.Millisecond)
	defer limiter.Stop()

	handler := limiter.Middleware("Rate limit exceeded")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	// Request 1
	req1 := httptest.NewRequest(http.MethodPost, "/test", nil)
	req1.RemoteAddr = "10.0.0.1:12345"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec1.Code)
	}

	// Request 2
	req2 := httptest.NewRequest(http.MethodPost, "/test", nil)
	req2.RemoteAddr = "10.0.0.1:12345"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec2.Code)
	}

	// Request 3 (blocked)
	req3 := httptest.NewRequest(http.MethodPost, "/test", nil)
	req3.RemoteAddr = "10.0.0.1:12345"
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected status 429, got %d", rec3.Code)
	}
}
