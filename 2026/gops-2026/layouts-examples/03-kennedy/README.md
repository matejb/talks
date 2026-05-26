# Layout #3 - William Kennedy / Ardan Labs Package-Oriented Design

> Origin: William Kennedy, "Package-Oriented Design" (2017)
> Very popular in enterprise Go.

Three layers: **domain packages**, **platform packages**, and a **composition root**.
All application code lives under `internal/`, enforced by the Go compiler.

## Directory Structure

```shell
.
├── cmd/
│   └── library/
│       └── main.go              # Composition root
├── internal/
│   ├── platform/
│   │   ├── database/
│   │   │   └── db.go            # Database wrapper / helper
│   │   └── http/
│   │       └── server.go        # HTTP server helper
│   ├── users/
│   │   ├── model.go             # User domain type
│   │   ├── service.go           # User service (in-memory store)
│   │   ├── handler.go           # HTTP handlers for /users
│   │   └── service_test.go      # Tests
│   ├── books/
│   │   ├── model.go             # Book domain type
│   │   ├── service.go           # Book service + borrow/return logic
│   │   ├── handler.go           # HTTP handlers for /books
│   │   └── service_test.go      # Tests
│   └── reviews/
│       ├── model.go             # Review domain type
│       ├── service.go           # Review service
│       ├── handler.go           # HTTP handlers for /reviews
│       └── service_test.go      # Tests
└── go.mod
```

## Three Layers

| Layer | Packages | Purpose |
|-------|----------|---------|
| **Domain** | `internal/users/`, `internal/books/`, `internal/reviews/` | Business logic, models, HTTP handlers |
| **Platform** | `internal/platform/` | Reusable infrastructure - database, HTTP helpers |
| **Composition root** | `cmd/library/` | Wires domain and platform together |

## Run

```bash
go run ./cmd/library/
```

Server starts on `:8080`.

## Test

```bash
go test -v ./...
```

## Key Characteristics

| Aspect | Detail |
|--------|--------|
| Package | Each domain has its own package under `internal/` |
| Abstractions | Minimal - services directly use concrete stores |
| Testability | Each domain package is independently testable |
| Encapsulation | `internal/` enforced by the Go compiler |

## When This Works

- Medium-to-large applications
- Multi-team, multi-repo setups - `internal/` enforces boundaries
- Enterprise / corporate Go

## When This Breaks

- `internal/platform/` can become a "junk drawer"
- Domain packages may directly import infrastructure - nothing prevents it
- No prescribed way to share libraries between repos
- Overkill for small services

## Related

- Simpler variant: [02-ben-johnson/](../02-ben-johnson/) - dependency inversion via interfaces
- Similar `pkg/` structure: [04-bourgon-pkg/](../04-bourgon-pkg/) - same idea with `pkg/` prefix
