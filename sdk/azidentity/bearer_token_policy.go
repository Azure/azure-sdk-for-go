// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/atomic"
)

type bearerTokenPolicy struct {
	empty     chan bool    // indicates we don't have a token yet
	updating  atomic.Int64 // atomically set to 1 to indicate a refresh is in progress
	header    atomic.String
	expiresOn atomic.Time
	creds     azcore.TokenCredential     // R/O
	options   azcore.TokenRequestOptions // R/O
}

func newBearerTokenPolicy(creds azcore.TokenCredential, opts azcore.AuthenticationPolicyOptions) *bearerTokenPolicy {
	bt := bearerTokenPolicy{creds: creds, options: opts.Options}
	// prime channel indicating there is no token
	bt.empty = make(chan bool, 1)
	bt.empty <- true
	return &bt
}

func (b *bearerTokenPolicy) Do(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
	if req.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, &AuthenticationFailedError{inner: errors.New("token credentials require a URL using the HTTPS protocol scheme")}
	}
	// check if the token is empty.  this will return true for
	// the first read.  all other reads will block until the
	// token has been obtained at which point it returns false
	if <-b.empty {
		tk, err := b.creds.GetToken(ctx, azcore.TokenRequestOptions{Scopes: b.options.Scopes})
		if err != nil {
			// failed to get a token, let another go routine try
			b.empty <- true
			return nil, err
		}
		b.header.Store("Bearer " + tk.Token)
		b.expiresOn.Store(tk.ExpiresOn)
		// signal token has been initialized.  go routines
		// that read from b.empty will always get false.
		close(b.empty)
	}
	// create a "refresh window" before the token's real expiration date.
	// this allows callers to continue to use the old token while the
	// refresh is in progress.
	const window = 2 * time.Minute
	// check if the token has expired
	now := time.Now().UTC()
	exp := b.expiresOn.Load()
	if now.After(exp.Add(-window)) {
		// token has expired, set update flag and begin the refresh.
		// if no other go routine has initiated a refresh the calling
		// go routine will do it.
		if b.updating.CAS(0, 1) {
			tk, err := b.creds.GetToken(ctx, azcore.TokenRequestOptions{Scopes: b.options.Scopes})
			if err != nil {
				// clear updating flag before returning so other
				// go routines can attempt to refresh
				b.updating.Store(0)
				return nil, err
			}
			b.header.Store("Bearer " + tk.Token)
			// set expiresOn last since refresh is predicated on it
			b.expiresOn.Store(tk.ExpiresOn)
			// signal the refresh is complete. this must happen after all shared
			// state has been updated.
			b.updating.Store(0)
		} // else { another go routine is refreshing the token, use previous token }
	}
	req.Request.Header.Set(azcore.HeaderAuthorization, b.header.Load())
	return req.Do(ctx)
}
