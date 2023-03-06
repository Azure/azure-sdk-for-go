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
	pl := runtime.NewPipeline(ainternal.ModuleName, ainternal.Version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &Client{endpoint: vaultURL, pl: pl}, nil
}

// ServerError - Internal error from Azure Key Vault server.
type ServerError struct {
	// REQUIRED; A machine readable error code.
	Code string

	// full error message detailing why the operation failed.
	data []byte
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Error.
func (e *ServerError) UnmarshalJSON(data []byte) error {
	e.data = data
	ei := struct{ Code string }{}
	if err := json.Unmarshal(data, &ei); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	e.Code = ei.Code

	return nil
}

// Error implements a custom error for type Error.
func (e *ServerError) Error() string {
	return string(e.data)
}

// beginFullRestore is a custom implementation of BeginFullRestore
// Uses custom poller handler
func (client *Client) beginFullRestore(ctx context.Context, restoreBlobDetails RestoreOperationParameters, options *ClientBeginFullRestoreOptions) (*runtime.Poller[ClientFullRestoreResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.fullRestore(ctx, restoreBlobDetails, options)
		if err != nil {
			return nil, err
		}
		handler, err := newRestorePoller[ClientFullRestoreResponse](client.pl, resp, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[ClientFullRestoreResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Handler:       handler,
		})
	} else {
		handler, err := newRestorePoller[ClientFullRestoreResponse](client.pl, nil, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.pl, &runtime.NewPollerFromResumeTokenOptions[ClientFullRestoreResponse]{Handler: handler})
	}
}

func (client *Client) beginSelectiveKeyRestore(ctx context.Context, keyName string, restoreBlobDetails SelectiveKeyRestoreOperationParameters, options *ClientBeginSelectiveKeyRestoreOptions) (*runtime.Poller[ClientSelectiveKeyRestoreResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.selectiveKeyRestore(ctx, keyName, restoreBlobDetails, options)
		if err != nil {
			return nil, err
		}
		handler, err := newRestorePoller[ClientSelectiveKeyRestoreResponse](client.pl, resp, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[ClientSelectiveKeyRestoreResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Handler:       handler,
		})
	} else {
		handler, err := newRestorePoller[ClientSelectiveKeyRestoreResponse](client.pl, nil, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.pl, &runtime.NewPollerFromResumeTokenOptions[ClientSelectiveKeyRestoreResponse]{Handler: handler})
	}
}
