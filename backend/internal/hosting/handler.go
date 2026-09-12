package hosting

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/project"
	"ngumpul-host/backend/internal/response"
)

var (
	subdomainRegex     = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
	reservedSubdomains = map[string]bool{
		"api": true, "admin": true, "www": true, "mail": true, "status": true,
		"me": true, "console": true, "auth": true, "login": true, "register": true,
		"static": true, "uploads": true, "assets": true, "internal": true, "public": true,
	}
)

// Handler handles hosting request submission, listing, and administrative workflow.
type Handler struct {
	db *pgxpool.Pool
}

// NewHandler creates a new hosting handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// CheckSubdomain handles GET /api/hosting-requests/check-subdomain?subdomain=...
func (h *Handler) CheckSubdomain(w http.ResponseWriter, r *http.Request) {
	sub := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("subdomain")))
	if sub == "" {
		response.Error(w, http.StatusBadRequest, "Subdomain parameter is required")
		return
	}

	if len(sub) < 2 {
		response.JSON(w, http.StatusOK, map[string]any{
			"available": false,
			"subdomain": sub,
			"message":   "Subdomain must be at least 2 characters",
		})
		return
	}

	if !subdomainRegex.MatchString(sub) {
		response.JSON(w, http.StatusOK, map[string]any{
			"available": false,
			"subdomain": sub,
			"message":   "Only lowercase letters, numbers, and hyphens allowed (cannot start or end with hyphen)",
		})
		return
	}

	if reservedSubdomains[sub] {
		response.JSON(w, http.StatusOK, map[string]any{
			"available": false,
			"subdomain": sub,
			"message":   "This subdomain is reserved for system services",
		})
		return
	}

	// Check if already taken in projects (as slug or in public_url)
	var count int
	err := h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM projects
		WHERE LOWER(slug) = $1 OR public_url ILIKE '%' || $1 || '.%'
	`, sub).Scan(&count)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	if count > 0 {
		response.JSON(w, http.StatusOK, map[string]any{
			"available": false,
			"subdomain": sub,
			"message":   "Subdomain is already taken by an existing project",
		})
		return
	}

	// Check if already requested in pending/approved/setup requests
	err = h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM hosting_requests
		WHERE LOWER(subdomain) = $1 AND status IN ('PENDING', 'REVIEWING', 'APPROVED', 'SETUP')
	`, sub).Scan(&count)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	if count > 0 {
		response.JSON(w, http.StatusOK, map[string]any{
			"available": false,
			"subdomain": sub,
			"message":   "Subdomain is currently reserved in a pending hosting request",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"available": true,
		"subdomain": sub,
		"message":   "Subdomain is available",
	})
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
	req.Subdomain = strings.ToLower(strings.TrimSpace(req.Subdomain))
	if req.Subdomain == "" {
		req.Subdomain = project.GenerateSlug(req.ProjectName)
	}
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
			requester_id, project_name, subdomain, description, readme, repository_url,
			documentation_url, deployment_notes, technology_stack, request_type, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'NEW_PROJECT', 'PENDING')
		RETURNING id
	`, user.ID, req.ProjectName, req.Subdomain, req.Description, req.Readme, req.RepositoryURL, req.DocumentationURL, req.DeploymentNotes, req.TechnologyStack).Scan(&requestID)
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
		SELECT id, requester_id, project_id, request_type, project_name, subdomain, description, readme, repository_url,
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
			&hr.ID, &hr.RequesterID, &hr.ProjectID, &hr.RequestType, &hr.ProjectName, &hr.Subdomain, &hr.Description, &hr.Readme, &hr.RepositoryURL,
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
		SELECT hr.id, hr.requester_id, hr.project_id, hr.request_type, hr.project_name, hr.subdomain, hr.description, hr.readme, hr.repository_url,
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
			&hr.ID, &hr.RequesterID, &hr.ProjectID, &hr.RequestType, &hr.ProjectName, &hr.Subdomain, &hr.Description, &hr.Readme, &hr.RepositoryURL,
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

// RequestSubdomainChange handles POST /api/projects/{id}/request-subdomain-change
func (h *Handler) RequestSubdomainChange(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	projectID := chi.URLParam(r, "id")

	var payload SubdomainChangeRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newSub := strings.ToLower(strings.TrimSpace(payload.NewSubdomain))
	if len(newSub) < 2 || !subdomainRegex.MatchString(newSub) {
		response.Error(w, http.StatusBadRequest, "Invalid subdomain format. Must be 2-63 lowercase alphanumeric characters or hyphens.")
		return
	}
	if reservedSubdomains[newSub] {
		response.Error(w, http.StatusBadRequest, "This subdomain is reserved for system services.")
		return
	}

	// Fetch project and verify ownership
	var projectName, currentSlug string
	var ownerID string
	err := h.db.QueryRow(r.Context(), `
		SELECT name, slug, owner_id FROM projects WHERE id = $1
	`, projectID).Scan(&projectName, &currentSlug, &ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if ownerID != user.ID && user.Role != "ADMIN" {
		response.Error(w, http.StatusForbidden, "You are not authorized to modify this project")
		return
	}

	if newSub == currentSlug {
		response.Error(w, http.StatusBadRequest, "New subdomain is identical to current subdomain")
		return
	}

	// Collision check against projects
	var projectCount int
	_ = h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM projects WHERE slug = $1 AND id != $2`, newSub, projectID).Scan(&projectCount)
	if projectCount > 0 {
		response.Error(w, http.StatusConflict, "Subdomain is already in use by another project")
		return
	}

	// Collision check against active pending hosting requests
	var pendingCount int
	_ = h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM hosting_requests 
		WHERE subdomain = $1 AND status = 'PENDING'
	`, newSub).Scan(&pendingCount)
	if pendingCount > 0 {
		response.Error(w, http.StatusConflict, "Subdomain is already reserved by an active request")
		return
	}

	// Check if this project already has a pending change request
	var existingPending int
	_ = h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM hosting_requests
		WHERE project_id = $1 AND request_type = 'SUBDOMAIN_CHANGE' AND status = 'PENDING'
	`, projectID).Scan(&existingPending)
	if existingPending > 0 {
		response.Error(w, http.StatusBadRequest, "A subdomain change request is already pending for this project")
		return
	}

	var requestID string
	err = h.db.QueryRow(r.Context(), `
		INSERT INTO hosting_requests (
			requester_id, project_id, request_type, project_name, subdomain,
			description, status
		) VALUES ($1, $2, 'SUBDOMAIN_CHANGE', $3, $4, $5, 'PENDING')
		RETURNING id
	`, user.ID, projectID, projectName, newSub, strings.TrimSpace(payload.Reason)).Scan(&requestID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to submit request")
		return
	}

	// Notify admins
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO notifications (user_id, type, title, body, data)
		SELECT id, 'SUBDOMAIN_CHANGE_REQUESTED', 'Subdomain Change Requested: ' || $1,
		       $2 || ' requested to change subdomain for project ' || $1 || ' to ' || $3,
		       $4::jsonb
		FROM users WHERE role = 'ADMIN'
	`, projectName, user.DisplayName, newSub, map[string]string{"request_id": requestID, "project_id": projectID, "subdomain": newSub})

	response.JSON(w, http.StatusCreated, map[string]any{
		"message":    "Subdomain change request submitted successfully. Awaiting administrator review.",
		"request_id": requestID,
		"subdomain":  newSub,
	})
}

// AdminApprove handles POST /api/admin/hosting-requests/{id}/approve
func (h *Handler) AdminApprove(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	requestID := chi.URLParam(r, "id")

	var req ReviewActionPayload
	_ = json.NewDecoder(r.Body).Decode(&req)

	req.Subdomain = strings.ToLower(strings.TrimSpace(req.Subdomain))

	var requesterID, projectName, reqType, existingSubdomain string
	var projID *string
	err := h.db.QueryRow(r.Context(), `
		SELECT requester_id, project_name, request_type, project_id, subdomain
		FROM hosting_requests
		WHERE id = $1
	`, requestID).Scan(&requesterID, &projectName, &reqType, &projID, &existingSubdomain)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Request not found")
		return
	}

	subdomainToUse := existingSubdomain
	if req.Subdomain != "" {
		subdomainToUse = req.Subdomain
	}

	// If this is a SUBDOMAIN_CHANGE request for an existing project, approve and immediately apply
	if reqType == "SUBDOMAIN_CHANGE" && projID != nil {
		newURL := "https://" + subdomainToUse + ".ngumpul.local"
		_, err = h.db.Exec(r.Context(), `
			UPDATE projects
			SET slug = $1, public_url = $2, updated_at = NOW()
			WHERE id = $3
		`, subdomainToUse, newURL, *projID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update project subdomain: "+err.Error())
			return
		}

		_, _ = h.db.Exec(r.Context(), `
			UPDATE hosting_requests
			SET status = 'COMPLETED',
			    subdomain = $1,
			    admin_notes = $2,
			    reviewed_by = $3,
			    reviewed_at = NOW(),
			    updated_at = NOW()
			WHERE id = $4
		`, subdomainToUse, req.AdminNotes, adminUser.ID, requestID)

		// Notify requester
		_, _ = h.db.Exec(r.Context(), `
			INSERT INTO notifications (user_id, type, title, body, data)
			VALUES ($1, 'SUBDOMAIN_CHANGED', 'Subdomain Change Approved', 'Your project ' || $2 || ' is now available at ' || $3, $4::jsonb)
		`, requesterID, projectName, newURL, map[string]string{"project_id": *projID, "subdomain": subdomainToUse, "public_url": newURL})

		// Audit
		_, _ = h.db.Exec(r.Context(), `
			INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
			VALUES ($1, 'ADMIN_APPROVED_SUBDOMAIN_CHANGE', 'project', $2, $3)
		`, adminUser.ID, *projID, map[string]string{"subdomain": subdomainToUse, "public_url": newURL})

		response.JSON(w, http.StatusOK, map[string]string{"message": "Subdomain change approved and applied"})
		return
	}

	err = h.db.QueryRow(r.Context(), `
		UPDATE hosting_requests
		SET status = 'APPROVED',
		    subdomain = CASE WHEN $1 != '' THEN $1 ELSE subdomain END,
		    admin_notes = $2,
		    reviewed_by = $3,
		    reviewed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $4
		RETURNING requester_id, project_name
	`, req.Subdomain, req.AdminNotes, adminUser.ID, requestID).Scan(&requesterID, &projectName)
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
		SELECT id, requester_id, project_name, subdomain, description, readme, repository_url, documentation_url, technology_stack
		FROM hosting_requests
		WHERE id = $1
	`, requestID).Scan(&hr.ID, &hr.RequesterID, &hr.ProjectName, &hr.Subdomain, &hr.Description, &hr.Readme, &hr.RepositoryURL, &hr.DocumentationURL, &hr.TechnologyStack)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Request not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if req.Subdomain != "" {
		hr.Subdomain = strings.ToLower(strings.TrimSpace(req.Subdomain))
	}

	// Mark request COMPLETED
	_, _ = h.db.Exec(r.Context(), `
		UPDATE hosting_requests
		SET status = 'COMPLETED',
		    subdomain = CASE WHEN $1 != '' THEN $1 ELSE subdomain END,
		    admin_notes = $2,
		    reviewed_by = $3,
		    reviewed_at = NOW(),
		    updated_at = NOW()
		WHERE id = $4
	`, hr.Subdomain, req.AdminNotes, adminUser.ID, requestID)

	slug := hr.Subdomain
	if slug == "" {
		slug = project.GenerateSlug(hr.ProjectName)
	}

	// Create published project
	var projectID string
	err = h.db.QueryRow(r.Context(), `
		INSERT INTO projects (
			owner_id, name, slug, description, readme, repository_url, documentation_url,
			technology_stack, hosting_type, public_url, status, visibility, published_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, 'HOSTED_HERE', $9, 'ONLINE', 'PUBLIC', NOW()
		)
		ON CONFLICT (slug) DO UPDATE
		SET public_url = EXCLUDED.public_url,
		    readme = CASE WHEN EXCLUDED.readme != '' THEN EXCLUDED.readme ELSE projects.readme END,
		    status = 'ONLINE',
		    visibility = 'PUBLIC',
		    updated_at = NOW()
		RETURNING id
	`, hr.RequesterID, hr.ProjectName, slug, hr.Description, hr.Readme, hr.RepositoryURL, hr.DocumentationURL,
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
