package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/auth"
)

// testSecret is a 32-byte key used across security tests.
var testSecret = []byte("this-is-a-valid-secret-32-bytes!")

func TestSecurityHandlerHandleCookieAuth(t *testing.T) {
	cfg, err := auth.NewTokenConfig(testSecret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}
	sh := auth.NewSecurityHandler(cfg)

	// Helper to create a valid JWT for the given role.
	makeToken := func(t *testing.T, role string) string {
		t.Helper()
		tok, err := cfg.CreateToken(1, role)
		if err != nil {
			t.Fatalf("CreateToken: %v", err)
		}
		return tok
	}

	tests := []struct {
		name      string
		operation api.OperationName
		token     string
		wantErr   error
	}{
		{
			name:      "valid token, no role required",
			operation: api.LogoutUserOperation,
			token:     makeToken(t, "customer"),
			wantErr:   nil,
		},
		{
			name:      "valid token, admin op as admin",
			operation: api.AddPetOperation,
			token:     makeToken(t, "admin"),
			wantErr:   nil,
		},
		{
			name:      "valid token, admin op as customer",
			operation: api.AddPetOperation,
			token:     makeToken(t, "customer"),
			wantErr:   auth.ErrForbidden,
		},
		{
			name:      "valid token, delete op as customer",
			operation: api.DeletePetOperation,
			token:     makeToken(t, "customer"),
			wantErr:   auth.ErrForbidden,
		},
		{
			name:      "invalid token",
			operation: api.LogoutUserOperation,
			token:     "garbage.token.string",
			wantErr:   auth.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cookieAuth := api.CookieAuth{APIKey: tt.token}

			got, err := sh.HandleCookieAuth(
				ctx, tt.operation, cookieAuth,
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal(
						"expected error, got nil",
					)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"error = %v, want %v",
						err, tt.wantErr,
					)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			claims, ok := auth.ClaimsFromContext(got)
			if !ok {
				t.Fatal(
					"expected claims in context",
				)
			}
			if claims.UserID != 1 {
				t.Errorf(
					"UserID = %d, want 1",
					claims.UserID,
				)
			}
		})
	}
}
