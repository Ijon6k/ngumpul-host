package comment

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/response"
)

// Handler handles comment retrieval, creation, deletion, and moderation.
type Handler struct {
	db          *pgxpool.Pool
	rateLimitMu sync.Mutex
	userPostMap map[string][]time.Time
}

// NewHandler creates a comment handler with rate-limiting support.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		db:          db,
		userPostMap: make(map[string][]time.Time),
	}
}

// checkRateLimit enforces maximum 5 comments per minute per user.
func (h *Handler) checkRateLimit(userID string) bool {
	h.rateLimitMu.Lock()
	defer h.rateLimitMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)

	times, exists := h.userPostMap[userID]
	if !exists {
		h.userPostMap[userID] = []time.Time{now}
		return true
	}

	validTimes := make([]time.Time, 0, len(times))
	for _, t := range times {
		if t.After(cutoff) {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= 5 {
		h.userPostMap[userID] = validTimes
		return false
	}

	validTimes = append(validTimes, now)
	h.userPostMap[userID] = validTimes
	return true
}

// ListByProject handles GET /api/projects/{slug}/comments
func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.Error(w, http.StatusBadRequest, "Project slug is required")
		return
	}

	var projectID string
	err := h.db.QueryRow(r.Context(), "SELECT id FROM projects WHERE slug = $1", slug).Scan(&projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT c.id, c.project_id, c.author_id, c.content, c.created_at, c.updated_at, c.deleted_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.project_id = $1
		ORDER BY c.created_at ASC
	`, projectID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to load comments")
		return
	}
	defer rows.Close()

	list := make([]Comment, 0)
	for rows.Next() {
		var c Comment
		var u auth.PublicUser
		if err := rows.Scan(
			&c.ID, &c.ProjectID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
			&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
		); err == nil {
			c.Author = &u
			if c.DeletedAt != nil {
				c.IsDeleted = true
				c.Content = "This comment was removed."
			}
			list = append(list, c)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"comments": list})
}

// Create handles POST /api/projects/{slug}/comments
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUser(r.Context())
	if caller == nil {
		response.Error(w, http.StatusUnauthorized, "Authentication required to comment")
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.Error(w, http.StatusBadRequest, "Project slug is required")
		return
	}

	if !h.checkRateLimit(caller.ID) {
		response.Error(w, http.StatusTooManyRequests, "You are posting comments too quickly. Please wait a moment.")
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		response.Error(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}
	if len(content) > 1000 {
		response.Error(w, http.StatusBadRequest, "Comment exceeds maximum length of 1000 characters")
		return
	}

	// Sanitize plain text (escape dangerous HTML)
	sanitized := html.EscapeString(content)

	var projectID, projectName string
	err := h.db.QueryRow(r.Context(), "SELECT id, name FROM projects WHERE slug = $1", slug).Scan(&projectID, &projectName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	var c Comment
	err = h.db.QueryRow(r.Context(), `
		INSERT INTO comments (project_id, author_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, project_id, author_id, content, created_at, updated_at
	`, projectID, caller.ID, sanitized).Scan(&c.ID, &c.ProjectID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to save comment")
		return
	}

	c.Author = &auth.PublicUser{
		ID:          caller.ID,
		Username:    caller.Username,
		DisplayName: caller.DisplayName,
		AvatarURL:   caller.AvatarURL,
		Bio:         caller.Bio,
		Role:        caller.Role,
		CreatedAt:   caller.CreatedAt,
	}

	// Record public activity
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO activities (actor_id, project_id, type, metadata, visibility)
		VALUES ($1, $2, 'COMMENT_ADDED', $3, 'PUBLIC')
	`, caller.ID, projectID, map[string]string{
		"title":   fmt.Sprintf("%s commented on %s", caller.DisplayName, projectName),
		"message": sanitized,
	})

	response.JSON(w, http.StatusCreated, map[string]any{"comment": c})
}

// Delete handles DELETE /api/comments/{id}
// Authorization: Caller must be ADMIN or the OWNER of the project hosting the comment.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUser(r.Context())
	if caller == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentID := chi.URLParam(r, "id")
	if commentID == "" {
		response.Error(w, http.StatusBadRequest, "Comment ID is required")
		return
	}

	// Lookup comment and the owner of the hosting project
	var projectOwnerID, projectID string
	var alreadyDeleted bool
	err := h.db.QueryRow(r.Context(), `
		SELECT p.owner_id, c.project_id, (c.deleted_at IS NOT NULL)
		FROM comments c
		JOIN projects p ON p.id = c.project_id
		WHERE c.id = $1
	`, commentID).Scan(&projectOwnerID, &projectID, &alreadyDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Comment not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if alreadyDeleted {
		response.JSON(w, http.StatusOK, map[string]string{"message": "Comment already removed"})
		return
	}

	// Check permissions: Admin or Project Owner
	if caller.Role != "ADMIN" && caller.ID != projectOwnerID {
		response.Error(w, http.StatusForbidden, "Only administrators or the project owner may delete this comment")
		return
	}

	// Soft delete comment
	_, err = h.db.Exec(r.Context(), `
		UPDATE comments
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, commentID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to remove comment")
		return
	}

	// Audit trail record
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'DELETE_COMMENT', 'comment', $2, $3)
	`, caller.ID, commentID, map[string]string{
		"project_id": projectID,
		"role":       caller.Role,
	})

	response.JSON(w, http.StatusOK, map[string]string{"message": "Comment removed successfully"})
}

// AdminList handles GET /api/admin/comments
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	projectID := r.URL.Query().Get("project_id")

	query := `
		SELECT c.id, c.project_id, p.name, p.slug, c.author_id, c.content, c.created_at, c.updated_at, c.deleted_at,
		       (SELECT COUNT(*) FROM reports r WHERE r.target_type = 'COMMENT' AND r.target_id = c.id AND r.status = 'OPEN') as report_count,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at
		FROM comments c
		JOIN projects p ON p.id = c.project_id
		JOIN users u ON u.id = c.author_id
		WHERE 1=1
	`
	args := make([]any, 0)
	argIdx := 1

	if projectID != "" {
		query += fmt.Sprintf(" AND c.project_id = $%d", argIdx)
		args = append(args, projectID)
		argIdx++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (c.content ILIKE $%d OR u.username ILIKE $%d OR p.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += " ORDER BY c.created_at DESC LIMIT 100"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	list := make([]AdminCommentRow, 0)
	for rows.Next() {
		var row AdminCommentRow
		var u auth.PublicUser
		if err := rows.Scan(
			&row.ID, &row.ProjectID, &row.ProjectName, &row.ProjectSlug, &row.AuthorID,
			&row.Content, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt, &row.ReportCount,
			&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
		); err == nil {
			row.Author = &u
			row.IsDeleted = (row.DeletedAt != nil)
			list = append(list, row)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"comments": list})
}
