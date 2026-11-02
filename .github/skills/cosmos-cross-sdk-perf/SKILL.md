---
name: cosmos-cross-sdk-perf
description: Use when provisioning, running, validating, or aggregating the matched Cosmos DB Go v1, Go v2 native-driver, and Rust SDK performance campaign.
---

# Cosmos Cross-SDK Performance Campaign

This skill reproduces a matched, closed-loop Cosmos DB point-read comparison. Read
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
   available four-vCPU Linux VM SKU, and writes ownership state. A managed identity may administer
   campaign storage, but measured SDK legs use the Cosmos account key.
6. On the VM, run `python3 scripts/campaign.py prerequisites --phase vm`, then build each adapter
   with profiling-quality symbols. The native Rust driver build must
   include `-C debuginfo=2`, `-C force-frame-pointers=yes`, and an ELF build ID, with stripping
   disabled. Build its pinned source-of-record commit from `Azure/azure-sdk-for-rust`, and record
   the pinned `Azure/azure-cosmos-driver` distribution commit consumed by Go v2. Record the complete
   build provenance in each leg's `provenance.json`.
7. Seed exactly 1,024 items outside measurement. Their serialized UTF-8 JSON payloads are exactly
   1,024 bytes and the partition key path is `/pk`. Cleanup writes are also outside measurement.
8. Before every measured process, prove classic Gateway V1, negotiated HTTP/2, and binary encoding
   disabled. Fail before warmup if any proof is absent. Every timed operation is a point read using
   the deterministic round-robin key sequence. Adapters consume raw response bytes and perform
   equivalent minimal validation outside the timed hot path.
9. Run the workers in this order to reduce temporal bias: workers 8 uses Go v1, Go v2/native,
   Rust; workers 16 uses Rust, Go v2/native, Go v1. Use one long-lived runtime/client/pool per
   process, a 30-second excluded warmup, and a 300-second measurement.
10. Keep run classes separate:
    - `clean-score`: at least two A/A repeats, external process/cgroup sampling only. No perf,
      pprof, runtime trace, uprobes, or new hot-path diagnostic counters.
    - `perf`: Linux perf only with profiling-quality optimized binaries. Calibrate frame-pointer
      versus DWARF first and use one passing mode consistently. No Go pprof.
    - `go-diagnostic`: Go v1/v2 CPU and heap pprof plus short runtime trace, block, and mutex
      evidence. These metrics are never authoritative scores.
11. Use `scripts/run-clean-leg.sh`, `scripts/profile-leg.sh`, and
    `scripts/run-go-diagnostic-leg.sh` for their respective classes. Keep every `raw/` directory
    immutable.
12. Each adapter writes `metrics.json` matching the schema; perf runs additionally write
    `perf-self.csv`. Run
    `python3 scripts/campaign.py validate-leg --leg-dir <leg> --config <config>`.
13. Generate whole-process, CPU-per-read, forward/reverse differential, and separately labeled
    composition-normalized flamegraphs with `scripts/generate-flamegraphs.sh`.
14. Render the comparison only after the complete workers 8 and 16 pilot validates:
    `python3 scripts/campaign.py report --run-dir <run> --config <config>`.
15. Create a reviewable result bundle with
    `python3 scripts/campaign.py sanitize --run-dir <run> --output <bundle.tar.gz>`.
16. If cleanup was selected, preview and then execute
    `python3 scripts/campaign.py teardown --state <state.json> --dry-run` followed by the same
    command without `--dry-run`.

## Adapter contract

Adapters may be maintained in their authoritative SDK repositories, but a run is comparable only
when all three emit the same contract. `metrics.json` must include operation counts, errors,
retries, throttles, wall and CPU seconds, latency percentiles, queue delay when available, RSS/HWM
and anonymous/file-backed memory, cgroup current/peak, workload fields, and exposed queue/in-flight/
completion high-water gauges. `perf-self.csv` contains `symbol,dso,self_percent`.

Use the Rust repository's `sdk/cosmos/azure_data_cosmos_perf` and observability/native-driver
tooling as the Rust source of truth, but adapt it to raw point reads only. Use Go v1 from a recorded stable module version.
Use Go v2 from a recorded commit containing diagnostics fix
`7201f4b70fac1184e66379278566f1ac8d299109`.

## Accounting rules

- `perf` Self rows are a composition and must normalize to 100% within the configured tolerance.
- Every row must map through `classifiers/cpu-categories.json`; unmatched rows fail validation.
- Category share of total four-vCPU capacity is:
  `self_share * process_cpu_seconds / (measurement_wall_seconds * 4)`.
- Keep process CPU, Go pprof, and `perf` percentages as separate denominators. Never add them.
- Required categories are Go runtime/SDK, cgo/completion FFI, Tokio, HTTP/2+HPACK, TLS,
  serialization/body/allocation/memcopy, shared-key signing/HMAC/headers, Cosmos
  routing/session/transport, kernel scheduling/syscalls, kernel network/memory, resolved-other,
  and residual-unknown.
- Fail report rendering for schema errors, nonzero retry/throttle counts, workload mismatch,
  missing raw evidence, unmapped symbols, or accounting drift.

## Secret handling

Retrieve the Cosmos key immediately before each measured process, pass it without a command-line
argument, and unset it immediately afterward. Do not place tokens, keys, connection strings,
environment dumps, command histories, home
directories, subscription display names, account names, endpoints, database/container names, or
public IPs in a submitted bundle. The sanitize command retains hashed resource identifiers and
removes known sensitive fields, but inspect the bundle before sharing it.

Stop after the workers 8 and 16 pilot report. Do not run larger matrices or Ashley-specific
second-pass CQ/FFI/copy/Tokio counters without explicit user approval after first-pass review.
