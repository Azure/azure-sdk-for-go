# Azure Cosmos DB Performance Testing Tool

A CLI tool for performance and scale testing the Azure Cosmos DB Go SDK. It runs
point reads, single-partition queries, read-many calls, change feed reads,
upserts, and creates concurrently and reports latency statistics at configurable
intervals.

## Prerequisites

- Go toolchain (Go 1.24 or later recommended)
- An Azure Cosmos DB account
  - The database, container, and results container will be created automatically if they don't exist (with `/partition_key` as the partition key path)
- For key auth: a Cosmos DB account key
- For AAD auth: workload identity or managed identity available in the hosting environment (e.g., AKS, Azure VM, App Service)

## Building

From this package directory:

```bash
go build ./cmd/cosmos-perf
```

## Usage

### Key Authentication

```bash
go run ./cmd/cosmos-perf -- \
  --endpoint https://<account>.documents.azure.com:443/ \
  --auth key \
  --key "<your-account-key>" \
  --application-region "East US 2"
```

Or use the `AZURE_COSMOS_KEY` environment variable:

```bash
export AZURE_COSMOS_KEY="<your-account-key>"
go run ./cmd/cosmos-perf -- \
  --endpoint https://<account>.documents.azure.com:443/ \
  --auth key \
  --application-region "East US 2"
```

### AAD (Entra ID) Authentication

Uses `WorkloadIdentityCredential` first, then falls back to
`ManagedIdentityCredential`, which requires the tool to run in an Azure
environment with an assigned identity.

```bash
go run ./cmd/cosmos-perf -- \
  --endpoint https://<account>.documents.azure.com:443/ \
  --auth aad \
  --application-region "East US 2"
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--endpoint` | *required* | Cosmos DB account endpoint URL |
| `--database` | `perfdb` | Database name |
| `--container` | `perfcontainer` | Container name (partition key path must be `/partition_key`) |
| `--auth` | *required* | Authentication method: `key` or `aad` |
| `--key` | — | Account key (or set `AZURE_COSMOS_KEY` env var) |
| `--application-region` | *required* | Azure region where the application is running; it is the first preferred region |
| `--preferred-regions` | — | Additional comma-separated preferred regions appended after `--application-region` |
| `--excluded-regions` | — | Comma-separated excluded regions |
| `--exclude-regions-for` | `both` | Scope for excluded regions: `reads`, `writes`, or `both` |
| `--concurrency` | `50` | Number of concurrent operations |
| `--duration` | indefinite | Run duration in seconds |
| `--seed-count` | `1000` | Number of items to pre-seed |
| `--throughput` | `100000` | Throughput (RU/s) when creating the container |
| `--default-ttl` | `3600` | Default TTL in seconds for items (0 to disable) |
| `--report-interval` | `300` | Stats reporting interval in seconds |
| `--results-container` | `perfresults` | Container for storing perf results and error documents |
| `--results-endpoint` | — | Cosmos DB endpoint for results (omit to use same account as `--endpoint`) |
| `--results-database` | `perfdb` | Database name on the results account |
| `--results-auth` | same as `--auth` | Authentication method for the results account: `key` or `aad` |
| `--results-key` | — | Account key for results account (or set `AZURE_COSMOS_RESULTS_KEY` env var) |
| `--workload-id` | random UUID | Unique identifier for this workload instance (for multi-VM correlation) |
| `--commit-sha` | auto-detected | Git commit SHA stamped on result documents (auto-detected from `git rev-parse --short HEAD` if omitted) |
| `--no-reads` | `false` | Disable point read operations |
| `--no-queries` | `false` | Disable query operations |
| `--no-upserts` | `false` | Disable upsert operations |
| `--no-creates` | `false` | Disable create operations |
| `--no-read-many` | `false` | Disable read-many operations |
| `--no-change-feed` | `false` | Disable incremental change feed operations |
| `--read-many-batch-size` | `20` | Number of items per ReadManyItems call |
| `--change-feed-max-items` | `100` | MaxItemCount for change feed pages |

### Examples

Run reads only with 100 concurrent operations for 60 seconds:

```bash
go run ./cmd/cosmos-perf -- \
  --endpoint https://myaccount.documents.azure.com:443/ \
  --auth key --key "$AZURE_COSMOS_KEY" \
  --application-region "East US 2" \
  --no-queries --no-upserts --no-creates --no-read-many --no-change-feed \
  --concurrency 100 --duration 60 --report-interval 10
```

Run all operations with application region and custom database:

```bash
go run ./cmd/cosmos-perf -- \
  --endpoint https://myaccount.documents.azure.com:443/ \
  --database mydb --container mycontainer \
  --auth aad \
  --application-region "East US 2" \
  --concurrency 200 --seed-count 5000
```

## Output

The tool prints periodic latency summaries like:

```text
--- Interval Report ---
  Process: CPU 45.2%, Memory 128.3 MB
  System:  CPU 12.8%, Memory 5.2 GB/16.0 GB
  Operation         Count   Errors        Min        Max       Mean        P50        P90        P99 BackendP99
  ------------------------------------------------------------------------------------------------------------------
  ChangeFeedItems     120        0      2.0ms     25.2ms      8.8ms      6.5ms     18.0ms     23.1ms      4.2ms
  CreateItem          280        0      4.0ms     55.2ms     16.8ms     12.5ms     35.0ms     50.1ms     11.7ms
  QueryItems          312        0      3.2ms     45.1ms     12.4ms      9.8ms     28.3ms     41.2ms      9.3ms
  ReadItem            298        2      1.8ms     38.7ms      8.2ms      6.1ms     19.5ms     35.4ms      6.4ms
  ReadManyItems       160        0      3.8ms     60.0ms     19.0ms     14.1ms     39.2ms     55.0ms        —
  UpsertItem          325        0      4.5ms     52.3ms     15.1ms     11.2ms     32.1ms     48.7ms     10.8ms
```

### Results Container

Periodic summary documents and individual error documents are written to the
results container (`--results-container`, default `perfresults`).

- **Summary documents**: Upserted at each reporting interval with latency
  percentiles, process metrics (CPU/memory), system metrics (CPU/memory), and
  workload ID per operation.
- **Error documents**: Written for each individual operation failure with the
  operation name, error message, source error chain, workload ID, and timestamp.
  Errors during the perf run never stop the workload — they are captured and
  reported but execution continues.

If the tool cannot write a result or error document (e.g., the results container
is temporarily unavailable), a warning is printed to stderr and the workload
continues unaffected.

### Separate Results Account

By default, results are stored on the same account being tested. To avoid adding
noise to your workload, use `--results-endpoint` to direct result/error documents
to a different Cosmos DB account:

```bash
go run ./cmd/cosmos-perf -- \
  --endpoint https://workload.documents.azure.com:443/ \
  --auth key --key "$AZURE_COSMOS_KEY" \
  --application-region "East US 2" \
  --results-endpoint https://results.documents.azure.com:443/ \
  --results-auth key --results-key "$AZURE_COSMOS_RESULTS_KEY"
```

### TTL

Containers are created with a default TTL (`--default-ttl`, default 1 hour).
Items automatically expire after this duration, keeping the container from
growing unboundedly during long or repeated runs. Set `--default-ttl 0` to
disable TTL.

### Create Operation

When enabled (the default), the `CreateItem` operation generates new items with
unique IDs and partition keys. Successfully created items are added to the
shared item pool so they become targets for subsequent read, query, read-many,
and upsert operations.

### ReadMany Operation

When enabled (the default), `ReadManyItems` samples `--read-many-batch-size`
items from the shared item pool and calls the Go SDK `ReadManyItems` API. The Go
SDK implementation may internally split work across physical partition ranges;
backend latency is collected best-effort from the SDK pipeline response headers.

### ChangeFeed Operation

When enabled (the default), `ChangeFeedItems` reads one incremental change feed
page per operation execution. On first use it discovers all container feed
ranges, then keeps an in-memory continuation token per feed range. Continuations
are not persisted across process restarts; a restarted process begins again from
the start of each feed range. When a range returns `304 Not Modified`, the tool
clears that range's continuation so future reads restart for that range.

All-versions-and-deletes / full-fidelity change feed mode is out of scope until
the Go SDK exposes a public change feed mode option.

### Multi-Process Launcher

The `run_perf.sh` script launches multiple OS processes of the perf tool in
parallel. This is useful for saturating a Cosmos DB account beyond what a single
process can achieve.

```bash
# Launch 4 parallel processes, each with 50 concurrent goroutines
./run_perf.sh --processes 4 \
  --endpoint https://myaccount.documents.azure.com:443 \
  --auth key --key "$AZURE_COSMOS_KEY" \
  --application-region "East US 2" \
  --concurrency 50 --duration 600

# All standard perf tool flags are passed through to each process
./run_perf.sh --processes 8 \
  --endpoint https://myaccount.documents.azure.com:443 \
  --auth aad --application-region "East US 2" --no-queries --no-creates
```

The script builds the binary, spawns the requested number of processes, and
forwards `Ctrl+C` to all children for graceful shutdown.

The launcher supports these additional flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--processes N` | `1` | Number of OS processes to spawn |
| `--cosmos-commit REF` | — | Build against a specific SDK commit/branch/tag |
| `--poll-branch BRANCH` | — | Continuously poll a remote branch for new commits |
| `--poll-interval SECS` | `43200` (12h) | Seconds between branch polls |
| `--stagger-ms MS` | `200` | Milliseconds between launching each process (0 = simultaneous) |

### Staggered Launch

Use `--stagger-ms` to introduce a delay between starting each process, which can
help avoid thundering-herd effects during container seeding:

```bash
./run_perf.sh --processes 8 --stagger-ms 500 \
  --endpoint https://myaccount.documents.azure.com:443 \
  --auth aad --application-region "East US 2" --concurrency 50
```

### Testing Against a Specific SDK Commit

Use `--cosmos-commit` to build and run against a specific version of
`azcosmos`. This is useful for A/B performance comparisons across SDK changes.

```bash
# Test against a specific commit
./run_perf.sh --cosmos-commit abc123 --processes 4 \
  --endpoint https://myaccount.documents.azure.com:443 \
  --auth key --key "$AZURE_COSMOS_KEY" --application-region "East US 2"

# Test against a branch
./run_perf.sh --cosmos-commit upstream/main --processes 2 \
  --endpoint https://myaccount.documents.azure.com:443 \
  --auth aad --application-region "East US 2"
```

The script checks out the SDK source at the given ref before building, then
restores the original source after the build completes (or on Ctrl+C/error).

### Continuous Branch Polling

Use `--poll-branch` to automatically detect new commits on a remote branch,
rebuild, and restart processes. This is useful for continuous performance
regression testing:

```bash
# Poll for new commits on main, rebuild and restart automatically
./run_perf.sh --poll-branch main \
  --poll-interval 120 --processes 4 \
  --endpoint https://myaccount.documents.azure.com:443 \
  --auth aad --application-region "East US 2"
```

The script fetches the remote every `--poll-interval` seconds (default: 12 hours),
and when a new commit is detected it stops all running processes, checks out the
new code, rebuilds, and restarts. The commit SHA is automatically passed to each
process via `--commit-sha` for Kusto correlation.
