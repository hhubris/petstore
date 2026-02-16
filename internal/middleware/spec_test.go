package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/middleware"
)

func TestSpec(t *testing.T) {
	passthrough := http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("passthrough"))
		},
	)

	mw := middleware.Spec()
	h := mw(passthrough)

	specEnabled := api.Spec() != nil

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantContain string
		wantCT      string
		needsSpec   bool
	}{
		{
			name:        "spec returns YAML",
			path:        "/docs/openapi.yml",
			wantStatus:  http.StatusOK,
			wantContain: "openapi",
			wantCT:      "application/x-yaml",
			needsSpec:   true,
		},
		{
			name:        "docs returns HTML",
			path:        "/docs",
			wantStatus:  http.StatusOK,
			wantContain: "<html",
			wantCT:      "text/html",
			needsSpec:   true,
		},
		{
			name:        "other paths pass through",
			path:        "/pets",
			wantStatus:  http.StatusOK,
			wantContain: "passthrough",
			wantCT:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.needsSpec && !specEnabled {
				t.Skip(
					"spec disabled via build tag",
				)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet, tc.path, nil,
			)
			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf(
					"status = %d, want %d",
					rec.Code, tc.wantStatus,
				)
			}

			body := rec.Body.String()
			if !strings.Contains(body, tc.wantContain) {
				t.Errorf(
					"body missing %q; got %s",
					tc.wantContain,
					body[:min(len(body), 200)],
				)
			}

			if tc.wantCT != "" {
				ct := rec.Header().Get("Content-Type")
				if !strings.Contains(ct, tc.wantCT) {
					t.Errorf(
						"Content-Type = %q, want %q",
						ct, tc.wantCT,
					)
				}
			}
		})
	}
}

func TestSpecDisabledPassthrough(t *testing.T) {
	if api.Spec() != nil {
		t.Skip("spec is enabled; testing bypass path")
	}

	passthrough := http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("passthrough"))
		},
	)

	mw := middleware.Spec()
	h := mw(passthrough)

	for _, path := range []string{
		"/docs", "/docs/openapi.yml", "/pets",
	} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet, path, nil,
			)
			h.ServeHTTP(rec, req)

			if body := rec.Body.String(); body != "passthrough" {
				t.Errorf(
					"got %q, want passthrough", body,
				)
			}
		})
	}
}
