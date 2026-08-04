# Guide to migrate from `azcosmos` v1 to `azcosmos` v2

> **This guide is a stub.** It is filled in by
> [issue 27329](https://github.com/Azure/azure-sdk-for-go/issues/27329), which owns the v1→v2
> breaking-change list and migration guidance. The sections below record the shape it will take so
> that reviewers of the v2 surface know where each decision is written down.

This guide is intended to assist in the migration from
`github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos` to
`github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2`.

## Why v2 exists

v1 is a pure-Go implementation of the Cosmos DB data plane. v2 binds to the shared Rust Cosmos
driver instead, so that routing, retries, session handling, failover behavior and query fan-out
are consistent across the Cosmos DB SDKs rather than being reimplemented per language. The
decision record is ADR 0001, added by
[PR 27238](https://github.com/Azure/azure-sdk-for-go/pull/27238).

## Sections to be written

- **Import path and module.** Moving from `.../azcosmos` to `.../azcosmos/v2`.
- **Surface classification.** Every v1 export, marked kept as-is / changed / removed / new, with a
  replacement or a rationale for each removal.
- **Side-by-side samples.** v1 and v2 versions of client construction, point operations, query,
  change feed and transactional batch.
- **Build requirements.** v2 binds to a native driver, so it requires cgo (`CGO_ENABLED=1`). This
  changes the cross-compilation story, and v1's WebAssembly support cannot carry over to v2.
