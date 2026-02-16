package handler

import (
	"context"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/auth"
)

// GetCurrentUser handles GET /auth/me.
func (h *Handler) GetCurrentUser(
	ctx context.Context,
) (*api.AuthUser, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, auth.ErrUnauthorized
	}
	u, err := h.auth.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	au := userToAPI(u)
	return &au, nil
}
