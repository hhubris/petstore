package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned when a token cannot be parsed
// or fails validation.
var ErrInvalidToken = errors.New("invalid token")

// Claims holds the application-level claims extracted from
// a validated JWT.
type Claims struct {
	UserID int64
	Role   string
}

// TokenConfig holds the signing key and expiry duration used
// to create and parse JWTs.
type TokenConfig struct {
	signingKey []byte
	expiry     time.Duration
	// timeNow is used for testing; defaults to time.Now.
	timeNow func() time.Time
}

// NewTokenConfig validates the secret and returns a
// TokenConfig with a default 1-hour expiry.
func NewTokenConfig(secret []byte) (*TokenConfig, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf(
			"jwt secret must be at least 32 bytes, got %d",
			len(secret),
		)
	}
	return &TokenConfig{
		signingKey: secret,
		expiry:     time.Hour,
		timeNow:    time.Now,
	}, nil
}

// CreateToken signs a JWT containing the given user ID and
// role. The token uses HS256 and includes sub, role, iat,
// and exp claims.
func (tc *TokenConfig) CreateToken(
	userID int64,
	role string,
) (string, error) {
	now := tc.timeNow()
	claims := jwt.MapClaims{
		"sub":  strconv.FormatInt(userID, 10),
		"role": role,
		"iat":  jwt.NewNumericDate(now),
		"exp":  jwt.NewNumericDate(now.Add(tc.expiry)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(tc.signingKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}

// ParseToken validates the token string and returns the
// extracted claims. It checks the signing method, signature,
// and expiry.
func (tc *TokenConfig) ParseToken(
	tokenString string,
) (Claims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					t.Header["alg"],
				)
			}
			return tc.signingKey, nil
		},
		jwt.WithTimeFunc(tc.timeNow),
	)
	if err != nil {
		return Claims{}, fmt.Errorf(
			"%w: %w", ErrInvalidToken, err,
		)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	sub, err := mapClaims.GetSubject()
	if err != nil {
		return Claims{}, fmt.Errorf(
			"%w: missing sub claim", ErrInvalidToken,
		)
	}
	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return Claims{}, fmt.Errorf(
			"%w: invalid sub claim: %w", ErrInvalidToken, err,
		)
	}

	role, ok := mapClaims["role"].(string)
	if !ok {
		return Claims{}, fmt.Errorf(
			"%w: missing role claim", ErrInvalidToken,
		)
	}

	return Claims{UserID: userID, Role: role}, nil
}
