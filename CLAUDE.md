# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based REST API for library management built with Gin framework, GORM ORM, SQLite database, and Swagger documentation. The application manages books and readers with full CRUD operations, in-memory caching, and runtime-configurable endpoint access control.

## Development Commands

### Running the Application
```bash
go run main.go
```
The server starts on `localhost:8080`. The SQLite database file `library.db` is automatically created in the project root.

### Installing Dependencies
```bash
go mod download
```

### Generating Swagger Documentation
After modifying handler annotations, regenerate Swagger docs:
```bash
swag init
```
This updates files in the `docs/` directory.

### Accessing Swagger UI
Navigate to `http://localhost:8080/swagger` while the server is running.

## Architecture

### Dependency Injection via Container Pattern

The `container` package implements a dependency injection container that initializes and wires up all application dependencies:

- **Database**: SQLite via GORM with auto-migration for models
- **Cache**: In-memory cache with configurable TTL
- **Repositories**: Data access layer with cache integration
- **Validator**: go-playground/validator for DTO validation
- **Config**: JSON-based runtime configuration

The container is instantiated in `main.go` and passed to handlers.

### Layered Architecture

1. **Models** (`models/`): GORM database entities (Book, Reader) with gorm.Model embedding (provides ID, CreatedAt, UpdatedAt, DeletedAt)

2. **DTOs** (`dto/`): Data Transfer Objects for API requests/responses
   - CreateDTO: Input validation for POST requests
   - UpdateDTO: Input validation for PUT requests
   - ResponseDTO: Clean API responses (no GORM fields)

3. **Repository** (`repository/`): Data access interfaces and implementations
   - Implements caching strategy: cache on read, invalidate on write
   - Both single-item and list caching with pattern-based invalidation

4. **Handlers** (`handlers/`): HTTP request handlers with Swagger annotations
   - Endpoint availability controlled by runtime config (config.json)
   - Each handler checks `config.Enable*` flags and returns 403 if disabled
   - Validation errors formatted with detailed field-level messages

5. **Validation** (`validation/`): Centralized validation using go-playground/validator
   - Custom error formatting for user-friendly API responses

6. **Cache** (`cache/`): Thread-safe in-memory cache with TTL
   - Helper functions generate standardized cache keys (e.g., "books:list", "books:id:5")
   - Pattern-based invalidation for bulk operations

7. **Config** (`config/`): JSON configuration loader with defaults
   - Enables/disables individual endpoints at runtime without recompilation
   - Configures cache TTL in seconds

### Configuration System

Edit `config.json` to enable/disable endpoints without code changes:
```json
{
  "cache_ttl_seconds": 300,
  "enable_get_books": false,    // Disables GET /books and GET /books/:id
  "enable_post_books": true,
  "enable_put_books": true,
  "enable_delete_books": true,
  "enable_get_readers": true,
  "enable_post_readers": true,
  "enable_put_readers": true,
  "enable_delete_readers": true
}
```

If config.json is missing or invalid, the application falls back to `DefaultConfig()` with all endpoints enabled.

### Cache Invalidation Strategy

- **Create**: Invalidates list cache (e.g., "books:list")
- **Update**: Invalidates both item cache (e.g., "books:id:5") and list cache
- **Delete**: Invalidates both item cache and list cache
- **DeleteAll**: Uses pattern-based invalidation to clear all related entries

### Module Name

The Go module is named `lab1` in go.mod. All imports use `lab1/` prefix (e.g., `lab1/models`, `lab1/handlers`).

## API Endpoints

All endpoints follow RESTful conventions:

### Books
- `GET /books/` - List all books
- `POST /books/` - Create a book
- `DELETE /books/` - Delete all books
- `GET /books/:id` - Get book by ID
- `PUT /books/:id` - Update book by ID
- `DELETE /books/:id` - Delete book by ID

### Readers
Same pattern as books, under `/readers/` route group.

### Swagger
- `GET /swagger` - Redirects to Swagger UI
- `GET /swagger/*any` - Serves Swagger UI and documentation

## Key Implementation Patterns

### Handler Pattern
Handlers receive injected dependencies (repository, validator, config) and follow this flow:
1. Check if endpoint is enabled via config
2. Parse and validate request
3. Call repository method
4. Convert model to DTO for response

### Repository Pattern
Repositories encapsulate database operations and caching:
- Interface-based for testability
- Cache-first reads with fallback to database
- Automatic cache invalidation on writes

### Error Handling
- 400: Invalid input (malformed JSON or validation errors)
- 403: Endpoint disabled via config
- 404: Resource not found
- 500: Internal server errors

Validation errors return structured JSON with field-level details.
