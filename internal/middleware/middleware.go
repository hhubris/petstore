package middleware

import "net/http"

// Middleware is an HTTP middleware function.
type Middleware func(http.Handler) http.Handler

// Chain applies middleware in order so the first argument
// is the outermost wrapper. For example:
//
//	Chain(h, A, B, C)
//
// produces A(B(C(h))).
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
