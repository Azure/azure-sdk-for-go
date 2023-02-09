//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
	"reflect"
	"strings"
	"testing"
)

func Test_authenticationClient_ExchangeAADAccessTokenForACRRefreshToken(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client := newAuthenticationClient(endpoint, &authenticationClientOptions{ClientOptions: options})
	ctx := context.Background()
	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	accessToken, err := cred.GetToken(
		ctx,
		policy.TokenRequestOptions{
			Scopes: []string{options.Cloud.Services[ServiceName].Audience + "/.default"},
		})
	require.NoError(t, err)
	resp, err := client.ExchangeAADAccessTokenForACRRefreshToken(ctx, postContentSchemaGrantTypeAccessToken, strings.TrimPrefix(endpoint, "https://"), &authenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
		AccessToken: &accessToken.Token,
	})
	require.NoError(t, err)
	require.NotEmpty(t, *resp.acrRefreshToken.RefreshToken)
}

func Test_authenticationClient_ExchangeACRRefreshTokenForACRAccessToken(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client := newAuthenticationClient(endpoint, &authenticationClientOptions{ClientOptions: options})
	ctx := context.Background()
	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	accessToken, err := cred.GetToken(
		ctx,
		policy.TokenRequestOptions{
			Scopes: []string{options.Cloud.Services[ServiceName].Audience + "/.default"},
		})
	require.NoError(t, err)
	refreshResp, err := client.ExchangeAADAccessTokenForACRRefreshToken(ctx, postContentSchemaGrantTypeAccessToken, strings.TrimPrefix(endpoint, "https://"), &authenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
		AccessToken: &accessToken.Token,
	})
	require.NoError(t, err)
	require.NotEmpty(t, *refreshResp.acrRefreshToken.RefreshToken)
	accessResp, err := client.ExchangeACRRefreshTokenForACRAccessToken(ctx, strings.TrimPrefix(endpoint, "https://"), "registry:catalog:*", *refreshResp.acrRefreshToken.RefreshToken, &authenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(tokenGrantTypeRefreshToken)})
	require.NoError(t, err)
	require.NotEmpty(t, *accessResp.acrAccessToken.AccessToken)
}
