# Pet Store

A full-stack pet store application demonstrating best
practices with a Go backend, React frontend, and
PostgreSQL database, following
[12-Factor App](https://12factor.net/) principles.

## Tech Stack

- **Backend:** Go with [ogen](https://ogen.dev/)-generated
  server from OpenAPI 3.0
- **Go Client:** ogen-generated HTTP client from the same
  spec
- **Frontend:** React, TypeScript, Vite, Tailwind CSS
- **Database:** PostgreSQL
- **Tooling:** [mise](https://mise.jdx.dev/) for task
  running, tool management, and secrets

## Prerequisites

- [mise](https://mise.jdx.dev/) (installs Go and other
  tools automatically)
- PostgreSQL
- Node.js / npm (for frontend)

## Quick Start

```bash
# Install tools
mise install

# Generate server and client code from OpenAPI spec
mise run generate

# Run the backend
mise run api

# Run the frontend (in a separate terminal)
mise run ui
```

## Project Structure

```
cmd/
  server/         # Go server entrypoint
  client/         # Go client CLI entrypoint
client/           # ogen-generated Go API client
internal/
  api/            # OpenAPI spec, ogen configs, generated
                  #   server code
  handler/        # Server handler implementations
  auth/           # JWT auth, security, user repository
  pet/            # Pet service and repository
frontend/         # React application
  src/
    components/
    pages/
    hooks/
    services/     # API client layer
migrations/       # SQL migration files
docs/
  REQUIREMENTS.md # What the system does
  DESIGN.md       # How it's built
```

## Common Tasks

| Task                  | Command              |
|-----------------------|----------------------|
| Generate API code     | `mise run generate`  |
| Run backend           | `mise run api`       |
| Run frontend          | `mise run ui`        |
| Test Go               | `mise run api:test`  |
| Test React            | `mise run ui:test`   |
| Test all              | `mise run test`      |
| Lint Go               | `mise run go:lint`   |
| Lint React            | `mise run ui:lint`   |
| Lint all              | `mise run lint`      |

## API

The API is defined in `internal/api/api.yml` (OpenAPI 3.0)
and serves as the single source of truth. Both server and
client code are generated from this spec using ogen.

Endpoints are served under `/api/v1/` with JSON
request/response bodies. See
[docs/REQUIREMENTS.md](docs/REQUIREMENTS.md) for the full
operation list and authorization matrix.

## Environment Variables

| Variable            | Description                      |
|---------------------|----------------------------------|
| `PORT`              | Server port (default: 8080)      |
| `DATABASE_URL`      | PostgreSQL connection string     |
| `FRONTEND_URL`      | Frontend origin for CORS         |
| `POSTGRES_PASSWORD` | postgres superuser password      |
| `PETSTORE_PASSWORD` | petstore app user password       |
| `JWT_SECRET`        | JWT signing key (min 32 bytes)   |

Secrets are stored in `.config/mise/mise.local.toml`
(gitignored) using mise's age encryption — never in
plaintext files or version control.

## Documentation

- [Requirements](docs/REQUIREMENTS.md) — what the system
  does (features, API contract, auth rules)
- [Design](docs/DESIGN.md) — how it's built (architecture,
  patterns, decision log)
- [CLAUDE.md](CLAUDE.md) — conventions for AI-assisted
  development
