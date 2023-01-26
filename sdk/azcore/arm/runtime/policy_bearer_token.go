// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

const headerAuxiliaryAuthorization = "x-ms-authorization-auxiliary"

// acquiringResourceState holds data for an auxiliary token request
type acquiringResourceState struct {
	ctx    context.Context
	p      *BearerTokenPolicy
	tenant string
}

// acquire acquires or updates the resource; only one
// thread/goroutine at a time ever calls this function
func acquire(state acquiringResourceState) (newResource azcore.AccessToken, newExpiration time.Time, err error) {
	tk, err := state.p.cred.GetToken(state.ctx, azpolicy.TokenRequestOptions{
		Scopes:   state.p.scopes,
		TenantID: state.tenant,
	})
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}
	return tk, tk.ExpiresOn, nil
}

// BearerTokenPolicy authorizes requests with bearer tokens acquired from a TokenCredential.
type BearerTokenPolicy struct {
	auxResources map[string]*temporal.Resource[azcore.AccessToken, acquiringResourceState]
	btp          *azruntime.BearerTokenPolicy
	cred         azcore.TokenCredential
	scopes       []string
}

// NewBearerTokenPolicy creates a policy object that authorizes requests with bearer tokens.
// cred: an azcore.TokenCredential implementation such as a credential object from azidentity
// opts: optional settings. Pass nil to accept default values; this is the same as passing a zero-value options.
func NewBearerTokenPolicy(cred azcore.TokenCredential, opts *armpolicy.BearerTokenOptions) *BearerTokenPolicy {
	if opts == nil {
		opts = &armpolicy.BearerTokenOptions{}
	}
	p := &BearerTokenPolicy{cred: cred}
	p.auxResources = make(map[string]*temporal.Resource[azcore.AccessToken, acquiringResourceState], len(opts.AuxiliaryTenants))
	for _, t := range opts.AuxiliaryTenants {
		p.auxResources[t] = temporal.NewResource(acquire)
	}
	p.scopes = make([]string, len(opts.Scopes))
	copy(p.scopes, opts.Scopes)
	p.btp = azruntime.NewBearerTokenPolicy(cred, opts.Scopes, &azpolicy.BearerTokenOptions{
		AuthorizationHandler: azpolicy.AuthorizationHandler{
			OnChallenge: p.onChallenge,
			OnRequest:   p.onRequest,
		},
	})
	return p
}

func (b *BearerTokenPolicy) onChallenge(req *azpolicy.Request, res *http.Response, authNZ func(azpolicy.TokenRequestOptions) error) error {
	challenge := res.Header.Get(shared.HeaderWWWAuthenticate)
	if claims := parseChallenge(challenge); claims != "" {
		return authNZ(azpolicy.TokenRequestOptions{Claims: claims, Scopes: b.scopes})
	}
	return fmt.Errorf("failed to parse authentication challenge %q", challenge)
}

// onRequest authorizes requests with one or more bearer tokens
func (b *BearerTokenPolicy) onRequest(req *azpolicy.Request, authNZ func(azpolicy.TokenRequestOptions) error) error {
	// authorize the request with a token for the primary tenant
	err := authNZ(azpolicy.TokenRequestOptions{Scopes: b.scopes})
	if err != nil || len(b.auxResources) == 0 {
		return err
	}
	// add tokens for auxiliary tenants
	as := acquiringResourceState{
		ctx: req.Raw().Context(),
		p:   b,
	}
	auxTokens := make([]string, 0, len(b.auxResources))
	for tenant, er := range b.auxResources {
		as.tenant = tenant
		auxTk, err := er.Get(as)
		if err != nil {
			return err
		}
		auxTokens = append(auxTokens, fmt.Sprintf("%s%s", shared.BearerTokenPrefix, auxTk.Token))
	}
	req.Raw().Header.Set(headerAuxiliaryAuthorization, strings.Join(auxTokens, ", "))
	return nil
}

// Do authorizes a request with a bearer token
func (b *BearerTokenPolicy) Do(req *azpolicy.Request) (*http.Response, error) {
	return b.btp.Do(req)
}

// parseChallenge parses claims from an authentication challenge so a client can request a token that will satisfy
// conditional access policies. Returns an empty string when no claims are found or claims aren't in ARM's format.
// This function isn't universally applicable because RPs differ in challenge format and content.
func parseChallenge(wwwAuthenticate string) string {
	claims := ""
	for _, param := range strings.Split(wwwAuthenticate, ",") {
		if _, after, found := strings.Cut(param, "claims="); found {
			if claims != "" {
				// The header contains multiple challenges, at least two of which specify claims. The specs allow this
				// but it's unclear what a client should do in this case and there's as yet no concrete example of it.
				return ""
			}
			// trim stuff that would get an error from RawURLEncoding; claims may or may not be padded
			claims = strings.Trim(after, `\"=`)
			// we don't return the error because when not nil it's something unhelpful like "illegal base64 data at input byte 42"
			if b, err := base64.RawURLEncoding.DecodeString(claims); err == nil {
				claims = string(b)
			}
		}
	}
	return claims
}
