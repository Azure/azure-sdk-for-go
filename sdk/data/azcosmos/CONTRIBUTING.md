# Contributing to azcosmos

> These instructions are for contributors to the `azcosmos` module. To use the module, see the
> [README](README.md). For general guidance, see the repository
> [contributing guide](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md) and
> this module's [AGENTS.md](AGENTS.md).

Emulator-based end-to-end tests are described in the README's
[Running the end-to-end tests](README.md#running-the-end-to-end-tests) section.

## Change-feed contract tests

The scripted HTTP tests in `change_feed_contract_native_test.go` run the published
native driver against controlled responses. They check ownership of raw JSON
envelopes, checkpoint compatibility, and routing across scripted physical ranges.
They do not establish service-side prefix filtering or historical change fidelity.

## Azure Live smoke tests

`ci.yml` keeps unit and emulator validation unchanged and adds an opt-in, Linux amd64
cgo live stage using the published native driver and Microsoft Entra ID. Ordinary PR and
CI runs do not provision live resources or require live credentials. Manual runs opt in
with `RunLiveTests=true`; scheduled runs and definitions ending in `- weekly` select the
same bounded smoke suite. Separate definitions can reuse this `ci.yml`; no separate
`tests.yml` is required. Live execution is restricted to branch builds, not PR merge refs.
Only those live-run contexts opt in to the additional Cosmos outbound-network policy;
ordinary Cosmos validation and other Go pipelines retain the shared policy defaults.

The selection in [`internal/e2e/live.go`](internal/e2e/live.go) contains six shared
scenario IDs: `bootstrap.primary-success`, `item.lifecycle`,
`item.optimistic-concurrency`, `diagnostics.success-and-error`,
`query.pagination-resume`, and `changefeed.pagination-resume`. The two pagination
scenarios select `coreOperations`; the others select `smokeTests`. The change-feed case
checks fresh-client resume after idle, point-in-time filtering, and prefix isolation and
resume within one verified physical range. Multi-range change-feed behavior is not
covered by this lane. Extended, multi-region, and long-running partition-split
profiles are not implemented by this lane, including on scheduled runs.
The pinned catalog is not refreshed. Its single-region East US profiles are explicitly
adapted to East US 2 for live smoke; all other profile requirements remain unchanged.
Logs and metadata identify that region-only adaptation rather than claiming an
unmodified profile run. Unsupported deployment regions, setup requirements, and omitted
required live scenarios fail validation before network I/O.

For local execution, use a dedicated single-region East US 2 Session account and a
database with a `/pk` container and a single-range hierarchical container with paths
`/tenant` and `/user`. The identity needs Cosmos DB
data-plane permissions for metadata, items, and queries at database scope; management
permissions alone are insufficient. Live tests authenticate only with Microsoft Entra ID;
disable local (key) authentication on the account (`disableLocalAuth: true`) so account
keys cannot be used. Tests never create/delete databases or containers.
Do not use a container containing customer data.

Run in a native Linux amd64 or macOS arm64 shell with Go 1.27+, a C compiler,
PowerShell 7, and a development credential:

```powershell
$env:AZURE_RECORD_MODE = 'live'
$env:AZURE_COSMOS_E2E_BACKEND = 'azureLive'
$env:AZCOSMOS_ENDPOINT = 'https://<test-account>.documents.azure.com/'
$env:AZCOSMOS_DATABASE = '<test-database>'
$env:AZCOSMOS_CONTAINER = '<test-container>'
$env:AZCOSMOS_HPK_CONTAINER = '<hierarchical-test-container>'
$env:AZCOSMOS_REGION = 'eastus2'
$env:CGO_ENABLED = '1'
./sdk/data/azcosmos/eng/scripts/Invoke-LiveTests.ps1 -OutputDirectory <temporary-results-directory>
```

Leave `EMULATOR` and `AZURE_COSMOS_E2E_PROFILE/ACCOUNT/RUNTIME/CLIENT` overrides unset.
The runner rejects unsupported and cross-compiled targets, clears stale results, and
runs only `TestAzureLive` with race detection and a 12-minute timeout. Missing credentials,
failed setup, skipped/missing scenarios, and incomplete results fail the run.
`Test-LiveResults.ps1` emits JUnit and JSON with selected scenario IDs, profiles,
durations, and outcomes; `metadata.json` records credential-free build/driver metadata
and both the expected (not independently measured) account configuration and original
catalog region. Never enable
credential/body logging or publish `.env` files, tokens, or raw native diagnostics.

For client-generated transport failures, the live helper also logs `causeHints` reconstructed
from a fixed vocabulary: connection setup/refusal/reset, permission denial, unreachable
network/host, DNS, TLS/certificate, timeout, and unexpected end-of-file hints. No original
message text, URLs, addresses, headers, credentials, or diagnostic payloads are copied.
Unrecognized transport messages produce `causeHints=unclassified`; service and authentication
details remain suppressed. These hints aid debugging but are not a stable error contract
or proof of a particular root cause.

Each scenario uses an isolated runtime and random item IDs/partition keys. Cleanup is
registered before writes with independent timeouts. Setup retries only newly deployed
RBAC's 403/5302 for at most 90 seconds; scenarios are not retried. Explicit session-token
reuse is covered without claiming broader session-management profile coverage.

[`test-resources.bicep`](test-resources.bicep) uses the deployment location (East US 2
in CI) to create a run-owned Session account,
database, two containers at 400 RU/s each, and a database-scoped data role assignment. The account
disables local (key) authentication, so only Entra ID tokens are accepted and no account
keys are emitted. Live jobs allow 45 minutes for provisioning and execution, with a 15-minute
cancellation cleanup grace period. Teardown runs on success, failure, and cancellation;
resource-group expiry tags provide a fallback for interrupted jobs. Container TTL
expires items after one day but does not stop account/throughput billing. Hard agent
termination cannot guarantee immediate teardown. Hosted provisioning, cancellation
cleanup, runtime, and cost must be verified before expanding this lane.

The result gate also has credential-free fixtures:

```powershell
./sdk/data/azcosmos/eng/scripts/Test-LiveResults.Tests.ps1
```

### Manual live change-feed split

Use a fresh disposable `/pk` Cosmos DB for NoSQL container with dedicated manual
throughput (400 RU/s to start) and one physical range. Use the same Entra setup as
the smoke tests, but set `AZCOSMOS_SPLIT_CONTAINER` to this dedicated container
and `AZCOSMOS_LIVE_SPLIT=1`. From `sdk/data/azcosmos`, run:

```shell
go test -v -race -count=1 -timeout=22m -run '^TestAzureLiveChangeFeedSplit$' ./e2e
```

Wait for `checkpoint saved` before scaling. From another terminal, with an
identity authorized to change throughput, raise **only the disposable container**
to 13,000 RU/s:

```shell
az cosmosdb sql container throughput update --subscription <test-subscription> \
  --resource-group <disposable-group> --account-name <test-account> \
  --database-name <test-database> --name <split-container> --throughput 13000
```

The test waits up to 15 minutes for Azure to report the child ranges. It checks
idle resume before any new writes, then checks that the original checkpoint
returns exactly the new items across the child ranges, without replay. A split
that never appears fails setup, not the recovery assertion. Delete the
disposable resources after the run, including on failure; scaling down does
not restore a single range. This test does not run in the smoke CI lane.
