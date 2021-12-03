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
	// If true, it will not assume the first successful credential should be always used.
	RetryAllSources bool
}

// ChainedTokenCredential is a chain of credentials that enables fallback behavior when a credential can't authenticate.
type ChainedTokenCredential struct {
	sources              []azcore.TokenCredential
	successfulCredential azcore.TokenCredential
	retryAllSources      bool
}

// NewChainedTokenCredential creates a ChainedTokenCredential.
// By default, this credential will assume that the first successful credential should be the only credential used on future requests.
// If the `RetryAllSources` option is set to true, it will always try to get a token using all of the originally provided credentials.
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
	credentialOptions := ChainedTokenCredentialOptions{}
	if options != nil {
		credentialOptions = *options
	}
	return &ChainedTokenCredential{sources: cp, retryAllSources: credentialOptions.RetryAllSources}, nil
}

// GetToken calls GetToken on the chained credentials in turn, stopping when one returns a token. This method is called automatically by Azure SDK clients.
// ctx: Context controlling the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	var errList []CredentialUnavailableError

	formatError := func(err error) error {
		var authFailed AuthenticationFailedError
		if errors.As(err, &authFailed) {
			err = fmt.Errorf("Authentication failed:\n%s\n%s"+createChainedErrorMessage(errList), err)
			authErr := newAuthenticationFailedError(err, authFailed.RawResponse())
			return authErr
		}
		return err
	}

	if c.successfulCredential != nil && !c.retryAllSources {
		token, err = c.successfulCredential.GetToken(ctx, opts)
		if err != nil {
			return nil, formatError(err)
		}
		return token, nil
	}
	for _, cred := range c.sources {
		token, err = cred.GetToken(ctx, opts)
		var credErr CredentialUnavailableError
		if errors.As(err, &credErr) {
			errList = append(errList, credErr)
		} else if err != nil {
			return nil, formatError(err)
		} else {
			logGetTokenSuccess(c, opts)
			c.successfulCredential = cred
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
