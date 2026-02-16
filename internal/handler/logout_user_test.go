package handler_test

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestLogoutUser(t *testing.T) {
	h := newHandler(t, nil, nil)
	rec := httptest.NewRecorder()
	ctx := ctxWithResponseWriter(rec)

	err := h.LogoutUser(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cookies := rec.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "access_token" {
			found = true
			if c.MaxAge != -1 {
				t.Errorf("got MaxAge %d, want -1", c.MaxAge)
			}
		}
	}
	if !found {
		t.Error("access_token cookie not cleared")
	}
}

func TestLogoutUser_NoResponseWriter(t *testing.T) {
	h := newHandler(t, nil, nil)
	err := h.LogoutUser(context.Background())
	if err == nil {
		t.Fatal("expected error when response writer missing")
	}
}
