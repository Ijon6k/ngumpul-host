package activity

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/pagination"
	"ngumpul-host/backend/internal/response"
)

// Handler handles public and administrative activity feeds.
type Handler struct {
	db *pgxpool.Pool
}

// NewHandler creates a new activity handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// ListPublic handles GET /api/activity
func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	p := pagination.Parse(r, 20, 100)
	page, limit, offset := p.Page, p.Limit, p.Offset

	var total int
	_ = h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM activities WHERE visibility = 'PUBLIC'`).Scan(&total)

	rows, err := h.db.Query(r.Context(), `
		SELECT a.id, a.actor_id, a.project_id, a.type, a.metadata, a.visibility, a.created_at,
		       u.display_name, u.username, p.name, p.slug
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		LEFT JOIN projects p ON p.id = a.project_id
		WHERE a.visibility = 'PUBLIC'
		ORDER BY a.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch activities")
		return
	}
	defer rows.Close()

	list := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		var metaJSON []byte
		if err := rows.Scan(
			&a.ID, &a.ActorID, &a.ProjectID, &a.Type, &metaJSON, &a.Visibility, &a.CreatedAt,
			&a.ActorName, &a.ActorUser, &a.ProjName, &a.ProjSlug,
		); err == nil {
			_ = json.Unmarshal(metaJSON, &a.Metadata)
			list = append(list, a)
		}
	}

	totalPages := 1
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"activities":  list,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// ListAdmin handles GET /api/admin/activity
func (h *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT a.id, a.actor_id, a.project_id, a.type, a.metadata, a.visibility, a.created_at,
		       u.display_name, u.username, p.name, p.slug
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		LEFT JOIN projects p ON p.id = a.project_id
		ORDER BY a.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch activities")
		return
	}
	defer rows.Close()

	list := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		var metaJSON []byte
		if err := rows.Scan(
			&a.ID, &a.ActorID, &a.ProjectID, &a.Type, &metaJSON, &a.Visibility, &a.CreatedAt,
			&a.ActorName, &a.ActorUser, &a.ProjName, &a.ProjSlug,
		); err == nil {
			_ = json.Unmarshal(metaJSON, &a.Metadata)
			list = append(list, a)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"activities": list})
}

// ListMyActivity handles GET /api/me/activity
func (h *Handler) ListMyActivity(w http.ResponseWriter, r *http.Request) {
	usr := auth.GetUser(r.Context())
	if usr == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	p := pagination.Parse(r, 15, 100)
	page, limit, offset := p.Page, p.Limit, p.Offset

	var total int
	_ = h.db.QueryRow(r.Context(), `
		SELECT COUNT(*)
		FROM activities a
		LEFT JOIN projects p ON p.id = a.project_id
		WHERE a.actor_id = $1 OR p.owner_id = $1
	`, usr.ID).Scan(&total)

	rows, err := h.db.Query(r.Context(), `
		SELECT a.id, a.actor_id, a.project_id, a.type, a.metadata, a.visibility, a.created_at,
		       u.display_name, u.username, p.name, p.slug
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		LEFT JOIN projects p ON p.id = a.project_id
		WHERE a.actor_id = $1 OR p.owner_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`, usr.ID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch activities")
		return
	}
	defer rows.Close()

	list := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		var metaJSON []byte
		if err := rows.Scan(
			&a.ID, &a.ActorID, &a.ProjectID, &a.Type, &metaJSON, &a.Visibility, &a.CreatedAt,
			&a.ActorName, &a.ActorUser, &a.ProjName, &a.ProjSlug,
		); err == nil {
			_ = json.Unmarshal(metaJSON, &a.Metadata)
			list = append(list, a)
		}
	}

	totalPages := 1
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"activities":  list,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

