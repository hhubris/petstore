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
    /components/
    /pages/
    /hooks/
    /services/  # API client code
    /types/
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
- State management: React Context + useReducer
- Use named exports, not default exports
- CSS approach: Tailwind CSS
- API calls go through `frontend/src/services/` — components don't call fetch directly

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
- `PORT` — server port (default: 8080)
- `DATABASE_URL` — database connection string
- `FRONTEND_URL` — frontend origin for CORS
- `POSTGRES_PASSWORD` — password for the postgres superuser
- `PETSTORE_PASSWORD` — password for the petstore application database user
- `JWT_SECRET` — JWT signing key (min 32 bytes)

### Secrets Management
- Secrets (`POSTGRES_PASSWORD`, `PETSTORE_PASSWORD`, `JWT_SECRET`) are stored in `.config/mise/mise.local.toml` (gitignored) using mise's age encryption — never in plaintext files or version control
- Access secrets via `mise env` at runtime
