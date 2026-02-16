package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hhubris/petstore/internal/middleware"
)

func TestRecovery(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.Handler
		wantStatus int
		wantJSON   bool
	}{
		{
			name: "no panic",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			),
			wantStatus: http.StatusOK,
			wantJSON:   false,
		},
		{
			name: "panic string",
			handler: http.HandlerFunc(
				func(_ http.ResponseWriter, _ *http.Request) {
					panic("something broke")
				},
			),
			wantStatus: http.StatusInternalServerError,
			wantJSON:   true,
		},
		{
			name: "panic error",
			handler: http.HandlerFunc(
				func(_ http.ResponseWriter, _ *http.Request) {
					panic(errors.New("error value"))
				},
			),
			wantStatus: http.StatusInternalServerError,
			wantJSON:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mw := middleware.Recovery()
			h := mw(tc.handler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet, "/", nil,
			)
			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf(
					"status = %d, want %d",
					rec.Code, tc.wantStatus,
				)
			}

			if tc.wantJSON {
				ct := rec.Header().Get("Content-Type")
				if ct != "application/json" {
					t.Errorf(
						"Content-Type = %q, want application/json",
						ct,
					)
				}

				var body map[string]any
				if err := json.NewDecoder(
					rec.Body,
				).Decode(&body); err != nil {
					t.Fatalf("invalid JSON body: %v", err)
				}
				if body["code"] != float64(500) {
					t.Errorf(
						"code = %v, want 500", body["code"],
					)
				}
				if body["message"] != "internal server error" {
					t.Errorf(
						"message = %v, want %q",
						body["message"],
						"internal server error",
					)
				}
			}
		})
	}
}
