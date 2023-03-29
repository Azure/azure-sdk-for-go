//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const fakeSecret = "secret"

func TestClientSecretCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewClientSecretCredential(badTenantID, fakeClientID, fakeSecret, nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestClientSecretCredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewClientSecretCredential(fakeTenantID, fakeClientID, fakeSecret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client = fakeConfidentialClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
}

func TestClientSecretCredential_Live(t *testing.T) {
	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			opts, stop := initRecording(t)
			defer stop()
			o := ClientSecretCredentialOptions{ClientOptions: opts, DisableInstanceDiscovery: disabledID}
			cred, err := NewClientSecretCredential(liveSP.tenantID, liveSP.clientID, liveSP.secret, &o)
			if err != nil {
				t.Fatalf("failed to construct credential: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestClientSecretCredential_Live_IPv6(t *testing.T) {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip()
	}

	ipv6Client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}
				return dialer.DialContext(ctx, "tcp6", addr)
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			o := ClientSecretCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					Transport: ipv6Client,
				},
				DisableInstanceDiscovery: disabledID}
			cred, err := NewClientSecretCredential(liveSP.tenantID, liveSP.clientID, liveSP.secret, &o)
			if err != nil {
				t.Fatalf("failed to construct credential: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestClientSecretCredentialADFS_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if adfsLiveSP.clientID == "" || adfsLiveSP.secret == "" || adfsScope == "" {
			t.Skip("set ADFS_SP_* environment variables to run this test live")
		}
	}
	opts, stop := initRecording(t)
	defer stop()
	opts.Cloud.ActiveDirectoryAuthorityHost = adfsAuthority
	o := ClientSecretCredentialOptions{ClientOptions: opts, DisableInstanceDiscovery: true}
	cred, err := NewClientSecretCredential("adfs", adfsLiveSP.clientID, adfsLiveSP.secret, &o)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}

func TestClientSecretCredential_InvalidSecretLive(t *testing.T) {
	opts, stop := initRecording(t)
	defer stop()
	o := ClientSecretCredentialOptions{ClientOptions: opts}
	cred, err := NewClientSecretCredential(liveSP.tenantID, liveSP.clientID, "invalid secret", &o)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
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
	if !strings.HasPrefix(err.Error(), credNameSecret) {
		t.Fatal("missing credential type prefix")
	}
}
