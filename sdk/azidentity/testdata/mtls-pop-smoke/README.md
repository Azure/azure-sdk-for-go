# Managed identity mTLS PoP smoke test

This standalone, read-only app verifies the local `azidentity` and `azcore` changes
against [MSAL's IMDSv2/mTLS PoP flow](https://github.com/AzureAD/microsoft-authentication-library-for-go/pull/662)
using MSAL v1.11.0-beta.1, overridden by the local feature branch for reboot-recovery
validation. The identity, core, SDK-internal, and Key Vault shared modules are replaced with this checkout;
replacements in a dependency's `go.mod` aren't inherited by an application.
The executable checks those four SDK replacements in its build metadata. The optional
attestation integration is also replaced locally and blank-imported by the default
build. A [separate negative-control build](#separate-activation-and-opt-out-checks)
omits that import.
The app's `go.mod` also replaces MSAL with
`C:\git\AzureAD\microsoft-authentication-library-for-go`. This replacement is
local to the standalone test app; SDK dependency manifests are unchanged.
Update or remove it when building on another machine or validating a published
MSAL version.

The live smoke test is run on demand when a compatible VM and service fixture
are available. It is not intended for normal CI or weekly live CI, which cannot
assume those prerequisites. Fixture-free tests in the SDK module suites provide
the routine regression coverage described below.

## Prerequisites

- Go 1.26 or later to build.
- Run on an Azure Windows VM with Credential Guard/VBS and KeyGuard available,
  native IMDSv2, and an assigned managed identity. Trusted Launch or Confidential
  VM configuration alone isn't proof that all these capabilities are enabled.
- An existing **public Azure** Key Vault configured to accept certificate-bound
  tokens and an existing test secret. Ordinary Key Vault availability does not
  establish that token binding is enabled. Use an approved mTLS PoP fixture.
- The identity needs permission to get that secret, for example the Key Vault
  Secrets User role, and network access to the vault, IMDS, and Entra.
- No managed identity source override environment variables. The app refuses
  `IDENTITY_*`, `IMDS_ENDPOINT`, `MSI_*`, `AZURE_POD_IDENTITY_AUTHORITY_HOST`,
  and `MSAL_MI_DISABLE_IMDS_V2` settings that could select another source or
  disable the feature. Run in a clean shell; the app doesn't clear variables.
  Configure `NO_PROXY` for `169.254.169.254` if your environment uses a proxy.

The app doesn't provision or modify Azure resources. MSAL may create a
VBS-isolated key and persist its issued certificate in the current user's
Windows certificate store. Don't export the private key, close its handle,
or delete a shared certificate/key container to obtain a cold-cache run.

## Test tenant and identity setup

Tenant access alone isn't sufficient. Follow the feature team's onboarding
instructions to obtain an IMDSv2-capable Windows VM and a token-binding-enabled
Key Vault in the test tenant. Generic VM or vault provisioning doesn't establish
that these preview capabilities are enabled.

From an administrative workstation with Azure CLI installed, select the test
tenant and subscription:

```powershell
az login --tenant "<test-tenant-id>"
az account set --subscription "<test-subscription-id>"
```

For a system-assigned identity, enable it on the test VM if necessary and obtain
its principal ID:

```powershell
az vm identity assign -g "<vm-resource-group>" -n "<vm-name>" --output none

$principalID = az vm show `
    -g "<vm-resource-group>" -n "<vm-name>" `
    --query identity.principalId -o tsv

$vaultID = az keyvault show --name "<vault-name>" --query id -o tsv
```

For an **RBAC-enabled** test vault, grant that identity permission to read secrets:

```powershell
az role assignment create `
    --assignee-object-id $principalID `
    --assignee-principal-type ServicePrincipal `
    --role "Key Vault Secrets User" `
    --scope $vaultID `
    --output none
```

Check that each Azure CLI command succeeds before proceeding. Allow time for RBAC
propagation. If the vault uses access policies instead of RBAC, grant the identity
the secret **Get** permission through an access policy. Use an existing,
non-sensitive test secret.

For a user-assigned identity, ensure it is attached to the VM and grant its
principal ID the same vault access instead. Obtain its `principalId` and
`clientId` with `az identity show -g "<identity-resource-group>" -n "<identity-name>"`.
Pass the **client ID**, not the principal ID, to the app's `-client-id` flag.

The app runs **on the VM** and authenticates as its managed identity, not the
Azure CLI account. It doesn't need Azure CLI on the VM, a client secret, or a
tenant-ID argument.

The default build blank-imports `azidentity/attestation`, automatically enabling attested issuance
for eligible mTLS PoP requests; there is no `-attestation` flag. Deploy `AttestationClientLib.dll`
from the `Microsoft.Azure.Security.KeyGuardAttestation` package
(`runtimes/win-x64/native`) alongside the executable or into System32, and install
the [x64 Microsoft Visual C++ 2015-2022 Redistributable](https://learn.microsoft.com/cpp/windows/latest-supported-vc-redist).
The DLL isn't loaded from the working directory or `PATH`. Deploying the DLL alone
doesn't activate attestation in an application without the import. Once selected,
attestation failures fail acquisition without falling back to an unattested credential.
Ordinary bearer requests remain unchanged.

## Build and run

From this directory in the current checkout, including the challenge
authentication fix:

```powershell
$env:GOWORK = "off" # Use this app's replacements, not a parent development workspace.
go test -count=1 ./...
go build -trimpath -o .\mtls-pop-smoke.exe .
if ($LASTEXITCODE -ne 0) { throw "Build failed" }

# System-assigned identity
.\mtls-pop-smoke.exe `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>"

# User-assigned identity attached to this VM
.\mtls-pop-smoke.exe `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>" `
    -client-id "<managed-identity-client-id>"
```

Alternatively select a user-assigned identity with `-object-id` or `-resource-id`.
Only one identity selector is permitted. Attestation selection is automatic.
The default overall timeout is five
minutes; override with `-timeout 10m`. To build elsewhere and copy to the VM,
preserve the five relative SDK module replacements and build with `GOOS=windows`
and the VM's `GOARCH`. Also point the MSAL replacement at the feature checkout on
the build machine. The VM doesn't need Go to run the executable.
Copy the newly built executable to the VM and run it there from a clean
PowerShell session. Check `$LASTEXITCODE` immediately after the app exits.

### Reboot recovery

After a successful run, rerun the same command after an owner-approved VM reboot,
without deleting keys or certificates. Every check must pass again. The local
MSAL fix recognizes `NCryptOpenKey` status `0x8009003A`
(`NTE_VBS_CANNOT_DECRYPT_KEY`) as an unusable per-boot key and recreates it through
the existing race-aware provisioning path. Access-denied and other unrelated
open failures are not treated as missing keys. Offline MSAL tests cover this
classification and concurrent-replacement decisions without requiring hardware
or a reboot.

### Optional concurrency and identity isolation

These are bounded, on-demand additions, disabled unless explicitly selected:

```powershell
.\mtls-pop-smoke.exe `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>" `
    -concurrency 4 `
    -isolation-client-id "<user-assigned-client-id-guid>"
```

`-concurrency` accepts 2-8 workers (zero disables it). After the ordinary smoke
checks, it runs one concurrent wave against a shared SDK client, then another
against two SDK clients sharing the same credential. Each wave has an authorized
baseline. A failure cancels sibling requests, and all workers finish before the
app returns. The overall `-timeout` also applies to these checks.
For multi-process coverage, run two instances concurrently under the same Windows
account, selecting the system identity in one and a UAMI in the other. Keep both
processes bounded by `-timeout` and require both exit codes to be zero. This uses
normal key provisioning; do not clear shared keys or certificates beforehand.

`-isolation-client-id` independently creates a system-assigned credential and a
credential for the specified UAMI, regardless of the primary identity selector.
Both identities must be assigned to the VM and authorized to read the same
secret. It alternates system/UAMI SDK reads twice, with separate SDK clients and
credential caches, and directly reenters each credential between reads. Every
acquisition must retain its identity claims and valid certificate binding. The
UAMI's `appid` (or `azp`) must match the requested GUID; the two principals must
be distinct and belong to the same tenant. Identity claims aren't printed.
The identities may share the underlying CNG key or certificate; differing keys
are not required. These checks do not force cold native key creation.

Offline tests exercise actual concurrent local TLS reads, cancellation and
failure propagation, and alternating SDK reads with deliberately crossed
identity tokens. Like the rest of this app, they require an explicit test run
in this nested module, not ordinary SDK CI discovery.

### DefaultAzureCredential (DAC)

Add `-dac` to run the same smoke checks through `DefaultAzureCredential` instead
of directly constructing `ManagedIdentityCredential`. For this test, set
`AZURE_TOKEN_CREDENTIALS=prod`: it retains the production credential chain and
its IMDS probe, while excluding developer-credential fallback. Selecting only
`ManagedIdentityCredential` disables the probe and is rejected by this mode.

Run in a clean shell without environment/workload credentials. The app rejects
`AZURE_TENANT_ID`, client secrets/certificate paths, username/password, and
federated token files so another credential cannot satisfy the test. It doesn't
change the environment itself.

```powershell
$env:AZURE_TOKEN_CREDENTIALS = "prod"
# AZURE_CLIENT_ID must be unset (not merely empty) for system-assigned identity.
.\mtls-pop-smoke.exe -dac `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>"

$env:AZURE_CLIENT_ID = "<user-assigned-client-id-guid>"
.\mtls-pop-smoke.exe -dac `
    -client-id $env:AZURE_CLIENT_ID `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>"
```

DAC supports UAMI selection by client ID only. The environment value must match
`-client-id`; the app also checks the returned token's identity. Run
`-isolation-client-id` separately without `-dac`. `-concurrency` works with DAC.
The blank import activates the same automatic selection for DAC's managed identity credential.

The first call on a cold DAC is an ordinary Key Vault SDK read. The shared policy
requests PoP automatically, exercising the managed identity probe and bound-token
acquisition before the explicit bearer and PoP regression checks. Both the SDK
challenge cache and credential chain start cold; neither is pre-seeded.
The probe is delegated to MSAL's `Capabilities` API under a one-second context
deadline, or the caller's earlier deadline, with retries disabled. Acquisition
then uses the original context. This is a probe/network deadline, not a
one-second limit on attestation or the entire token/resource flow.

Fixture-free DAC regressions in the ordinary `azidentity` suite verify real
MSAL v2/v1 probe requests, one-second timeout/cancellation, no retries,
recovery after timeout, caller-deadline preservation, no repeated discovery
after success, reuse of bounded-probe capabilities for PoP selection, and propagation
of attestation and minimum KeyGuard strength to MSAL. Real-MSAL probe tests simulate
IMDSv1 discovery to avoid creating native keys. An unrelated HTTP response is no
longer treated as proof that an IMDS source exists.

### Certificate reissuance and cross-language coexistence

Use MSAL's existing Windows live test
`TestIMDSv2BoundTokenIsRejectedWithoutCertificate` for real certificate
reissuance. Configure `IMDSV2_E2E_VAULT` (hostname only) and `IMDSV2_E2E_SECRET`,
and select that test with the `e2e` build tag in the MSAL checkout. It requires a
working system-assigned identity, including Graph token acquisition in its
preflight. Keep `AttestationClientLib.dll` alongside the compiled test executable
(not merely in the current directory). This test must prove both original and
reissued token/certificate pairs work, the certificates differ, and presenting
the original token with the new certificate gets a binding-specific rejection.
An arbitrary 403 or transport failure is not proof.

For coexistence, run this Go smoke app, the supplied C# live runner with
`--expect pop`, then this Go app again, sequentially under the same Windows
account and with the same identity. Preserve the C# runner's verified bundle
and inspect its per-case results and skips. Never delete the shared
`KeyGuardRSAKey` container or certificates between runs. This checks sequential
cross-language reuse, not concurrent Go/C# key replacement.

MSAL force refresh reissues a certificate; it does not prove natural
lifetime-driven renewal or rotation of the underlying key. Neither scenario
requires exposing force refresh through `azidentity`.

## Required results

For the default positive run, every check must pass, the process must exit zero, and the final message must be
`PASS: all mTLS PoP smoke checks`:

1. A fresh `azsecrets.Client` with no PoP enable flag reads the secret twice, before
   any direct token acquisition. A credential wrapper requires propagation of the
   request-level PoP signal, validates every returned binding, and checks a selected
   UAMI's client ID. The client starts cold, including its initial unauthenticated
   challenge; the app doesn't pre-seed it or substitute a hand-built pipeline.
2. Ordinary token acquisition returns a nonempty, unexpired bearer token without
   a binding certificate.
3. `GetToken` with `policy.TokenRequestOptions.EnableMTLSPoP` returns a live
   binding certificate and matching private-key public component.
   The token's `cnf.x5t#S256` must equal the leaf certificate's SHA-256 thumbprint.
4. A direct HTTPS secret GET succeeds with `Authorization: mtls_pop`,
   `x-ms-tokenboundauth: true`, and the binding certificate.
5. The **same token**, same URL, and same headers on a fresh connection with
   **no certificate** produce a 4xx response naming `MissingClientCertificate`.
   A timeout, TLS failure, redirect, unrelated 403, or 5xx doesn't pass.
6. A subsequent ordinary acquisition on the same credential remains unbound,
   checking bearer/PoP cache isolation without assuming exact token reuse.

The SDK now supports narrowly classified tenant-denial fallback, but the positive
eligible-host scenario deliberately requires bound tokens. A bearer result isn't a passing PoP smoke check.
Binding-validation retries use the normal outer retry policy; this app continues to disable
automatic retries to make its live checks explicit.

JWT decoding checks the binding claim, **not the signature**. The resource
requests establish server acceptance and enforcement. Token strings, private
keys, secret values, and HTTP bodies aren't printed. A public certificate
thumbprint and token expiry are printed for correlation. Failed direct resource
requests report only the HTTP status and Key Vault error and inner error codes,
not the error message or response body.

The direct resource transport uses TLS 1.2, client renegotiation, and
`GetClientCertificate`, as required by the MSAL sample. Server certificate
verification remains enabled and redirects are refused. The automatic path
leaves this transport configuration to `azcore`.

## Failure interpretation

- Acquisition failure: inspect Windows Credential Guard/KeyGuard, IMDSv2,
  identity assignment, Entra connectivity, and resource opt-in. The SDK selects bearer
  on a known unsupported host; the positive scenario requires a bound token and reports that as a
  smoke-test failure. Selected attestation failures don't continue to another credential.
- IMDS reports "Attestation Token is missing / empty": rebuild from the current checkout,
  retaining the attestation blank import and local module replacement, with the
  attestation DLL and x64 Visual C++ runtime installed.
  If attestation itself fails, inspect the reported native-library or attestation
  error rather than disabling attestation.
- Direct request HTTP 403/404: check token-binding opt-in, network restrictions,
  RBAC propagation, and secret existence. Use the reported Key Vault error and
  inner error codes to distinguish authorization or network denials from
  certificate-binding rejection. A successful token request alone does not
  verify resource access.
- Missing-certificate request succeeds: token binding isn't being enforced.
- Missing-certificate request fails for another reason: inconclusive, not a pass.
- Automatic SDK request fails before requesting a token: the resource client's
  challenge/pipeline integration is failing, even if a separate direct request succeeds.

Offline unit tests exercise certificate/claim validation, rejection criteria,
TLS and headers, redirects, SDK flag checking, and local-module guards. They
don't prove that a machine has IMDSv2, a real KeyGuard key, or a token-binding
enabled vault. A live run on the fixture is required for end-to-end verification.

### Challenge authentication regression coverage

The offline `TestKeyVaultSDKColdChallenge` test runs the real `azsecrets.Client`
against a local TLS server. It requires one certificate-free, unauthenticated
challenge followed by two successful, certificate-authenticated secret reads,
with the PoP flag reaching the credential and its token cached.

The [mTLS transport](../../../azcore/runtime/policy_mtls.go) leaves unauthenticated
challenges on the caller's original transport. With ordinary transport defaults,
no client certificate is sent. An explicitly configured client certificate remains
the caller's choice; no binding certificate is installed before token acquisition.
Bound requests use separate certificate-specific connection pools with HTTPS,
server certificate verification, redirect rejection, and transport validation.
This covers the failure originally reproduced on HEAD `dd45fa1ed0`, where the
initial challenge failed before the credential was called.

### Separate activation and opt-out checks

The default executable includes the blank import and requires bound tokens.
Build a separate executable without the import, and run explicit MIC opt-out on
the imported executable:

```powershell
$env:GOWORK = "off"
go test ./... -count=1
go test -tags noattestation ./... -count=1
go build -trimpath -tags noattestation -o .\mtls-pop-smoke-noimport.exe .
if ($LASTEXITCODE -ne 0) { throw "No-import build failed" }

.\mtls-pop-smoke-noimport.exe `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>"

# Use the default imported executable built above.
.\mtls-pop-smoke.exe -disable-mtls-pop `
    -vault-url "https://<vault-name>.vault.azure.net" `
    -secret-name "<existing-test-secret>"
```

Both controls require ordinary and PoP-requested direct acquisitions to return
bearer tokens without binding certificates. A fresh SDK client must then read the
secret twice with bearer authorization and no token-binding header. The
fixture-free TLS control tests additionally reject presentation of a client certificate.
Each control must exit zero with its own final message:

- No import: `PASS: no-import control returns bearer tokens and SDK reads remain bearer`
- Imported MIC opt-out: `PASS: imported integration with explicit MIC opt-out returns bearer tokens and SDK reads remain bearer`

Repeat with `-client-id` for UAMI. The no-import executable also accepts `-dac`
with the [DAC environment requirements](#defaultazurecredential-dac); the opt-out
flag is MIC-only. Don't combine the no-import build with explicit opt-out, or
combine either control with concurrency/identity-isolation flags.
The `noattestation` build tag is private to this smoke app, not the SDK activation
contract. It removes the optional package from the compiled dependency graph;
merely keeping that module in `go.mod` does not activate it. Both variants check
their actual registration state before running.

Keep the parent identity suite and optional integration suite as separate binaries
as well. The [integration test commands](../../attestation/README.md#offline-activation-checks)
provide routine imported/unimported and opt-out coverage without hardware.

### Qualification record

The [phase-6 record](../../design.md#6-run-parity-validation-and-coordinate-releases)
documents passing Windows MIC/DAC system/UAMI runs, shared/multiple-client and
multiple-process concurrency, identity isolation, real certificate reissuance,
sequential Go -> C# -> Go for both identities, and all six live bearer controls.
It also records C# skips, unavailable live fault scenarios, compilation-only
platform checks, and deferred release prerequisites. These results qualify local
source builds, not published dependency versions.

### Offline regression coverage in normal CI

The standalone app and its tests are intentionally in a nested module under
`testdata`. Ordinary module `go test ./...` and SDK CI discovery do not include
that module; the hardware-backed smoke test remains an on-demand check, not a
CI prerequisite.
The following fixture-free regressions are in the normal SDK module suites:

- `azidentity`: automatic selection gates (including no import, opt-out, unsupported
  hosts, and discovery failures), MSAL attestation and minimum-strength options,
  binding-certificate propagation, failure without downgrade, and certificate-scoped
  acquisition pipelines that preserve middleware and caller-owned transports.
- `azidentity/attestation`: a separate imported test binary verifies activation,
  unsupported-host bearer selection, opt-out, ordinary bearer requests, and
  construction without discovery. These tests run from the optional module, not
  as part of the parent module's `go test ./...`.
- `azidentity`: on Windows, real MSAL IMDSv2 metadata requests select the requested
  client, object, or resource ID (or the system-assigned identity). A mocked
  metadata rejection stops the test before native key creation, attestation, or
  certificate-store access, and must not trigger an IMDSv1 fallback.
- `azcore/runtime`: certificate-free cold challenges followed by bound requests,
  authorization scheme and binding headers, CAE challenges, concurrent cached
  requests, bearer/PoP cache isolation, certificate expiration, redirects,
  server trust, and unsupported transports.
- `keyvault/internal`: automatic request signaling, POST body replay, strict tenant-denial
  classification, concurrent tenant memory, cancellation, and bounded binding-validation
  retries that preserve response bodies and caller overrides.
- Key Vault clients: real TLS secrets reads, real credential-chain fallback, late
  bound-to-bearer tenant transitions, and request signaling through every shared-policy
  consumer (Secrets, Keys, Certificates, and the four Administration clients).
- `azcore/runtime`: forcing SDK cache reacquisition switches the certificate
  actually presented on the TLS connection. A retried 503 and a subsequent
  cached request retain the renewed token/certificate pair. This is synthetic
  certificate replacement, not proof of real IMDS/Entra renewal or rotation.

Run the focused suites from their respective module directories:

```powershell
# sdk\azidentity
go test -count=1 -run 'MTLSPoP|Attestation|BindingCertificate|MTLSFactory' .

# sdk\azcore
go test -count=1 -run MTLSPoP ./runtime ./arm/runtime
```

These tests need no managed identity, Azure permissions, attestation DLL, or
token-binding-enabled vault. They guard SDK integration behavior without relying
on the availability of the live fixture. They cannot establish native KeyGuard
availability, successful MAA attestation, real certificate issuance/rotation, or
service-side binding enforcement. Preserve a successful live smoke run as
separate evidence, and repeat it for release or authentication-dependency
changes when the fixture is available.
