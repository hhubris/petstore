package handler_test

import (
	"context"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/db"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name     string
		req      *api.RegisterRequest
		auths    *mockAuthService
		wantName string
		wantErr  error
	}{
		{
			name: "success",
			req: &api.RegisterRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Password: "secret123",
			},
			auths: &mockAuthService{
				registerFn: func(_ context.Context, name, email, _ string) (auth.User, error) {
					return auth.User{
						ID:    1,
						Name:  name,
						Email: email,
						Role:  "customer",
					}, nil
				},
			},
			wantName: "Alice",
		},
		{
			name: "conflict",
			req: &api.RegisterRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Password: "secret123",
			},
			auths: &mockAuthService{
				registerFn: func(context.Context, string, string, string) (auth.User, error) {
					return auth.User{}, db.ErrConflict
				},
			},
			wantErr: db.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, nil, tt.auths)
			got, err := h.RegisterUser(
				context.Background(), tt.req,
			)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			au, ok := got.(*api.AuthUser)
			if !ok {
				t.Fatal("expected *api.AuthUser response")
			}
			if au.Name != tt.wantName {
				t.Errorf("got name %q, want %q",
					au.Name, tt.wantName)
			}
		})
	}
}
