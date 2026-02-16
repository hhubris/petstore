package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hhubris/petstore/internal/api"
)

// LoginUser handles POST /auth/login.
func (h *Handler) LoginUser(
	ctx context.Context, req *api.LoginRequest,
) (api.LoginUserRes, error) {
	token, u, err := h.auth.Login(
		ctx, req.Email, req.Password,
	)
	if err != nil {
		return nil, err
	}

	w, ok := responseWriterFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf(
			"response writer not in context",
		)
	}
	http.SetCookie(w, newCookie(
		"access_token", token, 86400, h.secure,
	))

	au := userToAPI(u)
	return &au, nil
}
