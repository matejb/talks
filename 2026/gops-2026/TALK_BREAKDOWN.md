# Go Project Structures - 2026 Edition

## Detailed Talk Breakdown

**Target:** 60-minute talk, beginner-friendly with useful content for seasoned developers  
**Format:** Survey of options (not opinionated recommendation)  
**Example domain:** Office library - users browse books, borrow them, write reviews  
**Demo:** Each layout discussed will have a runnable code example  
**Slide format:** TBD later

---

## Section 0: Opening Hook (2 min)

**Purpose:** Grab attention immediately with a relatable pain point.

- Open with a screenshot of a confused developer staring at 47 empty folders in a new Go project
- Title card: "Go has no official project structure. So… what do I do?"
- Quick preview: "We'll look at 9 approaches - from zero structure to full hexagonal architecture - using the same example app so you can compare them side by side."

**Speaker notes:**
- Relate to the audience: everyone has asked "where do I put this file?" in Go
- Emphasize that Go's lack of a prescribed structure is intentional - it's about choosing what fits YOUR project

---

## Section 1: The Example Application (3 min)

**Purpose:** Establish the domain model that will be used in every layout demo.

### The Library App

A web application for a small office library:
- Users can browse book titles
- Users can borrow and return books
- Users can write reviews on books

### Domain Types (used across all layouts)

```go
// User represents a library member
type User struct {
    ID    int
    Name  string
    Email string
}

// Book represents a library book
type Book struct {
    ID       int
    Title    string
    Author   string
    Borrowed *User
}

// Review represents a book review
type Review struct {
    ID     int
    BookID int
    UserID int
    Text   string
    Rating int
}
```

### Key Operations

```go
// Core interfaces (will appear in different packages depending on layout)
type UserService interface {
    User(id int) (User, error)
    Users() ([]User, error)
    CreateUser(user User) error
}

type BookService interface {
    Book(id int) (Book, error)
    Books() ([]Book, error)
    Borrow(bookID int, userID int) error
    Return(bookID int) error
}

type ReviewService interface {
    Reviews(bookID int) ([]Review, error)
    CreateReview(review Review) error
}
```

**Speaker notes:**
- Keep this simple - it's the same domain from the 2019 version of this talk
- Point out that the same business logic will be organized differently in each layout
- The app is intentionally small enough to understand quickly but has enough pieces to show structural differences

---

## Section 2: Layout #1 - Flat / No Structure (5 min)

**Also known as:** "The Peter Bourgon Simple Layout", "The Weekend Project Layout"  
**Origin:** Peter Bourgon, "Go in Production" (2014, updated 2016)  
**2024 reinforcement:** LaurentSV, "No-Nonsense Go Package Layout" (Oct 2024)

### Directory Structure

```
.
├── main.go
├── user.go
├── book.go
├── review.go
├── user_service.go
├── book_service.go
├── review_service.go
├── handler.go
└── go.mod
```

### Key Code

```go
// main.go - everything in package main
package main

import (
    "database/sql"
    "log"
    "net/http"
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "connection string")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    us := &userService{DB: db}
    bs := &bookService{DB: db}
    rs := &reviewService{DB: db}

    h := handler{
        Users:   us,
        Books:   bs,
        Reviews: rs,
    }

    log.Fatal(http.ListenAndServe(":8080", h))
}
```

### When This Works

- Prototypes and hackathons
- Small CLI tools (under ~10 files)
- Learning projects and exercises
- When you genuinely don't know the domain yet

### When This Breaks

- More than ~10 files - hard to find things
- Multiple developers - merge conflicts on every file
- Need for test isolation - can't swap implementations
- Any reusable logic - everything is `package main`, can't import from outside

### 2026 Perspective

> "Don't structure for a future that may never arrive."  
> - LaurentSV, "No-Nonsense Go Package Layout" (2024)

> "Start with a very small structure, often just main.go plus go.mod, and add folders only when a real need appears."  
> - Reddit r/golang consensus (2024)

**Speaker notes:**
- This is NOT a "bad" layout - it's the right layout for the right situation
- Many popular Go tools (gofmt, staticcheck early versions) started flat
- The Go standard library itself uses this pattern in many packages
- LaurentSV's 2024 article is worth reading - argues that most Go projects never need more than this

---

## Section 3: Layout #2 - Ben Johnson's Standard Package Layout (6 min)

**Origin:** Ben Johnson, "Standard Package Layout" (2016)  
**URL:** https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

### Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── user.go            # domain types & interfaces (root package)
├── book.go
├── review.go
├── http/              # HTTP transport (depends on root package)
│   ├── handler.go
│   └── server.go
├── postgres/          # storage implementation (depends on root package)
│   ├── connection.go
│   ├── user_service.go
│   ├── book_service.go
│   └── review_service.go
└── mock/              # test doubles (depends on root package)
    ├── user.go
    ├── book.go
    └── review.go
```

> **Note:** The root package IS the domain. In Ben Johnson's original article, if the
> module is `github.com/myorg/library`, the root package is `package library` - containing
> all domain types and interfaces with zero external dependencies.

### Key Idea: Dependency Inversion

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  http/   │     │  (root)  │     │ postgres/ │
│ (handler)│────▶│ (domain) │◀────│ (storage) │
└──────────┘     └──────────┘     └──────────┘
     │                │                  │
     │    depends on  │   implements     │
     └────────────────┘─────────────────┘
```

**Both `http/` and `postgres/` depend on the root package (the domain).  
The root package depends on NOTHING.**

### Key Code

```go
// user.go - domain types and interfaces in root package, no external imports
package library

type User struct {
    ID    int
    Name  string
    Email string
}

type UserService interface {
    User(id int) (User, error)
    Users() ([]User, error)
    CreateUser(user User) error
}
```

```go
// postgres/user_service.go - implements the interface
package postgres

import (
    "database/sql"
    "github.com/myorg/library"
)

type UserService struct {
    DB *sql.DB
}

func (us *UserService) User(id int) (library.User, error) {
    // SQL query here
    return library.User{}, nil
}
```

```go
// cmd/server/main.go - composition root, wires everything together
package main

import (
    "log"
    "net/http"
    "github.com/myorg/library/http"
    "github.com/myorg/library/postgres"
)

func main() {
    db, _ := postgres.Open("connection")
    us := &postgres.UserService{DB: db}
    bs := &postgres.BookService{DB: db}

    h := &http.Handler{
        UserService: us,
        BookService: bs,
    }
    log.Fatal(http.ListenAndServe(":8080", h))
}
```

### When This Works

- Small-to-medium web services
- When you want testability via interfaces
- When you want domain logic separated from infrastructure
- The "next step up" from flat layout

### When This Breaks

- Root package (`library/`) can get large
- No `internal/` - all packages are importable by external code
- Multiple domains start competing in the same root package
- The `mock/` package grows unbounded as the system grows

### Historical Significance

- **This is the single most influential Go project layout.** Almost every modern layout is a variation of this pattern.
- The core insight - "domain in root, implementations in leaf packages, main wires them" - appears in hexagonal, clean architecture, and feature-based layouts.
- Ben Johnson first introduced this approach in "Structuring Applications in Go" (2014) and refined it in this 2016 article.

**Speaker notes:**
- Emphasize the dependency arrow diagram - this is THE key insight
- Show that the root package has zero imports from `http/` or `postgres/`
- This is still the most common "first structured layout" Go developers learn

---

## Section 4: Layout #3 - Ben Johnson's Original Article (5 min)

**Origin:** Ben Johnson, "Structuring Applications in Go" (2014)  
**URL:** https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091

### Directory Structure

```
.
├── book.go            # domain types in ROOT package
├── user.go
├── review.go
├── cmd/
│   └── server/
│       └── main.go
├── http/
│   ├── handler.go
│   └── server.go
├── postgres/
│   ├── connection.go
│   ├── user_service.go
│   ├── book_service.go
│   └── review_service.go
└── mock/
    ├── user.go
    └── book.go
```

### Relationship to Layout #2

> **⚠️ Important:** This article (2014) predates "Standard Package Layout" (2016). Both articles
> describe essentially the same core structure - domain types in the root package, implementations
> in subpackages. Consider merging this section with Section 3 for the actual talk.
>
> The 2016 "Standard Package Layout" article is the more refined and better-known version.
> This earlier article introduced the same concepts in a more exploratory way.

```go
// user.go - in the root package
package library

type User struct {
    ID    int
    Name  string
    Email string
}

type UserService interface {
    User(id int) (User, error)
    Users() ([]User, error)
    CreateUser(user User) error
}
```

### When This Works

- Single-domain applications (most Go services)
- When you want maximum simplicity with dependency inversion
- The root package name IS the domain name - very Go-idiomatic

### When This Breaks

- Multiple distinct domains - the root package becomes a dumping ground
- Can cause naming conflicts when the root package name is generic
- Still no `internal/` encapsulation

> **⚠️ Author's note:** Since both Ben Johnson articles (2014 and 2016) describe the same
> core structure (domain types in root package, implementations in subpackages), consider
> presenting them as ONE layout in the actual talk rather than two separate layouts.

### 2026 Perspective

> "The root package should contain your domain types. If it's getting too big, your domain might be too big for one module."  
> - Common interpretation of this pattern

**Speaker notes:**
- Both this article (2014) and "Standard Package Layout" (2016) describe essentially the same approach
- The 2014 article introduced the concepts; the 2016 article refined and formalized them
- Consider merging this section into Section 3 for the actual talk

---

## Section 5: Layout #4 - William Kennedy / Ardan Labs (6 min)

**Origin:** William Kennedy, "Package-Oriented Design" (2017)  
**URL:** https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

### Directory Structure

```
.
├── cmd/
│   └── libraryd/
│       └── main.go
├── internal/
│   ├── platform/
│   │   ├── database/
│   │   │   └── postgres.go
│   │   └── http/
│   │       └── server.go
│   ├── users/
│   │   ├── model.go
│   │   └── service.go
│   ├── books/
│   │   ├── model.go
│   │   └── service.go
│   └── reviews/
│       ├── model.go
│       └── service.go
└── go.mod
```

### Key Ideas

**Three layers of packages:**

1. **Domain packages** (`internal/users/`, `internal/books/`, `internal/reviews/`)  
   Business logic, models, and service interfaces.  
   These are the core of the application.

2. **Platform packages** (`internal/platform/`)  
   Foundational, reusable infrastructure code - database connections, HTTP server setup, middleware.  
   NOT business-specific.

3. **Composition root** (`cmd/libraryd/`)  
   Wires domain and platform together.

### Key Code

```go
// internal/users/model.go
package users

type User struct {
    ID    int
    Name  string
    Email string
}
```

```go
// internal/users/service.go
package users

import "database/sql"

type Service struct {
    DB *sql.DB
}

func (s *Service) User(id int) (User, error) {
    return User{}, nil
}
```

```go
// internal/platform/database/postgres.go
package database

import "database/sql"

func Open(connString string) (*sql.DB, error) {
    return sql.Open("postgres", connString)
}
```

```go
// cmd/libraryd/main.go
package main

import (
    "log"
    "os"
    "github.com/example/internal/platform/database"
    "github.com/example/internal/users"
)

func main() {
    db, _ := database.Open(os.Getenv("DB"))
    userService := users.Service{DB: db}
    // ... wire handlers, start server
}
```

### When This Works

- Medium-to-large applications
- Teams with multiple developers - `internal/` enforces boundaries
- When you have reusable infrastructure code (database, HTTP, logging)
- Enterprise/corporate Go - very popular in this segment

### When This Breaks

- `internal/platform/` can become a "junk drawer" over time
- Domain packages directly import `database/sql` - coupling to infrastructure
- The three-layer model doesn't naturally support multiple transport layers (gRPC + HTTP)
- Can be overkill for small services

### 2026 Perspective

- Still widely used, especially in enterprise Go
- The `internal/` convention is universally accepted now
- The `platform/` concept has evolved - many teams now split platform into separate packages or even modules

**Speaker notes:**
- The `internal/` directory is enforced by the Go compiler - external packages literally cannot import from it
- This was one of the first layouts to use `internal/` systematically
- Compare the domain coupling: here `users.Service` directly uses `*sql.DB`, unlike Ben Johnson's approach where domain defines interfaces

---

## Section 6: Layout #5 - Peter Bourgon Complex Layout (5 min)

**Origin:** Peter Bourgon, "Go Best Practices 2016"  
**URL:** https://peter.bourgon.org/go-best-practices-2016/#repository-structure

### Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── pkg/
│   ├── http/
│   │   ├── handler.go
│   │   └── server.go
│   ├── library/
│   │   ├── book.go
│   │   ├── user.go
│   │   └── review.go
│   └── postgres/
│       ├── connection.go
│       ├── user_service.go
│       ├── book_service.go
│       └── review_service.go
└── go.mod
```

### Key Idea: The `pkg/` Convention

- Everything importable goes under `pkg/`
- `pkg/` signals "this is public, reusable code"
- Very similar to Ben Johnson's layout but with an explicit `pkg/` namespace
- Popular in monorepos where multiple binaries share packages

### The `pkg/` Debate (2026)

> "`pkg/` doesn't make much sense with Go modules. In the GOPATH era, `pkg/` separated pre-compiled packages. With modules, this distinction doesn't exist."  
> - Eli Bendersky, "Simple Go Project Layout with Modules" (2019)

> "Don't use a `pkg/` directory. It's a convention from pre-modules Go that doesn't provide any benefit today."  
> - Multiple community members on Reddit r/golang (2024)

**The counter-argument:**

> "`pkg/` is useful in monorepos to distinguish between internal application code and shared library code."  
> - Some teams still defend it for organizational clarity

### When This Works

- Monorepos with shared packages
- When multiple binaries need the same domain logic
- Teams accustomed to the `pkg/` convention from other languages

### When This Breaks

- Single-binary applications - `pkg/` is an unnecessary extra level
- Modern Go modules make `pkg/` semantically meaningless
- Can confuse newcomers: "do I use `pkg/` or not?"

**Speaker notes:**
- Show this as a historical stepping stone
- The `pkg/` debate is one of the most heated in Go project structure discussions
- Peter Bourgon himself has moved toward simpler layouts in recent years

---

## Section 7: Layout #6 - golang-standards/project-layout (7 min)

**Origin:** Community repo on GitHub (2018+)  
**URL:** https://github.com/golang-standards/project-layout  
**Stars:** ~56k (one of the most-starred Go repos on GitHub)

### Directory Structure

```
.
├── cmd/
│   └── library/
│       └── main.go
├── internal/
│   └── app/
│       └── library/
│           ├── book.go
│           ├── user.go
│           └── user_service.go
├── pkg/                  # "Reusable packages" (debated)
│   └── ...               # (often empty in practice)
├── api/                  # API protocol definitions
├── web/                  # Web assets
├── configs/              # Configuration files
├── scripts/              # Build/deploy scripts
├── test/                 # Additional integration tests
├── docs/                 # Documentation
├── build/                # Build output
└── deployments/          # Docker, K8s manifests
```

### What It Defines

A comprehensive reference of conventional directories:

| Directory | Purpose |
|-----------|---------|
| `/cmd` | Main applications - each subdirectory = one binary |
| `/internal` | Private application and library code |
| `/pkg` | Library code ok to use by external applications |
| `/api` | API protocol definitions (OpenAPI, gRPC .proto) |
| `/web` | Web application specific assets |
| `/configs` | Configuration file templates |
| `/scripts` | Scripts for build, install, analysis, etc. |
| `/test` | Additional external test apps and test data |
| `/docs` | Design and user documentation |
| `/build` | Packaging and CI output |
| `/deployments` | IaaS, Docker, Kubernetes configs |

### The Critical 2026 Perspective

> "It is not an official Go standard, despite the name, so beginners can mistake it for endorsed guidance."  
> - Hacker News discussion (2021, still relevant in 2026)

> "The layout can encourage cargo-cult programming: people copy the structure because it looks professional, not because they understand the tradeoffs."  
> - kau.sh, "Cargo Culting" (2024)

> "The repo is best treated as a template to adapt, not a rule to follow blindly."  
> - waken.dev, "Go Project Layout Controversy" (2024)

> "Beginners see 15 directories and think they need all of them. You don't."  
> - Reddit r/golang (2024)

**Specific criticisms (with sources):**

1. **Misleading name** - "golang-standards" implies official endorsement; it's a community repo  
   Source: https://news.ycombinator.com/item?id=25651407

2. **`pkg/` is outdated** - Go modules made this convention unnecessary  
   Source: https://github.com/golang-standards/project-layout/issues/41

3. **Over-engineered for most projects** - 15+ directories when 3 would do  
   Source: https://www.reddit.com/r/golang/comments/1fi3xa9/seeking_advice_on_go_project_structure/

4. **Cargo cult adoption** - copied without understanding why  
   Source: https://kau.sh/blog/cargo-culting/

5. **Still useful as a REFERENCE** - when you DO need `/api/`, `/deployments/`, etc., this shows you where they conventionally go  
   Source: https://eventsandstuff.substack.com/p/go-standard-project-layout-a-mildly

### When This Works

- Large, complex projects that actually need all these directories
- When you need a shared vocabulary across a large team
- Reference for "where would X go?" - use it as a menu, not a blueprint

### When This Breaks

- Small-to-medium projects (most Go projects)
- When treated as a template to copy wholesale
- When the team doesn't understand why each directory exists

**Speaker notes:**
- BE BALANCED - don't just bash it. Acknowledge 50k stars means many people find it useful
- The key critique is the NAME and the CARGO CULT behavior, not necessarily the content
- Many individual directory conventions are fine on their own
- The problem is copying ALL of them at project start
- Show the "after" picture: a real project that uses only 4-5 of these directories

---

## Section 8: The 2020–2026 Shift - What Changed (5 min)

**Purpose:** A bridge section explaining WHY new approaches emerged.

### Timeline

```
2016 ─── Ben Johnson, Peter Bourgon, William Kennedy publish their approaches
 │
2017 ─── Community refines and debates these patterns
 │
2018 ─── golang-standards/project-layout repo becomes popular
 │       Go 1.11 introduces modules (preview)
 │
2019 ─── Go 1.12+ modules go mainstream
 │       Eli Bendersky: "Simple Go Project Layout with Modules"
 │       → GOPATH is dead, `pkg/` loses meaning
 │
2020 ─── Go 1.14-1.15 - modules are stable
 │       "Start simple" philosophy gains traction
 │
2021 ─── Go 1.18 generics preview
 │       Hexagonal architecture examples start appearing more in Go
 │
2022 ─── Go 1.18-1.19 - generics, workspaces (go.work)
 │       Multi-module monorepos become practical
 │
2023 ─── Alex Edwards publishes "Let's Go" and "Let's Go Further"
 │       Pragmatic middle-ground structure becomes popular
 │
2024 ─── Go 1.22 - stdlib HTTP routing eliminates need for 3rd-party routers
 │       Go 1.23 - iter package, improved slog (August 2024)
 │       LaurentSV publishes "No-Nonsense Go Package Layout"
 │       Mat Ryer publishes "How I Write HTTP Services in Go After 13 Years" (February 2024)
 │       Domain-driven / feature-based layouts gain momentum
 │
2025 ─── Go 1.24 - Swiss tables, refined tooling (February 2025)
 │       Community consensus: "start flat, grow when needed"
 │       Hexagonal / ports-and-adapters well-established pattern
 │
2026 ─── This talk! 🎉
```

### Key Changes That Drove New Approaches

1. **Go Modules (2019)** made `pkg/` largely unnecessary - packages are already namespaced by module path
2. **Go 1.22 ServeMux (2024)** removed the need for 3rd-party routers → simpler handler organization
3. **Go 1.23 iter & slog (2024)** continued enriching the standard library
4. **`go.work` (2022)** made multi-module monorepos practical → new layout patterns
5. **Community maturity** - the Go community has had 10 years to see what works and what doesn't
6. **Hexagonal / DDD patterns** - proven in Java/C#, now adapted (carefully) for Go

**Speaker notes:**
- This timeline slide grounds the audience - explains WHY we're talking about "new" structures
- Emphasize that the OLD structures aren't "wrong" - they were right for their era
- The Go ecosystem evolves, and so should our structural patterns

---

## Section 9: Layout #7 - Feature-Based / Domain-Driven (7 min)

**Origin:** Multiple sources - Alex Edwards (2023), community consensus (2024–2026)  
**Influenced by:** Domain-Driven Design concepts adapted for Go

### Core Idea: Organize by FEATURE, not by LAYER

**Instead of this (layer-based):**
```
internal/
├── handlers/       # all HTTP handlers mixed together
├── services/       # all business logic mixed together
├── models/         # all data types mixed together
└── repositories/   # all database code mixed together
```

**Do this (feature-based):**
```
internal/
├── user/           # everything about users
│   ├── handler.go  # HTTP handlers for users
│   ├── service.go  # user business logic
│   ├── repository.go  # user data access (interface + impl)
│   └── model.go    # user types
├── book/           # everything about books
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── model.go
├── review/         # everything about reviews
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── model.go
└── platform/       # shared infrastructure only
    ├── database/
    └── middleware/
```

### Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── book/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── review/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   └── platform/
│       ├── database/
│       │   └── postgres.go
│       └── middleware/
│           └── auth.go
└── go.mod
```

### Key Code

```go
// internal/user/model.go
package user

type User struct {
    ID    int
    Name  string
    Email string
}
```

```go
// internal/user/repository.go
package user

import "context"

type Repository interface {
    ByID(ctx context.Context, id int) (User, error)
    All(ctx context.Context) ([]User, error)
    Create(ctx context.Context, u User) error
}
```

```go
// internal/user/service.go
package user

import "context"

type Service struct {
    Repo Repository
}

func (s *Service) User(ctx context.Context, id int) (User, error) {
    return s.Repo.ByID(ctx, id)
}
```

```go
// internal/user/handler.go - using Go 1.22+ stdlib routing
package user

import (
    "encoding/json"
    "net/http"
)

type Handler struct {
    Service *Service
}

func (h *Handler) Routes() map[string]http.HandlerFunc {
    return map[string]http.HandlerFunc{
        "GET /users":      h.list,
        "GET /users/{id}": h.get,
        "POST /users":     h.create,
    }
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
    // r.PathValue("id") available in Go 1.22+
    id := r.PathValue("id")
    user, err := h.Service.User(r.Context(), atoi(id))
    // ... encode response
}
```

```go
// cmd/server/main.go - composition root
package main

import (
    "net/http"
    "github.com/example/internal/user"
    "github.com/example/internal/book"
    "github.com/example/internal/review"
    "github.com/example/internal/platform/database"
)

func main() {
    db := database.MustOpen("connection")

    userRepo := user.NewPostgresRepository(db)
    userSvc := user.NewService(userRepo)
    userHandler := user.NewHandler(userSvc)

    bookRepo := book.NewPostgresRepository(db)
    bookSvc := book.NewService(bookRepo)
    bookHandler := book.NewHandler(bookSvc)

    mux := http.NewServeMux()
    for pattern, handler := range userHandler.Routes() {
        mux.HandleFunc(pattern, handler)
    }
    for pattern, handler := range bookHandler.Routes() {
        mux.HandleFunc(pattern, handler)
    }

    http.ListenAndServe(":8080", mux)
}
```

### When This Works

- Medium-to-large services with multiple distinct features/domains
- Teams where different developers own different features
- When you want each feature to be independently understandable
- Microservice candidates - each feature folder could become its own service

### When This Breaks

- Cross-cutting concerns (e.g., a "borrow book" operation touches user + book + review)
- Small applications - overhead of multiple packages for 3 files each
- Requires discipline to keep `platform/` from becoming a junk drawer

### 2026 Perspective

> "Organize packages around features or domains, because that tends to scale better than purely layer-based folders."  
> - Alex Edwards, "11 Tips for Structuring Your Go Projects"

> "A good mental model is: domain folders group behavior by feature rather than by layer."  
> - Community consensus (2024–2025)

**Speaker notes:**
- This is arguably the MOST POPULAR layout for new Go web services in 2025–2026
- It combines the best ideas from Ben Johnson (dependency inversion) and William Kennedy (`internal/` + platform)
- Show the "before and after" comparison with layer-based structure
- The key test: "can I understand the 'book' feature by reading just the `book/` folder?"

---

## Section 10: Layout #8 - Hexagonal Architecture / Ports & Adapters (8 min)

**Origin:** Alistair Cockburn (2005), adapted for Go by many (2022–2026)
**Also known as:** Ports & Adapters
**Related patterns:** Onion Architecture (Jeffrey Palermo, 2008), Clean Architecture (Robert C. Martin, 2012) - share the dependency-inversion principle but are distinct patterns; Cockburn himself notes that the primary/secondary actor asymmetry "makes this fundamentally different from neighboring patterns such as the onion architecture."

### Core Principle

> "Code pertaining to the 'inside' part should not leak into the 'outside' part."
> - Alistair Cockburn, original article (2005)

In practice, this means the domain core has ZERO external dependencies. Everything else is an adapter.

```
                    ┌─────────────────────┐
                    │    HTTP Handler      │  Inbound Adapter
                    │  (adapter/in/http)   │  (driving)
                    └─────────┬───────────┘
                              │ depends on
                              ▼
                    ┌─────────────────────┐
                    │    Domain Core       │  "Inside" - no knowledge
                    │  (core/domain +     │  of outside technologies
                    │   core/ports)        │
                    └─────────▲───────────┘
                              │ implements
                    ┌─────────┴───────────┐
                    │    PostgreSQL        │  Outbound Adapter
                    │  (adapter/out/pg)    │  (driven)
                    └─────────────────────┘
```

### Why "Hexagonal"?

> "The hexagon is not a hexagon because the number six is important, but rather to allow the people doing the drawing to have room to insert ports and adapters as they need, not being constrained by a one-dimensional layered drawing."
> - Alistair Cockburn

The hexagon shape is a visual metaphor to break free from the layered (top-to-bottom) diagram mindset. It makes the **inside-outside** asymmetry obvious and leaves room for multiple ports - most applications have two, three, or four ports (Cockburn has never encountered more than four).

### Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Composition root - wires adapters to core
├── internal/
│   ├── core/
│   │   ├── domain/              # Entities & value objects
│   │   │   ├── book.go
│   │   │   ├── user.go
│   │   │   └── review.go
│   │   ├── port/                # Interfaces (contracts)
│   │   │   ├── book_repository.go    # outbound port
│   │   │   ├── user_repository.go
│   │   │   ├── review_repository.go
│   │   │   └── book_service.go       # inbound port (use case)
│   │   └── service/             # Business logic / use cases
│   │       ├── book_service.go
│   │       └── user_service.go
│   └── adapter/
│       ├── in/                  # Driving adapters (inbound)
│       │   └── http/
│       │       ├── book_handler.go
│       │       ├── user_handler.go
│       │       └── router.go
│       └── out/                 # Driven adapters (outbound)
│           ├── postgres/
│           │   ├── connection.go
│           │   ├── book_repository.go
│           │   ├── user_repository.go
│           │   └── review_repository.go
│           └── memory/          # In-memory implementation for tests
│               ├── book_repository.go
│               └── user_repository.go
└── go.mod
```

### Key Code - Domain Core (zero external dependencies)

```go
// internal/core/domain/book.go
package domain

type Book struct {
    ID       int
    Title    string
    Author   string
    Borrower *User
}
```

```go
// internal/core/port/book_repository.go - outbound port
package port

import (
    "context"
    "github.com/example/internal/core/domain"
)

type BookRepository interface {
    ByID(ctx context.Context, id int) (domain.Book, error)
    All(ctx context.Context) ([]domain.Book, error)
    Save(ctx context.Context, book domain.Book) error
    Update(ctx context.Context, book domain.Book) error
}
```

```go
// internal/core/service/book_service.go - business logic
package service

import (
    "context"
    "fmt"
    "github.com/example/internal/core/domain"
    "github.com/example/internal/core/port"
)

type BookService struct {
    books port.BookRepository
    users port.UserRepository
}

func NewBookService(books port.BookRepository, users port.UserRepository) *BookService {
    return &BookService{books: books, users: users}
}

func (s *BookService) Borrow(ctx context.Context, bookID, userID int) error {
    book, err := s.books.ByID(ctx, bookID)
    if err != nil {
        return fmt.Errorf("find book: %w", err)
    }
    if book.Borrower != nil {
        return fmt.Errorf("book already borrowed")
    }

    user, err := s.users.ByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("find user: %w", err)
    }

    book.Borrower = &user
    return s.books.Update(ctx, book)
}
```

### Key Code - Outbound Adapter (PostgreSQL)

```go
// internal/adapter/out/postgres/book_repository.go
package postgres

import (
    "context"
    "database/sql"
    "github.com/example/internal/core/domain"
    "github.com/example/internal/core/port"
)

// Compile-time check: BookRepository implements port.BookRepository
var _ port.BookRepository = (*BookRepository)(nil)

type BookRepository struct {
    db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
    return &BookRepository{db: db}
}

func (r *BookRepository) ByID(ctx context.Context, id int) (domain.Book, error) {
    var b domain.Book
    err := r.db.QueryRowContext(ctx,
        "SELECT id, title, author FROM books WHERE id = $1", id,
    ).Scan(&b.ID, &b.Title, &b.Author)
    return b, err
}

func (r *BookRepository) Save(ctx context.Context, b domain.Book) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO books (title, author) VALUES ($1, $2)", b.Title, b.Author,
    )
    return err
}
// ... All, Update omitted for brevity
```

### Key Code - Inbound Adapter (HTTP)

```go
// internal/adapter/in/http/book_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/example/internal/core/service"
)

type BookHandler struct {
    svc *service.BookService
}

func NewBookHandler(svc *service.BookService) *BookHandler {
    return &BookHandler{svc: svc}
}

func (h *BookHandler) Routes() map[string]http.HandlerFunc {
    return map[string]http.HandlerFunc{
        "GET /books/{id}": h.get,
        "POST /books/{id}/borrow": h.borrow,
    }
}

func (h *BookHandler) get(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.PathValue("id"))
    book, err := h.svc.Get(r.Context(), id)
    if err != nil {
        http.Error(w, "not found", 404)
        return
    }
    json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) borrow(w http.ResponseWriter, r *http.Request) {
    bookID, _ := strconv.Atoi(r.PathValue("id"))
    // parse userID from auth context or body
    if err := h.svc.Borrow(r.Context(), bookID, 1); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
```

### Key Code - Composition Root

```go
// cmd/server/main.go - THE ONLY PLACE that knows about concrete implementations
package main

import (
    "log"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"

    "github.com/example/internal/adapter/in/http"
    "github.com/example/internal/adapter/out/postgres"
    "github.com/example/internal/adapter/out/memory"
    "github.com/example/internal/core/service"
)

func main() {
    db, _ := sql.Open("postgres", "connection")

    // Outbound adapters (choose implementation)
    bookRepo := postgres.NewBookRepository(db)
    userRepo := postgres.NewUserRepository(db)
    // OR for testing: bookRepo := memory.NewBookRepository()

    // Core services (pure business logic)
    bookSvc := service.NewBookService(bookRepo, userRepo)

    // Inbound adapters
    bookHandler := http.NewBookHandler(bookSvc)

    // Wire routes (Go 1.22+ pattern-based routing)
    mux := http.NewServeMux()
    for pattern, handler := range bookHandler.Routes() {
        mux.HandleFunc(pattern, handler)
    }

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Testing Advantage - Swap Postgres for Memory

```go
// internal/adapter/out/memory/book_repository.go - test double
package memory

type BookRepository struct {
    books map[int]domain.Book
}

func NewBookRepository() *BookRepository {
    return &BookRepository{books: make(map[int]domain.Book)}
}

func (r *BookRepository) ByID(_ context.Context, id int) (domain.Book, error) {
    b, ok := r.books[id]
    if !ok {
        return domain.Book{}, fmt.Errorf("not found")
    }
    return b, nil
}
// ... Save, All, Update are simple map operations
```

```go
// Test - no database needed!
func TestBorrowBook(t *testing.T) {
    bookRepo := memory.NewBookRepository()
    userRepo := memory.NewUserRepository()
    svc := service.NewBookService(bookRepo, userRepo)

    // seed data
    bookRepo.Save(context.Background(), domain.Book{ID: 1, Title: "Go 101"})

    err := svc.Borrow(context.Background(), 1, 42)
    assert.NoError(t, err)
}
```

### When This Works

- Services with clear business rules that need to be protected from infrastructure changes
- When you need multiple implementations of the same port (Postgres + Memory + gRPC)
- When you want tests that don't touch any external system
- APIs that might change transport (HTTP → gRPC → CLI)
- Teams where different people own different adapters

### When This Breaks

- CRUD-only services - the abstraction overhead isn't worth it
- Small teams - the extra indirection can slow development
- Simple scripts and tools
- Risk of "architecture astronaut" behavior - too many layers for simple logic

### 2026 Perspective

> "Hexagonal architecture in Go is about keeping business logic pure while making infrastructure pluggable. It's not about folder names - it's about dependency direction."  
> - Community interpretation (2024–2025)

> "Clean Architecture struggles in Golang because Go's philosophy favors simplicity over abstraction layers."  
> - dev.to, "Why Clean Architecture Struggles in Golang" (2024)

**Speaker notes:**
- DEMO the testing advantage - this is the killer feature of hexagonal
- Show the compile-time interface check: `var _ port.BookRepository = (*BookRepository)(nil)`
- Emphasize: `internal/core/` has ZERO imports from `adapter/` - that's the whole point
- Be honest about the overhead - this isn't for every project
- Compare to Ben Johnson's Layout #2 - same dependency inversion principle, just more explicit

---

## Section 11: Layout #9 - Monorepo / Go Workspaces (6 min)

**Origin:** Go 1.18 workspaces (2022), matured through Go 1.23–1.24  
**URL:** https://go.dev/doc/tutorial/workspaces

### Directory Structure

```
library-monorepo/
├── go.work
├── services/
│   ├── api/                    # HTTP API service
│   │   ├── go.mod
│   │   ├── cmd/
│   │   │   └── api/
│   │   │       └── main.go
│   │   └── internal/
│   │       ├── handler/
│   │       └── middleware/
│   └── worker/                 # Background job worker
│       ├── go.mod
│       ├── cmd/
│       │   └── worker/
│       │       └── main.go
│       └── internal/
│           └── jobs/
├── libs/                       # Shared libraries
│   ├── library/                # Domain types & interfaces
│   │   ├── go.mod
│   │   ├── book.go
│   │   ├── user.go
│   │   └── review.go
│   └── dbutil/                 # Database utilities
│       ├── go.mod
│       └── dbutil.go
└── tools/                      # Development tools
    └── migrate/
        ├── go.mod
        └── main.go
```

### Key Code

```go
// go.work
go 1.24

use (
    ./services/api
    ./services/worker
    ./libs/library
    ./libs/dbutil
    ./tools/migrate
)
```

```go
// services/api/go.mod
module github.com/example/services/api

go 1.24

require (
    github.com/example/libs/library v0.0.0
    github.com/example/libs/dbutil v0.0.0
)
```

```go
// services/api/cmd/api/main.go
package main

import (
    "github.com/example/libs/library"   // resolves to local checkout!
    "github.com/example/libs/dbutil"
)

func main() {
    db := dbutil.MustOpen("connection")
    _ = library.Book{Title: "Go 101"}  // use shared domain types
    // ...
}
```

### When to Use Workspaces vs. Single Module

| Situation | Single Module | Workspaces |
|-----------|:---:|:---:|
| Single binary service | ✅ | |
| Service + shared library | | ✅ |
| Multiple services, one team | | ✅ |
| Multiple services, multiple teams | | ✅ |
| Independent release cycles | | ✅ |
| Microservices | | ✅ |

### Two Monorepo Strategies

**Strategy A: Multi-module (go.work)**
- Each service and library has its own `go.mod`
- `go.work` ties them together for local development
- Pros: Independent versioning, clear boundaries, CI builds only what changed
- Cons: More setup, module management overhead

**Strategy B: Single-module monorepo**
- One `go.mod` at root, everything is a package
- Use `internal/` for boundaries
- Pros: Simpler, no module management
- Cons: Everything rebuilds together, no independent versioning

```
# Single-module monorepo
library-monorepo/
├── go.mod              # ONE module for everything
├── services/
│   ├── api/
│   │   └── main.go     # package main
│   └── worker/
│       └── main.go     # package main
├── internal/
│   ├── user/
│   ├── book/
│   └── review/
└── pkg/
    └── dbutil/
```

### When This Works

- Multiple deployable services sharing domain logic
- Organization with 3+ Go services
- When different services have different release cycles
- When you want shared types (gRPC protos, domain models) across services

### When This Breaks

- Single service - monorepo adds complexity for no benefit
- Very large organizations - multi-repo with published modules may be better
- When teams can't agree on shared library versions

**Speaker notes:**
- `go.work` was THE biggest structural addition to Go since modules
- Before `go.work`, monorepos used `replace` directives in go.mod - messy
- Many large Go shops (Uber, Dropbox, Stripe) use monorepo patterns
- Show the workflow: edit shared lib → both services pick up changes immediately
- Keep this practical - don't get into CI/CD pipeline details

---

## Section 12: Bonus - Go 1.22+ Standard Library Routing Impact (3 min)

**Purpose:** Show how Go 1.22's routing improvements simplify project structure.

### Before Go 1.22 - You Needed a Third-Party Router

```go
// Required: chi, gorilla/mux, echo, gin, fiber, etc.
r := chi.NewRouter()
r.Get("/users/{id}", getUser)
r.Post("/users", createUser)
```

**Impact on structure:** You needed a `router/` or `transport/` package to encapsulate the third-party dependency.

### After Go 1.22 - Standard Library Does It All

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("GET /books/{path...}", getBook)  // wildcard

// In handlers:
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // built-in path parameter access
    // ...
}
```

### Structural Impact

- **No third-party router dependency** → one less package to organize
- **Method-based patterns** (`"GET /users"`) → handlers can be plain functions, no method dispatch needed
- **Path parameters built-in** → no need for request context middleware
- **Handlers can live closer to domain** → fewer indirection layers

```go
// internal/user/handler.go - clean, no external router dependency
package user

func Routes(svc *Service) map[string]http.HandlerFunc {
    return map[string]http.HandlerFunc{
        "GET /users":      List(svc),
        "GET /users/{id}": Get(svc),
        "POST /users":     Create(svc),
    }
}

func Get(svc *Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.PathValue("id"))
        user, err := svc.User(r.Context(), id)
        if err != nil {
            http.Error(w, "not found", 404)
            return
        }
        json.NewEncoder(w).Encode(user)
    }
}
```

**Speaker notes:**
- This is a quick but important point - Go's standard library keeps eating the ecosystem
- Many projects that used chi/gorilla can now use stdlib with zero structural changes
- Show the "before and after" - same functionality, fewer packages

---

## Section 13: Comparison & Decision Guide (5 min)

### Side-by-Side Comparison

| Layout | Complexity | Best For | Domain Isolation | Testability |
|--------|:----------:|----------|:----------------:|:-----------:|
| Flat | ⭐ | Prototypes, CLIs | ❌ | ⭐ |
| Ben Johnson (Standard) | ⭐⭐ | Small web services | ✅ | ⭐⭐⭐ |
| Ben Johnson (Original) | ⭐⭐ | Single-domain services | ✅ | ⭐⭐⭐ |
| William Kennedy | ⭐⭐⭐ | Enterprise services | ⚠️ | ⭐⭐ |
| Bourgon (pkg/) | ⭐⭐ | Monorepos (legacy) | ✅ | ⭐⭐ |
| Standard Layout | ⭐⭐⭐⭐ | Large complex projects | ⚠️ | ⭐⭐ |
| Feature-Based | ⭐⭐⭐ | Multi-feature services | ✅ | ⭐⭐⭐ |
| Hexagonal | ⭐⭐⭐⭐ | Core business logic | ✅✅ | ⭐⭐⭐⭐ |
| Monorepo | ⭐⭐⭐⭐ | Multi-service orgs | ✅ | ⭐⭐⭐ |

### Decision Flowchart

```
How many services are you building?
├── ONE service
│   └── How big is it?
│       ├── < 10 files → FLAT layout
│       ├── 10-50 files → BEN JOHNSON layout
│       ├── 50-200 files, single domain → FEATURE-BASED
│       └── 50+ files, complex business rules → HEXAGONAL
│
└── MULTIPLE services
    └── Do they share code?
        ├── No → Independent repos, whatever layout fits each
        └── Yes → MONOREPO with go.work
            └── Each service uses its own layout from above
```

### Universal Principles (regardless of layout)

1. **Start simple** - add structure when you feel pain, not before
2. **Use `cmd/`** for entry points when you have multiple binaries
3. **Use `internal/`** early - it costs nothing and prevents accidental coupling
4. **Keep dependency arrows pointing inward** - domain shouldn't know about HTTP or Postgres
5. **One `main.go` wires everything** - dependency injection at the composition root
6. **Organize by domain**, not by technical layer - `user/`, `book/`, not `handlers/`, `models/`
7. **Test from the outside** - test behavior, not structure

**Speaker notes:**
- The flowchart should be a visual slide
- Emphasize: you can MIX these approaches - hexagonal inside a monorepo, feature-based with flat sub-packages, etc.
- The universal principles are more important than any specific layout choice
- Good structure is about making the NEXT developer productive, not about architectural purity

---

## Section 14: Demo Walkthrough (8 min)

**Purpose:** Show all layouts side by side with the library app.

### Demo Setup

All demo code lives in a single repo under `gops2/demo/`:

```
demo/
├── 01-flat/              # Flat layout
├── 02-benjohnson-v1/     # Ben Johnson standard package layout
├── 03-benjohnson-v2/     # Ben Johnson revised
├── 04-kennedy/           # William Kennedy / Ardan Labs
├── 05-standard/          # golang-standards/project-layout
├── 06-feature/           # Feature-based / domain-driven
├── 07-hexagonal/         # Hexagonal / ports & adapters
└── 08-monorepo/          # go.work monorepo
    ├── go.work
    ├── services/
    └── libs/
```

### Demo Script

1. **Show the domain** - same types in each layout (30 sec each)
2. **Run each layout** - `go run ./cmd/server/main.go` and curl the API
3. **Highlight structural differences** - where does `User` live? Where does `Borrow` live?
4. **Show the dependency arrows** - which package imports which?
5. **Testing comparison** - show how hexagonal's memory adapter makes testing trivial
6. **Monorepo demo** - edit shared lib, show both services pick up changes

### Key Comparison Points (live demo)

For each layout, answer:
- Where is the `User` type defined?
- Where is the `Borrow(bookID, userID)` logic?
- Where is the HTTP handler?
- How do I test `Borrow` without a database?
- How many packages does a change to "add ISBN to Book" touch?

**Speaker notes:**
- This is the payoff - seeing the SAME app in 8 different layouts makes the tradeoffs tangible
- Run through quickly - the audience should be comparing in their heads
- The testing comparison (hexagonal memory adapter vs. others) is the most impactful moment

---

## Section 15: Summary & Resources (3 min)

### Key Takeaways

1. **There is no One True Layout™** - Go gives you freedom, use it wisely
2. **The community has evolved** - from "copy this template" to "start simple, grow when needed"
3. **Dependency direction matters more than directory names** - domain logic should never import infrastructure
4. **The best layout is the one your team can maintain** - consistency beats perfection
5. **Go 1.22+ changes the game** - stdlib routing means fewer third-party dependencies and simpler structures

### Resources

**Classic Articles (still relevant):**
- Ben Johnson - Standard Package Layout: https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
- Ben Johnson - Structuring Applications in Go: https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091
- William Kennedy - Package-Oriented Design: https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html
- Peter Bourgon - Go Best Practices 2016: https://peter.bourgon.org/go-best-practices-2016/
- Peter Bourgon - Go for Industrial Programming: https://peter.bourgon.org/go-for-industrial-programming/
- Dave Cheney - Five Suggestions for Setting Up a Go Project: https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project

**Modern Resources (2023–2026):**
- Alex Edwards - 11 Tips for Structuring Your Go Projects: https://www.alexedwards.net/blog/11-tips-for-structuring-your-go-projects
- LaurentSV - No-Nonsense Go Package Layout (2024): https://laurentsv.com/blog/2024/10/19/no-nonsense-go-package-layout.html
- Mat Ryer - How I Write HTTP Services in Go After 13 Years: https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/
- Eli Bendersky - Simple Go Project Layout with Modules: https://eli.thegreenplace.net/2019/simple-go-project-layout-with-modules/
- Go Blog - Routing Enhancements (Go 1.22): https://go.dev/blog/routing-enhancements
- Go Docs - Tutorial: Workspaces: https://go.dev/doc/tutorial/workspaces

**Critical Perspectives:**
- waken.dev - Go Project Layout Controversy: https://waken.dev/go-project-layout-controversy/
- kau.sh - Cargo Culting: https://kau.sh/blog/cargo-culting/
- dev.to - Why Clean Architecture Struggles in Golang: https://dev.to/lucasdeataides/why-clean-architecture-struggles-in-golang-and-what-works-better-m4g

**Reference Repos:**
- golang-standards/project-layout: https://github.com/golang-standards/project-layout
- This talk's demo repo: https://github.com/matejb/gops (2019) → gops2 (2026)

### Closing Quote

> "A good project structure is like good error handling - you only notice it when it's missing."  
> - Adapted for this talk

---

## Appendix A: Time Budget

| # | Section | Time |
|---|---------|------|
| 0 | Opening Hook | 2 min |
| 1 | The Example Application | 3 min |
| 2 | Layout #1 - Flat | 5 min |
| 3 | Layout #2 - Ben Johnson Standard Layout | 6 min |
| 4 | Layout #3 - Ben Johnson Original Article | 5 min |
| 5 | Layout #4 - William Kennedy | 6 min |
| 6 | Layout #5 - Bourgon / pkg/ | 5 min |
| 7 | Layout #6 - Standard Layout | 7 min |
| 8 | The 2020–2026 Shift | 5 min |
| 9 | Layout #7 - Feature-Based | 7 min |
| 10 | Layout #8 - Hexagonal | 8 min |
| 11 | Layout #9 - Monorepo | 6 min |
| 12 | Go 1.22+ Routing Impact | 3 min |
| 13 | Comparison & Decision Guide | 5 min |
| 14 | Demo Walkthrough | 8 min |
| 15 | Summary & Resources | 3 min |
|   | **Total** | **~84 min** |

**Note:** Total exceeds 60 min. Suggested cuts for 60-minute version:
- **Trim Section 4** (Ben Johnson Original) to 2 min - mention briefly, it's essentially the same as Section 3
- **Trim Section 6** (Bourgon pkg/) to 2 min - mainly historical
- **Shorten Section 7** (Standard Layout) to 5 min - be punchier with criticism
- **Shorten Section 12** (Go 1.22 routing) to 2 min - show one quick example
- **Demo** to 5 min - pre-selected highlights, not full walkthrough
- **Tighten transitions** between sections - save 3 min
- **Revised total: ~62 min** (with Q&A buffer)

## Appendix B: What's New vs. the 2019 Talk

| Aspect | 2019 Talk | 2026 Talk |
|--------|-----------|-----------|
| Layouts covered | 5 (Ben Johnson A+B, Kennedy, Bourgon simple+complex, Standard Layout) | 9 (adds feature-based, hexagonal, monorepo) |
| Domain example | Library app (users, books) | Same library app (users, books, reviews) |
| HTTP routing | Third-party assumed | Go 1.22+ stdlib routing highlighted |
| Monorepo support | Not covered | go.work workspace pattern |
| Critical perspective | Neutral survey | Critical analysis with cited sources |
| Testing comparison | Not covered | Per-layout testability comparison |
| Modern patterns | N/A | Feature-based, hexagonal, domain-driven |
| Decision framework | Not included | Flowchart + comparison matrix |

## Appendix C: Source Citations for Critical Claims

| Claim | Source | URL |
|-------|--------|-----|
| Standard Layout is not official | HN discussion | https://news.ycombinator.com/item?id=25651407 |
| Standard Layout encourages cargo cult | kau.sh | https://kau.sh/blog/cargo-culting/ |
| Standard Layout is over-engineered | Reddit r/golang | https://www.reddit.com/r/golang/comments/1fi3xa9/ |
| `pkg/` is outdated with modules | GitHub issue | https://github.com/golang-standards/project-layout/issues/41 |
| Start simple, grow when needed | LaurentSV (2024) | https://laurentsv.com/blog/2024/10/19/no-nonsense-go-package-layout.html |
| Feature-based > layer-based | Alex Edwards | https://www.alexedwards.net/blog/11-tips-for-structuring-your-go-projects |
| Clean Architecture struggles in Go | dev.to (2024) | https://dev.to/lucasdeataides/why-clean-architecture-struggles-in-golang-and-what-works-better-m4g |
| Mat Ryer handler pattern update | Reddit (2024) | https://www.reddit.com/r/golang/comments/1bg4d5e/ |
| Go 1.22 routing eliminates 3rd-party routers | Go Blog | https://go.dev/blog/routing-enhancements |
| Eli Bendersky on `pkg/` being unnecessary | eli.thegreenplace.net | https://eli.thegreenplace.net/2019/simple-go-project-layout-with-modules/ |
