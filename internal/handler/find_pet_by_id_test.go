package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/db"
	"github.com/hhubris/petstore/internal/pet"
)

func TestFindPetByID(t *testing.T) {
	tests := []struct {
		name     string
		params   api.FindPetByIDParams
		pets     *mockPetService
		wantName string
		wantCode int
	}{
		{
			name:   "found",
			params: api.FindPetByIDParams{ID: 1},
			pets: &mockPetService{
				getPetFn: func(_ context.Context, id int64) (pet.Pet, error) {
					return pet.Pet{ID: id, Name: "Fido"}, nil
				},
			},
			wantName: "Fido",
		},
		{
			name:   "not found",
			params: api.FindPetByIDParams{ID: 99},
			pets: &mockPetService{
				getPetFn: func(context.Context, int64) (pet.Pet, error) {
					return pet.Pet{}, db.ErrNotFound
				},
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler(t, tt.pets, nil)
			got, err := h.FindPetByID(context.Background(), tt.params)
			if tt.wantCode != 0 {
				if err == nil {
					t.Fatal("expected error")
				}
				apiErr := h.NewError(context.Background(), err)
				if apiErr.StatusCode != tt.wantCode {
					t.Errorf("got status %d, want %d",
						apiErr.StatusCode, tt.wantCode)
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

func TestNewError(t *testing.T) {
	h := newHandler(t, &mockPetService{}, &mockAuthService{})
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"not found", db.ErrNotFound, 404},
		{"conflict", db.ErrConflict, 409},
		{"unknown", errors.New("boom"), 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.NewError(context.Background(), tt.err)
			if got.StatusCode != tt.wantCode {
				t.Errorf("got %d, want %d",
					got.StatusCode, tt.wantCode)
			}
		})
	}
}
