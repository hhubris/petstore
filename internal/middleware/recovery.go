package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery returns middleware that recovers from panics,
// logs the error with a stack trace, and writes a 500 JSON
// response matching the ogen Error schema.
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			defer func() {
				if v := recover(); v != nil {
					slog.Error("panic recovered",
						"error", v,
						"stack", string(debug.Stack()),
					)
					w.Header().Set(
						"Content-Type",
						"application/json",
					)
					w.WriteHeader(
						http.StatusInternalServerError,
					)
					_ = json.NewEncoder(w).Encode(map[string]any{
						"code":    500,
						"message": "internal server error",
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
