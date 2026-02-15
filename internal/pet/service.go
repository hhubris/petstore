package pet

import "context"

// Repository is the persistence interface the service
// depends on. PetRepository satisfies it via duck typing.
type Repository interface {
	Create(ctx context.Context,
		name string, tag *string,
	) (Pet, error)
	FindByID(ctx context.Context,
		id int64,
	) (Pet, error)
	FindAll(ctx context.Context,
		tags []string, limit *int32,
	) ([]Pet, error)
	Delete(ctx context.Context,
		id int64,
	) error
}

// Service implements pet business logic on top of a
// Repository.
type Service struct {
	repo Repository
}

// NewService returns a Service wired to the given
// repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreatePet creates a new pet and returns it with the
// generated ID.
func (s *Service) CreatePet(
	ctx context.Context,
	name string,
	tag *string,
) (Pet, error) {
	return s.repo.Create(ctx, name, tag)
}

// GetPet returns the pet with the given ID.
func (s *Service) GetPet(
	ctx context.Context,
	id int64,
) (Pet, error) {
	return s.repo.FindByID(ctx, id)
}

// ListPets returns pets, optionally filtered by tags and
// limited to a maximum number of results.
func (s *Service) ListPets(
	ctx context.Context,
	tags []string,
	limit *int32,
) ([]Pet, error) {
	return s.repo.FindAll(ctx, tags, limit)
}

// DeletePet removes the pet with the given ID.
func (s *Service) DeletePet(
	ctx context.Context,
	id int64,
) error {
	return s.repo.Delete(ctx, id)
}
