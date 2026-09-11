package activity

import (
	"time"
)

// Activity represents an event or audit record visible in the activity feed.
type Activity struct {
	ID         string         `json:"id"`
	ActorID    *string        `json:"actor_id"`
	ProjectID  *string        `json:"project_id"`
	Type       string         `json:"type"`
	Metadata   map[string]any `json:"metadata"`
	Visibility string         `json:"visibility"`
	CreatedAt  time.Time      `json:"created_at"`
	ActorName  *string        `json:"actor_name"`
	ActorUser  *string        `json:"actor_username"`
	ProjName   *string        `json:"project_name"`
	ProjSlug   *string        `json:"project_slug"`
}
