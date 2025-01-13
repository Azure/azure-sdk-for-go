//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package rbac

// this file contains handwritten additions to the generated code

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	ainternal "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions

	// DisableChallengeResourceVerification controls whether the policy requires the
	// authentication challenge resource to match the Key Vault or Managed HSM domain.
	// See https://aka.ms/azsdk/blog/vault-uri for more information.
	DisableChallengeResourceVerification bool
}

// NewClient creates a client that manages role-based access for a Managed HSM.
// You should validate that vaultURL references a valid Managed HSM. See https://aka.ms/azsdk/blog/vault-uri for details.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	authPolicy := internal.NewKeyVaultChallengePolicy(
		credential,
		&internal.KeyVaultChallengePolicyOptions{
			DisableChallengeResourceVerification: options.DisableChallengeResourceVerification,
		},
	)
	azcoreClient, err := azcore.NewClient(ainternal.ModuleName, ainternal.Version, runtime.PipelineOptions{
		APIVersion: runtime.APIVersionOptions{
			Location: runtime.APIVersionLocationQueryParam,
			Name:     "api-version",
		},
		PerRetry: []policy.Policy{authPolicy},
		Tracing: runtime.TracingOptions{
			Namespace: "Microsoft.KeyVault",
		},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: vaultURL, internal: azcoreClient}, nil
}
