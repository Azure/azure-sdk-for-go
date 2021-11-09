//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"fmt"
	"net/http"
	"testing"

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
	sampleURL := "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000"
	tenant := parseTenant(sampleURL)
	if *tenant != fakeTenant {
		t.Fatalf("tenant was not properly parsed, got %s, expected %s", *tenant, fakeTenant)
	}
}

func TestFindScopeAndTenant(t *testing.T) {
	p := KeyVaultChallengePolicy{}
	resp := http.Response{}
	resp.Header = http.Header{}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(authResource, fakeTenant, mhsmResource),
	)
	err := p.findScopeAndTenant(&resp)
	require.NoError(t, err)
	if *p.scope != mhsmScope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, mhsmScope)
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		fmt.Sprintf(authResourceScope, fakeTenant, resource, scope),
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
		fmt.Sprintf(resourceScopeAuth, mhsmResource, mhsmScope, fakeTenant),
	)
	err = p.findScopeAndTenant(&resp)
	require.NoError(t, err)
	if *p.scope != mhsmScope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, "https://vault.azure.net/.default")
	}
	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set(
		"WWW-Authenticate",
		"Bearer authorization=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\", unimportantkey=\"unimportantvalue\" resource=\"https://vault.azure.net/.default\"",
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
		"Bearer   authorization=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\",    unimportantkey=\"unimportantvalue\"   resource=\"https://vault.azure.net/.default\"    fakekey=\"fakevalue\"			",
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
