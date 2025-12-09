// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func Test_authenticationClient_ExchangeAADAccessTokenForACRRefreshToken(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{ClientOptions: options})
	require.NoError(t, err)
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
	resp, err := client.ExchangeAADAccessTokenForACRRefreshToken(ctx, PostContentSchemaGrantTypeAccessToken, strings.TrimPrefix(endpoint, "https://"), &AuthenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
		AccessToken: &accessToken.Token,
	})
	require.NoError(t, err)
	require.NotEmpty(t, *resp.RefreshToken)
}

func Test_authenticationClient_ExchangeAADAccessTokenForACRRefreshToken_fail(t *testing.T) {
	startRecording(t)
	endpoint, _, options := getEndpointCredAndClientOptions(t)
	client, err := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{ClientOptions: options})
	require.NoError(t, err)
	ctx := context.Background()
	_, err = client.ExchangeAADAccessTokenForACRRefreshToken(ctx, PostContentSchemaGrantTypeAccessToken, strings.TrimPrefix(endpoint, "https://"), &AuthenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
		Tenant:       to.Ptr("wrong tenant"),
		RefreshToken: to.Ptr("wrong token"),
		AccessToken:  to.Ptr("wrong token"),
	})
	require.Error(t, err)
}

func Test_authenticationClient_ExchangeAADAccessTokenForACRRefreshToken_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("wrong response")))
	client, err := NewAuthenticationClient(srv.URL(), &AuthenticationClientOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	require.NoError(t, err)
	ctx := context.Background()
	_, err = client.ExchangeAADAccessTokenForACRRefreshToken(ctx, "grantType", "service", nil)
	require.Error(t, err)
}

func Test_authenticationClient_ExchangeACRRefreshTokenForACRAccessToken(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	client, err := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{ClientOptions: options})
	require.NoError(t, err)
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
	refreshResp, err := client.ExchangeAADAccessTokenForACRRefreshToken(ctx, PostContentSchemaGrantTypeAccessToken, strings.TrimPrefix(endpoint, "https://"), &AuthenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
		AccessToken: &accessToken.Token,
	})
	require.NoError(t, err)
	require.NotEmpty(t, *refreshResp.RefreshToken)
	accessResp, err := client.ExchangeACRRefreshTokenForACRAccessToken(ctx, strings.TrimPrefix(endpoint, "https://"), "registry:catalog:*", *refreshResp.RefreshToken, &AuthenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(TokenGrantTypeRefreshToken)})
	require.NoError(t, err)
	require.NotEmpty(t, *accessResp.AccessToken)
}

func Test_authenticationClient_ExchangeACRRefreshTokenForACRAccessToken_fail(t *testing.T) {
	startRecording(t)
	endpoint, _, options := getEndpointCredAndClientOptions(t)
	client, err := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{ClientOptions: options})
	require.NoError(t, err)
	ctx := context.Background()
	_, err = client.ExchangeACRRefreshTokenForACRAccessToken(ctx, strings.TrimPrefix(endpoint, "https://"), "registry:catalog:*", "wrong token", &AuthenticationClientExchangeACRRefreshTokenForACRAccessTokenOptions{GrantType: to.Ptr(TokenGrantTypeRefreshToken)})
	require.Error(t, err)
}

func Test_authenticationClient_ExchangeACRRefreshTokenForACRAccessToken_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("wrong response")))
	client, err := NewAuthenticationClient(srv.URL(), &AuthenticationClientOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	require.NoError(t, err)
	ctx := context.Background()
	_, err = client.ExchangeACRRefreshTokenForACRAccessToken(ctx, "service", "scope", "refresh token", nil)
	require.Error(t, err)
}
