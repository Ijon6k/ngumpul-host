package admin

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/response"
	"ngumpul-host/backend/internal/system"
)

var startTime = time.Now()

// Handler handles administrative statistics, moderation, audit, and deep system telemetry.
type Handler struct {
	db *pgxpool.Pool
}

// NewHandler creates a new administrative handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// GetStats handles GET /api/admin/stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	var pendingRequests, totalProjects, onlineProjects, totalMembers int

	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM hosting_requests WHERE status = 'PENDING'").Scan(&pendingRequests)
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM projects").Scan(&totalProjects)
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM projects WHERE status = 'ONLINE'").Scan(&onlineProjects)
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM users").Scan(&totalMembers)

	response.JSON(w, http.StatusOK, map[string]any{
		"pending_requests": pendingRequests,
		"total_projects":   totalProjects,
		"online_projects":  onlineProjects,
		"offline_projects": totalProjects - onlineProjects,
		"total_members":    totalMembers,
		"uptime_seconds":   int(time.Since(startTime).Seconds()),
	})
}

// ListUsers handles GET /api/admin/users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT u.id, u.username, u.email, u.display_name, u.avatar_url, u.bio,
		       u.role, u.status, u.email_verified, u.created_at,
		       COUNT(p.id) as projects_count
		FROM users u
		LEFT JOIN projects p ON p.owner_id = u.id
		GROUP BY u.id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	usersList := make([]AdminUserRow, 0)
	for rows.Next() {
		var u AdminUserRow
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Bio,
			&u.Role, &u.Status, &u.EmailVerified, &u.CreatedAt, &u.ProjectsCount,
		); err == nil {
			usersList = append(usersList, u)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"users": usersList})
}

// UpdateUserRole handles PATCH /api/admin/users/{id}/role
func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	targetUserID := chi.URLParam(r, "id")

	var req UpdateRolePayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Role != "USER" && req.Role != "ADMIN") {
		response.Error(w, http.StatusBadRequest, "Valid role (USER, ADMIN) is required")
		return
	}

	if targetUserID == adminUser.ID && req.Role != "ADMIN" {
		response.Error(w, http.StatusBadRequest, "You cannot demote your own administrator account")
		return
	}

	// Prevent demoting the last active administrator
	if req.Role != "ADMIN" {
		var currentRole string
		err := h.db.QueryRow(r.Context(), "SELECT role FROM users WHERE id = $1", targetUserID).Scan(&currentRole)
		if err == nil && currentRole == "ADMIN" {
			var adminCount int
			err := h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM users WHERE role = 'ADMIN' AND status = 'ACTIVE'").Scan(&adminCount)
			if err == nil && adminCount <= 1 {
				response.Error(w, http.StatusBadRequest, "Cannot demote the last remaining administrator on the system")
				return
			}
		}
	}

	_, err := h.db.Exec(r.Context(), "UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2", req.Role, targetUserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_CHANGED_ROLE', 'user', $2, $3)
	`, adminUser.ID, targetUserID, map[string]string{"new_role": req.Role})

	response.JSON(w, http.StatusOK, map[string]string{"message": "User role updated"})
}

// UpdateUserStatus handles PATCH /api/admin/users/{id}/status
func (h *Handler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	targetUserID := chi.URLParam(r, "id")

	var req UpdateStatusPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Status != "ACTIVE" && req.Status != "SUSPENDED") {
		response.Error(w, http.StatusBadRequest, "Valid status (ACTIVE, SUSPENDED) is required")
		return
	}

	if targetUserID == adminUser.ID {
		response.Error(w, http.StatusBadRequest, "You cannot suspend your own administrator account")
		return
	}

	// Prevent suspending the last active administrator
	if req.Status == "SUSPENDED" {
		var currentRole string
		err := h.db.QueryRow(r.Context(), "SELECT role FROM users WHERE id = $1", targetUserID).Scan(&currentRole)
		if err == nil && currentRole == "ADMIN" {
			var adminCount int
			err := h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM users WHERE role = 'ADMIN' AND status = 'ACTIVE'").Scan(&adminCount)
			if err == nil && adminCount <= 1 {
				response.Error(w, http.StatusBadRequest, "Cannot suspend the last active administrator on the system")
				return
			}
		}
	}

	_, err := h.db.Exec(r.Context(), "UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2", req.Status, targetUserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// If suspended, invalidate all active sessions
	if req.Status == "SUSPENDED" {
		_, _ = h.db.Exec(r.Context(), "DELETE FROM sessions WHERE user_id = $1", targetUserID)
	}

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_MODERATED_USER', 'user', $2, $3)
	`, adminUser.ID, targetUserID, map[string]string{"new_status": req.Status})

	response.JSON(w, http.StatusOK, map[string]string{"message": "User status updated"})
}

// GetSystemHealth handles GET /api/admin/system
func (h *Handler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbPingErr := h.db.Ping(r.Context())
	dbStatus := "healthy"
	if dbPingErr != nil {
		dbStatus = "unreachable"
	}

	stat := h.db.Stat()

	response.JSON(w, http.StatusOK, map[string]any{
		"application": map[string]any{
			"status":            "healthy",
			"uptime_seconds":    int(time.Since(startTime).Seconds()),
			"goroutines":        runtime.NumGoroutine(),
			"memory_alloc_mb":   m.Alloc / 1024 / 1024,
			"memory_sys_mb":     m.Sys / 1024 / 1024,
			"garbage_collector": m.NumGC,
		},
		"database": map[string]any{
			"status":         dbStatus,
			"total_conns":    stat.TotalConns(),
			"idle_conns":     stat.IdleConns(),
			"acquired_conns": stat.AcquiredConns(),
			"max_conns":      stat.MaxConns(),
		},
		"host_specs": system.GetHostSpecs(),
	})
}

// ListAuditLogs handles GET /api/admin/audit
func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT a.id, a.actor_id, a.action, a.target_type, a.target_id, a.metadata, a.created_at,
		       u.display_name, u.username
		FROM audit_logs a
		LEFT JOIN users u ON u.id = a.actor_id
		ORDER BY a.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	list := make([]AuditEntry, 0)
	for rows.Next() {
		var entry AuditEntry
		var metaJSON []byte
		if err := rows.Scan(
			&entry.ID, &entry.ActorID, &entry.Action, &entry.TargetType, &entry.TargetID,
			&metaJSON, &entry.CreatedAt, &entry.ActorName, &entry.ActorUser,
		); err == nil {
			_ = json.Unmarshal(metaJSON, &entry.Metadata)
			list = append(list, entry)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{"audit_logs": list})
}
