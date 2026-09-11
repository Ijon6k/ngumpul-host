package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"ngumpul-host/backend/internal/auth"
)

func TestUpdateUserRole_SelfDemotionRejected(t *testing.T) {
	h := &Handler{} // db not queried when self-demotion is intercepted

	admin := &auth.User{
		ID:       "admin-123",
		Username: "superadmin",
		Role:     "ADMIN",
		Status:   "ACTIVE",
	}

	payload, _ := json.Marshal(UpdateRolePayload{Role: "USER"})
	req := httptest.NewRequest(http.MethodPatch, "/api/admin/users/admin-123/role", bytes.NewReader(payload))

	// Setup chi URL param and auth context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "admin-123")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.UserCtxKey, admin)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	h.UpdateUserRole(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for self-demotion, got %d", rec.Code)
	}

	var res map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&res)
	if res["error"] != "You cannot demote your own administrator account" {
		t.Errorf("unexpected error message: %s", res["error"])
	}
}

func TestUpdateUserStatus_SelfSuspensionRejected(t *testing.T) {
	h := &Handler{}

	admin := &auth.User{
		ID:       "admin-123",
		Username: "superadmin",
		Role:     "ADMIN",
		Status:   "ACTIVE",
	}

	payload, _ := json.Marshal(UpdateStatusPayload{Status: "SUSPENDED"})
	req := httptest.NewRequest(http.MethodPatch, "/api/admin/users/admin-123/status", bytes.NewReader(payload))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "admin-123")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.UserCtxKey, admin)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	h.UpdateUserStatus(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for self-suspension, got %d", rec.Code)
	}

	var res map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&res)
	if res["error"] != "You cannot suspend your own administrator account" {
		t.Errorf("unexpected error message: %s", res["error"])
	}
}
