// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	bearerTokenPrefix = "Bearer "
)

type bearerTokenPolicy struct {
	// cond is used to synchronize token refresh.  the locker
	// must be locked when updating the following shared state.
	cond *sync.Cond

	// renewing indicates that the token is in the process of being refreshed
	renewing bool

	// header contains the authorization header value
	header string

	// auxHeader contains the authorization header value
	auxHeader string

	// expiresOn is when the token will expire
	expiresOn time.Time

	// the following fields are read-only
	creds   azcore.TokenCredential
	options ARMAuthenticationPolicyOptions
}

func NewMultiTenantAuthenticationPolicy(creds azcore.TokenCredential, opts ARMAuthenticationPolicyOptions) azcore.Policy {
	return &bearerTokenPolicy{
		cond:    sync.NewCond(&sync.Mutex{}),
		creds:   creds,
		options: opts,
	}
}

func (b *bearerTokenPolicy) Do(req *azcore.Request) (*azcore.Response, error) {
	if req.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, &ARMAuthenticationFailedError{msg: "token credentials require a URL using the HTTPS protocol scheme"}
	}
	// create a "refresh window" before the token's real expiration date.
	// this allows callers to continue to use the old token while the
	// refresh is in progress.
	const window = 2 * time.Minute
	now, getToken, header := time.Now(), false, ""
	// acquire exclusive lock
	b.cond.L.Lock()
	for {
		if b.expiresOn.IsZero() || b.expiresOn.Before(now) {
			// token was never obtained or has expired
			if !b.renewing {
				// another go routine isn't refreshing the token so this one will
				b.renewing = true
				getToken = true
				break
			}
			// getting here means this go routine will wait for the token to refresh
		} else if b.expiresOn.Add(-window).Before(now) {
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
		// this go routine has been elected to refresh the token
		tk, err := b.creds.GetToken(req.Context(), b.options.Options)
		auxH := []string{}
		for _, id := range b.options.MultiTenantIDs {
			opts := b.options.Options
			opts.TenantIDHint = id
			tk, err = b.creds.GetToken(req.Context(), opts)
			if err != nil {
				return nil, err
			}
			auxH = append(auxH, bearerTokenPrefix+tk.Token)
		}
		// update shared state
		b.cond.L.Lock()
		// to avoid a deadlock if GetToken() fails we MUST reset b.renewing to false before returning
		b.renewing = false
		if err != nil {
			b.unlock()
			return nil, err
		}
		header = bearerTokenPrefix + tk.Token
		b.header = header
		b.expiresOn = tk.ExpiresOn
		auxHeader := strings.Join(auxH, ", ")
		b.auxHeader = auxHeader
		b.unlock()
	}
	req.Request.Header.Set(azcore.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
	req.Request.Header.Set(azcore.HeaderAuthorization, header)
	req.Request.Header.Set("x-ms-authorization-auxiliary", b.auxHeader)
	return req.Next()
}

// signal any waiters that the token has been refreshed
func (b *bearerTokenPolicy) unlock() {
	b.cond.Broadcast()
	b.cond.L.Unlock()
}

// AuthenticationFailedError is returned when the authentication request has failed.
type ARMAuthenticationFailedError struct {
	inner error
	msg   string
}

// Unwrap method on AuthenticationFailedError provides access to the inner error if available.
func (e *ARMAuthenticationFailedError) Unwrap() error {
	return e.inner
}

// NonRetriable indicates that this error should not be retried.
func (e *ARMAuthenticationFailedError) NonRetriable() {
	// marker method
}

func (e *ARMAuthenticationFailedError) Error() string {
	if e.inner == nil {
		return e.msg
	} else if e.msg == "" {
		return e.inner.Error()
	}
	return e.msg + " details: " + e.inner.Error()
}
