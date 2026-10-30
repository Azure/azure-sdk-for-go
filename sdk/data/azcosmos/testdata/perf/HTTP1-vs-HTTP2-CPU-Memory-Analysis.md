# Go v2 gateway transport CPU and memory analysis

> [!NOTE]
> This report records engineering evidence from the specific environment and workload described below. It is not a release performance guarantee.

## Executive summary

The earlier CPU flywheel reported **CPU microseconds per successful operation**, not CPU seconds. This report presents the same experiments in two complementary forms:

1. Each CPU stage as a percentage of the entire four-vCPU VM capacity.
2. CPU microseconds per operation, which makes protocol efficiency comparable across worker counts.

Standard-Gateway HTTP/1.1 was more CPU-efficient than standard-Gateway HTTP/2 throughout the workload progression. Gateway V2 over HTTP/2 introduced a crossover: it was slower at 8-32 workers, where request latency governed closed-loop throughput, but became the best configuration at 64 workers as the four-core VM approached CPU saturation.

At 64 workers, Gateway V2 delivered 9,330 operations/second at 382.2 CPU microseconds per operation. Standard-Gateway HTTP/1.1 delivered 9,098 operations/second at 397.3 microseconds, and standard-Gateway HTTP/2 delivered 7,981 operations/second at 454.2 microseconds.

A fresh Go v1.5.0 reference run remained substantially more efficient than every v2 transport. At 64 workers it delivered 13,895 operations/second at 173.1 CPU microseconds per operation while using only 60.1% of the four-vCPU machine. The transport recommendations in this report therefore describe the best choices **within Go v2**; they do not close the remaining v1-to-v2 efficiency gap.

All Go v2 results in this report include two completion-path optimizations:

- [#27841](https://github.com/Azure/azure-sdk-for-go/pull/27841), `azcosmos/v2: prepare completions on waiting callers`, which is merged. It moves completion preparation and copying from the single queue-consumer path to the waiting operation callers.
- [#27873](https://github.com/Azure/azure-sdk-for-go/pull/27873), `azcosmos/v2: lazily render diagnostics JSON`, which remains open at the time of this report. The benchmark build includes this diagnostics fix, so the reported Go v2 numbers are **after** eager diagnostics JSON rendering was removed from the response hot path.

Memory was captured throughout:

- Whole-process resident memory and high-water RSS, including Go, Rust, TLS/HTTP buffers, native allocators, and runtime state.
- Go heap profiles for cumulative allocation (`alloc_space`) and live heap (`inuse_space`).

At saturation, Gateway V2 used 98.6 MiB high-water RSS: between standard-Gateway HTTP/2 at 90.6 MiB and HTTP/1.1 at 101.9 MiB.

## Experiment

- VM: `Standard_D4s_v5`, four vCPUs
- SDK: Go v2 lazy-diagnostics build
- Completion baseline: includes merged caller-side completion preparation from [#27841](https://github.com/Azure/azure-sdk-for-go/pull/27841)
- Diagnostics baseline: includes the lazy diagnostics implementation from open [#27873](https://github.com/Azure/azure-sdk-for-go/pull/27873)
- Workload: closed loop, 5 Read : 2 PATCH : 2 Upsert
- Worker progression: 8, 16, 32, 64
- Warmup: 30 seconds
- Measurement: 60 seconds
- Binary encoding: disabled
- Standard-Gateway HTTP/2: Gateway V2 disabled, using the Rust driver's sharded HTTP/2 transport
- Standard-Gateway HTTP/1.1: `AZURE_COSMOS_CONNECTION_POOL_IS_HTTP2_ALLOWED=false`
- Gateway V2 HTTP/2: HTTP/2 allowed, Gateway V2 enabled, override unset
- Outcomes: zero failures, zero throttles, one attempt per operation

The HTTP/1.1 switch also changes connection-pool strategy. Both HTTP/2 configurations use sharded transports; HTTP/1.1 uses an unsharded keep-alive pool. This is therefore a comparison of three supported transport configurations, not a protocol-only microbenchmark.

## Outcome progression

| Workers | Transport | Throughput | VM CPU | CPU/op | Mean | p95 | p99 |
|---:|---|---:|---:|---:|---:|---:|---:|
| 8 | Standard Gateway HTTP/1.1 | 1,728 ops/s | 20.8% | 481.4 us | 4.62 ms | 7.79 ms | 10.39 ms |
| 8 | Standard Gateway HTTP/2 | 1,690 ops/s | 25.0% | 590.8 us | 4.73 ms | 7.77 ms | 10.08 ms |
| 8 | Gateway V2 HTTP/2 | 1,427 ops/s | 19.9% | 558.0 us | 5.60 ms | 7.85 ms | 10.05 ms |
| 16 | Standard Gateway HTTP/1.1 | 3,371 ops/s | 40.8% | 484.4 us | 4.74 ms | 7.87 ms | 10.41 ms |
| 16 | Standard Gateway HTTP/2 | 3,307 ops/s | 45.8% | 554.3 us | 4.83 ms | 7.87 ms | 10.53 ms |
| 16 | Gateway V2 HTTP/2 | 2,840 ops/s | 38.4% | 541.3 us | 5.62 ms | 8.91 ms | 10.80 ms |
| 32 | Standard Gateway HTTP/1.1 | 6,479 ops/s | 68.5% | 423.0 us | 4.93 ms | 8.06 ms | 10.55 ms |
| 32 | Standard Gateway HTTP/2 | 6,108 ops/s | 76.2% | 499.0 us | 5.23 ms | 8.27 ms | 10.79 ms |
| 32 | Gateway V2 HTTP/2 | 5,422 ops/s | 62.6% | 462.0 us | 5.90 ms | 9.42 ms | 11.28 ms |
| 64 | Standard Gateway HTTP/1.1 | 9,098 ops/s | 90.4% | 397.3 us | 7.03 ms | 10.65 ms | 12.92 ms |
| 64 | Standard Gateway HTTP/2 | 7,981 ops/s | 90.6% | 454.2 us | 8.01 ms | 12.19 ms | 14.59 ms |
| 64 | Gateway V2 HTTP/2 | 9,330 ops/s | 89.1% | 382.2 us | 6.85 ms | 10.61 ms | 12.64 ms |

| Workers | HTTP/1.1 throughput change | CPU/op change | Mean-latency change |
|---:|---:|---:|---:|
| 8 | +2.3% | -18.5% | -2.3% |
| 16 | +1.9% | -12.6% | -1.9% |
| 32 | +6.1% | -15.2% | -5.7% |
| 64 | +14.0% | -12.5% | -12.3% |

### Gateway V2 crossover

| Workers | Throughput vs standard HTTP/2 | CPU/op vs standard HTTP/2 | Throughput vs standard HTTP/1.1 | CPU/op vs standard HTTP/1.1 |
|---:|---:|---:|---:|---:|
| 8 | -15.6% | -5.6% | -17.5% | +15.9% |
| 16 | -14.1% | -2.3% | -15.8% | +11.8% |
| 32 | -11.2% | -7.4% | -16.3% | +9.2% |
| 64 | +16.9% | -15.9% | +2.6% | -3.8% |

Gateway V2 is CPU-efficient throughout, but its dispatch path adds request latency. At 8-32 workers the host has idle CPU and closed-loop throughput is latency-bound, so Gateway V2 completes fewer operations. At 64 workers, CPU efficiency becomes dominant and Gateway V2 crosses over both standard-Gateway configurations.

## Go v1 reference baseline

A fresh Go v1.5.0 matrix was captured on the same four-vCPU VM with the same dataset, 5:2:2 workload, worker progression, text JSON, standard Gateway, and HTTP/2. It ran immediately before the detailed Go v2 transport-profile campaign. It is therefore a useful reference baseline, although it was not the same timed invocation as the profile runs below.

| Workers | Throughput | VM CPU | CPU/op | Mean | p95 | p99 | Current RSS |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 8 | 1,889 ops/s | 11.29% | 239.0 us | 4.23 ms | 7.19 ms | 9.58 ms | 27.3 MiB |
| 16 | 3,763 ops/s | 21.22% | 225.6 us | 4.25 ms | 7.27 ms | 9.96 ms | 40.1 MiB |
| 32 | 7,397 ops/s | 39.14% | 211.7 us | 4.32 ms | 7.33 ms | 9.97 ms | 50.0 MiB |
| 64 | 13,895 ops/s | 60.14% | 173.1 us | 4.60 ms | 7.93 ms | 10.70 ms | 76.1 MiB |

Go v1 continued scaling through 64 workers without exhausting the VM. The Go v2 transport configurations approached the CPU ceiling earlier, which is why their throughput flattened and latency rose. The retained v1 matrix did not include a whole-process stage profile or high-water RSS sampling, so v1 is intentionally omitted from the stage-attribution and high-water-memory flywheels.

## CPU flywheel: percentage of total four-vCPU capacity

Every value below is a percentage of the VM's entire four-vCPU capacity. The sum of all categories plus idle capacity is 100% for each workload.

```text
Total VM CPU progression

Go v1      8  ██████                                             11.29%
           16 ███████████                                        21.22%
           32 ████████████████████                               39.14%
           64 ██████████████████████████████                     60.14%

HTTP/1.1   8  ██████████                                         20.80%
           16 ████████████████████                               40.82%
           32 ██████████████████████████████████                 68.53%
           64 █████████████████████████████████████████████      90.36%

HTTP/2     8  ████████████                                       24.96%
           16 ███████████████████████                            45.83%
           32 ██████████████████████████████████████             76.20%
           64 █████████████████████████████████████████████      90.62%

Gateway V2 8  ██████████                                         19.90%
           16 ███████████████████                                38.43%
           32 ███████████████████████████████                    62.62%
           64 █████████████████████████████████████████████      89.14%

              0%          25%          50%          75%         100%
```

### Throughput and mean-latency progression

The throughput and latency panels below use the same 8, 16, 32, and 64-worker progression as the CPU flywheel. Their scales are independent: longer is better for throughput, while a shorter latency bar is better.

```text
Successful throughput progression

Go v1      8  ███████                                            1,889 ops/s
           16 █████████████                                      3,763 ops/s
           32 ██████████████████████████                         7,397 ops/s
           64 ██████████████████████████████████████████████████ 13,895 ops/s

HTTP/1.1   8  ██████                                             1,728 ops/s
           16 ████████████                                       3,371 ops/s
           32 ███████████████████████                            6,479 ops/s
           64 █████████████████████████████████                  9,098 ops/s

HTTP/2     8  ██████                                             1,690 ops/s
           16 ████████████                                       3,307 ops/s
           32 ██████████████████████                             6,108 ops/s
           64 █████████████████████████████                      7,981 ops/s

Gateway V2 8  █████                                              1,427 ops/s
           16 ██████████                                         2,840 ops/s
           32 ███████████████████                                5,422 ops/s
           64 █████████████████████████████████                  9,330 ops/s

Mean SDK latency progression — lower is better

Go v1      8  █████████████████████                              4.23 ms
           16 █████████████████████                              4.25 ms
           32 ██████████████████████                             4.32 ms
           64 ███████████████████████                            4.60 ms

HTTP/1.1   8  ███████████████████████                            4.62 ms
           16 ████████████████████████                           4.74 ms
           32 █████████████████████████                          4.93 ms
           64 ███████████████████████████████████                7.03 ms

HTTP/2     8  ████████████████████████                           4.73 ms
           16 ████████████████████████                           4.83 ms
           32 ██████████████████████████                         5.23 ms
           64 ████████████████████████████████████████           8.01 ms

Gateway V2 8  ████████████████████████████                       5.60 ms
           16 ████████████████████████████                       5.62 ms
           32 ██████████████████████████████                     5.90 ms
           64 ██████████████████████████████████                 6.85 ms
```

This view makes the crossover visible: Gateway V2 trails at lower worker counts, then overtakes both standard-Gateway v2 transports at 64 workers while retaining lower mean latency than either v2 alternative.

### HTTP/1.1

| CPU stage | 8 workers | 16 workers | 32 workers | 64 workers |
|---|---:|---:|---:|---:|
| Native transport pipeline | 3.14% | 6.51% | 12.30% | 18.21% |
| Signing + Cosmos headers | 1.88% | 4.29% | 8.25% | 11.87% |
| Completion queue + FFI | 2.22% | 4.47% | 7.26% | 9.83% |
| Operation planning/pipeline | 2.12% | 4.10% | 7.08% | 9.64% |
| Kernel scheduling/syscalls | 3.22% | 5.55% | 7.58% | 8.67% |
| Go application/runtime | 2.91% | 5.37% | 7.93% | 8.52% |
| TLS/crypto | 1.33% | 2.56% | 4.69% | 6.18% |
| HTTP/1 protocol | 0.94% | 2.11% | 3.92% | 5.14% |
| Kernel memory/other | 0.74% | 1.40% | 2.31% | 3.36% |
| Kernel network | 0.41% | 0.91% | 1.52% | 2.29% |
| Diagnostics snapshot/render | 0.50% | 0.99% | 1.65% | 2.10% |
| Routing + session state | 0.46% | 0.86% | 1.39% | 1.90% |
| Async runtime | 0.67% | 1.24% | 1.80% | 1.44% |
| HTTP client/protocol other | 0.08% | 0.12% | 0.30% | 0.39% |
| Other/unresolved | 0.15% | 0.27% | 0.37% | 0.62% |
| Native/libc remainder | 0.02% | 0.04% | 0.07% | 0.11% |
| **Profiled process total** | **20.80%** | **40.82%** | **68.53%** | **90.36%** |
| **Idle VM capacity** | **79.20%** | **59.18%** | **31.47%** | **9.64%** |

```text
HTTP/1.1 stage flywheel at 64 workers
Percentage of total four-vCPU capacity

Native transport pipeline     ██████████████████ 18.21%
Signing + Cosmos headers      ████████████       11.87%
Completion queue + FFI        ██████████          9.83%
Operation planning/pipeline   ██████████          9.64%
Kernel scheduling/syscalls    █████████           8.67%
Go application/runtime        █████████           8.52%
TLS/crypto                    ██████              6.18%
HTTP/1 protocol               █████               5.14%
Kernel memory/other           ███                 3.36%
Kernel network                ██                  2.29%
Diagnostics snapshot/render   ██                  2.10%
Routing + session state       ██                  1.90%
Async runtime                 █                   1.44%
Other/unresolved              ▏                   0.62%
HTTP client/protocol other    ▏                   0.39%
Native/libc remainder         ▏                   0.11%
```

### HTTP/2

| CPU stage | 8 workers | 16 workers | 32 workers | 64 workers |
|---|---:|---:|---:|---:|
| HTTP/2 + HPACK | 3.62% | 6.97% | 13.34% | 16.36% |
| Native transport pipeline | 3.27% | 6.59% | 12.33% | 16.06% |
| Signing + Cosmos headers | 1.92% | 3.99% | 7.83% | 10.04% |
| Completion queue + FFI | 2.40% | 4.41% | 7.10% | 8.43% |
| Operation planning/pipeline | 2.27% | 4.04% | 7.08% | 8.30% |
| Go application/runtime | 3.03% | 5.28% | 7.38% | 7.93% |
| Kernel scheduling/syscalls | 3.77% | 5.82% | 7.71% | 7.78% |
| TLS/crypto | 1.23% | 2.50% | 3.91% | 4.93% |
| Kernel memory/other | 0.71% | 1.50% | 2.38% | 2.74% |
| Diagnostics snapshot/render | 0.59% | 0.94% | 1.50% | 1.91% |
| Kernel network | 0.41% | 0.85% | 1.56% | 1.86% |
| Routing + session state | 0.43% | 0.94% | 1.35% | 1.70% |
| Async runtime | 0.98% | 1.39% | 1.71% | 1.44% |
| HTTP client/protocol | 0.14% | 0.34% | 0.58% | 0.61% |
| Other/unresolved | 0.16% | 0.26% | 0.41% | 0.46% |
| Native/libc remainder | 0.03% | 0.03% | 0.06% | 0.07% |
| **Profiled process total** | **24.96%** | **45.83%** | **76.20%** | **90.62%** |
| **Idle VM capacity** | **75.04%** | **54.17%** | **23.80%** | **9.38%** |

```text
Standard-Gateway HTTP/2 stage flywheel at 64 workers
Percentage of total four-vCPU capacity

HTTP/2 + HPACK                ████████████████ 16.36%
Native transport pipeline     ████████████████ 16.06%
Signing + Cosmos headers      ██████████       10.04%
Completion queue + FFI        ████████          8.43%
Operation planning/pipeline   ████████          8.30%
Go application/runtime        ████████          7.93%
Kernel scheduling/syscalls    ████████          7.78%
TLS/crypto                    █████             4.93%
Kernel memory/other           ███               2.74%
Diagnostics snapshot/render   ██                1.91%
Kernel network                ██                1.86%
Routing + session state       ██                1.70%
Async runtime                 █                 1.44%
HTTP client/protocol          █                 0.61%
Other/unresolved              ▏                 0.46%
Native/libc remainder         ▏                 0.07%
```

### Gateway V2 over HTTP/2

| CPU stage | 8 workers | 16 workers | 32 workers | 64 workers |
|---|---:|---:|---:|---:|
| Native transport pipeline | 3.30% | 6.46% | 11.92% | 20.26% |
| Completion queue + FFI | 2.06% | 3.86% | 6.64% | 10.20% |
| Operation planning/pipeline | 1.63% | 3.56% | 6.17% | 9.72% |
| Go application/runtime | 2.70% | 5.07% | 8.17% | 9.52% |
| Kernel scheduling/syscalls | 3.58% | 6.54% | 8.78% | 9.18% |
| HTTP/2 + HPACK | 1.43% | 3.06% | 4.83% | 6.33% |
| TLS/crypto | 1.02% | 1.87% | 3.11% | 4.70% |
| Signing + Cosmos headers | 0.70% | 1.44% | 2.73% | 4.50% |
| Kernel memory/other | 0.72% | 1.50% | 2.54% | 3.94% |
| Kernel network | 0.43% | 0.87% | 1.42% | 2.43% |
| Diagnostics snapshot/render | 0.48% | 0.97% | 1.51% | 2.24% |
| Routing + session state | 0.45% | 0.78% | 1.38% | 2.22% |
| Async runtime | 1.00% | 1.73% | 2.26% | 2.06% |
| HTTP client/protocol | 0.22% | 0.42% | 0.68% | 1.02% |
| Other/unresolved | 0.16% | 0.27% | 0.45% | 0.74% |
| Native/libc remainder | 0.02% | 0.04% | 0.04% | 0.08% |
| **Profiled process total** | **19.90%** | **38.43%** | **62.62%** | **89.14%** |
| **Idle VM capacity** | **80.10%** | **61.57%** | **37.38%** | **10.86%** |

```text
Gateway V2 HTTP/2 stage flywheel at 64 workers
Percentage of total four-vCPU capacity

Native transport pipeline     ████████████████████ 20.26%
Completion queue + FFI        ██████████           10.20%
Operation planning/pipeline   ██████████            9.72%
Go application/runtime        ██████████            9.52%
Kernel scheduling/syscalls    █████████             9.18%
HTTP/2 + HPACK                ██████                6.33%
TLS/crypto                    █████                 4.70%
Signing + Cosmos headers      █████                 4.50%
Kernel memory/other           ████                  3.94%
Kernel network                ██                    2.43%
Diagnostics snapshot/render   ██                    2.24%
Routing + session state       ██                    2.22%
Async runtime                 ██                    2.06%
HTTP client/protocol          █                     1.02%
Other/unresolved              █                     0.74%
Native/libc remainder         ▏                     0.08%
```

## CPU interpretation

The HTTP protocol stage is the differentiator:

| CPU/op | 8 workers | 16 workers | 32 workers | 64 workers |
|---|---:|---:|---:|---:|
| HTTP/1 protocol | 21.7 us | 25.0 us | 24.2 us | 22.6 us |
| HTTP/2 + HPACK | 85.7 us | 84.3 us | 87.4 us | 82.0 us |
| HTTP/1.1 saving | 64.0 us | 59.3 us | 63.2 us | 59.4 us |

Gateway V2 changes where the HTTP/2 work occurs:

| CPU/op at 64 workers | Standard Gateway HTTP/2 | Gateway V2 HTTP/2 |
|---|---:|---:|
| HTTP/2 + HPACK | 82.0 us | 27.1 us |
| Signing + Cosmos headers (profile attribution) | 50.3 us | 19.3 us |
| Native transport pipeline | 80.5 us | 86.9 us |
| Completion queue + FFI | 42.3 us | 43.7 us |
| Total CPU/op | 454.2 us | 382.2 us |

Gateway V2 introduces `gateway_v2_dispatch` wrapping and unwrapping, increasing native transport work slightly. The normal Cosmos request is still signed before wrapping, so Gateway V2 does **not** eliminate account-key HMAC or AAD token acquisition. It moves the resulting authorization, date, routing, consistency, session, continuation, query, and partition metadata into typed RNTBD tokens inside the HTTP DATA payload. The smaller outer HTTP header set reduces HPACK and HTTP header-map processing; some CPU previously attributed to the broad signing/header category also moves into native dispatch encoding. Whole-process profiles contain `gateway_v2_dispatch::wrap_request_for_gateway_v2`, `gateway_v2_dispatch::unwrap_response_for_gateway_v2`, and `ShardedHttpTransport::send`, confirming the intended path was active.

The safe conclusion is that Gateway V2 reduces **outer HTTP representation and processing** for signed Cosmos metadata. The profile does not show that cryptographic signing itself was removed.

HTTP/2 sharding still provides useful connection-level concurrency:

- independent connections and HPACK state;
- independent flow-control windows;
- independent TCP congestion and recovery state;
- distribution of send/receive work;
- reduced impact from a single connection failure.

However, sharding cannot eliminate the per-request HPACK and frame-processing cost. On this four-core workload, the CPU cost outweighed the multiplexing benefit.

The completion queue did not become the choke point. At 64 workers it used:

- HTTP/1.1: 9.83% of total VM capacity, approximately 0.39 core.
- Standard Gateway HTTP/2: 8.43% of total VM capacity, approximately 0.34 core.
- Gateway V2 HTTP/2: 10.20% of total VM capacity, approximately 0.41 core.

The queue remained well below one saturated core.

## Memory flywheel

### Whole-process resident memory

High-water RSS includes Go and Rust heaps, native allocators, TLS and HTTP buffers, stacks, and other resident process pages.

| Workers | Standard HTTP/1.1 current/high-water | Standard HTTP/2 current/high-water | Gateway V2 HTTP/2 current/high-water |
|---:|---:|---:|---:|
| 8 | 39.3 / 39.3 MiB | 36.6 / 38.5 MiB | **33.4 / 35.5 MiB** |
| 16 | 47.6 / 50.8 MiB | 45.1 / 48.7 MiB | **44.9 / 45.8 MiB** |
| 32 | 63.8 / 73.1 MiB | 66.0 / 66.9 MiB | **65.9 / 65.9 MiB** |
| 64 | 92.7 / 101.9 MiB | **79.9 / 90.6 MiB** | 98.6 / 98.6 MiB |

```text
High-water RSS progression

HTTP/1.1   8 ████████████████████                              39.3 MiB
           16 █████████████████████████                         50.8 MiB
           32 ████████████████████████████████████              73.1 MiB
           64 █████████████████████████████████████████████████ 101.9 MiB

HTTP/2     8 ███████████████████                               38.5 MiB
           16 ████████████████████████                          48.7 MiB
           32 █████████████████████████████████                 66.9 MiB
           64 █████████████████████████████████████████████     90.6 MiB

Gateway V2 8 ██████████████████                                35.5 MiB
           16 ███████████████████████                           45.8 MiB
           32 █████████████████████████████████                 65.9 MiB
           64 ███████████████████████████████████████████████   98.6 MiB
```

At 64 workers, Gateway V2 used 8.0 MiB more high-water RSS than standard-Gateway HTTP/2, but 3.3 MiB less than HTTP/1.1. Gateway V2 therefore sits between the two standard-Gateway configurations in resident memory while delivering the highest saturated throughput.

### Go allocation and live-heap progression

`alloc_space` is cumulative Go-managed allocation over the profiled run. It is allocation churn, not retained memory. `inuse_space` is the sampled live Go heap at profile capture time.

| Workers | Standard HTTP/1.1 allocation/live | Standard HTTP/2 allocation/live | Gateway V2 HTTP/2 allocation/live |
|---:|---:|---:|---:|
| 8 | 964.7 / 1.50 MiB | 1,050.1 / 1.50 MiB | **913.0 / 1.52 MiB** |
| 16 | 1,957.2 / 0.50 MiB | 2,008.6 / 1.51 MiB | **1,683.6 / 0.50 MiB** |
| 32 | 3,660.7 / 1.00 MiB | 3,706.9 / 1.52 MiB | **3,251.4 / 2.01 MiB** |
| 64 | 5,481.4 / 2.52 MiB | **4,889.2 / 1.00 MiB** | 5,614.0 / 2.02 MiB |

```text
Cumulative Go allocation during the profile

HTTP/1.1   8  █████████                                  0.96 GiB
           16 ██████████████████                          1.91 GiB
           32 █████████████████████████████████           3.57 GiB
           64 ███████████████████████████████████████████ 5.35 GiB

HTTP/2     8  ██████████                                 1.03 GiB
           16 ███████████████████                         1.96 GiB
           32 ██████████████████████████████████          3.62 GiB
           64 █████████████████████████████████████████   4.77 GiB

Gateway V2 8  █████████                                   0.89 GiB
           16 ████████████████                            1.64 GiB
           32 ██████████████████████████████               3.18 GiB
           64 ████████████████████████████████████████████ 5.48 GiB
```

Gateway V2 and HTTP/1.1 complete more operations at saturation, so their cumulative allocation is correspondingly higher. Normalized allocation per completed operation is the more useful comparison:

| Workers | Standard HTTP/1.1 allocation/op | Standard HTTP/2 allocation/op | Gateway V2 allocation/op |
|---:|---:|---:|---:|
| 8 | 9.53 KiB | 10.60 KiB | 10.92 KiB |
| 16 | 9.91 KiB | 10.37 KiB | 10.12 KiB |
| 32 | 9.64 KiB | 10.36 KiB | 10.24 KiB |
| 64 | 10.28 KiB | 10.45 KiB | 10.27 KiB |

At 64 workers, Gateway V2 and HTTP/1.1 have nearly identical Go allocation per operation and are both slightly lower than standard-Gateway HTTP/2. The differences in resident memory therefore appear primarily in native transport, connection, buffer, and runtime state rather than retained Go heap.

## Conclusions

1. These Go v2 results already include caller-side completion preparation from #27841 and lazy diagnostics from #27873. The remaining gap is therefore downstream of those improvements, not a measurement of the earlier eager-diagnostics or single-consumer-copying implementations.
2. Go v1 remains the most CPU-efficient implementation and continues scaling after the v2 paths approach the four-core ceiling. At 64 workers it delivered approximately 49% more throughput than Gateway V2 while consuming less total CPU.
3. Within Go v2, standard-Gateway HTTP/1.1 is 12.5-18.5% cheaper in CPU per operation than standard-Gateway HTTP/2.
4. Gateway V2 is the most CPU-efficient saturated v2 configuration: 382.2 microseconds per operation at 64 workers.
5. Gateway V2 has a latency-driven crossover. It is 11-16% slower at 8-32 workers, but 17% faster than standard-Gateway HTTP/2 and 2.6% faster than HTTP/1.1 at 64 workers.
6. Gateway V2 reduces HTTP/2/HPACK and outer Cosmos-header processing while shifting CPU into native RNTBD dispatch encoding/decoding; client-side signing still occurs.
7. HTTP/2 sharding improves connection-level scalability, but standard Gateway still pays substantial per-request HPACK/frame CPU.
8. The completion queue remains below half a core in every v2 configuration and does not explain the throughput plateau.
9. At 64 workers, high-water RSS is 90.6 MiB for standard HTTP/2, 98.6 MiB for Gateway V2, and 101.9 MiB for HTTP/1.1.
10. Live Go heap remains approximately 1-3 MiB; resident-memory differences are primarily native/runtime state.
11. The best v2 configuration depends on load: standard-Gateway HTTP/1.1 at low-to-medium concurrency, Gateway V2 HTTP/2 near CPU saturation.
12. These results apply to the tested four-vCPU, text-JSON workload. Larger hosts, higher latency, more endpoints, or different payloads may move the crossover.

## Evidence

The summarized measurements used by this report are checked in at
[`data/gateway-transport-results.csv`](data/gateway-transport-results.csv).
Raw performance captures, profiles, binaries, archives, and environment configuration are intentionally not committed.
