//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

const (
	headerAuthorization = "Authorization"
	bearerHeader        = "Bearer "
)

type authenticationPolicyOptions struct {
	*authenticationTokenCacheOptions
}

// authenticationPolicy is a policy to do the challenge-based authentication for container registry service. The authorization flow is as follows:
// Step 1: GET /api/v1/acr/repositories
// Return Header: 401: www-authenticate header - Bearer realm="{url}",service="{serviceName}",scope="{scope}",error="invalid_token"
// Step 2: Retrieve the serviceName, scope from the WWW-Authenticate header.
// Step 3: POST /api/oauth2/exchange
// Request Body : { service, scope, grant-type, aadToken with ARM scope }
// Response Body: { refreshToken }
// Step 4: POST /api/oauth2/token
// Request Body: { refreshToken, scope, grant-type }
// Response Body: { accessToken }
// Step 5: GET /api/v1/acr/repositories
// Request Header: { Bearer acrTokenAccess }
// Each registry service shares one refresh token, it will be cached in refreshTokenCache until expire time.
// Since the scope will be different for different API/repository/artifact, accessTokenCache will only work when continuously calling same API.
type authenticationPolicy struct {
	accessTokenCache *authenticationTokenCache
}

func newAuthenticationPolicy(cred azcore.TokenCredential, scopes []string, authClient *authenticationClient, opts *authenticationPolicyOptions) *authenticationPolicy {
	if opts == nil {
		opts = &authenticationPolicyOptions{}
	}
	return &authenticationPolicy{
		accessTokenCache: newAuthenticationTokenCache(cred, scopes, authClient, opts.authenticationTokenCacheOptions),
	}
}

func (p *authenticationPolicy) Do(req *policy.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	if req.Raw().Header.Get(headerAuthorization) != "" {
		// retry request could do the request with existed token directly
		resp, err = req.Next()
	} else if accessToken := p.accessTokenCache.Load(); accessToken != "" {
		// if there is a previous access token, then we try to use this token to do the request
		req.Raw().Header.Set(
			headerAuthorization,
			fmt.Sprintf("%s%s", bearerHeader, accessToken),
		)
		resp, err = req.Next()
	} else {
		// do challenge process for the initial request
		var challengeReq *policy.Request
		challengeReq, err = getChallengeRequest(*req)
		if err != nil {
			return nil, err
		}
		resp, err = challengeReq.Next()
	}
	if err != nil {
		return nil, err
	}

	// if 401 response, then try to get access token
	if resp.StatusCode == http.StatusUnauthorized {
		var service, scope, accessToken string
		if service, scope, err = findServiceAndScope(resp); err != nil {
			return nil, err
		}
		if accessToken, err = p.accessTokenCache.AcquireAccessToken(req.Raw().Context(), service, scope); err != nil {
			return nil, err
		}
		req.Raw().Header.Set(
			headerAuthorization,
			fmt.Sprintf("%s%s", bearerHeader, accessToken),
		)
		// since the request may already been used once, body should be rewound
		if err = req.RewindBody(); err != nil {
			return nil, err
		}
		return req.Next()
	}

	return resp, nil
}

func (c *authenticationTokenCache) AcquireAccessToken(ctx context.Context, service, scope string) (string, error) {
	// anonymous access
	if c.cred == nil {
		resp, err := c.authClient.ExchangeACRRefreshTokenForACRAccessToken(ctx, service, scope, "", &authenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(tokenGrantTypePassword)})
		if err != nil {
			return "", err
		}
		c.accessTokenCache.Store(*resp.acrAccessToken.AccessToken)
		return *resp.acrAccessToken.AccessToken, nil
	}

	// access with token
	// get refresh token from cache/request
	refreshToken, err := c.refreshTokenCache.Get(acquiringResourceState{
		ctx:           ctx,
		aadCredential: c.cred,
		aadScopes:     c.aadScopes,
		authClient:    c.authClient,
		service:       service,
	})
	if err != nil {
		return "", err
	}

	// get access token from request
	resp, err := c.authClient.ExchangeACRRefreshTokenForACRAccessToken(ctx, service, scope, refreshToken.Token, &authenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(tokenGrantTypeRefreshToken)})
	if err != nil {
		return "", err
	}
	c.accessTokenCache.Store(*resp.acrAccessToken.AccessToken)
	return *resp.acrAccessToken.AccessToken, nil
}

func findServiceAndScope(resp *http.Response) (string, string, error) {
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return "", "", errors.New("response has no WWW-Authenticate header for challenge authentication")
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

	if _, ok := valuesMap["service"]; !ok {
		return "", "", errors.New("could not find a valid service in the WWW-Authenticate header")
	}

	if _, ok := valuesMap["scope"]; !ok {
		return "", "", errors.New("could not find a valid scope in the WWW-Authenticate header")
	}

	return valuesMap["service"], valuesMap["scope"], nil
}

func getChallengeRequest(oriReq policy.Request) (*policy.Request, error) {
	copied := oriReq.Clone(oriReq.Raw().Context())
	err := copied.SetBody(nil, "")
	if err != nil {
		return nil, err
	}
	copied.Raw().Header.Del("Content-Type")
	return copied, nil
}

type acquiringResourceState struct {
	ctx context.Context

	aadCredential azcore.TokenCredential
	aadScopes     []string

	authClient *authenticationClient
	service    string
}

// acquireRefreshToken acquires or updates the refresh token of ACR service; only one thread/goroutine at a time ever calls this function
func acquireRefreshToken(state acquiringResourceState) (newResource azcore.AccessToken, newExpiration time.Time, err error) {
	// get AAD token from credential
	aadToken, err := state.aadCredential.GetToken(
		state.ctx,
		policy.TokenRequestOptions{
			Scopes: state.aadScopes,
		},
	)
	if err != nil {
		return azcore.AccessToken{}, time.Time{}, err
	}

	// exchange refresh token with AAD token
	refreshResp, err := state.authClient.ExchangeAADAccessTokenForACRRefreshToken(state.ctx, postContentSchemaGrantTypeAccessToken, state.service, &authenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
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
