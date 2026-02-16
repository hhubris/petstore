package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hhubris/petstore/internal/middleware"
)

func TestCorrelationID(t *testing.T) {
	tests := []struct {
		name      string
		headerVal string
		wantReuse bool
	}{
		{
			name:      "no header generates ULID",
			headerVal: "",
			wantReuse: false,
		},
		{
			name:      "existing header is reused",
			headerVal: "my-custom-id-123",
			wantReuse: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ctxID string
			inner := http.HandlerFunc(
				func(_ http.ResponseWriter, r *http.Request) {
					ctxID = middleware.GetCorrelationID(
						r.Context(),
					)
				},
			)

			mw := middleware.CorrelationID()
			h := mw(inner)

			req := httptest.NewRequest(
				http.MethodGet, "/", nil,
			)
			if tc.headerVal != "" {
				req.Header.Set(
					"X-Correlation-ID", tc.headerVal,
				)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			respID := rec.Header().Get("X-Correlation-ID")

			if respID == "" {
				t.Fatal("response missing X-Correlation-ID")
			}
			if ctxID == "" {
				t.Fatal("context missing correlation ID")
			}
			if respID != ctxID {
				t.Errorf(
					"response ID %q != context ID %q",
					respID, ctxID,
				)
			}

			if tc.wantReuse {
				if respID != tc.headerVal {
					t.Errorf(
						"got %q, want reused %q",
						respID, tc.headerVal,
					)
				}
			} else {
				// Should be a valid ULID (26 chars).
				if len(respID) != 26 {
					t.Errorf(
						"ULID length = %d, want 26",
						len(respID),
					)
				}
			}
		})
	}
}

func TestGetCorrelationIDEmpty(t *testing.T) {
	// Context without correlation ID returns empty string.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	id := middleware.GetCorrelationID(req.Context())
	if id != "" {
		t.Errorf("got %q, want empty string", id)
	}
}
