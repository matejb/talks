# Layout #8 - Monorepo / Go Workspaces

> Origin: Go 1.18 workspaces (2022), matured through Go 1.23–1.24
> `go.work` was the biggest structural addition to Go since modules.

Multiple independent modules tied together with a `go.work` file for local development.
Each service and library has its own `go.mod`.

## Directory Structure

```shell
.
├── go.work                         # Workspace - ties modules together
├── libs/
│   ├── library/
│   │   ├── go.mod                  # Shared domain library
│   │   ├── user.go                 # User type
│   │   ├── book.go                 # Book type
│   │   ├── review.go               # Review type
│   │   ├── store.go                # In-memory store with CRUD + Borrow/Return
│   │   └── store_test.go           # Tests
│   └── dbutil/
│       ├── go.mod                  # Database utility library
│       └── dbutil.go               # DB helpers
├── services/
│   ├── api/
│   │   ├── go.mod                  # HTTP API service (depends on libs/library)
│   │   ├── cmd/api/main.go         # API entrypoint
│   │   └── internal/handler/
│   │       ├── handler.go          # HTTP handlers (Go 1.22+ routing)
│   │       └── handler_test.go     # Handler tests
│   └── worker/
│       ├── go.mod                  # Background worker (depends on libs/library)
│       └── cmd/worker/main.go      # Worker CLI - prints borrowed books report
└── tools/                          # (placeholder for dev tools)
```

## go.work

```go
go 1.24

use (
    ./libs/library
    ./libs/dbutil
    ./services/api
    ./services/worker
)
```

## Run

```bash
# API service
go run ./services/api/cmd/api/

# Worker
go run ./services/worker/cmd/worker/
```

## Test

```bash
# From the monorepo root:
go test ./libs/library/... ./services/api/internal/handler/...
```

## Module Dependency Graph

```
  services/api/  ──▶  libs/library/
  services/worker/ ──▶ libs/library/
  services/api/  ──▶  libs/dbutil/
```

Each service's `go.mod` references the shared libraries. `go.work` resolves them to local checkouts during development.

## Key Characteristics

| Aspect | Detail |
|--------|--------|
| Modules | Each service and library has its own `go.mod` |
| Workspace | `go.work` ties them together for local dev |
| Sharing | Domain types in `libs/library/` shared across services |
| Independence | Each module can be versioned and released separately |

## When This Works

- Multiple deployable services sharing domain logic
- Organization with 3+ Go services
- Different services with different release cycles
- Shared types (gRPC protos, domain models) across services

## When This Breaks

- Single service - monorepo adds complexity for no benefit
- Very large organizations - multi-repo with published modules may be better
- Teams can't agree on shared library versions

## Related

- Single-service alternative: [06-feature-based/](../06-feature-based/) or [07-hexagonal/](../07-hexagonal/)
- `pkg/` variant: [04-bourgon-pkg/](../04-bourgon-pkg/) - shared packages without go.work
