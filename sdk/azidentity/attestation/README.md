# Azure Identity attestation integration for Go

This optional module enables automatic KeyGuard attestation for eligible managed identity
mTLS proof-of-possession (PoP) requests made through `ManagedIdentityCredential` or
`DefaultAzureCredential`. Native attestation is implemented by MSAL.

This feature is unreleased. Dependency floors and versions remain unchanged by agreement;
the [qualification and release-ordering record](../design.md#6-run-parity-validation-and-coordinate-releases)
lists passing live scenarios, coverage limits, and the remaining human publication steps.
For local builds, use application-level replacements for this
module, `azidentity`, `azcore`, SDK internals, and shared Key Vault internals, as shown by the
[on-demand smoke app](../testdata/mtls-pop-smoke/README.md). The releasable module contains
no local replacements.

## Activate the integration

Add a documented blank import to your application and construct credentials normally:

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    _ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation"
    "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

cred, err := azidentity.NewDefaultAzureCredential(nil)
if err != nil {
    // handle error
}
client, err := azsecrets.NewClient("https://<vault-name>.vault.azure.net/", cred, nil)
if err != nil {
    // handle error
}
```

The [compiled examples](example_test.go) show the complete default flow and per-credential opt-out.

Activation is process-wide, including when imported transitively. Adding a requirement to
`go.mod` alone doesn't activate the integration. Initialization only records availability:
it doesn't load native libraries, contact endpoints, or create keys. Credentials constructed
during package initialization also consult availability lazily when requesting tokens.

## Request selection

Attested acquisition requires all of the following:

- The integration is imported.
- The token request signals mTLS PoP support.
- Managed identity PoP isn't disabled.
- MSAL successfully discovers a host supporting mTLS PoP with at least KeyGuard binding strength.

Updated Key Vault clients signal PoP support automatically; construct them with ordinary
options, for example `azsecrets.NewClient(vaultURL, cred, nil)`. Other service clients aren't
automatically opted in; service authentication policies must signal support per token request.
There is no client-level PoP enable flag.
Direct token callers set `policy.TokenRequestOptions.EnableMTLSPoP` and inspect the returned
`AccessToken.BindingCertificate` to determine whether bearer or mTLS PoP authorization is required.
See the [identity documentation](../README.md#managed-identity-mtls-proof-of-possession)
for transport requirements.

Ordinary bearer requests are unchanged. Without activation, when opted out, or when discovery
successfully identifies an unsupported host, managed identity uses ordinary bearer acquisition.
The package can be imported on unsupported platforms. Discovery failures aren't treated as
known lack of support. Once attested acquisition is selected, the credential returns failures
without retrying unattested acquisition or continuing to another credential.

The [Key Vault policy](../../security/keyvault/internal/README.md) handles the specifically
classified tenant denial `AADSTS3921996` by acquiring a fresh bearer token and remembering
the decision per tenant for that client. This isn't permission to downgrade for native
attestation failures, cancellation, RBAC denials, network failures, or unsupported transports.
The specific bound-token validation HTTP 401 is retried under the configured outer retry policy.

To opt out on an individual managed identity credential, set
`ManagedIdentityCredentialOptions.DisableMTLSPoP: true`. No additional public attestation
option is required on `DefaultAzureCredential`.

## Source-only development

Until matching preview dependency versions are published, both builds and editor
diagnostics need the local modules selected together. If you don't already have
a workspace, create one from `sdk\azidentity`:

```powershell
go work init . ..\azcore ..\internal .\attestation ..\security\keyvault\internal ..\security\keyvault\azsecrets
go work edit '-replace=github.com/Azure/azure-sdk-for-go/sdk/azidentity@v1.15.1-beta.1=.'
```

The version-specific replacement prevents attempts to fetch the unpublished parent
identity version. Keep workspace files local and out of commits, and restart the
Go language server so it uses the workspace. Don't add local replacements to the
SDK modules themselves. The standalone smoke app instead uses its own application
replacements with `GOWORK=off`.

## Offline activation checks

The parent identity suite doesn't import this module. Run its selection tests separately
from this module's suite to retain distinct imported and unimported binaries. With matching
unreleased modules selected in an application workspace:

```powershell
# From sdk\azidentity; live mode avoids starting the recording proxy for these unit tests.
$env:AZURE_RECORD_MODE = "live"
go test . -run '^(TestManagedIdentityCredential_Attestation|TestDefaultAzureCredential_Attestation)$' -count=1

# A separate binary includes the blank import, compiled examples, and SDK integration checks.
go test github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation -count=1
```

These tests use fake capability clients or local HTTP endpoints. They verify ordinary requests,
explicit MIC opt-out, unsupported hosts, and automatic Key Vault signaling without loading
native attestation libraries or contacting Azure. Actual hardware-backed parity remains an
on-demand validation step.

## Native prerequisites

Attested acquisition currently requires a Windows VM with Credential Guard/KeyGuard and IMDSv2.
Deploy `AttestationClientLib.dll` from the `Microsoft.Azure.Security.KeyGuardAttestation`
package (`runtimes/win-x64/native`) alongside the executable or into System32. Install the
[x64 Microsoft Visual C++ 2015-2022 Redistributable](https://learn.microsoft.com/cpp/windows/latest-supported-vc-redist).
The DLL isn't loaded from the working directory or `PATH`. Deploying the DLL alone doesn't
activate the integration. A missing or unusable native dependency after selection is an
acquisition error, not permission to fall back to bearer.
