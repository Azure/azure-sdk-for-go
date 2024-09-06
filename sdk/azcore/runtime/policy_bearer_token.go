// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

// BearerTokenPolicy authorizes requests with bearer tokens acquired from a TokenCredential.
type BearerTokenPolicy struct {
	// mainResource is the resource to be retreived using the tenant specified in the credential
	mainResource *temporal.Resource[exported.AccessToken, acquiringResourceState]
	// the following fields are read-only
	authzHandler policy.AuthorizationHandler
	cred         exported.TokenCredential
	scopes       []string
	allowHTTP    bool
}

type acquiringResourceState struct {
	req *policy.Request
	p   *BearerTokenPolicy
	tro policy.TokenRequestOptions
}

// acquire acquires or updates the resource; only one
// thread/goroutine at a time ever calls this function
func acquire(state acquiringResourceState) (newResource exported.AccessToken, newExpiration time.Time, err error) {
	tk, err := state.p.cred.GetToken(&shared.ContextWithDeniedValues{Context: state.req.Raw().Context()}, state.tro)
	if err != nil {
		return exported.AccessToken{}, time.Time{}, err
	}
	return tk, tk.ExpiresOn, nil
}

// NewBearerTokenPolicy creates a policy object that authorizes requests with bearer tokens.
// cred: an azcore.TokenCredential implementation such as a credential object from azidentity
// scopes: the list of permission scopes required for the token.
// opts: optional settings. Pass nil to accept default values; this is the same as passing a zero-value options.
func NewBearerTokenPolicy(cred exported.TokenCredential, scopes []string, opts *policy.BearerTokenOptions) *BearerTokenPolicy {
	if opts == nil {
		opts = &policy.BearerTokenOptions{}
	}
	b := &BearerTokenPolicy{
		authzHandler: opts.AuthorizationHandler,
		cred:         cred,
		scopes:       scopes,
		mainResource: temporal.NewResource(acquire),
		allowHTTP:    opts.InsecureAllowCredentialWithHTTP,
	}
	if b.authzHandler.OnChallenge == nil {
		b.authzHandler.OnChallenge = b.handleClaimsChallenge
		b.authzHandler.SupportsCAE = true
	}
	return b
}

// authenticateAndAuthorize returns a function which authorizes req with a token from the policy's credential
func (b *BearerTokenPolicy) authenticateAndAuthorize(req *policy.Request) func(policy.TokenRequestOptions) error {
	return func(tro policy.TokenRequestOptions) error {
		tro.EnableCAE = b.authzHandler.SupportsCAE
		as := acquiringResourceState{p: b, req: req, tro: tro}
		tk, err := b.mainResource.Get(as)
		if err != nil {
			return err
		}
		req.Raw().Header.Set(shared.HeaderAuthorization, shared.BearerTokenPrefix+tk.Token)
		return nil
	}
}

// Do authorizes a request with a bearer token
func (b *BearerTokenPolicy) Do(req *policy.Request) (*http.Response, error) {
	// skip adding the authorization header if no TokenCredential was provided.
	// this prevents a panic that might be hard to diagnose and allows testing
	// against http endpoints that don't require authentication.
	if b.cred == nil {
		return req.Next()
	}

	if err := checkHTTPSForAuth(req, b.allowHTTP); err != nil {
		return nil, err
	}

	var err error
	if b.authzHandler.OnRequest != nil {
		err = b.authzHandler.OnRequest(req, b.authenticateAndAuthorize(req))
	} else {
		err = b.authenticateAndAuthorize(req)(policy.TokenRequestOptions{Scopes: b.scopes})
	}
	if err != nil {
		return nil, errorinfo.NonRetriableError(err)
	}

	res, err := req.Next()
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		b.mainResource.Expire()
		if res.Header.Get(shared.HeaderWWWAuthenticate) != "" && b.authzHandler.OnChallenge != nil {
			if err = b.authzHandler.OnChallenge(req, res, b.authenticateAndAuthorize(req)); err == nil {
				res, err = req.Next()
			}
		}
	}
	if err != nil {
		err = errorinfo.NonRetriableError(err)
	}
	return res, err
}

func (b *BearerTokenPolicy) handleClaimsChallenge(_ *policy.Request, res *http.Response, authNZ func(policy.TokenRequestOptions) error) error {
	challenge := res.Header.Get(shared.HeaderWWWAuthenticate)
	claims, err := parseChallenge(challenge)
	if err != nil {
		// challenge contains claims we can't parse
		return err
	}
	if claims == "" {
		// no claims in challenge, so this is a simple authorization failure
		return NewResponseError(res)
	}
	// request a new token having the specified claims, send the request again
	return authNZ(policy.TokenRequestOptions{Claims: claims, EnableCAE: true, Scopes: b.scopes})
}

func checkHTTPSForAuth(req *policy.Request, allowHTTP bool) error {
	if strings.ToLower(req.Raw().URL.Scheme) != "https" && !allowHTTP {
		return errorinfo.NonRetriableError(errors.New("authenticated requests are not permitted for non TLS protected (https) endpoints"))
	}
	return nil
}

// parseChallenge parses claims from an authentication challenge so a client can request a token that
// will satisfy conditional access policies. It returns a non-nil error when the given value contains
// claims it can't parse. If the value contains no claims, it returns an empty string and a nil error.
func parseChallenge(wwwAuthenticate string) (string, error) {
	var (
		claims string
		err    error
	)
	for _, param := range strings.Split(wwwAuthenticate, ",") {
		if _, after, found := strings.Cut(param, "claims="); found {
			if claims != "" {
				// The header contains multiple challenges, at least two of which specify claims. The specs allow this
				// but it's unclear what a client should do in this case and there's as yet no concrete example of it.
				err = errors.New("found multiple claims parameters in challenge: " + wwwAuthenticate)
				break
			}
			// trim stuff that would get an error from RawURLEncoding; claims may or may not be padded
			claims = strings.Trim(after, `\"=`)
			// we don't return this error because it's something unhelpful like "illegal base64 data at input byte 42"
			if b, de := base64.RawURLEncoding.DecodeString(claims); de == nil {
				claims = string(b)
			} else {
				err = errors.New("challenge contains invalid claims: " + claims)
				break
			}
		}
	}
	return claims, err
}
