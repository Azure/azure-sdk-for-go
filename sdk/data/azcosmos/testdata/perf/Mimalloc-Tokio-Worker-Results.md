# mimalloc and Tokio worker count: Go v2 (native driver) and Rust SDK point reads

## Summary

- **mimalloc cuts CPU per read by 26–29%** at Tokio 2 / W64 (Go v2 285→211 µs, Rust SDK 172→122 µs) and raises throughput by 43–51%.
- **A third Tokio worker is a throughput lever, not an efficiency one.** It adds about 18% throughput at W64 but raises CPU per read.
- **The two changes compound.** Go v2 mimalloc T3 is +81% over the glibc T2 baseline, more than the sum of the separate gains (+43% and +19%).
- **Rust SDK mimalloc T2 matches Go v1 throughput at W64** (16,159 vs 16,405 reads/s) with 26% less CPU per read (122 vs 164 µs).
- **Cost:** mimalloc adds about 20–40 MiB RSS.

## W64 results (mean of 2 repeats)

| Variant | Reads/s | CPU µs/read | Busy cores | p50 / p99 (ms) | p99.9 (ms) | RSS (MiB) |
| --- | ---: | ---: | ---: | --- | ---: | ---: |
| Go v1 | 16,405 | 164 | 2.69 | 3.57 / 8.79 | 37.57 | 37 |
| Go v2 glibc T2 (baseline) | 8,630 | 285 | 2.46 | 7.27 / 12.83 | 43.14 | 56 |
| Go v2 mimalloc T2 | 12,366 | 211 | 2.61 | 5.02 / 9.40 | 25.67 | 75 |
| Go v2 glibc T3 | 10,243 | 315 | 3.22 | 6.11 / 11.11 | 34.78 | 62 |
| Go v2 mimalloc T3 | 15,599 | 216 | 3.37 | 3.89 / 7.84 | 17.14 | 88 |
| Go v2 glibc T4 ¹ | 10,852 | 331 | 3.60 | 5.72 / 10.65 | 35.95 | 65 |
| Go v2 mimalloc T4 ¹ | 16,139 | 228 | 3.67 | 3.76 / 7.68 | 15.92 | 102 |
| Rust SDK glibc T2 (baseline) | 10,671 | 172 | 1.83 | 5.86 / 10.48 | 37.87 | 29 |
| Rust SDK mimalloc T2 | 16,159 | **122** | 1.97 | 3.80 / 7.78 | 17.22 | 66 |
| Rust SDK glibc T3 | 12,534 | 216 | 2.71 | 4.96 / 9.15 | 29.88 | 31 |
| Rust SDK mimalloc T3 | 19,143 | 152 | 2.92 | 3.16 / 7.05 | 18.54 | 76 |
| Rust SDK glibc T4 ¹ | 14,799 | 241 | 3.57 | 4.17 / 8.24 | 22.48 | 32 |
| Rust SDK mimalloc T4 ¹ | **21,608** | 176 | 3.80 | 2.78 / 6.43 | 11.72 | 88 |

¹ The T4 rows come from a second run on the same VM. That run re-measured two drift anchors: Go v2 glibc T2 at 8,537 reads/s (−1.1%) and Rust glibc T2 at 10,777 reads/s (+1.0%).

### Gains relative to each baseline (W64)

| Change | Go v2 reads/s | Go v2 CPU/read | Rust SDK reads/s | Rust SDK CPU/read |
| --- | ---: | ---: | ---: | ---: |
| mimalloc (T2) | +43% | −26% | +51% | −29% |
| Tokio 3 (glibc) | +19% | +11% | +17% | +26% |
| mimalloc + Tokio 3 | +81% | −24% | +79% | −11% |

## W32 results (mean of 2 repeats)

| Variant | Reads/s | CPU µs/read | Busy cores | p50 / p99 (ms) | RSS (MiB) |
| --- | ---: | ---: | ---: | --- | ---: |
| Go v1 | 11,051 | 181 | 2.00 | 2.63 / 6.60 | 30 |
| Go v2 glibc T2 | 7,905 | 304 | 2.41 | 3.88 / 7.88 | 50 |
| Go v2 mimalloc T2 | 10,310 | 237 | 2.44 | 2.88 / 6.89 | 67 |
| Go v2 glibc T3 | 9,466 | 325 | 3.08 | 3.17 / 7.02 | 52 |
| Go v2 mimalloc T3 | 10,958 | 261 | 2.86 | 2.72 / 6.39 | 72 |
| Rust SDK glibc T2 | 9,760 | 185 | 1.81 | 3.09 / 6.85 | 25 |
| Rust SDK mimalloc T2 | 11,497 | 146 | 1.68 | 2.59 / 6.23 | 57 |
| Rust SDK glibc T3 | 10,830 | 230 | 2.49 | 2.73 / 6.56 | 25 |
| Rust SDK mimalloc T3 | 11,845 | 179 | 2.12 | 2.52 / 6.20 | 58 |

## Methodology

- **Host:** one 4-vCPU Linux VM in West US 2, the same region as the Cosmos account. Every variant ran on the same VM against the same account and container.
- **Workload:** 1 KB point reads over seeded keys, closed loop, W32 and W64 workers.
- **Timing:** 10 s warmup, then 60 s measured.
- **Wire:** Gateway V1 over HTTP/2, binary encoding disabled, identical wire configuration across variants.
- **Repeats and ordering:** 2 mirrored repeats (forward then reverse leg order) to cancel drift. Every leg was valid: 36/36 in the T2/T3 run and 24/24 in the T4 run.
- **Sources:** Rust repo at `77ebe786ad1012bc9ac56105e08c05fb4c65b321`; Rust 1.95.
- **Rust build flags:** `-C force-frame-pointers=yes -C debuginfo=2`, `--build-id=sha1`, release builds not stripped.
- **Variants:** each variant differs from the pristine source by only one or both of these patches:
  - mimalloc `=0.1.48` (`default-features = false`) set as `#[global_allocator]`, in the native driver crate for Go v2 and in the bench binary for the Rust SDK;
  - Tokio `worker_threads` changed from 2 to 3 (or 4).
- **Verification:** each staged binary was checked for `mi_malloc` symbol presence, and all variant SHA-256 hashes were confirmed distinct. Provenance JSON records each binary's SHA-256, its build ID, the baseline binary SHA-256 and the variant delta SHA-256.
- **Metrics:** CPU per read is process CPU time divided by successful reads. Busy cores is process CPU divided by wall time. RSS is sampled at the end of the measurement.

## Caveats

- n = 2 per cell on one VM, so treat differences under about 3% as noise.
- These are point reads only, with no upserts or queries.
- Authentication used the account key, passed through a file descriptor and never written to disk or logs. It did not use AAD.
- Go v2 means the Rust native driver plus the Go v2 binding. The allocator change applies to the native library inside the Go process.
- Raw per-leg evidence, binaries and provenance stayed on the campaign VM scratch disk (`/mnt/cosmos-ffi-flow-scratch-20261105-a/alloc-tokio-{iso,t4}-20261106-a`) and are not part of this document.

## Precursor: quick run and cross-run anchors

This quick run (Nov 6) measured only the mimalloc variants. It compared them with glibc baselines and Go v1 taken from other runs that used the same setup: Gateway V1, 4 vCPU, 1 KB point reads, HTTP/2, binary encoding off, and zero errors, retries or throttles. Every row is the mean of 2 repeats. The isolation run above supersedes these numbers because it measured every variant together.

Sources:
- **Nov 6, mimalloc rows:** the quick run (16/16 legs valid, 60 s legs).
- **Nov 6, Go v1 and Rust SDK rows:** the same quick run.
- **Nov 6, Go v2 glibc row:** a separate Nov 6 run (60 s legs).
- **Nov 5 rows:**
  - Go v1 comes from the Tokio 2/3/4 canonical run.
  - Rust SDK and Go v2 glibc come from the latest-main FFI-flow run (120 s legs).

The Nov 5 Go v2 glibc T3/T4 rows were dropped because service telemetry showed them routed through Gateway V2.

### W64

| Variant | Allocator | Tokio | Run | Reads/s | CPU µs/read | Busy cores | p50 (ms) | p99 (ms) | RSS (MiB) |
| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Go v1 | Go | – | Nov 6 | 17,492 | 159 | 2.78 | 3.35 | 8.27 | 43 |
| Go v1 | Go | – | Nov 5 | 18,508 | 153 | 2.83 | 3.15 | 7.93 | 39 |
| Rust SDK ² | glibc | default | Nov 6 | 10,154 | 179 | 1.81 | 6.17 | 10.99 | 29 |
| Rust SDK | glibc | default | Nov 5 | 10,631 | 172 | 1.83 | 5.88 | 10.61 | 25 |
| Go v2 | glibc | 2 | Nov 6 | 8,475 | 287 | 2.43 | 7.41 | 13.06 | 51 |
| Go v2 | glibc | 2 | Nov 5 | 8,293 | 290 | 2.40 | 7.57 | 13.57 | 52 |
| **Go v2** | **mimalloc** | **3** | Nov 6 | **15,582** | **216** | 3.37 | 3.90 | 7.85 | 92 |
| Go v2 | mimalloc | 4 | Nov 6 | 15,940 | 231 | 3.67 | 3.79 | 7.77 | 99 |

### W32

| Variant | Allocator | Tokio | Run | Reads/s | CPU µs/read | Busy cores | p50 (ms) | p99 (ms) | RSS (MiB) |
| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Go v1 | Go | – | Nov 6 | 10,972 | 181 | 1.98 | 2.64 | 6.57 | 30 |
| Go v1 | Go | – | Nov 5 | 12,302 | 170 | 2.09 | 2.36 | 6.16 | 30 |
| Rust SDK | glibc | default | Nov 6 | 9,481 | 189 | 1.79 | 3.20 | 6.97 | 24 |
| Rust SDK | glibc | default | Nov 5 | 9,742 | 186 | 1.81 | 3.12 | 6.78 | 22 |
| Go v2 | glibc | 2 | Nov 6 | 7,747 | 304 | 2.36 | 3.94 | 8.04 | 44 |
| Go v2 | glibc | 2 | Nov 5 | 8,105 | 298 | 2.42 | 3.76 | 7.68 | 44 |
| **Go v2** | **mimalloc** | **3** | Nov 6 | **10,787** | **262** | 2.82 | 2.77 | 6.48 | 72 |
| Go v2 | mimalloc | 4 | Nov 6 | 10,714 | 271 | 2.91 | 2.79 | 6.34 | 85 |

² The two repeats in this cell differ by about 10% (9,660 and 10,649 reads/s).

Go v1 throughput varies by run (16.4k–18.5k reads/s at W64). Compare variants within a single run rather than across runs.

## Outcome

mimalloc became the default global allocator in the native driver in Azure/azure-sdk-for-rust#5579.
