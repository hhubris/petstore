package auth

import "context"

// claimsKey is the unexported context key for storing
// Claims. Using a struct type avoids collisions with keys
// from other packages.
type claimsKey struct{}

// ContextWithClaims returns a child context carrying the
// given Claims.
func ContextWithClaims(
	ctx context.Context, c Claims,
) context.Context {
	return context.WithValue(ctx, claimsKey{}, c)
}

// ClaimsFromContext extracts Claims from the context.
// Returns the Claims and true if present, or zero Claims
// and false if not.
func ClaimsFromContext(
	ctx context.Context,
) (Claims, bool) {
	c, ok := ctx.Value(claimsKey{}).(Claims)
	return c, ok
}
