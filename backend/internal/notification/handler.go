package notification

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/response"
)

// Handler handles user notification endpoints.
type Handler struct {
	db *pgxpool.Pool
}

// NewHandler creates a new notification handler.
func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// ListMyNotifications handles GET /api/me/notifications
func (h *Handler) ListMyNotifications(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT id, user_id, type, title, body, data, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`, user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	list := make([]Notification, 0)
	unreadCount := 0
	for rows.Next() {
		var n Notification
		var dataJSON []byte
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &dataJSON, &n.ReadAt, &n.CreatedAt,
		); err == nil {
			_ = json.Unmarshal(dataJSON, &n.Data)
			if n.ReadAt == nil {
				unreadCount++
			}
			list = append(list, n)
		}
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"notifications": list,
		"unread_count":  unreadCount,
	})
}

// MarkRead handles PATCH /api/me/notifications/{id}/read
func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notifID := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `
		UPDATE notifications
		SET read_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, notifID, user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

// MarkAllRead handles POST /api/me/notifications/read-all
func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	_, err := h.db.Exec(r.Context(), `
		UPDATE notifications
		SET read_at = NOW()
		WHERE user_id = $1 AND read_at IS NULL
	`, user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}
