# CatCare — Domain Spec (Canonical)

This document is the canonical domain specification for the first ZeroApps domain: **CatCare**.

It is intended to be read and understood by both humans and AI agents.

CatCare is designed as an **event-sourced aggregate** with explicit invariants and deterministic command handling, safe to be operated through AI and other adapters.

---

## 1) Aggregate Root: `CatCare`

Why: a cat’s vaccines, appointments, weight, anomalies, and treatments all orbit a single longitudinal record + reminders.

**Identity**
- `CatID` (stable ULID/UUID; minted by the core)

**State (derived from events)**
- Profile: name, birth/adoption date (optional)
- Schedules (vaccines, vet, treatment steps)
- Measurements (weight history)
- Observations (anomalies)
- Treatments (optional richer semantics than “scheduled steps”)

**Event stream**
- `stream = "catcare/<CatID>"`

---

## 2) Invariants (v0 examples)

These are enforced by the model (not the adapter):

- Weight must be positive and within sane bounds (configurable, but deterministic).
- You cannot complete/reschedule/cancel a care item that does not exist.
- Dates must be valid and not absurd (e.g., outside an allowed range).
- Idempotency: the same `command_id` must not apply twice.
- Unknown commands are rejected; invalid parameters are rejected.
- If parsing produces ambiguity, the command must be rejected with “needs clarification” (no mutation).

---

## 3) Event Catalog (v0)

Keep events small, explicit, and auditable.

### Core events
- `CatRegistered {cat_id, name, birth_date?}`
- `CatRenamed {cat_id, new_name}`

### Scheduling / reminders
- `CareItemScheduled {item_id, kind, title, due_at, recurrence?, metadata?}`
- `CareItemRescheduled {item_id, new_due_at}`
- `CareItemCompleted {item_id, completed_at, notes?}`
- `CareItemCanceled {item_id, reason?}`

Where `kind ∈ {VACCINE, VET_APPOINTMENT, TREATMENT_STEP, OTHER}`.

### Weight
- `WeightLogged {entry_id, at, grams, notes?}`

### Anomaly tracking
- `AnomalyReported {anomaly_id, at, summary, severity, tags[], notes?, attachments?}`
- `AnomalyResolved {anomaly_id, resolved_at, notes?}`

### Treatments (plan-level, optional)
- `TreatmentPlanStarted {plan_id, title, started_at, protocol?, notes?}`
- `TreatmentPlanUpdated {plan_id, patch...}`
- `TreatmentPlanEnded {plan_id, ended_at, outcome?, notes?}`

Note: for an MVP, treatments can be modeled purely as scheduled care items (steps/doses) and skip plans entirely.

---

## 4) Command Interface (Core Contract)

### 4.1 Envelope (recommended)

Every command is:
- pure data
- validated
- applied to aggregate
- produces events (or a rejection)
- returns a deterministic result

Envelope fields:
- `command_id` (ULID) — idempotency key
- `aggregate` — `{ type: "CatCare", id: "<cat_id>" }`
- `expected_version` — optimistic concurrency (optional but recommended)
- `actor` — `{ type: "ai"|"human"|"system", id: "<caller-id>" }`
- `time` — client-proposed timestamp (core may also record authoritative time via adapter/service)
- `command` — `{ name: "...", payload: {...} }`

### 4.2 Command names (v0)

- `RegisterCat`
- `RenameCat`
- `ScheduleCareItem`
- `RescheduleCareItem`
- `CompleteCareItem`
- `CancelCareItem`
- `LogWeight`
- `ReportAnomaly`
- `ResolveAnomaly`
- `StartTreatmentPlan` (optional)

### 4.3 Result schema (v0)

- `ok: true|false`
- `new_version`
- `events_applied` (count + ids/types)
- `rejections` (if not ok): `{code, message, field?}`

---

## 5) AI-Safe Protocol (Adapter ↔ Core)

The AI is not allowed to “do things”; it must request structured commands.

### Threat model
AI may:
- hallucinate ids/dates
- skip confirmations
- spam duplicates
- attempt destructive commands

The protocol must enforce:
- strict typing + bounds
- idempotency
- no “free-form execution”
- confirm-before-commit for risky actions

### Two-phase commit for AI

#### Phase 1 — `PROPOSE`
AI proposes a full structured command.

Core returns:
- `ACCEPTED_FOR_CONFIRMATION` + `confirmation_token` + a human-readable summary, or
- `REJECTED` with validation errors / clarification needs.

#### Phase 2 — `CONFIRM`
Only a trusted human (or explicitly trusted actor) can confirm using the token.
Core applies the command and emits events.

### Safety tiers (suggested)
- Tier A (auto-commit): `LogWeight`, `ReportAnomaly` (bounds/sanity enforced)
- Tier B (confirm): scheduling/rescheduling/canceling/completing items; treatment plans

### Hard constraints
- IDs (`item_id`, `entry_id`, …) are minted by the core; AI can only reference IDs previously returned by the core.
- Time can be proposed by the AI, but the core must reject ambiguous formats.
- Free text is limited to specific fields (`notes`, `summary`) with max length and sanitization.
- Attachments are referenced by opaque IDs minted by the adapter (never invented by the AI).

---

## 6) Implementation Shape (Go, event sourcing)

Required pieces:
- `Apply(event)` mutates state
- `Decide(command)` returns events or a rejection
- `LoadFrom(events[])`

Deterministic test style:
- Given events → when command → expect events (or rejection)
