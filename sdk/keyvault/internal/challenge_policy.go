//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

const headerAuthorization = "Authorization"
const bearerHeader = "Bearer "

type KeyVaultChallengePolicy struct {
	// mainResource is the resource to be retrieved using the tenant specified in the credential
	mainResource *ExpiringResource
	cred         azcore.TokenCredential
	scope        *string
	tenantID     *string
}

func NewKeyVaultChallengePolicy(cred azcore.TokenCredential) *KeyVaultChallengePolicy {
	return &KeyVaultChallengePolicy{
		cred:         cred,
		mainResource: NewExpiringResource(acquire),
	}
}

func (k *KeyVaultChallengePolicy) Do(req *policy.Request) (*http.Response, error) {
	as := acquiringResourceState{
		p:   k,
		req: req,
	}

	if k.scope == nil || k.tenantID == nil {
		// First request, get both to get the token
		challengeReq, err := k.getChallengeRequest(*req)
		if err != nil {
			return nil, err
		}

		challengeResp, err := challengeReq.Next()
		if err != nil {
			return nil, err
		}

		err = k.findScopeAndTenant(challengeResp)
		if err != nil {
			return nil, err
		}
	}

	tk, err := k.mainResource.GetResource(as)
	if err != nil {
		return nil, err
	}

	if token, ok := tk.(*azcore.AccessToken); ok {
		req.Raw().Header.Set(
			headerAuthorization,
			fmt.Sprintf("%s %s", bearerHeader, token.Token),
		)
	}

	// send a copy of the request
	cloneReq := req.Clone(req.Raw().Context())
	resp, cloneReqErr := cloneReq.Next()
	if cloneReqErr != nil {
		return nil, cloneReqErr
	}

	// If it fails and has a 401, try it with a new token
	if resp.StatusCode == 401 {
		// Force a new token
		k.mainResource.Reset()

		// Find the scope and tenant again in case they have changed
		err := k.findScopeAndTenant(resp)
		if err != nil {
			// Error parsing challenge, doomed to fail. Return
			return resp, err
		}

		tk, err := k.mainResource.GetResource(as)
		if err != nil {
			return resp, err
		}

		if token, ok := tk.(*azcore.AccessToken); ok {
			req.Raw().Header.Set(
				headerAuthorization,
				bearerHeader+token.Token,
			)
		} else {
			// tk is not an azcore.AccessToken type, something went wrong and we should return the 401 and accompanying error
			return resp, cloneReqErr
		}

		// send the original request now
		return req.Next()
	}

	return resp, err
}

// parses Tenant ID from auth challenge
// https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000
func parseTenant(url string) *string {
	if url == "" {
		return to.StringPtr("")
	}
	parts := strings.Split(url, "/")
	tenant := parts[3]
	tenant = strings.ReplaceAll(tenant, ",", "")
	return &tenant
}

// sets the k.scope and k.tenantID from the WWW-Authenticate header
func (k *KeyVaultChallengePolicy) findScopeAndTenant(resp *http.Response) error {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return errors.New("response has no WWW-Authenticate header for challenge authentication")
	}

	// Strip down to auth and resource
	// Format is "Bearer authorization=\"<site>\" resource=\"<site>\"" OR
	// "Bearer authorization=\"<site>\" scope=\"<site>\" resource=\"<resource>\""
	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")

	parts := strings.Split(authHeader, " ")

	vals := map[string]string{}
	for _, part := range parts {
		subParts := strings.Split(part, "=")
		if len(subParts) == 2 {
			stripped := strings.ReplaceAll(subParts[1], "\"", "")
			stripped = strings.TrimSuffix(stripped, ",")
			vals[subParts[0]] = stripped
		}
	}

	k.tenantID = parseTenant(vals["authorization"])
	if scope, ok := vals["scope"]; ok {
		k.scope = &scope
	} else if resource, ok := vals["resource"]; ok {
		if !strings.HasSuffix(resource, "/.default") {
			resource += "/.default"
		}
		k.scope = &resource
	} else {
		return errors.New("could not find a valid resource in the WWW-Authenticate header")
	}

	return nil
}

func (k KeyVaultChallengePolicy) getChallengeRequest(orig policy.Request) (*policy.Request, error) {
	req, err := runtime.NewRequest(orig.Raw().Context(), orig.Raw().Method, orig.Raw().URL.String())
	if err != nil {
		return nil, err
	}

	req.Raw().Header = orig.Raw().Header
	req.Raw().Header.Set("Content-Length", "0")

	copied := orig.Clone(orig.Raw().Context())
	copied.Raw().Body = req.Body()

	return copied, err
}

type acquiringResourceState struct {
	req *policy.Request
	p   *KeyVaultChallengePolicy
}

// acquire acquires or updates the resource; only one
// thread/goroutine at a time ever calls this function
func acquire(state interface{}) (newResource interface{}, newExpiration time.Time, err error) {
	s := state.(acquiringResourceState)
	tk, err := s.p.cred.GetToken(
		s.req.Raw().Context(),
		policy.TokenRequestOptions{
			Scopes:   []string{*s.p.scope},
			TenantID: *s.p.scope,
		},
	)
	if err != nil {
		return nil, time.Time{}, err
	}
	return tk, tk.ExpiresOn, nil
}
