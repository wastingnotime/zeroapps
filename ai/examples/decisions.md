# decisions

## D-20260131-01 — License
- Status: accepted
- Context: Repo is open to the world; prefer keeping origin visible without restricting use.
- Decision: Use Apache-2.0.
- Why:
  - Preserves attribution/NOTICE expectations better than MIT
  - Still permissive and widely compatible
- Consequences:
  - + Clear lineage for forks and downstream use
  - - Longer license text than MIT

## D-20260131-02 — Simulation model scope
- Status: accepted
- Context: Visual demo-first; avoid building a physics engine.
- Decision: Use simple circle-circle collision with elastic impulse approximation; no rotation, no friction.
- Why:
  - Matches visible intent with low complexity
  - Keeps code small and stable
- Consequences:
  - + Predictable behavior, easier to maintain
  - - Not physically “complete”
