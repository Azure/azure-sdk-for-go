//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"net/http"
	"testing"
)

var fakeTenant = "00000000-0000-0000-0000-000000000000"
var scope = "https://managedhsm.azure.net/.default"

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
	resp.Header.Set("WWW-Authenticate", "Bearer authorization=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\", resource=\"https://managedhsm.azure.net\"")

	p.findScopeAndTenant(&resp)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}

	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}

	resp.Header.Set("WWW-Authenticate", "Bearer authorization=\"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000\", resource=\"https://managedhsm.azure.net\" scope=\"https://vault.azure.net/.default\"")

	p.findScopeAndTenant(&resp)
	if *p.scope != scope {
		t.Fatalf("scope was not properly parsed, got %s, expected %s", *p.scope, scope)
	}

	if *p.tenantID != fakeTenant {
		t.Fatalf("tenant ID was not properly parsed, got %s, expected %s", *p.tenantID, fakeTenant)
	}
}
