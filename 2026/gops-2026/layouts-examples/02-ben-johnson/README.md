# Layout #2 - Ben Johnson Standard Package Layout

> Origin: Ben Johnson, "Standard Package Layout" (2016)
> The single most influential Go project layout.

Domain types and interfaces live in the **root package** (`package library`).
Implementations live in **leaf packages** (`memory/`, `http/`).
The root package depends on **NOTHING**.

## Directory Structure

```shell
.
├── cmd/
│   └── server/
│       └── main.go          # Composition root - wires everything together
├── model.go                  # Domain types (root package, package library)
├── service.go                # Service interfaces (root package, package library)
├── http/
│   └── handler.go            # HTTP handlers - depends on root package
├── memory/
│   ├── store.go              # In-memory implementations - depends on root package
│   └── store_test.go         # Tests for memory implementations
└── go.mod
```

## Dependency Direction

```shell
  http/  ──depends on──▶  (root)  ◀──depends on──  memory/
```

**The root package has ZERO external imports.**

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
| Package | Root package is `package library` - the domain |
| Abstractions | Interfaces in root package, implementations in leaf packages |
| Testability | Swap `memory/` for any implementation that satisfies the interfaces |
| Dependency direction | Leaf packages depend on root; root depends on nothing |

## When This Works

- Small-to-medium web services
- When you want testability via interfaces
- When you want domain separated from infrastructure
- The "next step up" from flat layout

## When This Breaks

- Root package can get large
- No `internal/` - all packages are importable externally
- Multiple domains compete in the same root package
- `mock/` package grows unbounded

## Related

- Simpler alternative: [01-flat/](../01-flat/) - no structure at all
- Enterprise variant: [03-kennedy/](../03-kennedy/) - adds `internal/` and platform layer
