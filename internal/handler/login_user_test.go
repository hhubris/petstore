package handler_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/auth"
)

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name       string
		req        *api.LoginRequest
		auths      *mockAuthService
		wantName   string
		wantCookie bool
		wantErr    error
	}{
		{
			name: "success",
			req: &api.LoginRequest{
				Email:    "alice@example.com",
				Password: "secret123",
			},
			auths: &mockAuthService{
				loginFn: func(_ context.Context, email, _ string) (string, auth.User, error) {
					return "jwt-token", auth.User{
						ID:    1,
						Name:  "Alice",
						Email: email,
						Role:  "customer",
					}, nil
				},
			},
			wantName:   "Alice",
			wantCookie: true,
		},
		{
			name: "invalid credentials",
			req: &api.LoginRequest{
				Email:    "alice@example.com",
				Password: "wrong",
			},
			auths: &mockAuthService{
				loginFn: func(context.Context, string, string) (string, auth.User, error) {
					return "", auth.User{}, auth.ErrInvalidCredentials
				},
			},
			wantErr: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, nil, tt.auths)
			rec := httptest.NewRecorder()
			ctx := ctxWithResponseWriter(rec)

			got, err := h.LoginUser(ctx, tt.req)
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

			if tt.wantCookie {
				cookies := rec.Result().Cookies()
				found := false
				for _, c := range cookies {
					if c.Name == "access_token" {
						found = true
						if c.Value != "jwt-token" {
							t.Errorf("got cookie value %q, want %q",
								c.Value, "jwt-token")
						}
						if !c.HttpOnly {
							t.Error("expected HttpOnly cookie")
						}
					}
				}
				if !found {
					t.Error("access_token cookie not set")
				}
			}
		})
	}
}

func TestLoginUser_NoResponseWriter(t *testing.T) {
	h := newHandler(t, nil, &mockAuthService{
		loginFn: func(context.Context, string, string) (string, auth.User, error) {
			return "token", auth.User{ID: 1}, nil
		},
	})
	_, err := h.LoginUser(context.Background(), &api.LoginRequest{
		Email:    "a@b.com",
		Password: "pass",
	})
	if err == nil {
		t.Fatal("expected error when response writer missing")
	}
}
