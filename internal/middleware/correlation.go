package middleware

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

const correlationHeader = "X-Correlation-ID"

// correlationKey is the context key for the correlation ID.
type correlationKey struct{}

// CorrelationID returns middleware that reads or generates a
// correlation ID for each request. An incoming
// X-Correlation-ID header is reused; otherwise a new ULID
// is generated. The ID is stored in the request context and
// set on the response header.
func CorrelationID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			id := r.Header.Get(correlationHeader)
			if id == "" {
				id = ulid.Make().String()
			}
			ctx := context.WithValue(
				r.Context(), correlationKey{}, id,
			)
			w.Header().Set(correlationHeader, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetCorrelationID returns the correlation ID from the
// context, or an empty string if none is set.
func GetCorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationKey{}).(string)
	return id
}
