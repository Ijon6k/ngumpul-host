package pagination

import (
	"net/http"
	"strconv"
	"strings"
)

// Params holds normalized pagination parameters.
type Params struct {
	Page   int
	Limit  int
	Offset int
}

// Parse extracts and normalizes "page" and "limit" query parameters from an HTTP request.
func Parse(r *http.Request, defaultLimit, maxLimit int) Params {
	page := 1
	limit := defaultLimit
	if limit <= 0 {
		limit = 10
	}
	if maxLimit <= 0 {
		maxLimit = 100
	}

	pageStr := strings.TrimSpace(r.URL.Query().Get("page"))
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		if l > maxLimit {
			limit = maxLimit
		} else {
			limit = l
		}
	}

	offset := (page - 1) * limit

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}
