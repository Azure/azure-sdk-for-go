// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

// this file contains handwritten additions to the generated code

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions

	// ServiceVersion specifies the version of the service to use. The default is ServiceVersionLatest.
	ServiceVersion ServiceVersion

	// DisableChallengeResourceVerification controls whether the policy requires the
	// authentication challenge resource to match the Key Vault or Managed HSM domain.
	// See https://aka.ms/azsdk/blog/vault-uri for more information.
	DisableChallengeResourceVerification bool
}

// ServiceVersion identifies the version of the Key Vault Secrets service.
type ServiceVersion string

const (
	// ServiceVersion20260301Preview is the 2026-03-01-preview service version.
	ServiceVersion20260301Preview ServiceVersion = ServiceVersion(version20260301Preview)

	// ServiceVersionLatest is the latest service version supported by this client.
	ServiceVersionLatest ServiceVersion = ServiceVersion20260301Preview
)

// NewClient creates a client that accesses a Key Vault's secrets. You should validate that
// vaultURL references a valid Key Vault. See https://aka.ms/azsdk/blog/vault-uri for details.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	clientOptions := options.ClientOptions
	if options.ServiceVersion != "" {
		clientOptions.APIVersion = string(options.ServiceVersion)
	}
	authPolicy := internal.NewKeyVaultChallengePolicy(
		credential,
		&internal.KeyVaultChallengePolicyOptions{
			DisableChallengeResourceVerification: options.DisableChallengeResourceVerification,
		},
	)
	azcoreClient, err := azcore.NewClient(moduleName, version, runtime.PipelineOptions{
		APIVersion: runtime.APIVersionOptions{
			Location: runtime.APIVersionLocationQueryParam,
			Name:     "api-version",
		},
		PerRetry: []policy.Policy{authPolicy},
		Tracing: runtime.TracingOptions{
			Namespace: "Microsoft.KeyVault",
		},
	}, &clientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{vaultBaseUrl: vaultURL, internal: azcoreClient}, nil
}

// ID is a secret's unique ID, containing its name and version.
type ID string

// Name of the secret.
func (i *ID) Name() string {
	_, name, _ := internal.ParseID((*string)(i))
	return *name
}

// Version of the secret. This returns an empty string when the ID contains no version.
func (i *ID) Version() string {
	_, _, version := internal.ParseID((*string)(i))
	if version == nil {
		return ""
	}
	return *version
}
