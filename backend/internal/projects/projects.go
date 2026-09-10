package projects

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
)

type Project struct {
	ID                string     `json:"id"`
	OwnerID           string     `json:"owner_id"`
	Name              string     `json:"name"`
	Slug              string     `json:"slug"`
	Description       string     `json:"description"`
	CoverImageURL     string     `json:"cover_image_url"`
	RepositoryURL     string     `json:"repository_url"`
	DocumentationURL  string     `json:"documentation_url"`
	DemoURL           string     `json:"demo_url"`
	TechnologyStack   []string   `json:"technology_stack"`
	HostingType       string     `json:"hosting_type"`
	PublicURL         string     `json:"public_url"`
	Status            string     `json:"status"`
	Visibility        string     `json:"visibility"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	PublishedAt       *time.Time `json:"published_at"`
	Owner             *auth.PublicUser `json:"owner,omitempty"`
}

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

var slugRegex = regexp.MustCompile(`[^a-z0-9\-]+`)

func GenerateSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRegex.ReplaceAllString(s, "")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "project"
	}
	return s
}

func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	techFilter := r.URL.Query().Get("tech")

	query := `
		SELECT p.id, p.owner_id, p.name, p.slug, p.description, p.cover_image_url,
		       p.repository_url, p.documentation_url, p.demo_url, p.technology_stack,
		       p.hosting_type, p.public_url, p.status, p.visibility, p.created_at, p.updated_at, p.published_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		WHERE p.visibility = 'PUBLIC'
	`
	args := []any{}

	if statusFilter != "" {
		args = append(args, strings.ToUpper(statusFilter))
		query += ` AND p.status = $` + string(rune('0'+len(args)))
	}
	if techFilter != "" {
		args = append(args, techFilter)
		query += ` AND $` + string(rune('0'+len(args))) + ` = ANY(p.technology_stack)`
	}

	query += ` ORDER BY p.published_at DESC NULLS LAST, p.created_at DESC`

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list projects"})
		return
	}
	defer rows.Close()

	list := make([]Project, 0)
	for rows.Next() {
		var p Project
		var u auth.PublicUser
		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Name, &p.Slug, &p.Description, &p.CoverImageURL,
			&p.RepositoryURL, &p.DocumentationURL, &p.DemoURL, &p.TechnologyStack,
			&p.HostingType, &p.PublicURL, &p.Status, &p.Visibility, &p.CreatedAt, &p.UpdatedAt, &p.PublishedAt,
			&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
		); err == nil {
			p.Owner = &u
			list = append(list, p)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"projects": list,
		"total":    len(list),
	})
}

func (h *Handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := strings.ToLower(strings.TrimSpace(chi.URLParam(r, "slug")))

	currentUser := auth.GetUser(r.Context())

	var p Project
	var u auth.PublicUser
	err := h.db.QueryRow(r.Context(), `
		SELECT p.id, p.owner_id, p.name, p.slug, p.description, p.cover_image_url,
		       p.repository_url, p.documentation_url, p.demo_url, p.technology_stack,
		       p.hosting_type, p.public_url, p.status, p.visibility, p.created_at, p.updated_at, p.published_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		WHERE LOWER(p.slug) = $1
	`, slug).Scan(
		&p.ID, &p.OwnerID, &p.Name, &p.Slug, &p.Description, &p.CoverImageURL,
		&p.RepositoryURL, &p.DocumentationURL, &p.DemoURL, &p.TechnologyStack,
		&p.HostingType, &p.PublicURL, &p.Status, &p.Visibility, &p.CreatedAt, &p.UpdatedAt, &p.PublishedAt,
		&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	// Object-level authorization:
	// If project is UNPUBLISHED or ARCHIVED, only owner or ADMIN can view it
	if p.Visibility != "PUBLIC" {
		if currentUser == nil || (currentUser.ID != p.OwnerID && currentUser.Role != "ADMIN") {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
			return
		}
	}

	p.Owner = &u
	writeJSON(w, http.StatusOK, map[string]any{"project": p})
}

func (h *Handler) ListMyProjects(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT id, owner_id, name, slug, description, cover_image_url, repository_url,
		       documentation_url, demo_url, technology_stack, hosting_type, public_url,
		       status, visibility, created_at, updated_at, published_at
		FROM projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list projects"})
		return
	}
	defer rows.Close()

	list := make([]Project, 0)
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Name, &p.Slug, &p.Description, &p.CoverImageURL,
			&p.RepositoryURL, &p.DocumentationURL, &p.DemoURL, &p.TechnologyStack,
			&p.HostingType, &p.PublicURL, &p.Status, &p.Visibility, &p.CreatedAt, &p.UpdatedAt, &p.PublishedAt,
		); err == nil {
			list = append(list, p)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"projects": list,
		"total":    len(list),
	})
}

type UpdateMetadataRequest struct {
	Description      *string  `json:"description"`
	CoverImageURL    *string  `json:"cover_image_url"`
	RepositoryURL    *string  `json:"repository_url"`
	DocumentationURL *string  `json:"documentation_url"`
	DemoURL          *string  `json:"demo_url"`
	TechnologyStack  []string `json:"technology_stack"`
}

func (h *Handler) UpdateMyProject(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	projectID := chi.URLParam(r, "id")

	// Verify ownership
	var ownerID, projectName string
	err := h.db.QueryRow(r.Context(), "SELECT owner_id, name FROM projects WHERE id = $1", projectID).Scan(&ownerID, &projectName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if ownerID != user.ID && user.Role != "ADMIN" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "You do not own this project"})
		return
	}

	var req UpdateMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	_, err = h.db.Exec(r.Context(), `
		UPDATE projects
		SET description = COALESCE($1, description),
		    cover_image_url = COALESCE($2, cover_image_url),
		    repository_url = COALESCE($3, repository_url),
		    documentation_url = COALESCE($4, documentation_url),
		    demo_url = COALESCE($5, demo_url),
		    technology_stack = CASE WHEN $6::text[] IS NOT NULL THEN $6::text[] ELSE technology_stack END,
		    updated_at = NOW()
		WHERE id = $7
	`, req.Description, req.CoverImageURL, req.RepositoryURL, req.DocumentationURL, req.DemoURL, req.TechnologyStack, projectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update project"})
		return
	}

	// Create activity: PROJECT_UPDATED
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO activities (actor_id, project_id, type, metadata, visibility)
		VALUES ($1, $2, 'PROJECT_UPDATED', $3, 'PUBLIC')
	`, user.ID, projectID, map[string]string{
		"title":   projectName + " was updated",
		"message": user.DisplayName + " updated project metadata",
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "Project metadata updated"})
}

// Admin project handlers
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT p.id, p.owner_id, p.name, p.slug, p.description, p.cover_image_url,
		       p.repository_url, p.documentation_url, p.demo_url, p.technology_stack,
		       p.hosting_type, p.public_url, p.status, p.visibility, p.created_at, p.updated_at, p.published_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list projects"})
		return
	}
	defer rows.Close()

	list := make([]Project, 0)
	for rows.Next() {
		var p Project
		var u auth.PublicUser
		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Name, &p.Slug, &p.Description, &p.CoverImageURL,
			&p.RepositoryURL, &p.DocumentationURL, &p.DemoURL, &p.TechnologyStack,
			&p.HostingType, &p.PublicURL, &p.Status, &p.Visibility, &p.CreatedAt, &p.UpdatedAt, &p.PublishedAt,
			&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
		); err == nil {
			p.Owner = &u
			list = append(list, p)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"projects": list})
}

type AdminCreateProjectRequest struct {
	OwnerID          string   `json:"owner_id"`
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	Description      string   `json:"description"`
	CoverImageURL    string   `json:"cover_image_url"`
	RepositoryURL    string   `json:"repository_url"`
	DocumentationURL string   `json:"documentation_url"`
	DemoURL          string   `json:"demo_url"`
	TechnologyStack  []string `json:"technology_stack"`
	HostingType      string   `json:"hosting_type"`
	PublicURL        string   `json:"public_url"`
	Status           string   `json:"status"`
	Visibility       string   `json:"visibility"`
}

func (h *Handler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())

	var req AdminCreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.OwnerID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Project name and owner are required"})
		return
	}

	if req.Slug == "" {
		req.Slug = GenerateSlug(req.Name)
	} else {
		req.Slug = GenerateSlug(req.Slug)
	}

	if req.HostingType == "" {
		req.HostingType = "HOSTED_HERE"
	}
	if req.Status == "" {
		req.Status = "ONLINE"
	}
	if req.Visibility == "" {
		req.Visibility = "PUBLIC"
	}
	if req.TechnologyStack == nil {
		req.TechnologyStack = []string{}
	}

	var projectID string
	err := h.db.QueryRow(r.Context(), `
		INSERT INTO projects (
			owner_id, name, slug, description, cover_image_url, repository_url,
			documentation_url, demo_url, technology_stack, hosting_type, public_url,
			status, visibility, published_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			CASE WHEN $13 = 'PUBLIC' THEN NOW() ELSE NULL END
		) RETURNING id
	`, req.OwnerID, req.Name, req.Slug, req.Description, req.CoverImageURL, req.RepositoryURL,
		req.DocumentationURL, req.DemoURL, req.TechnologyStack, req.HostingType, req.PublicURL,
		req.Status, req.Visibility).Scan(&projectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Could not create project: " + err.Error()})
		return
	}

	// Create audit log
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_CREATE_PROJECT', 'project', $2, $3)
	`, adminUser.ID, projectID, map[string]string{"name": req.Name, "slug": req.Slug})

	if req.Visibility == "PUBLIC" {
		_, _ = h.db.Exec(r.Context(), `
			INSERT INTO activities (actor_id, project_id, type, metadata, visibility)
			VALUES ($1, $2, 'PROJECT_PUBLISHED', $3, 'PUBLIC')
		`, req.OwnerID, projectID, map[string]string{"title": req.Name + " was published", "message": "Now available on Ngumpul Host"})
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Project created successfully",
		"id":      projectID,
		"slug":    req.Slug,
	})
}

type AdminUpdateProjectRequest struct {
	Name             *string   `json:"name"`
	Slug             *string   `json:"slug"`
	Description      *string   `json:"description"`
	CoverImageURL    *string   `json:"cover_image_url"`
	RepositoryURL    *string   `json:"repository_url"`
	DocumentationURL *string   `json:"documentation_url"`
	DemoURL          *string   `json:"demo_url"`
	TechnologyStack  []string  `json:"technology_stack"`
	HostingType      *string   `json:"hosting_type"`
	PublicURL        *string   `json:"public_url"`
	Status           *string   `json:"status"`
	Visibility       *string   `json:"visibility"`
	OwnerID          *string   `json:"owner_id"`
}

func (h *Handler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	projectID := chi.URLParam(r, "id")

	var req AdminUpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	var currentStatus, currentVisibility, projectName string
	err := h.db.QueryRow(r.Context(), "SELECT status, visibility, name FROM projects WHERE id = $1", projectID).Scan(&currentStatus, &currentVisibility, &projectName)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
		return
	}

	_, err = h.db.Exec(r.Context(), `
		UPDATE projects
		SET name = COALESCE($1, name),
		    slug = COALESCE($2, slug),
		    description = COALESCE($3, description),
		    cover_image_url = COALESCE($4, cover_image_url),
		    repository_url = COALESCE($5, repository_url),
		    documentation_url = COALESCE($6, documentation_url),
		    demo_url = COALESCE($7, demo_url),
		    technology_stack = CASE WHEN $8::text[] IS NOT NULL THEN $8::text[] ELSE technology_stack END,
		    hosting_type = COALESCE($9, hosting_type),
		    public_url = COALESCE($10, public_url),
		    status = COALESCE($11, status),
		    visibility = COALESCE($12, visibility),
		    owner_id = COALESCE($13, owner_id),
		    published_at = CASE WHEN $12 = 'PUBLIC' AND published_at IS NULL THEN NOW() ELSE published_at END,
		    updated_at = NOW()
		WHERE id = $14
	`, req.Name, req.Slug, req.Description, req.CoverImageURL, req.RepositoryURL,
		req.DocumentationURL, req.DemoURL, req.TechnologyStack, req.HostingType, req.PublicURL,
		req.Status, req.Visibility, req.OwnerID, projectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update project: " + err.Error()})
		return
	}

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_UPDATE_PROJECT', 'project', $2, $3)
	`, adminUser.ID, projectID, map[string]string{"name": projectName})

	// If status changed to ONLINE from something else, create activity
	if req.Status != nil && *req.Status == "ONLINE" && currentStatus != "ONLINE" {
		_, _ = h.db.Exec(r.Context(), `
			INSERT INTO activities (actor_id, project_id, type, metadata, visibility)
			VALUES ($1, $2, 'PROJECT_ONLINE', $3, 'PUBLIC')
		`, adminUser.ID, projectID, map[string]string{"title": projectName + " is online", "message": "Project is operational"})
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Project updated successfully"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
