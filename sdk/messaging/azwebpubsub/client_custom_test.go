//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azwebpubsub_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type stubCredential struct{}

func (stubCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "stub-aad-token"}, nil
}

// stubTransport records the last request and returns a canned generateToken response.
type stubTransport struct{ lastRequest *http.Request }

func (s *stubTransport) Do(req *http.Request) (*http.Response, error) {
	s.lastRequest = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"token":"server-issued-token"}`)),
		Request:    req,
	}, nil
}

func newTokenCredentialClient(t *testing.T, endpoint string) (*azwebpubsub.Client, *stubTransport) {
	transport := &stubTransport{}
	options := &azwebpubsub.ClientOptions{}
	options.Transport = transport

	client, err := azwebpubsub.NewClient(endpoint, stubCredential{}, options)
	require.NoError(t, err)
	return client, transport
}

// The endpoint passed to NewClient is not required to have a trailing slash, and the
// audience and client URL must be well formed either way.
func TestClient_GenerateClientAccessURLEndpointTrailingSlash(t *testing.T) {
	hub := "chat"

	for _, endpoint := range []string{
		"https://host.webpubsub.azure.com",
		"https://host.webpubsub.azure.com/",
	} {
		t.Run(endpoint, func(t *testing.T) {
			client, _ := newTokenCredentialClient(t, endpoint)

			resp, err := client.GenerateClientAccessURL(context.Background(), hub, nil)
			require.NoError(t, err)

			require.Equal(t, "wss://host.webpubsub.azure.com/client/hubs/chat", resp.BaseURL)
			require.Equal(t, "wss://host.webpubsub.azure.com/client/hubs/chat?access_token=server-issued-token", resp.URL)
		})
	}
}

// The same normalization applies to the audience embedded in a key-signed token.
func TestClient_GenerateClientAccessURLAudienceTrailingSlash(t *testing.T) {
	client, err := azwebpubsub.NewClientFromConnectionString("Endpoint=https://host.webpubsub.azure.com;AccessKey=ABC;", nil)
	require.NoError(t, err)

	resp, err := client.GenerateClientAccessURL(context.Background(), "chat", nil)
	require.NoError(t, err)

	claims := jwt.MapClaims{}
	_, _, err = jwt.NewParser().ParseUnverified(resp.Token, claims)
	require.NoError(t, err)

	// The service validates the host and path of the audience, so a missing separator
	// between the endpoint and "client/hubs" makes the token unusable.
	require.Equal(t, "https://host.webpubsub.azure.com/client/hubs/chat", claims["aud"])
	require.Equal(t, "wss://host.webpubsub.azure.com/client/hubs/chat", resp.BaseURL)
}

// nil options must be treated the same as empty options on both credential paths.
func TestClient_GenerateClientAccessURLNilOptions(t *testing.T) {
	t.Run("TokenCredential", func(t *testing.T) {
		client, transport := newTokenCredentialClient(t, "https://host.webpubsub.azure.com/")

		require.NotPanics(t, func() {
			resp, err := client.GenerateClientAccessURL(context.Background(), "chat", nil)
			require.NoError(t, err)
			require.Equal(t, "server-issued-token", resp.Token)
		})

		require.NotNil(t, transport.lastRequest)
	})

	t.Run("ConnectionString", func(t *testing.T) {
		client, err := azwebpubsub.NewClientFromConnectionString("Endpoint=https://host.webpubsub.azure.com;AccessKey=ABC;", nil)
		require.NoError(t, err)

		require.NotPanics(t, func() {
			resp, err := client.GenerateClientAccessURL(context.Background(), "chat", nil)
			require.NoError(t, err)
			require.NotEmpty(t, resp.Token)
		})
	})
}

// The service rejects minutesToExpire=0, so an unset expiration must fall back to the
// default rather than being forwarded as a zero value.
func TestClient_GenerateClientAccessURLMinutesToExpire(t *testing.T) {
	minutesToExpire := func(t *testing.T, transport *stubTransport) string {
		t.Helper()
		require.NotNil(t, transport.lastRequest)
		query, err := url.ParseQuery(transport.lastRequest.URL.RawQuery)
		require.NoError(t, err)
		return query.Get("minutesToExpire")
	}

	t.Run("unset uses the default", func(t *testing.T) {
		client, transport := newTokenCredentialClient(t, "https://host.webpubsub.azure.com/")

		_, err := client.GenerateClientAccessURL(context.Background(), "chat", &azwebpubsub.GenerateClientAccessURLOptions{})
		require.NoError(t, err)
		require.Equal(t, "60", minutesToExpire(t, transport))
	})

	t.Run("nil options uses the default", func(t *testing.T) {
		client, transport := newTokenCredentialClient(t, "https://host.webpubsub.azure.com/")

		_, err := client.GenerateClientAccessURL(context.Background(), "chat", nil)
		require.NoError(t, err)
		require.Equal(t, "60", minutesToExpire(t, transport))
	})

	t.Run("explicit value is honored", func(t *testing.T) {
		client, transport := newTokenCredentialClient(t, "https://host.webpubsub.azure.com/")

		_, err := client.GenerateClientAccessURL(context.Background(), "chat", &azwebpubsub.GenerateClientAccessURLOptions{
			ExpirationTimeInMinutes: 5,
		})
		require.NoError(t, err)
		require.Equal(t, "5", minutesToExpire(t, transport))
	})
}

// A negative expiration is rejected on both credential paths.
func TestClient_GenerateClientAccessURLNegativeExpiration(t *testing.T) {
	options := &azwebpubsub.GenerateClientAccessURLOptions{ExpirationTimeInMinutes: -1}

	client, _ := newTokenCredentialClient(t, "https://host.webpubsub.azure.com/")
	_, err := client.GenerateClientAccessURL(context.Background(), "chat", options)
	require.ErrorContains(t, err, "out of range")

	keyClient, err := azwebpubsub.NewClientFromConnectionString("Endpoint=https://host.webpubsub.azure.com;AccessKey=ABC;", nil)
	require.NoError(t, err)
	_, err = keyClient.GenerateClientAccessURL(context.Background(), "chat", options)
	require.ErrorContains(t, err, "out of range")
}

// Options are forwarded to the service on the TokenCredential path.
func TestClient_GenerateClientAccessURLForwardsOptions(t *testing.T) {
	client, transport := newTokenCredentialClient(t, "https://host.webpubsub.azure.com/")

	_, err := client.GenerateClientAccessURL(context.Background(), "chat", &azwebpubsub.GenerateClientAccessURLOptions{
		UserID: "user1",
		Roles:  []string{"admin"},
		Groups: []string{"group1"},
	})
	require.NoError(t, err)

	query, err := url.ParseQuery(transport.lastRequest.URL.RawQuery)
	require.NoError(t, err)
	require.Equal(t, "user1", query.Get("userId"))
	require.Equal(t, []string{"admin"}, query["role"])
	require.Equal(t, []string{"group1"}, query["group"])
}
