package report

import (
	"time"

	"ngumpul-host/backend/internal/auth"
)

// Report represents a moderation flag against a project or a comment.
type Report struct {
	ID            string           `json:"id"`
	ReporterID    string           `json:"reporter_id"`
	TargetType    string           `json:"target_type"` // 'PROJECT' or 'COMMENT'
	TargetID      string           `json:"target_id"`
	ProjectID     *string          `json:"project_id,omitempty"`
	ProjectName   string           `json:"project_name,omitempty"`
	ProjectSlug   string           `json:"project_slug,omitempty"`
	Reason        string           `json:"reason"` // SPAM, ABUSE_HARASSMENT, INAPPROPRIATE, MALICIOUS_SUSPICIOUS, OTHER
	Details       string           `json:"details"`
	Status        string           `json:"status"` // OPEN, REVIEWED, DISMISSED, RESOLVED
	ResolvedBy    *string          `json:"resolved_by,omitempty"`
	ResolvedAt    *time.Time       `json:"resolved_at,omitempty"`
	ActionNotes   string           `json:"action_notes"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	TargetSnippet string           `json:"target_snippet,omitempty"`
	Reporter      *auth.PublicUser `json:"reporter,omitempty"`
}

// CreateReportPayload defines the request body for submitting a report.
type CreateReportPayload struct {
	Reason  string `json:"reason"`
	Details string `json:"details"`
}

// ResolveReportPayload defines the action taken by an admin.
type ResolveReportPayload struct {
	Action string `json:"action"` // "DELETE_COMMENT" | "RESOLVE"
	Notes  string `json:"notes"`
}

// DismissReportPayload defines the admin dismissal note.
type DismissReportPayload struct {
	Notes string `json:"notes"`
}
