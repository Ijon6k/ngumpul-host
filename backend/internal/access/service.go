package access

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/config"
)

var (
	ErrInvalidToken       = errors.New("invalid or unrecognised invitation token")
	ErrTokenRevoked       = errors.New("invitation token has been revoked")
	ErrTokenExpired       = errors.New("invitation token has expired")
	ErrTokenLimitExceeded = errors.New("invitation token usage limit reached")
	ErrEmailMismatch      = errors.New("invitation token is restricted to a specific email address")
	ErrInvalidMode        = errors.New("invalid registration mode")
)

type Service struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewService(db *pgxpool.Pool, cfg *config.Config) *Service {
	return &Service{
		db:  db,
		cfg: cfg,
	}
}

// HashToken calculates the SHA-256 hash of a raw invitation token.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(h[:])
}

// GenerateToken generates a cryptographically secure 32-byte hex-encoded random string.
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// GetRegistrationMode retrieves the current registration mode from instance_settings.
func (s *Service) GetRegistrationMode(ctx context.Context) (string, error) {
	var mode string
	err := s.db.QueryRow(ctx, `
		SELECT value FROM instance_settings WHERE key = 'registration_mode'
	`).Scan(&mode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RegistrationModeInviteOnly, nil
		}
		return "", fmt.Errorf("failed to read registration mode: %w", err)
	}

	mode = strings.ToUpper(strings.TrimSpace(mode))
	switch mode {
	case RegistrationModeOpen, RegistrationModeInviteOnly, RegistrationModeClosed:
		return mode, nil
	default:
		return RegistrationModeInviteOnly, nil
	}
}

// SetRegistrationMode updates the registration mode and records an audit log.
func (s *Service) SetRegistrationMode(ctx context.Context, mode string, actorID string) error {
	mode = strings.ToUpper(strings.TrimSpace(mode))
	if mode != RegistrationModeOpen && mode != RegistrationModeInviteOnly && mode != RegistrationModeClosed {
		return ErrInvalidMode
	}

	var oldMode string
	_ = s.db.QueryRow(ctx, "SELECT value FROM instance_settings WHERE key = 'registration_mode'").Scan(&oldMode)

	_, err := s.db.Exec(ctx, `
		INSERT INTO instance_settings (key, value, updated_at)
		VALUES ('registration_mode', $1, NOW())
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, mode)
	if err != nil {
		return fmt.Errorf("failed to update registration mode: %w", err)
	}

	_, _ = s.db.Exec(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'UPDATE_REGISTRATION_MODE', 'SETTING', 'registration_mode', $2)
	`, actorID, map[string]string{
		"old_mode": oldMode,
		"new_mode": mode,
	})

	return nil
}

// CreateInvitation generates a new invitation token and saves its hash.
func (s *Service) CreateInvitation(ctx context.Context, actorID string, email *string, maxUses int, expiresInDays int) (string, *Invitation, error) {
	if maxUses <= 0 {
		maxUses = 1
	}

	rawToken, err := GenerateToken()
	if err != nil {
		return "", nil, err
	}
	tokenHash := HashToken(rawToken)

	var expiresAt *time.Time
	if expiresInDays > 0 {
		exp := time.Now().AddDate(0, 0, expiresInDays)
		expiresAt = &exp
	}

	var cleanedEmail *string
	if email != nil && strings.TrimSpace(*email) != "" {
		em := strings.ToLower(strings.TrimSpace(*email))
		cleanedEmail = &em
	}

	var inv Invitation
	err = s.db.QueryRow(ctx, `
		INSERT INTO invitations (token_hash, token, created_by, invited_email, max_uses, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, token, created_by, invited_email, max_uses, used_count, expires_at, revoked_at, created_at, updated_at
	`, tokenHash, rawToken, actorID, cleanedEmail, maxUses, expiresAt).Scan(
		&inv.ID, &inv.Token, &inv.CreatedBy, &inv.InvitedEmail, &inv.MaxUses, &inv.UsedCount,
		&inv.ExpiresAt, &inv.RevokedAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	inv.IsActive = inv.RevokedAt == nil && inv.UsedCount < inv.MaxUses && (inv.ExpiresAt == nil || time.Now().Before(*inv.ExpiresAt))

	_, _ = s.db.Exec(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'CREATE_INVITATION', 'INVITATION', $2, $3)
	`, actorID, inv.ID, map[string]any{
		"max_uses":      inv.MaxUses,
		"invited_email": inv.InvitedEmail,
	})

	return rawToken, &inv, nil
}

// ValidateInvitation checks if a raw token is valid without consuming it.
func (s *Service) ValidateInvitation(ctx context.Context, rawToken string) (*ValidateInvitationResponse, error) {
	mode, err := s.GetRegistrationMode(ctx)
	if err == nil && mode == RegistrationModeClosed {
		return &ValidateInvitationResponse{
			Valid:   false,
			Message: "Registration is currently closed on this host. Invitations cannot be redeemed while registration is closed.",
		}, nil
	}

	tokenHash := HashToken(rawToken)

	var inv Invitation
	err = s.db.QueryRow(ctx, `
		SELECT id, invited_email, max_uses, used_count, expires_at, revoked_at
		FROM invitations
		WHERE token_hash = $1
	`, tokenHash).Scan(
		&inv.ID, &inv.InvitedEmail, &inv.MaxUses, &inv.UsedCount, &inv.ExpiresAt, &inv.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &ValidateInvitationResponse{Valid: false, Message: "Invalid invitation token"}, nil
		}
		return nil, fmt.Errorf("failed to validate invitation: %w", err)
	}

	if inv.RevokedAt != nil {
		return &ValidateInvitationResponse{Valid: false, Message: "Invitation token has been revoked"}, nil
	}

	if inv.ExpiresAt != nil && time.Now().After(*inv.ExpiresAt) {
		return &ValidateInvitationResponse{Valid: false, Message: "Invitation token has expired"}, nil
	}

	if inv.UsedCount >= inv.MaxUses {
		return &ValidateInvitationResponse{Valid: false, Message: "Invitation token usage limit reached"}, nil
	}

	return &ValidateInvitationResponse{
		Valid:        true,
		InvitedEmail: inv.InvitedEmail,
		ExpiresAt:    inv.ExpiresAt,
	}, nil
}

// ConsumeInvitationTx locks the invitation record and increments used_count within an ongoing transaction.
func (s *Service) ConsumeInvitationTx(ctx context.Context, tx pgx.Tx, rawToken string, applicantEmail string) error {
	mode, err := s.GetRegistrationMode(ctx)
	if err == nil && mode == RegistrationModeClosed {
		return errors.New("registration is currently closed on this host")
	}

	tokenHash := HashToken(rawToken)

	var inv Invitation
	err = tx.QueryRow(ctx, `
		SELECT id, invited_email, max_uses, used_count, expires_at, revoked_at
		FROM invitations
		WHERE token_hash = $1
		FOR UPDATE
	`, tokenHash).Scan(
		&inv.ID, &inv.InvitedEmail, &inv.MaxUses, &inv.UsedCount, &inv.ExpiresAt, &inv.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidToken
		}
		return fmt.Errorf("failed to query invitation in transaction: %w", err)
	}

	if inv.RevokedAt != nil {
		return ErrTokenRevoked
	}

	if inv.ExpiresAt != nil && time.Now().After(*inv.ExpiresAt) {
		return ErrTokenExpired
	}

	if inv.UsedCount >= inv.MaxUses {
		return ErrTokenLimitExceeded
	}

	if inv.InvitedEmail != nil && *inv.InvitedEmail != "" {
		if strings.ToLower(*inv.InvitedEmail) != strings.ToLower(applicantEmail) {
			return ErrEmailMismatch
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE invitations
		SET used_count = used_count + 1, updated_at = NOW()
		WHERE id = $1
	`, inv.ID)
	if err != nil {
		return fmt.Errorf("failed to increment invitation used count: %w", err)
	}

	return nil
}

// ListInvitations lists all invitations for administrative inspection.
func (s *Service) ListInvitations(ctx context.Context) ([]Invitation, error) {
	rows, err := s.db.Query(ctx, `
		SELECT i.id, i.token, i.created_by, u.display_name, i.invited_email,
		       i.max_uses, i.used_count, i.expires_at, i.revoked_at,
		       i.created_at, i.updated_at
		FROM invitations i
		LEFT JOIN users u ON u.id = i.created_by
		ORDER BY i.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list invitations: %w", err)
	}
	defer rows.Close()

	list := make([]Invitation, 0)
	now := time.Now()
	for rows.Next() {
		var inv Invitation
		if err := rows.Scan(
			&inv.ID, &inv.Token, &inv.CreatedBy, &inv.CreatorName, &inv.InvitedEmail,
			&inv.MaxUses, &inv.UsedCount, &inv.ExpiresAt, &inv.RevokedAt,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		inv.IsActive = inv.RevokedAt == nil && inv.UsedCount < inv.MaxUses && (inv.ExpiresAt == nil || now.Before(*inv.ExpiresAt))
		list = append(list, inv)
	}

	return list, nil
}

// RevokeInvitation invalidates an invitation by setting revoked_at.
func (s *Service) RevokeInvitation(ctx context.Context, id string, actorID string) error {
	res, err := s.db.Exec(ctx, `
		UPDATE invitations
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`, id)
	if err != nil {
		return fmt.Errorf("failed to revoke invitation: %w", err)
	}

	if res.RowsAffected() == 0 {
		return errors.New("invitation not found or already revoked")
	}

	_, _ = s.db.Exec(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'REVOKE_INVITATION', 'INVITATION', $2, '{}'::jsonb)
	`, actorID, id)

	return nil
}
