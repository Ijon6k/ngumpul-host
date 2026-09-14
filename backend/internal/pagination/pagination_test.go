package pagination

import (
	"net/http/httptest"
	"testing"
)

func TestParse(t *testing.T) {
	// 1. Defaults
	req1 := httptest.NewRequest("GET", "/test", nil)
	p1 := Parse(req1, 12, 100)
	if p1.Page != 1 || p1.Limit != 12 || p1.Offset != 0 {
		t.Errorf("expected default page=1 limit=12 offset=0, got %+v", p1)
	}

	// 2. Custom valid values
	req2 := httptest.NewRequest("GET", "/test?page=3&limit=25", nil)
	p2 := Parse(req2, 10, 100)
	if p2.Page != 3 || p2.Limit != 25 || p2.Offset != 50 {
		t.Errorf("expected page=3 limit=25 offset=50, got %+v", p2)
	}

	// 3. Clamped limit
	req3 := httptest.NewRequest("GET", "/test?page=2&limit=500", nil)
	p3 := Parse(req3, 10, 100)
	if p3.Page != 2 || p3.Limit != 100 || p3.Offset != 100 {
		t.Errorf("expected page=2 limit=100 offset=100, got %+v", p3)
	}

	// 4. Invalid numbers fall back to defaults
	req4 := httptest.NewRequest("GET", "/test?page=-5&limit=invalid", nil)
	p4 := Parse(req4, 15, 100)
	if p4.Page != 1 || p4.Limit != 15 || p4.Offset != 0 {
		t.Errorf("expected fallback to default, got %+v", p4)
	}
}
