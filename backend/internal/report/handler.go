package report

import (
	"encoding/json"
	"errors"
	"fmt"
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

var validReasons = map[string]bool{
	"SPAM":                 true,
	"ABUSE_HARASSMENT":     true,
	"INAPPROPRIATE":        true,
	"MALICIOUS_SUSPICIOUS": true,
	"OTHER":                true,
}

// Handler coordinates report submission and administrative moderation.
type Handler struct {
	db          *pgxpool.Pool
	rateLimitMu sync.Mutex
	userPostMap map[string][]time.Time
}

// NewHandler creates a report handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		db:          db,
		userPostMap: make(map[string][]time.Time),
	}
}

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

	if len(validTimes) >= 3 {
		h.userPostMap[userID] = validTimes
		return false
	}

	validTimes = append(validTimes, now)
	h.userPostMap[userID] = validTimes
	return true
}

// ReportComment handles POST /api/comments/{id}/report
func (h *Handler) ReportComment(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUser(r.Context())
	if caller == nil {
		response.Error(w, http.StatusUnauthorized, "Authentication required to submit report")
		return
	}

	commentID := chi.URLParam(r, "id")
	if commentID == "" {
		response.Error(w, http.StatusBadRequest, "Comment ID is required")
		return
	}

	if !h.checkRateLimit(caller.ID) {
		response.Error(w, http.StatusTooManyRequests, "Too many reports submitted. Please wait a moment.")
		return
	}

	var req CreateReportPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	reason := strings.ToUpper(strings.TrimSpace(req.Reason))
	if !validReasons[reason] {
		response.Error(w, http.StatusBadRequest, "Invalid report reason")
		return
	}

	// Verify comment exists and get project_id
	var projectID string
	err := h.db.QueryRow(r.Context(), "SELECT project_id FROM comments WHERE id = $1", commentID).Scan(&projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "Comment not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Prevent duplicate open reports from same user
	var existingCount int
	_ = h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM reports
		WHERE reporter_id = $1 AND target_type = 'COMMENT' AND target_id = $2 AND status = 'OPEN'
	`, caller.ID, commentID).Scan(&existingCount)
	if existingCount > 0 {
		response.JSON(w, http.StatusOK, map[string]string{"message": "Your report has already been received and is pending review."})
		return
	}

	details := strings.TrimSpace(req.Details)
	if len(details) > 500 {
		details = details[:500]
	}

	_, err = h.db.Exec(r.Context(), `
		INSERT INTO reports (reporter_id, target_type, target_id, project_id, reason, details, status)
		VALUES ($1, 'COMMENT', $2, $3, $4, $5, 'OPEN')
	`, caller.ID, commentID, projectID, reason, details)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to submit report")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "Thanks. Your report has been submitted."})
}

// ReportProject handles POST /api/projects/{slug}/report
func (h *Handler) ReportProject(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUser(r.Context())
	if caller == nil {
		response.Error(w, http.StatusUnauthorized, "Authentication required to submit report")
		return
	}

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.Error(w, http.StatusBadRequest, "Project slug is required")
		return
	}

	if !h.checkRateLimit(caller.ID) {
		response.Error(w, http.StatusTooManyRequests, "Too many reports submitted. Please wait a moment.")
		return
	}

	var req CreateReportPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	reason := strings.ToUpper(strings.TrimSpace(req.Reason))
	if !validReasons[reason] {
		response.Error(w, http.StatusBadRequest, "Invalid report reason")
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

	var existingCount int
	_ = h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM reports
		WHERE reporter_id = $1 AND target_type = 'PROJECT' AND target_id = $2 AND status = 'OPEN'
	`, caller.ID, projectID).Scan(&existingCount)
	if existingCount > 0 {
		response.JSON(w, http.StatusOK, map[string]string{"message": "Your report has already been received and is pending review."})
		return
	}

	details := strings.TrimSpace(req.Details)
	if len(details) > 500 {
		details = details[:500]
	}

	_, err = h.db.Exec(r.Context(), `
		INSERT INTO reports (reporter_id, target_type, target_id, project_id, reason, details, status)
		VALUES ($1, 'PROJECT', $2, $2, $3, $4, 'OPEN')
	`, caller.ID, projectID, reason, details)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to submit report")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "Thanks. Your report has been submitted."})
}

// AdminList handles GET /api/admin/reports
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	targetFilter := r.URL.Query().Get("target_type")

	query := `
		SELECT r.id, r.reporter_id, r.target_type, r.target_id, r.project_id,
		       COALESCE(p.name, ''), COALESCE(p.slug, ''),
		       r.reason, r.details, r.status, r.resolved_by, r.resolved_at, r.action_notes,
		       r.created_at, r.updated_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.bio, u.role, u.created_at,
		       CASE
		         WHEN r.target_type = 'COMMENT' THEN COALESCE((SELECT content FROM comments WHERE id = r.target_id), 'Comment removed')
		         ELSE COALESCE(p.description, '')
		       END as target_snippet
		FROM reports r
		JOIN users u ON u.id = r.reporter_id
		LEFT JOIN projects p ON p.id = r.project_id
		WHERE 1=1
	`
	args := make([]any, 0)
	argIdx := 1

	if statusFilter != "" && statusFilter != "ALL" {
		query += fmt.Sprintf(" AND r.status = $%d", argIdx)
		args = append(args, strings.ToUpper(statusFilter))
		argIdx++
	}

	if targetFilter != "" && targetFilter != "ALL" {
		query += fmt.Sprintf(" AND r.target_type = $%d", argIdx)
		args = append(args, strings.ToUpper(targetFilter))
		argIdx++
	}

	query += " ORDER BY r.created_at DESC LIMIT 100"

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	list := make([]Report, 0)
	for rows.Next() {
		var rep Report
		var u auth.PublicUser
		if err := rows.Scan(
			&rep.ID, &rep.ReporterID, &rep.TargetType, &rep.TargetID, &rep.ProjectID,
			&rep.ProjectName, &rep.ProjectSlug,
			&rep.Reason, &rep.Details, &rep.Status, &rep.ResolvedBy, &rep.ResolvedAt, &rep.ActionNotes,
			&rep.CreatedAt, &rep.UpdatedAt,
			&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt,
			&rep.TargetSnippet,
		); err == nil {
			rep.Reporter = &u
			list = append(list, rep)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"reports": list})
}

// AdminResolve handles POST /api/admin/reports/{id}/resolve
func (h *Handler) AdminResolve(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	reportID := chi.URLParam(r, "id")

	var req ResolveReportPayload
	_ = json.NewDecoder(r.Body).Decode(&req)

	var targetType, targetID string
	err := h.db.QueryRow(r.Context(), `
		UPDATE reports
		SET status = 'RESOLVED', resolved_by = $1, resolved_at = NOW(), action_notes = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING target_type, target_id
	`, adminUser.ID, req.Notes, reportID).Scan(&targetType, &targetID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Report not found")
		return
	}

	// If action requested is DELETE_COMMENT, soft-delete the comment
	if req.Action == "DELETE_COMMENT" && targetType == "COMMENT" {
		_, _ = h.db.Exec(r.Context(), `
			UPDATE comments
			SET deleted_at = NOW(), updated_at = NOW()
			WHERE id = $1
		`, targetID)

		_, _ = h.db.Exec(r.Context(), `
			INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
			VALUES ($1, 'ADMIN_DELETED_REPORTED_COMMENT', 'comment', $2, $3)
		`, adminUser.ID, targetID, map[string]string{"report_id": reportID, "notes": req.Notes})
	}

	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_RESOLVED_REPORT', 'report', $2, $3)
	`, adminUser.ID, reportID, map[string]string{"action": req.Action, "notes": req.Notes})

	response.JSON(w, http.StatusOK, map[string]string{"message": "Report resolved"})
}

// AdminDismiss handles POST /api/admin/reports/{id}/dismiss
func (h *Handler) AdminDismiss(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	reportID := chi.URLParam(r, "id")

	var req DismissReportPayload
	_ = json.NewDecoder(r.Body).Decode(&req)

	_, err := h.db.Exec(r.Context(), `
		UPDATE reports
		SET status = 'DISMISSED', resolved_by = $1, resolved_at = NOW(), action_notes = $2, updated_at = NOW()
		WHERE id = $3
	`, adminUser.ID, req.Notes, reportID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Report not found")
		return
	}

	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_DISMISSED_REPORT', 'report', $2, $3)
	`, adminUser.ID, reportID, map[string]string{"notes": req.Notes})

	response.JSON(w, http.StatusOK, map[string]string{"message": "Report dismissed"})
}
