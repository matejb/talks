# Go Project Structures — 2026 Edition

Talk for **GoTalks / Golang ZG meetup** (28.05.2026).

Go has no official project structure — so what do you do?
This talk surveys **8 common layout approaches**, from flat files to hexagonal architecture, all built against the same example domain so you can compare them side by side.

## Slides

Presented using [present](https://pkg.go.dev/golang.org/x/tools/present):

```
go install golang.org/x/tools/cmd/present@latest
present .
```

## Layout Examples

Runnable code for every layout lives in **[`layouts-examples/`](layouts-examples/README.md)** → see its README for details, API endpoints, and how to run & test each one.

| # | Layout | Key Idea |
|---|--------|-----------|
| 1 | Flat / No Structure | Everything in `package main` |
| 2 | Ben Johnson Standard Package | Domain in root, implementations in leaf packages |
| 3 | William Kennedy / Ardan Labs | `internal/` with domain + platform layers |
| 4 | Peter Bourgon `pkg/` | Everything importable under `pkg/` |
| 5 | golang-standards/project-layout | Comprehensive conventional directories |
| 6 | Feature-Based / Domain-Driven | Organize by feature, not by layer |
| 7 | Hexagonal / Ports & Adapters | Core has zero dependencies |
| 8 | Monorepo / Go Workspaces | Multiple modules with `go.work` |
