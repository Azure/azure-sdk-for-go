//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func TestInteractiveBrowserCredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewInteractiveBrowserCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client.noCAE = fakePublicClient{
		ar: public.AuthResult{
			AccessToken: tokenValue,
			ExpiresOn:   time.Now().Add(1 * time.Hour),
		},
	}
	tk, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if tk.Token != tokenValue {
		t.Fatal("Received unexpected token")
	}
}

// instanceDiscoveryPolicy fails the test when the client requests instance metadata
type instanceDiscoveryPolicy struct {
	t *testing.T
}

func (p *instanceDiscoveryPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if strings.Contains(req.Raw().URL.Path, "discovery/instance") {
		p.t.Fatal("client requested instance metadata")
	}
	return req.Next()
}

func TestInteractiveBrowserCredential_Live(t *testing.T) {
	if !runManualTests {
		t.Skipf("set %s to run this test", azidentityRunManualTests)
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
		fmt.Printf("\t%s: consider this test passing when %q appears in the login prompt\n", t.Name(), upn)
		cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{LoginHint: upn})
		if err != nil {
			t.Fatal(err)
		}
		testGetTokenSuccess(t, cred)
	})
	t.Run("RedirectURL", func(t *testing.T) {
		url := "http://localhost:8180"
		fmt.Printf("\t%s: consider this test passing when Microsoft Entra redirects to %s\n", t.Name(), url)
		cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{RedirectURL: url})
		if err != nil {
			t.Fatal(err)
		}
		testGetTokenSuccess(t, cred)
	})

	t.Run("instance discovery disabled", func(t *testing.T) {
		cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{
			ClientOptions: policy.ClientOptions{
				PerCallPolicies: []policy.Policy{
					&instanceDiscoveryPolicy{t},
				}},
			DisableInstanceDiscovery: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		testGetTokenSuccess(t, cred)
	})
}

func TestInteractiveBrowserCredentialADFS_Live(t *testing.T) {
	if !runManualTests {
		t.Skipf("set %s to run this test", azidentityRunManualTests)
	}
	if adfsLiveUser.clientID == fakeClientID {
		t.Skip("set ADFS_IDENTITY_TEST_CLIENT_ID environment variables to run this test live")
	}
	url := adfsLiveSP.redirectURL

	cloudConfig := cloud.Configuration{ActiveDirectoryAuthorityHost: adfsAuthority}

	clientOptions := policy.ClientOptions{Cloud: cloudConfig}

	cred, err := NewInteractiveBrowserCredential(&InteractiveBrowserCredentialOptions{
		ClientOptions:            clientOptions,
		ClientID:                 adfsLiveUser.clientID,
		DisableInstanceDiscovery: true,
		RedirectURL:              url,
		TenantID:                 "adfs",
	})
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}
