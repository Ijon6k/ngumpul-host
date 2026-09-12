package report

import (
	"testing"
	"time"
)

func TestReportValidationAndRateLimit(t *testing.T) {
	// 1. Valid reasons
	for reason := range validReasons {
		if !validReasons[reason] {
			t.Errorf("Expected reason %s to be valid", reason)
		}
	}

	// 2. Invalid reasons
	invalidReasons := []string{"HATE_SPEECH_UNKNOWN", "RANDOM", "", "NOT_REAL"}
	for _, ir := range invalidReasons {
		if validReasons[ir] {
			t.Errorf("Expected reason %s to be invalid", ir)
		}
	}

	// 3. Rate limiting (max 3/min)
	h := NewHandler(nil)
	userID := "user-report-1"

	for i := 0; i < 3; i++ {
		if !h.checkRateLimit(userID) {
			t.Fatalf("Expected report %d to pass rate limit", i+1)
		}
	}

	if h.checkRateLimit(userID) {
		t.Fatalf("Expected 4th report to be blocked by rate limit")
	}

	// Window expiration
	h.userPostMap[userID] = []time.Time{time.Now().Add(-2 * time.Minute)}
	if !h.checkRateLimit(userID) {
		t.Fatalf("Expected rate limit to reset after window expires")
	}
}
