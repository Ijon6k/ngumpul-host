package access

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"ngumpul-host/backend/internal/auth"
	"ngumpul-host/backend/internal/config"
	"ngumpul-host/backend/internal/instance"
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

// GetRegistrationMode handles public GET /api/auth/mode
func (h *Handler) GetRegistrationMode(w http.ResponseWriter, r *http.Request) {
	mode, err := h.svc.GetRegistrationMode(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to determine registration mode")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{
		"registration_mode": mode,
	})
}

// ValidateInvitation handles public GET /api/invitations/validate?token=...
func (h *Handler) ValidateInvitation(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		response.Error(w, http.StatusBadRequest, "Invitation token is required")
		return
	}

	result, err := h.svc.ValidateInvitation(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to validate invitation token")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// GetSettings handles admin GET /api/admin/settings
func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	mode, err := h.svc.GetRegistrationMode(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to load instance settings")
		return
	}

	domain, _ := instance.GetDomain(r.Context(), h.svc.db)

	response.JSON(w, http.StatusOK, InstanceSettingsResponse{
		RegistrationMode: mode,
		Domain:           domain,
		UpdatedAt:        time.Now(),
	})
}

// UpdateSettings handles admin PATCH /api/admin/settings
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.RegistrationMode != "" {
		req.RegistrationMode = strings.ToUpper(strings.TrimSpace(req.RegistrationMode))
		if err := h.svc.SetRegistrationMode(r.Context(), req.RegistrationMode, user.ID); err != nil {
			if err == ErrInvalidMode {
				response.Error(w, http.StatusBadRequest, "Invalid registration mode. Valid modes: OPEN, INVITE_ONLY, CLOSED")
				return
			}
			response.Error(w, http.StatusInternalServerError, "Failed to update instance settings")
			return
		}
	}

	cleanDomain := strings.TrimSpace(req.Domain)
	if cleanDomain != "" {
		cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
		cleanDomain = strings.TrimPrefix(cleanDomain, "https://")
		cleanDomain = strings.TrimRight(cleanDomain, "/")

		_, err := h.svc.db.Exec(r.Context(), `
			INSERT INTO instance_settings (key, value, updated_at)
			VALUES ('domain', $1, NOW())
			ON CONFLICT (key) DO UPDATE SET value = $1, updated_at = NOW()
		`, cleanDomain)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update domain setting")
			return
		}

		appURL := fmt.Sprintf("https://%s", cleanDomain)
		if strings.Contains(cleanDomain, "localhost") || strings.Contains(cleanDomain, "127.0.0.1") {
			appURL = fmt.Sprintf("http://%s", cleanDomain)
		}
		_, _ = h.svc.db.Exec(r.Context(), `
			INSERT INTO instance_settings (key, value, updated_at)
			VALUES ('app_url', $1, NOW())
			ON CONFLICT (key) DO UPDATE SET value = $1, updated_at = NOW()
		`, appURL)
	}

	currentDomain, _ := instance.GetDomain(r.Context(), h.svc.db)

	response.JSON(w, http.StatusOK, map[string]any{
		"message":           "Settings updated successfully",
		"registration_mode": req.RegistrationMode,
		"domain":            currentDomain,
	})
}

// ListInvitations handles admin GET /api/admin/invitations
func (h *Handler) ListInvitations(w http.ResponseWriter, r *http.Request) {
	invitations, err := h.svc.ListInvitations(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to list invitations")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"invitations": invitations,
	})
}

// CreateInvitation handles admin POST /api/admin/invitations
func (h *Handler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	rawToken, inv, err := h.svc.CreateInvitation(r.Context(), user.ID, req.InvitedEmail, req.MaxUses, req.ExpiresInDays)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate invitation")
		return
	}

	// Prefer the canonical domain stored in instance_settings (set on the setup
	// page). Fall back to cfg.AppURL only when no instance domain is configured.
	baseURL := strings.TrimRight(h.cfg.AppURL, "/")
	if appURL, appErr := instance.GetAppURL(r.Context(), h.svc.db); appErr != nil {
		log.Printf("failed to read stored app_url for invite link: %v", appErr)
	} else if appURL != "" {
		baseURL = strings.TrimRight(appURL, "/")
	}
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	inviteURL := fmt.Sprintf("%s/invite/%s", baseURL, rawToken)

	response.JSON(w, http.StatusCreated, CreateInvitationResponse{
		Invitation: *inv,
		RawToken:   rawToken,
		InviteURL:  inviteURL,
	})
}

// RevokeInvitation handles admin POST /api/admin/invitations/{id}/revoke
func (h *Handler) RevokeInvitation(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "Invitation ID is required")
		return
	}

	if err := h.svc.RevokeInvitation(r.Context(), id, user.ID); err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Invitation revoked successfully",
	})
}
