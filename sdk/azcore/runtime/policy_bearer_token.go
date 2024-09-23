// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"encoding/base64"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"
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
		b.authzHandler.OnChallenge = b.handleCAEChallenge
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

func (b *BearerTokenPolicy) handleCAEChallenge(_ *policy.Request, res *http.Response, authNZ func(policy.TokenRequestOptions) error) error {
	claims, err := parseCAEChallenge(res)
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

var (
	claimsChallenge *regexp.Regexp
	once            = &sync.Once{}
)

// parseCAEChallenge extracts claims from the first CAE challenge in the response's WWW-Authenticate header. If
// the response doesn't contain a CAE challenge, it returns an empty string and nil error. It returns a non-nil
// error only when it identifies a CAE challenge whose claims aren't valid base64.
func parseCAEChallenge(res *http.Response) (string, error) {
	once.Do(func() {
		// This expression matches CAE claims challenges and captures their parameters. It doesn't
		// correspond precisely to the challenge grammar in RFC 7235 appendix C because CAE challenges
		// are more narrowly defined at
		// https://learn.microsoft.com/entra/identity-platform/claims-challenge#claims-challenge-header-format
		// It matches only Bearer challenges because this policy doesn't support other auth schemes.
		claimsChallenge = regexp.MustCompile(`(?:Bearer ((?:\w+="[^"]*",?\s*)+))`)
	})
	// WWW-Authenticate can have multiple values, each containing multiple challenges
	for _, h := range res.Header.Values(shared.HeaderWWWAuthenticate) {
		for _, sm := range claimsChallenge.FindAllStringSubmatch(h, -1) {
			if len(sm) < 2 {
				continue
			}
			for _, params := range sm[1:] {
				cae := false
				claims := ""
				for _, param := range strings.Split(params, ", ") {
					k, v, found := strings.Cut(param, `="`)
					if !found {
						// should never happen given the regex that got us here
						continue
					}
					switch strings.Trim(k, " ") {
					case "claims":
						claims = v
					case "error":
						// v has an irrelevant trailing quote
						if strings.HasPrefix(v, "insufficient_claims") {
							cae = true
						}
					}
				}
				if cae && claims != "" {
					// remove the trailing quote and any base64 padding
					if end := strings.IndexAny(claims, `="`); end > 0 {
						claims = claims[:end]
					}
					if b, de := base64.RawURLEncoding.DecodeString(claims); de == nil {
						return string(b), nil
					}
					// we don't include the decoding error because it's something
					// unhelpful like "illegal base64 data at input byte 42"
					return "", errors.New("challenge contains invalid claims: " + claims)
				}
			}
		}
	}
	return "", nil
}
