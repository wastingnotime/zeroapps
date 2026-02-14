# ZeroApps

ZeroApps is a model-first software laboratory for building **deterministic, invariant-preserving domain cores** that can be safely consumed by AI agents.

This repository’s first domain is **CatCare**: an event-sourced aggregate for a cat’s longitudinal care record (vaccines, appointments, weight history, anomalies, and treatments).

## What this repo is (and is not)

Is:
- A small Go codebase focused on domain modeling (commands → decisions → events → state)
- Deterministic tests that prove invariants and legal transitions
- Thin adapters (optional) that only call into the core

Is not:
- A SaaS/product roadmap
- An “LLM orchestration framework”
- A UI-first project where business logic lives in the interface

## Canonical docs (humans + AI)

Read these to understand the project:
- `docs/context.md` — ZeroApps architecture and boundaries (“domain core is the authority; AI is an adapter”)
- `docs/domains/catcare.md` — CatCare domain specification (aggregate, invariants, event catalog, command patterns)

## AI contract (agent operating rules)

Codex (and other agents) must follow the repo’s operating rules:
- Start at `ai/bootstrap.md`
- Then follow the linked AI resources (`ai/contract.md`, `ai/purpose.md`, `ai/groundrules.md`, `ai/workflow.md`, `ai/decisions.md`)

Key constraints:
- No direct state mutation: changes must happen via commands that produce events
- Deterministic core: no hidden time/randomness/IO in domain logic
- Decisions are recorded only after explicit human approval

## How humans and AI should work here

Humans:
- Treat `docs/` as the shared spec (design intent, invariants, boundaries).
- Use `task.md` to define the current work and “done when”.
- Approve (or reject) irreversible choices before they are added to `ai/decisions.md`.

AI agents:
- Follow `ai/bootstrap.md` (role + operating rules for working in this repo).
- If input is ambiguous, ask one clarifying question before coding.
- Keep changes small, compile-safe, and test-backed.

## Suggested workflow

1. Update spec (if needed): `docs/context.md`, `docs/domains/catcare.md`
2. Implement domain: `Decide(command) -> events | rejection`, then `Apply(event)`
3. Add deterministic tests: given events → when command → expect events/rejection
4. Add adapters last (CLI/HTTP/bot), keeping them thin

## Go quickstart

Run tests:
```bash
go test ./...
```

## Repo map (current + intended)

- `docs/` — stable docs for humans + AI
- `ai/` — agent operating contract and decision log
- `task.md` — current task definition
- (intended) `core/`, `svc/`, `store/`, `adapters/`, `cmd/` — as the implementation grows
