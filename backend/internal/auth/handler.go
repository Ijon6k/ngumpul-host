package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/config"
)

type Handler struct {
	db  *pgxpool.Pool
	sm  *SessionManager
	cfg *config.Config
}

func NewHandler(db *pgxpool.Pool, sm *SessionManager, cfg *config.Config) *Handler {
	return &Handler{
		db:  db,
		sm:  sm,
		cfg: cfg,
	}
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.cfg.SessionSecure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cfg.SessionSecure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	req.Username = strings.TrimSpace(strings.ToLower(req.Username))
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.Username == "" || req.Email == "" || len(req.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username, valid email, and password (minimum 8 characters) are required"})
		return
	}

	if req.DisplayName == "" {
		req.DisplayName = req.Username
	}

	passHash, err := HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process credentials"})
		return
	}

	var userID string
	err = h.db.QueryRow(r.Context(), `
		INSERT INTO users (username, email, password_hash, display_name, role, status)
		VALUES ($1, $2, $3, $4, 'USER', 'ACTIVE')
		RETURNING id
	`, req.Username, req.Email, passHash, req.DisplayName).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Username already taken"})
			return
		}
		if strings.Contains(err.Error(), "users_email_key") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Email already registered"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Could not register account"})
		return
	}

	// Create activity event: MEMBER_JOINED
	_, _ = h.db.Exec(r.Context(), `
		INSERT INTO activities (actor_id, type, metadata, visibility)
		VALUES ($1, 'MEMBER_JOINED', $2, 'PUBLIC')
	`, userID, map[string]string{
		"title":   req.DisplayName + " joined Ngumpul",
		"message": "Welcome " + req.DisplayName + " to the community",
	})

	// Create session
	session, err := h.sm.CreateSession(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Account created, but session could not be initialized"})
		return
	}

	h.setSessionCookie(w, session.ID, session.ExpiresAt)

	user, _ := h.sm.GetUserBySession(r.Context(), session.ID)
	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Account created successfully",
		"user":    user,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	req.EmailOrUsername = strings.TrimSpace(strings.ToLower(req.EmailOrUsername))
	if req.EmailOrUsername == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Email/username and password are required"})
		return
	}

	var u User
	var passHash string
	err := h.db.QueryRow(r.Context(), `
		SELECT id, username, email, password_hash, display_name, avatar_url, bio, role, status, email_verified, created_at
		FROM users
		WHERE LOWER(username) = $1 OR LOWER(email) = $1
	`, req.EmailOrUsername).Scan(
		&u.ID, &u.Username, &u.Email, &passHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.Role, &u.Status, &u.EmailVerified, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Authentication error"})
		return
	}

	if u.Status == "SUSPENDED" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "This account is suspended"})
		return
	}

	valid, err := VerifyPassword(req.Password, passHash)
	if err != nil || !valid {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	session, err := h.sm.CreateSession(r.Context(), u.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Could not create session"})
		return
	}

	h.setSessionCookie(w, session.ID, session.ExpiresAt)
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Login successful",
		"user":    u,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(SessionCookieName); err == nil && cookie.Value != "" {
		_ = h.sm.DeleteSession(r.Context(), cookie.Value)
	}
	h.clearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": nil})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if strings.TrimSpace(req.DisplayName) != "" {
		user.DisplayName = strings.TrimSpace(req.DisplayName)
	}
	user.Bio = strings.TrimSpace(req.Bio)
	if req.AvatarURL != "" {
		user.AvatarURL = strings.TrimSpace(req.AvatarURL)
	}

	_, err := h.db.Exec(r.Context(), `
		UPDATE users
		SET display_name = $1, bio = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $4
	`, user.DisplayName, user.Bio, user.AvatarURL, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Profile updated",
		"user":    user,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
