<!-- Copyright (c) Microsoft Corporation. All rights reserved. -->
<!-- Licensed under the MIT License. -->

# Cosmos Go v2: 32-vCPU performance results

Results collected 26-28 October 2026; documentation prepared 1 November 2026.
These are local experimental results, not release or production certification.
This document does not include the separate 4-vCPU study.

## Summary

The caller-side completion optimization removed serialized response preparation
from the Go collector. Its direct before/after experiment was on **16 vCPU**:
at fixed 1,024 callers, observed throughput increased from 16,959 to 27,966
operations/second. No matched original-collector/caller-side A/B was performed
on 32 vCPU, so that gain must not be attributed to the 32-vCPU runs.

On **32 vCPU**, the optimized Go binding sustained 47,754 operations/second at
closed512 and 30K open-loop without demand loss in the binary-wire campaign.
Later shared-harness/text-wire short probes reached approximately 51K at
closed512. The normal-gateway 15-minute comparison matched Go v1 throughput
at open20K and closed128, but v2 consumed substantially more whole-VM CPU.
The higher-concurrency gap was not retested in that normal-gateway campaign.

CPU attribution reproduced approximately 87% VM utilization. A native
frame-pointer diagnostic build provided 99.99% parent-stack coverage:
diagnostics JSON rendering consumed approximately 10.1% of sampled process CPU,
GatewayV2 request wrapping 4.3%, and response unwrapping 3.2%.
These cumulative call-path shares overlap and are not guaranteed throughput
improvements. No functional optimization was implemented during that attribution.

## Conditions and metric definitions

| Condition | Tested configuration |
|---|---|
| Compute | Same West US 2 Standard_D32s_v5 VM; 32 vCPU, 128 GiB RAM |
| CPU/runtime | Intel Xeon Platinum 8370C; CPU affinity 0-31; GOMAXPROCS and Tokio worker count 32 where applicable |
| Software | Linux/amd64; Ubuntu 24.04; Go 1.26.8; local native source pin 62192bf |
| Service capacity | 1M manual shared database RU/s; no dedicated container offer; 100 active physical ranges in sampled campaigns |
| Workload | Point Read, server-side PATCH setting `/value`, Upsert; weights 50:20:20 normalize to 55.56%/22.22%/22.22% |
| Population | Normally 16,384 uniformly selected logical keys, one stable-read and one mutable 1-KiB item per key: 32,768 items plus marker |
| Lifetimes | One long-lived measured client/runtime; temporary detailed transport probe outside measurement where used |
| Operations | Session reads; preferred West US 2; write response content disabled |
| V2 retries/deadlines | Three native 429 retries/5s cumulative wait per transport invocation; native E2E 10s/caller 15s; session/failover retries and hedging explicitly disabled |
| Open loop | Independent scheduled demand, maximum 1,024 callers, queue 4,096, 1ms pacing, 20ms lateness and bounded catch-up |
| Closed loop | Fixed callers; each issues the next operation after the prior completion; no arrival-rate target |
| Isolation | SDK workloads sequential; fresh fixture namespaces; no simultaneous v1/v2/Rust load |
| Profiling | Instrumented windows separated from unprofiled results; normal diagnostics retained except explicitly labeled JSON-only ablation |

Successful throughput counts verified completions inside measurement, excluding
final-drain successes. Successful SDK means/tails include cohort drain and
measure public method entry-to-return, including retries, not generator queueing
or harness validation. Failed-call latency is not included in successful means.
P99 values are histogram upper bounds; different harnesses have different bins.

Whole-VM CPU includes application, harness and other host activity. Peak RSS
covers the process lifecycle, not just steady state. CPU percentages are not
one-core utilization. Per-allocated-vCPU rates divide throughput by 32; they
do not predict a smaller VM's capacity.

Fresh uniform keys do not guarantee identical physical placement or uniform
RU distribution. Shared database metrics can include other activity.
Service-reported processing duration is not network latency or complete native
elapsed time. Missing samples are not zero.

## Before and after caller-side preparation: 16-vCPU background

The original collector copied bodies/headers, translated diagnostics, rendered
JSON and freed native results before signaling each caller. The optimization
transfers exclusive completion ownership to the existing waiting goroutine;
that goroutine prepares and frees the result inside SDK code before returning.
It does not introduce customer callbacks or a goroutine per response.
Cancellation/close ownership barriers remain.

The following comparison is **16 vCPU, 500K shared RU/s, 8,192 items,
1,024 application callers**. SDK means use the first 29 complete unprofiled
seconds of a 60-second original run and corrected 300-second optimized run.
Queue and collector CPU measurements use separate diagnostic windows.

| Metric | Original collector | Caller-side preparation |
|---|---:|---:|
| Successful operations/second | 16,959 | 27,966 |
| SDK Read mean, ms | 58.65 | 34.68 |
| SDK PATCH mean, ms | 62.29 | 38.85 |
| SDK Upsert mean, ms | 62.24 | 38.77 |
| Weighted successful SDK mean, ms | 60.25 | 36.52 |
| Native-ready queue wait: resolution/item, ms | 26.86 / 26.82 | 0.211 / 0.313 |
| Collector CPU, approximate cores | 0.97 | 0.26 |

Observed short-window throughput improved approximately 65%; work was moved
and parallelized, not eliminated. The optimized five-minute run completed
8,366,721 successful cohort calls, zero final errors/drops/cancellations and
full validation/cleanup/drain. Whole-VM CPU reached approximately 90%.
The initial sampled cgo-pointer crash and its recovery are separate failures,
not passing results.

Thirty-minute qualification on the same 16-vCPU/500K setup confirmed original
open15K at 14,999.94/sec, optimized open20K at 19,999.93/sec and optimized
closed128 at 19,301.05/sec. Their Read means were 5.12/4.93/4.99ms respectively.
These differently loaded confirmations are not a matched throughput A/B.

## Resize reference: optimized closed256

| Metric | 16-vCPU reference | 32-vCPU follow-up |
|---|---:|---:|
| Measurement seconds | 900 | 900 |
| Successful operations/second | 27,284.94 | 37,724.99 |
| SDK Read/PATCH/Upsert mean, ms | 7.73 / 11.34 / 11.26 | 5.12 / 8.78 / 8.71 |
| Read p99 upper bound, ms | 13.67 | 9.25 |
| PATCH/Upsert p99 upper bounds, ms | 18.32 / 18.32 | 15.07 / 15.07 |
| Mean whole-VM CPU | 88.87% | 68.88% |
| Peak process RSS, MiB | 169.1 | 178.3 |

The 32-vCPU case completed 900 seconds, with 33,952,745 cohort successes,
zero final errors/cancellations/drops/mismatches/observed429s and full
validation/cleanup/drain. The busiest complete-minute range peaked at 50%.
The CPU count, Go/Tokio workers and CPU-derived maximum H2 shard ceiling
changed together. The 38.26% throughput increase is not CPU-only causality,
linear scaling or lower absolute CPU consumption.

## Binary-wire Go/Rust matrix: 32 vCPU

All nine cases ran 900 seconds after 30-second excluded warmup. The main matrix
used HTTP/2, binary wire and text responses returned to the application.
The Rust runner was a local matched-workload adaptation, not upstream defaults.
The HTTP/1.1 experiment also changed gateway implementation/pool defaults.

| Load | SDK/protocol | Success/sec | SDK Read/PATCH/Upsert mean, ms | Mean VM CPU | Demand loss/final errors |
|---|---|---:|---|---:|---|
| Closed512 | Go v2 H2 | 47,754.02 | 9.03 / 12.72 / 12.64 | 87.67% | None / none |
| Closed512 | Direct Rust H2 | 75,738.78 | 5.06 / 8.86 / 8.79 | 85.55% | None / none |
| Open20K | Go v2 H2 | 19,999.85 | 4.50 / 8.24 / 8.15 | 32.65% | None / none |
| Open20K | Direct Rust H2 | 19,996.48 | 4.25 / 7.91 / 7.83 | 11.40% | 3,026 late arrivals / none |
| Open30K | Go v2 H2 | 29,999.75 | 4.59 / 8.30 / 8.23 | 51.57% | None / none |
| Open30K | Direct Rust H2 | 29,998.13 | 4.32 / 7.97 / 7.89 | 19.66% | 1,453 late arrivals / none |
| Open50K | Go v2 H2 | 47,365.91 | 19.86 / 23.70 / 23.62 | 89.06% | 2,365,587 queue drops (5.25686%) / none |
| Open50K | Direct Rust H2 | 49,999.50 | 4.40 / 8.09 / 8.01 | 45.88% | 71 late arrivals / none |
| Closed256 | Go v2 H1.1 | 37,723.23 | 5.02 / 9.08 / 8.64 | 67.96% | None / two TransportFailure503/20003 |

| Load | SDK/protocol | Read/PATCH/Upsert p99 upper bounds, ms | Peak RSS, MiB |
|---|---|---|---:|
| Closed512 | Go v2 H2 | 16.62 / 22.27 / 22.27 | 310.7 |
| Closed512 | Direct Rust H2 | 8.96 / 22.82 / 22.53 | 101.6 |
| Open20K | Go v2 H2 | 8.81 / 14.35 / 14.35 | 317.9 |
| Open20K | Direct Rust H2 | 8.40 / 13.63 / 13.53 | 135.3 |
| Open30K | Go v2 H2 | 8.81 / 15.07 / 15.07 | 377.2 |
| Open30K | Direct Rust H2 | 8.33 / 13.91 / 13.75 | 142.0 |
| Open50K | Go v2 H2 | 31.33 / 38.09 / 38.09 | 494.9 |
| Open50K | Direct Rust H2 | 8.35 / 15.31 / 15.18 | 138.6 |
| Closed256 | Go v2 H1.1 | 8.81 / 15.83 / 15.07 | Not included in this summary |

The eight H2 cases had zero final errors/cancellations/mismatches/observed429s.
Demand loss remains visible: Go open50K is not a successful 50K offered-rate
pass, and Rust late arrivals are generator loss, not SDK response failures.
The first H1.1 failure was not wire-backed; its status does not establish that
Cosmos returned503.

Across nine cases: 322,725,637 cohort successes, 322,718,208 in-window successes
and two final failures. Every case validated all 32,768 items, deleted its
32,769 owned items/marker and completed native/task drain.
Reported RU including maintenance was 1,969,698,975.28, not billing total.

Complete-minute busiest-range peaks were Go/Rust closed512 64%/95%, open20K
27%/26%, open30K 39%/40%, open50K 60%/62%; H1.1 closed256 peaked at47%.
Rust's busiest closed512 range averaged90.64%. No complete-minute sample
reached100%; shorter spikes are not ruled out. Client efficiency, physical
placement and service work cannot be isolated by total provisioned RU alone.

## Released Go v1 initial reference

| Load | Success/sec | Actual seconds | SDK Read/PATCH/Upsert mean, ms | Mean VM CPU | Outcome |
|---|---:|---:|---|---:|---|
| Closed512 | 72,570.78 | 559.03 | 5.11 / 9.47 / 9.40 | 55.82% | Stopped on cumulative final429 budget; 102 final429s after drain |
| Open20K | 19,999.88 | 900 | 4.01 / 7.89 / 7.80 | 9.80% | One final429 and one500/uncertain write; data retained |
| Open30K | Not run | Not run | Not measured | Not measured | Campaign halted |

These are not clean equivalent v1/v2 comparisons: harness, wire encoding and
retry/gateway semantics differed. Initial v1 did not record backend durations.
Closed512 completed validation/cleanup/drain. Open20K validated/drained but
retained its whole fixture; validation does not establish the failed write's
commit outcome. Final429 count is distinct from recovered attempt429s.

## Shared-harness/text-wire short probes

Each case ran60seconds after30seconds warmup. The table uses the first two
complete unprofiled intervals, approximately20seconds; profiling occupied a
separate window. V2 retained caller-side preparation and GatewayV2.

| Load | SDK | Unprofiled success/sec | SDK Read/PATCH/Upsert mean, ms | Mean VM CPU | Backend mean, ms/attempt |
|---|---|---:|---|---:|---:|
| Open20K | v1 | 19,997.73 | 4.04 / 7.88 / 7.83 | 23.75% | 2.486 |
| Open20K | v2 | 19,996.82 | 4.48 / 8.21 / 8.08 | 36.78% | 3.050 |
| Closed512 | v1 | 70,954.58 | 5.21 / 9.69 / 9.63 | 56.77% | 2.534 |
| Closed512 | v2 | 51,086.99 | 8.38 / 12.05 / 11.94 | 87.56% | 3.041 |

All four60second cohorts were clean: 9,682,866 successes, zero final
errors/cancellations/drops, complete validation/cleanup/drain.
The gap survived shared harness/text-wire controls; GatewayV2 and retry
placement remained different. These short probes do not qualify endurance.

## Diagnostics JSON-only ablation

| V2 closed512 variant | Unprofiled success/sec | SDK Read/PATCH/Upsert mean, ms | Mean VM CPU | Backend mean, ms/attempt |
|---|---:|---|---:|---:|
| Normal JSON control | 50,570.30 | 8.48 / 12.14 / 12.03 | 87.00% | 3.055 |
| Experimental JSON render/copy skipped | 57,338.79 | 7.28 / 10.97 / 10.86 | 86.52% | 3.028 |

Both closed cases completed60seconds cleanly; values use approximately20seconds
unprofiled intervals. Throughput increased13.38% and Read mean fell14.22% in
one sequential pair. Structured/scalar/region/attempt diagnostics, bodies,
headers, errors, ownership and release remained enabled. This experimental
always-empty JSON behavior is not a supported product configuration.

The open20K skip case completed60seconds with1,199,999 successes, one final429
and391retained throttled attempts. Validation/cleanup/drain passed; the final
open control was not run because the zero-error progression gate stopped.
No paired open-loop benefit is established.

## Normal-gateway checks, failures and completed comparison

All use HTTP/2 and text wire; v2 normal Summary diagnostics JSON remained enabled.

| Check/campaign | What ran | Outcome |
|---|---|---|
| First normal-gateway comparison | v1 open20K,900seconds | 19,999.862/sec; means4.058/8.051/7.961ms;11.33%VMCPU; zero final errors/drops;50recovered429 attempts; validated/deleted/drained |
| Same campaign, v2 | Warmup only | One uncertain-write ServiceError500 after0.266seconds admission;1,062started/1,061success/1failure;4,094queued cancellations;134queue drops. No measurement or closed cases; validated/drained, fixture retained |
| Single-item v2 correctness check | Eleven sequential API checks | Create/read/PATCH/upsert/ETag delete/404 and transport checks passed; item deleted; native drain |
| Small v2 closed16 | 128items;30second warmup,60second measurement | 2,862.70/sec; means3.179/8.627/8.481ms; p99 6.576/28.420/28.420ms;4.35%VMCPU;54.46MiB peakRSS; zero failures;128validated/129deleted/drained |
| Fresh v1 comparison | Warmup only | Upsert408 after approximately0.204seconds admission;3,511started/3,510success/1failure;542queued cancellations; no drops/429s. Validated/drained, fixture retained; no remaining cases |
| Staged v1 comparison | First open5K warmup step | Upsert500 after7.933seconds admission;39,660started/39,659success/1failure; no drops/cancels/429s. Validated/drained, fixture retained; no measurement |

Those warmup stops came from the one-uncertain-write guard, not the numeric
100-failure limit or proven SDK capacity. V1 also failing weakens an exclusively
v2 or sudden20K-startup explanation.

The user then approved counting eligible wire-backed408/500 as failed calls
under numeric budgets, without replaying failed writes. Uncertain fixtures
remained retained; auth/cancellation/correctness/artifact/drain safeguards
remained. Three10second warmup stages plus drains were excluded from900seconds.

| Load | SDK | Success/sec | SDK Read/PATCH/Upsert mean, ms | Backend mean, ms/attempt |
|---|---|---:|---|---:|
| Open20K | v1 | 19,999.88 | 3.87 / 7.95 / 7.86 | 2.509 |
| Open20K | v2 | 19,999.88 | 3.75 / 7.59 / 7.49 | 2.276 |
| Closed128 | v1 | 24,144.35 | 3.53 / 7.51 / 7.43 | 2.539 |
| Closed128 | v2 | 24,082.95 | 3.64 / 7.40 / 7.31 | 2.272 |

| Load | SDK | Read/PATCH/Upsert p99 upper bounds, ms | Mean VM CPU | Peak RSS, MiB | Warmup/measurement final errors |
|---|---|---|---:|---:|---|
| Open20K | v1 | 8.81 / 15.07 / 15.07 | 10.82% | 141.08 | One500 / one408 |
| Open20K | v2 | 7.61 / 13.67 / 13.67 | 25.01% | 307.45 | One500 / one408 and one500 |
| Closed128 | v1 | 7.99 / 14.35 / 14.35 | 15.31% | 49.49 | None / none |
| Closed128 | v2 | 7.25 / 13.67 / 13.67 | 42.79% | 124.93 | None / none |

Both tables describe the same latest four900second runs: averages versus tails,
not newer versus older measurements. All cases had zero drops/cancellations/
mismatches/observed429s. Across measured cohorts:79,404,825successes and3failures;
warmup added1,536,966successes and2failures.
Each validated32,768items and drained. Closed cases deleted32,769items/marker
each; open cases retained32,769each. Reported445,664,505.53RU includes maintenance.

Throughput was effectively matched at these points; closed128 v2 was0.254%lower.
Whole-VM CPU was2.31x/2.79x higher for v2 at open/closed. This is one sequential
pair, not statistical equivalence or a resolution of the earlier closed512 gap.

## Backend processing and CPU normalization

V1 captures `x-ms-request-duration-ms`; v2 captures nonnegative per-attempt
`ServerDurationMS` from retained native diagnostics. The reported mean is
total duration divided by valid sample count, including available retry/failure
durations. It is not an independent backend-log measurement or full SDK latency.

| Latest load | SDK | Duration samples / observed attempts |
|---|---|---:|
| Open20K | v1 | 17,999,999 / 18,000,001 |
| Open20K | v2 | 18,016,979 / 18,016,981 |
| Closed128 | v1 | 21,730,043 / 21,730,043 |
| Closed128 | v2 | 21,699,527 / 21,699,527 |

Latest backend per-operation means and p99 were not recorded. The older binary
v2 campaign did record per-operation means: closed512 Read/PATCH/Upsert
1.473/5.104/5.045ms; open20K1.467/5.082/5.019ms. These are historical, not
substitutes for missing latest data. Do not subtract per-attempt backend mean
from per-success SDK mean and label the difference network or FFI overhead.

Latest throughput/allocated-vCPU was approximately625/sec for both SDKs at
open20K; closed128 v1/v2 was754.51/752.59. Open-loop process CPU estimates
exceeded VM-derived busy-CPU estimates by25-33%; absolute busy-core figures
remain provisional. The later aligned high-load CPU windows reconciled much
more closely, but do not retrospectively resolve the earlier open-loop counters.

Zero final429s is not zero retries: the small successful v2 case reported
9,288additional native attempts. Normal native gateway requests disable
gateway-side449RetryWith retry and perform client-side handling; v1 does not
send that opt-out. Matching gateway/H2 does not match every retry-placement detail.
No normal received408-to500SDKrewrite was found; terminal error causes remain
unresolved rather than attributed to service/client without request-level evidence.

## High-CPU attribution campaign

All three v2-only closed512 cases measured60seconds after30second warmup,
GatewayV2/HTTP2, text wire, normal Summary JSON and the same operation mix.
The initial pair used byte-identical executables. The approved follow-up changed
native compilation/instrumentation provenance and is diagnostic, not a stock
or isolated frame-pointer throughput A/B.

| Case | Success/sec | SDK Read/PATCH/Upsert mean, ms | Mean VM CPU | Peak RSS, MiB |
|---|---:|---|---:|---:|
| Baseline | 53,408.25 | 7.934 / 11.632 / 11.519 | 87.08% | 247.5 |
| Instrumented | 50,525.97 | 8.426 / 12.262 / 12.117 | 85.82% | 255.9 |
| Native frame-pointer follow-up | 51,928.38 | 8.193 / 11.923 / 11.805 | 86.64% | 299.7 |

All had zero final errors/cancels/drops, full validation,32,769owned deletions
each and native drain. The initial instrumented throughput difference was-5.4%,
not a measured functional regression or gain.

### Same-run latency, tails and native attempts

| CPU case | Mixed successful SDK mean, ms | Native elapsed mean, ms/call | Backend mean, ms/attempt | Read/PATCH/Upsert p99 upper bounds, ms | Native attempts / SDK calls |
|---|---:|---:|---:|---|---:|
| Baseline | 9.553 | 6.988 | 3.055 | 14.35 / 21.21 / 21.21 | 3,213,226 / 3,205,006 |
| Instrumented | 10.099 | 7.026 | 3.052 | 14.35 / 21.21 / 21.21 | 3,039,548 / 3,032,070 |
| Native frame-pointer follow-up | 9.825 | 7.078 | 3.057 | 15.07 / 21.21 / 21.21 | 3,123,978 / 3,116,215 |

SDK means are weighted by actual successful operation counts. Native elapsed
is driver-reported elapsed time, including native work/waiting; it is not
service-only time. Backend duration is per observed attempt and includes
available retry durations. These different denominators must not be subtracted
to manufacture network or FFI cost. More native attempts than calls despite
zero final errors demonstrates that these runs were not retry-free.

### Aligned process and host CPU accounting

| CPU case, approximately 20-second window | Process CPU-seconds | Process busy-CPU equivalents | Host busy-CPU equivalents | Process/host busy-CPU ratio |
|---|---:|---:|---:|---:|
| Baseline | 544.75 | 27.237 | 27.019 | 1.008 |
| Instrumented native-profile window | 536.60 | 26.830 | 27.062 | 0.991 |
| Native frame-pointer follow-up | 541.19 | 27.059 | 26.939 | 1.004 |

These use process and host counter deltas with separate monotonic read
boundaries. Process counts include Go/cgo/Rust threads; host counts include
other activity and sampling/accounting noise. Their approximately 1% agreement
is evidence for these high-load windows, not a retrospective correction of
the earlier open-loop discrepancy. They are CPU consumption, not CPU assigned
to each individual operation.

Native CPU sampling used20seconds near measurement+5..25seconds; separate Go
CPU/mutex/block/breadcrumb collection was near+30..50seconds. Native formatting
and compression waited until process exit. The baseline aligned CPU window
measured27.24busy CPUs for the process versus27.02host busy CPUs.

The original DWARF capture had50,733samples,0lost, but only23.67%caller-stack
coverage. Missing binutils/addr2line initially broke report formatting; raw
recordings and failed processing verdict were preserved and recovered offline.
Exact matching libc debug symbols were installed without upgrading libc.
Native frame pointers/unwind tables were enabled in a separate diagnostic
build; a no-client/no-service runtime probe showed100%parent coverage before
the approved live follow-up. That follow-up had51,138samples,0lost/empty
callchains and99.988%parent coverage.

| Native FP caller path | Cumulative share of sampled CPU |
|---|---:|
| Diagnostics JSON-render call path | 10.10% |
| GatewayV2 request wrapping | 4.30% |
| GatewayV2 response unwrapping | 3.22% |
| URL parsing/encoding | 2.87% |
| Native item submission | 2.50% |
| Native completion release | 1.35% |
| Attempt/region diagnostics bridge | 0.94% |
| Response-header construction | 0.63% |
| HTTP stack umbrella | 26.34% |
| Tokio runtime umbrella | 68.36% |

These shares overlap. Tokio's umbrella includes useful driver/HTTP/allocator/
kernel work, not68%pure scheduler overhead. Submission and `runtime.cgocall`
include native work, not pure FFI crossing cost.

Separately, exclusive leaf samples attributed18.33%to libc allocation/free,
5.72%to memory/string operations and23.15%to kernel execution. Most allocator
work was on Rust worker threads. The closest named allocator callers included
GatewayV2 wrapping/unwrapping, reqwest send/execute, h2 flushing and token
serialization; this is not exclusively diagnostics allocation.

Separate Go pprof captured287.36CPU-seconds versus approximately483.64process
CPU-seconds in its matching window. Its cumulative JSON path53.79seconds
contained53.18in native JSON rendering; a filtered structured-copy path
excluding JSON had6.40seconds. Do not add nested profiles or present18.7%of
Go-captured CPU as18.7%of all process CPU; normalized JSON cost was approximately11.1%.

The strongest candidates are safe optional/lazy diagnostics materialization and
bounded GatewayV2 intermediate-buffer/metadata-copy reductions. The current
Go-owned exported JSON string cannot execute lazy rendering on a field read.
A deferred native handle needs a deterministic retain/render/free contract;
holding borrowed pointers after completion free or relying only on finalizers
is not safe. No such API or functional optimization was implemented here.

## Build provenance and limitations

Tested Go base: `69428b7fb376cabf2c3647d51a04f3b63fa7fc0e`, with local caller-side
completion preparation and optional performance breadcrumbs. The reports do not
claim byte identity with the evolving caller-completion PR or current feature tip.
Runtime source/dependency fingerprints and raw result archives were verified.

Native pin: `62192bf69dac1a3384505971d7deef19eb8282df`, ABI0.2.0.
The main local archive was71,328,962bytes; published archive31,858,618bytes.
They had different hashes and identical headers; packaged-build functional/
compiler parity was not proved. Later detailed and frame-pointer builds are
distinct fingerprinted artifacts. Diagnostics copying remained enabled in
normal comparisons; only the isolated ablation skipped JSON.

All raw data/profiles and recovery records remain private retained artifacts.
No credentials, request bodies, raw stack dumps or account identifiers are
included in this report. Campaign IDs distinguish results without depending
on local file links:

| Campaign | Evidence covered here |
|---|---|
| go1m32-20261025a | First32-vCPU optimized closed256 |
| http1m-20261026b | Nine15-minute Go/Rust/protocol cases |
| gov1-20261026a | Initial released-v1 partial/error-bearing reference |
| gocmp-20261027a | Shared-harness/text-wire probes |
| gocmp-20261027b | JSON-only ablation and incomplete open pair |
| gocmp-20261027c | First normal-gateway comparison and v2 warmup failure |
| gov2-smoke-20261028a / gov2-small-20261028a | Sequential correctness and light-load v2 checks |
| gocmp-20261028a / gocmp-20261028b | Fresh/staged v1 warmup failures |
| gocmp-20261028c | Completed four-case counted-error comparison |
| cpu-20261028a / cpu-20261028b | Baseline/instrumented CPU pair and high-coverage native follow-up |

Original failed runs and retained uncertain-write fixtures were not erased,
replayed or labeled passing. Performance-test cleanup is limited to known
run-owned items and conditional deletion. Previously retained fixtures remain
untouched. Current capacity/VM state is not inferred from historical results.

These are single-host sequential observations, not randomized repeated trials,
cross-language rankings, multi-host scale tests, leak-freedom proof or a completed
continuous-hour qualification. The earlier hour goal remains unachieved.
This documentation change starts no load and modifies no SDK behavior,
resources, throughput, indexing or data.
