# task

## summary
Build a first projection for the CatCare event-sourced aggregate to prove read-model support by listing registered cats.

## done when
- [x] A `RegisteredCats` projection derives data from `core/catcare` events (at minimum `CatRegistered`).
- [x] The projection exposes a query to list cats (e.g., `ListRegisteredCats() -> []{CatID, Name, BirthDate}`).
- [x] The projection is updated from the event append path (service publishes newly appended events after a successful append).
- [x] Projection updates are idempotent per aggregate stream version (re-applying the same events does not duplicate cats).
- [x] Unit tests cover: empty projection, single cat, multiple cats, and idempotent re-apply.
- [x] `go test ./...` passes.
- [x] update catcare-cli to include a call for the projection query.
- [x] update today session file on ai/sessions with the current session id.
