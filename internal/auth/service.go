package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/hhubris/petstore/internal/db"
)

// ErrInvalidCredentials is returned when login fails due to
// an unknown email or wrong password. The message is
// intentionally vague to avoid leaking whether the email
// exists.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Repository is the persistence interface the service
// depends on. UserRepository satisfies it via duck typing.
type Repository interface {
	Create(ctx context.Context,
		name, email, passwordHash, role string,
	) (User, error)
	FindByEmail(ctx context.Context,
		email string,
	) (User, error)
	FindByID(ctx context.Context,
		id int64,
	) (User, error)
}

// Service implements authentication business logic on top
// of a Repository and TokenConfig.
type Service struct {
	repo  Repository
	token *TokenConfig
}

// NewService returns a Service wired to the given repository and
// token configuration.
func NewService(repo Repository, token *TokenConfig) *Service {
	return &Service{repo: repo, token: token}
}

// Register creates a new customer account. The plaintext
// password is hashed before storage. Returns db.ErrConflict
// if the email is already taken.
func (s *Service) Register(
	ctx context.Context,
	name, email, password string,
) (User, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return User{}, fmt.Errorf("hashing password: %w", err)
	}
	user, err := s.repo.Create(
		ctx, name, email, string(hash), "customer",
	)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Login authenticates by email and password. On success it
// returns a signed JWT and the User. Both unknown-email and
// wrong-password cases return ErrInvalidCredentials.
func (s *Service) Login(
	ctx context.Context,
	email, password string,
) (string, User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return "", User{}, ErrInvalidCredentials
		}
		return "", User{}, err
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	); err != nil {
		return "", User{}, ErrInvalidCredentials
	}
	token, err := s.token.CreateToken(user.ID, user.Role)
	if err != nil {
		return "", User{}, fmt.Errorf(
			"creating token: %w", err,
		)
	}
	return token, user, nil
}

// GetUser returns the user with the given ID. Returns
// db.ErrNotFound if the user does not exist.
func (s *Service) GetUser(
	ctx context.Context,
	id int64,
) (User, error) {
	return s.repo.FindByID(ctx, id)
}
