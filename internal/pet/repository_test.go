package pet_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/db"
	"github.com/hhubris/petstore/internal/pet"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	tagVal := "dog"

	tests := []struct {
		name    string
		petName string
		tag     *string
		mock    func(m pgxmock.PgxPoolIface)
		want    api.Pet
		wantErr bool
	}{
		{
			name:    "success with tag",
			petName: "Fido",
			tag:     &tagVal,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("INSERT INTO pets").
					WithArgs("Fido", &tagVal).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(1), "Fido", &tagVal),
					)
			},
			want: api.Pet{
				ID:   1,
				Name: "Fido",
				Tag:  api.NewOptString("dog"),
			},
		},
		{
			name:    "success without tag",
			petName: "Luna",
			tag:     nil,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("INSERT INTO pets").
					WithArgs("Luna", (*string)(nil)).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(2), "Luna", (*string)(nil)),
					)
			},
			want: api.Pet{
				ID:   2,
				Name: "Luna",
			},
		},
		{
			name:    "scan error",
			petName: "Bad",
			tag:     nil,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("INSERT INTO pets").
					WithArgs("Bad", (*string)(nil)).
					WillReturnError(errors.New("scan failed"))
			},
			wantErr: true,
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

			repo := pet.NewPetRepository(mock)
			got, err := repo.Create(ctx, tt.petName, tt.tag)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.want.ID ||
				got.Name != tt.want.Name ||
				got.Tag != tt.want.Tag {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()
	tagVal := "cat"

	tests := []struct {
		name    string
		id      int64
		mock    func(m pgxmock.PgxPoolIface)
		want    api.Pet
		wantErr error
	}{
		{
			name: "found",
			id:   1,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT id, name, tag FROM pets").
					WithArgs(int64(1)).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(1), "Whiskers", &tagVal),
					)
			},
			want: api.Pet{
				ID:   1,
				Name: "Whiskers",
				Tag:  api.NewOptString("cat"),
			},
		},
		{
			name: "not found",
			id:   999,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT id, name, tag FROM pets").
					WithArgs(int64(999)).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}),
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

			repo := pet.NewPetRepository(mock)
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
			if got.ID != tt.want.ID ||
				got.Name != tt.want.Name ||
				got.Tag != tt.want.Tag {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	limit10 := int32(10)
	tagDog := "dog"
	tagCat := "cat"

	tests := []struct {
		name    string
		tags    []string
		limit   *int32
		mock    func(m pgxmock.PgxPoolIface)
		want    []api.Pet
		wantErr bool
	}{
		{
			name:  "no filters",
			tags:  nil,
			limit: nil,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT id, name, tag FROM pets ORDER BY id").
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(1), "Fido", &tagDog).
							AddRow(int64(2), "Luna", (*string)(nil)),
					)
			},
			want: []api.Pet{
				{ID: 1, Name: "Fido", Tag: api.NewOptString("dog")},
				{ID: 2, Name: "Luna"},
			},
		},
		{
			name:  "with tags",
			tags:  []string{"dog", "cat"},
			limit: nil,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(
					"SELECT id, name, tag FROM pets WHERE tag IN").
					WithArgs("dog", "cat").
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(1), "Fido", &tagDog).
							AddRow(int64(3), "Mimi", &tagCat),
					)
			},
			want: []api.Pet{
				{ID: 1, Name: "Fido", Tag: api.NewOptString("dog")},
				{ID: 3, Name: "Mimi", Tag: api.NewOptString("cat")},
			},
		},
		{
			name:  "with limit",
			tags:  nil,
			limit: &limit10,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(
					"SELECT id, name, tag FROM pets ORDER BY id LIMIT").
					WithArgs(limit10).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(1), "Fido", &tagDog),
					)
			},
			want: []api.Pet{
				{ID: 1, Name: "Fido", Tag: api.NewOptString("dog")},
			},
		},
		{
			name:  "with tags and limit",
			tags:  []string{"dog"},
			limit: &limit10,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(
					"SELECT id, name, tag FROM pets WHERE tag IN").
					WithArgs("dog", limit10).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}).
							AddRow(int64(1), "Fido", &tagDog),
					)
			},
			want: []api.Pet{
				{ID: 1, Name: "Fido", Tag: api.NewOptString("dog")},
			},
		},
		{
			name:  "empty result",
			tags:  nil,
			limit: nil,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT id, name, tag FROM pets ORDER BY id").
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "name", "tag"}),
					)
			},
			want: nil,
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

			repo := pet.NewPetRepository(mock)
			got, err := repo.FindAll(ctx, tt.tags, tt.limit)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d pets, want %d",
					len(got), len(tt.want))
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID ||
					got[i].Name != tt.want[i].Name ||
					got[i].Tag != tt.want[i].Tag {
					t.Errorf("pet[%d]: got %+v, want %+v",
						i, got[i], tt.want[i])
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int64
		mock    func(m pgxmock.PgxPoolIface)
		wantErr error
	}{
		{
			name: "row affected",
			id:   1,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec("DELETE FROM pets").
					WithArgs(int64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
		{
			name: "no row affected",
			id:   999,
			mock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec("DELETE FROM pets").
					WithArgs(int64(999)).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
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

			repo := pet.NewPetRepository(mock)
			err = repo.Delete(ctx, tt.id)

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
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
