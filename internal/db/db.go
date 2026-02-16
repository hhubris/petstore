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

// connString builds a PostgreSQL connection string from
// individual environment variables.
func connString() (string, error) {
	user := os.Getenv("PETSTORE_USER")
	if user == "" {
		return "", fmt.Errorf("PETSTORE_USER is required")
	}

	password := os.Getenv("PETSTORE_PASSWORD")
	if password == "" {
		return "", fmt.Errorf("PETSTORE_PASSWORD is required")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	sslMode := "disable"
	if os.Getenv("DB_SSL_ENABLE") == "true" {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/petstore?sslmode=%s",
		user, password, host, port, sslMode,
	), nil
}

// New connects to the database using connection parameters
// from environment variables and verifies connectivity with
// a ping.
func New(ctx context.Context) (*DB, error) {
	url, err := connString()
	if err != nil {
		return nil, err
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

// Query executes a query that returns rows.
func (d *DB) Query(
	ctx context.Context, sql string, args ...any,
) (pgx.Rows, error) {
	return d.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that returns at most one row.
func (d *DB) QueryRow(
	ctx context.Context, sql string, args ...any,
) pgx.Row {
	return d.pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query that doesn't return rows.
func (d *DB) Exec(
	ctx context.Context, sql string, args ...any,
) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, sql, args...)
}

// Close releases all database resources.
func (d *DB) Close() {
	d.pool.Close()
}
