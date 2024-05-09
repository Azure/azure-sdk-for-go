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
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestClient_SendToAll(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	client := newClientWrapper(t)
	hub := "hub1"
	_, err := client.SendToAll(context.Background(), hub,
		azwebpubsub.ContentTypeTextPlain, newStream("Hello world!"),
		&azwebpubsub.SendToAllOptions{})
	require.NoError(t, err)

	_, err = client.SendToAll(context.Background(), hub,
		azwebpubsub.ContentTypeApplicationJSON, newStream("true"),
		&azwebpubsub.SendToAllOptions{})
	require.NoError(t, err)

	_, err = client.SendToAll(context.Background(), hub,
		azwebpubsub.ContentTypeApplicationOctetStream, newStream("true"),
		&azwebpubsub.SendToAllOptions{})
	require.NoError(t, err)
}

func TestClient_ManagePermissions(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	client := newClientWrapper(t)
	const hub = "chat"
	const conn1 = "conn1"
	group := "group1"
	_, err := client.GrantPermission(context.Background(), hub, azwebpubsub.PermissionJoinLeaveGroup, conn1, &azwebpubsub.GrantPermissionOptions{
		TargetName: &group,
	})
	require.ErrorContains(t, err, "404 Not Found")
	_, err = client.RevokePermission(context.Background(), hub, azwebpubsub.PermissionJoinLeaveGroup, conn1, &azwebpubsub.RevokePermissionOptions{
		TargetName: &group,
	})
	require.NoError(t, err)
}

func TestClient_CloseConnections(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	client := newClientWrapper(t)
	const hub = "chat"
	const conn1 = "conn1"
	const group1 = "group1"
	const user1 = "user1"
	reason := "TestClient_CloseConnections"
	_, err := client.CloseAllConnections(context.Background(),
		hub, &azwebpubsub.CloseAllConnectionsOptions{Excluded: []string{conn1}, Reason: &reason})
	require.NoError(t, err)
	_, err = client.CloseConnection(context.Background(),
		hub, conn1, &azwebpubsub.CloseConnectionOptions{Reason: &reason})
	require.NoError(t, err)
	_, err = client.CloseGroupConnections(context.Background(),
		hub, group1, &azwebpubsub.CloseGroupConnectionsOptions{Excluded: []string{conn1}, Reason: &reason})
	require.NoError(t, err)
	_, err = client.CloseUserConnections(context.Background(),
		hub, user1, &azwebpubsub.CloseUserConnectionsOptions{Excluded: []string{conn1}, Reason: &reason})
	require.NoError(t, err)
}

func TestClient_GenerateClientAccessURLFromConnectionString(t *testing.T) {
	_, err1 := azwebpubsub.NewClientFromConnectionString("Endpoint=http://test/subpath;;;;", nil)
	require.ErrorContains(t, err1, "connection string is either blank or malformed.")

	hub := "chat/go"
	client, err := azwebpubsub.NewClientFromConnectionString("Endpoint=http://test/subpath;AccessKey=ABC;;;", nil)
	require.NoError(t, err)

	token, err := client.GenerateClientAccessURL(context.Background(), hub, nil)

	require.NoError(t, err)
	extract := extractToken(t, token, "http://test/subpath", "ABC", hub)
	require.Nil(t, extract.Roles)
	require.Nil(t, extract.Groups)
	require.Empty(t, extract.UserID)

	user1 := "user1"
	token, err = client.GenerateClientAccessURL(context.Background(), hub, &azwebpubsub.GenerateClientAccessURLOptions{
		UserID: user1,
		Roles:  []string{"admin"},
		Groups: []string{"group1"},
	})
	require.NoError(t, err)
	parsedURL, err := url.Parse(token.URL)
	require.NoError(t, err)
	queryValues := parsedURL.Query()
	accessToken := queryValues.Get("access_token")
	require.NotEmpty(t, accessToken)
	extract = extractToken(t, token, "http://test/subpath", "ABC", hub)
	require.Equal(t, user1, extract.UserID)
	require.Equal(t, "admin", extract.Roles[0])
	require.Equal(t, "group1", extract.Groups[0])
}

func extractToken(t *testing.T, token *azwebpubsub.GenerateClientAccessURLResponse, endpoint string, key string, hub string) azwebpubsub.GenerateClientAccessURLOptions {
	expectedAudience := endpoint + "/client/hubs/" + url.PathEscape(hub)
	expectedBaseUrl := strings.Replace(expectedAudience, "http", "ws", 1)

	require.Equal(t, expectedBaseUrl, token.BaseURL)
	parsed, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key for validation
		return []byte(key), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid, "token is not valid")
	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok, "claims is not valid")
	audience, ok := claims["aud"].(string)
	require.True(t, ok, "audience is not valid")
	require.Equal(t, expectedAudience, audience)
	subject, ok := claims["sub"].(string)
	var userId string
	if ok {
		userId = subject
	} else {
		userId = ""
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

	return azwebpubsub.GenerateClientAccessURLOptions{
		UserID: userId,
		Roles:  roles,
		Groups: groups,
	}
}

func newStream(message string) io.ReadSeekCloser {
	return streaming.NopCloser(bytes.NewReader([]byte(message)))
}
