# Automatic attested managed identity mTLS PoP

## Status

- Last updated: 2026-11-18.
- Working branch: `mtls-pop-1`, based on the existing mTLS PoP implementation.
- Current phase: phases 2 through 5 implemented; phase 6 qualification passed for
  the available fixture. Release ordering is documented; version changes and
  publication are deferred by agreement.
- Automatic managed identity eligibility selection and blank-import activation
  are implemented, together with automatic Key Vault signaling, narrowly classified
  tenant fallback, and bounded binding-error retries. Native MSAL behavior is unchanged.
- None of the mTLS PoP work to date has been released in any form. Public
  surface added by this work can be renamed or removed without a deprecation
  period or compatibility shim. Previously released SDK APIs remain unaffected.

This is a living record of decisions, implementation order, validation, and
remaining work. Update the relevant status and decision entries as work lands.
Distinguish tested baseline behavior from proposed behavior.

## Goal

The tested C# implementation is the north star for behavior and customer
experience. Match it unless an existing Go SDK design or language constraint
requires a deviation; convenience or preserving our unreleased implementation
is not sufficient justification. Use a documented Go blank import instead of
an optional NuGet reference:

```go
import _ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation"
```

Applications construct ordinary `ManagedIdentityCredential` or
`DefaultAzureCredential` instances and PoP-aware service clients without setting
attestation or PoP enable flags. Eligible requests automatically use
attested mTLS PoP. Ordinary bearer requests remain ordinary bearer requests.

The import makes the integration available; it does not require every credential,
request, or service to use attestation. Native DLL and runtime deployment
prerequisites remain necessary.

## Confirmed decisions

| Decision | Status |
| --- | --- |
| Use a documented blank import rather than a build tag as the primary opt-in mechanism. | Agreed |
| The optional module is `github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation`. | Agreed |
| C# is the behavioral reference; deviate only for documented existing SDK design or language constraints. | Agreed |
| All surface introduced by this unreleased work may be renamed or removed; no legacy enable-flag compatibility layer is required. | Confirmed |
| Target the automatic C# experience rather than requiring customers to set multiple enable flags. | Agreed |
| Do not introduce a new customer-facing strict-mode option for this work. | Agreed |
| Use the zero-value-friendly MIC opt-out `DisableMTLSPoP bool`; do not add another public DAC option. Remove superseded unreleased enable flags. | Implemented; names retained in phase 5 |
| Use a minimal internal registration hook that records availability only; keep native work lazy and in MSAL. Do not add a public registration/provider framework. | Agreed |
| Preserve ordinary bearer requests and caller transport configuration. Return an explicit error if a bound request cannot safely apply its certificate; never silently downgrade or discard caller configuration. | Agreed |
| Initially enable automatic behavior for Key Vault clients using the shared authentication policy, with secret reads as the first live validation target. Other services remain unchanged. | Agreed |
| Choose reasonable initial public names during implementation; naming is not a blocker. Perform a final naming cleanup before release. | Agreed |
| Keep hardware-backed live tests on demand; use fixture-free tests for routine regression protection. | Established testing approach |
| Use minimal SDK-internal authentication-error and response-classification bridges, not a Key Vault runtime dependency on azidentity or combined-message parsing. Standard Go error methods are observable additions; no customer switches or new public error types. | Agreed and implemented in phase 4 |
| Phase 6 documents release ordering only. Leave dependency manifests and versions unchanged; human release owners perform version alignment and publication later. | Agreed |

Only an actual package import activates a Go package. Adding a requirement to
`go.mod` alone is not activation. A transitive import can activate the integration
as well, which must be documented.

## C#-aligned behavior contract

This table describes the agreed automatic SDK contract, checked against the
[pinned C# sources](#reference-implementation) and implemented through phase 5.
The [phase-6 qualification](#6-run-parity-validation-and-coordinate-releases)
records the hardware-backed results and the cases covered only offline.
The constraint-driven adaptations below are agreed; phase 5 retains the selected API names.

| Condition | Intended outcome |
| --- | --- |
| Integration not imported | Do not select the attested PoP path; use ordinary authentication. |
| Request does not ask for PoP, or managed identity PoP is explicitly disabled | Ordinary bearer acquisition; do not perform attestation. |
| Imported, PoP-aware client, managed identity PoP permitted, host reports mTLS PoP support and at least KeyGuard binding strength | Automatic attested mTLS PoP. |
| Capability discovery succeeds but the host does not meet the PoP/KeyGuard requirements | Ordinary managed identity bearer acquisition. |
| Non-managed-identity credential selected | Preserve that credential's supported authentication behavior; the service's PoP signal is not a universal bound-token requirement. |
| Capability discovery fails or is canceled | Preserve error, cancellation, and normal credential-chain classification; do not equate discovery failure with known lack of support. |
| Token acquisition fails with the specifically classified tenant denial `AADSTS3921996` | Acquire a fresh bearer token with PoP disabled and remember the decision for that tenant in the Key Vault policy's credential wrapper for the client lifetime. |
| Attestation selected but fails, including a missing or unusable native dependency | Surface the failure; do not silently retry without attestation. |
| Bound Key Vault request receives HTTP 401 with `[MtlsCnfClaimRequestDataValidationFailed]` | Classify as retriable under the configured outer retry policy; do not downgrade to bearer or add an unbounded authentication resend loop. |
| RBAC denial, other binding failure, or unrelated network/service failure | Preserve normal retry/error handling; do not treat the failure as permission to downgrade to bearer. |
| Custom transport cannot apply the binding certificate | Not a C# bearer-fallback trigger. Define the Go transport adaptation explicitly; do not silently change authentication strength. |

Attestation is performed when needed for credential issuance, not necessarily
on every token request or service operation. Existing valid caches remain usable
subject to their binding and attestation guarantees.

Fallback is not a generic catch-and-retry policy. In particular, Key Vault
authorization failures are not equivalent to an Entra tenant-eligibility denial.
The tenant classifier must preserve the C# distinctions for cancellation and
mixed credential-chain failures rather than accepting any error text containing
the denial code. Its memory is tenant-keyed and local to the wrapper, not
process-wide. Fallback requires a new token acquisition, not relabeling a bound
token's authorization scheme.

### C# reference details

- The service policy requests PoP; a direct token request defaults to not requesting
  PoP. Importing the integration must not change that request-level default.
- MIC permits PoP by default. In the pinned C# version,
  `ManagedIdentityCredentialOptions.EnableMtlsProofOfPossession` is public and
  defaults to `true`. DAC also defaults its managed identity path to enabled,
  but its corresponding option is internal; do not assume public DAC opt-out
  surface is needed for parity.
- The Azure.Core MSAL adapter selects attested PoP only when the request,
  credential setting, host capabilities, and optional managed extension permit it.
  An absent or unresolvable optional managed extension prevents selection. This
  differs from failure of native attestation after that path has been selected.
- The configured MSAL minimum binding strength is KeyGuard, not simply any
  certificate-bound or IMDSv2-capable host.
- C# updates the resource transport for a returned binding certificate. Its
  transport-update path does not request a bearer token when a custom transport
  cannot accept that update.
- The Key Vault policy adds `x-ms-tokenboundauth` for `mtls_pop` authorization
  and removes it for ordinary authorization.

### Constraint-driven adaptations

| Constraint | Adaptation | Status |
| --- | --- | --- |
| A Go module requirement does not initialize its packages; Go has no equivalent of the C# optional assembly lookup. | Documented blank import with process-wide registration. | Agreed |
| Go options must preserve equivalent behavior for nil and zero-valued options; a plain enable bool cannot distinguish unset from explicitly false. | Use a MIC `DisableMTLSPoP bool` opt-out: the zero value permits automatic PoP when eligible. | Implemented |
| Go and C# expose different transport extension mechanisms. | Preserve caller transport configuration for bearer and unauthenticated requests. Return an explicit error when a bound request cannot safely apply its certificate; do not silently downgrade. | Implemented in phases 2 and 5 |

Record any additional deviation here with its concrete constraint, C# reference
behavior, Go adaptation, behavioral impact, and approval status before adopting it.

## Implementation plan

Each stage includes focused offline tests. Keep intermediate changes testable
and avoid enabling new service defaults before the underlying paths are ready.

### 1. Finalize semantics and API shape

**Status: complete; implementation preparation was handled in phases 2 through 5.**

- MIC uses a zero-value-friendly opt-out, initially `DisableMTLSPoP bool`.
  Preserve DAC's enabled default and internal propagation without another public
  DAC option.
- Remove superseded unreleased enable fields, updating their tests, examples,
  and callers together. No deprecation layer is required.
- Use a minimal internal registration hook with initialization limited to
  recording availability. Resolve its exact internal signature during
  implementation; no public provider framework is needed.
- Start with Key Vault clients using the shared authentication policy and
  validate secret reads first.
- Keep ordinary bearer requests working with caller transports. If a bound
  request cannot safely apply its certificate, fail explicitly without
  downgrading or discarding caller configuration.
- Activation is process-wide and may occur through transitive imports. Document
  that behavior and test imported/unimported applications in separate processes.
- Initial public names are best guesses, not implementation blockers. Schedule
  a final naming cleanup rather than delaying behavioral work.

Carry the following preparation into the relevant implementation stages:

- Translate matrix rows into focused selection, fallback, and retry test cases.
- Inventory obsolete flags and dependent tests/examples before removing them.

**Exit criterion:** agreed API and selection/fallback rules, with every
intentional C# deviation tied to a documented constraint.

**Outcome:** satisfied for beginning phase 2. No automatic-path code has been
implemented as part of this design sign-off.

### 2. Add automatic token and transport handling in azcore

**Status: complete; focused offline regression suites, build, and vet passed.**

Update the [authentication policy](../azcore/runtime/policy_bearer_token.go),
[mTLS transport](../azcore/runtime/policy_mtls.go), and
[pipeline integration](../azcore/runtime/pipeline.go):

- Allow automatic mode to handle either a bearer or certificate-bound token.
- Use binding headers and a certificate-scoped transport for bound tokens.
- Preserve ordinary authorization and normal transport behavior for bearer tokens.
- Preserve certificate-free cold challenges, HTTPS enforcement, redirect
  safeguards, middleware, cancellation, and retry behavior.
- Preserve token-cache separation and certificate-scoped connection pools,
  including certificate replacement and expiration handling.
- Do not enable automatic PoP indiscriminately for all Azure services.
- Define unsupported custom-transport behavior without relaxing TLS validation
  or inventing an additional bearer-fallback condition.

**Exit criterion:** offline authentication and real local-TLS tests pass for
bearer and bound paths, with no default service behavior change yet.

**Implemented:**

- Every runtime pipeline can lazily select a certificate-scoped transport when
  a credential returns a binding certificate. Installing this support does not
  enable PoP token requests globally or create extra HTTP clients for bearer reads.
- Authorization handlers can signal PoP support per token request without the
  existing client enable flag. The returned token determines authorization and
  transport: a bearer response is permitted, and a bound response is never
  relabeled as bearer.
- Bearer authorization clears previous binding state and the token-bound header.
  It uses the original configured transport, preserving caller TLS settings.
- Request modes retain separate token caches. Bound connections retain the
  existing certificate-specific pools and certificate-expiration handling.
- Credential failures and incompatible bound transports surface errors without
  triggering a bearer acquisition. Existing HTTPS, trust, and redirect protections
  remain in place.
- At the phase-2 boundary, the unreleased client enable flag remained as a
  pipeline-wide preference with a certificate-free cold-challenge path. Phase 5
  removes that preference and leaves unauthenticated requests on the caller's transport.

[Offline and local-TLS regressions](../azcore/runtime/policy_mtls_test.go) cover
service-level signaling, ordinary bearer defaults, bearer and bound responses,
TLS configuration and middleware preservation, cancellation, cold challenges,
opaque custom transports, and a bound-to-bearer challenge followed by a retry and
cached read. Existing certificate replacement, concurrency, trust, redirect, and
expiration tests continue to pass.

Validation, run from the [azcore module](../azcore):

| Command | Result |
| --- | --- |
| `go test ./runtime ./arm/runtime ./policy ./internal/exported -count=1` | Passed |
| `go test ./runtime ./arm/runtime -run '^(TestMTLSPoP\|TestBearerPolicy_MTLSPoP)' -count=5` | Passed |
| `go build ./...` | Passed |
| `go vet ./runtime ./arm/runtime ./policy ./internal/exported` | Passed |

No live hardware tests were run for this phase; the native attestation
implementation was not changed. The race detector remains unavailable with the
current `CGO_ENABLED=0` toolchain configuration.

### 3. Add optional activation and credential selection

**Status: implemented and validated offline.**

- Add the `azidentity/attestation` module and documented blank-import package.
- Initialization only registers the integration. It must not load native
  libraries, make network calls, or create keys.
- Keep native attestation in MSAL; do not introduce a second implementation.
  Add an MSAL hook only if the existing adapter cannot express the agreed design.
- Keep dependency direction one-way: the parent identity module must not depend
  on its optional activation module.
- Have MIC and DAC consult registration and MSAL capabilities for eligible requests.
- Apply the agreed opt-out shape and preserve configuration immutability.
- Preserve the one-second DAC discovery context, disabled probe retries,
  earlier caller deadlines, and the original acquisition context.
- Keep the optional package importable on unsupported platforms without an
  initialization failure.

**Exit criterion:** imported and unimported applications exhibit the agreed
selection behavior. Test them in separate processes or binaries so registration
cannot contaminate the no-import case.

Implemented in the [optional module](attestation/README.md) and the existing MIC/MSAL
adapter. The parent module does not import the optional module. The registration hook
only records availability, and credentials read it lazily during acquisition so
construction during package initialization cannot permanently miss activation.
The MIC opt-out is copied at construction. The superseded unreleased
`EnableAttestation` fields were removed from MIC and DAC.

Selection requires the request signal, registration, permission, and successful MSAL
discovery of at least KeyGuard binding strength. Selected acquisitions pass PoP,
attestation support, and minimum KeyGuard strength together to MSAL. Known unsupported
hosts use bearer; discovery failures remain failures. Selected attestation errors
don't enter legacy IMDS bearer error classification or continue the credential chain.
DAC reuses its initial bounded-probe capabilities for eligibility and retains the
original acquisition context. No MSAL or native implementation changes were required.

Parent tests cover no-import behavior, selection gates, immutable opt-out, discovery
failure/cancellation, native-attestation failure, IMDS attestation rejection, binding
certificates, and DAC probe reuse. A separate imported-module test binary verifies real
registration and public API behavior against mocked IMDS, without native operations.
The smoke app now blank-imports the module and no longer accepts `-attestation`;
its explicit resource PoP signal remains until service integration is completed.

Validation from the relevant module, with `AZURE_RECORD_MODE=live` for fixture-free
parent tests to avoid starting the recording proxy:

| Command | Result |
| --- | --- |
| `go test .\ ./internal/attestation -run '^(TestDefaultAzureCredential\|TestManagedIdentityCredential\|TestManagedIdentityClient\|TestChainedTokenCredential)' -count=1` | Passed |
| `go test .\ -run '^(TestManagedIdentityCredential_(Attestation\|CapabilityFailure\|MTLSPoP)$\|TestDefaultAzureCredential_(Attestation\|AttestedProbe\|IMDS)$)' -count=3` | Passed |
| Optional module: `go test ./... -count=3` and `go vet ./...` | Passed |
| Parent module: `go build ./...` and `go vet ./...` | Passed |
| Smoke module: `go test ./... -count=1` and `go build -trimpath -o .\mtls-pop-smoke.exe .` | Passed |
| Smoke module: `go test -mod=readonly -count=1 -run '^(TestDefaultAzureCredential\|TestManagedIdentityCredential\|TestManagedIdentityClient\|TestChainedTokenCredential)' github.com/Azure/azure-sdk-for-go/sdk/azidentity` | Passed against local MSAL feature checkout |
| Parent and optional modules: `go build ./...` with `GOOS=linux`, `GOARCH=amd64`, then with `GOOS=darwin`, `GOARCH=arm64` | Passed |
| Optional module: Linux `go test -c` | Cross-compilation passed; not executed |

Optional-module validation used a temporary local workspace with a version-specific
replacement for the unreleased parent identity version. The releasable optional
module has no replacements. Final dependency versions and release ordering remain
phase 6 work.

No phase-3 hardware-backed live tests or race-detector runs were performed.
Cross-platform checks are compilation only. At the phase-3 boundary, automatic
Key Vault signaling, tenant fallback, and binding-error retries remained phase-4 work;
they are now implemented below.

### 4. Enable automatic behavior in Key Vault

**Status: implemented and validated offline.**

Update [shared Key Vault authentication](../security/keyvault/internal/challenge_policy.go):

- Signal PoP support automatically for authenticated requests.
- Preserve unauthenticated challenge behavior and challenge-resource verification.
- Support ordinary non-managed-identity credentials without breaking existing apps.
- Implement C#-equivalent `AADSTS3921996` classification, fresh bearer acquisition,
  and tenant-keyed memory scoped to the policy's credential wrapper/client lifetime.
  Preserve cancellation and mixed credential-chain failure distinctions.
- Match C# binding-error classification for bounded outer-policy retries on the
  specific HTTP 401 validation marker. A binding mismatch is not a reason to
  switch to bearer.
- Start with secret reads, then verify every Key Vault client consuming the
  shared authentication policy.
- Keep changes in shared or hand-maintained integration points where possible,
  rather than manually changing generated patterns.

**Exit criterion:** real Key Vault SDK clients pass offline positive, ineligible,
tenant-denial, and non-downgrade cases without customer enable flags.

The [shared policy](../security/keyvault/internal/README.md) now signals PoP on
authenticated requests, wraps the supplied credential for per-client tenant memory,
and installs an SDK-internal operation-scoped response classifier. All seven consumers
use their existing constructors: Secrets, Keys, Certificates, and Administration's
Backup, EKM, RBAC, and Settings clients. No generated implementation changes were needed.
The smoke app no longer sets the resource client's PoP enable flag.

#### Internal bridges and Go adaptations

- `sdk/internal/errorinfo` supplies authentication, unavailable, and credential-chain
  markers. Identity errors implement standard `Is` and multi-cause `Unwrap` methods.
  Existing error types and formatted messages are preserved. Nested chains explicitly
  preserve terminal-failure precedence over unavailable causes.
- Key Vault classifies individual credential-chain causes, never a flattened aggregate
  message. Ordinary single-error Go wrappers are supported; unclassified joins,
  cancellation, and unrelated terminal failures prevent fallback. Unlike C#'s exception
  inheritance, Go represents a terminal managed identity denial as
  `AuthenticationFailedError` within a chain. A known SDK chain permits fallback only
  when every cause is unavailable or the specifically classified tenant denial.
- `sdk/internal/exported.ResponseClassifier` connects authentication and outer retry
  decisions without a public response-classifier option. `azcore` doesn't contain
  Key Vault error codes. A classified binding response bypasses challenge resends
  and token expiration, leaving count, delay, retry-after, cancellation, and operation
  limits to the normal outer retry policy. Explicit `ShouldRetry` callbacks retain
  precedence; service classification can add responses beyond the status-code list.
- Tenant keys are case-insensitive, and memory is local to the client's credential
  wrapper, including after a failed bearer acquisition. Fresh acquisition preserves
  claims, scopes, tenant, CAE, and context while disabling PoP. Binding certificates
  and bound headers are cleared when the actual returned token is bearer.
- Classification matches the complete case-sensitive `AADSTS3921996` code with
  Unicode word boundaries. Binding retries require HTTP 401, a bound request, and
  the exact bracketed marker in JSON `error.message`. Malformed or differently shaped
  JSON remains an ordinary service error. Body-read errors aren't suppressed.

#### Validation

Validation used a temporary workspace containing `sdk/internal`, `azcore`, `azidentity`,
the optional integration, Key Vault shared internals, and all four Key Vault client
modules, including the version-specific replacement for unreleased identity.
No local replacements were added to releasable manifests. The standalone smoke app
has application-local replacements and was also tested with `GOWORK=off`.
Parent identity and public client unit runs used `AZURE_RECORD_MODE=live` only to
avoid starting the recording proxy; no Azure service or native attestation was called.

| Scope and command | Result |
| --- | --- |
| Core: `go test ./runtime ./arm/runtime ./policy ./internal/exported -count=1` | Passed |
| Core: `go test ./runtime -run '^(TestRetryPolicy_ResponseClassifierError\|TestMTLSPoP\|TestBearerPolicy_MTLSPoP)' -count=3` | Passed |
| Core: `go test ./runtime -run '^(TestBearerTokenPolicy\|TestRetryPolicy_ResponseClassifierError)' -count=1` | Passed |
| Identity: `go test .\ -run '^(TestAuthenticationFailedError\|TestChainedTokenCredential\|TestDefaultAzureCredential\|TestManagedIdentity)' -count=1` | Passed |
| Shared internals: `go test github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo github.com/Azure/azure-sdk-for-go/sdk/internal/exported -count=1` | Passed |
| Optional integration: `go test github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation -count=3` | Passed in a separate imported binary |
| Key Vault shared module: `go test ./... -count=3` | Passed, including cold POST replay through binding retry |
| Secrets: `go test .\ -run '^TestClientMTLSPoP$' -count=3` | Passed, including real TLS and late tenant transition |
| All Key Vault modules: `go test @kv -run '^(TestClientMTLSPoP\|TestClientSignalsMTLSPoP\|TestClientsSignalMTLSPoP\|TestDisableChallengeResourceVerification)$' -count=3` | Passed for all seven shared-policy consumers |
| Affected modules: `go build @packages` and `go vet @packages` | Passed |
| Smoke module, `GOWORK=off`: `go test ./... -count=1` and `go build -trimpath -o <temporary-executable> .` | Passed without replacing the existing smoke executable |
| Affected local modules: `go build ./...` with `GOOS=linux`, `GOARCH=amd64`, and `GOOS=darwin`, `GOARCH=arm64` | Passed; compilation only |
| Secrets: Linux `go test -c` | Cross-compilation passed; not executed |

`@kv` contains the recursive import patterns for `azsecrets`, `azkeys`, `azcertificates`,
and `azadmin`. `@packages` additionally includes `azcore`, `azidentity`, the optional
integration, Key Vault shared internals, and the two changed SDK-internal packages.
An initial broad macOS import wildcard also selected the unchanged optional identity
cache module, which requires unavailable cgo accessors. Restricting cross-builds to
the affected local module directories passed.

No phase-4 hardware-backed live tests or race-detector runs were performed.
At the phase-4 boundary, final API/docs cleanup remained phase 5; version floors,
coordinated publication, and the live parity matrix remain phase 6.

### 5. Finish public APIs, examples, and smoke coverage

**Status: complete; offline validation passed.**

- Finalize the minimal parity-aligned opt-out API and remove superseded unreleased
  enable flags rather than maintaining compatibility aliases.
- Add module documentation and examples showing the blank import with ordinary
  credential and service-client construction.
- Document process-wide activation, transitive imports, native prerequisites,
  unsupported environments, and exact fallback behavior.
- Update relevant changelogs.
- Make the primary [smoke app](testdata/mtls-pop-smoke/README.md) scenario use the
  blank import and default options.
- Preserve separate no-import and opt-out validation.
- Keep local dependency replacements out of releasable module manifests.

**Exit criterion:** the documented customer flow works without attestation or
PoP enable flags, and documentation accurately describes its limits.

#### Public surface and customer flow

| Surface | Phase-5 decision |
| --- | --- |
| `azidentity/attestation` blank import | Sole activation mechanism; process-wide and transitive, with no native work during initialization. |
| `ManagedIdentityCredentialOptions.DisableMTLSPoP` | Retained as the zero-value-friendly MIC opt-out. No additional public DAC switch. |
| `ClientOptions.EnableMTLSPoP` | Removed, without an alias. Service authentication policies own request selection. |
| MIC/DAC `EnableAttestation` | Remains removed by phase 3. |
| `TokenRequestOptions.EnableMTLSPoP` | Retained for direct callers and PoP-aware service authors; false by default. Not a customer option required for ordinary Key Vault calls. |
| `AccessToken.BindingCertificate` | Retained as the actual bearer/bound result discriminator. |
| Registration and response/error bridges | Remain SDK-internal; no customer provider or strict-mode API. |

The [compiled examples](attestation/example_test.go) demonstrate the blank import
with `NewDefaultAzureCredential(nil)` and `azsecrets.NewClient(..., nil)`, and
the MIC opt-out. README and package documentation cover lazy process-wide activation,
transitive imports, prerequisites, known unsupported hosts, and the exact tenant-denial
exception. Only the optional integration's tests/examples gain a Key Vault test
dependency; production activation remains a small registration hook.

Removing the old client flag also removes its eager preparation policy and
certificate-free transport pool. Bearer and unauthenticated requests now always use
the original transport, including an explicitly configured caller certificate.
With ordinary defaults, the cold challenge remains certificate-free. Bound requests
still require HTTPS, use certificate-specific pools, and reject incompatible transports
or redirects without reacquiring bearer tokens. Orphan bound authorization headers
are rejected even without any enable flag.

The [smoke app](testdata/mtls-pop-smoke/README.md) now starts with two real SDK reads
on a fresh client and credential, before direct token calls. This also makes DAC's
first selection follow the automatic service path. Its wrapper validates bindings
and selected UAMI client IDs. Explicit direct-token calls remain separate regression
checks. The smoke executable deliberately requires bound tokens; separate SDK test
binaries preserve unimported activation and MIC opt-out controls without converting
legitimate bearer outcomes into positive PoP smoke results.

#### Validation and source-only development

| Scope and command | Result |
| --- | --- |
| Core: `go test ./runtime ./arm/runtime -run '^(TestMTLSPoP\|TestBearerPolicy_MTLSPoP)' -count=1` | Passed after removal of the client flag |
| Core: `go test ./runtime ./arm/runtime ./policy -count=1` | Passed |
| Parent identity: `go test . -run '^(TestManagedIdentityCredential_Attestation\|TestDefaultAzureCredential_Attestation)$' -count=3` | Passed in the unimported binary |
| Optional integration: `go test ./... -count=3` and `go vet ./...` | Passed, including compiled examples, SDK-driven MIC/DAC discovery, and explicit MIC opt-out |
| Key Vault shared module: `go test ./... -count=1` | Passed |
| All Key Vault modules: `go test @kv -run '^(TestClientMTLSPoP\|TestClientSignalsMTLSPoP\|TestClientsSignalMTLSPoP\|TestDisableChallengeResourceVerification)$' -count=1` | Passed for all seven shared-policy consumers |
| Each affected local module: `go build ./...` and `go vet ./...` | Passed |
| Smoke module, `GOWORK=off`: `go test ./... -count=1` and `go build -trimpath -o <temporary-executable> .` | Passed; existing executable preserved |
| Each affected local module: `go build ./...` for Linux amd64 and Darwin arm64 | Passed; compilation only |
| Optional integration: `go test -c` for both cross targets | Tests and examples compiled; not executed |

`@kv` denotes the same four module patterns as phase 4. Builds and vet used the nine
affected local module directories rather than broad import patterns that select the
unrelated optional cache module. Parent and public-client unit runs used
`AZURE_RECORD_MODE=live` only to avoid the recording proxy; they didn't contact Azure.

The identity module's existing local core replacement was removed. Releasable manifests
now contain no local replacements; the standalone smoke app retains its application
replacements. The unpublished parent identity version prevents ordinary isolated
`go mod tidy` in the optional module. Dependency tidying therefore used a temporary
modfile with local replacements, stripped before writing the generated manifest back.
Published Key Vault test dependency checksums were recorded separately.

Source validation still requires matching local modules in a workspace or application
replacements until matching feature-bearing dependency versions are published. The current VS Code
language-server configuration still selects published core for identity and excludes the
nested optional module, so its diagnostics aren't clean despite the successful workspace
builds and vet. Configure the editor with the same local workspace before using its
diagnostics as a validation signal; the [source-only setup](attestation/README.md#source-only-development)
documents the commands. No editor settings or ignore rules were changed.

No phase-5 hardware-backed live tests or race-detector runs were performed.
Phase 6 retains dependency/version coordination, the live activation/opt-out controls,
and the full on-demand hardware parity matrix.

### 6. Run parity validation and coordinate releases

**Status: completed for the available fixture on 2026-11-18. Dependency version
alignment and publication are explicitly deferred, not completed release work.**

#### Activation controls

The standalone smoke app now has separate default-imported and `noattestation`
builds, plus an imported MIC `-disable-mtls-pop` control. These are test-app
controls, not another customer activation mechanism or a public DAC option.
The unimported dependency graph excludes the optional integration package, and
each binary verifies its registration state.

Controls require unbound bearer results for both ordinary and PoP-requested
direct acquisition, then two real SDK reads with bearer authorization and no
token-binding header. Their final success messages are distinct from the
positive PoP smoke result. Local TLS tests also reject client certificates,
bound results, and mutation of caller options. Concurrency and identity-isolation
flags cannot be combined with bearer controls.

#### Offline qualification

Source validation selected all nine affected SDK modules in a temporary local
workspace; the smoke app used its own replacements with `GOWORK=off`.
`AZURE_RECORD_MODE=live` avoided the recording proxy for fixture-free unit
tests; it did not turn those tests into Azure live tests.

| Scope | Validation and result |
| --- | --- |
| Parent identity | Focused MIC attestation/binding/PoP and transport-factory, DAC discovery/attestation, authentication-error, and credential-chain suites passed |
| Optional attestation module | `go test github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation -count=3` passed, including imported SDK MIC/DAC selection and MIC opt-out |
| Shared Key Vault module | `go test github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal/... -count=3` passed, including exact fallback classification, tenant isolation, cancellation, and bounded retries |
| Core and SDK internals | Full core runtime, ARM runtime, policy, core internal exported, SDK error-info, and SDK internal exported suites passed |
| Public Key Vault modules | Phase-4 selectors for `TestClientMTLSPoP`, request signaling, and disabled challenge verification passed across all seven consumers |
| Nine affected SDK modules | Directory-local `go build ./...` and `go vet ./...` passed |
| Smoke module | `go test ./... -count=1`, `go test -tags noattestation ./... -count=1`, and both corresponding `go vet` commands passed |
| Smoke no-import dependency graph | `go list -tags noattestation -deps -f '{{.ImportPath}}' .` confirmed the optional integration package was absent |
| MSAL managed identity | Cache separation, concurrent attested/plain acquisition, minimum-strength partitioning, C# cache-component parity, and force-refresh regressions passed three repetitions |
| Linux amd64 and Darwin arm64 | All nine SDK modules and both smoke variants cross-built successfully; compilation only, not execution |

The focused parent selector was:

```powershell
go test github.com/Azure/azure-sdk-for-go/sdk/azidentity `
    -run '^(TestManagedIdentityCredential_(Attestation|BindingCertificate|MTLSPoP.*)|TestManagedIdentityClient_MTLSFactory.*|TestDefaultAzureCredential_(IMDS|Attestation|AttestedProbe|UnexpectedIMDSResponse)|TestAuthenticationFailedError.*|TestChainedTokenCredential.*)$' -count=1
```

The additional MSAL selector, run from its checkout with `GOWORK=off`, was:

```powershell
go test .\apps\managedidentity `
    -run '^(TestIMDSv2(AttestedAndNonAttestedCertificatesDoNotShareCache|BoundAndBearerTokensDoNotShareCache|ConcurrentAttestedAndPlainAcquisitions)|TestForceRefresh(BypassesTokenCache|ReissuesCertificateAndReusesKeyAndAttestation)|TestCacheKeySeparatesAttestation|TestCacheComponentsPartitionBy(MinStrength|Attestation)|TestCacheComponentValuesMatchDotNet)$' -count=3
```

No race-detector run was performed (`CGO_ENABLED=0`). Portability checks used
module directories, avoiding unrelated cached optional modules.

#### Windows live qualification

Runs used the approved Windows KeyGuard/IMDSv2 VM, token-binding-enabled public
Azure Key Vault, and authorized system and user-assigned identities. Fixture
identifiers and secret values are not embedded in reusable documentation.
The Go smoke app used MSAL checkout `3a911b0`; its pre-existing modified live
test was reused without changes. This is local-source qualification, not
qualification of published package versions.

| Scenario | Result |
| --- | --- |
| MIC, system-assigned identity | Passed cold automatic SDK reads, bearer/PoP/bearer isolation, matching-certificate success, and missing-certificate enforcement |
| MIC, UAMI | Passed the same checks, including selected-client-ID validation |
| DAC, system-assigned identity and UAMI | Both passed the same checks with `AZURE_TOKEN_CREDENTIALS=prod`, including the cold managed identity discovery path |
| Shared and multiple clients | Four-way waves passed in each of the four MIC/DAC identity combinations |
| Same-process identity isolation | Alternating system/UAMI SDK and direct acquisitions passed without crossing principals |
| Multiple processes | System and UAMI smoke processes were observed running simultaneously; both completed their four-way client checks and exited zero |
| Real certificate reissuance | MSAL `TestIMDSv2BoundTokenIsRejectedWithoutCertificate` passed: original and reissued pairs worked; missing and mismatched certificates were rejected; a subsequent fresh Go SDK smoke run passed |
| Sequential Go -> C# -> Go | Passed separately for system identity and UAMI; each C# run reported 15 passed, zero failed, 11 skipped, and exit zero |
| No-import controls | MIC and DAC each passed for system identity and UAMI: four bearer-control runs |
| Imported explicit MIC opt-out | Passed for system identity and UAMI: two bearer-control runs |

The reissuance test requires distinct certificates and binding-specific mismatch
rejection, or the exact TLS bad-certificate alert after the client supplied its
certificate; an arbitrary 403 or transport failure is not accepted. The initial
multi-process orchestration completed both apps but could not report exit codes
reliably through PowerShell `Start-Process`; a repeat with retained .NET process
handles verified both zero exit codes. Neither that harness issue nor a corrected
test-argument quoting error was an SDK test failure.

The verified C# bundle used the reference source pin
`f75fadb8d36d8fdf46b03d11193b8dd6009f498f`, Azure.Core 1.63.0,
MSAL/KeyAttestation 4.90.0, and native KeyGuardAttestation 1.1.7.
Per-case reports and safe Go console output were retained in session evidence.
No shared keys/certificates were deleted and no reboot was performed.

#### Qualification limits

- Real denied-tenant fallback and forced binding-validation HTTP 401 retries were
  not exercised: no approved denied-tenant/fault-injection fixture was available.
  Their exact classification, fresh bearer acquisition, tenant memory, and retry
  bounds have offline coverage.
- Live resource reads covered Secrets, not optional Keys, Certificates, or Managed
  HSM fixtures. All seven shared-policy consumers passed offline integration checks.
- Native dependency failures, unsupported hosts/transports, and cancellation have
  offline coverage; no live dependency removal, outage, TLS/revocation fault, RBAC
  mutation, or denied-identity experiment was performed.
- C# skips are not passes. Its optional metadata, denial, rotation, and destructive
  fault scenarios were not qualified by the 15 passing cases.
- Certificate reissuance does not establish natural lifetime-driven renewal,
  rotation of the underlying key, or recovery after a new reboot. The earlier
  post-reboot baseline below is historical, not a new phase-6 result.
- Linux/macOS execution and the race detector remain unverified.

#### Release ordering and remaining human work

**Agreed scope: documentation only.** No dependency manifests, sums, package
versions, telemetry constants, or release automation were changed in phase 6.
The current dependency floors do not select all unpublished feature-bearing
modules in isolation; successful workspace builds are not release readiness.

Before publication, release owners must choose compatible versions and align
runtime dependency floors in this order:

1. SDK `internal`: publish the authentication-error and response-classification
   bridges.
2. `azcore`: require that SDK-internal version and publish adaptive token/transport
   handling and response classification.
3. `azidentity`: require matching core/internals and publish automatic selection
   and the optional-module registration contract.
4. Key Vault shared `internal`: require matching core/internals and publish the
   automatic signal, tenant fallback, and binding-response classification.
   This step can proceed independently of identity after its own prerequisites.
5. Public Key Vault clients: raise their shared/core floors and any test/client
   dependencies needed for the selected versions, then publish. Add new
   unreleased versions where necessary rather than editing released changelog entries.
6. Optional `azidentity/attestation`: require the published parent identity and
   matching core/internals. Its example/test dependency graph also includes
   `azsecrets` and Key Vault internals; publish only after every version selected
   by that graph is available. It need not wait for unrelated public clients.

This is a dependency graph, not a requirement to serialize independent branches.
SDK internals already have older core/identity test dependencies: do not blindly
raise them and introduce a release bootstrap cycle. Distinguish runtime floors
from test/example dependencies, especially in the optional module. Confirm the
shipping MSAL/native versions provide the behavior tested with the smoke app's
local MSAL replacement.

Human follow-up includes version/floor assignment, sums, changelogs and telemetry
constants as appropriate, isolated builds/tests with published prerequisites and
no local workspace/replacements, and the normal
[release checklist](../../documentation/development/release.md).
Release readiness, publishing, and approvals were not invoked.

**Exit criterion met for the agreed scope:** focused offline suites, builds, and
static checks passed; available live scenarios passed; unavailable scenarios and
the deferred publication prerequisites are explicitly recorded.

## Validated baseline

The following predates the automatic-selection work and remains the regression
baseline:

- Explicit managed identity attestation and mTLS PoP acquisition.
- MIC and DAC live success for system-assigned identity and UAMI.
- DAC discovery delegated to MSAL under a one-second context deadline.
- Direct and automatic Key Vault reads, bearer/PoP isolation, and
  missing-certificate rejection.
- Concurrent clients/processes and same-process identity isolation.
- Real MSAL certificate reissuance with matching-pair success and mismatched-pair
  rejection.
- Recovery from the observed post-reboot VBS-key failure and fresh-process reuse.
- Sequential Go -> C# -> Go coexistence. The C# runner reported 15 passes,
  zero failures, and 11 skips; skips are not verified coverage.
- Focused offline SDK/MSAL regressions, repeated smoke tests, builds, and vet.

These results used local feature builds and the available Windows fixture;
they are not qualification of published packages or a broader platform matrix.
The Go race detector was not run because the native compiler prerequisites were
unavailable. The nested smoke module is not discovered by ordinary SDK module
test runs; routine coverage belongs in normal module suites.

## Reference implementation

Use the tested versions as the initial comparison point rather than assuming
current upstream defaults are identical. These versioned sources define the
parity target; changing that target requires an explicit decision:

- [Azure.Core 1.63.0 managed identity selection](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Identity/ManagedIdentityClient.cs)
- [Azure.Core 1.63.0 MSAL adapter and optional attestation discovery](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Identity/MsalManagedIdentityClient.cs)
- [Azure.Core 1.63.0 credential option defaults](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Identity/Credentials/ManagedIdentityCredentialOptions.cs)
- [Azure.Core 1.63.0 DAC option visibility and defaults](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Identity/Credentials/DefaultAzureCredentialOptions.cs)
- [Azure.Core 1.63.0 DAC managed identity construction](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Identity/DefaultAzureCredentialFactory.cs)
- [Azure.Core 1.63.0 token request defaults](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/TokenRequestContext.cs)
- [Azure.Core 1.63.0 resource authentication and transport updates](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Pipeline/BearerTokenAuthenticationPolicy.cs)
- [Azure.Core 1.63.0 HTTP transport update behavior](https://github.com/Azure/azure-sdk-for-net/blob/Azure.Core_1.63.0/sdk/core/Azure.Core/src/Pipeline/HttpClientTransport.cs)
- [Test kit Key Vault request policy](https://github.com/Azure/azure-sdk-for-net/blob/f75fadb8d36d8fdf46b03d11193b8dd6009f498f/sdk/keyvault/Azure.Security.KeyVault.Shared/src/ChallengeBasedAuthenticationPolicy.cs)
- [Test kit tenant-denial classification and fallback scope](https://github.com/Azure/azure-sdk-for-net/blob/f75fadb8d36d8fdf46b03d11193b8dd6009f498f/sdk/keyvault/Azure.Security.KeyVault.Shared/src/TenantTokenBindingCredential.cs)

## Decision and implementation updates

| Date | Update |
| --- | --- |
| 2026-11-17 | Captured the staged plan. Confirmed documented blank-import activation and the `azidentity/attestation` module name. Implementation of the automatic path has not started. |
| 2026-11-17 | Confirmed all work is unreleased and its new public surface may be renamed or removed. Established C# as the north star with constraint-justified deviations only. Checked the behavior matrix against pinned sources; clarified KeyGuard eligibility, tenant fallback scope, binding-error retries, custom transports, and MIC/DAC option visibility. |
| 2026-11-17 | Agreed phase 1 resolutions: MIC opt-out without another public DAC option, minimal internal registration, explicit failure for unsupported bound transports, and initial Key Vault shared-policy scope. Public names are provisional and must not block implementation. Phase 1 design decisions are complete; phase 2 is ready to start. |
| 2026-11-17 | Completed phase 2: adaptive token authorization and lazy certificate-scoped transport selection, with bearer transport preservation and no automatic downgrade on errors. Offline runtime/ARM/policy/request suites, five focused repetitions, build, and vet passed. Service defaults, native attestation, and optional registration are unchanged; phase 3 is ready. |
| 2026-11-17 | Completed phase 3: optional blank-import registration, lazy availability lookup, automatic KeyGuard eligibility, MIC opt-out, and DAC bounded-probe reuse. Imported/unimported offline tests and portability checks passed. |
| 2026-11-17 | Approved and implemented phase 4 internal bridges. Automatic Key Vault request signaling, structured tenant-denial fallback, and bounded binding-validation retries passed focused and repeated offline tests across all seven consumers, real local TLS, builds, vet, and scoped cross-compilation. Hardware parity and release coordination remain pending. |
| 2026-11-17 | Completed phase 5: retained the MIC opt-out and direct-request signal, removed the unreleased client enable flag and eager transport preparation, added compiled blank-import/default-client examples and imported SDK selection tests, moved automatic SDK smoke checks first, and removed the identity manifest's local replacement. Separate imported/unimported controls, core and consumer regressions, builds, vet, and Linux/macOS compilation passed. Editor workspace selection and source-only dependency constraints are documented; hardware parity and publication remain phase 6. |
| 2026-11-18 | Completed phase 6 for the available fixture: added separate live bearer controls; repeated offline SDK/MSAL regressions, builds, vet, and cross-builds; passed MIC/DAC system/UAMI, client/process concurrency, identity isolation, real certificate reissuance, and Go/C# coexistence. Recorded live gaps and C# skips without counting them as passes. By agreement, documented release ordering only; dependency/version alignment and publication remain human follow-up. |
