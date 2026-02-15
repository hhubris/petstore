# Requirements: Pet Store

## Purpose

A reference implementation demonstrating best practices for a
full-stack application with a Go backend and React frontend,
following 12-Factor App principles.

## Tech Stack

- **Backend:** Go standard library, API generated with
  [ogen](https://ogen.dev/) from OpenAPI 3.0 spec
- **Go Client:** Generated Go HTTP client (ogen) from the
  same OpenAPI spec; CLI entrypoint in `cmd/client/`
- **Frontend:** React + TypeScript, Vite, Tailwind CSS
- **Database:** PostgreSQL
- **Tooling:** mise (task runner, software management,
  encrypted secrets)

## API Specification

The API is defined in `internal/api/api.yml` (OpenAPI 3.0) and
serves as the single source of truth. Both server and client
code are generated using ogen from this spec.

### Operations

| Operation      | Method | Path             | Description              |
|----------------|--------|------------------|--------------------------|
| findPets       | GET    | /pets            | List pets, filter/limit  |
| addPet         | POST   | /pets            | Create a new pet         |
| find pet by id | GET    | /pets/{id}       | Get a single pet by ID   |
| deletePet      | DELETE | /pets/{id}       | Delete a pet by ID       |
| registerUser   | POST   | /auth/register   | Register a new user      |
| loginUser      | POST   | /auth/login      | Log in, set cookie       |
| logoutUser     | POST   | /auth/logout     | Log out, clear cookie    |
| getCurrentUser | GET    | /auth/me         | Get current user details |

### Data Models

- **Pet:** `id` (int64, required), `name` (string, required),
  `tag` (string, optional)
- **NewPet:** `name` (string, required),
  `tag` (string, optional)
- **Error:** `code` (int32, required),
  `message` (string, required)
- **RegisterRequest:** `name` (string, required),
  `email` (string, email format, required),
  `password` (string, min 8 / max 72, required)
- **LoginRequest:** `email` (string, email format, required),
  `password` (string, required)
- **AuthUser:** `id` (int64, required),
  `name` (string, required), `email` (string, required),
  `role` (enum: admin | customer, required)

### Response Behavior

- Successful list returns `200` with JSON array of Pet
- Successful create returns `200` with the created Pet
- Successful get returns `200` with a single Pet
- Successful delete returns `204` with no body
- Successful register returns `201` with AuthUser
- Successful login returns `200` with AuthUser and sets
  `access_token` cookie
- Successful logout returns `204` and clears the
  `access_token` cookie
- Successful get current user returns `200` with AuthUser
- Duplicate email on register returns `409`
- Invalid credentials on login returns `401`
- All errors return the Error schema with an appropriate
  HTTP status

## Authentication & Authorization

### Roles

- **admin** — full CRUD access (create, read, delete pets)
- **customer** — read-only access (list pets, view pet by ID)

### Auth Mechanism: JWT via HttpOnly Cookie

- **Algorithm:** HMAC-SHA256 (HS256)
- **Library:** `github.com/golang-jwt/jwt/v5`
- **Claims:** `sub` (user ID), `role`, `exp` (1 hour), `iat`
- **Signing key:** `JWT_SECRET` env var (min 32 bytes),
  stored in encrypted mise config
- **Delivery:** HttpOnly cookie (`access_token`)
- **No refresh tokens** — users re-login after 1hr expiry

### Cookie Configuration

| Property | Value                           |
|----------|---------------------------------|
| Name     | `access_token`                  |
| HttpOnly | `true`                          |
| Secure   | `true` (configurable for dev)   |
| SameSite | `Strict`                        |
| Path     | `/`                             |
| MaxAge   | `3600` (1 hour)                 |

### CSRF Protection

`SameSite=Strict` prevents cross-origin cookie sending.
Additionally, validate the `Origin` header on state-changing
requests (POST, DELETE) against `FRONTEND_URL`. No separate
CSRF token is needed.

### Auth Endpoints

| Operation        | Method | Path             | Auth |
|------------------|--------|------------------|------|
| `registerUser`   | POST   | `/auth/register` | No   |
| `loginUser`      | POST   | `/auth/login`    | No   |
| `logoutUser`     | POST   | `/auth/logout`   | Yes  |
| `getCurrentUser` | GET    | `/auth/me`       | Yes  |

### Auth Data Models

- **RegisterRequest:** `name` (string), `email` (string,
  email format), `password` (string, min 8 / max 72 chars
  — bcrypt limit)
- **LoginRequest:** `email` (string), `password` (string)
- **AuthUser:** `id` (int64), `name` (string),
  `email` (string), `role` (enum: admin, customer)

### Authorization Matrix

| Endpoint            | Public | Customer | Admin |
|---------------------|--------|----------|-------|
| GET /pets           | Yes    | Yes      | Yes   |
| GET /pets/{id}      | Yes    | Yes      | Yes   |
| POST /pets          | No     | No       | Yes   |
| DELETE /pets/{id}   | No     | No       | Yes   |
| POST /auth/register | Yes    | —        | —     |
| POST /auth/login    | Yes    | —        | —     |
| POST /auth/logout   | —      | Yes      | Yes   |
| GET /auth/me        | —      | Yes      | Yes   |

### Password Hashing

- **Library:** `golang.org/x/crypto/bcrypt`, default cost
- Passwords are hashed before storage; plaintext is never
  persisted
- The 72-byte bcrypt input limit is enforced at the API
  schema level (`maxLength: 72` on password fields)

### Admin Account Creation

- New registrations always receive the `customer` role
- Admin accounts are created via seed script or direct
  database insert — no self-service admin promotion

### Users Table Schema (design reference)

```sql
CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    name          TEXT         NOT NULL,
    email         TEXT         NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          TEXT         NOT NULL DEFAULT 'customer'
                  CHECK (role IN ('admin', 'customer')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);
```

### Architecture

```
internal/
  db/
    db.go           # DBTX interface, sentinel errors
  auth/
    user.go         # User domain model (private fields)
    repository.go   # UserRepository (DB queries) ✓
    security.go     # ogen SecurityHandler (JWT validation) ✓
    service.go      # AuthService (register, login, get user) ✓
    authz.go        # RequireAdmin() helper ✓
    jwt.go          # Token creation and parsing ✓
    context.go      # Context key types, ClaimsFromContext() ✓
  pet/
    repository.go   # PetRepository (DB queries) ✓
    service.go      # PetService (CRUD logic) ✓
migrations/
  000001–000006     # Schema and role setup
```

Items marked ✓ are implemented; others are planned.

### Key Design Decisions

| Decision             | Choice             | Rationale                      |
|----------------------|--------------------|--------------------------------|
| GET /pets public     | No auth required   | Allows browsing without account|
| No refresh tokens    | 1hr access token   | Simpler; re-login on expiry    |
| SameSite=Strict      | No CSRF token      | Strongest browser protection   |
| Role in JWT claims   | Avoid DB lookup    | Role changes require re-login  |
| bcrypt default cost  | Standard, tested   | 72-byte limit in schema        |
| Admin creation       | Manual / seed      | No self-service admin promotion|

## Frontend Requirements

### Pages / Views

1. **Login** — Email and password form; redirects to Pet List
   on success
2. **Register** — Name, email, and password form; redirects
   to Login on success
3. **Pet List** — Browse all pets with tag filtering and
   pagination/limit
4. **Pet Detail** — View a single pet's information
5. **Add Pet** (admin only) — Form to create a new pet
6. **Delete confirmation** (admin only) — Confirm before
   deleting a pet

### Routing & Auth Guards

- Unauthenticated users can access Login, Register, Pet List,
  and Pet Detail
- Unauthenticated users attempting to access protected pages
  are redirected to Login
- After login, redirect to the originally requested page
  (or Pet List by default)
- Navigation bar shows current user name/role when logged in,
  and Login/Register links when logged out
- Logout clears the cookie (via POST /auth/logout) and
  redirects to Pet List

### UX Requirements

- Responsive layout (mobile + desktop)
- Loading and error states for all API calls
- Feedback on successful create/delete actions
  (toast or inline message)
- Admin-only actions are hidden or disabled for
  customer role
- Login/register forms show inline validation errors
  (invalid email, password too short, duplicate email,
  wrong credentials)

### State Management

- React Context + useReducer for auth/role state and UI state
- API calls go through a service layer
  (`frontend/src/services/`)

## Go Client Requirements

- Generated from the same `internal/api/api.yml` spec as
  the server using ogen
- Client code lives in the `client/` package at the
  project root
- CLI entrypoint in `cmd/client/main.go`
- Supports all API operations (pets CRUD, auth endpoints)
- Handles cookie-based authentication via ogen's
  `SecuritySource` interface
- Code generation via `mise run generate`
  (runs `go generate ./internal/api/...`)

### Code Generation

- Two ogen config files control generation:
  - `internal/api/ogen-server.yml` — server code →
    `internal/api/`
  - `internal/api/ogen-client.yml` — client code →
    `client/`
- Generation directives live in
  `internal/api/generate.go`
- Generated `oas_*.go` files must not be edited manually

## Database Requirements

- PostgreSQL with two database users:
  - `postgres` — superuser for admin/migration tasks
  - `petstore` — application user with limited privileges
    (SELECT, INSERT, UPDATE, DELETE on `pets` and `users`
    tables; USAGE on sequences)
- Passwords stored in encrypted mise config
  (never in plaintext or version control)
- Tables:
  - **pets:** `id` (bigserial primary key),
    `name` (text, not null), `tag` (text, nullable,
    indexed)
  - **users:** `id` (bigserial primary key),
    `name` (text, not null),
    `email` (text, not null, unique index),
    `password_hash` (text, not null),
    `role` (text, not null, default 'customer',
    check in ('admin', 'customer')),
    `created_at` (timestamptz, not null, default now()),
    `updated_at` (timestamptz, not null, default now())

### Migrations

- SQL migration files live in `migrations/` with separate
  `.up.sql` / `.down.sql` files, numbered sequentially
- Applied using `golang-migrate/migrate` with the
  PostgreSQL driver
- Table creation and index creation are kept in separate
  migrations
- Migration order:
  1. Create `petstore` role (password from env var via
     session variable)
  2. Create `pets` table
  3. Create `pets` indexes (`idx_pets_tag`)
  4. Create `users` table
  5. Create `users` indexes (`idx_users_email` unique)
  6. Grant privileges to `petstore` role
- A migration runner script (`scripts/migrate.sh`) wraps
  `golang-migrate` to inject the `petstore_password`
  session variable required by migration 000001
- The PostgreSQL container uses a named Docker volume
  (`petstore-data`) for data persistence across restarts
- Migrations run automatically on server startup in dev,
  and explicitly via CLI in production

## Environment Variables

| Variable           | Description                    |
|--------------------|--------------------------------|
| `PORT`             | Server port (default: 8080)    |
| `DATABASE_URL`     | PostgreSQL connection string   |
| `FRONTEND_URL`     | Frontend origin for CORS       |
| `POSTGRES_PASSWORD`| postgres superuser password    |
| `PETSTORE_PASSWORD`| petstore app user password     |
| `JWT_SECRET`       | JWT signing key (min 32 bytes) |

## Non-Functional Requirements

- Follow 12-Factor App principles
- Conventional Commits for all git history
- Structured logging with slog
- Table-driven tests for Go code
- Dev/prod parity — containerization-ready
- API spec (`api.yml`) is the single source of truth;
  code is generated, not hand-written

### Documentation

- All feature changes must update `docs/REQUIREMENTS.md`
- All implementation changes must update `docs/DESIGN.md`
- API changes require updating `internal/api/api.yml`
  before implementation

### Security

- Rate limiting on `/auth/login` and `/auth/register`
  to mitigate brute-force and credential-stuffing attacks
  (e.g., 10 requests per minute per IP)
- Log authentication events at INFO level: successful
  login, failed login (without password), registration,
  and logout
- Never log passwords, tokens, or password hashes
