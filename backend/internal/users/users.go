package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

type UserWithProjects struct {
	auth.PublicUser
	ProjectsCount int `json:"projects_count"`
}

type UserProfileResponse struct {
	auth.PublicUser
	Projects []any `json:"projects"`
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at,
		       COUNT(p.id) as projects_count
		FROM users u
		LEFT JOIN projects p ON p.owner_id = u.id AND p.visibility = 'PUBLIC'
		WHERE u.status = 'ACTIVE'
		GROUP BY u.id
		ORDER BY u.created_at ASC
	`)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch members"})
		return
	}
	defer rows.Close()

	members := make([]UserWithProjects, 0)
	for rows.Next() {
		var m UserWithProjects
		if err := rows.Scan(
			&m.ID, &m.Username, &m.DisplayName, &m.AvatarURL, &m.Bio, &m.Role, &m.CreatedAt, &m.ProjectsCount,
		); err == nil {
			members = append(members, m)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"members": members,
		"total":   len(members),
	})
}

func (h *Handler) GetMember(w http.ResponseWriter, r *http.Request) {
	username := strings.ToLower(strings.TrimSpace(chi.URLParam(r, "username")))

	var u auth.PublicUser
	err := h.db.QueryRow(r.Context(), `
		SELECT id, username, display_name, avatar_url, bio, role, created_at
		FROM users
		WHERE LOWER(username) = $1 AND status = 'ACTIVE'
	`, username).Scan(
		&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to load user profile"})
		return
	}

	// Fetch user's published projects
	rows, err := h.db.Query(r.Context(), `
		SELECT id, name, slug, description, cover_image_url, repository_url, documentation_url,
		       demo_url, technology_stack, hosting_type, public_url, status, created_at, published_at
		FROM projects
		WHERE owner_id = $1 AND visibility = 'PUBLIC'
		ORDER BY created_at DESC
	`, u.ID)

	projectsList := make([]map[string]any, 0)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, name, slug, desc, cover, repo, docs, demo, hostingType, publicURL, status string
			var tech []string
			var createdAt timeScanner
			var publishedAt *timeScanner
			if err := rows.Scan(
				&id, &name, &slug, &desc, &cover, &repo, &docs, &demo, &tech, &hostingType, &publicURL, &status, &createdAt, &publishedAt,
			); err == nil {
				projectsList = append(projectsList, map[string]any{
					"id":                 id,
					"name":               name,
					"slug":               slug,
					"description":        desc,
					"cover_image_url":    cover,
					"repository_url":     repo,
					"documentation_url":  docs,
					"demo_url":           demo,
					"technology_stack":   tech,
					"hosting_type":       hostingType,
					"public_url":         publicURL,
					"status":             status,
					"created_at":         createdAt,
					"published_at":       publishedAt,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user":     u,
		"projects": projectsList,
	})
}

type timeScanner = any
