package handler_test

import (
	"context"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/db"
	"github.com/hhubris/petstore/internal/pet"
)

func TestAddPet(t *testing.T) {
	tag := "dog"
	tests := []struct {
		name     string
		req      *api.NewPet
		pets     *mockPetService
		wantName string
		wantErr  error
	}{
		{
			name: "success without tag",
			req:  &api.NewPet{Name: "Fido"},
			pets: &mockPetService{
				createPetFn: func(_ context.Context, name string, tag *string) (pet.Pet, error) {
					if tag != nil {
						t.Error("expected nil tag")
					}
					return pet.Pet{ID: 1, Name: name}, nil
				},
			},
			wantName: "Fido",
		},
		{
			name: "success with tag",
			req: &api.NewPet{
				Name: "Buddy",
				Tag:  api.NewOptString("dog"),
			},
			pets: &mockPetService{
				createPetFn: func(_ context.Context, name string, tg *string) (pet.Pet, error) {
					if tg == nil || *tg != "dog" {
						t.Error("expected tag=dog")
					}
					return pet.Pet{ID: 2, Name: name, Tag: &tag}, nil
				},
			},
			wantName: "Buddy",
		},
		{
			name: "conflict error",
			req:  &api.NewPet{Name: "Fido"},
			pets: &mockPetService{
				createPetFn: func(context.Context, string, *string) (pet.Pet, error) {
					return pet.Pet{}, db.ErrConflict
				},
			},
			wantErr: db.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, tt.pets, nil)
			got, err := h.AddPet(context.Background(), tt.req)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tt.wantName {
				t.Errorf("got name %q, want %q",
					got.Name, tt.wantName)
			}
		})
	}
}
