package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/handler"
	"github.com/hhubris/petstore/internal/pet"
)

// mockPetService implements handler.PetService for testing.
type mockPetService struct {
	createPetFn func(ctx context.Context, name string, tag *string) (pet.Pet, error)
	getPetFn    func(ctx context.Context, id int64) (pet.Pet, error)
	listPetsFn  func(ctx context.Context, tags []string, limit *int32) ([]pet.Pet, error)
	deletePetFn func(ctx context.Context, id int64) error
}

func (m *mockPetService) CreatePet(ctx context.Context, name string, tag *string) (pet.Pet, error) {
	return m.createPetFn(ctx, name, tag)
}

func (m *mockPetService) GetPet(ctx context.Context, id int64) (pet.Pet, error) {
	return m.getPetFn(ctx, id)
}

func (m *mockPetService) ListPets(ctx context.Context, tags []string, limit *int32) ([]pet.Pet, error) {
	return m.listPetsFn(ctx, tags, limit)
}

func (m *mockPetService) DeletePet(ctx context.Context, id int64) error {
	return m.deletePetFn(ctx, id)
}

// mockAuthService implements handler.AuthService for testing.
type mockAuthService struct {
	registerFn func(ctx context.Context, name, email, password string) (auth.User, error)
	loginFn    func(ctx context.Context, email, password string) (string, auth.User, error)
	getUserFn  func(ctx context.Context, id int64) (auth.User, error)
}

func (m *mockAuthService) Register(ctx context.Context, name, email, password string) (auth.User, error) {
	return m.registerFn(ctx, name, email, password)
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, auth.User, error) {
	return m.loginFn(ctx, email, password)
}

func (m *mockAuthService) GetUser(ctx context.Context, id int64) (auth.User, error) {
	return m.getUserFn(ctx, id)
}

// newHandler is a test helper that constructs a Handler with
// the given mocks and secure=false.
func newHandler(
	t *testing.T,
	pets *mockPetService,
	auths *mockAuthService,
) *handler.Handler {
	t.Helper()
	return handler.New(pets, auths, false)
}

// ctxWithResponseWriter returns a context with an embedded
// http.ResponseWriter for tests that need cookie verification.
func ctxWithResponseWriter(
	w http.ResponseWriter,
) context.Context {
	return handler.WithResponseWriter(context.Background(), w)
}
