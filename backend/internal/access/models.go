package access

import "time"

const (
	RegistrationModeOpen       = "OPEN"
	RegistrationModeInviteOnly = "INVITE_ONLY"
	RegistrationModeClosed     = "CLOSED"
)

// Invitation represents an access grant token record.
type Invitation struct {
	ID           string     `json:"id"`
	Token        *string    `json:"token,omitempty"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreatorName  *string    `json:"creator_name,omitempty"`
	InvitedEmail *string    `json:"invited_email,omitempty"`
	MaxUses      int        `json:"max_uses"`
	UsedCount    int        `json:"used_count"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	IsActive     bool       `json:"is_active"`
}

type CreateInvitationRequest struct {
	InvitedEmail  *string `json:"invited_email"`
	MaxUses       int     `json:"max_uses"`
	ExpiresInDays int     `json:"expires_in_days"`
}

type CreateInvitationResponse struct {
	Invitation Invitation `json:"invitation"`
	RawToken   string     `json:"raw_token"`
	InviteURL  string     `json:"invite_url"`
}

type ValidateInvitationResponse struct {
	Valid        bool       `json:"valid"`
	InvitedEmail *string    `json:"invited_email,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	Message      string     `json:"message,omitempty"`
}

type UpdateSettingsRequest struct {
	RegistrationMode string `json:"registration_mode"`
}

type InstanceSettingsResponse struct {
	RegistrationMode string    `json:"registration_mode"`
	UpdatedAt        time.Time `json:"updated_at"`
}
