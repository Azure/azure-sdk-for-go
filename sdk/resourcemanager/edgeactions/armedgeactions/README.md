# Azure Edge Actions Module for Go

The `armedgeactions` module provides operations for managing Azure Edge Actions resources, including edge actions, versions, and execution filters. Edge Actions is part of the Azure Front Door platform and enables running custom logic at the network edge.

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/edgeactions/armedgeactions)

# Getting started

## Prerequisites

- An [Azure subscription](https://azure.microsoft.com/free/)
- [Supported](https://aka.ms/azsdk/go/supported-versions) version of Go (You could download and install the latest version of Go from [here](https://go.dev/doc/install). It will replace the existing Go on your machine. If you want to install multiple Go versions on the same machine, you could refer this [doc](https://go.dev/doc/manage-install).)

## Install the package

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Install the Azure Edge Actions module:

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions
```

## Authorization

When creating a client, you will need to provide a credential for authenticating with Azure Edge Actions.  The `azidentity` module provides facilities for various ways of authenticating with Azure including client/secret, certificate, managed identity, and more.

```go
cred, err := azidentity.NewDefaultAzureCredential(nil)
```

For more information on authentication, please see the documentation for `azidentity` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity).

## Client Factory

Azure Edge Actions module consists of one or more clients. We provide a client factory which could be used to create any client in this module.

```go
clientFactory, err := armedgeactions.NewClientFactory(<subscription ID>, cred, nil)
```

You can use `ClientOptions` in package `github.com/Azure/azure-sdk-for-go/sdk/azcore/arm` to set endpoint to connect with public and sovereign clouds as well as Azure Stack. For more information, please see the documentation for `azcore` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore).

```go
options := arm.ClientOptions {
    ClientOptions: azcore.ClientOptions {
        Cloud: cloud.AzureChina,
    },
}
clientFactory, err := armedgeactions.NewClientFactory(<subscription ID>, cred, &options)
```

## Clients

A client groups a set of related APIs, providing access to its functionality.  Create one or more clients to access the APIs you require using client factory.

```go
client := clientFactory.NewClient()
```

## Fakes

The fake package contains types used for constructing in-memory fake servers used in unit tests.
This allows writing tests to cover various success/error conditions without the need for connecting to a live service.

Please see https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/samples/fakes for details and examples on how to use fakes.

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Edge Actions` label.

## Contributor tests

Tests follow the repository's [management-plane testing and secret-scrubbing guidance](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/testing.md).
Assertions follow the repository's [current test standards](https://github.com/Azure/azure-sdk-for-go/blob/main/.github/instructions/go-tests.instructions.md):
`testify` assertions, `require` for dependent preconditions, `t.Helper()` on assertion
helpers, shared ARM credentials/proxy helpers, and playback-aware `recording.Sleep`.
Coverage is adapted from [Azure/azure-sdk-for-python#49467](https://github.com/Azure/azure-sdk-for-python/pull/49467)
and [Azure/azure-sdk-for-java#50810](https://github.com/Azure/azure-sdk-for-java/pull/50810).
Generated clients and examples are not modified by these tests.

| Tests | Coverage | Azure access |
| --- | --- | --- |
| `TestContract*` | Pinned response schemas, request bodies/API version, PATCH omission/empty/null semantics, pagination, async PUT 200/201, action/PATCH 202, status and final-result endpoints, typed results, failure/cancellation, and existing DELETE statuses | None; requests are routed to local mock servers |
| `TestHelper*` | Deterministic file/ZIP payloads, ownership checks, partial-create cleanup, bounded provisioning waits, and filter/non-default/default/parent deletion order | None |
| `TestSanitizer*` | Literal credential replacement and a local proxy record/playback round trip, including changed polling tokens and replayed sanitized polling URLs | None; the upstream is a local HTTP server |
| `TestRecorded*` | Tags and invalid updates, file/ZIP version deployment and retrieval, default-version swap, execution-filter CRUD, wire-response validation, and direct typed LRO-result assertions | Only in `record` or `live` mode |
| `TestRecordingResponseContracts` | Every saved Edge Actions response, including negative/error and cleanup responses; all schema failures per response, async endpoint association/completeness, and per-recording operation counts | None; reads local proxy assets without replaying requests |

`testdata/response_spec.json` is an offline response-only snapshot from the exact spec
commit in `tsp-location.yaml`, including successful/default error responses and referenced
common ARM schemas. Tests reject a
stale snapshot. Refresh it after a specification update with `go generate` from this
directory; only this refresh command needs network access to the public specification.
The Go validator supports the schema vocabulary in this snapshot (references, inheritance,
required properties, types, enums, arrays, dictionaries, and formats), rejects unsupported
validation keywords, and reports all schema paths/rules without response values.

Successful Edge Actions responses are checked before SDK deserialization. Async header
URLs must be absolute HTTPS URLs, preserve the pinned API version when present, and
distinguish the status monitor from the final result. Status monitors must return HTTP 200
with a nonempty string status and reach `Succeeded`. Advertised endpoints that the SDK
does not visit are checked explicitly through the same authenticated, recorded ARM pipeline,
including the monitor on an already-completed PUT. Final `Location` responses are checked
against the **initiating operation**, not the GET route that the URL happens to resemble.
The explicit exceptions do not modify the pinned snapshot:

- AsyncOperation response bodies are checked only for a nonempty string `status`; provider-specific
  fields are ignored. HTTP 200, URL checks, and successful completion remain required.
- Version `validationStatus` validation is temporarily omitted, including typed-result assertions,
  pending the RP fix to make the field optional, not nullable. This does not permit nulls in other fields.
- The known `swapDefault` final `Location` HTTP-200-with-body response is excluded until the RP fix;
  initial responses, header URLs, status polling, and subsequent default-version GETs remain checked.

Cleanup bypasses pipeline contract validation so a detected service defect cannot itself prevent
deletion. The recording-wide audit still validates those saved cleanup responses, including errors.
It inspects all recording entries even when a lifecycle test previously failed, reports missing
advertised monitor/final-result exchanges, and handles the proxy's sanitized request authority
without weakening validation of the original header URLs. Resource-group fixture operations and
their advertised polling endpoints are counted separately: they belong to a different ARM spec.
This audits **response** contracts for recorded operations; it does not claim request-schema
coverage or execution of operations absent from the recordings.

Run commands from this module's directory with a supported Go release (at least the
version specified in `go.mod`) and Git on `PATH`. Restore dependencies and build first:

```powershell
go mod download
go build ./...
```

Go's recording helpers install and start the
repository-pinned test proxy automatically. Playback is the default. Offline regressions
do not require service recordings or Azure credentials:

```powershell
$env:AZURE_RECORD_MODE = "playback"
go test ./... -run 'Test(Contract|Helper|Sanitizer)' -count=1
```

After capturing recordings, audit all stored responses independently of lifecycle playback:

```powershell
go test ./... -run '^TestRecordingResponseContracts$' -v -count=1 -timeout 5m
```

The audit locates assets with the repository-pinned proxy's `config locate` command and fails
when recordings are absent. Do not use `-failfast` for this audit: every file must be examined.
Running `go test ./...` also includes the audit.

`assets.json` pins the published Go service recordings. The repository test proxy restores
that tag for playback in another checkout or CI. Running `TestRecorded*` in playback
requires those recordings; missing recordings fail instead of silently skipping or making live requests.
The synthetic sanitizer recording is temporary and is not a substitute for service coverage.

To capture service recordings, install the Azure CLI if needed, select a test subscription
enabled for Microsoft.Cdn Edge Actions, and sign in with permission to create and delete
resource groups and Edge Actions. These tests create isolated resource groups in `LOCATION` (default
`eastus`) and Edge Actions in `global`; they refuse to overwrite existing resources.
Each LRO has a fresh 20-minute polling budget rather than sharing a cumulative scenario
deadline. HTTP attempts are bounded to two minutes. Cleanup runs on failure with its own
10-minute timeout, waits for non-terminal resources instead of retrying arbitrary DELETE
errors, deletes filters before versions, deletes the
current default version last, and then deletes the action and owned resource group.
Cleanup failures are reported and may require manual removal of test resources.

```powershell
az login
$env:AZURE_SUBSCRIPTION_ID = "<authorized-test-subscription-id>"
az account set --subscription $env:AZURE_SUBSCRIPTION_ID
$env:LOCATION = "eastus"
$env:AZURE_RECORD_MODE = "record"
go test ./... -run '^TestRecorded' -count=1 -timeout 150m

$env:AZURE_RECORD_MODE = "playback"
go test ./... -count=1 -timeout 10m
```

`live` mode calls Azure without saving recordings. The shared test credential supports
local CLI authentication and the repository's CI authentication; credentials are never
needed for playback. Set environment variables in the shell; these tests do not load
`.env` files automatically.

Sanitizers preserve resource names, SKU names, IDs, and LRO URL structure needed for
assertions and polling. They redact subscription/identity values, configured client
secrets, nested tenant/client/principal identity fields and `client_secret`, cookies,
OAuth tokens, creator identities, trace headers, and the values of
polling query parameters `t`, `c`, `s`, and `h` without erasing paths, API versions, or
unrelated parameters. Test resource names and the fixed JavaScript payload are deliberately
non-sensitive. Do not substitute confidential code or credentials for the synthetic payload.

After recording, inspect every recording for secrets and verify playback before publishing.
Use the repository-installed proxy's `config locate -a assets.json` command to find the
asset working directory. Publication is a separate, explicitly authorized step:

```powershell
..\..\..\..\.proxy\Azure.Sdk.Tools.TestProxy.exe push -a assets.json
```

Commit the resulting non-empty asset tag with the tests only after successful record/playback
validation. Never commit real credentials or use Python's asset tag for Go tests.

# Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.