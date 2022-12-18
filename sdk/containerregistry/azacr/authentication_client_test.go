//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthenticationClient_ExchangeAADAccessTokenForAcrRefreshToken(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	client := NewAuthenticationClient("https://azacrlivetest.azurecr.io", &AuthenticationClientOptions{ClientOptions: options})
	ctx := context.Background()
	accessToken, err := cred.GetToken(
		ctx,
		policy.TokenRequestOptions{
			Scopes: []string{"https://management.core.windows.net/.default"},
		})
	require.NoError(t, err)
	resp, err := client.ExchangeAADAccessTokenForAcrRefreshToken(ctx, PostContentSchemaGrantTypeAccessToken, "azacrlivetest.azurecr.io", &AuthenticationClientExchangeAADAccessTokenForAcrRefreshTokenOptions{
		AccessToken: &accessToken.Token,
	})
	require.NoError(t, err)
	require.NotEmpty(t, *resp.RefreshToken.RefreshToken)
}

func TestAuthenticationClient_ExchangeAcrRefreshTokenForAcrAccessToken(t *testing.T) {
	startRecording(t)
	cred, options := getCredAndClientOptions(t)
	client := NewAuthenticationClient("https://azacrlivetest.azurecr.io", &AuthenticationClientOptions{ClientOptions: options})
	ctx := context.Background()
	accessToken, err := cred.GetToken(
		ctx,
		policy.TokenRequestOptions{
			Scopes: []string{"https://management.core.windows.net/.default"},
		})
	require.NoError(t, err)
	refreshResp, err := client.ExchangeAADAccessTokenForAcrRefreshToken(ctx, PostContentSchemaGrantTypeAccessToken, "azacrlivetest.azurecr.io", &AuthenticationClientExchangeAADAccessTokenForAcrRefreshTokenOptions{
		AccessToken: &accessToken.Token,
	})
	require.NoError(t, err)
	require.NotEmpty(t, *refreshResp.RefreshToken.RefreshToken)
	accessResp, err := client.ExchangeAcrRefreshTokenForAcrAccessToken(ctx, "azacrlivetest.azurecr.io", "registry:catalog:*", *refreshResp.RefreshToken.RefreshToken, &AuthenticationClientExchangeAcrRefreshTokenForAcrAccessTokenOptions{GrantType: to.Ptr(TokenGrantTypeRefreshToken)})
	require.NoError(t, err)
	require.NotEmpty(t, *accessResp.AccessToken.AccessToken)
}
