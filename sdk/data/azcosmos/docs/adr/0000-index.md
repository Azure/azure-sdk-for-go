# Architecture Decision Records — Go v2 FFI direction

ADRs (Architecture Decision Records) capture **what we decided** and a brief
**why**, in a minimal, easy-to-reference form. They are **numbered and
immutable**: once accepted, an ADR is not edited; a later ADR may supersede it.
Detailed discussion and exploration live in
[`../go-v2-ffi-decision.md`](../go-v2-ffi-decision.md) and
[`../go-ffi-distribution-design.md`](../go-ffi-distribution-design.md), not here.

**Format (template):** Status (immediately under the title) · Context (2-4
sentences) · Decision (1-3 bullets) · Consequences (2-4 bullets) · Alternatives
considered (1 line each).

| # | Title | Status |
|---|-------|--------|
| [0001](0001-go-v2-uses-ffi.md) | Go v2 uses the Rust driver through FFI | Proposed |

> This ADR is proposed for review. Packaging mechanics are intentionally left to
> the distribution design document; this record captures the Go v2 implementation
> direction.
