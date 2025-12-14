# Library Management System

A full-stack web application for managing books and readers, built with Go (Gin) backend and vanilla JavaScript frontend.

## Tech Stack

**Backend**: Go 1.23+ | Gin | GORM | SQLite | JWT
**Frontend**: Vanilla JavaScript ES6+ (modules) | HTML5 | CSS3
**Testing**: Selenium WebDriver (Go)

## Quick Start

```bash
# Install dependencies
go mod download

# Run application
go run main.go
# Access at http://localhost:8080

# Run tests (requires ChromeDriver on port 4444)
go test ./tests/
```

## Project Structure

```
.
├── handlers/           # HTTP handlers (auth, books, readers)
├── models/            # Database models (User, Book, Reader)
├── repository/        # Data access layer
├── dto/              # Request/response structures
├── middleware/       # Auth middleware
├── container/        # Dependency injection
├── validation/       # Input validation
├── static/           # Frontend files
│   ├── js/          # Modular JavaScript
│   ├── index.html
│   └── style.css
├── tests/            # E2E Selenium tests
├── main.go          # Application entry point
└── config.json      # Configuration
```

## API Endpoints

**Authentication** (no auth required):
- `POST /auth/register` - Register user
- `POST /auth/login` - Login user

**Protected** (require Bearer token):
- `GET/POST/DELETE /books/` - Manage all books
- `GET/PUT/DELETE /books/:id` - Manage single book
- `GET/POST/DELETE /readers/` - Manage all readers
- `GET/PUT/DELETE /readers/:id` - Manage single reader
- `POST /readers/:id/books/:bookId` - Add book to reader's reading list
- `DELETE /readers/:id/books/:bookId` - Remove book from reader's reading list
- `GET /auth/profile` - Get user profile

## Architecture

See detailed documentation:
- [BACKEND.md](./BACKEND.md) - Backend architecture and patterns
- [FRONTEND.md](./FRONTEND.md) - Frontend module structure
- [TESTING.md](./TESTING.md) - Testing approach and test cases

## Key Features

- JWT authentication with bcrypt password hashing
- Role-based access control (admin vs general users)
- Book ownership - users can only edit/delete their own books (admins can edit all)
- CRUD for books (title, description, owner) and readers (name, surname)
- Readers can have "currently reading" lists (many-to-many with books)
- Client-side filtering, sorting, pagination
- Client-side CSV export
- Responsive UI with modal forms and custom confirmations
- Auto-seeded admin user (username: `admin`, password: `password`)
