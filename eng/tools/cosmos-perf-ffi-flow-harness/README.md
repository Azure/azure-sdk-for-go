# Cosmos Go v2 FFI flow harness (Nov 5 campaign)

This is the exact harness behind the Nov 5 Go v2 CPU analysis. It compares three point-read paths:

- `go-v2-native`: Go v2 calling the Rust native driver through cgo.
- `rust-sdk`: the standalone Rust SDK.
- `rust-driver`: the Rust driver called directly.

The harness is a snapshot for review, not a supported tool. For the reusable, parameterized workflow, see the `cosmos-cross-sdk-perf` skill.

## Pinned inputs

| Input | Value |
| --- | --- |
| Rust source (`Azure/azure-sdk-for-rust`) | `77ebe786ad1012bc9ac56105e08c05fb4c65b321` |
| Go source (`Azure/azure-sdk-for-go`) | `b20b0fd00ac20c78bdb2d138fb6960b35e428af0` |
| Native ABI | 0.3.0 |
| VM | 4 vCPU Linux, in the same region as the Cosmos account |

The workload and gates are defined in [`harness/run-spec.json`](harness/run-spec.json):

- 1,024 seeded 1 KiB documents with partition key `/pk`, read in round-robin order.
- Strict closed loop at 32 and 64 workers.
- 10 s warmup, then 60 s of measurement.
- Gateway V1 over HTTP/2, with binary encoding disabled.
- Zero errors, retries, and throttles required.

## Where the two Tokio worker threads come from

All three legs pin Tokio to two worker threads so they're compared on the same runtime shape:

- **Go v2 native driver:** [`harness/patch_native.py`](harness/patch_native.py) (around lines 35–38) inserts `.worker_threads(2)` into the driver's `src/runtime.rs`. The unpatched driver uses Tokio's default, which is one worker per core.
- **Rust SDK and Rust driver legs:** [`harness/rust_main.rs`](harness/rust_main.rs) uses `#[tokio::main(flavor = "multi_thread", worker_threads = 2)]`.

## Flow

[`harness/run_campaign.sh`](harness/run_campaign.sh) runs the campaign in this order:

1. Prerequisite checks (`check_run_prerequisites.sh`).
2. Stage the pinned sources (`stage_sources.sh`), patch the native driver and Go code (`patch_native.py`, `patch_go.py`, `patch_wire.py`), then build (`build_rust.sh`, `build_go.sh`).
3. Write build provenance: a SHA-256 hash and build ID for every binary (`write_build_provenance.py`).
4. Calibrate perf call chains (`calibrate.py`).
5. Run the separated classes. The clean-score order is mirrored to reduce time bias, and each leg runs through `run_leg.py`.
   - `clean-score`: external sampling only.
   - `perf`: per-PID perf at 49 Hz.
   - `ffi-flow`: Go v2 per-second lifecycle ledgers.
   - `go-diagnostic`: pprof and a short runtime trace.
6. Summarize and attribute the results:
   - `clean_summary.py` and `flow_summary.py`.
   - `analyze_go_diagnostics.py`.
   - `profile_attribution.py`: perf Self accounting, CPU flywheels, and flame graphs.
   - `report.py`.

Flame graphs are generated in `profile_attribution.py` from folded perf stacks using Brendan Gregg's `flamegraph.pl`. Output goes to `run/w{32,64}/perf/repeat-1/<leg>/derived-label-v2/`:

- `whole-process.{folded,svg,html}`
- `cpu-per-read.{folded,svg,html}`
- `cpu-flywheel.svg`

The `report/` scripts build the team HTML report from the exported results.

## Not included

- The VM key helper (`common.sh`). It retrieves the Cosmos key through managed identity and passes it on file descriptor 3.
- Local VM-access scripts, generated input bundles, and results.

No endpoints, keys, or resource names are included.

## Caveats

Each class ran once per configuration, with two clean-score repeats. Treat the deltas as directional.
