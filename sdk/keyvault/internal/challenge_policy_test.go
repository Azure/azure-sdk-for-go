//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

var fakeTenant = "00000000-0000-0000-0000-000000000000"
var mhsmScope = "https://managedhsm.azure.net/.default"
var mhsmResource = "https://managedhsm.azure.net"

var scope = "https://vault.azure.net/.default"
var resource = "https://vault.azure.net"

var authScope = "Bearer authorization=\"https://login.microsoftonline.com/%s\", scope=\"%s\""
var authResource = "Bearer authorization=\"https://login.microsoftonline.com/%s\", resource=\"%s\""
var authResourceScope = "Bearer authorization=\"https://login.microsoftonline.com/%s\", resource=\"%s\" scope=\"%s\""
var resourceScopeAuth = "Bearer resource=\"%s\" scope=\"%s\", authorization=\"https://login.microsoftonline.com/%s\""

func TestParseTenantID(t *testing.T) {
	tenant := parseTenant("")
	require.NotNil(t, tenant)
	require.Empty(t, *tenant)

	sampleURL := "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000"
	tenant = parseTenant(sampleURL)
	require.NotNil(t, tenant)
	require.Equal(t, fakeTenant, *tenant, "tenant was not properly parsed, got %s, expected %s", *tenant, fakeTenant)
}

type credentialFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (cf credentialFunc) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return cf(ctx, options)
}

func TestFindScopeAndTenant(t *testing.T) {
	p := NewKeyVaultChallengePolicy(credentialFunc(func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
		return azcore.AccessToken{
			Token:     "fake_token",
			ExpiresOn: time.Now().Add(time.Hour),
		}, nil
	}))
	resp := http.Response{}
	resp.Header = http.Header{}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(authResource, fakeTenant, mhsmResource),
	)
	req, err := http.NewRequest("GET", "https://42.managedhsm.azure.net", nil)
	require.NoError(t, err)
	resp.Request = req
	err = p.findScopeAndTenant(&resp, req)
	require.NoError(t, err)
	if *p.scope != mhsmScope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, mhsmScope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(resourceScopeAuth, mhsmResource, mhsmScope, fakeTenant),
	)
	err = p.findScopeAndTenant(&resp, req)
	require.NoError(t, err)
	if *p.scope != mhsmScope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, "https://vault.azure.net/.default")
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(authResourceScope, fakeTenant, resource, scope),
	)
	req, err = http.NewRequest("GET", "https://42.vault.azure.net", nil)
	require.NoError(t, err)
	resp.Request = req
	err = p.findScopeAndTenant(&resp, req)
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
	err = p.findScopeAndTenant(&resp, req)
	require.NoError(t, err)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		"Bearer authorization=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\", unimportantkey=\"unimportantvalue\" resource=\"https://vault.azure.net/.default\"",
	)
	err = p.findScopeAndTenant(&resp, req)
	require.NoError(t, err)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		"Bearer   authorization=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\",    unimportantkey=\"unimportantvalue\"   resource=\"https://vault.azure.net/.default\"    fakekey=\"fakevalue\"			",
	)
	err = p.findScopeAndTenant(&resp, req)
	require.NoError(t, err)
	if *p.scope != "https://vault.azure.net/.default" {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, "https://vault.azure.net/.default")
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set("WWW-Authenticate", "this is an invalid value")
	err = p.findScopeAndTenant(&resp, req)
	var challengeError *challengePolicyError
	require.ErrorAs(t, err, &challengeError)

	resp.Header = http.Header{}
	err = p.findScopeAndTenant(&resp, req)
	require.ErrorAs(t, err, &challengeError)
}

func TestResourceVerification(t *testing.T) {
	for _, test := range []struct {
		challenge, resource string
		err                 bool
	}{
		{authResource, "https://vault.azure.net", false},
		{authScope, "https://vault.azure.net/.default", false},

		// error cases: resource/scope doesn't match the requested vault's host (vault.azure.net)
		{authResource, "https://vault.azure.cn", true},
		{authResource, "https://myvault.azure.net", true},
		{authScope, "https://vault.azure.cn/.default", true},
		{authScope, "https://myvault.azure.net/.default", true},
	} {
		t.Run(test.resource, func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(
				mock.WithHeader("WWW-Authenticate", fmt.Sprintf(test.challenge, fakeTenant, test.resource)),
				mock.WithStatusCode(401),
			)
			srv.AppendResponse()

			cred := credentialFunc(func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
				return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
			})
			p := NewKeyVaultChallengePolicy(cred)
			pl := runtime.NewPipeline("", "",
				runtime.PipelineOptions{PerCall: []policy.Policy{p}},
				&policy.ClientOptions{Transport: srv},
			)
			req, err := runtime.NewRequest(context.Background(), "GET", "https://42.vault.azure.net")
			require.NoError(t, err)
			_, err = pl.Do(req)
			if test.err {
				expected := fmt.Sprintf(`the challenge resource "%s" doesn't match the requested domain`, test.resource)
				require.EqualError(t, err, expected)
				if _, ok := err.(*challengePolicyError); !ok {
					t.Fatalf("unexpected error type %T", err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TODO: add test coverage for the following
//   func (k *KeyVaultChallengePolicy) Do
//   func (k KeyVaultChallengePolicy) getChallengeRequest
//   func acquire
