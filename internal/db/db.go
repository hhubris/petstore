package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is the common interface satisfied by *pgxpool.Pool,
// pgx.Tx, and pgxmock. Repositories depend on this
// interface so they can be used with any of those
// implementations.
type DBTX interface {
	Query(ctx context.Context, sql string,
		args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string,
		args ...any) pgx.Row
	Exec(ctx context.Context, sql string,
		args ...any) (pgconn.CommandTag, error)
}

// Sentinel errors returned by repositories. Upper layers
// translate these into appropriate HTTP status codes.
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
