package handler_test

import (
	"context"
	"testing"

	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/db"
)

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name     string
		claims   *auth.Claims
		auths    *mockAuthService
		wantName string
		wantErr  error
	}{
		{
			name:   "success",
			claims: &auth.Claims{UserID: 1, Role: "customer"},
			auths: &mockAuthService{
				getUserFn: func(_ context.Context, id int64) (auth.User, error) {
					return auth.User{
						ID:    id,
						Name:  "Alice",
						Email: "alice@example.com",
						Role:  "customer",
					}, nil
				},
			},
			wantName: "Alice",
		},
		{
			name:    "no claims in context",
			claims:  nil,
			wantErr: auth.ErrUnauthorized,
		},
		{
			name:   "user not found",
			claims: &auth.Claims{UserID: 99, Role: "customer"},
			auths: &mockAuthService{
				getUserFn: func(context.Context, int64) (auth.User, error) {
					return auth.User{}, db.ErrNotFound
				},
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, nil, tt.auths)
			ctx := context.Background()
			if tt.claims != nil {
				ctx = auth.ContextWithClaims(ctx, *tt.claims)
			}

			got, err := h.GetCurrentUser(ctx)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tt.wantName {
				t.Errorf("got name %q, want %q",
					got.Name, tt.wantName)
			}
		})
	}
}
