# Project: Pet Store

## Overview
A pet store application with a Go backend API and a React frontend.

## Tech Stack
- **Backend:** Go standard library generated with ogen (https://ogen.dev/)
- **Frontend:** React with TypeScript
- **Database:** PostgreSQL
- **Build:** Go modules for backend, npm/Vite for frontend, mise for software and tasks

## Project Structure
```
/cmd/
  /server/      # Go server entrypoint
  /client/      # Go client CLI entrypoint
/client/        # ogen-generated Go API client (DO NOT EDIT)
/internal/      # Private Go packages (handlers, models, services)
  /api/         # OpenAPI spec, ogen configs, generated server code
/frontend/      # React application
  /src/
    /components/ # UI components (shadcn/ui + custom)
    /pages/      # Page components (one per route)
    /hooks/      # TanStack Query hooks (usePets, etc.)
    /services/   # API client layer (fetch wrappers)
    /types/      # Shared TypeScript types
```

## Build & Run
- **Backend:** `mise run api`
- **Frontend:** `mise run ui`
- **Generate (ogen):** `mise run generate`
- **Tests (Go):** `mise run api:test`
- **Tests (React):** `mise run ui:test`
- **Tests (All):** `mise run test`
- **Lint (Go):** `mise run go:lint`
- **Lint (React):** `mise run ui:lint`
- **Lint (All):** `mise run lint`

## Code Conventions

### Go
- Follow standard Go conventions (gofmt, go vet)
- Use structured logging (slog)
- Error handling: return errors, don't panic. Wrap errors with `fmt.Errorf("context: %w", err)`
- Use table-driven tests
- Keep handlers thin — business logic belongs in service layer
- Use context for cancellation and request-scoped values

### React / TypeScript
- Functional components only, use hooks
- Auth state: React Context + useReducer
- Server state (pets): TanStack Query v5
- Forms: React Hook Form + Zod
- Toasts: Sonner (call `toast()` directly, no hooks)
- Components: shadcn/ui for dialogs/dropdowns; hand-write
  simple components with Tailwind
- CSS: Tailwind CSS v4 (Vite plugin, `@theme` in CSS,
  no tailwind.config.js)
- Use named exports, not default exports
- API calls go through `frontend/src/services/` consumed
  by TanStack Query hooks in `frontend/src/hooks/`

### General
- No secrets or credentials in code — use environment variables
- Write tests for new functionality
- Keep PRs focused — one concern per change
- Wrap all markdown files at 80 characters per line
  (tables are exempt)
- Use Conventional Commits (e.g., `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`)
- Follow 12-Factor App principles (https://12factor.net/) — config via env vars, stateless processes, port binding, disposability, dev/prod parity, etc.
- Update `docs/REQUIREMENTS.md` and `docs/DESIGN.md` when making changes to features or implementation

## API Design
- RESTful endpoints under `/api/v1/`
- JSON request/response bodies
- Standard HTTP status codes (200, 201, 400, 404, 500)
- Consistent error response format: `{"error": "message"}`

## Environment Variables
- `ADDRESS` — listen address (default: `:8080`)
- `PETSTORE_USER` — database application user (required)
- `PETSTORE_PASSWORD` — database application user password
- `DB_HOST` — database host (default: `localhost`)
- `DB_PORT` — database port (default: `5432`)
- `DB_SSL_ENABLE` — set to `true` to require SSL
  (default: disabled)
- `FRONTEND_URL` — frontend origin for CORS
- `POSTGRES_PASSWORD` — password for the postgres superuser
- `JWT_SECRET` — JWT signing key (min 32 bytes)

### Secrets Management
- Secrets (`POSTGRES_PASSWORD`, `PETSTORE_PASSWORD`,
  `JWT_SECRET`) are stored in
  `.config/mise/mise.local.toml` (gitignored) using mise's
  age encryption — never in plaintext files or version
  control
- Access secrets via `mise env` at runtime
