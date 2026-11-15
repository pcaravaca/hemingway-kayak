package auth

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateAndParse(t *testing.T) {
	secret := []byte("test-secret")
	token, err := Generate(secret, "user-1", "a@b.com", time.Hour)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("Generate() returned an empty token")
	}

	claims, err := Parse(secret, token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-1")
	}
	if claims.Email != "a@b.com" {
		t.Errorf("Email = %q, want %q", claims.Email, "a@b.com")
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		t.Error("ExpiresAt should be in the future")
	}
}

func TestParseWrongSecret(t *testing.T) {
	token, err := Generate([]byte("secret-a"), "user-1", "a@b.com", time.Hour)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := Parse([]byte("secret-b"), token); err == nil {
		t.Fatal("Parse() with a different secret should fail")
	}
}

func TestParseExpired(t *testing.T) {
	token, err := Generate([]byte("secret"), "user-1", "a@b.com", -time.Minute)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := Parse([]byte("secret"), token); err == nil {
		t.Fatal("Parse() should reject an expired token")
	}
}

func TestParseTampered(t *testing.T) {
	token, err := Generate([]byte("secret"), "user-1", "a@b.com", time.Hour)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("malformed token: %d parts", len(parts))
	}
	if _, err := Parse([]byte("secret"), parts[0]+"."+parts[1]+".AAAA"); err == nil {
		t.Fatal("Parse() should reject a tampered signature")
	}
}

func TestRevocation(t *testing.T) {
	r := NewRevocation()
	expires := time.Now().Add(time.Hour)

	if r.IsRevoked("tok") {
		t.Fatal("unknown token should not be revoked")
	}

	r.Revoke("tok", expires)
	if !r.IsRevoked("tok") {
		t.Fatal("revoked token should be reported as revoked")
	}

	r.Revoke("expired", time.Now().Add(-time.Minute))
	if r.IsRevoked("expired") {
		t.Fatal("expired revocation should be ignored")
	}
}
