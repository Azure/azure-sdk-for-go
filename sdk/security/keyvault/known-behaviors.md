<!-- Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License. -->

# Key Vault issue investigation context

This is advisory service context, not an automatic closure list. Read it together with `sdk/security/keyvault/TROUBLESHOOTING.md` and the affected module's README, CHANGELOG, and troubleshooting guide. Establish the exact operation and service evidence; an HTTP status alone does not establish ownership.

## Soft-delete and name reuse

A deleted vault or object can remain recoverable during its retention period, preventing reuse of its name. Verify the resource type, deleted state, and retention policy before diagnosing a conflict.

Prefer recovery or waiting where appropriate. Purging is permanent, can destroy recoverable data, and may be blocked by purge protection; never recommend it as a routine retry or perform it automatically.

- https://learn.microsoft.com/azure/key-vault/general/soft-delete-overview
- https://learn.microsoft.com/azure/key-vault/general/key-vault-recovery

## Throttling

Key Vault enforces operation limits and may return HTTP 429. Consult the Go troubleshooting guide for client/credential reuse and caching. Verify workload and retry behavior before deciding the SDK has no defect; a throttling response alone does not justify closure.

- https://learn.microsoft.com/azure/key-vault/general/service-limits
- https://learn.microsoft.com/azure/key-vault/general/overview-throttling

## Certificate import requirements

Certificate import requires a supported PFX or PEM with the appropriate matching private key and certificate chain. Inspect the service error and the documented format requirements without requesting the customer's private key or certificate bundle.

- https://learn.microsoft.com/azure/key-vault/certificates/tutorial-import-certificate
- https://learn.microsoft.com/azure/key-vault/certificates/about-certificates

## Access policies in deployments

A deployment that replaces the vault's accessPolicies property can remove policies omitted from that property. Distinguish replacement from incremental accessPolicies subresource operations.

Check the deployment definition and active permission model. Do not attribute a documented replacement operation to a Go serialization defect without source evidence.

- https://learn.microsoft.com/azure/key-vault/general/assign-access-policy
- https://learn.microsoft.com/azure/key-vault/general/rbac-migration

## Firewall restrictions

Firewall and virtual-network rules can deny requests from disallowed networks. Verify the documented network path and the specific service error. Trusted-service access does not cover every Azure service.

Recommend reviewing the intended network policy, not disabling the firewall or broadening access indiscriminately.

- https://learn.microsoft.com/azure/key-vault/general/network-security
- https://learn.microsoft.com/azure/key-vault/general/overview-vnet-service-endpoints

## RBAC propagation

Role assignments can take time to propagate. Verify the principal, scope, and time since assignment before treating a 403 as propagation delay.

- https://learn.microsoft.com/azure/role-based-access-control/troubleshooting

## Token audience and cloud

The token audience and authority must match the resource and Azure cloud. Do not assume every deployment uses public-cloud endpoints or prescribe a single hardcoded scope across clouds.

Key Vault's initial unauthenticated request can legitimately receive a 401 challenge; distinguish that expected exchange from a failed authenticated operation using the Go troubleshooting guide.

- https://learn.microsoft.com/azure/key-vault/general/authentication

## Tenant and identity selection

Authentication can select an identity other than the one granted vault access. Check the credential type and intended principal, using sanitized azidentity diagnostics rather than token contents or secrets.

- https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/TROUBLESHOOTING.md
- https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md
- https://learn.microsoft.com/azure/key-vault/general/authentication

## Purge protection

Purge protection prevents permanent deletion during the applicable retention period. Immediate name reuse or purge may therefore be unavailable by design.

- https://learn.microsoft.com/azure/key-vault/general/soft-delete-overview

## Private endpoint DNS

Private endpoint access depends on DNS resolving the vault hostname to the intended private endpoint from the customer's network. Verify DNS and connectivity evidence before ruling out a client problem.

- https://learn.microsoft.com/azure/key-vault/general/private-link-service

## Managed identity availability and access

A managed identity must be enabled on the host, selected by the credential, and granted the required data-plane permissions. Management-plane access alone does not grant access to secrets, keys, or certificates.

- https://learn.microsoft.com/azure/key-vault/general/authentication
- https://learn.microsoft.com/azure/app-service/overview-managed-identity

## Authorization model

RBAC and vault access policies are distinct permission models. Permissions in the inactive model do not grant the expected data-plane access.

- https://learn.microsoft.com/azure/key-vault/general/rbac-guide
- https://learn.microsoft.com/azure/key-vault/general/rbac-migration

## Object attributes and permitted operations

Enabled state, activation/expiration times, and permitted key operations have operation-specific semantics. Check the relevant service documentation and exact operation; do not assume all operations reject every expired secret, key, or certificate.

- https://learn.microsoft.com/azure/key-vault/secrets/about-secrets
- https://learn.microsoft.com/azure/key-vault/keys/about-keys

## Managed HSM versus vault APIs

Managed HSM and vault endpoints support different API sets. Go's azkeys module supports keys and cryptographic operations for both; azadmin targets Managed HSM administration, while azsecrets and azcertificates target vault functionality. Verify the affected module's README and actual endpoint.

- https://learn.microsoft.com/azure/key-vault/managed-hsm/overview
- https://learn.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates
