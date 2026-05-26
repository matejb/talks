# Go Project Structures - Layout Examples

Runnable examples of **8 Go project layouts** for an office library application, demonstrating how the same domain is organized differently in each approach.

Companion to the talk _"Go Project Structures - 2026 Edition"_.

## Layouts

| # | Layout | Key Idea |
|---|--------|-----------|
| 1 | [Flat / No Structure](01-flat/) | Everything in `package main` |
| 2 | [Ben Johnson Standard Package](02-ben-johnson/) | Domain in root package, implementations in leaf packages |
| 3 | [William Kennedy / Ardan Labs](03-kennedy/) | `internal/` with domain + platform layers |
| 4 | [Peter Bourgon - `pkg/`](04-bourgon-pkg/) | Everything importable under `pkg/` |
| 5 | [golang-standards/project-layout](05-standard-layout/) | Comprehensive conventional directories |
| 6 | [Feature-Based / Domain-Driven](06-feature-based/) | Organize by feature, not by layer |
| 7 | [Hexagonal / Ports & Adapters](07-hexagonal/) | Core has zero dependencies; everything else is an adapter |
| 8 | [Monorepo / Go Workspaces](08-monorepo/) | Multiple modules tied together with `go.work` |

## The Example Domain

Every layout implements the same **office library** application:

- **Users** browse books, borrow and return them, and write reviews
- **Books** can be borrowed by one user at a time
- **Reviews** have a text comment and a 1–5 rating

### Domain Types

```go
type User    struct { ID int; Name string; Email string }
type Book    struct { ID int; Title string; Author string; Borrowed *User }
type Review  struct { ID int; BookID int; UserID int; Text string; Rating int }
```

### API Endpoints

All layouts expose the same HTTP API (Go 1.22+ pattern routing):

| Method | Pattern | Description |
|--------|---------|-------------|
| GET | `/users` | List all users |
| GET | `/users/{id}` | Get user by ID |
| POST | `/users` | Create a user |
| GET | `/books` | List all books |
| GET | `/books/{id}` | Get book by ID |
| POST | `/books/{id}/borrow` | Borrow a book (`{"user_id": 1}`) |
| POST | `/books/{id}/return` | Return a book |
| GET | `/books/{id}/reviews` | List reviews for a book |
| POST | `/reviews` | Create a review |

## Quick Start

Each layout is self-contained. Pick one and run it:

```bash
# Flat
cd 01-flat && go run .

# Any structured layout
cd 02-ben-johnson && go run ./cmd/server/

# Monorepo (uses go.work)
cd 08-monorepo && go run ./services/api/cmd/api/
```

## Running Tests

```bash
# Single layout
cd 03-kennedy && go test ./... -v

# All layouts (non-monorepo)
for d in 01-flat 02-ben-johnson 03-kennedy 04-bourgon-pkg 05-standard-layout 06-feature-based 07-hexagonal; do
  (cd "$d" && go test ./... -v && echo "✅ $d")
done

# Monorepo
cd 08-monorepo && go test -v ./libs/library/... ./services/api/internal/handler/...
```
