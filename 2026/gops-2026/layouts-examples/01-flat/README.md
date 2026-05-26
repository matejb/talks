# Layout #1 - Flat / No Structure

> Also known as: "Simple Layout", "The Weekend Project Layout"
> Origin: Peter Bourgon, "Go in Production" (GopherCon 2014)

Everything lives in `package main`. No subpackages, no layers, no abstractions.

## Directory Structure

```shell
.
├── main.go             # Entry point - wires services and starts HTTP server
├── model.go            # Domain types (User, Book, Review)
├── user_service.go     # In-memory user store
├── book_service.go     # In-memory book store + borrow/return logic
├── review_service.go   # In-memory review store
├── handler.go          # HTTP handlers (Go 1.22+ pattern routing)
├── main_test.go        # Tests
└── go.mod
```

## Run

```bash
go run .
```

Server starts on `:8080`.

## Test

```bash
go test -v ./...
```

## Key Characteristics

| Aspect | Detail |
|--------|--------|
| Package | Everything is `package main` |
| Abstractions | None - concrete types everywhere |
| Testability | Limited - can't swap implementations |
| Dependency direction | All in one package, no direction to enforce |

## When This Works

- Prototypes and hackathons
- Small CLI tools (under ~10 files)
- Learning projects
- When you genuinely don't know the domain yet

## When This Breaks

- More than ~10 files - hard to find things
- Multiple developers - merge conflicts on every file
- Need for test isolation - can't swap implementations
- Reusable logic - everything is `package main`, can't import from outside

## Related

- Next step up: [02-ben-johnson/](../02-ben-johnson/) - adds dependency inversion with a root package
