package auth_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/db"
)

// mockRepo is a hand-written mock of auth.Repository.
type mockRepo struct {
	createFn      func(ctx context.Context, name, email, passwordHash, role string) (auth.User, error)
	findByEmailFn func(ctx context.Context, email string) (auth.User, error)
	findByIDFn    func(ctx context.Context, id int64) (auth.User, error)
}

func (m *mockRepo) Create(
	ctx context.Context,
	name, email, passwordHash, role string,
) (auth.User, error) {
	return m.createFn(ctx, name, email, passwordHash, role)
}

func (m *mockRepo) FindByEmail(
	ctx context.Context,
	email string,
) (auth.User, error) {
	return m.findByEmailFn(ctx, email)
}

func (m *mockRepo) FindByID(
	ctx context.Context,
	id int64,
) (auth.User, error) {
	return m.findByIDFn(ctx, id)
}

// newTestService returns a Service wired to the given mock
// and a valid TokenConfig.
func newTestService(
	t *testing.T, repo *mockRepo,
) *auth.Service {
	t.Helper()
	secret := []byte("test-secret-that-is-at-least-32-bytes!")
	tc, err := auth.NewTokenConfig(secret)
	if err != nil {
		t.Fatalf("NewTokenConfig: %v", err)
	}
	return auth.NewService(repo, tc)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr error
	}{
		{
			name: "success",
			repo: &mockRepo{
				createFn: func(
					_ context.Context,
					name, email, hash, role string,
				) (auth.User, error) {
					// Verify the password was hashed.
					if err := bcrypt.CompareHashAndPassword(
						[]byte(hash), []byte("s3cret"),
					); err != nil {
						t.Errorf(
							"password not properly hashed: %v",
							err,
						)
					}
					if role != "customer" {
						t.Errorf(
							"role = %q, want %q",
							role, "customer",
						)
					}
					return auth.User{
						ID:    1,
						Name:  name,
						Email: email,
						Role:  role,
					}, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "duplicate email",
			repo: &mockRepo{
				createFn: func(
					context.Context, string, string,
					string, string,
				) (auth.User, error) {
					return auth.User{}, db.ErrConflict
				},
			},
			wantErr: db.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.repo)
			user, err := svc.Register(
				context.Background(),
				"Alice", "alice@example.com", "s3cret",
			)
			if tt.wantErr != nil {
				if !errorIs(err, tt.wantErr) {
					t.Fatalf(
						"err = %v, want %v",
						err, tt.wantErr,
					)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.Email != "alice@example.com" {
				t.Errorf(
					"email = %q, want %q",
					user.Email, "alice@example.com",
				)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("s3cret"), bcrypt.DefaultCost,
	)
	stored := auth.User{
		ID:           1,
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: string(hash),
		Role:         "customer",
	}

	tests := []struct {
		name     string
		email    string
		password string
		repo     *mockRepo
		wantErr  error
	}{
		{
			name:     "success",
			email:    "alice@example.com",
			password: "s3cret",
			repo: &mockRepo{
				findByEmailFn: func(
					_ context.Context, _ string,
				) (auth.User, error) {
					return stored, nil
				},
			},
			wantErr: nil,
		},
		{
			name:     "unknown email",
			email:    "nobody@example.com",
			password: "s3cret",
			repo: &mockRepo{
				findByEmailFn: func(
					_ context.Context, _ string,
				) (auth.User, error) {
					return auth.User{}, db.ErrNotFound
				},
			},
			wantErr: auth.ErrInvalidCredentials,
		},
		{
			name:     "wrong password",
			email:    "alice@example.com",
			password: "wrong",
			repo: &mockRepo{
				findByEmailFn: func(
					_ context.Context, _ string,
				) (auth.User, error) {
					return stored, nil
				},
			},
			wantErr: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.repo)
			token, user, err := svc.Login(
				context.Background(),
				tt.email, tt.password,
			)
			if tt.wantErr != nil {
				if !errorIs(err, tt.wantErr) {
					t.Fatalf(
						"err = %v, want %v",
						err, tt.wantErr,
					)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token == "" {
				t.Error("expected non-empty token")
			}
			if user.ID != stored.ID {
				t.Errorf(
					"user.ID = %d, want %d",
					user.ID, stored.ID,
				)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr error
	}{
		{
			name: "success",
			repo: &mockRepo{
				findByIDFn: func(
					_ context.Context, id int64,
				) (auth.User, error) {
					return auth.User{
						ID: id, Name: "Alice",
					}, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "not found",
			repo: &mockRepo{
				findByIDFn: func(
					context.Context, int64,
				) (auth.User, error) {
					return auth.User{}, db.ErrNotFound
				},
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.repo)
			user, err := svc.GetUser(
				context.Background(), 42,
			)
			if tt.wantErr != nil {
				if !errorIs(err, tt.wantErr) {
					t.Fatalf(
						"err = %v, want %v",
						err, tt.wantErr,
					)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.ID != 42 {
				t.Errorf(
					"user.ID = %d, want 42", user.ID,
				)
			}
		})
	}
}

// errorIs is a small helper so test table entries can use
// errors.Is semantics.
func errorIs(err, target error) bool {
	if target == nil {
		return err == nil
	}
	return errors.Is(err, target)
}
