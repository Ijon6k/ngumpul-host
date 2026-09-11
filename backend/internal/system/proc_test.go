package system

import (
	"testing"
	"time"
)

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		sec      int64
		expected string
	}{
		{0, "Just started"},
		{30, "30 mins"}, // or min calculation
		{3600, "1 hours, 0 mins"},
		{90000, "1 days, 1 hours"},
	}

	for _, tc := range tests {
		got := FormatUptime(tc.sec)
		if got == "" {
			t.Errorf("FormatUptime(%d) returned empty string", tc.sec)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	if got := FormatDuration(30); got != "30s" {
		t.Errorf("expected '30s', got %q", got)
	}
	if got := FormatDuration(90); got != "1m 30s" {
		t.Errorf("expected '1m 30s', got %q", got)
	}
	if got := FormatDuration(3660); got != "1h 1m" {
		t.Errorf("expected '1h 1m', got %q", got)
	}
}

func TestFormatRelativeTime(t *testing.T) {
	if got := FormatRelativeTime(time.Time{}); got != "Never" {
		t.Errorf("expected 'Never', got %q", got)
	}
	if got := FormatRelativeTime(time.Now()); got != "Just now" {
		t.Errorf("expected 'Just now', got %q", got)
	}
}
