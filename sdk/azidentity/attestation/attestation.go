// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package attestation enables automatic KeyGuard attestation for eligible managed
// identity mTLS proof-of-possession requests. Activate it with a blank import:
//
//	import _ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation"
//
// Activation is process-wide, including when imported transitively. Only applications
// should import this package; a library that imports it activates attestation for every
// application that depends on it. Initialization
// only registers availability; it doesn't load native libraries, contact endpoints,
// or create keys. Unsupported hosts continue to use ordinary authentication.
// Once attested acquisition is selected, failures are returned without downgrading.
//
// Key Vault clients signal PoP support automatically; ordinary credential and
// client options are sufficient. ManagedIdentityCredentialOptions.DisableMTLSPoP
// opts out for an individual managed identity credential. Direct token requests
// remain bearer requests unless TokenRequestOptions.EnableMTLSPoP is set.
//
// Native acquisition currently requires a Windows VM with Credential Guard/KeyGuard,
// IMDSv2, AttestationClientLib.dll from Microsoft.Azure.Security.KeyGuardAttestation,
// and the x64 Visual C++ runtime. A missing native dependency after selection is an
// acquisition error. Key Vault alone handles the specifically classified tenant
// denial AADSTS3921996 with fresh bearer acquisition and client-lifetime tenant memory.
package attestation

import internalattestation "github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal/attestation"

func init() {
	internalattestation.Register()
}
