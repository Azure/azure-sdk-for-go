//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_getJWTExpireTime(t *testing.T) {
	for _, test := range []struct {
		name   string
		token  string
		expire time.Time
		err    bool
	}{
		{
			"test1",
			".ewogICJqdGkiOiAiMzY1ZTNiNWItODQ0ZS00YTIxLWEzOGMtNGQ4YWViZGQ2YTA2IiwKICAic3ViIjogInVzZXJAY29udG9zby5jb20iLAogICJuYmYiOiAxNDk3OTg4NzEyLAogICJleHAiOiAxNDk3OTkwODAxLAogICJpYXQiOiAxNDk3OTg4NzEyLAogICJpc3MiOiAiQXp1cmUgQ29udGFpbmVyIFJlZ2lzdHJ5IiwKICAiYXVkIjogImNvbnRvc29yZWdpc3RyeS5henVyZWNyLmlvIiwKICAidmVyc2lvbiI6ICIxLjAiLAogICJncmFudF90eXBlIjogInJlZnJlc2hfdG9rZW4iLAogICJ0ZW5hbnQiOiAiNDA5NTIwZDQtODEwMC00ZDFkLWFkNDctNzI0MzJkZGNjMTIwIiwKICAicGVybWlzc2lvbnMiOiB7CiAgICAiYWN0aW9ucyI6IFsKICAgICAgIioiCiAgICBdLAogICAgIm5vdEFjdGlvbnMiOiBbXQogIH0sCiAgInJvbGVzIjogW10KfQ==.",
			time.Unix(1497990801, 0),
			false,
		},
		{
			"test2",
			".eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJuYmYiOjE2NzA0MTA1NDEsImV4cCI6MTY3MDQyMjI0MSwiaWF0IjoxNjcwNDEwNTQxLCJpc3MiOiJBenVyZSBDb250YWluZXIgUmVnaXN0cnkiLCJhdWQiOiJhemFjcmxpdmV0ZXN0LmF6dXJlY3IuaW8iLCJ2ZXJzaW9uIjoiMS4wIiwicmlkIjoiMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJncmFudF90eXBlIjoicmVmcmVzaF90b2tlbiIsImFwcGlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwicGVybWlzc2lvbnMiOnsiQWN0aW9ucyI6WyJyZWFkIiwid3JpdGUiLCJkZWxldGUiLCJkZWxldGVkL3JlYWQiLCJkZWxldGVkL3Jlc3RvcmUvYWN0aW9uIl0sIk5vdEFjdGlvbnMiOm51bGx9LCJyb2xlcyI6W119.",
			time.Unix(1670422241, 0),
			false,
		},
		{
			"test-padding",
			".ewogICJqdGkiOiAiMzY1ZTNiNWItODQ0ZS00YTIxLWEzOGMtNGQ4YWViZGQ2YTA2IiwKICAic3ViIjogInVzZXJAY29udG9zby5jb20iLAogICJuYmYiOiAxNDk3OTg4NzEyLAogICJleHAiOiAxNDk3OTkwODAxLAogICJpYXQiOiAxNDk3OTg4NzEyLAogICJpc3MiOiAiQXp1cmUgQ29udGFpbmVyIFJlZ2lzdHJ5IiwKICAiYXVkIjogImNvbnRvc29yZWdpc3RyeS5henVyZWNyLmlvIiwKICAidmVyc2lvbiI6ICIxLjAiLAogICJncmFudF90eXBlIjogInJlZnJlc2hfdG9rZW4iLAogICJ0ZW5hbnQiOiAiNDA5NTIwZDQtODEwMC00ZDFkLWFkNDctNzI0MzJkZGNjMTIwIiwKICAicGVybWlzc2lvbnMiOiB7CiAgICAiYWN0aW9ucyI6IFsKICAgICAgIioiCiAgICBdLAogICAgIm5vdEFjdGlvbnMiOiBbXQogIH0sCiAgInJvbGVzIjogW10KfQ=.",
			time.Unix(1497990801, 0),
			false,
		},
		{
			"test-error",
			".error.",
			time.Unix(1497990801, 0),
			true,
		},
		{
			"test-unmarshal-error",
			".ewogICJqdGkiOiAiMzY1ZTNiNWItODQ0ZS00YTIxLWEzOGMtNGQ4YWViZGQ2YTA2IiwKICAic3ViIjogInVzZXJAY29udG9zby5jb20iLAogICJuYmYiOiAxNDk3OTg4NzEyLAogICJleHAiOiAiMTQ5Nzk5MDgwMSIsCiAgImlhdCI6IDE0OTc5ODg3MTIsCiAgImlzcyI6ICJBenVyZSBDb250YWluZXIgUmVnaXN0cnkiLAogICJhdWQiOiAiY29udG9zb3JlZ2lzdHJ5LmF6dXJlY3IuaW8iLAogICJ2ZXJzaW9uIjogIjEuMCIsCiAgImdyYW50X3R5cGUiOiAicmVmcmVzaF90b2tlbiIsCiAgInRlbmFudCI6ICI0MDk1MjBkNC04MTAwLTRkMWQtYWQ0Ny03MjQzMmRkY2MxMjAiLAogICJwZXJtaXNzaW9ucyI6IHsKICAgICJhY3Rpb25zIjogWwogICAgICAiKiIKICAgIF0sCiAgICAibm90QWN0aW9ucyI6IFtdCiAgfSwKICAicm9sZXMiOiBbXQp9.",
			time.Unix(1497990801, 0),
			true,
		},
		{
			"test-length-error",
			".",
			time.Unix(1497990801, 0),
			true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			expire, err := getJWTExpireTime(test.token)
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expire, expire)
			}
		})
	}
}

func Test_authenticationPolicy_findServiceAndScope(t *testing.T) {
	resp1 := http.Response{}
	resp1.Header = http.Header{}
	resp1.Header.Set("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\",scope=\"registry:catalog:*\"")

	resp2 := http.Response{}
	resp2.Header = http.Header{}
	resp2.Header.Set("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\",scope=\"artifact-repository:repo:pull\"")

	for _, test := range []struct {
		acrScope   string
		acrService string
		resp       *http.Response
		err        bool
	}{
		{"registry:catalog:*", "contosoregistry.azurecr.io", &resp1, false},
		{"artifact-repository:repo:pull", "contosoregistry.azurecr.io", &resp2, false},
		{"error", "error", &http.Response{}, true},
	} {
		t.Run(fmt.Sprintf("%s-%s", test.acrService, test.acrScope), func(t *testing.T) {
			p := &authenticationPolicy{}
			err := p.findServiceAndScope(test.resp)
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.acrScope, p.acrScope)
				require.Equal(t, test.acrService, p.acrService)
			}
		})
	}
}

func Test_authenticationPolicy_getAccessToken_live(t *testing.T) {
	startRecording(t)
	endpoint, cred, options := getEndpointCredAndClientOptions(t)
	if reflect.ValueOf(options.Cloud).IsZero() {
		options.Cloud = cloud.AzurePublic
	}
	authClient := newAuthenticationClient(endpoint, &authenticationClientOptions{options})
	p := &authenticationPolicy{
		temporal.NewResource(acquire),
		cred,
		[]string{options.Cloud.Services[ServiceName].Audience + "/.default"},
		"registry:catalog:*",
		strings.TrimPrefix(endpoint, "https://"),
		authClient,
	}
	request, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://test.com")
	require.NoError(t, err)
	token, err := p.getAccessToken(request)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func Test_authenticationPolicy_getAccessToken_live_anonymous(t *testing.T) {
	startRecording(t)
	endpoint, _, options := getEndpointCredAndClientOptions(t)
	authClient := newAuthenticationClient(endpoint, &authenticationClientOptions{options})
	p := &authenticationPolicy{
		temporal.NewResource(acquire),
		nil,
		nil,
		"registry:catalog:*",
		strings.TrimPrefix(endpoint, "https://"),
		authClient,
	}
	request, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://test.com")
	require.NoError(t, err)
	token, err := p.getAccessToken(request)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func Test_authenticationPolicy_anonymousAccess(t *testing.T) {
	startRecording(t)
	endpoint, _, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, nil, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	ctx := context.Background()
	pager := client.NewListRepositoriesPager(nil)
	repositoryName := ""
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, page.Repositories.Names)
		if repositoryName == "" {
			repositoryName = *page.Repositories.Names[0]
		}
	}
	require.NotEmpty(t, repositoryName)
	_, err = client.UpdateRepositoryProperties(ctx, repositoryName, &ClientUpdateRepositoryPropertiesOptions{Value: &RepositoryWriteableProperties{CanDelete: to.Ptr(true)}})
	require.Error(t, err)
}
