package handler

import (
	"context"

	"github.com/hhubris/petstore/internal/api"
)

// FindPetByID handles GET /pets/{id}.
func (h *Handler) FindPetByID(
	ctx context.Context, params api.FindPetByIDParams,
) (*api.Pet, error) {
	p, err := h.pets.GetPet(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	ap := petToAPI(p)
	return &ap, nil
}
