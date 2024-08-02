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
	"strings"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

type authenticationTokenCacheOptions struct{}

type authenticationTokenCache struct {
	refreshTokenCache *temporal.Resource[azcore.AccessToken, acquiringResourceState]
	accessTokenCache  atomic.Value
	cred              azcore.TokenCredential
	aadScopes         []string
	authClient        *authenticationClient
}

func newAuthenticationTokenCache(cred azcore.TokenCredential, scopes []string, authClient *authenticationClient, opts *authenticationTokenCacheOptions) *authenticationTokenCache {
	return &authenticationTokenCache{
		cred:              cred,
		aadScopes:         scopes,
		authClient:        authClient,
		refreshTokenCache: temporal.NewResource(acquireRefreshToken),
	}
}

func (c *authenticationTokenCache) Load() string {
	value, ok := c.accessTokenCache.Load().(string)
	if !ok {
		return ""
	}
	return value
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
