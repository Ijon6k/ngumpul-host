package comment

import (
	"time"

	"ngumpul-host/backend/internal/auth"
)

// Comment represents a public user comment on a published project.
type Comment struct {
	ID        string           `json:"id"`
	ProjectID string           `json:"project_id"`
	AuthorID  string           `json:"author_id"`
	Content   string           `json:"content"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt *time.Time       `json:"deleted_at,omitempty"`
	IsDeleted bool             `json:"is_deleted"`
	Author    *auth.PublicUser `json:"author,omitempty"`
}

// CreateCommentRequest is the JSON payload for posting a new comment.
type CreateCommentRequest struct {
	Content string `json:"content"`
}

// AdminCommentRow provides rich moderation context for administrators.
type AdminCommentRow struct {
	ID          string           `json:"id"`
	ProjectID   string           `json:"project_id"`
	ProjectName string           `json:"project_name"`
	ProjectSlug string           `json:"project_slug"`
	AuthorID    string           `json:"author_id"`
	Content     string           `json:"content"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
	IsDeleted   bool             `json:"is_deleted"`
	ReportCount int              `json:"report_count"`
	Author      *auth.PublicUser `json:"author,omitempty"`
}
