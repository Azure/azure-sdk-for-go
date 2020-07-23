// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"net/url"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

func Test_NonNilTokenCredentialOptsNilAuthorityHost(t *testing.T) {
	opts := &TokenCredentialOptions{Retry: &azcore.RetryOptions{MaxRetries: 6}}
	opts, err := opts.setDefaultValues()
	if err != nil {
		t.Fatalf("Received an error: %v", err)
	}
	if opts.AuthorityHost == nil {
		t.Fatalf("Did not set default authority host")
	}
}

func Test_SetEnvAuthorityHost(t *testing.T) {
	err := os.Setenv("AZURE_AUTHORITY_HOST", envHostString)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}

	opts := &TokenCredentialOptions{}
	opts, err = opts.setDefaultValues()
	if opts.AuthorityHost.String() != envHostString {
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

	customHost, err := url.Parse(customHostString)
	if err != nil {
		t.Fatalf("Received an error: %v", err)
	}

	opts := &TokenCredentialOptions{AuthorityHost: customHost}
	opts, err = opts.setDefaultValues()
	if opts.AuthorityHost.String() != customHostString {
		t.Fatalf("Unexpected error when get host from environment vairable: %v", err)
	}

	// Unset that host environment vairable to avoid other tests failed.
	err = os.Unsetenv("AZURE_AUTHORITY_HOST")
	if err != nil {
		t.Fatalf("Unexpected error when unset environment vairable: %v", err)
	}
}

func Test_DefaultAuthorityHost(t *testing.T) {
	opts := &TokenCredentialOptions{}
	opts, err := opts.setDefaultValues()
	if opts.AuthorityHost.String() != AzurePublicCloud {
		t.Fatalf("Unexpected error when set default AuthorityHost: %v", err)
	}
}

func Test_AzureGermanyAuthorityHost(t *testing.T) {
	opts := &TokenCredentialOptions{}
	opts, err := opts.setDefaultValues()
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(AzureGermany)
	if err != nil {
		t.Fatal(err)
	}
	opts.AuthorityHost = u
	if opts.AuthorityHost.String() != AzureGermany {
		t.Fatalf("Did not retrieve expected authority host string")
	}
}

func Test_AzureChinaAuthorityHost(t *testing.T) {
	opts := &TokenCredentialOptions{}
	opts, err := opts.setDefaultValues()
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(AzureChina)
	if err != nil {
		t.Fatal(err)
	}
	opts.AuthorityHost = u
	if opts.AuthorityHost.String() != AzureChina {
		t.Fatalf("Did not retrieve expected authority host string")
	}
}

func Test_AzureGovernmentAuthorityHost(t *testing.T) {
	opts := &TokenCredentialOptions{}
	opts, err := opts.setDefaultValues()
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(AzureGovernment)
	if err != nil {
		t.Fatal(err)
	}
	opts.AuthorityHost = u
	if opts.AuthorityHost.String() != AzureGovernment {
		t.Fatalf("Did not retrieve expected authority host string")
	}
}
