# Layout #5 - golang-standards/project-layout

> Origin: Community repo on GitHub (~56k stars)
> **Not** an official Go standard despite the name.

A comprehensive reference of conventional directories. Best treated as a **menu**, not a blueprint - use the directories you need, ignore the rest.

## Directory Structure

```shell
.
├── cmd/
│   └── library/
│       └── main.go              # Application entry point
├── internal/
│   └── app/
│       └── library/
│           ├── user.go          # User type + UserService
│           ├── book.go          # Book type + BookService (borrow/return)
│           ├── review.go        # Review type + ReviewService (with validation)
│           ├── handler.go       # HTTP handlers (Go 1.22+ routing)
│           ├── memory.go        # In-memory store implementations
│           └── service_test.go  # Tests
├── api/
│   └── README.md                # API documentation
├── configs/
│   └── default.json             # Configuration file template
└── go.mod
```

## Key Directories

| Directory | Purpose | In this example |
|-----------|---------|-----------------|
| [`cmd/`](cmd/) | Main applications | `cmd/library/main.go` |
| [`internal/`](internal/) | Private application code | `internal/app/library/` |
| [`api/`](api/) | API definitions | [`api/README.md`](api/README.md) |
| [`configs/`](configs/) | Configuration templates | [`configs/default.json`](configs/default.json) |

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
| Package | All app code in `internal/app/library/` |
| Abstractions | Store interfaces for testability |
| Testability | In-memory implementations swappable for real stores |
| Scope | Demonstrates `cmd/`, `internal/`, `api/`, `configs/` conventions |

## Critical Perspective

> "It is not an official Go standard, despite the name." - Hacker News (2021)

> "The layout can encourage cargo-cult programming." - kau.sh (2024)

> "Beginners see 15 directories and think they need all of them. You don't." - Reddit r/golang (2024)

**But also:** Individual directory conventions are fine on their own. When you DO need `/api/`, `/configs/`, etc., this shows where they conventionally go.

## When This Works

- Large projects that actually need many conventional directories
- Teams needing a shared vocabulary for directory placement
- Reference for "where would X go?"

## When This Breaks

- Small-to-medium projects (most Go projects)
- When treated as a template to copy wholesale
- When the team doesn't understand why each directory exists

## Related

- Pragmatic alternative: [06-feature-based/](../06-feature-based/) - organizes by feature instead
