//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azadmin

// this file contains handwritten additions to the generated code

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
)

// AccessControlClientOptions contains optional settings for AccessControlClient.
type AccessControlClientOptions struct {
	azcore.ClientOptions

	// DisableChallengeResourceVerification controls whether the policy requires the
	// authentication challenge resource to match the Key Vault or Managed HSM domain.
	// See https://aka.ms/azsdk/blog/vault-uri for more information.
	DisableChallengeResourceVerification bool
}

// vaultURL references a valid Key Vault. See https://aka.ms/azsdk/blog/vault-uri for details.
func NewAccessControlClient(vaultURL string, credential azcore.TokenCredential, options *AccessControlClientOptions) (*AccessControlClient, error) {
	if options == nil {
		options = &AccessControlClientOptions{}
	}
	authPolicy := internal.NewKeyVaultChallengePolicy(
		credential,
		&internal.KeyVaultChallengePolicyOptions{
			DisableChallengeResourceVerification: options.DisableChallengeResourceVerification,
		},
	)
	pl := runtime.NewPipeline(shared.ModuleName, shared.Version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &AccessControlClient{endpoint: vaultURL, pl: pl}, nil
}

// AccessControlClientOptions contains optional settings for AccessControlClient.
type BackupClientOptions struct {
	azcore.ClientOptions

	// DisableChallengeResourceVerification controls whether the policy requires the
	// authentication challenge resource to match the Key Vault or Managed HSM domain.
	// See https://aka.ms/azsdk/blog/vault-uri for more information.
	DisableChallengeResourceVerification bool
}

// vaultURL references a valid Key Vault. See https://aka.ms/azsdk/blog/vault-uri for details.
func NewBackupClient(vaultURL string, credential azcore.TokenCredential, options *BackupClientOptions) (*BackupClient, error) {
	if options == nil {
		options = &BackupClientOptions{}
	}
	authPolicy := internal.NewKeyVaultChallengePolicy(
		credential,
		&internal.KeyVaultChallengePolicyOptions{
			DisableChallengeResourceVerification: options.DisableChallengeResourceVerification,
		},
	)
	pl := runtime.NewPipeline(shared.ModuleName, shared.Version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &BackupClient{endpoint: vaultURL, pl: pl}, nil
}

// AccessControlClientOptions contains optional settings for AccessControlClient.
type SettingsClientOptions struct {
	azcore.ClientOptions

	// DisableChallengeResourceVerification controls whether the policy requires the
	// authentication challenge resource to match the Key Vault or Managed HSM domain.
	// See https://aka.ms/azsdk/blog/vault-uri for more information.
	DisableChallengeResourceVerification bool
}

// vaultURL references a valid Key Vault. See https://aka.ms/azsdk/blog/vault-uri for details.
func NewSettingsClient(vaultURL string, credential azcore.TokenCredential, options *SettingsClientOptions) (*SettingsClient, error) {
	if options == nil {
		options = &SettingsClientOptions{}
	}
	authPolicy := internal.NewKeyVaultChallengePolicy(
		credential,
		&internal.KeyVaultChallengePolicyOptions{
			DisableChallengeResourceVerification: options.DisableChallengeResourceVerification,
		},
	)
	pl := runtime.NewPipeline(shared.ModuleName, shared.Version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	return &SettingsClient{endpoint: vaultURL, pl: pl}, nil
}

// Error - The code and message for an error.
type Error struct {
	// REQUIRED; A machine readable error code.
	Code string

	// full error message detailing why the operation failed.
	data []byte
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Error.
func (e *Error) UnmarshalJSON(data []byte) error {
	e.data = data
	ei := struct{ Code string }{}
	if err := json.Unmarshal(data, &ei); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	e.Code = ei.Code

	return nil
}

// Error implements a custom error for type Error.
func (e *Error) Error() string {
	return string(e.data)
}

// beginFullRestore is a custom implementation of BeginFullRestore
// Uses custom poller handler
func (client *BackupClient) beginFullRestore(ctx context.Context, restoreBlobDetails RestoreOperationParameters, options *BackupClientBeginFullRestoreOptions) (*runtime.Poller[BackupClientFullRestoreResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.fullRestore(ctx, restoreBlobDetails, options)
		if err != nil {
			return nil, err
		}
		handler, err := newRestorePoller[BackupClientFullRestoreResponse](client.pl, resp, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[BackupClientFullRestoreResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Handler:       handler,
		})
	} else {
		handler, err := newRestorePoller[BackupClientFullRestoreResponse](client.pl, nil, runtime.FinalStateViaAzureAsyncOp)
		if err != nil {
			return nil, err
		}
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.pl, &runtime.NewPollerFromResumeTokenOptions[BackupClientFullRestoreResponse]{Handler: handler})
	}
}
