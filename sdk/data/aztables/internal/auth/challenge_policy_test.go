//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package auth

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

var fakeTenant = "00000000-0000-0000-0000-000000000000"

var scope = "https://myaccount.table.core.windows.net/.default"
var resource = "https://table.core.windows.net"

var authScope = "Bearer authorization_uri=\"https://login.microsoftonline.com/%s\", scope=\"%s\""
var authResource = "Bearer authorization_uri=\"https://login.microsoftonline.com/%s\", resource_id=\"%s\""
var authResourceScope = "Bearer authorization_uri=\"https://login.microsoftonline.com/%s\", resource_id=\"%s\" scope=\"%s\""
var resourceScopeAuth = "Bearer resource_id=\"%s\" scope=\"%s\", authorization_uri=\"https://login.microsoftonline.com/%s\""

func TestParseTenantID(t *testing.T) {
	sampleURL := "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000"
	tenant := parseTenant(sampleURL)
	if *tenant != fakeTenant {
		t.Fatalf("tenant was not properly parsed, got %s, expected %s", *tenant, fakeTenant)
	}
}

func TestFindScopeAndTenant(t *testing.T) {
	p := StorageChallengePolicy{}
	resp := http.Response{}
	resp.Header = http.Header{}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(authResourceScope, fakeTenant, resource, scope),
	)
	err := p.findScopeAndTenant(&resp)
	require.NoError(t, err)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(authScope, fakeTenant, scope),
	)
	err = p.findScopeAndTenant(&resp)
	require.NoError(t, err)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		"Bearer authorization_uri=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\", unimportantkey=\"unimportantvalue\" resource_id=\"https://vault.azure.net/.default\"",
	)
	err = p.findScopeAndTenant(&resp)
	require.NoError(t, err)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		"Bearer   authorization_uri=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\",    unimportantkey=\"unimportantvalue\"   resource_id=\"https://vault.azure.net/.default\"    fakekey=\"fakevalue\"			",
	)
	err = p.findScopeAndTenant(&resp)
	require.NoError(t, err)
	if *p.scope != "https://vault.azure.net/.default" {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, "https://vault.azure.net/.default")
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}
}

type fakeCredential struct{}

func (f fakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{}, nil
}

func TestDo(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	fakeCred := fakeCredential{}
	p := NewStorageChallengePolicy(fakeCred)

	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized), mock.WithHeader("WWW-Authenticate", fmt.Sprintf(authResourceScope, fakeTenant, resource, scope)))
	pl := runtime.NewPipeline("test", "test", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport:       srv,
		PerCallPolicies: []policy.Policy{p},
	})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusUnauthorized)

	require.Equal(t, *p.scope, scope)
	require.Equal(t, *p.tenantID, fakeTenant)
}
