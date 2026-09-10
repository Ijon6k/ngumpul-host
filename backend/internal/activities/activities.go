package activities

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Activity struct {
	ID         string         `json:"id"`
	ActorID    *string        `json:"actor_id"`
	ProjectID  *string        `json:"project_id"`
	Type       string         `json:"type"`
	Metadata   map[string]any `json:"metadata"`
	Visibility string         `json:"visibility"`
	CreatedAt  time.Time      `json:"created_at"`
	ActorName  *string        `json:"actor_name"`
	ActorUser  *string        `json:"actor_username"`
	ProjName   *string        `json:"project_name"`
	ProjSlug   *string        `json:"project_slug"`
}

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT a.id, a.actor_id, a.project_id, a.type, a.metadata, a.visibility, a.created_at,
		       u.display_name, u.username, p.name, p.slug
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		LEFT JOIN projects p ON p.id = a.project_id
		WHERE a.visibility = 'PUBLIC'
		ORDER BY a.created_at DESC
		LIMIT 50
	`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch activities"})
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

	writeJSON(w, http.StatusOK, map[string]any{"activities": list})
}

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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch activities"})
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

	writeJSON(w, http.StatusOK, map[string]any{"activities": list})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
