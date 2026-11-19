# mTLS PoP review issues

Issues found while reviewing the staged mTLS PoP changes, which add support via
[AzureAD/microsoft-authentication-library-for-go#662](https://github.com/AzureAD/microsoft-authentication-library-for-go/pull/662).
Each entry records the problem, the evidence, a proposed resolution, and its status.

**Status values**

- **Open**: needs a fix or a decision.
- **Fixed**: resolved; the entry records what changed.
- **Not an issue**: investigation showed the concern doesn't apply; the entry records why.
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
| [AZC-1](#azc-1-mtls-client-cache-is-keyed-by-certificate-pointer) | azcore | mTLS client cache is keyed by certificate pointer | P1 | Fixed |
| [AZC-2](#azc-2-authorization-scheme-check-scans-the-whole-token) | azcore | Authorization scheme check scans the whole token | P1 | Fixed |
| [AZC-3](#azc-3-arm-auxiliary-tenant-error-is-retriable) | azcore/arm | ARM auxiliary-tenant error is retriable | P2 | Not an issue |
| [AZC-4](#azc-4-response-classifier-contract-is-implicit) | azcore/internal | Response classifier contract is implicit | P2 | Fixed |
| [AZC-5](#azc-5-bound-401-retries-reuse-the-cached-token) | azcore/keyvault | Bound 401 retries reuse the cached token | — | By design |
| [AZC-6](#azc-6-bound-connection-pools-arent-shared-across-clients) | azcore | Bound connection pools aren't shared across clients | P2 | Won't fix |
| [AZC-7](#azc-7-separate-bearer-and-pop-token-caches) | azcore | Separate bearer and PoP token caches | — | By design |
| [AZC-8](#azc-8-mtlspolicy-acts-as-a-transport-selector) | azcore | mtlsPolicy acts as a transport selector | P2 | Fixed |
| [AZC-9](#azc-9-x-ms-tokenboundauth-is-set-in-azcore-rather-than-per-service) | azcore | x-ms-tokenboundauth is set in azcore rather than per service | P2 | Fixed |
| [AZI-1](#azi-1-dac-imds-probe-performs-binding-strength-discovery) | azidentity | DAC IMDS probe performs binding-strength discovery | P1 | Deferred |
| [AZI-2](#azi-2-error-unwrapping-is-a-new-public-behavior) | azidentity | Error unwrapping is a new public behavior | P1 | Fixed (unwrapping removed) |
| [AZI-3](#azi-3-attestation-activation-is-process-wide) | azidentity | Attestation activation is process-wide | P2 | Fixed (documented) |
| [AZI-4](#azi-4-unreachable-and-redundant-checks-in-gettoken) | azidentity | Unreachable and redundant checks in GetToken | P2 | Fixed |
| [AZI-5](#azi-5-designmd-ships-in-the-module) | azidentity | design.md ships in the module | P1 | Open |
| [AZI-6](#azi-6-attestation-module-doesnt-isolate-native-code) | azidentity | Attestation module doesn't isolate native code | — | Won't fix |
| [KV-1](#kv-1-tenant-denial-detection-depends-on-error-text) | keyvault | Tenant-denial detection depends on error text | — | Won't fix |
| [KV-2](#kv-2-extra-bearer-acquisition-after-tenant-denial) | keyvault | Extra bearer acquisition after tenant denial | — | Won't fix |
| [KV-3](#kv-3-empty-tenant-id-shares-one-fallback-entry) | keyvault | Empty tenant ID shares one fallback entry | P2 | Not an issue |
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
  Neither has the new APIs, so `go build ./...` fails with undefined `errorinfo` symbols.
- [attestation/go.mod](attestation/go.mod) requires `azidentity v1.15.1-beta.1`, which hasn't been published.

**Resolution:** release in dependency order:
`sdk/internal` → `azcore` → `azidentity` → `azidentity/attestation` and `keyvault/internal` → Key Vault clients.
Remove each `replace` directive and bump the `require` versions at each step.
Until then, validate with a temporary `go.work`.
That workspace needs a `replace` for `azidentity@v1.15.1-beta.1`.

The Key Vault clients (azsecrets, azkeys, azcertificates, azadmin) must be released too, not just have their dependencies bumped.
Since [AZC-4](#azc-4-response-classifier-contract-is-implicit), each enables PoP by adding `internal.NewTokenBindingPolicy()` to its pipeline, and a client without that line stays bearer-only.
Each needs a changelog entry.
azsecrets and azcertificates currently have no Unreleased section.

---

## azcore and sdk/internal

### AZC-1: mTLS client cache is keyed by certificate pointer

**Priority:** P1 · **Status:** Fixed

`mtlsPolicy.clients` was a `map[*tls.Certificate]*http.Client`.
MSAL's `copyBindingCertificate` returns a new `*tls.Certificate` on every acquisition, including MSAL cache hits
(`managedidentity.go:946`, `imdsv2_api.go:379` in MSAL v1.11.0-beta.1).

As a result, every time the bearer policy acquires a token again (refresh, 401 expiry, CAE), the cache misses even when the certificate hasn't changed. Each miss:

- builds a new HTTP client and connection pool;
- forces a new TLS 1.2 handshake plus renegotiation to present the client certificate;
- leaves the old pool idle until eviction at 8 entries;
- keeps the old entry's `retainedSigner`, and so a reference to the binding key, until that eviction.

**Fix:** [transport_mtls.go](../azcore/internal/exported/transport_mtls.go) now keys `clients` by the leaf certificate's DER encoding (`string(cert.Certificate[0])`).
Each cache hit stores the newly acquired `*tls.Certificate` in the entry's `atomic.Pointer`, and `GetClientCertificate` reads it.
New handshakes therefore use the latest signer, for example a new OS key handle after the previous one was closed, while established connections stay pooled.
MSAL's own cache key also includes the signer's pointer; azcore can't, because every acquisition returns a new `retainedSigner`.

- An exact byte key needs no collision argument.
- Go's `m[string(b)]` lookup doesn't allocate, so a bound request avoids a per-request SHA-256 and allocation.
- The cache holds at most 8 leaf copies.

A rotated certificate has different DER and still gets a new pool.
The certificate validity check now runs before the cache lookup.

**Tests** ([policy_mtls_test.go](../azcore/runtime/policy_mtls_test.go)):

- New `TestMTLSPoPReacquiredCertificateReusesConnection` runs a pipeline whose credential returns a new `*tls.Certificate` on each of 3 acquisitions.
  It asserts the server sees 1 connection; before the fix it saw 3.
- `TestMTLSPoPTransport` now asserts that a copy shares the client and that `GetClientCertificate` then returns the copy.
  Its rotation loop uses re-issued certificates with distinct serial numbers instead of pointer copies.
- New `TestMTLSPoPReacquiredCertificateUsesLatestSigner`: the first acquisition's signer fails; the second returns the same certificate with a working signer, and the request must succeed.
  Both this and the `TestMTLSPoPTransport` assertion fail when the policy keeps the first certificate.

### AZC-2: Authorization scheme check scans the whole token

**Priority:** P1 · **Status:** Fixed (scan removed)

`mtlsPolicy` is in every azcore pipeline. When a request had no binding certificate, it scanned every header for `x-ms-tokenboundauth` and called `strings.Fields` on every `Authorization` value.
If either named PoP, it rejected the request with a non-retriable error.
`strings.Fields` walked the whole token to find the scheme.

Measured with a 1.5 KB bearer token:

| Case | Cost per request |
|---|---|
| `mtlsPolicy` with no certificate, original | ~1,600 ns, 4 allocations |
| scheme-only check (interim fix) | ~200 ns, 3 allocations |
| no-op policy (baseline) | ~92 ns, 2 allocations |

**Interim fix, since superseded:** a `hasMTLSPoPScheme` helper checked only the scheme.
It cut at the first `unicode.IsSpace` rune so that tab-separated values still matched.

**Final fix:** the scan is removed. With no certificate recorded, `mtlsPolicy` then called `req.Next()` directly; since [AZC-8](#azc-8-mtlspolicy-acts-as-a-transport-selector), `transportPolicy` sends the request with the pipeline's transport.
`hasMTLSPoPScheme`, `TestHasMTLSPoPScheme` and `TestMTLSPoPRejectsAuthenticationWithoutCertificate` are deleted.

**Rationale:**

- **The SDK's own requests can't trigger it.** `BearerTokenPolicy` records the certificate and sets the scheme and `x-ms-tokenboundauth` together on every try (bound: `mtls_pop` plus the header; bearer: `Bearer`, header deleted).
  The scan could only catch headers from other code: `PerRetryPolicies`, `policy.WithHTTPHeader`, or pipelines with no bearer policy.
- **It wasn't a security boundary.** A bound token sent without its certificate still goes only to the intended server over HTTPS, and it can't be used without the private key.
  The real guarantee, that a bound token is never sent as a bearer token, comes from the bearer policy choosing the scheme from `BindingCertificate`.
- **It blocked a legitimate setup:** an app doing its own PoP with a custom policy plus a transport configured with the client certificate.
- **C# has no equivalent check.** Azure.Core's `HttpClientTransport` never inspects `Authorization`.
  Key Vault's `UpdateTokenBoundAuthHeader` derives `x-ms-tokenboundauth` from the scheme rather than rejecting anything.

**Trade-off accepted:** custom code that sets PoP headers without a certificate now gets a 401 from the service instead of an immediate client-side error.
For Key Vault, a binding-validation 401 on such a request is classified as retriable, so it fails after the retry budget.
Revisit if this causes confusion in practice.

### AZC-3: ARM auxiliary-tenant error is retriable

**Priority:** P2 · **Status:** Not an issue

The concern was that the new auxiliary-tenant check in
[arm/runtime/policy_bearer_token.go](../azcore/arm/runtime/policy_bearer_token.go) returns a plain `fmt.Errorf`, so the retry policy would retry a deterministic failure.
This path is defensive: ARM never requests PoP, so it only triggers if a credential returns a bound token nobody asked for.

**Finding:** the error is already non-retriable.
ARM's `onRequest` is the `AuthorizationHandler.OnRequest` of azcore's `BearerTokenPolicy`, and [`BearerTokenPolicy.Do`](../azcore/runtime/policy_bearer_token.go) wraps every `OnRequest` error with `errorinfo.NonRetriableError`.

This was verified with a temporary pipeline test, since deleted:

- an ARM pipeline with an auxiliary tenant and `MaxRetries: 3`;
- a credential returning an already-expired bound token, so any retry would have to acquire a new one.

The credential was called once and the server received no requests.
Wrapping the error locally would duplicate the central handling, so no change was made.

### AZC-4: Response classifier contract is implicit

**Priority:** P2 · **Status:** Fixed

Key Vault's `ResponseClassifier` operation value was set by Key Vault's per-*retry* policy but read by the retry policy, which runs *before* it.
That only worked because every `Request.Next()` and `Clone()` copy shares one operation-value map.

Operation values are meant to flow *down* the pipeline. Every other use in the repo does this:

- the log policy's per-try state;
- `SkipBodyDownload`;
- `mtlsRequest`, from the bearer policy to `mtlsPolicy`;
- Cosmos `retryContext`.

So this reverse flow could easily break in a refactor.
The classifier's other requirements weren't written down either:

- it runs twice on a 401, once in the bearer policy and once in the retry policy;
- the retry policy calls it on every response;
- it must read the body with `runtime.Payload`.

**Compared with C#:** the flow is the same, including two evaluations and the classified 401 bypassing challenge handling.
C#, however, uses Azure.Core's existing public `HttpMessage.ResponseClassifier`.
Its Key Vault policy wraps that classifier, buffers the body explicitly, and rewinds the stream after reading.

**Alternatives considered:**

- **A. Public `Request.SetResponseClassifier`, like C#:** deferred. It would be new public API on GA azcore. It would need decisions on a public classifier type, replace-versus-compose semantics, auth-policy interaction, and `ShouldRetry` precedence, plus azcore review.
- **B. A classifier field on `runtime.PipelineOptions`:** still public, but narrower API for SDK authors.
- **C. Document the reverse flow and pin it with a test:** keeps the unusual direction.
- **D. Install the classifier from a per-call policy:** chosen.

**Fix (D):**

- **[challenge_policy.go](../security/keyvault/internal/challenge_policy.go):**
  - New `NewTokenBindingPolicy()` per-call policy sets the classifier operation value before the retry policy runs, so it flows downward like every other use.
  - The classifier is now a stateless `classifyTokenBindingFailure` that reads `x-ms-tokenboundauth` from `res.Request`, the request as sent, instead of capturing the per-try request.
  - `keyVaultAuthorizer.Do` simply delegates to the bearer policy.
  - azcore's retry and bearer policies are unchanged.
- **Compatibility gate:** the challenge policy requests PoP only for requests that passed through `NewTokenBindingPolicy`. `updateTokenRequestOptions` no longer stores `EnableMTLSPoP`; a new `tokenRequestOptions` helper sets it per request.
  Go selects one `keyvault/internal` version per build, so an already-released client can end up running this module's new challenge policy.
  Such a client stays bearer-only, rather than getting PoP with a challenge resend instead of a retry.
  `NewKeyVaultChallengePolicy`'s signature is unchanged; the module is released at v1.
  The [mtls-pop-smoke](testdata/mtls-pop-smoke) app hit exactly this case.
  It replaced keyvault/internal locally but built the released `azsecrets v1.5.0`, so the live run failed with "SDK didn't propagate EnableMTLSPoP to the credential".
  Its go.mod now also replaces azsecrets, and `verifyBuild` requires that, so a mismatched build fails at startup with a clear message.
  Its unit tests reproduce the failure without the replacement and pass with it.
- **Clients:** azsecrets, azkeys, azcertificates, and azadmin's backup, ekm, rbac and settings each add `PerCall: []policy.Policy{internal.NewTokenBindingPolicy()}`, one line in each handwritten `custom_client.go`.
- **[ResponseClassifier doc](../internal/exported/exported.go)** now states the contract:
  - set it before the retry policy;
  - it may run more than once and must return quickly when it doesn't apply;
  - read the body with `runtime.Payload`;
  - inspect `res.Request`;
  - an error replaces the response.
- **Tests:**
  - The binding-retry and automatic-PoP tests add the per-call policy.
  - New `TestChallengePolicy_MTLSPoPRequiresTokenBindingPolicy` covers the gate.
  - The `EnableMTLSPoP` assertion the baseline added to `TestChallengePolicy` is removed, because that pipeline has no per-call policy.
  - Every client's `mtls_pop_test.go` passes, confirming the wiring end to end.
- **Pre-existing bug fixed along the way:** [azadmin/mtls_pop_test.go](../security/keyvault/azadmin/mtls_pop_test.go) didn't compile. It called the beta method name `GetEkmConnection`, which was renamed to `GetConnection` before release.
- **Docs:**
  - the keyvault/internal CHANGELOG and README;
  - the azidentity CHANGELOG and README: "Updated Key Vault clients", not "shared authentication";
  - [design.md](design.md).

**Follow-up (optional):** a public response classifier on `policy.Request` (option A), matching C#, is a reasonable future feature on its own merits; it isn't required for this work.
The retry policy could also call the classifier only for error responses, as C# does; that wasn't done here.

### AZC-5: Bound 401 retries reuse the cached token

**Priority:** — · **Status:** By design

When the classifier returns true, the bearer policy hands the 401 back to the retry policy without expiring its token caches.
Each retry therefore sends the same bound token and certificate.

**Reason:** the Key Vault README says this is intended: a bound 401 whose message contains `[MtlsCnfClaimRequestDataValidationFailed]` "is retriable under the normal outer retry policy. It isn't an authentication challenge resend or bearer fallback."
It's treated as a transient service-side condition.
**Follow-up:** confirm against the C# behavior.

### AZC-6: Bound connection pools aren't shared across clients

**Priority:** P2 · **Status:** Won't fix (for now)

Each pipeline has its own certificate client cache (formerly `mtlsPolicy`, now in `transportPolicy`), so separate Key Vault clients don't share bound connections, even when they share a base transport.
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

### AZC-8: mtlsPolicy acts as a transport selector

**Priority:** P2 · **Status:** Fixed

When a request had a binding certificate, `mtlsPolicy` (formerly `azcore/runtime/policy_mtls.go`) sent it with a certificate-specific `*http.Client` and skipped the pipeline's `transportPolicy`.
In practice it was a transport selector placed as the last policy.

**Considered and rejected: moving it to azidentity.**
The certificate-specific client replaces the *service client's* transport, such as `azsecrets.ClientOptions.Transport`, not the credential's.
azidentity already handles its own connection, credential to IMDS/Entra, through MSAL's `NewMtlsClient`/`ConfigureClient`.
Moving the service-side switch into credentials would require:

- `AccessToken` to carry a transport, so every credential that returns a binding certificate builds a correct mTLS client;
- credentials to know each service client's base transport, which they don't, especially when one credential is shared;
- users to wire an azidentity transport into each service client, giving up automatic selection.

**Considered and rejected: enabling it only when `azidentity/attestation` is imported.**
That import is today's only source of bound tokens in azidentity, but `AccessToken.BindingCertificate` is public azcore API.
Custom credentials, future credentials (for example certificate-based PoP without attestation), and direct `EnableMTLSPoP` callers can all return bound tokens.
azcore can't see azidentity, so gating would need global state in `sdk/internal`, and it would make azcore's handling of `BindingCertificate` depend on an unrelated import.
The policy costs a few tens of nanoseconds per request when no certificate is present.

**Fix: the selection is folded into `transportPolicy`**, so one place chooses each request's client.

- **[pipeline.go](../azcore/internal/exported/pipeline.go):** `transportPolicy` is now a pointer type holding the base transport and the certificate client cache.
  `Do` reads the `MTLSRequest` operation value:
  - with no certificate, it uses the base transport;
  - with one, it requires HTTPS and uses that certificate's client.

  Bound requests now also get the "received nil response" check, which the separate policy skipped.
- **[transport_mtls.go](../azcore/internal/exported/transport_mtls.go):** new file holding the `MTLSRequest` operation value type (set by the bearer policy), the `mtlsClient` cache entry, and `transportPolicy.clientFor`.
  The latter is the former `mtlsPolicy.client`, unchanged apart from using the already-defaulted pipeline transport as its base.
- **runtime:** `policy_mtls.go` is deleted, `NewPipeline` no longer adds a separate policy, and the bearer policy sets `exported.MTLSRequest`.
- **Placement:** the code went into `azcore/internal/exported` rather than a new internal package, to keep client selection with `transportPolicy`.
- **Tests:** the four tests that call the cache directly moved to [transport_mtls_test.go](../azcore/internal/exported/transport_mtls_test.go): derivation and caching, unsupported configurations and HTTPS, redirects, and server trust.
  The pipeline-level tests stay in [policy_mtls_test.go](../azcore/runtime/policy_mtls_test.go) as regression coverage, unchanged.
  All of azcore, keyvault/internal, the Key Vault clients' PoP tests, azidentity, attestation, and the smoke app pass.

No public API changed.

### AZC-9: x-ms-tokenboundauth is set in azcore rather than per service

**Priority:** P2 · **Status:** Fixed

Go's `BearerTokenPolicy` set `x-ms-tokenboundauth: true` for every bound token, for any service, and deleted it for bearer tokens.
C# sets it only in Key Vault's `ChallengeBasedAuthenticationPolicy.UpdateTokenBoundAuthHeader`, which derives it from the `mtls_pop` authorization prefix after each authorization.
Azure.Core never touches it.

The header isn't part of mTLS PoP itself; RFC 8705 uses the `mtls_pop` scheme plus the client certificate.
Neither repo documents what the service does with it, so it's treated as a Key Vault protocol detail.

Before the fix, Go and C# behaved the same only because only Key Vault requests PoP.
Go's ARM pipeline would have sent the header with a bound token nobody asked for, and Go would send it to any future PoP service.

**Fix (C# parity):**

- **[policy_bearer_token.go](../azcore/runtime/policy_bearer_token.go):** no longer sets or deletes the header. It chooses only the `mtls_pop`/`Bearer` scheme and records the certificate.
- **[challenge_policy.go](../security/keyvault/internal/challenge_policy.go):**
  - A new `authorizeWith` helper wraps `authNZ` in both `authorize` and `authorizeOnChallenge`. After a successful authorization it sets the header for `mtls_pop ` authorization and deletes it otherwise.
  - This covers every authorization, including CAE: azcore's CAE path calls the handler's `OnRequest`.
  - The header name is now a `tokenBoundAuthHeader` constant shared with `classifyTokenBindingFailure`.
- **Tests:**
  - [policy_mtls_test.go](../azcore/runtime/policy_mtls_test.go): the bound-path assertions are removed, and one now asserts azcore does *not* set the header. The stale-header setup is removed, since clearing is no longer azcore's job.
  - `TestChallengePolicy_AutomaticMTLSPoP` starts each request with a stale header, and checks it's replaced with `true` for bound tokens and removed for bearer tokens. Disabling the delete makes this test fail.
  - The azsecrets "tenant transition" test still covers clearing within one request (a bound attempt, then a bearer resend after a challenge).
  - The smoke app's `TestAutomaticPipeline` builds a plain azcore pipeline, so it now asserts the header is absent. Its live checks run through azsecrets, or set the header manually in `resourceGET`, and are unaffected.
- **Docs:** the [keyvault/internal README](../security/keyvault/internal/README.md) and [design.md](design.md).

**Compatibility:** the new keyvault/internal already requires the new azcore. Released keyvault/internal versions never request PoP, so they behave the same with an azcore that no longer sets the header.

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

**Priority:** P1 · **Status:** Fixed (unwrapping removed)

The baseline gave azidentity's error types `Is` and `Unwrap() []error` methods for the first time.
That let `errors.Is` and `errors.As` reach MSAL error types, `context` errors, and each credential's error inside a chain.
Reachable errors are public API ([Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors), "Whether to Wrap").

**Why it was added:** the only consumer was Key Vault's `classifyAuthenticationFailure`. It needed to:

- recognize azidentity errors without importing azidentity;
- read the unavailable and chain flags;
- walk causes to classify each chain child and find `AADSTS3921996` in MSAL's error.

Applications never needed it.

**Alternatives considered:**

- **Keep `Unwrap` and document it:** rejected; it commits to MSAL types and chain structure.
- **Wrap returned errors in an `errorinfo` type:** rejected. Every credential's concrete return type would change, breaking `err.(*AuthenticationFailedError)` type assertions at runtime.
- **A classifier registered from `init`:** this would add no public method, but it needs global state.
- **Match on the top-level message only:** this loses per-child classification of chain errors.

**Fix:** a single `AuthenticationFailure` method, following the same pattern as the existing `NonRetriable` marker.

- **[sdk/internal/errorinfo](../internal/errorinfo/errorinfo.go):**
  - The `AuthenticationFailure` interface and the `AuthenticationFailureInfo{Unavailable, Chained, Causes}` struct replace the `AuthenticationError` interface and three sentinel errors.
  - The struct allows fields to be added later. Each field's zero value means "no" or "unknown", so a mismatched module version is safe.
  - Go's `internal` import rule (verified) keeps applications from naming the type or the interface.
- **[errors.go](errors.go):**
  - `AuthenticationFailedError` and `credentialUnavailableError` implement `AuthenticationFailure()`, documented as being for Azure SDK modules and not covered by compatibility guarantees.
  - `Is` and `Unwrap` are removed, so `errors.Is` and `errors.As` behave exactly as in 1.15.0.
  - `AuthenticationRequiredError` gets the method through embedding, as it does `NonRetriable`.
- **[chained_token_credential.go](chained_token_credential.go):** the extra `errors.As(err, &authFailedErr)` stop condition is reverted. It only existed because of `Unwrap`.
- **[tenant_token_binding_credential.go](../security/keyvault/internal/tenant_token_binding_credential.go):**
  - The classifier type-asserts each node it visits to `errorinfo.AuthenticationFailure` and reads the info struct; the classification rules are unchanged.
  - `errors.Is` can't see causes any more, so cancellation is checked on each node visited.
- **Changelogs:**
  - azidentity drops the unwrapping bullet and notes the SDK-only method under Other Changes.
  - sdk/internal describes the new `errorinfo` types.
  - [design.md](design.md) is updated to match.
- **Tests:**
  - `TestAuthenticationFailedErrorFromMSAL` and `TestChainedTokenCredential_AuthenticationFailure`, covering nested chains and all-unavailable chains, check the info returned and that causes are *not* reachable through `errors.Is`.
  - They replace the three unwrapping tests.
  - Key Vault's `authenticationTestError` now implements the new interface without `Unwrap`, matching azidentity, and the 23-case classification table passes unchanged.

**Standard library precedent for exported-but-internal API:**

- `testing.InternalTest`: "an internal type but exported because it is cross-package".
- `testing.MainStart`: "not subject to the Go 1 compatibility document"; its signature uses the unexported `testDeps` type.
- `net/http.NoBody`: an exported variable of an unexported type.

### AZI-3: Attestation activation is process-wide

**Priority:** P2 · **Status:** Fixed (documented)

Importing `azidentity/attestation` anywhere in the build turns on attested PoP for every managed identity credential in the process, even when the import is indirect.
`ManagedIdentityCredentialOptions.DisableMTLSPoP` is the only way to opt out, and it's per credential.
A library that imports attestation forces the behavior on every app that depends on it.

**Decision:** accept the design and add guidance.
A stray activation only affects requests that meet all of these:

- the request comes from a managed identity credential that hasn't opted out;
- a service policy signals PoP for it (today only updated Key Vault clients);
- the host reports KeyGuard binding.

Everything else stays bearer.
A requirement in `go.mod` alone doesn't activate the integration; the package must be imported.

**Fix:** the [package doc](attestation/attestation.go) and [README](attestation/README.md) now say that importing the package is an application-level decision and that libraries shouldn't import it.
The README also points to `DisableMTLSPoP` for keeping an individual credential on bearer authentication.

### AZI-4: Unreachable and redundant checks in GetToken

**Priority:** P2 · **Status:** Fixed

In the PoP path of `managedIdentityClient.GetToken`:

- `Capabilities()` returns an error only when the context is cancelled, so the branches for other errors can't currently run.
- `capabilities.IsMtlsPoPSupportedByHost()` was redundant.
  It's true whenever binding strength is above `None`, so the `>= MtlsBindingStrengthKeyGuard` check that followed already implied it.

**Fix ([managed_identity_client.go](managed_identity_client.go)):**

- **Redundant call removed.** Selection is now just `MaxSupportedBindingStrength >= MtlsBindingStrengthKeyGuard`.
  A comment explains why the code doesn't use `IsMtlsPoPSupportedByHost`: it's also true for software keys, and MSAL binds only KeyGuard keys.
  The existing table test covers `None` ("unsupported host"), `Software` and `KeyGuard`.
- **Error branches kept as defensive code, with a comment.** `Capabilities` is a beta MSAL API.
  Simplifying to handle only cancellation would, if MSAL ever reported another failure, return a raw MSAL error, and `DefaultAzureCredential` would stop instead of moving on.
  The branches keep the usual managed identity behavior: `credentialUnavailableError` when chained, `AuthenticationFailedError` otherwise.
  `TestManagedIdentityCredential_CapabilityFailure` already exercises them through a fake.

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

**Priority:** P2 · **Status:** Not an issue

The fallback memory is keyed by `strings.ToLower(tro.TenantID)` in [tenant_token_binding_credential.go](../security/keyvault/internal/tenant_token_binding_credential.go).
The concern was that if `parseTenant` returns `""`, every such tenant shares one entry, so a denial for one would switch all of them to bearer.

**Finding:** one shared entry is correct.

- **When it can happen:** only through `authorizeOnChallenge`, when Key Vault's challenge has no `authorization` parameter.
  `authorize` never authenticates with an empty tenant; it sends the unauthenticated request to get a challenge instead.
- **Why sharing is right:** an empty `TenantID` means "use the credential's own default tenant".
  The memory belongs to one client and wraps one credential, so every request recorded under `""` goes to the same tenant.
- **C# parity:** `TenantTokenBindingCredential.cs` keys both lookup and insert on `context.TenantId ?? string.Empty`.
- **The proposed fix would be worse:** skipping memory for `""` would make that tenant request PoP, and be denied, on every token acquisition.

No change was made.

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

- `transportPolicy`'s certificate client map;
- `tenantTokenBindingCredential.bearerOnlyTenants`;
- the two token caches;
- `TestTenantTokenBindingCredential_Concurrent`.

**Resolution:** run the race detector on a machine with cgo, or confirm CI does, for all of these modules:

- `azcore/runtime/...`
- `azcore/arm/runtime/...`
- `azidentity/...`
- `security/keyvault/internal/...`
- the Key Vault client `mtls_pop_test.go` files
