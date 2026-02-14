# purpose

ZeroApps is a model-first software laboratory for building **deterministic, invariant-preserving domain cores** that can be safely consumed by AI agents.

This repo’s first domain is **CatCare**: an event-sourced aggregate that represents one cat’s longitudinal care record (vaccines, appointments, weight history, anomalies, and treatments).

Primary outputs:
- A small Go domain core (commands → decisions → events → state)
- Deterministic tests that prove invariants and legal transitions
- Optional adapters (CLI/HTTP/bot) that *only* call the core

Non-goals:
- A SaaS/product roadmap
- An LLM orchestration framework
- UI-first “business logic in the interface”
- Hidden mutation via DB writes or ad-hoc scripts
