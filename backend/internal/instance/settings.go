package instance

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

// QueryRower is satisfied by *pgxpool.Pool and pgx.Tx.
type QueryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// GetValue reads a single scalar value from instance_settings, returning "" when the key is absent.
func GetValue(ctx context.Context, q QueryRower, key string) (string, error) {
	var value string
	err := q.QueryRow(ctx, "SELECT value FROM instance_settings WHERE key = $1", key).Scan(&value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(value), nil
}

// GetDomain returns the canonical node domain (scheme-less, no trailing slash).
func GetDomain(ctx context.Context, q QueryRower) (string, error) {
	return GetValue(ctx, q, "domain")
}

// GetAppURL returns the canonical app URL, preferring the stored app_url key and
// falling back to a scheme derived from the stored domain. Returns "" when unset.
func GetAppURL(ctx context.Context, q QueryRower) (string, error) {
	appURL, err := GetValue(ctx, q, "app_url")
	if err != nil {
		return "", err
	}
	if appURL != "" {
		return strings.TrimRight(appURL, "/"), nil
	}

	domain, err := GetDomain(ctx, q)
	if err != nil {
		return "", err
	}
	if domain == "" {
		return "", nil
	}
	if strings.Contains(domain, "localhost") || strings.Contains(domain, "127.0.0.1") {
		return "http://" + domain, nil
	}
	return "https://" + domain, nil
}