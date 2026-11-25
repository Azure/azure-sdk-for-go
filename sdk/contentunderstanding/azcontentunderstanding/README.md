# Azure AI Content Understanding client library for Go

**Prototype: not released or ready for production use.** This library combines
TypeSpec-generated service operations/models with handwritten Go conveniences.
The functionality reference is the .NET CU SDK at commit
`fe8ee421b95316cb301c3342112f496fb9dab90c`, not its C# API signatures or
platform-specific dependency-injection/compatibility APIs.

## Supported workflows

- GA `2025-11-01` and preview `2026-06-01-preview`; preview is the default,
  matching the .NET preview library.
- JSON URL/base64 and streaming binary analysis, inline analysis, processing
  location, model overrides, ranges and truncation options.
- Analyzer create/get/list/update/delete, same-/cross-resource copy and copy
  authorization; defaults get/merge-patch; result file download/deletion.
- Standard azcore polling, cancellation and cross-process resume tokens.
- Raw operation JSON, usage and operation identifiers; typed document/media
  results including preview metadata, diagnostics, signatures, semantic chunks,
  agentic workflows and in-page segments.
- Page/time range constructors, parsed document/media grounding, media durations
  and LLM-input formatting.

Generated models use normal Go structs, maps, slices and pointer-valued optional
properties instead of C# model factories, indexers and overloaded sync/async APIs.
This is not source-compatible with the earlier handwritten Go prototype.

## Client construction

Applications supply credentials and their resource endpoint:

```go
credential, err := azidentity.NewDefaultAzureCredential(nil)
if err != nil {
    return err
}
client, err := azcontentunderstanding.NewClient(endpoint, credential,
    &azcontentunderstanding.ClientOptions{
        ClientOptions: azcore.ClientOptions{APIVersion: "2025-11-01"},
    })
if err != nil {
    return err
}
```

`NewClientWithKeyCredential` accepts an `azcore.KeyCredential` for resource-key
authentication and rotation. Transport, retry, tracing and logging settings use
the embedded azcore options; custom cloud audiences use `ServiceName`.
The client does not discover credentials or modify resource defaults.

The constructors install service-specific policy that restores merge-patch
content types, rejects inline failed statuses even on HTTP 200, and prevents
continuation requests from sending credentials to another resource endpoint.
Inline status failures are checked outside the retry loop to avoid resubmission.

## Analysis, usage and raw output

```go
var raw *http.Response
capture := policy.WithCaptureResponse(ctx, &raw)
poller, err := client.BeginAnalyze(capture, "prebuilt-layout",
    []azcontentunderstanding.AnalysisInput{{URL: to.Ptr(documentURL)}},
    "utf16", nil)
if err != nil {
    return err
}
operationID, err := azcontentunderstanding.GetOperationID(raw)
if err != nil {
    return err
}
_, err = poller.PollUntilDone(capture,
    &runtime.PollUntilDoneOptions{Frequency: 3 * time.Second})
if err != nil {
    return err
}
result, err := azcontentunderstanding.GetAnalysisResult(raw)
if err != nil {
    return err
}
usage, err := azcontentunderstanding.GetUsageDetails(raw)
if err != nil {
    return err
}
```

The operation ID is read before polling replaces the captured response.
The final response can also be read through `GetRawJSON`, retaining all unknown
properties and numeric precision. `GetAnalysisResult` additionally accepts the
observed `KeyFrameTimesMs` service alias without changing the original raw JSON.
Use it for media results needing that compatibility, rather than reserializing
the generated poller's projected model.

Missing/null usage returns nil. Malformed present usage returns an error rather
than silently looking absent. This is an intentional difference from the .NET
best-effort usage accessor.

Binary methods accept `io.ReadSeekCloser`; `streaming.NopCloser` adapts an
in-memory reader for retryable streaming. An empty content type defaults to
`application/octet-stream`. Inline methods return status/result/
usage without polling, and require the preview service API.

String encoding is explicit: use `utf16` to match .NET span offsets, `codePoint`
for rune offsets, or `utf8` for byte offsets. `SliceContentSpan` validates and
slices all three encodings. Native azcore pollers retain their own default wait
frequency; pass three seconds explicitly to match .NET analyze convenience
waiting. Server retry instructions and context deadlines take precedence.
Resume through the generated `ResumeToken` option; tokens are Go/runtime
specific, not interchangeable with .NET tokens.

## Typed and raw request definitions

Generated operations accept typed analyzer definitions and inputs. For JSON
schemas/patches that must retain properties unknown to this library, use:

- `BeginCreateAnalyzerJSON`
- `UpdateAnalyzerJSON`
- `BeginAnalyzeJSON`
- `AnalyzeInlineJSON`

These reuse generated request builders and response/poller handling while
preserving the submitted JSON object, including null/false/zero values.
They do not perform application schema validation or merge model mappings.
For typed defaults updates, `ClientUpdateDefaultsOptions.ModelDeployments`
contains `RecordMergePatchUpdate.AdditionalProperties`; nil map values remove
individual mappings.

Result downloads return a caller-owned `Body` stream. Close it after reading.
Delete only the results/resources the application owns.

## Result conveniences

`ContentRange` supports raw expressions plus validated `Page`, `Pages`,
`PagesFrom`, `TimeRange`, `TimeRangeFrom` and `CombineRanges` constructors.
Its `Ptr()` method supplies generated range options.

`ParseContentSources` handles page-only/polygon `D(...)` and timestamp/rectangle
`AV(...)` regions. `ContentSource` exposes document properties through
`PageNumber`/`Polygon`, media time through `Time`, and optional `BoundingBox`.
`GetFieldSources` parses a generated field's grounding; `FormatContentSources`
preserves raw region expressions.

Media content, segments, transcript phrases and words expose duration methods.
Missing timestamps and values outside Go's `time.Duration` range return errors
instead of overflowing.

```go
text, err := azcontentunderstanding.ToLLMInput(result,
    map[string]any{"source": "invoice.pdf"}, nil)
```

`LLMInputOptions` controls fields and Markdown with pointer-valued booleans,
defaulting to true. Output includes ordered YAML front matter, service metadata
separate from `customMetadata`, nested typed fields/JSON, responsible-AI warnings,
original page markers, classification/routing and media time ranges.

Go map keys are sorted deterministically rather than preserving .NET dictionary
insertion order. Invalid spans/unsupported metadata return explicit errors.
JSON numbers and nested empty/null metadata remain lossless; YAML safety and
numeric-overflow fixes are not intended to reproduce unsafe reference output.
These output differences are deliberate and should be considered when comparing
literal strings instead of parsed content.

## Source layout and regeneration

Runtime/generated files are at the module root; tests are separated into
`tests/unit` and opt-in `tests/live`. Generation inputs are under `generation`.
No CLI/profile/provisioning or other application business logic is included.

The generator is `@azure-tools/typespec-go` **0.17.6**, pinned by the repository
emitter manifest/lockfile. `tsp-location.yaml` pins the CU TypeSpec to
`Azure/azure-rest-api-specs` commit
`ca6afc6f24a2ef9ed210a164504718507b2e3659`.

From a clean disposable spec checkout at that commit:

```sh
git -C "$SPEC_REPO" apply --unidiff-zero "$SDK_MODULE/generation/spec-go.patch"
# Run from SDK_MODULE, with Node 24, Go and tsp-client 0.30.0 on PATH:
tsp-client update \
  --local-spec-repo "$SPEC_REPO/specification/ai/ContentUnderstanding" \
  --no-prompt
```

The patch aligns base/derived Go discriminator names, enables raw JSON bytes,
pins the working azcore version and uses `omit-constructors` so handwritten
constructors install the required policies. Generated Go files are not manually
modified. The emitter's API-key omission warning is addressed by the handwritten
key constructor, not by pretending it generated key authentication.

This implementation uses **azcore v1.23.3**. The emitted serializers import
`azcore/runtime/datetime`; compatibility with v1.20.0 is not established.

## Validation

```sh
go build ./...
go test -race -coverpkg=. ./tests/unit/...
go vet -tags live ./...
go mod verify
```

Formatter validation also ran 84 expanded .NET reference assertions and 106 Go
differential checks, accounting for map ordering and the documented numeric
safety difference. The optional `TestToLLMInputDotNetParity` test accepts an
exported reference fixture through
`AZURE_CONTENT_UNDERSTANDING_DOTNET_LLM_CASES`; ordinary Go tests remain
self-contained.

Live tests require Azure CLI login, `CONTENTUNDERSTANDING_ENDPOINT`, and existing
`gpt-5.2`/embedding/prebuilt model mappings:

```sh
go test -tags live -run TestLive -count=1 -timeout 30m -v ./tests/live/...
```

The suite incurs charges, submits public reference invoices/documents/media and
creates uniquely named analyzers. It exercises both service versions, binary/
URL/base64/inline analysis, schema suggestion, field grounding, raw/usage/
operation resume, PDF/DOCX metadata, signatures, semantic chunks, agentic
workflows, in-page classification, media transcripts/keyframes, file downloads,
analyzer updates/copy/deletion and result deletion. Test-owned resources are
cleaned up; shared defaults are not modified.

Resource-key authentication can be verified separately with a caller-supplied
`CONTENTUNDERSTANDING_KEY`, without printing or recording it:

```sh
go test -tags live -run '^TestResourceKeyAuthentication$' -count=1 ./tests/live/...
```

This read-only key-authentication test also passed. Successful cross-resource
authorization/copy and defaults mutation have offline contract evidence but
remain unverified live; these are validation gaps, not passed scenarios.
