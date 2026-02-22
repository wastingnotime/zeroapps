# task

## summary
Add persistent storage for `catcare-cli` so status survives process restarts, using SQLite as the single storage mechanism for event streams and projection rebuild.

## done when
- [x] A SQLite-backed implementation of `store.EventStore` exists and is usable by adapters.
- [x] `catcare-cli` accepts a DB path flag and uses SQLite store instead of in-memory-only state.
- [x] On startup, `catcare-cli` replays persisted events into `RegisteredCats` so `list-registered` reflects prior runs.
- [x] Existing command handling semantics remain unchanged (`register`, `log-weight`, `list-registered`).
- [x] Unit tests cover SQLite store append/load behavior and optimistic concurrency conflict.
- [x] `go test ./...` passes.
