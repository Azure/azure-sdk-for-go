// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChainedTokenCredential provides a TokenCredential implementation which chains multiple TokenCredential implementations to be tried in order
// until one of the GetToken methods returns a non-default AccessToken.
type ChainedTokenCredential struct {
	sources []azcore.TokenCredential
}

// NewChainedTokenCredential creates an instance of ChainedTokenCredential with the specified TokenCredential sources.
func NewChainedTokenCredential(sources ...azcore.TokenCredential) (*ChainedTokenCredential, error) {
	if len(sources) == 0 {
		return nil, &CredentialUnavailableError{Message: "NewChainedTokenCredential: length of sources cannot be 0"}
	}
	// TODO: test for each of these conditions
	for _, source := range sources {
		if source == nil {
			return nil, &CredentialUnavailableError{Message: "NewChainedTokenCredential: sources cannot contain a nil TokenCredential"}
		}
	}

	return &ChainedTokenCredential{sources: sources}, nil
}

// GetToken sequentially calls TokenCredential.GetToken on all the specified sources, returning the first non default AccessToken.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	// CP: Notes from session with Jeff, have two funcs for the DefaultAzureCredential the main default func will issue the call that hits the wire, however internally first we create the chain and return a credUnavailable error if no credentials are found, if one works then we immediately return the one that works. End goal to always try to use the same credential type so we dont switch between credentials during execution, or at least reduce the possibility of that happening?
	var token *azcore.AccessToken
	var err error
	var errList []error

	for i := 0; i < len(c.sources) && token == nil; i++ {
		token, err = c.sources[i].GetToken(ctx, scopes)
		// CP: This currently works the same way as the other languages
		// TODO: instead of appending to an error slice, check if the error returned is an auth failure then stop the function completely and return the auth failure, credunavailable is still ok to continue looping
		if err != nil {
			errList = append(errList, err)
		}
	}

	if token == nil && len(errList) > 0 {
		// Here we dont want to return the slice of errors, instead ideally just return an auth failed error since the credential unavailable error should already have been discarded.
		err = &AggregateError{ErrList: errList}
		return nil, err
	}

	return token, nil
}

// AuthenticationPolicy implements the azcore.Credential interface on ChainedTokenCredential.
func (c *ChainedTokenCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options.Scopes)
}
