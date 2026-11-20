# Pinned shared E2E catalog

This directory is a vendored, read-only snapshot of the language-neutral Cosmos SDK E2E
catalog (schema, profiles, and scenario metadata) from `Azure/azure-sdk-for-rust`.

Pinned revision: `bfb35725eb9b02ab845c37c781db2ae364d26622` (`sdk/cosmos/e2e_tests`).

This revision is **not** the same as the `EmulatorSourceRevision` pinned in `../../ci.yml`
(`408c4eb145b189cb14771ed7ae1f52626e5873ef`). The emulator exposes a versioned management
capabilities endpoint and deterministic region/partition topology controls, but it does not expose
the deterministic fault-injection API required by the resilience scenarios. The Go fixture
realizes the `sessionSingleRegion` and `sessionTwoRegionsDelayed` account profiles with the public
azcosmos/v2 options; dynamic-topology scenario ports and a topology-control client remain out of
scope. Check each scenario's backend requirements and Go implementation entry rather than
assuming every blocker is cleared by changing the emulator revision.

Do not hand-edit these JSON files. To refresh the pin:

1. Update the revision above and re-copy `schema/`, `profiles/`, and `scenarios/` verbatim from
   that revision of `Azure/azure-sdk-for-rust`'s `sdk/cosmos/e2e_tests`.
2. Re-run `go test ./internal/e2e/...` to confirm the Go implementation map
   (`../implementations.go`) still accounts for every scenario ID (active, blocked, or
   not-applicable).
3. Re-run the `../../e2e` package against the hosted emulator before merging.

This package only vendors `schema/`, `profiles/`, and `scenarios/`. `implementations/rust.json`
is Rust's own SDK implementation map and is intentionally not vendored; Go's equivalent lives in
`../implementations.go` as a compiled, type-checked map rather than parallel JSON.
