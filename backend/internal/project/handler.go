package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/response"
)

// ProjectActivity represents an event specifically associated with a project.
type ProjectActivity struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
	ActorName *string        `json:"actor_name"`
	ActorUser *string        `json:"actor_username"`
}

// Handler handles public, user, and administrative project management endpoints.
type Handler struct {
	db *pgxpool.Pool
}

// NewHandler creates a new project handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// ListPublic handles GET /api/projects
func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	techFilter := r.URL.Query().Get("tech")
	typeFilter := r.URL.Query().Get("type")
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	pageStr := strings.TrimSpace(r.URL.Query().Get("page"))
	limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))

	baseWhere := "WHERE p.visibility = 'PUBLIC'"
	args := []any{}
	argIdx := 1

	if statusFilter != "" {
		args = append(args, strings.ToUpper(statusFilter))
		baseWhere += fmt.Sprintf(" AND p.status = $%d", argIdx)
		argIdx++
	}
	if typeFilter != "" {
		args = append(args, strings.ToUpper(typeFilter))
		baseWhere += fmt.Sprintf(" AND p.hosting_type = $%d", argIdx)
		argIdx++
	}
	if techFilter != "" {
		args = append(args, techFilter)
		baseWhere += fmt.Sprintf(" AND $%d = ANY(p.technology_stack)", argIdx)
		argIdx++
	}
	if searchQuery != "" {
		term := "%" + strings.ToLower(searchQuery) + "%"
		args = append(args, term)
		baseWhere += fmt.Sprintf(" AND (LOWER(p.name) LIKE $%d OR LOWER(p.description) LIKE $%d OR LOWER(u.username) LIKE $%d OR LOWER(u.display_name) LIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		argIdx++
	}

	// Count total records matching filter
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		%s
	`, baseWhere)

	var totalCount int
	if err := h.db.QueryRow(r.Context(), countQuery, args...).Scan(&totalCount); err != nil {
		response.Error(w, http.StatusInternalServerError, "Database count error")
		return
	}

	dataQuery := fmt.Sprintf(`
		SELECT p.id, p.owner_id, p.name, p.slug, p.description, p.cover_image_url,
		       p.repository_url, p.documentation_url, p.demo_url, p.technology_stack,
		       p.hosting_type, p.public_url, p.status, p.visibility, p.created_at, p.updated_at, p.published_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM projects p
		JOIN users u ON u.id = p.owner_id
		%s
		ORDER BY p.published_at DESC NULLS LAST, p.created_at DESC
	`, baseWhere)

	page := 1
	limit := 12
	if pageStr != "" || limitStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 50 {
				limit = 50
			} else {
				limit = l
			}
		}
		offset := (page - 1) * limit
		dataQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := h.db.Query(r.Context(), dataQuery, args...)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to list projects")
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

	totalPages := 1
	if limit > 0 {
		totalPages = (totalCount + limit - 1) / limit
		if totalPages < 1 {
			totalPages = 1
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"projects":    list,
		"total":       totalCount,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetBySlug handles GET /api/projects/{slug}
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
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Authorization check for unpublished projects
	if p.Visibility != "PUBLIC" {
		if currentUser == nil || (currentUser.ID != p.OwnerID && currentUser.Role != "ADMIN") {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
	}

	p.Owner = &u

	// Query project-specific public activity stream
	activities := make([]ProjectActivity, 0)
	actRows, err := h.db.Query(r.Context(), `
		SELECT a.id, a.type, a.metadata, a.created_at, u.display_name, u.username
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		WHERE a.project_id = $1 AND a.visibility = 'PUBLIC'
		ORDER BY a.created_at DESC
		LIMIT 15
	`, p.ID)
	if err == nil {
		defer actRows.Close()
		for actRows.Next() {
			var act ProjectActivity
			var metaJSON []byte
			if err := actRows.Scan(
				&act.ID, &act.Type, &metaJSON, &act.CreatedAt,
				&act.ActorName, &act.ActorUser,
			); err == nil {
				_ = json.Unmarshal(metaJSON, &act.Metadata)
				activities = append(activities, act)
			}
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"project":    p,
		"activities": activities,
	})
}

// ListMyProjects handles GET /api/me/projects
func (h *Handler) ListMyProjects(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
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
		response.Error(w, http.StatusInternalServerError, "Failed to list projects")
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

	response.JSON(w, http.StatusOK, map[string]any{
		"projects": list,
		"total":    len(list),
	})
}

// UpdateMyProject handles PATCH /api/me/projects/{id}
func (h *Handler) UpdateMyProject(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	var ownerID, projectName string
	err := h.db.QueryRow(r.Context(), "SELECT owner_id, name FROM projects WHERE id = $1", projectID).Scan(&ownerID, &projectName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if ownerID != user.ID && user.Role != "ADMIN" {
		response.Error(w, http.StatusForbidden, "You do not own this project")
		return
	}

	var req UpdateMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
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
		response.Error(w, http.StatusInternalServerError, "Failed to update project")
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

	response.JSON(w, http.StatusOK, map[string]string{"message": "Project metadata updated"})
}

// AdminList handles GET /api/admin/projects
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
		response.Error(w, http.StatusInternalServerError, "Failed to list projects")
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

	response.JSON(w, http.StatusOK, map[string]any{"projects": list})
}

// AdminCreate handles POST /api/admin/projects
func (h *Handler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())

	var req AdminCreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.OwnerID == "" {
		response.Error(w, http.StatusBadRequest, "Project name and owner are required")
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
		response.Error(w, http.StatusInternalServerError, "Could not create project: "+err.Error())
		return
	}

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

	response.JSON(w, http.StatusCreated, map[string]any{
		"message": "Project created successfully",
		"id":      projectID,
		"slug":    req.Slug,
	})
}

// AdminUpdate handles PATCH /api/admin/projects/{id}
func (h *Handler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	projectID := chi.URLParam(r, "id")

	var req AdminUpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var currentStatus, currentVisibility, projectName string
	err := h.db.QueryRow(r.Context(), "SELECT status, visibility, name FROM projects WHERE id = $1", projectID).Scan(&currentStatus, &currentVisibility, &projectName)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Project not found")
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
		response.Error(w, http.StatusInternalServerError, "Failed to update project: "+err.Error())
		return
	}

	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_UPDATE_PROJECT', 'project', $2, $3)
	`, adminUser.ID, projectID, map[string]string{"name": projectName})

	if req.Status != nil && *req.Status == "ONLINE" && currentStatus != "ONLINE" {
		_, _ = h.db.Exec(r.Context(), `
			INSERT INTO activities (actor_id, project_id, type, metadata, visibility)
			VALUES ($1, $2, 'PROJECT_ONLINE', $3, 'PUBLIC')
		`, adminUser.ID, projectID, map[string]string{"title": projectName + " is online", "message": "Project is operational"})
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Project updated successfully"})
}
