// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

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

	for _, source := range sources {
		if source == nil {
			return nil, &CredentialUnavailableError{Message: "NewChainedTokenCredential: sources cannot contain a nil TokenCredential"}
		}
	}

	return &ChainedTokenCredential{sources: sources}, nil
}

// GetToken sequentially calls TokenCredential.GetToken on all the specified sources, returning the first non default AccessToken.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts *azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	var token *azcore.AccessToken
	var err error
	var errList []error

	for i := 0; i < len(c.sources); i++ {
		token, err = c.sources[i].GetToken(ctx, opts)
		// TODO: check if the err is an auth failure and stop loop, if cred unavailable you can continue
		if errors.Is(err, &CredentialUnavailableError{}) {
			errList = append(errList, err)
		} else { // TODO fix this
			return token, &AuthenticationFailedError{Message: err.Error()}
		}
	}

	if token == nil && len(errList) > 0 {
		// TODO err message should include cred name and failure reason
		// Pass back slice of errors here
		err = &CredentialUnavailableError{ErrList: errList}
		return nil, err
	}

	return token, nil
}

// AuthenticationPolicy implements the azcore.Credential interface on ChainedTokenCredential.
func (c *ChainedTokenCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options.Scopes)
}
