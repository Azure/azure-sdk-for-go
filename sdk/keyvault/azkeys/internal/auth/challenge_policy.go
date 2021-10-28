//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type KeyVaultChallengePolicy struct {
	Cred        azcore.TokenCredential
	Transport   policy.Transporter
	cachedToken *azcore.AccessToken
	scope       *string
	tenantID    *string
}

func (k *KeyVaultChallengePolicy) Do(req *policy.Request) (*http.Response, error) {
	if k.scope == nil || k.tenantID == nil {
		// First request, get both to get the token
		challengeReq, err := getChallengeRequest(req)
		if err != nil {
			return nil, err
		}

		challengeResp, err := k.Transport.Do(challengeReq)
		if err != nil {
			return nil, err
		}

		resource, scope, err := parseAuthChallenge(challengeResp)
		if err != nil {
			return nil, err
		}
		k.tenantID = parseTenant(resource)

		if strings.HasSuffix(scope, "/") {
			scope += ".default"
		} else {
			scope += "/.default"
		}
		k.scope = &scope
	}

	token, err := k.Cred.GetToken(req.Raw().Context(), policy.TokenRequestOptions{
		Scopes:   []string{*k.scope},
		TenantID: *k.tenantID,
	})
	if err != nil {
		return nil, err
	}
	k.cachedToken = token
	k.decorateRequest(req)

	resp, err := req.Next()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		// Check for a new auth policy
		newTenant, newScope, err := parseAuthChallenge(resp)

		// Error parsing challenge, doomed to fail. Return
		if err != nil {
			return nil, err
		}

		k.tenantID = &newTenant
		k.scope = &newScope

		token, err := k.Cred.GetToken(
			req.Raw().Context(),
			policy.TokenRequestOptions{
				Scopes:   []string{*k.scope},
				TenantID: *k.tenantID,
			},
		)
		if err != nil {
			return nil, err
		}

		k.cachedToken = token
		k.decorateRequest(req)

		resp, err = http.DefaultClient.Do(req.Raw())
		if err != nil {
			// A second request failed, return error
			return nil, err
		}
	}

	return resp, err
}

// parses Tenant ID from auth challenge
// https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000
func parseTenant(url string) *string {
	parts := strings.Split(url, "/")
	tenant := parts[len(parts)-1]
	return &tenant
}

func (k KeyVaultChallengePolicy) decorateRequest(req *policy.Request) {
	req.Raw().Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.cachedToken.Token))
}

func getChallengeRequest(orig *policy.Request) (*http.Request, error) {
	req, err := http.NewRequest(orig.Raw().Method, orig.Raw().URL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = orig.Raw().Header.Clone()
	req.Header.Set("Content-Length", "0")

	return req, err
}

func parseAuthChallenge(resp *http.Response) (string, string, error) {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return "", "", errors.New("response has no WWW-Authenticate header for challenge authentication")
	}

	// Strip down to auth and resource
	// Format is "Bearer authorization=\"<site>\" resource=\"<site>\""
	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("could not understand WWW-Authenticate header: %s", authHeader)
	}

	vals := make([]string, 0)
	for _, part := range parts {
		subParts := strings.Split(part, "=")
		if len(subParts) != 2 {
			return "", "", fmt.Errorf("could not understand WWW-Authenticate header: %s", authHeader)
		}
		url := subParts[1]

		url = strings.ReplaceAll(url, "\"", "")
		vals = append(vals, url)
	}

	return vals[0], vals[1], nil
}
