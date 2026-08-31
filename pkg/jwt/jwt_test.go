package jwt

import (
	"errors"
	"testing"
	"time"
)

func newTestManager() *Manager {
	return NewManager(Options{
		Issuer:        "user-service-test",
		AccessSecret:  "access-secret-access-secret-1234",
		RefreshSecret: "refresh-secret-refresh-secret-56",
		AccessTTL:     time.Minute,
		RefreshTTL:    time.Hour,
	})
}

func TestGenerateAndParseAccess(t *testing.T) {
	m := newTestManager()

	token, err := m.GenerateAccess(42, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := m.ParseAccess(token.Value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "alice" {
		t.Errorf("unexpected claims: %+v", claims)
	}
	if claims.ID != token.TokenID {
		t.Errorf("token id mismatch: %q vs %q", claims.ID, token.TokenID)
	}
}

func TestRefreshKeepsFamilyID(t *testing.T) {
	m := newTestManager()

	first, err := m.GenerateRefresh(42, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.FamilyID == "" {
		t.Fatal("family id should be generated when empty")
	}

	rotated, err := m.GenerateRefresh(42, first.FamilyID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated.FamilyID != first.FamilyID {
		t.Errorf("family id should be preserved: %q vs %q", rotated.FamilyID, first.FamilyID)
	}
	if rotated.TokenID == first.TokenID {
		t.Error("rotated token should get a new token id")
	}
}

// refresh 令牌必须无法当作 access 令牌使用，反之亦然。
func TestTokenTypeIsolation(t *testing.T) {
	m := newTestManager()

	refresh, err := m.GenerateRefresh(42, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := m.ParseAccess(refresh.Value); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("refresh token must be rejected by ParseAccess, got %v", err)
	}

	access, err := m.GenerateAccess(42, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := m.ParseRefresh(access.Value); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("access token must be rejected by ParseRefresh, got %v", err)
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	m := NewManager(Options{
		Issuer:       "user-service-test",
		AccessSecret: "access-secret-access-secret-1234",
		AccessTTL:    -time.Minute,
	})

	token, err := m.GenerateAccess(42, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := m.ParseAccess(token.Value); !errors.Is(err, ErrTokenExpired) {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestParseRejectsForeignSignature(t *testing.T) {
	issued := newTestManager()
	other := NewManager(Options{
		Issuer:       "user-service-test",
		AccessSecret: "another-secret-another-secret-78",
		AccessTTL:    time.Minute,
	})

	token, err := issued.GenerateAccess(42, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := other.ParseAccess(token.Value); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
