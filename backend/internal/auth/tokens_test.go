package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var testSecret = []byte("test-secret-for-auth-package")

// TestMakeAndVerifyJWT tests the happy path: make a token, verify it,
// confirm the subject (user ID) round-trips correctly.
func TestMakeAndVerifyJWT(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT returned an empty string")
	}

	claims, err := Verifyer(token, testSecret)
	if err != nil {
		t.Fatalf("Verifyer rejected a valid token: %v", err)
	}

	if claims.Subject != userID.String() {
		t.Errorf("Subject mismatch: got %s, want %s", claims.Subject, userID.String())
	}

	// JTI (JWT ID) must be present — the denylist depends on it
	if claims.ID == "" {
		t.Error("Token is missing a JTI (JWT ID) — denylist will not work")
	}
}

// TestVerifyerRejectsWrongSecret confirms that a token signed with one
// secret cannot be verified with a different secret.
func TestVerifyerRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	wrongSecret := []byte("completely-different-secret")
	_, err = Verifyer(token, wrongSecret)
	if err == nil {
		t.Fatal("Verifyer accepted a token signed with a different secret — this is a security bug")
	}
}

// TestVerifyerRejectsExpiredToken confirms that a token whose ExpiresAt
// is in the past is rejected rather than silently accepted.
func TestVerifyerRejectsExpiredToken(t *testing.T) {
	userID := uuid.New()

	// Negative duration means the token expired before it was even created
	token, err := MakeJWT(userID, testSecret, -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = Verifyer(token, testSecret)
	if err == nil {
		t.Fatal("Verifyer accepted an already-expired token — sessions will never actually expire")
	}
}

// TestVerifyerRejectsTamperedToken confirms that modifying the token
// string after signing invalidates the signature check.
func TestVerifyerRejectsTamperedToken(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Flip the last character to simulate tampering
	tampered := token[:len(token)-1] + "X"

	_, err = Verifyer(tampered, testSecret)
	if err == nil {
		t.Fatal("Verifyer accepted a tampered token — signature check is not working")
	}
}

// TestMakeRefreshTokenIsUnique confirms that two calls never return the
// same token (collision would let one user steal another's session).
func TestMakeRefreshTokenIsUnique(t *testing.T) {
	a := MakeRefreshToken()
	b := MakeRefreshToken()

	if a == "" || b == "" {
		t.Fatal("MakeRefreshToken returned an empty string")
	}
	if a == b {
		t.Fatal("MakeRefreshToken returned the same value twice — tokens are not random")
	}
}

// TestMakeRefreshTokenLength confirms the token is the expected 64-char
// hex string (32 random bytes → 64 hex characters).
func TestMakeRefreshTokenLength(t *testing.T) {
	token := MakeRefreshToken()
	if len(token) != 64 {
		t.Errorf("Expected refresh token length 64, got %d", len(token))
	}
}
