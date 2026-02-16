package pet

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hhubris/petstore/internal/db"
)

// dbtx is the database interface required by
// PetRepository. Satisfied by *pgxpool.Pool, pgx.Tx, and
// pgxmock.
type dbtx interface {
	Query(ctx context.Context, sql string,
		args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string,
		args ...any) pgx.Row
	Exec(ctx context.Context, sql string,
		args ...any) (pgconn.CommandTag, error)
}

// PetRepository provides database access for pets.
type PetRepository struct {
	db dbtx
}

// NewPetRepository returns a PetRepository backed by the
// given database connection.
func NewPetRepository(conn dbtx) *PetRepository {
	return &PetRepository{db: conn}
}

// Create inserts a new pet and returns it with the
// generated ID.
func (r *PetRepository) Create(
	ctx context.Context,
	name string,
	tag *string,
) (Pet, error) {
	var pet Pet
	err := r.db.QueryRow(ctx,
		"INSERT INTO pets (name, tag) VALUES ($1, $2) "+
			"RETURNING id, name, tag",
		name, tag,
	).Scan(&pet.ID, &pet.Name, &pet.Tag)
	if err != nil {
		return Pet{}, fmt.Errorf("create pet: %w", err)
	}
	return pet, nil
}

// FindByID returns the pet with the given ID, or
// db.ErrNotFound if it does not exist.
func (r *PetRepository) FindByID(
	ctx context.Context,
	id int64,
) (Pet, error) {
	var pet Pet
	err := r.db.QueryRow(ctx,
		"SELECT id, name, tag FROM pets WHERE id = $1",
		id,
	).Scan(&pet.ID, &pet.Name, &pet.Tag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Pet{}, db.ErrNotFound
		}
		return Pet{}, fmt.Errorf("find pet by id: %w", err)
	}
	return pet, nil
}

// FindAll returns pets, optionally filtered by tags and
// limited to a maximum number of results.
func (r *PetRepository) FindAll(
	ctx context.Context,
	tags []string,
	limit *int32,
) ([]Pet, error) {
	var (
		query strings.Builder
		args  []any
		argN  int
	)
	query.WriteString("SELECT id, name, tag FROM pets")

	if len(tags) > 0 {
		query.WriteString(" WHERE tag IN (")
		for i, t := range tags {
			if i > 0 {
				query.WriteString(", ")
			}
			argN++
			query.WriteString("$" + strconv.Itoa(argN))
			args = append(args, t)
		}
		query.WriteString(")")
	}

	query.WriteString(" ORDER BY id")

	if limit != nil {
		argN++
		query.WriteString(" LIMIT $" + strconv.Itoa(argN))
		args = append(args, *limit)
	}

	rows, err := r.db.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("find pets: %w", err)
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Tag); err != nil {
			return nil, fmt.Errorf("scan pet: %w", err)
		}
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pets: %w", err)
	}
	return pets, nil
}

// Delete removes the pet with the given ID, or returns
// db.ErrNotFound if it does not exist.
func (r *PetRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	tag, err := r.db.Exec(ctx,
		"DELETE FROM pets WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("delete pet: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return db.ErrNotFound
	}
	return nil
}
