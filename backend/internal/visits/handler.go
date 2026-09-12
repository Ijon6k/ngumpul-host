package visits

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/response"
)

// Handler handles outbound redirects and privacy-conscious visit statistics.
type Handler struct {
	db        *pgxpool.Pool
	viewCache sync.Map // key: "projectID:clientKey", value: time.Time
}

// NewHandler creates a visits handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		db: db,
	}
}

// Redirect handles GET /go/{slug}
// Records an outbound click ("Visits from Ngumpul") and redirects HTTP 302 to project's public_url.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Project slug is required", http.StatusBadRequest)
		return
	}

	var projectID, publicURL, status, visibility string
	err := h.db.QueryRow(r.Context(), `
		SELECT id, public_url, status, visibility
		FROM projects
		WHERE slug = $1
	`, slug).Scan(&projectID, &publicURL, &status, &visibility)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if publicURL == "" {
		http.Error(w, "Project has no configured public URL", http.StatusNotFound)
		return
	}

	// Record outbound visit asynchronously
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = h.db.Exec(ctx, `
			INSERT INTO project_visits (project_id, visited_at, date_bucket)
			VALUES ($1, NOW(), CURRENT_DATE)
		`, projectID)
	}()

	// 302 Found redirect
	http.Redirect(w, r, publicURL, http.StatusFound)
}

// RecordPageView tracks a lightweight page view for /projects/:slug with 30-min deduplication.
func (h *Handler) RecordPageView(ctx context.Context, projectID string, clientKey string) {
	if projectID == "" {
		return
	}

	cacheKey := fmt.Sprintf("%s:%s", projectID, clientKey)
	now := time.Now()

	if lastSeen, ok := h.viewCache.Load(cacheKey); ok {
		if t, ok := lastSeen.(time.Time); ok && now.Sub(t) < 30*time.Minute {
			return // Deduplicated: view already counted within last 30 minutes
		}
	}

	h.viewCache.Store(cacheKey, now)

	// Clean cache occasionally or let it be bounded
	_, _ = h.db.Exec(ctx, `
		INSERT INTO project_page_views (project_id, date, views_count)
		VALUES ($1, CURRENT_DATE, 1)
		ON CONFLICT (project_id, date) DO UPDATE
		SET views_count = project_page_views.views_count + 1
	`, projectID)
}

// GetProjectVisits handles GET /api/me/projects/{id}/visits
// Authorization: Must be owner of the project or an administrator.
func (h *Handler) GetProjectVisits(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUser(r.Context())
	if caller == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	// Verify project ownership
	var ownerID string
	err := h.db.QueryRow(r.Context(), "SELECT owner_id FROM projects WHERE id = $1", projectID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if caller.Role != "ADMIN" && caller.ID != ownerID {
		response.Error(w, http.StatusForbidden, "Forbidden: You do not own this project")
		return
	}

	// 30 days series map
	now := time.Now()
	startDate := now.AddDate(0, 0, -29)
	dailyMap := make(map[string]*DailyMetric, 30)

	for i := 0; i < 30; i++ {
		d := startDate.AddDate(0, 0, i).Format("2006-01-02")
		dailyMap[d] = &DailyMetric{
			Date:           d,
			Views:          0,
			OutboundVisits: 0,
		}
	}

	// 1. Query daily views for last 30 days
	viewRows, err := h.db.Query(r.Context(), `
		SELECT to_char(date, 'YYYY-MM-DD'), views_count
		FROM project_page_views
		WHERE project_id = $1 AND date >= CURRENT_DATE - INTERVAL '29 days'
	`, projectID)
	if err == nil {
		defer viewRows.Close()
		for viewRows.Next() {
			var dateStr string
			var views int
			if err := viewRows.Scan(&dateStr, &views); err == nil {
				if m, ok := dailyMap[dateStr]; ok {
					m.Views = views
				}
			}
		}
	}

	// 2. Query outbound visits for last 30 days
	visitRows, err := h.db.Query(r.Context(), `
		SELECT to_char(date_bucket, 'YYYY-MM-DD'), COUNT(*)
		FROM project_visits
		WHERE project_id = $1 AND date_bucket >= CURRENT_DATE - INTERVAL '29 days'
		GROUP BY date_bucket
	`, projectID)
	if err == nil {
		defer visitRows.Close()
		for visitRows.Next() {
			var dateStr string
			var count int
			if err := visitRows.Scan(&dateStr, &count); err == nil {
				if m, ok := dailyMap[dateStr]; ok {
					m.OutboundVisits = count
				}
			}
		}
	}

	totalViews := 0
	totalOutbound := 0
	breakdown := make([]DailyMetric, 0, 30)

	for i := 0; i < 30; i++ {
		d := startDate.AddDate(0, 0, i).Format("2006-01-02")
		m := dailyMap[d]
		totalViews += m.Views
		totalOutbound += m.OutboundVisits
		breakdown = append(breakdown, *m)
	}

	response.JSON(w, http.StatusOK, ProjectVisitsResponse{
		TotalViews30d:    totalViews,
		TotalOutbound30d: totalOutbound,
		DailyBreakdown:   breakdown,
	})
}
