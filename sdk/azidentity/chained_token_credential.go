// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

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
	cond                 *sync.Cond
	iterating            bool
	name                 string
	retrySources         bool
	sources              []azcore.TokenCredential
	successfulCredential azcore.TokenCredential
}

// NewChainedTokenCredential creates a ChainedTokenCredential.
// sources: Credential instances to comprise the chain. GetToken() will invoke them in the given order.
// options: Optional configuration. Pass nil to accept default settings.
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
	return &ChainedTokenCredential{
		cond:         sync.NewCond(&sync.Mutex{}),
		name:         "ChainedTokenCredential",
		retrySources: options.RetrySources,
		sources:      cp,
	}, nil
}

// GetToken calls GetToken on the chained credentials in turn, stopping when one returns a token. This method is called automatically by Azure SDK clients.
// ctx: Context controlling the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	if !c.retrySources {
		// ensure only one goroutine at a time iterates the sources and perhaps sets c.successfulCredential
		c.cond.L.Lock()
		for {
			if c.successfulCredential != nil {
				c.cond.L.Unlock()
				return c.successfulCredential.GetToken(ctx, opts)
			}
			if !c.iterating {
				c.iterating = true
				// allow other goroutines to wait while this one iterates
				c.cond.L.Unlock()
				break
			}
			c.cond.Wait()
		}
	}

	var err error
	var errs []error
	var token *azcore.AccessToken
	var successfulCredential azcore.TokenCredential
	for _, cred := range c.sources {
		token, err = cred.GetToken(ctx, opts)
		if err == nil {
			log.Writef(EventAuthentication, "%s authenticated with %s", c.name, extractCredentialName(cred))
			successfulCredential = cred
			break
		}
		errs = append(errs, err)
		if _, ok := err.(credentialUnavailableError); !ok {
			break
		}
	}
	if c.iterating {
		c.cond.L.Lock()
		c.successfulCredential = successfulCredential
		c.iterating = false
		c.cond.L.Unlock()
		c.cond.Broadcast()
	}
	// err is the error returned by the last GetToken call. It will be nil when that call succeeds
	if err != nil {
		// return credentialUnavailableError iff all sources did so; return AuthenticationFailedError otherwise
		msg := createChainedErrorMessage(errs)
		if _, ok := err.(credentialUnavailableError); ok {
			err = newCredentialUnavailableError(c.name, msg)
		} else {
			res := getResponseFromError(err)
			err = newAuthenticationFailedError(c.name, msg, res)
		}
	}
	return token, err
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
