# Layout #7 - Hexagonal Architecture / Ports & Adapters

> Origin: Alistair Cockburn (2005), adapted for Go (2022–2026)
> Also known as: Ports & Adapters

The domain core has **ZERO external dependencies**. Everything else is an adapter.
The killer feature: test business logic with in-memory adapters - no database needed.

## Directory Structure

```shell
.
├── cmd/
│   └── server/
│       └── main.go                    # Composition root - wires adapters to core
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   │   ├── user.go                # User entity
│   │   │   ├── book.go                # Book entity
│   │   │   └── review.go              # Review entity
│   │   ├── port/
│   │   │   ├── user_repository.go     # Outbound port (interface)
│   │   │   ├── book_repository.go     # Outbound port (interface)
│   │   │   └── review_repository.go   # Outbound port (interface)
│   │   └── service/
│   │       ├── user_service.go        # User use cases
│   │       ├── book_service.go        # Book use cases (Borrow, Return)
│   │       ├── review_service.go      # Review use cases
│   │       └── service_test.go        # Tests - uses memory adapters!
│   └── adapter/
│       ├── in/
│       │   └── http/
│       │       ├── user_handler.go    # Driving adapter (inbound)
│       │       ├── book_handler.go    # Driving adapter (inbound)
│       │       └── review_handler.go  # Driving adapter (inbound)
│       └── out/
│           └── memory/
│               ├── user_repository.go # Driven adapter (outbound)
│               ├── book_repository.go # Driven adapter (outbound)
│               └── review_repository.go
└── go.mod
```

## Architecture

```shell
               ┌─────────────────────┐
               │   adapter/in/http   │   Driving (inbound) adapter
               └─────────┬───────────┘
                         │ depends on
                         ▼
               ┌─────────────────────┐
               │     core/domain     │   "Inside" - no knowledge
               │     core/port       │   of outside technologies
               │     core/service    │
               └─────────▲───────────┘
                         │ implements
               ┌─────────┴───────────┐
               │  adapter/out/memory │   Driven (outbound) adapter
               └─────────────────────┘
```

**`internal/core/` has ZERO imports from `internal/adapter/`.** That's the whole point.

## Run

```bash
go run ./cmd/server/
```

Server starts on `:8080`.

## Test

```bash
go test -v ./...
```

## The Killer Feature: Testing Without a Database

```go
func TestBorrowBook(t *testing.T) {
    // Use in-memory adapters - no database, no HTTP server needed
    bookRepo := memory.NewBookRepository()
    userRepo := memory.NewUserRepository()
    svc := service.NewBookService(bookRepo, userRepo)

    bookRepo.Save(context.Background(), domain.Book{ID: 1, Title: "Go 101"})
    err := svc.Borrow(context.Background(), 1, 42)
    // ✅ Test pure business logic in isolation
}
```

## Key Characteristics

| Aspect | Detail |
|--------|--------|
| Package | `domain`, `port`, `service` in core; `http`, `memory` in adapters |
| Abstractions | Ports (interfaces) in `core/port/` - adapters implement them |
| Testability | ⭐⭐⭐⭐ - swap any adapter without touching core |
| Dependency direction | Adapters depend on core; core depends on nothing |

Compile-time interface checks ensure adapters satisfy ports:

```go
var _ port.BookRepository = (*BookRepository)(nil)
```

## When This Works

- Services with clear business rules protected from infrastructure
- Multiple implementations of the same port (Postgres, in-memory, gRPC)
- APIs that might change transport (HTTP → gRPC → CLI)
- Teams where different people own different adapters

## When This Breaks

- CRUD-only services - abstraction overhead isn't worth it
- Small teams - extra indirection slows development
- Risk of "architecture astronaut" behavior

## Related

- Simpler alternative: [06-feature-based/](../06-feature-based/) - same feature cohesion, less ceremony
- Foundation: [02-ben-johnson/](../02-ben-johnson/) - same dependency inversion principle
