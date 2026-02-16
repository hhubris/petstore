package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hhubris/petstore/internal/middleware"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		wantLevel string
	}{
		{
			name:      "200 logs at info",
			status:    http.StatusOK,
			wantLevel: "INFO",
		},
		{
			name:      "404 logs at info",
			status:    http.StatusNotFound,
			wantLevel: "INFO",
		},
		{
			name:      "500 logs at error",
			status:    http.StatusInternalServerError,
			wantLevel: "ERROR",
		},
		{
			name:      "503 logs at error",
			status:    http.StatusServiceUnavailable,
			wantLevel: "ERROR",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(
				slog.NewJSONHandler(&buf, nil),
			)
			slog.SetDefault(logger)
			t.Cleanup(func() {
				slog.SetDefault(
					slog.New(slog.NewTextHandler(
						&bytes.Buffer{}, nil,
					)),
				)
			})

			inner := http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(tc.status)
				},
			)

			mw := middleware.Logging()
			h := mw(inner)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet, "/test", nil,
			)
			h.ServeHTTP(rec, req)

			var entry map[string]any
			if err := json.Unmarshal(
				buf.Bytes(), &entry,
			); err != nil {
				t.Fatalf(
					"failed to parse log: %v\nraw: %s",
					err, buf.String(),
				)
			}

			if entry["level"] != tc.wantLevel {
				t.Errorf(
					"level = %v, want %s",
					entry["level"], tc.wantLevel,
				)
			}
			if entry["method"] != "GET" {
				t.Errorf(
					"method = %v, want GET",
					entry["method"],
				)
			}
			if entry["path"] != "/test" {
				t.Errorf(
					"path = %v, want /test",
					entry["path"],
				)
			}
			if entry["status"] != float64(tc.status) {
				t.Errorf(
					"status = %v, want %d",
					entry["status"], tc.status,
				)
			}
			if _, ok := entry["duration"]; !ok {
				t.Error("missing duration attribute")
			}
		})
	}
}
