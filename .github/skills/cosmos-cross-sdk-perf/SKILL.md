---
name: cosmos-cross-sdk-perf
description: Use when provisioning, running, validating, or aggregating the matched Cosmos DB Go v1, Go v2 native-driver, and Rust SDK performance campaign.
---

# Cosmos Cross-SDK Performance Campaign

This skill reproduces a matched, closed-loop Cosmos DB performance comparison. Read
[`README.md`](README.md), then use only the checked-in configuration, schemas, classifier, and
scripts in this directory.

## Required workflow

1. Copy `config/default.json` into a new directory outside the repository. Fill in subscription
   and resource values without adding secrets.
2. Run `python3 scripts/campaign.py validate --config <config>`.
3. On the control machine, run
   `python3 scripts/campaign.py prerequisites --phase control`.
4. Preview Azure changes with
   `python3 scripts/campaign.py provision --config <config> --state <state.json> --dry-run`.
5. Provision only after the preview is accepted. The script verifies the Cosmos region, chooses an
   available four-vCPU Linux VM SKU, enables a system-assigned identity, grants Cosmos DB Built-in
   Data Contributor, and writes ownership state.
6. On the VM, run `python3 scripts/campaign.py prerequisites --phase vm`, then build each adapter
   with profiling-quality symbols. The native Rust driver build must
   include `-C debuginfo=2`, `-C force-frame-pointers=yes`, and an ELF build ID, with stripping
   disabled. Build its pinned source-of-record commit from `Azure/azure-sdk-for-rust`, and record
   the pinned `Azure/azure-cosmos-driver` distribution commit consumed by Go v2. Record the complete
   build provenance in each leg's `provenance.json`.
7. Seed exactly 1,024 items whose serialized UTF-8 JSON payloads are exactly 1,024 bytes and whose
   partition key path is `/pk`.
8. Run each leg separately with the effective configuration in `config/effective.json`: eight
   strict closed-loop workers, deterministic `Read,Read,Read,Read,Read,Upsert,Upsert` cycles,
   30-second excluded warmup, and 300-second measurement. Use the same endpoint, database,
   container, keys, identity, VM, and process CPU affinity. Gateway V1 HTTP/2 is required, binary
   encoding is disabled, and retries/throttles must both be zero.
9. Wrap each adapter command with `scripts/profile-leg.sh`. It captures per-PID `perf` at 49 Hz,
   call chains, synchronized `/proc` memory snapshots, cgroup memory, per-thread stacks, Go pprof
   locations, binary/debug files, and checksums. Keep its `raw/` directory immutable.
10. Each adapter writes `metrics.json` and `perf-self.csv` matching the schemas. Run
    `python3 scripts/campaign.py validate-leg --leg-dir <leg> --config <config>`.
11. Render the comparison only after all three legs validate:
    `python3 scripts/campaign.py report --run-dir <run> --config <config>`.
12. Create a reviewable result bundle with
    `python3 scripts/campaign.py sanitize --run-dir <run> --output <bundle.tar.gz>`.
13. If cleanup was selected, preview and then execute
    `python3 scripts/campaign.py teardown --state <state.json> --dry-run` followed by the same
    command without `--dry-run`.

## Adapter contract

Adapters may be maintained in their authoritative SDK repositories, but a run is comparable only
when all three emit the same contract. `metrics.json` must include operation counts, errors,
retries, throttles, wall and CPU seconds, latency percentiles, queue delay when available, RSS/HWM
and anonymous/file-backed memory, cgroup current/peak, workload fields, and exposed queue/in-flight/
completion high-water gauges. `perf-self.csv` contains `symbol,dso,self_percent`.

Use the Rust repository's `sdk/cosmos/azure_data_cosmos_perf` and observability/native-driver
tooling as the Rust source of truth, but do not use its default mixed workload unless it has been
configured or adapted to the exact 5:2 schedule. Use Go v1 from a recorded stable module version.
Use Go v2 from a recorded commit containing diagnostics fix
`7201f4b70fac1184e66379278566f1ac8d299109`.

## Accounting rules

- `perf` Self rows are a composition and must normalize to 100% within the configured tolerance.
- Every row must map through `classifiers/cpu-categories.json`; unmatched rows fail validation.
- Category share of total four-vCPU capacity is:
  `self_share * process_cpu_seconds / (measurement_wall_seconds * 4)`.
- Keep process CPU, Go pprof, and `perf` percentages as separate denominators. Never add them.
- Required categories are Go runtime/SDK, cgo/completion FFI, Tokio, HTTP/2+HPACK, TLS,
  serialization/body/allocation/memcopy, AAD bearer auth/headers, Cosmos
  routing/session/transport, kernel scheduling/syscalls, kernel network/memory, resolved-other,
  and residual-unknown.
- Fail report rendering for schema errors, nonzero retry/throttle counts, workload mismatch,
  missing raw evidence, unmapped symbols, or accounting drift.

## Secret handling

Do not place tokens, keys, connection strings, environment dumps, command histories, home
directories, subscription display names, account names, endpoints, database/container names, or
public IPs in a submitted bundle. The sanitize command retains hashed resource identifiers and
removes known sensitive fields, but inspect the bundle before sharing it.
