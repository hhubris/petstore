package auth

import "time"

// User is the domain model for a user account. It contains
// fields like PasswordHash and timestamps that must not
// leak to the API layer.
type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
