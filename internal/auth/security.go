package auth

import (
	"context"
	"errors"

	"github.com/hhubris/petstore/internal/api"
)

// ErrForbidden is returned when a user lacks the required
// role for an operation.
var ErrForbidden = errors.New("forbidden: admin required")

// adminOperations lists operations that require the admin
// role. ogen does not populate CookieAuth.Roles from the
// x-required-role vendor extension, so we maintain this
// map ourselves.
var adminOperations = map[api.OperationName]bool{
	api.AddPetOperation:    true,
	api.DeletePetOperation: true,
}

// SecurityHandler implements api.SecurityHandler by
// validating JWTs from cookies and enforcing role
// requirements.
type SecurityHandler struct {
	token *TokenConfig
}

// NewSecurityHandler returns a SecurityHandler that uses
// the given TokenConfig for JWT validation.
func NewSecurityHandler(
	token *TokenConfig,
) *SecurityHandler {
	return &SecurityHandler{token: token}
}

// HandleCookieAuth validates the JWT from the cookie,
// checks role requirements, and stores Claims in ctx.
func (sh *SecurityHandler) HandleCookieAuth(
	ctx context.Context,
	operationName api.OperationName,
	t api.CookieAuth,
) (context.Context, error) {
	claims, err := sh.token.ParseToken(t.APIKey)
	if err != nil {
		return ctx, ErrInvalidToken
	}

	if adminOperations[operationName] &&
		claims.Role != "admin" {
		return ctx, ErrForbidden
	}

	return ContextWithClaims(ctx, claims), nil
}
