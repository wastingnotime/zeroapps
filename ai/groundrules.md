# ground rules

These rules are non-negotiable.

## Domain authority
- The **domain core is the authority**. The AI (and any UI/adapter) is an adapter.
- Adapters may translate/format/route, but must not “implement business rules”.

## State changes
- No direct state mutation. All changes happen via **commands** that produce **events**, which are then **applied**.
- Invalid commands must be rejected. Ambiguous input must result in clarification (no mutation).
- IDs for domain entities (`cat_id`, `item_id`, `entry_id`, etc.) are minted by the core (adapters can only reference IDs previously returned by the core).
- Idempotency is required: the same `command_id` must not apply twice.

## Determinism
- The core must be deterministic: time/randomness/IO are explicit inputs, never hidden globals.
- Keep the core pure: no network/filesystem/env access inside domain logic (put side effects behind adapters/services).

## Engineering constraints
- Prefer Go + standard library for the core.
- Do not introduce new dependencies without approval.
- Prefer explicit code over clever abstraction; minimal surface area.
- Avoid broad refactors unless requested or clearly required to maintain invariants.

## Ambiguity and safety
- If a requirement is unclear, ask **one** clarifying question before coding.

## Decisions
- Codex may suggest that something is a decision, but must not record one without explicit approval.
