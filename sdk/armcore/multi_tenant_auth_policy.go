// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	bearerTokenPrefix            = "Bearer "
	headerAuthorizationAuxiliary = "x-ms-authorization-auxiliary"
)

type multiTenantPolicy struct {
	// cond is used to synchronize token refresh.  the locker
	// must be locked when updating the following shared state.
	cond *sync.Cond

	// renewing indicates that the token is in the process of being refreshed
	renewing bool

	// header contains the authorization header value
	header string

	// expiresOn is when the token will expire
	expiresOn map[string]time.Time

	// the following fields are read-only
	creds            azcore.TokenCredential
	options          azcore.TokenRequestOptions
	auxiliaryTenants []string
}

func newMultiTenantPolicy(cred azcore.TokenCredential, opts authenticationPolicyOptions) azcore.Policy {
	return &multiTenantPolicy{
		cond:             sync.NewCond(&sync.Mutex{}),
		creds:            cred,
		options:          opts.Options,
		auxiliaryTenants: opts.AuxiliaryTenants,
	}
}

func (b *multiTenantPolicy) Do(req *azcore.Request) (*azcore.Response, error) {
	if req.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, errors.New("token credentials require a URL using the HTTPS protocol scheme")
	}
	// create a "refresh window" before the token's real expiration date.
	// this allows callers to continue to use the old token while the
	// refresh is in progress.
	const window = 2 * time.Minute
	now, getToken, header := time.Now(), false, ""
	auxTenants := []string{}
	// acquire exclusive lock
	b.cond.L.Lock()
	for {
		if b.expiresOn == nil || len(b.expiresOn) == 0 {
			// token was never obtained
			if !b.renewing {
				// another go routine isn't refreshing the token so this one will
				b.renewing = true
				getToken = true
				for _, tenant := range b.auxiliaryTenants {
					auxTenants = append(auxTenants, tenant)
				}
				break
			}
			// getting here means this go routine will wait for the token to refresh
		} else {
			// this go routine will use the existing token while another refreshes it
			// unless one of the auxiliary tokens needs to be updated
			header = b.header
			renew := false
			for id, eo := range b.expiresOn {
				if eo.Add(-window).Before(now) || eo.IsZero() {
					// token is within the expiration window or has expired
					// TODO check this logic
					if !b.renewing {
						// another go routine isn't refreshing the token so this one will
						renew = true
						getToken = true
						auxTenants = append(auxTenants, id)
					}
				}
			}
			if !b.renewing && renew {
				b.renewing = renew
			}
			break
		}
		// wait for the token to refresh
		b.cond.Wait()
	}
	b.cond.L.Unlock()
	if getToken {
		auxH := []string{}
		for _, tenant := range auxTenants {
			opts := b.options
			opts.TenantID = tenant
			tk, err := b.creds.GetToken(req.Context(), opts)
			if err != nil {
				// update shared state
				b.cond.L.Lock()
				// to avoid a deadlock if GetToken() fails we MUST reset b.renewing to false before returning
				b.renewing = false
				b.unlock()
				return nil, err
			}
			// update shared state
			b.cond.L.Lock()
			if b.expiresOn == nil {
				b.expiresOn = map[string]time.Time{}
			}
			b.expiresOn[tenant] = tk.ExpiresOn
			b.unlock()
			auxH = append(auxH, bearerTokenPrefix+tk.Token)
		}
		header = strings.Join(auxH, ", ")
		// update shared state
		b.cond.L.Lock()
		b.renewing = false
		b.header = header
		b.unlock()
	}
	req.Request.Header.Set(headerAuthorizationAuxiliary, header)
	return req.Next()
}

// signal any waiters that the token has been refreshed
func (b *multiTenantPolicy) unlock() {
	b.cond.Broadcast()
	b.cond.L.Unlock()
}

type authenticationPolicyOptions struct {
	Options          azcore.TokenRequestOptions
	AuxiliaryTenants []string
}
