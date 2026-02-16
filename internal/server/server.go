package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hhubris/petstore/internal/api"
	"github.com/hhubris/petstore/internal/auth"
	"github.com/hhubris/petstore/internal/db"
	"github.com/hhubris/petstore/internal/handler"
	"github.com/hhubris/petstore/internal/pet"
)

// shutdownTimeout is the maximum time to wait for in-flight
// requests to complete during graceful shutdown.
const shutdownTimeout = 10 * time.Second

// Run is the public entry point for the server. It reads
// configuration from environment variables, wires up all
// dependencies, and starts the HTTP server. It blocks until
// ctx is cancelled, then performs a graceful shutdown.
func Run(ctx context.Context) error {
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	env := os.Getenv("ENVIRONMENT")
	secure := env != "development"

	database, err := db.New(ctx)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer database.Close()

	h, err := build(database.DBTX(), jwtSecret, secure)
	if err != nil {
		return fmt.Errorf("building server: %w", err)
	}

	slog.Info("server starting", "addr", addr)

	if err := serve(ctx, addr, h); err != nil {
		return fmt.Errorf("serving: %w", err)
	}

	slog.Info("server stopped")
	return nil
}

// build wires up all dependencies and returns an
// http.Handler ready to serve requests.
func build(
	dbtx db.DBTX,
	jwtSecret string,
	secure bool,
) (http.Handler, error) {
	tc, err := auth.NewTokenConfig([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf(
			"creating token config: %w", err,
		)
	}

	userRepo := auth.NewUserRepository(dbtx)
	authSvc := auth.NewService(userRepo, tc)
	secHandler := auth.NewSecurityHandler(tc)

	petRepo := pet.NewPetRepository(dbtx)
	petSvc := pet.NewService(petRepo)

	h := handler.New(petSvc, authSvc, secure)

	srv, err := api.NewServer(h, secHandler)
	if err != nil {
		return nil, fmt.Errorf(
			"creating ogen server: %w", err,
		)
	}

	return handler.WrapWithResponseWriter(srv), nil
}

// serve starts an HTTP server and blocks until ctx is
// cancelled, then gracefully shuts down.
func serve(
	ctx context.Context, addr string, h http.Handler,
) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), shutdownTimeout,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	return nil
}
