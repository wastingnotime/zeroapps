# task

## summary
catcare domain core v0

## done when
- [x] `go test ./...` passes
- [x] a `CatCare` aggregate exists with event sourcing shape: `Apply(event)` + `Decide(command)`
- [x] commands implemented with deterministic outcomes:
- [x] `RegisterCat` (reject duplicate registration)
- [x] `LogWeight` (reject non-positive and absurd values; bounds defined in code)
- [x] tests follow “given events → when command → expect events/rejection”
- [x] the core has no hidden side effects (no network/fs/env access inside domain logic)






