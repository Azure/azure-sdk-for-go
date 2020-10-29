// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"net/url"
	"os"
	"testing"
)

const (
	envHostString    = "https://mock.com/"
	customHostString = "https://custommock.com/"
)

func Test_AuthorityHost_Parse(t *testing.T) {
	_, err := url.Parse(AzurePublicCloud)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

func Test_SetEnvAuthorityHost(t *testing.T) {
	err := os.Setenv("AZURE_AUTHORITY_HOST", envHostString)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	authorityHost, err := setAuthorityHost("")
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != envHostString {
		t.Fatalf("Unexpected error when get host from environment vairable: %v", err)
	}

	// Unset that host environment vairable to avoid other tests failed.
	err = os.Unsetenv("AZURE_AUTHORITY_HOST")
	if err != nil {
		t.Fatalf("Unexpected error when unset environment vairable: %v", err)
	}
}

func Test_CustomAuthorityHost(t *testing.T) {
	err := os.Setenv("AZURE_AUTHORITY_HOST", envHostString)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}

	authorityHost, err := setAuthorityHost(customHostString)
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != customHostString {
		t.Fatalf("Unexpected error when get host from environment vairable: %v", err)
	}

	// Unset that host environment vairable to avoid other tests failed.
	err = os.Unsetenv("AZURE_AUTHORITY_HOST")
	if err != nil {
		t.Fatalf("Unexpected error when unset environment vairable: %v", err)
	}
}

func Test_DefaultAuthorityHost(t *testing.T) {
	authorityHost, err := setAuthorityHost("")
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != AzurePublicCloud {
		t.Fatalf("Unexpected error when set default AuthorityHost: %v", err)
	}
}

func Test_AzureGermanyAuthorityHost(t *testing.T) {
	authorityHost, err := setAuthorityHost(AzureGermany)
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != AzureGermany {
		t.Fatalf("Did not retrieve expected authority host string")
	}
}

func Test_AzureChinaAuthorityHost(t *testing.T) {
	authorityHost, err := setAuthorityHost(AzureChina)
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != AzureChina {
		t.Fatalf("Did not retrieve expected authority host string")
	}
}

func Test_AzureGovernmentAuthorityHost(t *testing.T) {
	authorityHost, err := setAuthorityHost(AzureGovernment)
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != AzureGovernment {
		t.Fatalf("Did not retrieve expected authority host string")
	}
}

func Test_NonHTTPSAuthorityHost(t *testing.T) {
	authorityHost, err := setAuthorityHost("http://foo.com")
	if err == nil {
		t.Fatal("Expected an error but did not receive one.")
	}
	if authorityHost != "" {
		t.Fatalf("Unexpected value in authority host string: %s", authorityHost)
	}
}

func Test_ValidTenantIDFalse(t *testing.T) {
	if validTenantID("bad@tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad/tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad(tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad)tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
	if validTenantID("bad:tenant") {
		t.Fatal("Expected to receive false, but received true")
	}
}

func Test_ValidTenantIDTrue(t *testing.T) {
	if !validTenantID("goodtenant") {
		t.Fatal("Expected to receive true, but received false")
	}
	if !validTenantID("good-tenant") {
		t.Fatal("Expected to receive true, but received false")
	}
	if !validTenantID("good.tenant") {
		t.Fatal("Expected to receive true, but received false")
	}
}
