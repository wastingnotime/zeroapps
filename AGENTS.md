# Repository Guidelines

## Project Structure & Module Organization

- `core/` — deterministic domain cores (currently `core/catcare/`).
- `docs/` — canonical architecture + domain specs (start with `docs/context.md`).
- `ai/` — agent operating contract, workflow, and decision log.
- `task.md` — current work definition and “done when” checklist.

## Build, Test, and Development Commands

- `go test ./...` — run the full test suite.
- `go test ./core/catcare -run TestLogWeight` — run a focused subset while iterating.
- `gofmt -w core/` — format Go code (keep diffs gofmt-clean).
- `go vet ./...` — basic static analysis.

## Coding Style & Naming Conventions

- Formatting: standard Go (`gofmt`), tabs for indentation, no manual alignment.
- Domain patterns: implement behavior as `Decide(command) -> events | rejection`, and state changes only in `Apply(event)`.
- Determinism: do not use hidden time/randomness/IO inside `core/` (pass timestamps/IDs via commands/events).
- Naming patterns used here:
  - Commands: `RegisterCat`, `LogWeight` (verb + noun).
  - Events: `CatRegistered`, `WeightLogged` (past tense).
  - Rejection codes: `snake_case` constants like `CodeInvalidWeight`.

## Testing Guidelines

- Framework: Go’s built-in `testing` package (`*_test.go` files live next to code).
- Tests should be deterministic: fixed timestamps (e.g. `"2026-02-14T10:00:00Z"`) and predictable IDs.
- Naming convention used in this repo: `TestXGivenYWhenZThenW` (behavioral, scenario-based).

## Commit & Pull Request Guidelines

- Commits: prefer Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`).
- PRs: include a short “why”, link/update `task.md` when relevant, update `docs/` for spec changes, and add tests that enforce invariants.

## Agent-Specific Instructions

Agents should start at `ai/bootstrap.md` and follow the contract; record irreversible decisions in `ai/decisions.md` only after explicit human approval.

