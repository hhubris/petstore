package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hhubris/petstore/internal/middleware"
)

func TestChain(t *testing.T) {
	tests := []struct {
		name string
		mws  []middleware.Middleware
		want string
	}{
		{
			name: "empty chain",
			mws:  nil,
			want: "handler",
		},
		{
			name: "single middleware",
			mws: []middleware.Middleware{
				tagMiddleware("A"),
			},
			want: "A-handler",
		},
		{
			name: "ordering first is outermost",
			mws: []middleware.Middleware{
				tagMiddleware("A"),
				tagMiddleware("B"),
				tagMiddleware("C"),
			},
			want: "A-B-C-handler",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("handler"))
				},
			)
			chain := middleware.Chain(h, tc.mws...)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet, "/", nil,
			)
			chain.ServeHTTP(rec, req)

			if got := rec.Body.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// tagMiddleware writes a tag before calling next.
func tagMiddleware(
	tag string,
) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tag + "-"))
				next.ServeHTTP(w, r)
			},
		)
	}
}
