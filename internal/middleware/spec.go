package middleware

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/hhubris/petstore/internal/api"
)

// Spec returns middleware that serves the OpenAPI spec and
// Swagger UI documentation. Requests to /docs serve the
// Swagger UI and /docs/openapi.yml serves the raw spec.
// All other requests pass through to the next handler.
//
// If api.Spec() returns nil (binary built with
// -tags=disable_spec), the middleware is a no-op
// passthrough.
func Spec() Middleware {
	return func(next http.Handler) http.Handler {
		spec := api.Spec()
		if spec == nil {
			return next
		}

		mux := http.NewServeMux()

		mux.HandleFunc(
			"GET /docs/openapi.yml",
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/x-yaml",
				)
				_, _ = w.Write(spec)
			},
		)

		ui := middleware.SwaggerUI(middleware.SwaggerUIOpts{
			SpecURL: "/docs/openapi.yml",
			Path:    "docs",
		}, next)

		mux.Handle("GET /docs", ui)
		mux.Handle("/", next)

		return mux
	}
}
