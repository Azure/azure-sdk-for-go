//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azwebpubsub_test

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub/internal"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestClient_SendToAll(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	client := newClientWrapper(t)

	_, err := client.SendToAll(context.Background(),
		client.TestVars.Hub, azwebpubsub.ContentTypeTextPlain, newStream("Hello world!"),
		&azwebpubsub.ClientSendToAllOptions{})
	require.NoError(t, err)

	_, err = client.SendToAll(context.Background(),
		client.TestVars.Hub, azwebpubsub.ContentTypeApplicationJSON, newStream("true"),
		&azwebpubsub.ClientSendToAllOptions{})
	require.NoError(t, err)

	_, err = client.SendToAll(context.Background(),
		client.TestVars.Hub, azwebpubsub.ContentTypeApplicationOctetStream, newStream("true"),
		&azwebpubsub.ClientSendToAllOptions{})
	require.NoError(t, err)
}

func TestClient_CloseConnections(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	client := newClientWrapper(t)
	const hub = "chat"
	const conn1 = "conn1"
	const group1 = "group1"
	const user1 = "user1"
	reason := "TestClient_CloseConnections"
	_, err := client.CloseAllConnections(context.Background(),
		hub, &azwebpubsub.ClientCloseAllConnectionsOptions{Excluded: []string{conn1}, Reason: &reason})
	require.NoError(t, err)
	_, err = client.CloseConnection(context.Background(),
		hub, conn1, &azwebpubsub.ClientCloseConnectionOptions{Reason: &reason})
	require.NoError(t, err)
	_, err = client.CloseGroupConnections(context.Background(),
		hub, group1, &azwebpubsub.ClientCloseGroupConnectionsOptions{Excluded: []string{conn1}, Reason: &reason})
	require.NoError(t, err)
	_, err = client.CloseUserConnections(context.Background(),
		hub, user1, &azwebpubsub.ClientCloseUserConnectionsOptions{Excluded: []string{conn1}, Reason: &reason})
	require.NoError(t, err)
}

func TestClient_GenerateClientAccessUrl(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	client := newClientWrapper(t)
	token, err := client.GenerateClientAccessUrl(context.Background(), nil)
	require.NoError(t, err)
	extract := extractToken(t, token, client)
	require.Nil(t, extract.Roles)
	require.Nil(t, extract.Groups)
	require.Nil(t, extract.UserID)

	user1 := "user1"
	token, err = client.GenerateClientAccessUrl(context.Background(), &azwebpubsub.GenerateClientAccessUrlOptions{
		UserID: &user1,
		Roles:  []string{"admin"},
		Groups: []string{"group1"},
	})
	require.NoError(t, err)
	parsedURL, err := url.Parse(token.URL)
	require.NoError(t, err)
	queryValues := parsedURL.Query()
	accessToken := queryValues.Get("access_token")
	require.NotEmpty(t, accessToken)
	extract = extractToken(t, token, client)
	require.Equal(t, user1, *extract.UserID)
	require.Equal(t, "admin", extract.Roles[0])
	require.Equal(t, "group1", extract.Groups[0])
}

func extractToken(t *testing.T, token *azwebpubsub.GenerateClientAccessUrlResponse, client clientWrapper) azwebpubsub.GenerateClientAccessUrlOptions {
	key, err := internal.ParseConnectionString(client.TestVars.ConnectionString)
	require.NoError(t, err)
	expectedAudience := key.Endpoint + "client/hubs/" + client.TestVars.Hub
	expectedBaseUrl := regexp.MustCompile(`$(http)(s?://)`).ReplaceAllString(expectedAudience, "ws$2")

	require.Equal(t, expectedBaseUrl, token.BaseURL)
	parsed, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key for validation
		return []byte(key.AccessKey), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid, "token is not valid")
	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok, "claims is not valid")
	audience, ok := claims["aud"].(string)
	require.True(t, ok, "audience is not valid")
	require.Equal(t, expectedAudience, audience)
	subject, ok := claims["sub"].(string)
	var userId *string
	if ok {
		userId = &subject
	} else {
		userId = nil
	}
	rawRoles, ok := claims["role"].([]interface{})
	var roles []string
	if ok {
		// Convert the interface slice to a slice of strings
		for _, role := range rawRoles {
			if r, ok := role.(string); ok {
				roles = append(roles, r)
			}
		}
	} else {
		roles = nil
	}

	rawGroups, ok := claims["webpubsub.group"].([]interface{})
	var groups []string
	if ok {
		// Convert the interface slice to a slice of strings
		for _, group := range rawGroups {
			if r, ok := group.(string); ok {
				groups = append(groups, r)
			}
		}
	} else {
		groups = nil
	}

	return azwebpubsub.GenerateClientAccessUrlOptions{
		UserID: userId,
		Roles:  roles,
		Groups: groups,
	}
}

func newStream(message string) io.ReadSeekCloser {
	return streaming.NopCloser(bytes.NewReader([]byte(message)))
}
