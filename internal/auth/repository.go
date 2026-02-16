package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hhubris/petstore/internal/db"
)

// dbtx is the database interface required by
// UserRepository. Satisfied by *pgxpool.Pool, pgx.Tx, and
// pgxmock.
type dbtx interface {
	QueryRow(ctx context.Context, sql string,
		args ...any) pgx.Row
}

// UniqueViolation is the PostgreSQL error code for a
// unique constraint violation.
const uniqueViolation = "23505"

// UserRepository provides database access for users.
type UserRepository struct {
	db dbtx
}

// NewUserRepository returns a UserRepository backed by the
// given database connection.
func NewUserRepository(conn dbtx) *UserRepository {
	return &UserRepository{db: conn}
}

// Create inserts a new user and returns it with the
// generated ID and timestamps. Returns db.ErrConflict if
// the email already exists.
func (r *UserRepository) Create(
	ctx context.Context,
	name, email, passwordHash, role string,
) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		"INSERT INTO users (name, email, password_hash, role) "+
			"VALUES ($1, $2, $3, $4) "+
			"RETURNING id, name, email, password_hash, role, "+
			"created_at, updated_at",
		name, email, passwordHash, role,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == uniqueViolation {
			return User{}, db.ErrConflict
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

// FindByEmail returns the user with the given email, or
// db.ErrNotFound if no such user exists.
func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		"SELECT id, name, email, password_hash, role, "+
			"created_at, updated_at "+
			"FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, db.ErrNotFound
		}
		return User{}, fmt.Errorf("find user by email: %w", err)
	}
	return u, nil
}

// FindByID returns the user with the given ID, or
// db.ErrNotFound if no such user exists.
func (r *UserRepository) FindByID(
	ctx context.Context,
	id int64,
) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		"SELECT id, name, email, password_hash, role, "+
			"created_at, updated_at "+
			"FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, db.ErrNotFound
		}
		return User{}, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}
