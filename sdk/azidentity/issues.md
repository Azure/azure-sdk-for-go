# mTLS PoP review issues

Issues found while reviewing the staged mTLS PoP changes, which add support via
[AzureAD/microsoft-authentication-library-for-go#662](https://github.com/AzureAD/microsoft-authentication-library-for-go/pull/662).
Each entry records the problem, the evidence, a proposed resolution, and its status.

**Status values**

- **Open**: needs a fix or a decision.
- **Won't fix**: deliberately not changing; the reason is recorded.
- **By design**: the behavior is intended; documentation may still be needed.
- **Deferred**: tracked elsewhere or after this change.

**Priority values**

- **P0**: blocks release.
- **P1**: should fix before release.
- **P2**: nice to have.

## Summary

| ID | Area | Title | Priority | Status |
|---|---|---|---|---|
| [REL-1](#rel-1-module-manifests-arent-releasable) | all | Module manifests aren't releasable | P0 | Open |
| [AZC-1](#azc-1-mtls-client-cache-is-keyed-by-certificate-pointer) | azcore | mTLS client cache is keyed by certificate pointer | P1 | Open |
| [AZC-2](#azc-2-authorization-scheme-check-scans-the-whole-token) | azcore | Authorization scheme check scans the whole token | P1 | Open |
| [AZC-3](#azc-3-arm-auxiliary-tenant-error-is-retriable) | azcore/arm | ARM auxiliary-tenant error is retriable | P2 | Open |
| [AZC-4](#azc-4-response-classifier-contract-is-implicit) | azcore/internal | Response classifier contract is implicit | P2 | Open |
| [AZC-5](#azc-5-bound-401-retries-reuse-the-cached-token) | azcore/keyvault | Bound 401 retries reuse the cached token | — | By design |
| [AZC-6](#azc-6-bound-connection-pools-arent-shared-across-clients) | azcore | Bound connection pools aren't shared across clients | P2 | Won't fix |
| [AZC-7](#azc-7-separate-bearer-and-pop-token-caches) | azcore | Separate bearer and PoP token caches | — | By design |
| [AZI-1](#azi-1-dac-imds-probe-performs-binding-strength-discovery) | azidentity | DAC IMDS probe performs binding-strength discovery | P1 | Deferred |
| [AZI-2](#azi-2-error-unwrapping-is-a-new-public-behavior) | azidentity | Error unwrapping is a new public behavior | P1 | Open |
| [AZI-3](#azi-3-attestation-activation-is-process-wide) | azidentity | Attestation activation is process-wide | P2 | Open |
| [AZI-4](#azi-4-unreachable-and-redundant-checks-in-gettoken) | azidentity | Unreachable and redundant checks in GetToken | P2 | Open |
| [AZI-5](#azi-5-designmd-ships-in-the-module) | azidentity | design.md ships in the module | P1 | Open |
| [AZI-6](#azi-6-attestation-module-doesnt-isolate-native-code) | azidentity | Attestation module doesn't isolate native code | — | Won't fix |
| [KV-1](#kv-1-tenant-denial-detection-depends-on-error-text) | keyvault | Tenant-denial detection depends on error text | — | Won't fix |
| [KV-2](#kv-2-extra-bearer-acquisition-after-tenant-denial) | keyvault | Extra bearer acquisition after tenant denial | — | Won't fix |
| [KV-3](#kv-3-empty-tenant-id-shares-one-fallback-entry) | keyvault | Empty tenant ID shares one fallback entry | P2 | Open |
| [KV-4](#kv-4-go-accepts-wrappers-that-c-rejects) | keyvault | Go accepts wrappers that C# rejects | — | By design |
| [TEST-1](#test-1-race-detector-not-run) | all | Race detector not run | P1 | Open |

---

## Release

### REL-1: Module manifests aren't releasable

**Priority:** P0 · **Status:** Open

The modules only build together through local replacements or a temporary workspace.

- [azcore/go.mod](../azcore/go.mod) replaces `sdk/internal` with `../internal`.
- [keyvault/internal/go.mod](../security/keyvault/internal/go.mod) replaces both `azcore` and `sdk/internal`.
- [azidentity/go.mod](go.mod) requires released `azcore v1.23.3` and `internal v1.13.0`.
  Neither has the new APIs, so `go build ./...` fails with
  `undefined: errorinfo.ErrAuthenticationFailed`.
- [attestation/go.mod](attestation/go.mod) requires `azidentity v1.15.1-beta.1`, which hasn't been published.

**Resolution:** release in dependency order:
`sdk/internal` → `azcore` → `azidentity` → `azidentity/attestation` and `keyvault/internal` → Key Vault clients.
Remove each `replace` directive and bump the `require` versions at each step.
Until then, validate with a temporary `go.work`.
That workspace needs a `replace` for `azidentity@v1.15.1-beta.1`.

---

## azcore and sdk/internal

### AZC-1: mTLS client cache is keyed by certificate pointer

**Priority:** P1 · **Status:** Open

`mtlsPolicy.clients` is a `map[*tls.Certificate]*http.Client`
([policy_mtls.go](../azcore/runtime/policy_mtls.go#L27), [#L61](../azcore/runtime/policy_mtls.go#L61)).
MSAL's `copyBindingCertificate` returns a new `*tls.Certificate` on every acquisition, including MSAL cache hits
(`managedidentity.go:946`, `imdsv2_api.go:379` in MSAL v1.11.0-beta.1).

As a result, every time the bearer policy acquires a token again (refresh, 401 expiry, CAE), the cache misses even when the certificate hasn't changed. Each miss:

- builds a new HTTP client and connection pool;
- forces a new TLS 1.2 handshake plus renegotiation to present the client certificate;
- leaves the old pool idle until eviction at 8 entries;
- keeps the old entry's `retainedSigner`, and so a reference to the binding key, until that eviction.

**Resolution:** key the cache by certificate identity, for example the SHA-256 of `cert.Certificate[0]`.
Keep the first `*tls.Certificate` for `GetClientCertificate`.
A rotated certificate has a different hash and still gets a new pool, so the existing rotation guarantee holds.
Add a test showing that two copies of one certificate share a client and that a different certificate doesn't.

### AZC-2: Authorization scheme check scans the whole token

**Priority:** P1 · **Status:** Open

`mtlsPolicy` is in every azcore pipeline. When a request has no binding certificate, it calls `strings.Fields` on every `Authorization` value
([policy_mtls.go](../azcore/runtime/policy_mtls.go#L40)).
That walks the whole token to find the scheme.

Measured with a 1.5 KB bearer token:

| Case | Cost per request |
|---|---|
| `mtlsPolicy` with no certificate | ~1,600 ns, 4 allocations |
| no-op policy (baseline) | ~92 ns, 2 allocations |
| `strings.Fields` alone | ~1,690 ns |
| `TrimLeft` + `strings.Cut` + `EqualFold` | ~16 ns, 0 allocations |

The scan also changes behavior. It rejects requests where a custom policy sets an `mtls_pop` header itself, such as an app doing its own PoP with a custom transport.

**Resolution:** replace `strings.Fields` with a `strings.Cut` check of the scheme.
Separately, decide whether the scan is needed at all: azcore's bearer policy always records the certificate as an operation value.
If the scan stays, document the rejection of custom `mtls_pop` headers.

### AZC-3: ARM auxiliary-tenant error is retriable

**Priority:** P2 · **Status:** Open

The new auxiliary-tenant check in
[arm/runtime/policy_bearer_token.go](../azcore/arm/runtime/policy_bearer_token.go#L82)
returns a plain `fmt.Errorf`, so the retry policy retries a deterministic failure.
`mtlsPolicy` uses `errorinfo.NonRetriableError` for similar errors.
This path is defensive: ARM never requests PoP, so it only triggers if a credential returns a bound token nobody asked for.

**Resolution:** wrap the error with `errorinfo.NonRetriableError`.

### AZC-4: Response classifier contract is implicit

**Priority:** P2 · **Status:** Open

`ResponseClassifier` ([sdk/internal/exported](../internal/exported/exported.go#L16)) can run twice on the same 401:

- once in the bearer policy ([policy_bearer_token.go](../azcore/runtime/policy_bearer_token.go#L170));
- once in the retry policy ([policy_retry.go](../azcore/runtime/policy_retry.go#L173)).

It also runs on every response in the retry policy.
The Key Vault classifier works only because `runtime.Payload` buffers the body.

The retry policy can see a classifier set by a policy further down the pipeline only because `opValues` is a shared map.
Every `Request.Next()` and `Clone()` copy carries that same map.

**Resolution:** document three requirements in the `ResponseClassifier` doc comment:

- classifiers must be idempotent;
- they must be cheap for responses that don't apply;
- they must read the body with `runtime.Payload` or otherwise make it re-readable.

Add a short comment at both call sites noting that this relies on the shared operation-value map.

### AZC-5: Bound 401 retries reuse the cached token

**Priority:** — · **Status:** By design

When the classifier returns true, the bearer policy hands the 401 back to the retry policy without expiring its token caches.
Each retry therefore sends the same bound token and certificate.

**Reason:** the Key Vault README says this is intended: a bound 401 whose message contains `[MtlsCnfClaimRequestDataValidationFailed]` "is retriable under the normal outer retry policy. It isn't an authentication challenge resend or bearer fallback."
It's treated as a transient service-side condition.
**Follow-up:** confirm against the C# behavior.

### AZC-6: Bound connection pools aren't shared across clients

**Priority:** P2 · **Status:** Won't fix (for now)

Each pipeline has its own `mtlsPolicy`, so separate Key Vault clients don't share bound connections, even when they share a base transport.
With bearer tokens they share the base transport's pool.

**Reason:** a process-wide cache keyed by base transport and certificate would add lifetime and eviction complexity for little demonstrated benefit.
Revisit if bound connection counts become a problem.

### AZC-7: Separate bearer and PoP token caches

**Priority:** — · **Status:** By design

`BearerTokenPolicy` now has both `mainResource` and `mtlsResource`
([policy_bearer_token.go](../azcore/runtime/policy_bearer_token.go#L100)).

**Reason:** `EnableMTLSPoP=false` means the caller can't accept a bound token.
With a single cache, a handler that sets the flag per request could be given a cached bound token on a request that must stay bearer.
The separation upholds the public contract. Its cost is covered in [KV-2](#kv-2-extra-bearer-acquisition-after-tenant-denial).

---

## azidentity

### AZI-1: DAC IMDS probe performs binding-strength discovery

**Priority:** P1 · **Status:** Deferred (raise with MSAL owners)

`DefaultAzureCredential`'s first IMDS probe now calls MSAL's `Capabilities()`
([managed_identity_client.go](managed_identity_client.go#L244)).
It does so even when `azidentity/attestation` isn't imported, where only `Source` is used.

Discovery also works out the binding strength. What that costs depends on the host:

| Host | Before | Now (no attestation import) |
|---|---|---|
| IMDSv1-only VM | 1 request | 3 requests: IMDSv2 probe, IMDSv1 probe, compute metadata |
| IMDSv2 Windows VM | 1 request | IMDSv2 probe, then creates and persists a KeyGuard CNG key |
| IMDSv2 Linux VM | 1 request | IMDSv2 probe; the key provider fails immediately |

All of this runs under the same 1-second `imdsProbeTimeout`.
A timeout makes DAC skip managed identity and move to the next credential.

**Decision:** rely on `Capabilities()`, which MSAL documents as intended for credential chains.
Ask the MSAL owners for an option to discover only the source and skip the binding-strength work.
If they add it, use it when attestation isn't registered.

### AZI-2: Error unwrapping is a new public behavior

**Priority:** P1 · **Status:** Open

Two azidentity error types now implement `Unwrap() []error`:

- `AuthenticationFailedError` ([errors.go](errors.go#L81)) returns the MSAL error.
- `credentialUnavailableError` ([errors.go](errors.go#L191)) returns the chain's child errors when it comes from a chain
  ([chained_token_credential.go](chained_token_credential.go#L136)).

As a result, `errors.As` and `errors.Is` now reach MSAL error types and individual chain failures.
For example, `errors.As(dacErr, &*AuthenticationRequiredError)` is now true when any child returned one.
This is useful, but it becomes a de facto contract.

**Resolution:** add explicit changelog entries describing what callers can now match.
Note that wrapped MSAL error types aren't covered by azidentity's compatibility guarantees.

### AZI-3: Attestation activation is process-wide

**Priority:** P2 · **Status:** Open

Importing `azidentity/attestation` anywhere in the build turns on attested PoP for every managed identity credential in the process, even when the import is indirect.
`ManagedIdentityCredentialOptions.DisableMTLSPoP` is the only way to opt out, and it's per credential.
A library that imports attestation forces the behavior on every app that depends on it.

**Resolution:** choose one:

- Add guidance to the attestation README and package doc that only applications, not libraries, should import the package.
- Accept the current behavior as documented.

### AZI-4: Unreachable and redundant checks in GetToken

**Priority:** P2 · **Status:** Open

In the PoP path of `managedIdentityClient.GetToken`:

- `Capabilities()` returns an error only when the context is cancelled.
  The branches for other errors after [managed_identity_client.go](managed_identity_client.go#L265) can't run.
- `capabilities.IsMtlsPoPSupportedByHost()` ([#L283](managed_identity_client.go#L283)) is redundant.
  It's true whenever binding strength is above `None`, so the `>= MtlsBindingStrengthKeyGuard` check that follows already implies it.

**Resolution:** remove the redundant call.
Either keep the error branches as defensive code with a comment, or simplify them to handle only cancellation.

### AZI-5: design.md ships in the module

**Priority:** P1 · **Status:** Open

[design.md](design.md) is a 693-line implementation plan with sections such as "Implementation plan", "Validated baseline", and "Decision and implementation updates".
It sits at the module root and would ship in the module zip.
The [keyvault/internal README](../security/keyvault/internal/README.md#L33) links to it as the "implementation plan".

**Resolution:** move the content to the PR description or an engineering doc outside the module.
Update or remove the Key Vault README link at the same time.
This file (`issues.md`) needs the same treatment before merge.

### AZI-6: Attestation module doesn't isolate native code

**Priority:** — · **Status:** Won't fix

`azidentity/attestation` only sets a process-wide flag.
The native attestation loader (`LoadLibraryEx` in `attestation_windows.go`), the IMDSv2 code, and the CNG key provider all live in MSAL's core `managedidentity` package.
They're compiled into every azidentity build whether or not the module is imported.

**Reason:** this comes from MSAL's packaging, not this change.
The module is an opt-in switch; it was never meant to isolate dependencies.
The DLL is loaded only when an attested acquisition actually runs.

---

## keyvault/internal

### KV-1: Tenant-denial detection depends on error text

**Priority:** — · **Status:** Won't fix

`classifyAuthenticationFailure` matches `AADSTS3921996` in `err.Error()` text
([tenant_token_binding_credential.go](../security/keyvault/internal/tenant_token_binding_credential.go#L98), [#L112](../security/keyvault/internal/tenant_token_binding_credential.go#L112)).

**Reason:** C# does the same.

- `Azure.Security.KeyVault.Shared/src/TenantTokenBindingCredential.cs` matches `\bAADSTS3921996\b` against each `Exception.Message`.
- Its tests deliberately reject structured `error_codes` JSON.

On the real managed identity path, Go finds the code at three levels:

1. azidentity's formatting of the response body in `AuthenticationFailedError.Error()`;
2. MSAL's `CallErr`;
3. MSAL's `entraTokenError`, which includes Entra's `error_description`.

A change to azidentity's error formatting therefore wouldn't break detection.

**Residual risk:** errors with no MSAL cause, such as third-party credentials, depend only on azidentity's formatting.
The azsecrets end-to-end test exercises that more fragile shape.
The keyvault/internal "inner message" test case covers the MSAL shape.

### KV-2: Extra bearer acquisition after tenant denial

**Priority:** — · **Status:** Won't fix

After a denial, the fallback bearer token is stored in azcore's PoP cache.
The next request has the flag cleared, so it uses the empty bearer cache and calls `GetToken` again.
The azsecrets test records the sequence `[true, false, false, …]`.

**Reason:** the second call is served from MSAL's in-memory token cache, with no network request.
Avoiding it would weaken the cache separation in [AZC-7](#azc-7-separate-bearer-and-pop-token-caches).

### KV-3: Empty tenant ID shares one fallback entry

**Priority:** P2 · **Status:** Open

The fallback memory is keyed by `strings.ToLower(tro.TenantID)`
([tenant_token_binding_credential.go](../security/keyvault/internal/tenant_token_binding_credential.go#L28), [#L46](../security/keyvault/internal/tenant_token_binding_credential.go#L46)).
If `parseTenant` returns `""`, every such tenant shares one entry.
A denial for one would switch all of them to bearer.

**Resolution:** don't record or look up an empty tenant ID; treat it as "no memory".
Alternatively, document why an empty tenant can't happen once a challenge has been parsed.

### KV-4: Go accepts wrappers that C# rejects

**Priority:** — · **Status:** By design

Go falls back through ordinary single-error wrappers above the authentication error, such as `fmt.Errorf("…: %w", authErr)`.
C#'s `arbitrary-wrapper` test requires the top-level exception to be `AuthenticationFailedException`.

**Reason:** wrapping with `%w` is normal in Go, and the code comment says so.
Text in a wrapper above the authentication error is never matched, so the divergence doesn't widen what can trigger a fallback.

---

## Validation

### TEST-1: Race detector not run

**Priority:** P1 · **Status:** Open

During review, the tests passed for azcore, sdk/internal, azidentity (through a temporary workspace), azidentity/attestation, and keyvault/internal.
They weren't run with `-race`, because cgo wasn't available.
The new code has concurrent state:

- the `mtlsPolicy` client map;
- `tenantTokenBindingCredential.bearerOnlyTenants`;
- the two token caches;
- `TestTenantTokenBindingCredential_Concurrent`.

**Resolution:** run the race detector on a machine with cgo, or confirm CI does, for all of these modules:

- `azcore/runtime/...`
- `azcore/arm/runtime/...`
- `azidentity/...`
- `security/keyvault/internal/...`
- the Key Vault client `mtls_pop_test.go` files
