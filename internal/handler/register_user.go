package handler

import (
	"context"

	"github.com/hhubris/petstore/internal/api"
)

// RegisterUser handles POST /auth/register.
func (h *Handler) RegisterUser(
	ctx context.Context, req *api.RegisterRequest,
) (api.RegisterUserRes, error) {
	u, err := h.auth.Register(
		ctx, req.Name, req.Email, req.Password,
	)
	if err != nil {
		return nil, err
	}
	au := userToAPI(u)
	return &au, nil
}
