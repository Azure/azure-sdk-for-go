// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// ChainedTokenCredentialOptions contains optional parameters for ChainedTokenCredential.
type ChainedTokenCredentialOptions struct {
	// placeholder for future options
}

// ChainedTokenCredential is a chain of credentials that enables fallback behavior when a credential can't authenticate.
type ChainedTokenCredential struct {
	sources []azcore.TokenCredential
}

// NewChainedTokenCredential creates a ChainedTokenCredential.
// sources: Credential instances to comprise the chain. GetToken() will invoke them in the given order.
// options: Optional configuration.
func NewChainedTokenCredential(sources []azcore.TokenCredential, options *ChainedTokenCredentialOptions) (*ChainedTokenCredential, error) {
	if len(sources) == 0 {
		return nil, errors.New("sources must contain at least one TokenCredential")
	}
	for _, source := range sources {
		if source == nil { // cannot have a nil credential in the chain or else the application will panic when GetToken() is called on nil
			return nil, errors.New("sources cannot contain nil")
		}
	}
	cp := make([]azcore.TokenCredential, len(sources))
	copy(cp, sources)
	return &ChainedTokenCredential{sources: cp}, nil
}

// GetToken calls GetToken on the chained credentials in turn, stopping when one returns a token. This method is called automatically by Azure SDK clients.
// ctx: Context controlling the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	var errList []CredentialUnavailableError
	for _, cred := range c.sources {
		token, err = cred.GetToken(ctx, opts)
		var credErr CredentialUnavailableError
		if errors.As(err, &credErr) {
			errList = append(errList, credErr)
		} else if err != nil {
			var authFailed AuthenticationFailedError
			if errors.As(err, &authFailed) {
				err = fmt.Errorf("Authentication failed:\n%s\n%s"+createChainedErrorMessage(errList), err)
				authErr := newAuthenticationFailedError(err, authFailed.RawResponse())
				return nil, authErr
			}
			return nil, err
		} else {
			logGetTokenSuccess(c, opts)
			return token, nil
		}
	}
	// if we reach this point it means that all of the credentials in the chain returned CredentialUnavailableError
	credErr := newCredentialUnavailableError("Chained Token Credential", createChainedErrorMessage(errList))
	// skip adding the stack trace here as it was already logged by other calls to GetToken()
	addGetTokenFailureLogs("Chained Token Credential", credErr, false)
	return nil, credErr
}

func createChainedErrorMessage(errList []CredentialUnavailableError) string {
	msg := ""
	for _, err := range errList {
		msg += err.Error()
	}

	return msg
}

var _ azcore.TokenCredential = (*ChainedTokenCredential)(nil)
