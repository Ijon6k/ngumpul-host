package access

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestHashToken(t *testing.T) {
	raw := "ngumpul-inv-secret-token-12345"
	expectedHashBytes := sha256.Sum256([]byte(raw))
	expectedHash := hex.EncodeToString(expectedHashBytes[:])

	hash := HashToken(raw)
	if hash != expectedHash {
		t.Fatalf("expected hash %q, got %q", expectedHash, hash)
	}

	// Trailing/leading whitespace trimming check
	hashWithSpaces := HashToken("  " + raw + " \n")
	if hashWithSpaces != expectedHash {
		t.Fatalf("expected trimmed hash %q, got %q", expectedHash, hashWithSpaces)
	}
}

func TestGenerateToken(t *testing.T) {
	tok1, err := GenerateToken()
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}
	if len(tok1) != 64 { // 32 bytes hex encoded = 64 characters
		t.Fatalf("expected 64 characters, got %d", len(tok1))
	}

	tok2, err := GenerateToken()
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}
	if tok1 == tok2 {
		t.Fatalf("two independently generated tokens should not collide: %s == %s", tok1, tok2)
	}
}
