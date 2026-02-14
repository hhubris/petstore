package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"

	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/db"
)

func TestUserCreate(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	tests := []struct {
		name    string
		mock    func(m pgxmock.PgxPoolIface)
		want    auth.User
		wantErr error
	}{
		{
			name: "success",
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs("Alice", "alice@example.com",
						"hashed", "customer").
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "name", "email",
							"password_hash", "role",
							"created_at", "updated_at",
						}).AddRow(
							int64(1), "Alice",
							"alice@example.com", "hashed",
							"customer", now, now,
						),
					)
			},
			want: auth.User{
				ID:           1,
				Name:         "Alice",
				Email:        "alice@example.com",
				PasswordHash: "hashed",
				Role:         "customer",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name: "duplicate email",
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs("Alice", "alice@example.com",
						"hashed", "customer").
					WillReturnError(&pgconn.PgError{
						Code: "23505",
					})
			},
			wantErr: db.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.mock(mock)

			repo := auth.NewUserRepository(mock)
			got, err := repo.Create(ctx,
				"Alice", "alice@example.com",
				"hashed", "customer")

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got error %v, want %v",
						err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestUserFindByEmail(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	tests := []struct {
		name    string
		email   string
		mock    func(m pgxmock.PgxPoolIface)
		want    auth.User
		wantErr error
	}{
		{
			name:  "found",
			email: "alice@example.com",
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT .+ FROM users WHERE email").
					WithArgs("alice@example.com").
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "name", "email",
							"password_hash", "role",
							"created_at", "updated_at",
						}).AddRow(
							int64(1), "Alice",
							"alice@example.com", "hashed",
							"customer", now, now,
						),
					)
			},
			want: auth.User{
				ID:           1,
				Name:         "Alice",
				Email:        "alice@example.com",
				PasswordHash: "hashed",
				Role:         "customer",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name:  "not found",
			email: "nobody@example.com",
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT .+ FROM users WHERE email").
					WithArgs("nobody@example.com").
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "name", "email",
							"password_hash", "role",
							"created_at", "updated_at",
						}),
					)
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.mock(mock)

			repo := auth.NewUserRepository(mock)
			got, err := repo.FindByEmail(ctx, tt.email)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got error %v, want %v",
						err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestUserFindByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Microsecond)

	tests := []struct {
		name    string
		id      int64
		mock    func(m pgxmock.PgxPoolIface)
		want    auth.User
		wantErr error
	}{
		{
			name: "found",
			id:   1,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT .+ FROM users WHERE id").
					WithArgs(int64(1)).
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "name", "email",
							"password_hash", "role",
							"created_at", "updated_at",
						}).AddRow(
							int64(1), "Alice",
							"alice@example.com", "hashed",
							"customer", now, now,
						),
					)
			},
			want: auth.User{
				ID:           1,
				Name:         "Alice",
				Email:        "alice@example.com",
				PasswordHash: "hashed",
				Role:         "customer",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name: "not found",
			id:   999,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT .+ FROM users WHERE id").
					WithArgs(int64(999)).
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "name", "email",
							"password_hash", "role",
							"created_at", "updated_at",
						}),
					)
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			tt.mock(mock)

			repo := auth.NewUserRepository(mock)
			got, err := repo.FindByID(ctx, tt.id)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got error %v, want %v",
						err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
