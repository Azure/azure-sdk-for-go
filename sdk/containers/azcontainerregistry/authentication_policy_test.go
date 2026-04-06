// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
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

func Test_findServiceAndScope(t *testing.T) {
	resp1 := http.Response{}
	resp1.Header = http.Header{}
	resp1.Header.Set("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\",scope=\"registry:catalog:*\"")

	resp2 := http.Response{}
	resp2.Header = http.Header{}
	resp2.Header.Set("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\",scope=\"artifact-repository:repo:pull\"")

	resp3 := http.Response{}
	resp3.Header = http.Header{}
	resp3.Header.Set("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",scope=\"artifact-repository:repo:pull\"")

	resp4 := http.Response{}
	resp4.Header = http.Header{}
	resp4.Header.Set("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\"")

	for _, test := range []struct {
		acrScope   string
		acrService string
		resp       *http.Response
		err        bool
	}{
		{"registry:catalog:*", "contosoregistry.azurecr.io", &resp1, false},
		{"artifact-repository:repo:pull", "contosoregistry.azurecr.io", &resp2, false},
		{"error", "error", &http.Response{}, true},
		{"error2", "error", &resp3, true},
		{"error3", "error", &resp4, true},
	} {
		t.Run(fmt.Sprintf("%s-%s", test.acrService, test.acrScope), func(t *testing.T) {
			service, scope, err := findServiceAndScope(test.resp)
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.acrService, service)
				require.Equal(t, test.acrScope, scope)
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
	authClient, err := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{options})
	require.NoError(t, err)
	p := &authenticationPolicy{
		temporal.NewResource(acquireRefreshToken),
		atomic.Value{},
		cred,
		[]string{options.Cloud.Services[ServiceName].Audience + "/.default"},
		authClient,
	}
	request, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://test.com")
	require.NoError(t, err)
	token, err := p.getAccessToken(request, strings.TrimPrefix(endpoint, "https://"), "registry:catalog:*")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func Test_authenticationPolicy_getAccessToken_error(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("wrong response")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("{\"refresh_token\": \"test\"}")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("{\"refresh_token\": \".eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJuYmYiOjQ2NzA0MTEyMTIsImV4cCI6NDY3MDQyMjkxMiwiaWF0Ijo0NjcwNDExMjEyLCJpc3MiOiJBenVyZSBDb250YWluZXIgUmVnaXN0cnkiLCJhdWQiOiJhemFjcmxpdmV0ZXN0LmF6dXJlY3IuaW8iLCJ2ZXJzaW9uIjoiMS4wIiwicmlkIjoiMDAwMCIsImdyYW50X3R5cGUiOiJyZWZyZXNoX3Rva2VuIiwiYXBwaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJwZXJtaXNzaW9ucyI6eyJBY3Rpb25zIjpbInJlYWQiLCJ3cml0ZSIsImRlbGV0ZSIsImRlbGV0ZWQvcmVhZCIsImRlbGV0ZWQvcmVzdG9yZS9hY3Rpb24iXSwiTm90QWN0aW9ucyI6bnVsbH0sInJvbGVzIjpbXX0.\"}")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("wrong response")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("wrong response")))
	authClient, err := NewAuthenticationClient(srv.URL(), &AuthenticationClientOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	require.NoError(t, err)

	p := &authenticationPolicy{
		temporal.NewResource(acquireRefreshToken),
		atomic.Value{},
		&credential.Fake{},
		[]string{"test"},
		authClient,
	}
	request, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	_, err = p.getAccessToken(request, "service", "scope")
	require.Error(t, err)
	_, err = p.getAccessToken(request, "service", "scope")
	require.Error(t, err)
	_, err = p.getAccessToken(request, "service", "scope")
	require.Error(t, err)
	p.cred = nil
	_, err = p.getAccessToken(request, "service", "scope")
	require.Error(t, err)
}

func Test_authenticationPolicy_getAccessToken_live_anonymous(t *testing.T) {
	startRecording(t)
	endpoint, _, options := getEndpointCredAndClientOptions(t)
	authClient, err := NewAuthenticationClient(endpoint, &AuthenticationClientOptions{options})
	require.NoError(t, err)
	p := &authenticationPolicy{
		refreshTokenCache: temporal.NewResource(acquireRefreshToken),
		authClient:        authClient,
	}
	request, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://test.com")
	require.NoError(t, err)
	token, err := p.getAccessToken(request, strings.TrimPrefix(endpoint, "https://"), "registry:catalog:*")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func Test_authenticationPolicy_anonymousAccess(t *testing.T) {
	startRecording(t)
	endpoint, _, options := getEndpointCredAndClientOptions(t)
	client, err := NewClient(endpoint, nil, &ClientOptions{ClientOptions: options})
	require.NoError(t, err)
	pager := client.NewListRepositoriesPager(nil)
	for pager.More() {
		_, err = pager.NextPage(ctx)
		require.NoError(t, err)
	}
}

func Test_getChallengeRequest(t *testing.T) {
	oriReq, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://test.com")
	require.NoError(t, err)
	testBody := []byte("test")
	err = oriReq.SetBody(streaming.NopCloser(bytes.NewReader(testBody)), "text/plain")
	require.NoError(t, err)
	challengeReq, err := getChallengeRequest(*oriReq)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", len(testBody)), oriReq.Raw().Header.Get("Content-Length"))
	require.Equal(t, "", challengeReq.Raw().Header.Get("Content-Length"))
}

func Test_authenticationPolicy(t *testing.T) {
	srv, closeServer := mock.NewServer()
	defer closeServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized), mock.WithHeader("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\",scope=\"registry:catalog:*\""))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized), mock.WithHeader("WWW-Authenticate", "Bearer realm=\"https://contosoregistry.azurecr.io/oauth2/token\",service=\"contosoregistry.azurecr.io\",scope=\"registry:catalog:*\""))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("{\"refresh_token\": \".eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJuYmYiOjQ2NzA0MTEyMTIsImV4cCI6NDY3MDQyMjkxMiwiaWF0Ijo0NjcwNDExMjEyLCJpc3MiOiJBenVyZSBDb250YWluZXIgUmVnaXN0cnkiLCJhdWQiOiJhemFjcmxpdmV0ZXN0LmF6dXJlY3IuaW8iLCJ2ZXJzaW9uIjoiMS4wIiwicmlkIjoiMDAwMCIsImdyYW50X3R5cGUiOiJyZWZyZXNoX3Rva2VuIiwiYXBwaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJwZXJtaXNzaW9ucyI6eyJBY3Rpb25zIjpbInJlYWQiLCJ3cml0ZSIsImRlbGV0ZSIsImRlbGV0ZWQvcmVhZCIsImRlbGV0ZWQvcmVzdG9yZS9hY3Rpb24iXSwiTm90QWN0aW9ucyI6bnVsbH0sInJvbGVzIjpbXX0.\"}")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("{\"access_token\": \"test\"}")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	authClient, err := NewAuthenticationClient(srv.URL(), &AuthenticationClientOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	require.NoError(t, err)
	authPolicy := &authenticationPolicy{
		temporal.NewResource(acquireRefreshToken),
		atomic.Value{},
		&credential.Fake{},
		[]string{"test"},
		authClient,
	}
	pl := runtime.NewPipeline("testmodule", "v0.1.0", runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &policy.ClientOptions{Transport: srv})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	req.Raw().Header.Set(headerAuthorization, "test")
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err = pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
