package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/pet"
)

func TestFindPets(t *testing.T) {
	tag := "dog"
	tests := []struct {
		name    string
		params  api.FindPetsParams
		pets    *mockPetService
		want    int
		wantErr bool
	}{
		{
			name:   "no filters",
			params: api.FindPetsParams{},
			pets: &mockPetService{
				listPetsFn: func(_ context.Context, tags []string, limit *int32) ([]pet.Pet, error) {
					return []pet.Pet{
						{ID: 1, Name: "Fido", Tag: &tag},
						{ID: 2, Name: "Rex"},
					}, nil
				},
			},
			want: 2,
		},
		{
			name: "with limit",
			params: api.FindPetsParams{
				Limit: api.NewOptInt32(1),
			},
			pets: &mockPetService{
				listPetsFn: func(_ context.Context, _ []string, limit *int32) ([]pet.Pet, error) {
					if limit == nil || *limit != 1 {
						t.Error("expected limit=1")
					}
					return []pet.Pet{{ID: 1, Name: "Fido"}}, nil
				},
			},
			want: 1,
		},
		{
			name:   "service error",
			params: api.FindPetsParams{},
			pets: &mockPetService{
				listPetsFn: func(context.Context, []string, *int32) ([]pet.Pet, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, tt.pets, nil)
			got, err := h.FindPets(context.Background(), tt.params)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.want {
				t.Errorf("got %d pets, want %d", len(got), tt.want)
			}
		})
	}
}
