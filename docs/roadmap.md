# Roadmap — ZeroApps (AI-Consumed Domain Core)

Just for keep in mind what path we are following.
Below some list of very tiny milestones.


## micro-milestones

[x] AI augmented development.
[x] Aggregate roots and aggregates (domain cores).
[x] Unit tests for Aggregate roots and aggregates (domain cores).
[x] Concurrency for Aggregate roots.
[x] A cli to help understand the usage of ZeroApps.
[x] Projections for each domain core.
[ ] catcare-cli: persist local status (event store + replay/snapshots) across runs.
[ ] CatCare: validate timestamps as RFC3339 (`BirthDate`, `At`) + tests for invalid formats.
[ ] CatCare: `RenameCat` command + `CatRenamed` event + invariants + tests.
[ ] CatCare: `CareItemScheduled` (mint `item_id`) + state + tests.
[ ] CatCare: `RescheduleCareItem` + invariants (exists, not canceled/completed) + tests.
[ ] CatCare: `CompleteCareItem` + invariants + tests.
[ ] CatCare: `CancelCareItem` + invariants + tests.
[ ] CatCare: `ReportAnomaly` / `ResolveAnomaly` + invariants + tests.
[ ] Projections: minimal projector interface + in-memory projection store.
[ ] Projections: `CatCareSummary` (name, last weight, unresolved anomalies, next due care items).
[ ] Projections: `UpcomingCareItems` (sorted by due date) + CLI query command.
[ ] svc/store: optional snapshot loading + tail events + tests.
[ ] TBD...
