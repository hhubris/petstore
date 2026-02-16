package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseCapture wraps http.ResponseWriter to capture the
// status code written by downstream handlers.
type responseCapture struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rc *responseCapture) WriteHeader(code int) {
	if !rc.wroteHeader {
		rc.status = code
		rc.wroteHeader = true
	}
	rc.ResponseWriter.WriteHeader(code)
}

func (rc *responseCapture) Write(b []byte) (int, error) {
	if !rc.wroteHeader {
		rc.status = http.StatusOK
		rc.wroteHeader = true
	}
	return rc.ResponseWriter.Write(b)
}

// Logging returns middleware that logs each request after
// it completes. Requests resulting in status < 500 are
// logged at Info; 5xx responses are logged at Error.
func Logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			rc := &responseCapture{
				ResponseWriter: w,
				status:         http.StatusOK,
			}
			start := time.Now()

			next.ServeHTTP(rc, r)

			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rc.status),
				slog.Duration("duration", time.Since(start)),
				slog.String(
					"correlation_id",
					GetCorrelationID(r.Context()),
				),
			}

			if rc.status >= 500 {
				slog.LogAttrs(
					r.Context(), slog.LevelError,
					"request completed", attrs...,
				)
			} else {
				slog.LogAttrs(
					r.Context(), slog.LevelInfo,
					"request completed", attrs...,
				)
			}
		})
	}
}
