package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const UserCtxKey contextKey = "user"

const SessionCookieName = "ngumpul_session"

func AuthMiddleware(sm *SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				// No session cookie, continue as unauthenticated visitor
				next.ServeHTTP(w, r)
				return
			}

			user, err := sm.GetUserBySession(r.Context(), cookie.Value)
			if err != nil || user == nil {
				// Invalid or expired session
				next.ServeHTTP(w, r)
				return
			}

			if user.Status == "SUSPENDED" {
				http.Error(w, `{"error":"Account is suspended"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) *User {
	if u, ok := ctx.Value(UserCtxKey).(*User); ok {
		return u
	}
	return nil
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Authentication required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Authentication required",
			})
			return
		}
		if user.Role != "ADMIN" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Administrator privilege required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
