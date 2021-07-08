// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	bearerTokenPrefix = "Bearer "
)

// DO NOT USE. This is for internal SDK usage. The TokenFetcher is used to provide
// custom implementations for the bearer token policy.
type TokenFetcher interface {
	// IsZeroOrExpired returns a bool if there is a tenant that needs to be initialized
	IsZeroOrExpired() bool
	// ShouldRefresh returns a bool if a token needs to be refreshed
	ShouldRefresh() bool
	// Fetch performs the GetToken call for the credential and returns the header needed for authentication
	Fetch(ctx context.Context, cred TokenCredential, opts TokenRequestOptions) (string, error)
	// Header returns the key for the header in the request
	Header() string
}

type tokenRefreshPolicy struct {
	// cond is used to synchronize token refresh.  the locker
	// must be locked when updating the following shared state.
	cond *sync.Cond

	// renewing indicates that the token is in the process of being refreshed
	renewing bool

	// header contains the authorization header value
	header string

	// the following fields are read-only
	creds       TokenCredential
	options     TokenRequestOptions
	implementer TokenFetcher
}

// NewTokenRefreshPolicy instantiates a thread safe token refresh policy.
// Pass in the credential to use, along with a TokenFetcher that will be
// used to provide a custom implementation for checking the expires_on
// time of the token, when tokens should be refreshed, and setting custom headers.
// Additionally specify options for the token to be retreived.
func NewTokenRefreshPolicy(cred TokenCredential, p TokenFetcher, opts AuthenticationPolicyOptions) Policy {
	return &tokenRefreshPolicy{
		cond:        sync.NewCond(&sync.Mutex{}),
		creds:       cred,
		options:     opts.Options,
		implementer: p,
	}
}

func (b *tokenRefreshPolicy) Do(req *Request) (*Response, error) {
	if req.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, errors.New("token credentials require a URL using the HTTPS protocol scheme")
	}
	getToken, header := false, ""
	// acquire exclusive lock
	b.cond.L.Lock()
	for {
		if b.implementer.IsZeroOrExpired() {
			// token was never obtained or has expired
			if !b.renewing {
				// another go routine isn't refreshing the token so this one will
				b.renewing = true
				getToken = true
				break
			}
			// getting here means this go routine will wait for the token to refresh
		} else if b.implementer.ShouldRefresh() {
			// token is within the expiration window
			if !b.renewing {
				// another go routine isn't refreshing the token so this one will
				b.renewing = true
				getToken = true
				break
			}
			// this go routine will use the existing token while another refreshes it
			header = b.header
			break
		} else {
			// token is not expiring yet so use it as-is
			header = b.header
			break
		}
		// wait for the token to refresh
		b.cond.Wait()
	}
	b.cond.L.Unlock()
	if getToken {
		var err error
		// this go routine has been elected to refresh the token
		header, err = b.implementer.Fetch(req.Context(), b.creds, b.options)
		// update shared state
		b.cond.L.Lock()
		// to avoid a deadlock if GetToken() fails we MUST reset b.renewing to false before returning
		b.renewing = false
		if err != nil {
			b.unlock()
			return nil, err
		}
		b.header = header
		b.unlock()
	}
	req.Request.Header.Set(HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
	req.Request.Header.Set(b.implementer.Header(), header)
	return req.Next()
}

// signal any waiters that the token has been refreshed
func (b *tokenRefreshPolicy) unlock() {
	b.cond.Broadcast()
	b.cond.L.Unlock()
}
