package auth

import (
	"context"
	"net/http"

	"ngumpul-host/backend/internal/response"
)

type contextKey string

const UserCtxKey contextKey = "user"

const SessionCookieName = "ngumpul_session"

// Middleware extracts the session cookie and attaches the authenticated user to the request context.
func Middleware(sm *SessionManager) func(http.Handler) http.Handler {
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
				response.Error(w, http.StatusForbidden, "Account is suspended")
				return
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthMiddleware is an alias for Middleware to support existing callers.
var AuthMiddleware = Middleware

// GetUser retrieves the authenticated user from the context.
func GetUser(ctx context.Context) *User {
	if u, ok := ctx.Value(UserCtxKey).(*User); ok {
		return u
	}
	return nil
}

// RequireAuth enforces that a valid user is logged in.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil {
			response.Error(w, http.StatusUnauthorized, "Authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin enforces that the logged in user has the ADMIN role.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil {
			response.Error(w, http.StatusUnauthorized, "Authentication required")
			return
		}
		if user.Role != "ADMIN" {
			response.Error(w, http.StatusForbidden, "Administrator privilege required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
