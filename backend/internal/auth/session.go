package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionManager struct {
	db *pgxpool.Pool
}

func NewSessionManager(db *pgxpool.Pool) *SessionManager {
	return &SessionManager{db: db}
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID string) (*Session, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	sessionID := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days

	_, err := sm.db.Exec(ctx, `
		INSERT INTO sessions (id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
	`, sessionID, userID, expiresAt)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

func (sm *SessionManager) GetUserBySession(ctx context.Context, sessionID string) (*User, error) {
	if sessionID == "" {
		return nil, errors.New("empty session token")
	}

	var u User
	var expiresAt time.Time
	err := sm.db.QueryRow(ctx, `
		SELECT u.id, u.username, u.email, u.display_name, u.avatar_url, u.bio, u.role, u.status, u.email_verified, u.created_at, s.expires_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = $1
	`, sessionID).Scan(
		&u.ID, &u.Username, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.Role, &u.Status, &u.EmailVerified, &u.CreatedAt, &expiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if time.Now().After(expiresAt) {
		// Session expired, delete it
		_ = sm.DeleteSession(ctx, sessionID)
		return nil, errors.New("session expired")
	}

	return &u, nil
}

func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := sm.db.Exec(ctx, "DELETE FROM sessions WHERE id = $1", sessionID)
	return err
}
