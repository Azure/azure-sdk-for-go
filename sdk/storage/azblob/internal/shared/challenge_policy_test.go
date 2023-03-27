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
	"net/http"
	"strings"
	"testing"
	"time"
)

type credentialFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (cf credentialFunc) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return cf(ctx, options)
}

func TestChallengePolicy(t *testing.T) {
	accessToken := "***"
	storageResource := "https://storage.azure.com"
	storageScope := "https://storage.azure.com/.default"
	challenge := `Bearer authorization_uri="https://login.microsoftonline.com/{tenant}", resource_id="{storageResource}"`
	diskResource := "https://disk.azure.com/"
	diskScope := "https://disk.azure.com//.default"

	for _, test := range []struct {
		expectedScope, format, resource string
	}{
		{format: challenge, resource: storageResource, expectedScope: storageScope},
		{format: challenge, resource: diskResource, expectedScope: diskScope},
	} {
		t.Run("", func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(
				mock.WithHeader("WWW-Authenticate", strings.ReplaceAll(test.format, "{storageResource}", test.resource)),
				mock.WithStatusCode(401),
			)
			srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
				if authz := r.Header.Values("Authorization"); len(authz) != 1 || authz[0] != "Bearer "+accessToken {
					t.Errorf(`unexpected Authorization "%s"`, authz)
				}
				return true
			}))
			srv.AppendResponse()
			authenticated := false
			cred := credentialFunc(func(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				authenticated = true
				require.Equal(t, []string{test.expectedScope}, tro.Scopes)
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
		})
	}
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
