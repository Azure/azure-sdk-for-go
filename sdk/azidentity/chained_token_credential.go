// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
)

// ChainedTokenCredential provides a TokenCredential implementation which chains multiple TokenCredential implementations to be tried in order
// until one of the GetToken methods returns a non-default AccessToken.
type ChainedTokenCredential struct {
	TokenCredential
	sources []TokenCredential
}

// NewChainedTokenCredential creates an instance with the specified TokenCredential sources.
func NewChainedTokenCredential(sources ...TokenCredential) (ChainedTokenCredential, error) {
	if len(sources) == 0 {
		return ChainedTokenCredential{sources: nil}, errors.New("NewChainedTokenCredential: length of sources cannot be 0")
	}

	for _, source := range sources {
		if source == nil {
			return ChainedTokenCredential{sources: nil}, errors.New("NewChainedTokenCredential: sources cannot contain a nil TokenCredential")
		}
	}

	return ChainedTokenCredential{sources: sources}, nil
}

// GetToken sequentially calls TokenCredential.GetToken on all the specified sources, returning the first non default AccessToken.
// If all credentials in the chain return default, a default AccessToken is returned.
func (c ChainedTokenCredential) GetToken(ctx context.Context, scopes []string) (*AccessToken, error) {
	var token *AccessToken
	var err error
	for i := 0; i < len(c.sources) && token == nil; i++ {
		token, err = c.sources[i].GetToken(ctx, scopes)
		if err != nil {
			return nil, fmt.Errorf("ChainedTokenCredential GetToken(): %w", err)
		}
	}
	return token, nil
}

// NewDefaultTokenCredential provides a default ChainedTokenCredential configuration for applications that will be deployed to Azure.  The following credential
// types will be tried, in order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// Consult the documentation of these credential types for more information on how they attempt authentication.
func NewDefaultTokenCredential(o *IdentityClientOptions) (ChainedTokenCredential, error) {
	// CP: This is fine because we are not calling GetToken we are simple creating the new EnvironmentClient
	envClient, err := NewEnvironmentCredential(o)
	if err != nil {
		return ChainedTokenCredential{sources: nil}, fmt.Errorf("NewDefaultTokenCredential: %w", err)
	}
	// TODO: check this implementation:
	// 1. params for constructor should be nilable
	// 2. Should this func ask for a client id? or get it from somewhere else?
	msiClient, err := NewManagedIdentityCredential("", o)
	if err != nil {
		return ChainedTokenCredential{sources: nil}, fmt.Errorf("NewDefaultTokenCredential: %w", err)
	}

	return NewChainedTokenCredential(
		envClient,
		msiClient,
		credentialNotFoundGuard{})
}

type credentialNotFoundGuard struct {
	TokenCredential
}

func (c credentialNotFoundGuard) GetToken(ctx context.Context, scopes []string) (*AccessToken, error) {
	return &AccessToken{}, errors.New("Failed to find a credential to use for authentication.  If running in an environment where a managed identity is not available ensure the environment variables AZURE_TENANT_ID, AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET are set")
}
