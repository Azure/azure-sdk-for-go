# Cosmos Go v2 FFI flow harness (Nov 5 campaign)

This is the exact harness behind the Nov 5 Go v2 CPU analysis. It is a review snapshot, not a supported tool. For the reusable workflow, see the `cosmos-cross-sdk-perf` skill.

It compares three point-read paths on the same VM, account, and container:

| Leg | What it measures |
| --- | --- |
| `go-v2-native` | Go v2 calling the Rust native driver through cgo |
| `rust-driver` | The Rust driver called directly |
| `rust-sdk` | The standalone Rust SDK |

## Layout

| Folder | Purpose | Needed for a smoke test? |
| --- | --- | --- |
| [`bench/`](bench) | Workload spec, the Rust benchmark, source pinning and patching, builds, and the campaign runner | Yes |
| [`instrumentation/`](instrumentation) | Optional FFI lifecycle counters, wire proofs, and their tests | No |
| [`analysis/`](analysis) | Perf Self attribution, calibration, summaries, and provenance | No |
| [`ops/`](ops) | VM plumbing: scratch disk, resume, inspection, export, and input bundling | No |
| [`report/`](report) | Builds the team HTML report from exported results | No |

The folders exist for review only. [`ops/package_inputs.py`](ops/package_inputs.py) bundles everything into one flat directory, and the scripts run from there on the VM. Moving the files into folders didn't change the bundle; its SHA-256 is identical.

## Read these first

1. [`bench/run-spec.json`](bench/run-spec.json): the workload and the gates.
2. [`bench/run_campaign.sh`](bench/run_campaign.sh): the end-to-end order.
3. [`bench/patch_native.py`](bench/patch_native.py): how the native driver is patched, including Tokio threads.
4. [`bench/rust_main.rs`](bench/rust_main.rs): the Rust SDK and Rust driver legs.
5. [`bench/run_leg.py`](bench/run_leg.py): how one leg is launched and sampled.

## Pinned inputs and workload

| Input | Value |
| --- | --- |
| Rust source (`Azure/azure-sdk-for-rust`) | `77ebe786ad1012bc9ac56105e08c05fb4c65b321` |
| Go source (`Azure/azure-sdk-for-go`) | `b20b0fd00ac20c78bdb2d138fb6960b35e428af0` |
| Native ABI | 0.3.0 |
| VM | 4 vCPU Linux, in the same region as the Cosmos account |

The workload:

- 1,024 seeded 1 KiB documents with partition key `/pk`, read round-robin.
- Strict closed loop at 32 and 64 workers.
- 10 s warmup, then 60 s of measurement.
- Gateway V1 over HTTP/2, with binary encoding disabled.
- Zero errors, retries, and throttles required.

## Why Tokio uses two worker threads

All legs use the same runtime shape:

- **Go v2 native driver:** [`bench/patch_native.py`](bench/patch_native.py) (around lines 35–38) inserts `.worker_threads(2)` into the driver's `src/runtime.rs`. The unpatched driver defaults to one worker per core.
- **Rust legs:** [`bench/rust_main.rs`](bench/rust_main.rs) uses `#[tokio::main(flavor = "multi_thread", worker_threads = 2)]`.

## Flow

[`bench/run_campaign.sh`](bench/run_campaign.sh) runs these steps:

1. Check prerequisites.
2. Stage pinned sources, patch them, then build.
3. Record each binary's SHA-256 and build ID.
4. Calibrate perf call chains.
5. Run the run classes, mirroring order to reduce time bias:
   - `clean-score`: external sampling only.
   - `perf`: per-PID perf at 49 Hz.
   - `ffi-flow`: Go v2 lifecycle ledgers.
   - `go-diagnostic`: pprof and a short runtime trace.
6. Summarize and attribute CPU.

Flame graphs come from [`analysis/profile_attribution.py`](analysis/profile_attribution.py), which uses `flamegraph.pl`. Output goes to `run/w{32,64}/perf/repeat-1/<leg>/derived-label-v2/`:

- `whole-process.{folded,svg,html}`
- `cpu-per-read.{folded,svg,html}`
- `cpu-flywheel.svg`

## Tests

The tests assume the flat VM layout. To run them locally, copy the four source folders into one directory and run `python -m unittest test_flow` there. Four tests need the pinned driver sources, which aren't included, so they error locally.

## Not included

- The VM key helper, which passes the Cosmos key on file descriptor 3.
- Local VM-access scripts, generated bundles, and results.

No endpoints, keys, or resource names are included.

## Caveats

Each class ran once per configuration, with two clean-score repeats. Treat the deltas as directional.