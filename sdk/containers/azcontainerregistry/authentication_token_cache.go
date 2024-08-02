//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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
