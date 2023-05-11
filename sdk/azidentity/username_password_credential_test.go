//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

func TestUsernamePasswordCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewUsernamePasswordCredential(badTenantID, fakeClientID, "username", "password", nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestUsernamePasswordCredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewUsernamePasswordCredential(fakeTenantID, fakeClientID, "username", "password", nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client = fakePublicClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestUsernamePasswordCredential_Live(t *testing.T) {
	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			o, stop := initRecording(t)
			defer stop()
			opts := UsernamePasswordCredentialOptions{ClientOptions: o, DisableInstanceDiscovery: disabledID}
			cred, err := NewUsernamePasswordCredential(liveUser.tenantID, developerSignOnClientID, liveUser.username, liveUser.password, &opts)
			if err != nil {
				t.Fatalf("Unable to create credential. Received: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestUsernamePasswordCredentialADFS_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if adfsLiveUser.clientID == "" || adfsLiveUser.username == "" || adfsLiveUser.password == "" {
			t.Skip("set ADFS_IDENTITY_TEST_* environment variables to run this test live")
		}
	}
	o, stop := initRecording(t)
	o.Cloud.ActiveDirectoryAuthorityHost = adfsAuthority
	defer stop()
	opts := UsernamePasswordCredentialOptions{ClientOptions: o, DisableInstanceDiscovery: true}
	cred, err := NewUsernamePasswordCredential("adfs", adfsLiveUser.clientID, adfsLiveUser.username, adfsLiveUser.password, &opts)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}

func TestUsernamePasswordCredential_InvalidPasswordLive(t *testing.T) {
	o, stop := initRecording(t)
	defer stop()
	opts := UsernamePasswordCredentialOptions{ClientOptions: o}
	cred, err := NewUsernamePasswordCredential(liveUser.tenantID, developerSignOnClientID, liveUser.username, "invalid password", &opts)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if !reflect.ValueOf(tk).IsZero() {
		t.Fatal("expected a zero value AccessToken")
	}
	if e, ok := err.(*AuthenticationFailedError); ok {
		if e.RawResponse == nil {
			t.Fatal("expected a non-nil RawResponse")
		}
	} else {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	if !strings.HasPrefix(err.Error(), credNameUserPassword) {
		t.Fatal("missing credential type prefix")
	}
}
