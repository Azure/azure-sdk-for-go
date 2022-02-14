// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// ChainedTokenCredentialOptions contains optional parameters for ChainedTokenCredential.
type ChainedTokenCredentialOptions struct {
	// RetrySources configures how the credential uses its sources.
	// When true, the credential will always request a token from each source in turn,
	// stopping when one provides a token. When false, the credential requests a token
	// only from the source that previously retrieved a token--it never again tries the sources which failed.
	RetrySources bool
}

// ChainedTokenCredential is a chain of credentials that enables fallback behavior when a credential can't authenticate.
// By default, this credential will assume that the first successful credential should be the only credential used on future requests.
// If the `RetrySources` option is set to true, it will always try to get a token using all of the originally provided credentials.
type ChainedTokenCredential struct {
	sources              []azcore.TokenCredential
	successfulCredential azcore.TokenCredential
	retrySources         bool
	name                 string
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
	if options == nil {
		options = &ChainedTokenCredentialOptions{}
	}
	return &ChainedTokenCredential{sources: cp, name: "ChainedTokenCredential", retrySources: options.RetrySources}, nil
}

// GetToken calls GetToken on the chained credentials in turn, stopping when one returns a token. This method is called automatically by Azure SDK clients.
// ctx: Context controlling the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	if c.successfulCredential != nil && !c.retrySources {
		return c.successfulCredential.GetToken(ctx, opts)
	}

	var errs []error
	for _, cred := range c.sources {
		token, err := cred.GetToken(ctx, opts)
		if err == nil {
			log.Writef(EventAuthentication, "Azure Identity => %s authenticated with %s", c.name, extractCredentialName(cred))
			c.successfulCredential = cred
			return token, nil
		}
		errs = append(errs, err)
		if _, ok := err.(credentialUnavailableError); !ok {
			res := getResponseFromError(err)
			msg := createChainedErrorMessage(errs)
			return nil, newAuthenticationFailedError(c.name, msg, res)
		}
	}
	// if we get here, all credentials returned credentialUnavailableError
	msg := createChainedErrorMessage(errs)
	err := newCredentialUnavailableError(c.name, msg)
	log.Write(EventAuthentication, "Azure Identity => ERROR: "+err.Error())
	return nil, err
}

func createChainedErrorMessage(errs []error) string {
	msg := "failed to acquire a token.\nAttempted credentials:"
	for _, err := range errs {
		msg += fmt.Sprintf("\n\t%s", err.Error())
	}
	return msg
}

func extractCredentialName(credential azcore.TokenCredential) string {
	return strings.TrimPrefix(fmt.Sprintf("%T", credential), "*azidentity.")
}

var _ azcore.TokenCredential = (*ChainedTokenCredential)(nil)
