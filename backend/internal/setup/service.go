package setup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/config"
	"ngumpul-host/backend/internal/instance"
)

var (
	ErrAlreadyInitialized = errors.New("node setup has already been completed")
	ErrInvalidInput       = errors.New("invalid setup parameters")
)

var validUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)

type SetupRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

type Service struct {
	db  *pgxpool.Pool
	sm  *auth.SessionManager
	cfg *config.Config
}

func NewService(db *pgxpool.Pool, sm *auth.SessionManager, cfg *config.Config) *Service {
	return &Service{
		db:  db,
		sm:  sm,
		cfg: cfg,
	}
}

// IsInitialized checks if at least one administrator account already exists.
func (s *Service) IsInitialized(ctx context.Context) (bool, error) {
	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'ADMIN'").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check admin existence: %w", err)
	}
	return count > 0, nil
}

// GetDomain retrieves the stored node domain or falls back to AppURL hostname.
func (s *Service) GetDomain(ctx context.Context) (string, error) {
	domain, err := instance.GetDomain(ctx, s.db)
	if err != nil {
		return "", fmt.Errorf("failed to load domain setting: %w", err)
	}
	return domain, nil
}

// GetAppURL retrieves the canonical stored app URL (app_url key, or scheme + domain).
func (s *Service) GetAppURL(ctx context.Context) (string, error) {
	appURL, err := instance.GetAppURL(ctx, s.db)
	if err != nil {
		return "", fmt.Errorf("failed to load app URL setting: %w", err)
	}
	return appURL, nil
}

// ExecuteSetup performs initial node provisioning by creating the first admin user and recording the domain.
func (s *Service) ExecuteSetup(ctx context.Context, req SetupRequest) (*auth.User, *auth.Session, error) {
	// 1. Validation
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 128 {
		return nil, nil, fmt.Errorf("%w: operator name must be between 1 and 128 characters", ErrInvalidInput)
	}

	username := strings.ToLower(strings.TrimSpace(req.Username))
	if !validUsernameRegex.MatchString(username) {
		return nil, nil, fmt.Errorf("%w: username must be 3-30 characters and contain only alphanumeric, underscore, or hyphens", ErrInvalidInput)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, nil, fmt.Errorf("%w: invalid email address", ErrInvalidInput)
	}

	if len(req.Password) < 8 {
		return nil, nil, fmt.Errorf("%w: password must be at least 8 characters long", ErrInvalidInput)
	}

	domain := strings.TrimSpace(req.Domain)
	// Clean protocol if provided (e.g. https://domain.com -> domain.com)
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimRight(domain, "/")
	if domain == "" {
		return nil, nil, fmt.Errorf("%w: domain cannot be empty", ErrInvalidInput)
	}

	// 2. Hash Password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash operator password: %w", err)
	}

	// 3. Database Transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Guard against race condition: check if any admin exists
	var adminCount int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'ADMIN'").Scan(&adminCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify setup state: %w", err)
	}
	if adminCount > 0 {
		return nil, nil, ErrAlreadyInitialized
	}

	var user auth.User
	err = tx.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name, role, status, email_verified, bio)
		VALUES ($1, $2, $3, $4, 'ADMIN', 'ACTIVE', true, 'Primary node operator and community organizer.')
		RETURNING id, username, email, display_name, avatar_url, bio, role, status, email_verified, created_at
	`, username, email, passwordHash, name).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.Bio, &user.Role, &user.Status,
		&user.EmailVerified, &user.CreatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to insert operator account: %w", err)
	}

	// Save domain into instance_settings
	_, err = tx.Exec(ctx, `
		INSERT INTO instance_settings (key, value, updated_at)
		VALUES ('domain', $1, NOW())
		ON CONFLICT (key) DO UPDATE SET value = $1, updated_at = NOW()
	`, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save domain setting: %w", err)
	}

	// Save app_url into instance_settings for system referencing
	appURL := fmt.Sprintf("http://%s", domain)
	if strings.Contains(domain, "localhost") || strings.Contains(domain, "127.0.0.1") {
		appURL = fmt.Sprintf("http://%s", domain)
	} else if !strings.HasPrefix(domain, "http") {
		appURL = fmt.Sprintf("https://%s", domain)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO instance_settings (key, value, updated_at)
		VALUES ('app_url', $1, NOW())
		ON CONFLICT (key) DO UPDATE SET value = $1, updated_at = NOW()
	`, appURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save app_url setting: %w", err)
	}

	// Record audit log
	metaJSON, _ := json.Marshal(map[string]any{
		"operator": user.Username,
		"domain":   domain,
	})
	_, err = tx.Exec(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata)
		VALUES ($1, 'INITIAL_NODE_SETUP', 'SYSTEM', 'node', $2)
	`, user.ID, metaJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to record setup audit log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to commit setup transaction: %w", err)
	}

	// 4. Create active session for automatic login
	session, err := s.sm.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create operator session: %w", err)
	}

	return &user, session, nil
}
