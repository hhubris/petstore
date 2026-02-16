package db

import (
	"context"
	"testing"
)

func TestNewMissingPetstoreUser(t *testing.T) {
	t.Setenv("PETSTORE_USER", "")

	_, err := New(context.Background())
	if err == nil {
		t.Fatal("expected error for missing PETSTORE_USER")
	}
	want := "PETSTORE_USER is required"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestNewMissingPetstorePassword(t *testing.T) {
	t.Setenv("PETSTORE_USER", "petstore")
	t.Setenv("PETSTORE_PASSWORD", "")

	_, err := New(context.Background())
	if err == nil {
		t.Fatal(
			"expected error for missing PETSTORE_PASSWORD",
		)
	}
	want := "PETSTORE_PASSWORD is required"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestNewBadConnection(t *testing.T) {
	t.Setenv("PETSTORE_USER", "petstore")
	t.Setenv("PETSTORE_PASSWORD", "secret")
	t.Setenv("DB_HOST", "invalid")
	t.Setenv("DB_PORT", "5432")

	_, err := New(context.Background())
	if err == nil {
		t.Fatal("expected error for bad connection")
	}
}
