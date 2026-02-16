package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

// DB encapsulates database connectivity so callers never
// import pgx directly.
type DB struct {
	pool *pgxpool.Pool
}

// New reads DATABASE_URL from the environment, connects to
// the database, and verifies connectivity with a ping.
func New(ctx context.Context) (*DB, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf(
			"creating connection pool: %w", err,
		)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf(
			"pinging database: %w", err,
		)
	}

	return &DB{pool: pool}, nil
}

// DBTX returns the underlying DBTX interface for use by
// repositories.
func (d *DB) DBTX() DBTX {
	return d.pool
}

// Close releases all database resources.
func (d *DB) Close() {
	d.pool.Close()
}
