// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"os"
	"testing"
)

const (
	envHostString    = "https://mock.com/"
	customHostString = "https://custommock.com/"
)

// Set AZURE_AUTHORITY_HOST for the duration of a test. Restore its prior value
// after the test completes. Prevents tests which set the variable from breaking live
// tests in sovereign clouds. Obviated by 1.17's T.Setenv
func setEnvAuthorityHost(host string, t *testing.T) {
	originalHost := os.Getenv("AZURE_AUTHORITY_HOST")
	err := os.Setenv("AZURE_AUTHORITY_HOST", host)
	if err != nil {
		t.Fatalf("Unexpected error setting AZURE_AUTHORITY_HOST: %v", err)
	}
	t.Cleanup(func() {
		err = os.Setenv("AZURE_AUTHORITY_HOST", originalHost)
		if err != nil {
			t.Fatalf("Unexpected error resetting AZURE_AUTHORITY_HOST: %v", err)
		}
	})
}

func Test_SetEnvAuthorityHost(t *testing.T) {
	setEnvAuthorityHost(envHostString, t)
	authorityHost, err := setAuthorityHost("")
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != envHostString {
		t.Fatalf("Unexpected error when get host from environment variable: %v", err)
	}
}

func Test_CustomAuthorityHost(t *testing.T) {
	setEnvAuthorityHost(envHostString, t)
	authorityHost, err := setAuthorityHost(customHostString)
	if err != nil {
		t.Fatal(err)
	}
	// ensure env var doesn't override explicit value
	if authorityHost != customHostString {
		t.Fatalf("Unexpected host when get host from environment variable: %v", authorityHost)
	}
}

func Test_DefaultAuthorityHost(t *testing.T) {
	setEnvAuthorityHost("", t)
	authorityHost, err := setAuthorityHost("")
	if err != nil {
		t.Fatal(err)
	}
	if authorityHost != string(AzurePublicCloud) {
		t.Fatalf("Unexpected host when set default AuthorityHost: %v", authorityHost)
	}
}

func Test_NonHTTPSAuthorityHost(t *testing.T) {
	setEnvAuthorityHost("", t)
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
