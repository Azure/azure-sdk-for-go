//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

type credentialFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (cf credentialFunc) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return cf(ctx, options)
}

func TestChallengePolicyStorage(t *testing.T) {
	accessToken := "***"
	storageScope := "https://storage.azure.com/.default"

	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(200),
	)
	authenticated := false
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		authenticated = true
		require.Equal(t, []string{storageScope}, tro.Scopes)
		return azcore.AccessToken{Token: accessToken, ExpiresOn: time.Now().Add(time.Hour)}, nil
	})
	p := NewStorageChallengePolicy(cred)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)
	req, err := runtime.NewRequest(context.Background(), "GET", "https://localhost")
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.True(t, authenticated, "policy should have authenticated")
}

func TestChallengePolicyStorageTenantID(t *testing.T) {
	accessToken := "***"
	storageResource := "https://storage.azure.com"
	storageScope := "https://storage.azure.com/.default"
	tenantID := "faketenantid"
	challenge := `Bearer authorization_uri="https://login.microsoftonline.com/{tenant}", resource_id="{storageResource}"`

	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", strings.ReplaceAll(strings.ReplaceAll(challenge, "{storageResource}", storageResource), "{tenant}", tenantID)),
		mock.WithStatusCode(401),
	)
	srv.AppendResponse(
		mock.WithStatusCode(200),
	)
	attemptedAuthentication := false
	authenticated := false
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		if attemptedAuthentication {
			authenticated = true
			require.Equal(t, []string{storageScope}, tro.Scopes)
			require.Equal(t, tenantID, tro.TenantID)
			return azcore.AccessToken{Token: accessToken, ExpiresOn: time.Now().Add(time.Hour)}, nil
		}
		attemptedAuthentication = true
		return azcore.AccessToken{}, nil
	})
	p := NewStorageChallengePolicy(cred)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)
	req, err := runtime.NewRequest(context.Background(), "GET", "https://localhost")
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.True(t, authenticated, "policy should have authenticated")
}

func TestChallengePolicyDisk(t *testing.T) {
	accessToken := "***"
	diskResource := "https://disk.azure.com/"
	diskScope := "https://disk.azure.com//.default"
	tenantID := "faketenantid"
	challenge := `Bearer authorization_uri="https://login.microsoftonline.com/{tenant}", resource_id="{storageResource}"`

	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithHeader("WWW-Authenticate", strings.ReplaceAll(strings.ReplaceAll(challenge, "{storageResource}", diskResource), "{tenant}", tenantID)),
		mock.WithStatusCode(401),
	)
	srv.AppendResponse(
		mock.WithStatusCode(200),
	)
	attemptedAuthentication := false
	authenticated := false
	cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
		if attemptedAuthentication {
			authenticated = true
			require.Equal(t, []string{diskScope}, tro.Scopes)
			require.Equal(t, tenantID, tro.TenantID)
			return azcore.AccessToken{Token: accessToken, ExpiresOn: time.Now().Add(time.Hour)}, nil
		}
		attemptedAuthentication = true
		return azcore.AccessToken{}, nil
	})
	p := NewStorageChallengePolicy(cred)
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerRetry: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)
	req, err := runtime.NewRequest(context.Background(), "GET", "https://localhost")
	require.NoError(t, err)
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.True(t, authenticated, "policy should have authenticated")
}

func TestParseTenant(t *testing.T) {
	actual := parseTenant("")
	require.Empty(t, actual)

	expected := "00000000-0000-0000-0000-000000000000"
	sampleURL := "https://login.microsoftonline.com/" + expected
	actual = parseTenant(sampleURL)
	require.Equal(t, expected, actual, "tenant was not properly parsed")
}

func TestParseTenantNegative(t *testing.T) {
	actual := parseTenant("")
	require.Empty(t, actual)

	expected := ""
	sampleURL := "https://login.microsoftonline.com/" + expected
	actual = parseTenant(sampleURL)
	require.Equal(t, expected, actual)

	sampleURL = ""
	actual = parseTenant(sampleURL)
	require.Equal(t, expected, actual)
}
