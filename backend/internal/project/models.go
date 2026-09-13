package project

import (
	"regexp"
	"strings"
	"time"

	"ngumpul-host/backend/internal/auth"
)

// Project represents an application or experiment hosted on Ngumpul Host.
type Project struct {
	ID               string           `json:"id"`
	OwnerID          string           `json:"owner_id"`
	Name             string           `json:"name"`
	Slug             string           `json:"slug"`
	Description      string           `json:"description"`
	CoverImageURL    string           `json:"cover_image_url"`
	RepositoryURL    string           `json:"repository_url"`
	DocumentationURL string           `json:"documentation_url"`
	DemoURL          string           `json:"demo_url"`
	Readme           string           `json:"readme"`
	TechnologyStack  []string         `json:"technology_stack"`
	HostingType      string           `json:"hosting_type"`
	PublicURL          string           `json:"public_url"`
	Status             string           `json:"status"`
	LifecycleStatus    string           `json:"lifecycle_status"`
	Availability       string           `json:"availability"`
	AvailabilityReason string           `json:"availability_reason"`
	Visibility         string           `json:"visibility"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	PublishedAt        *time.Time       `json:"published_at"`
	Owner              *auth.PublicUser `json:"owner,omitempty"`
}

// UpdateMetadataRequest payload for updating personal project metadata.
type UpdateMetadataRequest struct {
	Description      *string  `json:"description"`
	Readme           *string  `json:"readme"`
	CoverImageURL    *string  `json:"cover_image_url"`
	RepositoryURL    *string  `json:"repository_url"`
	DocumentationURL *string  `json:"documentation_url"`
	DemoURL          *string  `json:"demo_url"`
	TechnologyStack  []string `json:"technology_stack"`
}

// AdminCreateProjectRequest payload for creating a project by administrator.
type AdminCreateProjectRequest struct {
	OwnerID          string   `json:"owner_id"`
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	Description      string   `json:"description"`
	Readme           string   `json:"readme"`
	CoverImageURL    string   `json:"cover_image_url"`
	RepositoryURL    string   `json:"repository_url"`
	DocumentationURL string   `json:"documentation_url"`
	DemoURL          string   `json:"demo_url"`
	TechnologyStack  []string `json:"technology_stack"`
	HostingType      string   `json:"hosting_type"`
	PublicURL        string   `json:"public_url"`
	Status           string   `json:"status"`
	LifecycleStatus  string   `json:"lifecycle_status"`
	Availability     string   `json:"availability"`
	Visibility       string   `json:"visibility"`
}

// AdminUpdateProjectRequest payload for updating a project by administrator.
type AdminUpdateProjectRequest struct {
	Name               *string  `json:"name"`
	Slug               *string  `json:"slug"`
	Description        *string  `json:"description"`
	Readme             *string  `json:"readme"`
	CoverImageURL      *string  `json:"cover_image_url"`
	RepositoryURL      *string  `json:"repository_url"`
	DocumentationURL   *string  `json:"documentation_url"`
	DemoURL            *string  `json:"demo_url"`
	TechnologyStack    []string `json:"technology_stack"`
	HostingType        *string  `json:"hosting_type"`
	PublicURL          *string  `json:"public_url"`
	Status             *string  `json:"status"`
	LifecycleStatus    *string  `json:"lifecycle_status"`
	Availability       *string  `json:"availability"`
	AvailabilityReason *string  `json:"availability_reason"`
	Visibility         *string  `json:"visibility"`
	OwnerID            *string  `json:"owner_id"`
}

var (
	slugRegex        = regexp.MustCompile(`[^a-z0-9]+`)
	multiHyphenRegex = regexp.MustCompile(`-+`)
)

// GenerateSlug normalizes a string into a URL-safe lowercase slug.
func GenerateSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = slugRegex.ReplaceAllString(s, "-")
	s = multiHyphenRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "project"
	}
	return s
}
