# Architecture Context — ZeroApps (AI-Consumed Domain Core)

This file is a stable reference for **ZeroApps architecture and interaction boundaries**.

Domain specs live in:
- `docs/domains/catcare.md` (CatCare)

---

## 0) One-line purpose

ZeroApps is a model-first software laboratory focused on building deterministic, invariant-preserving domain cores designed to be safely consumed by AI agents.

---

## 1) Core Intent

ZeroApps exists to design domain systems that:

* Enforce invariants at the model level
* Prevent illegal state transitions
* Remain deterministic under evolution
* Stay small and composable
* Can be safely operated by an AI adapter

It is not a product.
It is not a SaaS.
It is not market-driven.

It is structural depth engineering.

---

## 2) Core Architectural Principle

**The domain core is the authority.
The AI is an adapter.**

The AI:

* Translates human language into structured commands
* Proposes actions
* Interprets results for the user

The core:

* Validates invariants
* Executes state transitions
* Rejects illegal operations
* Produces deterministic outputs

The AI never mutates state directly.

---

## 3) Problem ZeroApps Addresses

Modern systems often:

* Expose raw database mutation
* Allow uncontrolled state edits
* Depend on infrastructure for correctness
* Let UI layers perform business logic
* Allow LLMs to “hallucinate” domain behavior

ZeroApps exists to invert that.

It defines a strict boundary:

> LLM proposes.
> Core disposes.

---

## 4) Philosophical Principles (Non-Negotiable)

### 1️⃣ Invariants First

Every domain must:

* Explicitly define legal transitions
* Prevent invalid states
* Reject illegal commands
* Return structured outcomes

If correctness depends on UI validation or database constraints alone, the model is incomplete.

---

### 2️⃣ Behavior Over Raw State

Direct mutation is forbidden.

Instead of:

```
balance = 0
```

Prefer:

```
account.close()
```

Actions express intent.
State transitions are consequences of behavior.

---

### 3️⃣ Determinism

Given the same:

* state
* command
* inputs

The system must produce the same output.

Time, randomness, and external dependencies must be explicit inputs — never hidden globals.

---

### 4️⃣ Small Surface Area

ZeroApps avoids:

* Premature scalability layers
* Over-abstraction
* Framework dependency inside the domain
* Feature creep

Minimalism is structural.

---

### 5️⃣ Infrastructure Is Replaceable

Frameworks, transports, and UIs are adapters.

Examples of adapters:

* Telegram bot
* WebSocket SPA
* CLI
* LLM interface

None of these define domain truth.

The domain must survive their removal.

---

## 5) AI Consumption Model

ZeroApps is intentionally designed to be consumed by an AI agent.

The interaction model is:

Human → LLM → Structured Command → ZeroApps Core → Structured Result → LLM → Human

The LLM is:

* A translator
* A planner
* A UI layer

It is not:

* A state authority
* A rule engine
* A source of business logic

---

## 6) Command Contract (Critical Boundary)

The core exposes explicit commands.

Examples (illustrative):

* `CreateEntity(...)`
* `ScheduleAction(...)`
* `CloseAggregate(...)`
* `ListUpcoming(...)`

The AI must emit structured commands (e.g., JSON schema or tool calls).

The core:

* Validates schema
* Validates invariants
* Executes or rejects

Unknown commands must be rejected.
Invalid parameters must be rejected.
Ambiguous input must not mutate state.

---

## 7) Result Contract

Every command execution returns structured output:

* `Status` (Success | Rejected)
* `UserMessage` (calm explanation)
* `MachineDetails` (optional, structured)
* `StateSummary`
* Optional `Events[]` (for audit/log)

The LLM may transform the UserMessage into conversational language,
but the domain message originates from the core.

---

## 8) Initial Architecture (Exploratory Phase)

Early architecture candidates:

Option A:
Telegram → Bot → LLM → Command API → ZeroApps Core

Option B:
SPA → WebSocket → LLM → Command API → ZeroApps Core

The transport layer is secondary.

The domain contract remains identical.

---

## 9) LLM Constraints

The AI layer must:

* Only interact via defined commands
* Never generate direct database writes
* Never bypass the core
* Log attempted commands for audit

LLM ambiguity must result in clarification, not mutation.

The system must remain safe under hallucination.

---

## 10) Strategic Role in Henrique’s Ecosystem

ZeroApps is:

* The invariant engine behind applied systems (e.g., CatCare)
* A structural laboratory for AI-safe domain modeling
* A depth counterpoint to stack exploration (Contacts)

If CalmApps is applied calm software,
ZeroApps is structural calm software.

It explores how AI can operate deterministic systems without corrupting them.

---

## 11) What ZeroApps Is NOT

* Not an AI chatbot framework
* Not an LLM orchestration toolkit
* Not an automation engine
* Not a microservices template
* Not a productivity stack

AI consumption is a surface.
The domain is the substance.

---

## 12) Success Criteria

ZeroApps succeeds if:

* Invariants are explicit and enforced
* AI interaction cannot corrupt domain state
* Infrastructure can be swapped without rewriting the model
* Adding features does not explode complexity
* The model remains understandable years later

It does not require users.
It does not require revenue.
It requires coherence.

---

## 13) AI Assistance Guidelines

When assisting on ZeroApps:

* Prioritize domain modeling over infrastructure
* Challenge uncontrolled mutation
* Design commands before endpoints
* Design invariants before persistence
* Avoid startup/growth framing
* Avoid LLM-first logic

Always ask:

> Can this invariant be enforced in the model itself?



---

## 14) Language + repo shape (current preference)

Go fits ZeroApps well because:
* strong typing for schemas/invariants
* easy deterministic tests on event streams
* adapters can be thin and replaceable

Suggested repo shape (not mandatory):
* `core/` (pure domain: aggregate, commands, events, invariants)
* `store/` (event store, snapshots)
* `svc/` (application service: command handling, idempotency, proposal/confirm)
* `adapters/` (CLI/HTTP/bot)
* `cmd/` (binaries)

## 15) Event sourcing paradigm (aggregate)

Required pieces:
* `Apply(event)` mutates state
* `Decide(command)` produces events or rejection
* `LoadFrom(events[])`

Deterministic test style:
* Given events → when command → expect events (or rejection)
