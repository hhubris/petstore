package auth

import (
	"context"
	"errors"
)

// ErrUnauthorized is returned when no valid claims are
// present in the request context.
var ErrUnauthorized = errors.New("unauthorized")

// RequireAdmin extracts Claims from ctx and returns
// ErrForbidden if the user does not have the "admin" role,
// or ErrUnauthorized if no claims are present.
func RequireAdmin(ctx context.Context) error {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}
	if claims.Role != "admin" {
		return ErrForbidden
	}
	return nil
}
