# Layout #4 - Peter Bourgon Complex Layout (`pkg/`)

> Origin: Peter Bourgon, "Go Best Practices 2016"
> The `pkg/` convention - everything importable goes under `pkg/`.

Very similar to Ben Johnson's layout but with an explicit `pkg/` namespace.
Popular in older monorepos where multiple binaries share packages.

## Directory Structure

```shell
.
├── cmd/
│   └── server/
│       └── main.go              # Composition root
├── pkg/
│   ├── library/
│   │   ├── user.go              # User type + UserService interface
│   │   ├── book.go              # Book type + BookService interface
│   │   └── review.go            # Review type + ReviewService interface
│   ├── memory/
│   │   ├── user_service.go      # In-memory UserService
│   │   ├── book_service.go      # In-memory BookService
│   │   ├── review_service.go    # In-memory ReviewService
│   │   └── service_test.go      # Tests
│   └── httphandler/
│       └── handler.go           # HTTP handlers
└── go.mod
```

> **Note:** The HTTP handler package is named `httphandler` (not `http`) to avoid collision
> with `net/http` in imports. This is a practical consideration when using `pkg/http/`.

## Dependency Direction

```shell
  pkg/httphandler/  ──▶  pkg/library/  ◀──  pkg/memory/
```

`pkg/library/` has **zero** external imports.

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
| Package | Domain types and interfaces in `pkg/library/` |
| Abstractions | Interfaces in `pkg/library/`, implementations in `pkg/memory/` |
| Testability | Swap `pkg/memory/` for any implementation |
| Naming | `pkg/` signals "this is public, reusable code" |

## The `pkg/` Debate (2026)

- **Go modules** made `pkg/` semantically meaningless - packages are already namespaced by module path
- Single-binary applications don't need the extra level
- Peter Bourgon himself has moved toward simpler layouts
- Still useful in **monorepos** to distinguish shared library code from application code

## When This Works

- Monorepos with shared packages
- When multiple binaries need the same domain logic
- Teams accustomed to the `pkg/` convention

## When This Breaks

- Single-binary applications - `pkg/` is unnecessary
- Modern Go modules make `pkg/` redundant
- Can confuse newcomers

## Related

- Without `pkg/`: [02-ben-johnson/](../02-ben-johnson/) - same structure at module root
- Full monorepo: [08-monorepo/](../08-monorepo/) - go.work with multiple modules
