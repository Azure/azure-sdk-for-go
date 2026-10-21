# Content Understanding Go emitter iteration

Local prototype generated with `@azure-tools/typespec-go` **0.17.6** using the
repository's pinned emitter manifest/lockfile and `tsp-client` **0.30.0**.
This is not a published SDK or a drop-in replacement for the handwritten
prototype. The original prototype branch is unchanged.

## Source and regeneration

The service source is `Azure/azure-rest-api-specs` commit
`ca6afc6f24a2ef9ed210a164504718507b2e3659`, directory
`specification/ai/ContentUnderstanding`, recorded in `tsp-location.yaml`.
The local `generation/spec-go.patch` contains the complete generation input
changes:

- Add Go emitter output/module options, raw-JSON bytes and an explicit working
  azcore version.
- Align the base `ContentField.type` Go client name with the existing derived
  `fieldType` names. Without this Go-scoped customization, the emitter fails
  with `unsupported type kind enum for discriminator property type`.

To reproduce, start with a clean disposable checkout of the spec commit, apply
the patch at its repository root, and run from this SDK module:

```sh
git -C "$SPEC_REPO" apply --unidiff-zero "$SDK_MODULE/generation/spec-go.patch"
tsp-client update \
  --local-spec-repo "$SPEC_REPO/specification/ai/ContentUnderstanding" \
  --no-prompt
```

`SPEC_REPO` and `SDK_MODULE` are absolute paths to those local directories.
For tsp-client 0.30.0, the local-spec argument must name the CU project directory,
not the specification repository root. Use Node 24 and Go on PATH. The SDK repo
must retain its pinned `eng/emitter-package.json` and lockfile.
No generated Go file is manually patched.

## Generated surface and compatibility

Ten generated runtime files provide clients, typed/polymorphic models,
serialization, responses, options and constants. All twelve service operations
required by the earlier prototype are present; generation also includes
analyzer update and result deletion/download operations.

Important differences from the handwritten prototype:

- The emitted client defaults to `2026-06-01-preview`. Select GA through
  `ClientOptions{ClientOptions: azcore.ClientOptions{APIVersion: "2025-11-01"}}`.
- String encoding is an explicit argument. Tests use `codePoint`; Go consumers
  must use rune-aware slicing for these offsets.
- The generated constructor currently accepts an Entra `TokenCredential`; it
  does not generate the earlier key-credential constructor. The emitter reports
  `UnsupportedAuthenticationScheme: unsupported authentication scheme apiKey
  will be omitted`.
- Requests/responses use generated typed structs rather than raw definitions.
  Unknown properties are not retained by typed reserialization.
- The completed LRO result projects `AnalysisResult` and does not expose its
  top-level usage. Use `policy.WithCaptureResponse` for the begin and polling
  context, then `runtime.Payload` on the captured response to retain the original
  envelope, unknown fields and usage. Offline and live tests verify this path.
- Inline response models include status/result/usage, but the generated method
  does not convert HTTP-200 failed operation statuses into Go errors. Consumers
  must inspect status; this is tested explicitly, not treated as parity.
- The emitter currently generates `application/json` for defaults PATCH even
  though the pinned REST contract specifies `application/merge-patch+json`.
  The offline test uses the supported `policy.WithHTTPHeader` override to
  restore the required content type. This still needs a durable generator or
  TypeSpec-level fix before treating the generated API as a direct replacement.
- Omitted JSON model-deployment overrides currently serialize as
  `"modelDeployments": null`; the live analysis paths accept this shape.
- Despite an initial `azcore-version: "1.20.0"` attempt, the emitter imports
  `azcore/runtime/datetime`, which that version does not provide. Dependency
  resolution upgraded azcore. This iteration explicitly pins the verified
  **azcore v1.23.3**; it does not establish compatibility with v1.20.0.

Applications own credentials, input/output planning, command policy and result
formatting. This iteration adds no replacement business-logic module or
handwritten wrapper client.

## Layout and validation

Generated runtime files are in the module root. Public-boundary tests remain
separate under `tests/unit` and `tests/live`; generation inputs are under
`generation`.

```sh
go build .
go test -race -coverpkg=. ./tests/unit/...
go vet -tags live ./...
```

Offline coverage is intentionally focused on consumed operations, not all
generated document/media models. It is not comparable to coverage of the much
smaller handwritten client.

Live tests use Azure CLI authentication and require
`CONTENTUNDERSTANDING_ENDPOINT`:

```sh
go test -tags live -run TestLiveGenerated -count=1 -timeout 30m -v ./tests/live/...
```

They submit a public sample invoice and create/copy/delete uniquely named test
analyzers, requiring existing `gpt-5.2` defaults. They exercise GA/preview binary,
URL and base64 analysis, preview inline calls, raw-envelope/usage capture,
sample schema suggestion, field extraction, analyzer paging/conflict/copy and
deletion. Cleanup failures fail the tests.

Defaults mutations and successful cross-resource authorization/copy are not
live-verified. Shared defaults are not modified. The generated defaults PATCH
header issue and inline failure semantics remain explicit integration gaps.

The local validation run passed build, vet, race tests, module verification,
and the live scenarios above. Focused offline runtime coverage was 20.4%.
Repeated generation produced byte-identical output for all ten generated Go
files. Compilation also reports the existing C# `Response` suffix naming
warning from the shared spec; it does not prevent Go generation.
