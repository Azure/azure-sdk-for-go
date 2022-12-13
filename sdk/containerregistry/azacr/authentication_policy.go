//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
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

type AuthenticationPolicyOptions struct {
}

type AuthenticationPolicy struct {
	mainResource *temporal.Resource[azcore.AccessToken, acquiringResourceState]
	cred         azcore.TokenCredential
	aadScopes    []string
	acrScope     string
	acrService   string
	authClient   *AuthenticationClient
}

func NewAuthenticationPolicy(cred azcore.TokenCredential, scopes []string, authClient *AuthenticationClient, opts *AuthenticationPolicyOptions) *AuthenticationPolicy {
	return &AuthenticationPolicy{
		cred:         cred,
		aadScopes:    scopes,
		authClient:   authClient,
		mainResource: temporal.NewResource(acquire),
	}
}

func (p *AuthenticationPolicy) Do(req *policy.Request) (*http.Response, error) {
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

func (p *AuthenticationPolicy) getAccessToken(req *policy.Request) (string, error) {
	// anonymous access
	if p.cred == nil {
		resp, err := p.authClient.ExchangeAcrRefreshTokenForAcrAccessToken(req.Raw().Context(), p.acrService, p.acrScope, "", &AuthenticationClientExchangeAcrRefreshTokenForAcrAccessTokenOptions{GrantType: to.Ptr(TokenGrantTypePassword)})
		if err != nil {
			return "", err
		}
		return *resp.AccessToken.AccessToken, nil
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
	resp, err := p.authClient.ExchangeAcrRefreshTokenForAcrAccessToken(req.Raw().Context(), p.acrService, p.acrScope, refreshToken.Token, &AuthenticationClientExchangeAcrRefreshTokenForAcrAccessTokenOptions{GrantType: to.Ptr(TokenGrantTypeRefreshToken)})
	if err != nil {
		return "", err
	}
	return *resp.AccessToken.AccessToken, nil
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

func (p *AuthenticationPolicy) findServiceAndScope(resp *http.Response) error {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return &challengePolicyError{err: errors.New("response has no WWW-Authenticate header for challenge authentication")}
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
	//if p.acrScope == "" {
	//	return &challengePolicyError{err: errors.New("could not find a valid scope in the WWW-Authenticate header")}
	//}
	if p.acrScope == "" {
		p.acrScope = "registry:catalog:*"
	}

	if v, ok := valuesMap["service"]; ok {
		p.acrService = v
	}
	if p.acrService == "" {
		return &challengePolicyError{err: errors.New("could not find a valid service in the WWW-Authenticate header")}
	}

	return nil
}

func (p AuthenticationPolicy) getChallengeRequest(orig policy.Request) (*policy.Request, error) {
	req, err := runtime.NewRequest(orig.Raw().Context(), orig.Raw().Method, orig.Raw().URL.String())
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
	req    *policy.Request
	policy *AuthenticationPolicy
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
	refreshResp, err := state.policy.authClient.ExchangeAADAccessTokenForAcrRefreshToken(state.req.Raw().Context(), PostContentSchemaGrantTypeAccessToken, state.policy.acrService, &AuthenticationClientExchangeAADAccessTokenForAcrRefreshTokenOptions{
		AccessToken: &aadToken.Token,
	})
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	refreshToken := azcore.AccessToken{
		Token: *refreshResp.RefreshToken.RefreshToken,
	}

	// get refresh token expire time
	refreshToken.ExpiresOn, err = getJWTExpireTime(*refreshResp.RefreshToken.RefreshToken)
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
			for i := 0; i < padding; i++ {
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

	return time.Time{}, &challengePolicyError{err: errors.New("could not parse refresh token expire time")}
}

type jwtOnlyWithExp struct {
	Exp int64 `json:"exp"`
}
