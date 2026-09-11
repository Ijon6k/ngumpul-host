package hosting

import (
	"time"

	"ngumpul-host/backend/internal/auth"
)

// Request represents a user's submission requesting project hosting on Ngumpul Host.
type Request struct {
	ID               string           `json:"id"`
	RequesterID      string           `json:"requester_id"`
	ProjectName      string           `json:"project_name"`
	Description      string           `json:"description"`
	RepositoryURL    string           `json:"repository_url"`
	DocumentationURL string           `json:"documentation_url"`
	DeploymentNotes  string           `json:"deployment_notes"`
	TechnologyStack  []string         `json:"technology_stack"`
	Status           string           `json:"status"`
	AdminNotes       string           `json:"admin_notes"`
	ReviewedBy       *string          `json:"reviewed_by"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	ReviewedAt       *time.Time       `json:"reviewed_at"`
	Requester        *auth.PublicUser `json:"requester,omitempty"`
}

// SubmitRequestPayload payload for submitting a new hosting request.
type SubmitRequestPayload struct {
	ProjectName      string   `json:"project_name"`
	Description      string   `json:"description"`
	RepositoryURL    string   `json:"repository_url"`
	DocumentationURL string   `json:"documentation_url"`
	DeploymentNotes  string   `json:"deployment_notes"`
	TechnologyStack  []string `json:"technology_stack"`
}

// ReviewActionPayload payload for admin approval, rejection, or completion.
type ReviewActionPayload struct {
	AdminNotes string `json:"admin_notes"`
	PublicURL  string `json:"public_url"`
}
