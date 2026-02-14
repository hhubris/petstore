package pet

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/db"
)

// PetRepository provides database access for pets.
type PetRepository struct {
	db db.DBTX
}

// NewPetRepository returns a PetRepository backed by the
// given DBTX.
func NewPetRepository(dbtx db.DBTX) *PetRepository {
	return &PetRepository{db: dbtx}
}

// Create inserts a new pet and returns it with the
// generated ID.
func (r *PetRepository) Create(
	ctx context.Context,
	name string,
	tag *string,
) (api.Pet, error) {
	var pet api.Pet
	err := r.db.QueryRow(ctx,
		"INSERT INTO pets (name, tag) VALUES ($1, $2) "+
			"RETURNING id, name, tag",
		name, tag,
	).Scan(&pet.ID, &pet.Name, &tag)
	if err != nil {
		return api.Pet{}, fmt.Errorf("create pet: %w", err)
	}
	if tag != nil {
		pet.Tag = api.NewOptString(*tag)
	}
	return pet, nil
}

// FindByID returns the pet with the given ID, or
// db.ErrNotFound if it does not exist.
func (r *PetRepository) FindByID(
	ctx context.Context,
	id int64,
) (api.Pet, error) {
	var (
		pet api.Pet
		tag *string
	)
	err := r.db.QueryRow(ctx,
		"SELECT id, name, tag FROM pets WHERE id = $1",
		id,
	).Scan(&pet.ID, &pet.Name, &tag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Pet{}, db.ErrNotFound
		}
		return api.Pet{}, fmt.Errorf("find pet by id: %w", err)
	}
	if tag != nil {
		pet.Tag = api.NewOptString(*tag)
	}
	return pet, nil
}

// FindAll returns pets, optionally filtered by tags and
// limited to a maximum number of results.
func (r *PetRepository) FindAll(
	ctx context.Context,
	tags []string,
	limit *int32,
) ([]api.Pet, error) {
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

	var pets []api.Pet
	for rows.Next() {
		var (
			pet api.Pet
			tag *string
		)
		if err := rows.Scan(&pet.ID, &pet.Name, &tag); err != nil {
			return nil, fmt.Errorf("scan pet: %w", err)
		}
		if tag != nil {
			pet.Tag = api.NewOptString(*tag)
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
