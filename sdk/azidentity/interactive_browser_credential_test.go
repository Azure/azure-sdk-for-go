//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

var _, runManualBrowserTests = os.LookupEnv("AZIDENTITY_RUN_MANUAL_BROWSER_TESTS")

func TestInteractiveBrowserCredential_InvalidTenantID(t *testing.T) {
	options := InteractiveBrowserCredentialOptions{}
	options.TenantID = badTenantID
	cred, err := NewInteractiveBrowserCredential(&options)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestInteractiveBrowserCredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewInteractiveBrowserCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client = fakePublicClient{
		ar: public.AuthResult{
			AccessToken: tokenValue,
			ExpiresOn:   time.Now().Add(1 * time.Hour),
		},
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if tk.Token != tokenValue {
		t.Fatal("Received unexpected token")
	}
}

func TestInteractiveBrowserCredential_CreateWithNilOptions(t *testing.T) {
	cred, err := NewInteractiveBrowserCredential(nil)
	if err != nil {
		t.Fatalf("Failed to create interactive browser credential: %v", err)
	}
	if cred.options.ClientID != developerSignOnClientID {
		t.Fatalf("Wrong clientID set. Expected: %s, Received: %s", developerSignOnClientID, cred.options.ClientID)
	}
	if cred.options.TenantID != organizationsTenantID {
		t.Fatalf("Wrong tenantID set. Expected: %s, Received: %s", organizationsTenantID, cred.options.TenantID)
	}
}

func TestInteractiveBrowserCredential_Live(t *testing.T) {
	if !runManualBrowserTests {
		t.Skip("set AZIDENTITY_RUN_MANUAL_BROWSER_TESTS to run this test")
	}
	t.Run("defaults", func(t *testing.T) {
		cred, err := NewInteractiveBrowserCredential(nil)
		if err != nil {
			t.Fatal(err)
		}
		testGetTokenSuccess(t, cred)
	})
	t.Run("LoginHint", func(t *testing.T) {
		upn := "test@pass"
		t.Logf("consider this test passing when %q appears in the login prompt", upn)
		cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{LoginHint: upn})
		if err != nil {
			t.Fatal(err)
		}
		testGetTokenSuccess(t, cred)
	})
	t.Run("RedirectURL", func(t *testing.T) {
		url := "http://localhost:8180"
		t.Logf("consider this test passing when AAD redirects to %s", url)
		cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{RedirectURL: url})
		if err != nil {
			t.Fatal(err)
		}
		testGetTokenSuccess(t, cred)
	})

	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{DisableInstanceDiscovery: disabledID})
			if err != nil {
				t.Fatal(err)
			}

			testGetTokenSuccess(t, cred)
		})
	}
}

func TestInteractiveBrowserCredentialADFS_Live(t *testing.T) {
	if !runManualBrowserTests {
		t.Skip("set AZIDENTITY_RUN_MANUAL_BROWSER_TESTS to run this test")
	}
	if adfsLiveUser.clientID == fakeClientID {
		t.Skip("set ADFS_IDENTITY_TEST_CLIENT_ID environment variables to run this test live")
	}
	//Redirect URL is necessary
	url := adfsLiveSP.redirectURL
	vars := map[string]string{
		azureAuthorityHost: adfsAuthority,
	}
	setEnvironmentVariables(t, vars)
	cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{ClientID: adfsLiveUser.clientID, TenantID: "adfs", RedirectURL: url, DisableInstanceDiscovery: true})
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}
