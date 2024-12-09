//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

// this file contains handwritten additions to the generated code

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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

// NewClient creates a client that accesses a Key Vault's certificates. You should validate that
// vaultURL references a valid Key Vault. See https://aka.ms/azsdk/blog/vault-uri for details.
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
	azcoreClient, err := azcore.NewClient(moduleName, version, runtime.PipelineOptions{
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
	return &Client{vaultBaseUrl: vaultURL, internal: azcoreClient}, nil
}

// ID is a certificate's unique ID, containing its name and version.
type ID string

// Name of the certificate.
func (i *ID) Name() string {
	_, name, _ := internal.ParseID((*string)(i))
	return *name
}

// Version of the certificate. Returns an empty string when the ID contains no version.
func (i *ID) Version() string {
	_, _, version := internal.ParseID((*string)(i))
	if version == nil {
		return ""
	}
	return *version
}

// ErrorInfo - Internal error from Azure Key Vault server.
type ErrorInfo struct {
	// REQUIRED; A machine readable error code.
	Code string

	// full error message detailing why the operation failed.
	data []byte
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorInfo.
func (e *ErrorInfo) UnmarshalJSON(data []byte) error {
	e.data = data
	ei := struct{ Code string }{}
	if err := json.Unmarshal(data, &ei); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	e.Code = ei.Code

	return nil
}

// Error implements a custom error for type ErrorInfo.
// Returns full error message
func (e *ErrorInfo) Error() string {
	return string(e.data)
}
