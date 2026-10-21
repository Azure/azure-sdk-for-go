# Azure AI Content Understanding client prototype for Go

**Prototype. Not ready for review or production use.** This module has not been
released. Its API is intentionally small and may change without compatibility
guarantees.

This handwritten client explores Content Understanding analysis, analyzer
creation/listing/deletion/copy, and model deployment defaults. Wire contracts
are based on `Azure/azure-rest-api-specs` commit
`ca6afc6f24a2ef9ed210a164504718507b2e3659`, under
`specification/ai/ContentUnderstanding`. It is not a complete service SDK and
does not include a formatter, resource provisioning, or command-line policy.

## Authentication and analysis

Supply a resource endpoint such as `https://<resource>.services.ai.azure.com`
and an `azcore.TokenCredential`. Applications can use `azidentity`; the client
does not discover credentials or read application configuration.

```go
credential, err := azidentity.NewDefaultAzureCredential(nil)
if err != nil {
    return err
}
client, err := azcontentunderstanding.NewClient(endpoint, credential, nil)
if err != nil {
    return err
}
poller, err := client.BeginAnalyze(ctx, "prebuilt-layout",
    azcontentunderstanding.AnalyzeRequest{
        Inputs: []azcontentunderstanding.AnalysisInput{{URL: documentURL}},
    })
if err != nil {
    return err
}
response, err := poller.PollUntilDone(ctx,
    &runtime.PollUntilDoneOptions{Frequency: time.Second})
if err != nil {
    return err
}
// response.Result is the analysis JSON; response.Usage is the usage subtree.
// response.RawJSON is the complete service envelope, including unknown fields.
```

Resource keys are supported through `NewClientWithKeyCredential`.
`ClientOptions.APIVersion` defaults to `2025-11-01`. Select
`2026-06-01-preview` for `AnalyzeInline` and `AnalyzeBinaryInline`, which do
not poll. Analysis requests use `codePoint` span offsets: Go callers must
use rune-aware rather than byte slicing.

Binary methods require an `io.ReadSeekCloser` for retryable streaming; callers
can use `streaming.NopCloser` from `azcore/streaming` with an in-memory reader. Do not close the stream
until the submission returns. No input files are opened by the client.
Pollers use azcore's default interval unless the caller sets `Frequency`.
One second is suitable for interactive use; service retry instructions and
context deadlines take precedence.

`Analyzer.RawJSON` preserves the complete definition. Supply a JSON object to
`BeginCreateAnalyzer`; create/copy do not replace existing analyzers. The
application owns sorting, dependency preflight and other lifecycle policy.
Cross-resource copies require a source-side authorization grant followed by
a destination-side copy request.

`UpdateDefaults` accepts a map of string pointers: nil removes a mapping.
It does not populate aliases or merge application settings automatically.
Service errors retain `*azcore.ResponseError` status/code information. Inline
HTTP 200 responses with failed operation status return an error.

## Validation

Offline:

```sh
go test -race ./...
go vet ./...
go build ./...
```

Live tests require Azure CLI login, `CONTENTUNDERSTANDING_ENDPOINT`, and the
same resource's `CONTENTUNDERSTANDING_RESOURCE_ID` and
`CONTENTUNDERSTANDING_REGION` for authorization verification:

```sh
go test -tags live -run TestLive -count=1 -timeout 30m -v ./...
```

These tests incur service charges, submit a public sample invoice, and create,
copy and delete uniquely named analyzers, and verify that the service rejects
a same-resource authorization request for a test-owned analyzer. The resource must already have the
`gpt-5.2` completion mapping configured. They do not modify resource defaults.
Credentials and endpoint values are not recorded in this repository.

Successful cross-resource authorization/copy and defaults mutation have offline contract
coverage; the default live suite does not mutate shared defaults or assume
access to a second resource. Passing it alone does not establish those live
contracts. There is no promise of complete API or SDK feature parity.
