# CLAUDE.md

## Project Overview

**Mindscape** is a full-stack web application built with the **Egg Framework** (Go-based rapid development framework).

### Tech Stack
- **Backend**: Go + Echo, PostgreSQL, Redis, MinIO, JWT auth
- **Frontend**: SolidJS + TypeScript, TailwindCSS 4.x, Rsbuild
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Docs**: Swagger/OpenAPI with auto-generated TypeScript client

## Architecture

```
├── cmd/              # Cobra CLI (migrate, swagger, serve)
├── controllers/      # Echo REST controllers (User, Folder, Bookmark, Widget)
├── services/         # Business logic (Auth, User, Folder, Bookmark, Redis, MinIO, Validator)
├── models/
│   ├── requests/     # Input validation models
│   ├── responses/    # API response models
│   └── handlers/     # IHandler pattern (user_handlers, folder_handlers, bookmark_handlers, widget_handlers)
├── db/
│   ├── migrations/   # SQL schema files
│   ├── queries/      # SQLC query files
│   └── repository/   # Auto-generated from SQLC
├── middlewares/      # Auth, logging, recovery
└── web/              # SolidJS frontend
    └── src/
        ├── api/      # Auto-generated API client from OpenAPI
        ├── components/
        ├── contexts/
        ├── hooks/
        └── pages/
```

## Commands

```bash
# Backend
make dev              # Hot reload with air
make build            # Build binary
make test             # Run all tests
make test-coverage    # Tests with coverage
make swagger          # Generate Swagger docs
go run main.go migrate up/down  # Database migrations
sqlc generate         # Regenerate DB code after schema changes

# Frontend (in web/)
pnpm dev              # Dev server at localhost:5173
pnpm build            # Production build
pnpm check            # Biome lint + format
```

## Testing

Tests use Go's testing package with testify/assert. Pattern:
- **Fixtures**: `NewUserRequestParams()`, `WithBadUsername()`
- **Runners**: `Run_Create_ValidRequest()`
- **Evaluators**: `EvaluateCreateSuccess()`
- **Maps**: Test case organization for iteration

Mock services (`MockUserService`, `MockAuthService`) provide in-memory implementations with failure injection for testing.

```bash
go test ./services -v                    # Service tests
go test ./models/handlers/... -v         # Handler tests
```

## Configuration

- Config files: `config/development.yaml`, `config/production.yaml`
- Server default: `0.0.0.0:60000`
- Frontend proxy configured to backend
- Path alias: `@/*` → `./src/*`

## Workflow

1. Backend changes → `make swagger` → regenerates OpenAPI spec
2. Frontend build → auto-generates TypeScript client from spec
3. Database schema changes → update migrations → `sqlc generate`
