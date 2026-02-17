# workflow

- Work in small, incremental steps; keep the repo compiling and tests green.
- Use `docs/domains/catcare.md` as the canonical domain reference (CatCare aggregate, event catalog, command/result patterns, AI-safety protocol).
- Use `docs/context.md` as the canonical ZeroApps architecture reference (domain-core authority, adapters, command/result boundaries).

For domain work, prefer this order:
1. Define/adjust invariants and terminology
2. Define the command payload schema (what the adapter can ask for)
3. Implement `Decide(command) -> events | rejection`
4. Implement `Apply(event)` and derived state
5. Add deterministic tests: given events → when command → expect events/rejection
6. Only then add/adjust adapters (CLI/HTTP/bot), keeping them thin

Changes are considered done when:
- Invariants are enforced in the model (not in UI/adapters)
- Behavior is covered by deterministic tests
- The core has no hidden side effects

About decisions:
- When a task forces an irreversible choice, propose a decision entry using `ai/templates/decision.md`.
- Record it in `ai/decisions.md` only after explicit approval (see `ai/groundrules.md`).

The current task, if any, is defined in `task.md` at the repo root.

Session tracking during execution:
- For each Codex session working on this repo, create/update `ai/sessions/YYYY-MM-DD-<session_id>.md`.
- Keep `Handoff Notes` concise and task-focused; append key completed steps and current status before finishing.
