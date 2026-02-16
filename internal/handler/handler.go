package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/db"
	"github.com/hhubris/petstore/internal/pet"
)

// PetService defines the pet operations the handler depends on.
type PetService interface {
	CreatePet(ctx context.Context, name string, tag *string) (pet.Pet, error)
	GetPet(ctx context.Context, id int64) (pet.Pet, error)
	ListPets(ctx context.Context, tags []string, limit *int32) ([]pet.Pet, error)
	DeletePet(ctx context.Context, id int64) error
}

// AuthService defines the auth operations the handler depends on.
type AuthService interface {
	Register(ctx context.Context, name, email, password string) (auth.User, error)
	Login(ctx context.Context, email, password string) (string, auth.User, error)
	GetUser(ctx context.Context, id int64) (auth.User, error)
}

// Handler implements the ogen api.Handler interface.
type Handler struct {
	pets   PetService
	auth   AuthService
	secure bool
}

// New creates a Handler. The secure flag controls the
// Secure attribute on cookies (true in production).
func New(
	pets PetService,
	auth AuthService,
	secure bool,
) *Handler {
	return &Handler{
		pets:   pets,
		auth:   auth,
		secure: secure,
	}
}

// NewError maps service-layer errors to HTTP error responses.
func (h *Handler) NewError(
	_ context.Context, err error,
) *api.ErrorStatusCode {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, db.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, db.ErrConflict):
		code = http.StatusConflict
	case errors.Is(err, auth.ErrInvalidCredentials):
		code = http.StatusUnauthorized
	case errors.Is(err, auth.ErrUnauthorized):
		code = http.StatusUnauthorized
	case errors.Is(err, auth.ErrForbidden):
		code = http.StatusForbidden
	case errors.Is(err, auth.ErrInvalidToken):
		code = http.StatusUnauthorized
	}
	return &api.ErrorStatusCode{
		StatusCode: code,
		Response: api.Error{
			Code:    int32(code),
			Message: err.Error(),
		},
	}
}

// responseWriterKey is the context key for the http.ResponseWriter.
type responseWriterKey struct{}

// WithResponseWriter stores an http.ResponseWriter in the context.
func WithResponseWriter(
	ctx context.Context, w http.ResponseWriter,
) context.Context {
	return context.WithValue(ctx, responseWriterKey{}, w)
}

// responseWriterFromContext retrieves the http.ResponseWriter
// stored by WithResponseWriter.
func responseWriterFromContext(
	ctx context.Context,
) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseWriterKey{}).(http.ResponseWriter)
	return w, ok
}

// WrapWithResponseWriter returns middleware that injects the
// http.ResponseWriter into each request's context so handlers
// can set cookies. Used by internal/server when wiring up
// the ogen server.
func WrapWithResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithResponseWriter(r.Context(), w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// newCookie creates an HTTP cookie with common defaults.
func newCookie(
	name, value string, maxAge int, secure bool,
) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

// petToAPI converts a domain Pet to an API Pet.
func petToAPI(p pet.Pet) api.Pet {
	ap := api.Pet{
		ID:   p.ID,
		Name: p.Name,
	}
	if p.Tag != nil {
		ap.Tag = api.NewOptString(*p.Tag)
	}
	return ap
}

// userToAPI converts a domain User to an API AuthUser.
func userToAPI(u auth.User) api.AuthUser {
	return api.AuthUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  api.AuthUserRole(u.Role),
	}
}
