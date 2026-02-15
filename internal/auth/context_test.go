package auth_test

import (
	"context"
	"testing"

	"github.com/hhubris/petstore/internal/auth"
)

func TestClaimsContext(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() context.Context
		wantOK    bool
		wantClaim auth.Claims
	}{
		{
			name: "round-trip stores and retrieves claims",
			setup: func() context.Context {
				return auth.ContextWithClaims(
					context.Background(),
					auth.Claims{UserID: 42, Role: "admin"},
				)
			},
			wantOK:    true,
			wantClaim: auth.Claims{UserID: 42, Role: "admin"},
		},
		{
			name:      "missing returns zero claims and false",
			setup:     context.Background,
			wantOK:    false,
			wantClaim: auth.Claims{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got, ok := auth.ClaimsFromContext(ctx)
			if ok != tt.wantOK {
				t.Fatalf(
					"ClaimsFromContext() ok = %v, want %v",
					ok, tt.wantOK,
				)
			}
			if got != tt.wantClaim {
				t.Fatalf(
					"ClaimsFromContext() = %+v, want %+v",
					got, tt.wantClaim,
				)
			}
		})
	}
}
