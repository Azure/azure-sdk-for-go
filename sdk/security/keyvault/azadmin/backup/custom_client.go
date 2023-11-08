//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package backup

// this file contains handwritten additions to the generated code

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup/internal/pollers"
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

// NewClient creates a client that performs backup and restore operations for a Managed HSM.
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

// ErrorInfo - Internal error from Azure Key Vault server.
type ErrorInfo struct {
	// REQUIRED; A machine readable error code.
	Code string

	// full error message detailing why the operation failed.
	data []byte
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Error.
func (e *ErrorInfo) UnmarshalJSON(data []byte) error {
	e.data = data
	ei := struct{ Code string }{}
	if err := json.Unmarshal(data, &ei); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	e.Code = ei.Code

	return nil
}

// Error implements a custom error for type ServerError.
// Returns full error message
func (e *ErrorInfo) Error() string {
	return string(e.data)
}

// beginFullRestore is a custom implementation of BeginFullRestore
// Uses custom poller handler
func (client *Client) beginFullRestore(ctx context.Context, restoreBlobDetails RestoreOperationParameters, options *BeginFullRestoreOptions) (*runtime.Poller[FullRestoreResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.fullRestore(ctx, restoreBlobDetails, options)
		if err != nil {
			return nil, err
		}
		handler, err := pollers.NewRestorePoller[FullRestoreResponse](client.internal.Pipeline(), resp, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[FullRestoreResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Handler:       handler,
			Tracer:        client.internal.Tracer(),
		})
	} else {
		handler, err := pollers.NewRestorePoller[FullRestoreResponse](client.internal.Pipeline(), nil, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[FullRestoreResponse]{Handler: handler})
	}
}

// beginSelectiveKeyRestore is a custom implementation of BeginFullRestore
// Uses custom poller handler
func (client *Client) beginSelectiveKeyRestore(ctx context.Context, keyName string, restoreBlobDetails SelectiveKeyRestoreOperationParameters, options *BeginSelectiveKeyRestoreOptions) (*runtime.Poller[SelectiveKeyRestoreResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.selectiveKeyRestore(ctx, keyName, restoreBlobDetails, options)
		if err != nil {
			return nil, err
		}
		handler, err := pollers.NewRestorePoller[SelectiveKeyRestoreResponse](client.internal.Pipeline(), resp, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SelectiveKeyRestoreResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Handler:       handler,
			Tracer:        client.internal.Tracer(),
		})
	} else {
		handler, err := pollers.NewRestorePoller[SelectiveKeyRestoreResponse](client.internal.Pipeline(), nil, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SelectiveKeyRestoreResponse]{Handler: handler})
	}
}
