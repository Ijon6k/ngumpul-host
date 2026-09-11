package project

import "testing"

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"My  Awesome Project! 123", "my-awesome-project-123"},
		{"---Already-Clean---", "already-clean"},
		{"", "project"},
		{"Special #$% Chars", "special-chars"},
	}

	for _, tc := range tests {
		got := GenerateSlug(tc.input)
		if got != tc.expected {
			t.Errorf("GenerateSlug(%q) = %q, expected %q", tc.input, got, tc.expected)
		}
	}
}
