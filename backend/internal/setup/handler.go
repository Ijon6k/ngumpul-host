package setup

import (
	"encoding/json"
	"errors"
	"net/http"

	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/config"
	"ngumpul-host/backend/internal/response"
)

type Handler struct {
	svc *Service
	cfg *config.Config
}

func NewHandler(svc *Service, cfg *config.Config) *Handler {
	return &Handler{
		svc: svc,
		cfg: cfg,
	}
}

// GetStatus handles GET /api/setup/status
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	initialized, err := h.svc.IsInitialized(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to determine node initialization state")
		return
	}

	domain, err := h.svc.GetDomain(r.Context())
	if err != nil {
		domain = ""
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"initialized": initialized,
		"domain":      domain,
	})
}

// Setup handles POST /api/setup
func (h *Handler) Setup(w http.ResponseWriter, r *http.Request) {
	var req SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, session, err := h.svc.ExecuteSetup(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrAlreadyInitialized) {
			response.Error(w, http.StatusForbidden, "Node setup has already been completed")
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set session cookie for seamless operator authentication
	cookie := &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   h.cfg.SessionSecure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "Node setup completed successfully",
		"user":    user,
	})
}
