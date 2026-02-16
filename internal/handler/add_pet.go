package handler

import (
	"context"

	"github.com/hhubris/petstore/internal/api"
)

// AddPet handles POST /pets.
func (h *Handler) AddPet(
	ctx context.Context, req *api.NewPet,
) (*api.Pet, error) {
	var tag *string
	if v, ok := req.Tag.Get(); ok {
		tag = &v
	}

	p, err := h.pets.CreatePet(ctx, req.Name, tag)
	if err != nil {
		return nil, err
	}
	ap := petToAPI(p)
	return &ap, nil
}
