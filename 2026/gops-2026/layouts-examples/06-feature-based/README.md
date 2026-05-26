# Layout #6 - Feature-Based / Domain-Driven

> Origin: Multiple sources - Alex Edwards (2023), community consensus (2024–2026)
> Arguably the **most popular** layout for new Go web services in 2025–2026.

Organize by **feature**, not by layer. Each feature folder contains everything about that domain: types, repository, service, and HTTP handlers.

## Directory Structure

```shell
.
├── cmd/
│   └── server/
│       └── main.go              # Composition root
├── internal/
│   ├── user/
│   │   ├── model.go             # User type
│   │   ├── repository.go        # Repository interface + in-memory impl
│   │   ├── service.go           # Business logic
│   │   ├── handler.go           # HTTP handlers for /users
│   │   └── service_test.go      # Tests
│   ├── book/
│   │   ├── model.go             # Book type
│   │   ├── repository.go        # Repository interface + in-memory impl
│   │   ├── service.go           # Business logic (borrow/return)
│   │   ├── handler.go           # HTTP handlers for /books
│   │   └── service_test.go      # Tests
│   ├── review/
│   │   ├── model.go             # Review type
│   │   ├── repository.go        # Repository interface + in-memory impl
│   │   ├── service.go           # Business logic
│   │   ├── handler.go           # HTTP handlers for /reviews
│   │   └── service_test.go      # Tests
│   └── platform/
│       └── database/
│           └── database.go      # Shared database infrastructure
└── go.mod
```

## Feature Package Anatomy

Each feature folder follows the same pattern:

```shell
internal/<feature>/
├── model.go        # Domain type
├── repository.go   # Interface + in-memory implementation
├── service.go      # Business logic (depends on Repository interface)
├── handler.go      # HTTP handlers - Routes() map[string]http.HandlerFunc
└── service_test.go # Tests using in-memory repository
```

## Cross-Feature Dependencies

Cross-feature lookups (e.g., book service checking if a user exists) are handled via **function types** injected at the composition root, avoiding circular imports:

```go
// internal/book/service.go
type UserLookup func(id int) (user.User, error)

type Service struct {
    Repo        Repository
    LookupUser  UserLookup
}
```

## Run

```bash
go run ./cmd/server/
```

Server starts on `:8080`.

## Test

```bash
go test -v ./...
```

## Key Characteristics

| Aspect | Detail |
|--------|--------|
| Package | One package per feature (`user`, `book`, `review`) |
| Abstractions | Repository interface per feature |
| Testability | Each feature independently testable with its own in-memory repo |
| Cohesion | Everything about a feature lives in one folder |

## When This Works

- Medium-to-large services with multiple distinct features
- Teams where different developers own different features
- Each feature is independently understandable
- Microservice candidates - each feature folder could become its own service

## When This Breaks

- Cross-cutting concerns (e.g., "borrow book" touches user + book)
- Small applications - overhead of multiple packages for a few files
- Requires discipline to keep `platform/` clean

## Related

- With full dependency inversion: [07-hexagonal/](../07-hexagonal/) - core/adapter separation
- Simpler variant: [03-kennedy/](../03-kennedy/) - similar `internal/` structure without interfaces
