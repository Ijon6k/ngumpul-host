package hosting

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	subdomainRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

	reservedSubdomains = map[string]bool{
		"api": true, "admin": true, "www": true, "mail": true, "status": true,
		"me": true, "console": true, "auth": true, "login": true, "register": true,
		"static": true, "uploads": true, "assets": true, "internal": true, "public": true,
		"setup": true, "health": true, "app": true, "node": true, "system": true,
	}
)

var (
	ErrSubdomainEmpty    = errors.New("subdomain is required")
	ErrSubdomainTooShort = errors.New("subdomain must be at least 2 characters")
	ErrSubdomainTooLong  = errors.New("subdomain cannot exceed 63 characters")
	ErrSubdomainInvalid  = errors.New("only lowercase letters, numbers, and hyphens allowed (cannot start or end with hyphen)")
	ErrSubdomainReserved = errors.New("this subdomain is reserved for system services")
)

// NormalizeSubdomain trims whitespace and converts to lowercase.
func NormalizeSubdomain(sub string) string {
	return strings.ToLower(strings.TrimSpace(sub))
}

// IsReservedSubdomain checks if a subdomain name is reserved.
func IsReservedSubdomain(sub string) bool {
	return reservedSubdomains[NormalizeSubdomain(sub)]
}

// ValidateSubdomain checks length, syntax, and reserved list for a subdomain.
func ValidateSubdomain(sub string) error {
	normalized := NormalizeSubdomain(sub)
	if normalized == "" {
		return ErrSubdomainEmpty
	}
	if len(normalized) < 2 {
		return ErrSubdomainTooShort
	}
	if len(normalized) > 63 {
		return ErrSubdomainTooLong
	}
	if !subdomainRegex.MatchString(normalized) {
		return ErrSubdomainInvalid
	}
	if reservedSubdomains[normalized] {
		return fmt.Errorf("%w: %q", ErrSubdomainReserved, normalized)
	}
	return nil
}
