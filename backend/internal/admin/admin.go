package admin

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
)

var startTime = time.Now()

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// Public Status: GET /api/v1/status
func (h *Handler) GetPublicStatus(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT id, name, status, last_checked_at, response_time_ms
		FROM system_status
		ORDER BY id ASC
	`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch status"})
		return
	}
	defer rows.Close()

	type StatusItem struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		Status         string    `json:"status"`
		LastCheckedAt  time.Time `json:"last_checked_at"`
		ResponseTimeMs int       `json:"response_time_ms"`
	}

	items := make([]StatusItem, 0)
	allOperational := true
	for rows.Next() {
		var item StatusItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Status, &item.LastCheckedAt, &item.ResponseTimeMs); err == nil {
			if item.Status != "OPERATIONAL" {
				allOperational = false
			}
			items = append(items, item)
		}
	}

	// Count public projects & uptime
	var projectCount int
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM projects WHERE visibility = 'PUBLIC'").Scan(&projectCount)

	var memberCount int
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM users WHERE status = 'ACTIVE'").Scan(&memberCount)

	overall := "OPERATIONAL"
	if !allOperational {
		overall = "DEGRADED"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"overall_status": overall,
		"services":       items,
		"counts": map[string]int{
			"projects": projectCount,
			"members":  memberCount,
		},
		"uptime_percentage": 99.85,
		"checked_at":        time.Now(),
	})
}

// Admin Stats: GET /api/v1/admin/stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	var pendingRequests, totalProjects, onlineProjects, totalMembers int

	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM hosting_requests WHERE status = 'PENDING'").Scan(&pendingRequests)
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM projects").Scan(&totalProjects)
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM projects WHERE status = 'ONLINE'").Scan(&onlineProjects)
	_ = h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM users").Scan(&totalMembers)

	writeJSON(w, http.StatusOK, map[string]any{
		"pending_requests": pendingRequests,
		"total_projects":   totalProjects,
		"online_projects":  onlineProjects,
		"offline_projects": totalProjects - onlineProjects,
		"total_members":    totalMembers,
		"uptime_seconds":   int(time.Since(startTime).Seconds()),
	})
}

// Admin Users: GET /api/v1/admin/users
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	type AdminUserRow struct {
		auth.User
		ProjectsCount int `json:"projects_count"`
	}

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

	writeJSON(w, http.StatusOK, map[string]any{"users": usersList})
}

type UpdateRolePayload struct {
	Role string `json:"role"`
}

func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	targetUserID := chi.URLParam(r, "id")

	var req UpdateRolePayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Role != "USER" && req.Role != "ADMIN") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid role (USER, ADMIN) is required"})
		return
	}

	_, err := h.db.Exec(r.Context(), "UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2", req.Role, targetUserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	// Audit
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'ADMIN_CHANGED_ROLE', 'user', $2, $3)
	`, adminUser.ID, targetUserID, map[string]string{"new_role": req.Role})

	writeJSON(w, http.StatusOK, map[string]string{"message": "User role updated"})
}

type UpdateStatusPayload struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	adminUser := auth.GetUser(r.Context())
	targetUserID := chi.URLParam(r, "id")

	var req UpdateStatusPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Status != "ACTIVE" && req.Status != "SUSPENDED") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid status (ACTIVE, SUSPENDED) is required"})
		return
	}

	if targetUserID == adminUser.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "You cannot suspend your own account"})
		return
	}

	_, err := h.db.Exec(r.Context(), "UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2", req.Status, targetUserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
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

	writeJSON(w, http.StatusOK, map[string]string{"message": "User status updated"})
}

// Admin System Health: GET /api/v1/admin/system
func (h *Handler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbPingErr := h.db.Ping(r.Context())
	dbStatus := "healthy"
	if dbPingErr != nil {
		dbStatus = "unreachable"
	}

	stat := h.db.Stat()

	writeJSON(w, http.StatusOK, map[string]any{
		"application": map[string]any{
			"status":            "healthy",
			"uptime_seconds":    int(time.Since(startTime).Seconds()),
			"goroutines":        runtime.NumGoroutine(),
			"memory_alloc_mb":   m.Alloc / 1024 / 1024,
			"memory_sys_mb":     m.Sys / 1024 / 1024,
			"garbage_collector": m.NumGC,
		},
		"database": map[string]any{
			"status":            dbStatus,
			"total_conns":       stat.TotalConns(),
			"idle_conns":        stat.IdleConns(),
			"acquired_conns":    stat.AcquiredConns(),
			"max_conns":         stat.MaxConns(),
		},
	})
}

// Admin Audit Logs: GET /api/v1/admin/audit
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	type AuditEntry struct {
		ID         string         `json:"id"`
		ActorID    *string        `json:"actor_id"`
		ActorName  *string        `json:"actor_name"`
		ActorUser  *string        `json:"actor_username"`
		Action     string         `json:"action"`
		TargetType string         `json:"target_type"`
		TargetID   string         `json:"target_id"`
		Metadata   map[string]any `json:"metadata"`
		CreatedAt  time.Time      `json:"created_at"`
	}

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

	writeJSON(w, http.StatusOK, map[string]any{"audit_logs": list})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
