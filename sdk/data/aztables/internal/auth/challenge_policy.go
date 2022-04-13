//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package auth

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
)

const headerAuthorization = "Authorization"
const bearerHeader = "Bearer "

type TablesChallengePolicy struct {
	// mainResource is the resource to be retrieved using the tenant specified in the credential
	mainResource *ExpiringResource
	cred         azcore.TokenCredential
	scope        *string
	tenantID     *string
}

func NewTablesChallengePolicy(cred azcore.TokenCredential) *TablesChallengePolicy {
	return &TablesChallengePolicy{
		cred:         cred,
		mainResource: NewExpiringResource(acquire),
	}
}

func (scp *TablesChallengePolicy) Do(req *policy.Request) (*http.Response, error) {
	as := acquiringResourceState{
		p:   scp,
		req: req,
	}

	if scp.scope == nil || scp.tenantID == nil {
		// First request, get both to get the token
		challengeReq, err := scp.getChallengeRequest(*req)
		if err != nil {
			return nil, err
		}

		resp, err := challengeReq.Next()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode > 399 && resp.StatusCode != http.StatusUnauthorized {
			// the request failed for some other reason, don't try any further
			return resp, nil
		}
		err = scp.findScopeAndTenant(resp)
		if err != nil {
			return nil, err
		}
	}

	tk, err := scp.mainResource.GetResource(as)
	if err != nil {
		return nil, err
	}

	if token, ok := tk.(*azcore.AccessToken); ok {
		req.Raw().Header.Set(
			headerAuthorization,
			fmt.Sprintf("%s%s", bearerHeader, token.Token),
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
		scp.mainResource.Reset()

		// Find the scope and tenant again in case they have changed
		err := scp.findScopeAndTenant(resp)
		if err != nil {
			// Error parsing challenge, doomed to fail. Return
			return resp, cloneReqErr
		}

		tk, err := scp.mainResource.GetResource(as)
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
		return to.Ptr("")
	}
	parts := strings.Split(url, "/")
	tenant := parts[3]
	tenant = strings.ReplaceAll(tenant, ",", "")
	return &tenant
}

type challengePolicyError struct {
	err error
}

func (c *challengePolicyError) Error() string {
	return c.err.Error()
}

func (*challengePolicyError) NonRetriable() {
	// marker method
}

func (c *challengePolicyError) Unwrap() error {
	return c.err
}

var _ errorinfo.NonRetriable = (*challengePolicyError)(nil)

// sets the scp.scope and scp.tenantID from the WWW-Authenticate header
func (scp *TablesChallengePolicy) findScopeAndTenant(resp *http.Response) error {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return &challengePolicyError{err: errors.New("response has no WWW-Authenticate header for challenge authentication")}
	}

	// Strip down to auth and resource
	// Format is "Bearer authorization=<site> resource=<site>" OR
	// "Bearer authorization=<site> scope=<site> resource=<resource>"
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

	scp.tenantID = parseTenant(vals["authorization_uri"])
	if scope, ok := vals["scope"]; ok {
		scp.scope = &scope
	} else if resource, ok := vals["resource_id"]; ok {
		if !strings.HasSuffix(resource, "/.default") {
			resource += "/.default"
		}
		scp.scope = &resource
	} else {
		return &challengePolicyError{err: errors.New("could not find a valid resource in the WWW-Authenticate header")}
	}

	return nil
}

func (k TablesChallengePolicy) getChallengeRequest(orig policy.Request) (*policy.Request, error) {
	req, err := runtime.NewRequest(orig.Raw().Context(), "GET", orig.Raw().URL.String())
	if err != nil {
		return nil, &challengePolicyError{err: err}
	}

	req.Raw().Header = orig.Raw().Header
	req.Raw().Header.Set("Content-Length", "0")
	req.Raw().ContentLength = 0

	copied := orig.Clone(orig.Raw().Context())
	copied.Raw().Body = req.Body()
	copied.Raw().ContentLength = 0
	copied.Raw().Header.Set("Content-Length", "0")
	err = copied.SetBody(streaming.NopCloser(bytes.NewReader([]byte{})), "application/json")
	if err != nil {
		return nil, &challengePolicyError{err: err}
	}
	copied.Raw().Header.Del("Content-Type")

	return copied, err
}

type acquiringResourceState struct {
	req *policy.Request
	p   *TablesChallengePolicy
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
