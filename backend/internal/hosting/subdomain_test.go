package hosting

import (
	"testing"
)

func TestSubdomainValidation(t *testing.T) {
	validCases := []string{
		"my-app",
		"project123",
		"cool-bot-2",
		"asdasd",
		"ab",
	}

	for _, sub := range validCases {
		if !subdomainRegex.MatchString(sub) {
			t.Errorf("expected %q to be valid subdomain, but regex failed", sub)
		}
		if reservedSubdomains[sub] {
			t.Errorf("expected %q not to be reserved", sub)
		}
	}

	invalidCases := []string{
		"-leading-dash",
		"trailing-dash-",
		"under_score",
		"UPPERCASE",
		"space in it",
		"a", // less than 2 chars
		"dot.subdomain",
	}

	for _, sub := range invalidCases {
		if len(sub) >= 2 && subdomainRegex.MatchString(sub) && !reservedSubdomains[sub] {
			t.Errorf("expected %q to be invalid, but passed", sub)
		}
	}

	reservedCases := []string{
		"api", "admin", "www", "status", "auth", "login",
	}

	for _, sub := range reservedCases {
		if !reservedSubdomains[sub] {
			t.Errorf("expected %q to be in reservedSubdomains list", sub)
		}
	}
}
