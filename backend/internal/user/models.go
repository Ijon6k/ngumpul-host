package user

import (
	"ngumpul-host/backend/internal/auth"
)

// UserWithProjects represents a public member profile augmented with total projects count.
type UserWithProjects struct {
	auth.PublicUser
	ProjectsCount int `json:"projects_count"`
}

// ProfileResponse represents the full public profile of a user along with their public projects.
type ProfileResponse struct {
	auth.PublicUser
	Projects []any `json:"projects"`
}
