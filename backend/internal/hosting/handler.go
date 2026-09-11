package hosting

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/project"
	"ngumpul-host/backend/internal/response"
)

// Handler handles hosting request submission, listing, and administrative workflow.
type Handler struct {
	db *pgxpool.Pool
}

// NewHandler creates a new hosting handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// Submit handles POST /api/hosting-requests
func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req SubmitRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.ProjectName = strings.TrimSpace(req.ProjectName)
	if req.ProjectName == "" || req.RepositoryURL == "" {
		response.Error(w, http.StatusBadRequest, "Project name and repository URL are required")
		return
	}

	if req.TechnologyStack == nil {
		req.TechnologyStack = []string{}
	}

	var requestID string
	err := h.db.QueryRow(r.Context(), `
		INSERT INTO hosting_requests (
			requester_id, project_name, description, repository_url,
			documentation_url, deployment_notes, technology_stack, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING')
		RETURNING id
	`, user.ID, req.ProjectName, req.Description, req.RepositoryURL, req.DocumentationURL, req.DeploymentNotes, req.TechnologyStack).Scan(&requestID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to submit hosting request")
		return
	}

	// Notify admins
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO notifications (user_id, type, title, body, data)
		SELECT id, 'HOSTING_REQUEST_SUBMITTED', 'New Hosting Request: ' || $1,
		       $2 || ' submitted a hosting request for project ' || $1,
		       $3::jsonb
		FROM users WHERE role = 'ADMIN'
	`, req.ProjectName, user.DisplayName, map[string]string{"request_id": requestID})

	response.JSON(w, http.StatusCreated, map[string]any{
		"message":    "Hosting request submitted successfully. The administrator will review your project.",
		"request_id": requestID,
	})
}

// ListMyRequests handles GET /api/me/hosting-requests
func (h *Handler) ListMyRequests(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT id, requester_id, project_name, description, repository_url,
		       documentation_url, deployment_notes, technology_stack, status,
		       admin_notes, reviewed_by, created_at, updated_at, reviewed_at
		FROM hosting_requests
		WHERE requester_id = $1
		ORDER BY created_at DESC
	`, user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	list := make([]Request, 0)
	for rows.Next() {
		var hr Request
		if err := rows.Scan(
			&hr.ID, &hr.RequesterID, &hr.ProjectName, &hr.Description, &hr.RepositoryURL,
			&hr.DocumentationURL, &hr.DeploymentNotes, &hr.TechnologyStack, &hr.Status,
			&hr.AdminNotes, &hr.ReviewedBy, &hr.CreatedAt, &hr.UpdatedAt, &hr.ReviewedAt,
		); err == nil {
			list = append(list, hr)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"requests": list})
}

// AdminListRequests handles GET /api/admin/hosting-requests
func (h *Handler) AdminListRequests(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")

	query := `
		SELECT hr.id, hr.requester_id, hr.project_name, hr.description, hr.repository_url,
		       hr.documentation_url, hr.deployment_notes, hr.technology_stack, hr.status,
		       hr.admin_notes, hr.reviewed_by, hr.created_at, hr.updated_at, hr.reviewed_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM hosting_requests hr
		JOIN users u ON u.id = hr.requester_id
	`
	args := []any{}
	if statusFilter != "" {
		args = append(args, strings.ToUpper(statusFilter))
		query += ` WHERE hr.status = $1`
	}
	query += ` ORDER BY hr.created_at DESC`

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	list := make([]Request, 0)
	for rows.Next() {
		var hr Request
		var u auth.PublicUser
		if err := rows.Scan(
			&hr.ID, &hr.RequesterID, &hr.ProjectName, &hr.Description, &hr.RepositoryURL,
			&hr.DocumentationURL, &hr.DeploymentNotes, &hr.TechnologyStack, &hr.Status,
			&hr.AdminNotes, &hr.ReviewedBy, &hr.CreatedAt, &hr.UpdatedAt, &hr.ReviewedAt,
			&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
		); err == nil {
			hr.Requester = &u
			list = append(list, hr)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"requests": list})
}

// AdminApprove handles POST /api/admin/hosting-requests/{id}/approve
func (h *Handler) AdminApprove(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	requestID := chi.URLParam(r, "id")

	var req ReviewActionPayload
	_ = json.NewDecoder(r.Body).Decode(&req)

	var requesterID, projectName string
	err := h.db.QueryRow(r.Context(), `
		UPDATE hosting_requests
		SET status = 'APPROVED', admin_notes = $1, reviewed_by = $2, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $3
		RETURNING requester_id, project_name
	`, req.AdminNotes, adminUser.ID, requestID).Scan(&requesterID, &projectName)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Request not found")
		return
	}

	// Notify requester
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO notifications (user_id, type, title, body, data)
		VALUES ($1, 'HOSTING_REQUEST_APPROVED', 'Hosting Request Approved', 'Your request for ' || $2 || ' was approved. Setup will begin shortly.', $3::jsonb)
	`, requesterID, projectName, map[string]string{"request_id": requestID})

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_APPROVED_REQUEST', 'hosting_request', $2, $3)
	`, adminUser.ID, requestID, map[string]string{"project_name": projectName})

	response.JSON(w, http.StatusOK, map[string]string{"message": "Request approved"})
}

// AdminReject handles POST /api/admin/hosting-requests/{id}/reject
func (h *Handler) AdminReject(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	requestID := chi.URLParam(r, "id")

	var req ReviewActionPayload
	_ = json.NewDecoder(r.Body).Decode(&req)

	var requesterID, projectName string
	err := h.db.QueryRow(r.Context(), `
		UPDATE hosting_requests
		SET status = 'REJECTED', admin_notes = $1, reviewed_by = $2, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $3
		RETURNING requester_id, project_name
	`, req.AdminNotes, adminUser.ID, requestID).Scan(&requesterID, &projectName)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Request not found")
		return
	}

	// Notify requester
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO notifications (user_id, type, title, body, data)
		VALUES ($1, 'HOSTING_REQUEST_REJECTED', 'Hosting Request Rejected', 'Your request for ' || $2 || ' was not approved.', $3::jsonb)
	`, requesterID, projectName, map[string]string{"request_id": requestID, "notes": req.AdminNotes})

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_REJECTED_REQUEST', 'hosting_request', $2, $3)
	`, adminUser.ID, requestID, map[string]string{"project_name": projectName, "reason": req.AdminNotes})

	response.JSON(w, http.StatusOK, map[string]string{"message": "Request rejected"})
}

// AdminComplete handles POST /api/admin/hosting-requests/{id}/complete
func (h *Handler) AdminComplete(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	requestID := chi.URLParam(r, "id")

	var req ReviewActionPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.PublicURL == "" {
		response.Error(w, http.StatusBadRequest, "Public URL is required to complete hosting setup")
		return
	}

	var hr Request
	err := h.db.QueryRow(r.Context(), `
		SELECT id, requester_id, project_name, description, repository_url, documentation_url, technology_stack
		FROM hosting_requests
		WHERE id = $1
	`, requestID).Scan(&hr.ID, &hr.RequesterID, &hr.ProjectName, &hr.Description, &hr.RepositoryURL, &hr.DocumentationURL, &hr.TechnologyStack)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Request not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Mark request COMPLETED
	_, _ = h.db.Exec(r.Context(), `
		UPDATE hosting_requests
		SET status = 'COMPLETED', admin_notes = $1, reviewed_by = $2, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`, req.AdminNotes, adminUser.ID, requestID)

	slug := project.GenerateSlug(hr.ProjectName)

	// Create published project
	var projectID string
	err = h.db.QueryRow(r.Context(), `
		INSERT INTO projects (
			owner_id, name, slug, description, repository_url, documentation_url,
			technology_stack, hosting_type, public_url, status, visibility, published_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, 'HOSTED_HERE', $8, 'ONLINE', 'PUBLIC', NOW()
		)
		ON CONFLICT (slug) DO UPDATE
		SET public_url = EXCLUDED.public_url, status = 'ONLINE', visibility = 'PUBLIC', updated_at = NOW()
		RETURNING id
	`, hr.RequesterID, hr.ProjectName, slug, hr.Description, hr.RepositoryURL, hr.DocumentationURL,
		hr.TechnologyStack, req.PublicURL).Scan(&projectID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to register project: "+err.Error())
		return
	}

	// Notify requester
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO notifications (user_id, type, title, body, data)
		VALUES ($1, 'PROJECT_PUBLISHED', 'Your project ' || $2 || ' is now live!', 'Accessible at ' || $3, $4::jsonb)
	`, hr.RequesterID, hr.ProjectName, req.PublicURL, map[string]string{"project_slug": slug, "public_url": req.PublicURL})

	// Public Activity
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO activities (actor_id, project_id, type, metadata, visibility)
		VALUES ($1, $2, 'PROJECT_PUBLISHED', $3, 'PUBLIC')
	`, hr.RequesterID, projectID, map[string]string{"title": hr.ProjectName + " went online", "message": "Hosted at " + req.PublicURL})

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_PUBLISHED_PROJECT', 'project', $2, $3)
	`, adminUser.ID, projectID, map[string]string{"name": hr.ProjectName, "public_url": req.PublicURL})

	response.JSON(w, http.StatusOK, map[string]any{
		"message":    "Hosting setup completed and project published",
		"project_id": projectID,
		"slug":       slug,
		"public_url": req.PublicURL,
	})
}
