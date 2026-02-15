package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// validSecret is a 32-byte key used across tests.
var validSecret = []byte("this-is-a-valid-secret-32-bytes!")

func TestNewTokenConfig(t *testing.T) {
	tests := []struct {
		name    string
		secret  []byte
		wantErr bool
	}{
		{
			name:    "valid 32-byte secret",
			secret:  validSecret,
			wantErr: false,
		},
		{
			name:    "too short secret",
			secret:  []byte("short"),
			wantErr: true,
		},
		{
			name:    "empty secret",
			secret:  []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewTokenConfig(tt.secret)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg == nil {
				t.Fatal("expected non-nil config")
			}
		})
	}
}

func TestJWTRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		userID int64
		role   string
	}{
		{
			name:   "admin user",
			userID: 42,
			role:   "admin",
		},
		{
			name:   "regular user",
			userID: 1,
			role:   "user",
		},
		{
			name:   "large user ID",
			userID: 9999999,
			role:   "editor",
		},
	}

	cfg, err := NewTokenConfig(validSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := cfg.CreateToken(tt.userID, tt.role)
			if err != nil {
				t.Fatalf("CreateToken: %v", err)
			}

			claims, err := cfg.ParseToken(token)
			if err != nil {
				t.Fatalf("ParseToken: %v", err)
			}

			if claims.UserID != tt.userID {
				t.Errorf(
					"UserID = %d, want %d",
					claims.UserID, tt.userID,
				)
			}
			if claims.Role != tt.role {
				t.Errorf(
					"Role = %q, want %q",
					claims.Role, tt.role,
				)
			}
		})
	}
}

func TestJWTExpiredToken(t *testing.T) {
	cfg, err := NewTokenConfig(validSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	// Create a token with time in the past.
	past := time.Now().Add(-2 * time.Hour)
	cfg.timeNow = func() time.Time { return past }

	token, err := cfg.CreateToken(1, "user")
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	// Parse with current time — token should be expired.
	cfg.timeNow = time.Now
	_, err = cfg.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestJWTTamperedToken(t *testing.T) {
	cfg, err := NewTokenConfig(validSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	token, err := cfg.CreateToken(1, "user")
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	// Tamper with the token by flipping a character.
	tampered := token[:len(token)-1] + "X"

	_, err = cfg.ParseToken(tampered)
	if err == nil {
		t.Fatal("expected error for tampered token, got nil")
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestJWTWrongSigningKey(t *testing.T) {
	cfg1, err := NewTokenConfig(validSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	otherSecret := []byte("another-valid-secret-32-bytes!ab")
	cfg2, err := NewTokenConfig(otherSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	token, err := cfg1.CreateToken(1, "user")
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	_, err = cfg2.ParseToken(token)
	if err == nil {
		t.Fatal(
			"expected error for wrong signing key, got nil",
		)
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestJWTInvalidTokenStrings(t *testing.T) {
	cfg, err := NewTokenConfig(validSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	tests := []struct {
		name  string
		token string
	}{
		{name: "empty string", token: ""},
		{name: "garbage", token: "not.a.jwt"},
		{name: "partial", token: "eyJhbGciOiJIUzI1NiJ9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cfg.ParseToken(tt.token)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrInvalidToken) {
				t.Errorf(
					"expected ErrInvalidToken, got: %v", err,
				)
			}
		})
	}
}

func TestJWTSigningMethodEnforced(t *testing.T) {
	// Create a token signed with "none" method — should be
	// rejected.
	cfg, err := NewTokenConfig(validSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}

	claims := jwt.MapClaims{
		"sub":  "1",
		"role": "admin",
		"iat":  jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(
			time.Now().Add(time.Hour),
		),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := token.SignedString(
		jwt.UnsafeAllowNoneSignatureType,
	)
	if err != nil {
		t.Fatalf("signing none token: %v", err)
	}

	_, err = cfg.ParseToken(signed)
	if err == nil {
		t.Fatal("expected error for none alg, got nil")
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got: %v", err)
	}
}
