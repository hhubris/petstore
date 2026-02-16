package handler

import (
	"context"

	"github.com/hhubris/petstore/internal/api"
)

// FindPets handles GET /pets.
func (h *Handler) FindPets(
	ctx context.Context, params api.FindPetsParams,
) ([]api.Pet, error) {
	var limit *int32
	if v, ok := params.Limit.Get(); ok {
		limit = &v
	}

	pets, err := h.pets.ListPets(ctx, params.Tags, limit)
	if err != nil {
		return nil, err
	}

	out := make([]api.Pet, len(pets))
	for i, p := range pets {
		out[i] = petToAPI(p)
	}
	return out, nil
}
