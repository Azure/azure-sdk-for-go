//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

const (
	headerAuthorization = "Authorization"
	bearerHeader        = "Bearer "
)

type authenticationPolicyOptions struct {
}

type authenticationPolicy struct {
	mainResource *temporal.Resource[azcore.AccessToken, acquiringResourceState]
	cred         azcore.TokenCredential
	aadScopes    []string
	acrScope     string
	acrService   string
	authClient   *authenticationClient
}

func newAuthenticationPolicy(cred azcore.TokenCredential, scopes []string, authClient *authenticationClient, opts *authenticationPolicyOptions) *authenticationPolicy {
	return &authenticationPolicy{
		cred:         cred,
		aadScopes:    scopes,
		authClient:   authClient,
		mainResource: temporal.NewResource(acquire),
	}
}

func (p *authenticationPolicy) Do(req *policy.Request) (*http.Response, error) {
	// send a copy of the original request without body content
	challengeReq, err := p.getChallengeRequest(*req)
	if err != nil {
		return nil, err
	}
	resp, err := challengeReq.Next()
	if err != nil {
		return nil, err
	}

	// do challenge process
	if resp.StatusCode == 401 {
		err := p.findServiceAndScope(resp)
		if err != nil {
			return nil, err
		}

		accessToken, err := p.getAccessToken(req)
		if err != nil {
			return nil, err
		}

		req.Raw().Header.Set(
			headerAuthorization,
			fmt.Sprintf("%s%s", bearerHeader, accessToken),
		)

		// send the original request with auth
		return req.Next()
	}

	return resp, nil
}

func (p *authenticationPolicy) getAccessToken(req *policy.Request) (string, error) {
	// anonymous access
	if p.cred == nil {
		resp, err := p.authClient.ExchangeACRRefreshTokenForACRAccessToken(req.Raw().Context(), p.acrService, p.acrScope, "", &authenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(tokenGrantTypePassword)})
		if err != nil {
			return "", err
		}
		return *resp.acrAccessToken.AccessToken, nil
	}

	// access with token
	as := acquiringResourceState{
		policy: p,
		req:    req,
	}

	// get refresh token from cache/request
	refreshToken, err := p.mainResource.Get(as)
	if err != nil {
		return "", err
	}

	// get access token from request
	resp, err := p.authClient.ExchangeACRRefreshTokenForACRAccessToken(req.Raw().Context(), p.acrService, p.acrScope, refreshToken.Token, &authenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(tokenGrantTypeRefreshToken)})
	if err != nil {
		return "", err
	}
	return *resp.acrAccessToken.AccessToken, nil
}

func (p *authenticationPolicy) findServiceAndScope(resp *http.Response) error {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return errors.New("response has no WWW-Authenticate header for challenge authentication")
	}

	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")
	parts := strings.Split(authHeader, "\",")
	valuesMap := map[string]string{}
	for _, part := range parts {
		subParts := strings.Split(part, "=")
		if len(subParts) == 2 {
			valuesMap[subParts[0]] = strings.ReplaceAll(subParts[1], "\"", "")
		}
	}

	if v, ok := valuesMap["scope"]; ok {
		p.acrScope = v
	}
	if p.acrScope == "" {
		return errors.New("could not find a valid scope in the WWW-Authenticate header")
	}

	if v, ok := valuesMap["service"]; ok {
		p.acrService = v
	}
	if p.acrService == "" {
		return errors.New("could not find a valid service in the WWW-Authenticate header")
	}

	return nil
}

func (p authenticationPolicy) getChallengeRequest(oriReq policy.Request) (*policy.Request, error) {
	copied := oriReq.Clone(oriReq.Raw().Context())
	err := copied.SetBody(nil, "")
	if err != nil {
		return nil, err
	}
	copied.Raw().Header.Del("Content-Type")
	return copied, nil
}

type acquiringResourceState struct {
	req    *policy.Request
	policy *authenticationPolicy
}

// acquire acquires or updates the resource; only one
// thread/goroutine at a time ever calls this function
func acquire(state acquiringResourceState) (newResource azcore.AccessToken, newExpiration time.Time, err error) {
	// get AAD token from credential
	aadToken, err := state.policy.cred.GetToken(
		state.req.Raw().Context(),
		policy.TokenRequestOptions{
			Scopes: state.policy.aadScopes,
		},
	)
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	// exchange refresh token with AAD token
	refreshResp, err := state.policy.authClient.ExchangeAADAccessTokenForACRRefreshToken(state.req.Raw().Context(), postContentSchemaGrantTypeAccessToken, state.policy.acrService, &authenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
		AccessToken: &aadToken.Token,
	})
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	refreshToken := azcore.AccessToken{
		Token: *refreshResp.acrRefreshToken.RefreshToken,
	}

	// get refresh token expire time
	refreshToken.ExpiresOn, err = getJWTExpireTime(*refreshResp.acrRefreshToken.RefreshToken)
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	// return refresh token
	return refreshToken, refreshToken.ExpiresOn, nil
}

func getJWTExpireTime(token string) (time.Time, error) {
	values := strings.Split(token, ".")
	if len(values) > 2 {
		value := values[1]
		padding := len(value) % 4
		if padding > 0 {
			for i := 0; i < 4-padding; i++ {
				value += "="
			}
		}
		parsedValue, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return time.Time{}, err
		}

		var jsonValue *jwtOnlyWithExp
		err = json.Unmarshal(parsedValue, &jsonValue)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(jsonValue.Exp, 0), nil
	}

	return time.Time{}, errors.New("could not parse refresh token expire time")
}

type jwtOnlyWithExp struct {
	Exp int64 `json:"exp"`
}
