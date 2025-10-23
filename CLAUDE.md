# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Mindscape** is a full-stack web application built using the **Egg Framework** - a Go-based framework designed for rapid development by solo developers. The project follows a clean architecture pattern with clear separation between frontend, backend, and data layers.

### Tech Stack
- **Backend**: Go with Echo framework, PostgreSQL, Redis, MinIO
- **Frontend**: SolidJS with TypeScript, TailwindCSS 4.x, Rsbuild
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Authentication**: JWT tokens
- **File Storage**: MinIO (S3-compatible)
- **Documentation**: Swagger/OpenAPI

## Backend Architecture

### Core Components

1. **CLI Commands (`cmd/`)**
   - Cobra-based CLI with commands for migration, swagger generation, version management
   - Main entry point serves the application on configurable port (default: 60000)
   - Environment-specific configuration loading

2. **Controllers (`controllers/`)**
   - Echo-based REST API controllers
   - Route registration with middleware (logging, recovery, request ID)
   - Controllers: User, Folder, Bookmark
   - Static file serving for assets and public files
   - Health check endpoint (`/api/_health`)

3. **Services Layer (`services/`)**
   - Interface-driven design with implementations and mocks
   - Services: Auth, User, Folder, Bookmark, Notes, Redis, MinIO, Validator
   - Clean separation of business logic from controllers

4. **Models (`models/`)**
   - **Requests**: Structured input validation models
   - **Responses**: API response models
   - **Handlers**: Business logic handlers with IHandler interface pattern
   - Handler categories: user_handlers, folder_handlers, bookmark_handlers

5. **Database Layer (`db/`)**
   - **Migrations**: SQL migration files for schema management  
   - **Queries**: SQLC-generated type-safe database queries
   - **Repository**: Auto-generated database access layer from SQLC

6. **Middleware (`middlewares/`)**
   - Custom middleware configurations
   - Authentication and authorization handling

### Key Features
- User management with admin capabilities
- Folder/bookmark organization system
- Profile picture upload/download
- JWT-based authentication
- SQL injection prevention
- Input validation and sanitization

## Frontend Architecture

### Tech Stack Details
- **SolidJS**: Reactive UI framework with TypeScript
- **Rsbuild**: Fast build tool with Solid plugin
- **TailwindCSS 4.x**: Utility-first CSS framework
- **Biome**: Code formatting and linting
- **Auto-generated API Client**: TypeScript client from OpenAPI spec

### Project Structure
```
web/
├── src/
│   ├── api/           # Auto-generated API client
│   │   ├── apis/      # API endpoint classes
│   │   ├── models/    # TypeScript interfaces
│   │   └── runtime.ts # Base API runtime
│   ├── App.tsx        # Main application
│   └── index.tsx      # Entry point
├── package.json       # Dependencies and scripts
└── rsbuild.config.ts  # Build configuration
```

### API Integration
- Backend proxy: `http://0.0.0.0:60000`
- JWT authorization for authenticated endpoints
- CORS-enabled development server
- Path alias `@/*` maps to `./src/*`

## Testing Architecture

### Framework & Location
- **Framework**: Go's built-in testing with testify/assert
- **Test Files**: 
  - `services/ValidatorService_test.go` - Input validation testing
  - `services/UserService_test.go` - Service layer testing with MockUserService
  - `models/handlers/user_handlers/*_test.go` - Handler layer testing with functional structure
- **Purpose**: Comprehensive validation, service layer, and handler layer testing

### Testing Pattern Structure
All test files follow this consistent structured pattern:

```go
// <method var=services.ServiceName.MethodName>
// <fixtures/>     - Test data builders and mutations
// <runners/>      - Test execution functions  
// <tests>
// <evaluators>    - Assertion functions
// <map/>         - Test case organization
// <hook/>        - Main test runner
// </method>
```

#### Key Testing Components
1. **Fixtures**: Base parameter builders (`NewUserRequestParams()`) and mutation functions (`WithBadUsername()`, `WithDuplicateEmail()`)
2. **Runners**: Execute specific test scenarios (`Run_Create_ValidRequest`, `Run_Login_WrongPassword`)
3. **Evaluators**: Assert expected outcomes using testify (`EvaluateCreateSuccess`, `EvaluateLoginFailure`)
4. **Maps**: Organize test cases for dynamic execution (`CreateTestMap`, `LoginTestMap`)
5. **Hooks**: Main test functions that iterate through test maps (`Test_MockUserService_Create`)

#### Security Testing
- SQL injection prevention tests
- Input sanitization validation  
- Password complexity requirements (7+ chars, upper/lower/number/special)

#### Current Test Coverage

**ValidatorService_test.go** (Input Validation)
- `ValidateNewUserRequest` - User registration validation (14 test cases)
- `ValidateLoginFormRequest` - Login credential validation (7 test cases)
- `ValidateUpdateUserCredentialRequest` - User update validation (15 test cases)
- `ValidateCreateFolderRequest` - Folder creation validation (14 test cases)
- `CreateBookmarkRequest` - Bookmark creation validation (17 test cases)

**UserService_test.go** (Service Layer with MockUserService)
- `Create` - User creation with validation and duplication checks (6 test cases)
- `Login` - Authentication with email/username and password verification (7 test cases)
- `Get` - User retrieval by ID (3 test cases)
- `Remove` - User deletion with verification (3 test cases)
- `GetAll` - Bulk user retrieval (2 test cases)
- `Update` - Profile picture updates (3 test cases)
- `UpdateUserCredentials` - Credential changes with validation (7 test cases)

**Handler Layer Tests** (User Handlers with functional structure)
- `LoginHandler_test.go` - User authentication testing (11 test cases)
- `RegisterHandler_test.go` - User registration testing (18 test cases)
- `UpdateUserHandler_test.go` - User credential updates testing (8 test cases)
- `GetUsersHandler_test.go` - Admin user listing testing (4 test cases)
- `DeleteUserHandler_test.go` - Admin user deletion testing (7 test cases)
- `GetCurrentLoggedInUserHandler_test.go` - Current user retrieval testing (3 test cases)

**Folder Handler Tests** (Handler Layer Integration Testing)
- `CreateFolderHandler_test.go` - Folder creation with validation, authentication, and service coordination (7 test cases)
- `GetRootHandler_test.go` - Root folder retrieval with multi-service coordination (6 test cases)
- `GetByIDHandler_test.go` - Folder retrieval by ID with authorization and service failures (8 test cases)
- `DeleteFolderHandler_test.go` - Folder deletion with ownership verification (6 test cases)

**Total Test Coverage**: 176 test cases across all test files

### Handler Layer Testing

The handler layer tests utilize the same functional testing structure as service tests, providing comprehensive coverage for HTTP request handling, authentication, authorization, and business logic validation.

#### Key Handler Testing Features
- **JWT Authentication Testing**: Proper JWT token setup and context creation for authenticated endpoints
- **HTTP Context Simulation**: Echo framework context with proper request/response setup
- **Embedded Utility Functions**: Reusable test utilities embedded in each test file to avoid package conflicts
- **Mock Service Integration**: Utilizes MockUserService and MockAuthService with failure injection capabilities
- **Edge Case Coverage**: Tests for malformed JSON, invalid parameters, non-existent users, and permission issues

#### Handler Test Structure
Each handler test follows the functional pattern:
```go
// Test fixtures and data builders
func HandlerRequestParams() map[string]any { ... }
func WithBadRequest(req map[string]any) map[string]any { ... }

// Test execution runners  
func Run_Handler_ValidRequest(t *testing.T) { ... }
func Run_Handler_ValidationFailure(t *testing.T) { ... }

// Test case organization
var HandlerTestMap = map[string]func(*testing.T){
    "ValidRequest": Run_Handler_ValidRequest,
    "ValidationFailure": Run_Handler_ValidationFailure,
}

// Main test runner
func Test_Handler(t *testing.T) {
    for name, testFunc := range HandlerTestMap {
        t.Run(name, testFunc)
    }
}
```

#### Bug Fixes During Testing
Several critical bugs were discovered and fixed during handler testing implementation:

1. **UpdateUserHandler.go**: Fixed nil pointer dereference in authorization failure
2. **GetUsersHandler.go**: Fixed nil pointer dereference in admin privilege check
3. **DeleteUserHandler.go**: Fixed multiple critical issues:
   - Incorrect variable references in error checking
   - Missing return statements for error conditions
   - Nil pointer dereference in authorization logic
   - SetError method parameter bug

### MockUserService Architecture

The `MockUserService` provides a complete in-memory implementation of `IUserService` for database-independent testing:

#### Key Features
- **Thread-Safe Operations**: All data operations protected by `sync.RWMutex`
- **In-Memory Storage**: Three concurrent maps for efficient lookups:
  - `users` - ID-based user storage
  - `usersByEmail` - Email-based user lookup
  - `usersByUsername` - Username-based user lookup
- **Behavioral Controls**: Test failure flags (`ShouldFailCreate`, `ShouldFailLogin`, etc.)
- **Call Tracking**: Captures method invocations and parameters for verification
- **Realistic Business Logic**: 
  - bcrypt password hashing and verification
  - Duplicate email/username detection
  - Basic field validation
- **Seeded Test Data**: Pre-populated with test and admin users using known credentials
- **Reset Functionality**: `Reset()` method for clean test isolation

#### Usage Pattern
```go
// Setup
service := CreateMockUserService(context.Background(), nil)
service.Reset()

// Configure behavior for testing
service.ShouldFailCreate = true
service.CreateErrorMessage = "Custom error message"

// Execute and verify
user, err := service.Create(params)
assert.Error(t, err)
assert.Equal(t, 1, service.CreateCallCount)
```

#### Test Data
- **Regular User**: email: `test@example.com`, username: `testuser`, password: `password123`
- **Admin User**: email: `admin@example.com`, username: `adminuser`, password: `password123`
- **Handler Test Users**: Additional test users for handler layer testing with complex passwords
  - Handler User: email: `handler@example.com`, username: `handleruser`, password: `passwordABC123!`
  - Delete Test Users: Various users for delete handler testing scenarios

## Development Commands

### Backend
```bash
# Development with hot reload
make dev
air

# Build and run
make build
make dev-backend

# Testing
make test
make test-coverage
make test-race

# Database migrations
go run main.go migrate up
go run main.go migrate down

# Swagger documentation
make swagger
go run main.go swag
```

### Frontend
```bash
# Development server (http://localhost:5173)
pnpm dev

# Build and preview
pnpm build
pnpm preview

# Code quality
pnpm format
pnpm check
```

### Testing Commands
```bash
# All tests
go test ./...

# Run with coverage report
make test-coverage

# Run with race detection
make test-race

# Service layer tests
go test ./services -v

# Specific validation tests
go test ./services -run Test_Validate -v

# Mock service tests  
go test ./services -run Test_MockUserService -v

# Handler layer tests
go test ./models/handlers/user_handlers -v
go test ./models/handlers/folder_handlers -v
go test ./models/handlers/folder_handlers -run Test_CreateFolderHandler -v
go test ./models/handlers/folder_handlers -run Test_GetRootHandler -v

# Test with coverage (use Makefile commands)
make test-coverage
make test-race
```

### Database Operations
```bash
# Run migrations up
go run main.go migrate up

# Run migrations down  
go run main.go migrate down

# Generate database code from queries (after schema changes)
sqlc generate
```

## SQLC Database Code Generation

This project uses SQLC for type-safe database queries. Key configuration:

- **Schema location**: `db/migrations/` (SQL migration files)
- **Queries location**: `db/queries/` (SQL query files) 
- **Generated code**: `db/repository/` (Go structs and query functions)
- **Configuration**: `sqlc.yml` with PostgreSQL engine and pgx/v5 driver

After making schema or query changes, regenerate code with `sqlc generate`.

## Egg Framework Context

Mindscape is built on the **Egg Framework** - a Go framework optimized for solo developers to rapidly build full-stack applications. Key characteristics:

- **Rapid Development**: Pre-configured with common patterns and dependencies
- **Solo Developer Focus**: Minimal boilerplate, fast iteration cycles  
- **Hot Reload**: Uses `air` for automatic rebuilds during development
- **Docker Ready**: Multiple Dockerfiles for different deployment scenarios
- **CI/CD Ready**: Coolify deployment configuration included

## Prerequisites & Setup

- **PostgreSQL database** running with connection string configured in `config/development.yaml`
- **S3-compatible storage** (MinIO) configured in development config
- **Go 1.24.5+** and **pnpm** for frontend development
- Run `air` for hot reload development or `go run main.go` for standard server start

## Configuration

- **Environment-specific configs**: `config/development.yaml`, `config/production.yaml`
- **Database**: PostgreSQL connection configuration with SQLC for type-safe queries
- **Storage**: MinIO/S3 configuration for file uploads
- **Cache**: Redis configuration for session management
- **Server**: Configurable host/port (default: 0.0.0.0:60000)
- **JWT**: Authentication token configuration
- **API Generation**: Auto-generated TypeScript client in `web/src/api/`

## Development Workflow

### Database Operations
```bash
# Run database migrations
go run main.go migrate up
go run main.go migrate down

# Generate SQLC code from SQL queries
sqlc generate
```

### Frontend API Client Generation
The frontend API client is auto-generated from OpenAPI specs. After backend changes:
```bash
# Generate new Swagger docs
make swagger
go run main.go swag

# Frontend will automatically regenerate API client on next build
cd web && pnpm build
```

## Deployment

- Docker support with multiple Dockerfiles (Local, Production)
- Docker Compose for local development
- CI/CD ready for Coolify deployment
- Makefile automation for build/test/deploy workflows
