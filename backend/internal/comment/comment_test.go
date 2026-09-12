package comment

import (
	"testing"
	"time"
)

func TestCommentRateLimit(t *testing.T) {
	h := NewHandler(nil)
	userID := "user-123"

	// Posting 5 times in quick succession should be allowed
	for i := 0; i < 5; i++ {
		if !h.checkRateLimit(userID) {
			t.Fatalf("Expected comment %d to be allowed within rate limit", i+1)
		}
	}

	// 6th attempt should be blocked
	if h.checkRateLimit(userID) {
		t.Fatalf("Expected 6th comment to be rejected by rate limit")
	}

	// Other user should still be allowed
	if !h.checkRateLimit("other-user") {
		t.Fatalf("Expected different user to be allowed")
	}

	// Expired timestamps should be pruned
	h.userPostMap[userID] = []time.Time{
		time.Now().Add(-2 * time.Minute),
		time.Now().Add(-70 * time.Second),
	}
	if !h.checkRateLimit(userID) {
		t.Fatalf("Expected rate limit to allow post after window expiration")
	}
}
