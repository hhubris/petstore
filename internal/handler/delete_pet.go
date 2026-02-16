package handler

import (
	"context"

	"github.com/hhubris/petstore/internal/api"
)

// DeletePet handles DELETE /pets/{id}.
func (h *Handler) DeletePet(
	ctx context.Context, params api.DeletePetParams,
) error {
	return h.pets.DeletePet(ctx, params.ID)
}
