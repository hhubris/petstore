package pet_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hhubris/petstore/internal/db"
	"github.com/hhubris/petstore/internal/pet"
)

// mockRepo is a hand-written mock of pet.Repository.
type mockRepo struct {
	createFn   func(ctx context.Context, name string, tag *string) (pet.Pet, error)
	findByIDFn func(ctx context.Context, id int64) (pet.Pet, error)
	findAllFn  func(ctx context.Context, tags []string, limit *int32) ([]pet.Pet, error)
	deleteFn   func(ctx context.Context, id int64) error
}

func (m *mockRepo) Create(
	ctx context.Context,
	name string,
	tag *string,
) (pet.Pet, error) {
	return m.createFn(ctx, name, tag)
}

func (m *mockRepo) FindByID(
	ctx context.Context,
	id int64,
) (pet.Pet, error) {
	return m.findByIDFn(ctx, id)
}

func (m *mockRepo) FindAll(
	ctx context.Context,
	tags []string,
	limit *int32,
) ([]pet.Pet, error) {
	return m.findAllFn(ctx, tags, limit)
}

func (m *mockRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	return m.deleteFn(ctx, id)
}

func TestServiceCreatePet(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				createFn: func(
					_ context.Context,
					name string, tag *string,
				) (pet.Pet, error) {
					return pet.Pet{
						ID: 1, Name: name, Tag: tag,
					}, nil
				},
			},
		},
		{
			name: "repo error",
			repo: &mockRepo{
				createFn: func(
					context.Context, string, *string,
				) (pet.Pet, error) {
					return pet.Pet{},
						errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := pet.NewService(tt.repo)
			tag := ptrStr("dog")
			got, err := svc.CreatePet(
				context.Background(), "Fido", tag,
			)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != "Fido" {
				t.Errorf(
					"Name = %q, want %q",
					got.Name, "Fido",
				)
			}
			if !ptrStrEq(got.Tag, tag) {
				t.Errorf("Tag = %v, want %v",
					got.Tag, tag)
			}
		})
	}
}

func TestServiceGetPet(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr error
	}{
		{
			name: "found",
			repo: &mockRepo{
				findByIDFn: func(
					_ context.Context, id int64,
				) (pet.Pet, error) {
					return pet.Pet{
						ID: id, Name: "Fido",
					}, nil
				},
			},
		},
		{
			name: "not found",
			repo: &mockRepo{
				findByIDFn: func(
					context.Context, int64,
				) (pet.Pet, error) {
					return pet.Pet{}, db.ErrNotFound
				},
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := pet.NewService(tt.repo)
			got, err := svc.GetPet(
				context.Background(), 42,
			)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
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
			if got.ID != 42 {
				t.Errorf(
					"ID = %d, want 42", got.ID,
				)
			}
		})
	}
}

func TestServiceListPets(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockRepo
		wantLen int
	}{
		{
			name: "success",
			repo: &mockRepo{
				findAllFn: func(
					_ context.Context,
					_ []string, _ *int32,
				) ([]pet.Pet, error) {
					return []pet.Pet{
						{ID: 1, Name: "Fido"},
						{ID: 2, Name: "Luna"},
					}, nil
				},
			},
			wantLen: 2,
		},
		{
			name: "empty",
			repo: &mockRepo{
				findAllFn: func(
					context.Context, []string, *int32,
				) ([]pet.Pet, error) {
					return nil, nil
				},
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := pet.NewService(tt.repo)
			got, err := svc.ListPets(
				context.Background(), nil, nil,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Errorf(
					"len = %d, want %d",
					len(got), tt.wantLen,
				)
			}
		})
	}
}

func TestServiceDeletePet(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr error
	}{
		{
			name: "success",
			repo: &mockRepo{
				deleteFn: func(
					context.Context, int64,
				) error {
					return nil
				},
			},
		},
		{
			name: "not found",
			repo: &mockRepo{
				deleteFn: func(
					context.Context, int64,
				) error {
					return db.ErrNotFound
				},
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := pet.NewService(tt.repo)
			err := svc.DeletePet(
				context.Background(), 1,
			)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
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
		})
	}
}
