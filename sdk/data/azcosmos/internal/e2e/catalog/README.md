# Pinned shared E2E catalog

This directory is a vendored, read-only snapshot of the language-neutral Cosmos SDK E2E
catalog (schema, profiles, and scenario metadata) from `Azure/azure-sdk-for-rust`.

Pinned revision: `bfb35725eb9b02ab845c37c781db2ae364d26622` (`sdk/cosmos/e2e_tests`).

This revision is **not** the same as the `EmulatorSourceRevision` pinned in `../../ci.yml`
(`a05f024d22cd0a70f6ccac5e310cbbbaf40bf20e`). That emulator build predates the catalog
(2026-09-02, versus the catalog's introduction on 2026-09-11) and has no hosted-emulator
management API for deterministic fault/topology control. Smoke-tier scenarios that only need
item CRUD and query run against the currently pinned emulator; scenarios whose
`implementations/go.go` entry is `StatusBlocked` with a reason mentioning "emulator capability"
require bumping `EmulatorSourceRevision` first. Do not silently widen scope by relaxing that gap.

Do not hand-edit these JSON files. To refresh the pin:

1. Update the revision above and re-copy `schema/`, `profiles/`, and `scenarios/` verbatim from
   that revision of `Azure/azure-sdk-for-rust`'s `sdk/cosmos/e2e_tests`.
2. Re-run `go test ./internal/e2e/...` to confirm the Go implementation map
   (`../implementations/go.go`) still accounts for every scenario ID (active, blocked, or
   not-applicable).
3. Re-run the `../../e2e` package against the hosted emulator before merging.

This package only vendors `schema/`, `profiles/`, and `scenarios/`. `implementations/rust.json`
is Rust's own SDK implementation map and is intentionally not vendored; Go's equivalent lives in
`../implementations/go.go` as a compiled, type-checked map rather than parallel JSON.
