package admin

import (
	"time"

	"ngumpul-host/backend/internal/auth"
)

// AdminUserRow represents a user record with project counts for moderation.
type AdminUserRow struct {
	auth.User
	ProjectsCount int `json:"projects_count"`
}

// UpdateRolePayload payload for changing a user's role (USER or ADMIN).
type UpdateRolePayload struct {
	Role string `json:"role"`
}

// UpdateStatusPayload payload for changing a user's status (ACTIVE or SUSPENDED).
type UpdateStatusPayload struct {
	Status string `json:"status"`
}

// AuditEntry represents an audit log entry for administrative traceability.
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
