# Key Vault Internal Module for Go

This module contains shared code for all the Key Vault SDKs, mainly the challenge authentication policy.

## Unreleased mTLS proof-of-possession integration

The shared policy automatically signals mTLS PoP support for Secrets, Keys, Certificates,
and Administration clients, including Managed HSM. It preserves cold unauthenticated
challenges, request-body protection, challenge-resource verification, and CAE handling.
Credentials returning ordinary bearer tokens continue to work without a special transport.

Eligible managed identity requests use the optional
[attestation integration](../../../azidentity/attestation/README.md).
Customer code doesn't need a Key Vault PoP enable flag. Bound tokens use their binding
certificates through the adaptive `azcore` transport.

Only a specifically classified token-acquisition denial containing the complete error code
`AADSTS3921996` permits a fresh acquisition with PoP disabled. The policy remembers this
decision case-insensitively per tenant, within that client, until the client is recreated,
even if the bearer acquisition fails. Authentication logging reports the first decision
without tenant IDs or token contents. Cancellation, unclassified aggregate messages,
unrelated terminal authentication failures, and attestation/native dependency errors
don't permit fallback.

A bound HTTP 401 whose JSON `error.message` contains
`[MtlsCnfClaimRequestDataValidationFailed]` is retriable under the normal outer retry
policy. It isn't an authentication challenge resend or bearer fallback. Configured retry
counts, delays, operation overrides, and explicit `ShouldRetry` callbacks still apply.
RBAC denials, other binding failures, and unsupported transports aren't fallback signals.

The internal bridges require matching unreleased local SDK modules for development.
Final dependency versions and release coordination are tracked in the
[implementation plan](../../../azidentity/design.md).

## Contributing
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
