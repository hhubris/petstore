package handler_test

import (
	"context"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/db"
)

func TestDeletePet(t *testing.T) {
	tests := []struct {
		name    string
		params  api.DeletePetParams
		pets    *mockPetService
		wantErr error
	}{
		{
			name:   "success",
			params: api.DeletePetParams{ID: 1},
			pets: &mockPetService{
				deletePetFn: func(context.Context, int64) error {
					return nil
				},
			},
		},
		{
			name:   "not found",
			params: api.DeletePetParams{ID: 99},
			pets: &mockPetService{
				deletePetFn: func(context.Context, int64) error {
					return db.ErrNotFound
				},
			},
			wantErr: db.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, tt.pets, nil)
			err := h.DeletePet(context.Background(), tt.params)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
