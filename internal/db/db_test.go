package db

import (
	"context"
	"testing"
)

func TestNewMissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := New(context.Background())
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
	want := "DATABASE_URL is required"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestNewBadDatabaseURL(t *testing.T) {
	t.Setenv(
		"DATABASE_URL",
		"postgres://invalid:5432/nonexistent?connect_timeout=1",
	)

	_, err := New(context.Background())
	if err == nil {
		t.Fatal("expected error for bad database URL")
	}
}
