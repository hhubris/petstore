package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hhubris/petstore/internal/auth"
)

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
	}{
		{
			name: "admin user",
			ctx: auth.ContextWithClaims(
				context.Background(),
				auth.Claims{UserID: 1, Role: "admin"},
			),
			wantErr: nil,
		},
		{
			name: "non-admin user",
			ctx: auth.ContextWithClaims(
				context.Background(),
				auth.Claims{UserID: 2, Role: "customer"},
			),
			wantErr: auth.ErrForbidden,
		},
		{
			name:    "no claims in context",
			ctx:     context.Background(),
			wantErr: auth.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.RequireAdmin(tt.ctx)

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
		})
	}
}
