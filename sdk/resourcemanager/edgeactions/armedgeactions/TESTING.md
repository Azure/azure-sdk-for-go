# Testing Edge Actions

For development setup, authentication, recording modes, sanitization, and publishing assets,
follow the repository's [testing guide](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/testing.md).
Run the commands below from this module's directory with Go and Git on `PATH`.
Configure the HTTPS proxy certificate in the same PowerShell session before recording
or playback; the tests start the proxy and register sanitizers automatically:

```powershell
$env:PROXY_CERT = Join-Path (git rev-parse --show-toplevel) 'eng\common\testproxy\dotnet-devcert.crt'
```

## Local checks

Playback requires no Azure credentials and creates no resources. The shared test proxy
restores the recordings pinned by `assets.json`; missing recordings fail rather than
falling back to live requests.

```powershell
$env:AZURE_RECORD_MODE = "playback"
go test ./... -v -p 1 -count=1 -timeout 10m
```

For mock/helper/sanitizer regressions without service recordings:

```powershell
go test ./... -run 'Test(Contract|Helper|Sanitizer)' -count=1
```

To audit every stored Edge Actions response independently of lifecycle playback:

```powershell
go test ./... -run '^TestRecordingResponseContracts$' -v -count=1 -timeout 5m
```

The full suite includes this audit in record/playback modes; live mode skips it
because live runs do not require a local proxy or stored recordings.
Avoid `-failfast` when auditing so every recording
is examined, including error and cleanup responses.

## Recording requirements

Use an authorized subscription enabled for Microsoft.Cdn Edge Actions, with permission
to create and delete resource groups and Edge Actions. Set `AZURE_SUBSCRIPTION_ID` and
optionally `LOCATION` (default `eastus`) in the shell or a module-local `.env`.
If you do not already have `.env`, copy [sample.env](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/edgeactions/armedgeactions/sample.env):

```powershell
Copy-Item sample.env .env
```

Replace the subscription placeholder in `.env` with your authorized test subscription.
Keep an existing `.env` rather than overwriting your local configuration.

Record and live modes load `.env` before proxy, sanitizer, and client setup.
Existing shell/CI environment variables take precedence. A missing `.env` is
allowed; unreadable or malformed files fail setup without printing their contents.
Playback does not load `.env`. The repository ignores `.env`; do not commit local
configuration or credentials. An existing Azure CLI login can supply authentication.

For local Azure CLI authentication, ensure `az` is on `PATH`, sign in if needed, and
select the same subscription configured in `.env` or the shell:

```powershell
az login
az account set --subscription '<your-test-subscription-id>'
if ($LASTEXITCODE -ne 0) { throw 'Subscription selection failed' }
```

Set `AZURE_RECORD_MODE` and, when needed, `PROXY_CERT` in the shell before running
`go test`, not in `.env`: the recording package reads these during initialization,
before test setup. The package-wide `-timeout` remains a command-line argument.

```powershell
$env:AZURE_RECORD_MODE = "record"
go test . -run '^TestRecordedEdgeActions$' -v -count=1 -timeout 3h
if ($LASTEXITCODE -ne 0) { throw 'Recording failed; inspect operation and cleanup errors before proceeding' }
```

The seven scenarios are eligible to run concurrently within `TestRecordedEdgeActions`.
The module defaults to seven parallel tests even on machines with fewer CPUs. Use
Go's `-parallel N` flag to override this limit; `-parallel 1` runs sequentially.
Lower concurrency needs a longer package timeout, for example `-parallel 1 -timeout 12h`.
Individual scenario deadlines do not change.

Use `go test . -v` for streamed console output while recording. With multiple packages,
`go test ./... -v` can buffer output until a package finishes; use `-p 1` or `-json`
when streaming output from the entire module. The `-p` flag controls packages,
not the parallel scenarios within this package. Long-running service polls can still
be quiet between test start and completion messages.

To record only the six file/ZIP scenarios while retaining a passing tags capture
at the current recording path:

```powershell
$env:AZURE_RECORD_MODE = "record"
go test . -run '^TestRecordedEdgeActions$/^(VersionDeployment|VersionSwap|ExecutionFilters)$' -v -count=1 -timeout 3h
if ($LASTEXITCODE -ne 0) { throw 'Recording failed; inspect operation and cleanup errors before proceeding' }
```

Scenarios cover tags, file/ZIP version deployment and retrieval, default-version swaps,
and execution-filter CRUD. They create isolated resource groups and Edge Actions in
`global`, and refuse to overwrite existing resources. Resource groups use
`go-sdk-test-<digits>`; the shared ARM sanitizer replaces their names with `Sanitized`
in recordings, and playback uses that placeholder. Edge Action names retain the `eago` prefix.
Each scenario owns its resources and recording, so it can run or be retried independently.
Scenario deadlines cover setup and all test operations:

| Test | Formats | Scenario deadline |
| --- | --- | --- |
| `TestRecordedEdgeActions/Tags` | No code deployment | 20 minutes |
| `TestRecordedEdgeActions/VersionDeployment` | File and ZIP | 50 minutes each |
| `TestRecordedEdgeActions/VersionSwap` | File and ZIP; deploy v1/v2, then swap | 130 minutes each |
| `TestRecordedEdgeActions/ExecutionFilters` | File and ZIP; deploy v1/v2, then filter CRUD | 90 minutes each |

The deployment scenarios check code retrieval and deletion. The swap scenarios
include version listing, tag/property updates, rejected updates, and default-state
checks. Filter scenarios do not perform a version swap; v1 remains the default,
the filter changes from v1 to v2, and teardown deletes v2 before v1.

For example, re-record only the ZIP swap scenario:

```powershell
go test . -run '^TestRecordedEdgeActions$/^VersionSwap$/^zip$' -v -count=1 -timeout 160m
```

Recordings now live under `TestRecordedEdgeActions`: `Tags.json`, plus `file.json`
and `zip.json` in each of `VersionDeployment`, `VersionSwap`, and `ExecutionFilters`.
Captures from the earlier separate test names must be moved to the corresponding
new paths or re-recorded, then verified with playback; do not edit captured responses.
The previous combined version/filter lifecycle recordings cannot replace the split
scenario recordings.

Each LRO can poll for up to 40 minutes, subject to the remaining scenario deadline;
each HTTP attempt is bounded to two minutes. With seven-way parallelism, the longest
scenario allows 130 minutes plus up to 20 minutes of independent cleanup; the
three-hour package timeout leaves additional time for startup and overhead.
Cleanup has independent 10-minute budgets, removes filters before versions and the default version last, then removes the
action and resource group. Inspect reported cleanup failures for resources needing manual removal.

## Validate replacement recordings

Do not run overlapping recording jobs against the same asset directory. Once all six
split scenarios pass and no old combined run is active, back up the recordings and move
the retired `TestRecordedEdgeActionVersionsAndFilters` folder and any superseded
top-level scenario captures outside the active recording directory. Keep the passing
tags capture and all six split captures under `TestRecordedEdgeActions`. The response audit
examines every stored recording, including obsolete captures unless archived.
Do not remove failing captures for active scenarios to make the audit pass.

Run full playback, which includes the helper/sanitizer regressions and exhaustive
recorded-response audit:

```powershell
$env:AZURE_RECORD_MODE = "playback"
go test ./... -v -p 1 -count=1 -timeout 10m
if ($LASTEXITCODE -ne 0) { throw 'Playback or recording audit failed; do not publish' }
```

Keep the synthetic JavaScript payload and generated test names non-sensitive. Inspect
sanitized captures for secrets and reported cleanup failures before publishing.
The test commands do not publish assets or commit changes. Follow the shared guide to
locate/publish assets, and commit the resulting `assets.json` tag with the test changes.

## Response-contract maintenance

`testdata/response_spec.json` snapshots response schemas from the commit in
`tsp-location.yaml`. After updating the spec pin, run `go generate` to refresh the snapshot;
this command downloads the pinned specification. Tests reject stale snapshots and
unsupported validation keywords.

The audit checks response contracts, not complete request schemas or operations absent
from recordings. Final `Location` responses must have a successful terminal HTTP status
and are validated against the initiating operation. Error bodies are still checked,
but cannot satisfy successful completion.
Resource-group fixture operations use a separate ARM API and are counted outside the
Edge Actions audit. Cleanup bypasses pipeline validation to preserve resource removal;
its recorded responses are still audited.

For DELETE operations that advertise both `Azure-AsyncOperation` and `Location`,
a successful status monitor is sufficient evidence of completion: the SDK's async
DELETE poller does not fetch a final `Location`. Any recorded `Location` response
is still checked for terminal HTTP success and response-schema validity.
Location-only DELETE recordings must include a successful terminal response.

Current exceptions, implemented in the response-validation helpers, leave the snapshot unchanged:

- AsyncOperation bodies validate only a nonempty string `status`; HTTP 200, URL checks,
  and successful completion remain required.
- Version `validationStatus` schema and typed-result assertions are temporarily omitted
  pending the service fix to make the field optional, not nullable.
- The known `swapDefault` final `Location` HTTP-200-with-body response is excluded pending
  the service fix. Initial responses, polling, and subsequent resource GETs remain checked.
  A missing final response is allowed only when an Azure-AsyncOperation monitor proves
  completion; Location-only polling must reach a successful terminal HTTP response.
