package handler

import (
	"context"
	"fmt"
	"net/http"
)

// LogoutUser handles POST /auth/logout.
func (h *Handler) LogoutUser(ctx context.Context) error {
	w, ok := responseWriterFromContext(ctx)
	if !ok {
		return fmt.Errorf("response writer not in context")
	}
	http.SetCookie(w, newCookie(
		"access_token", "", -1, h.secure,
	))
	return nil
}
