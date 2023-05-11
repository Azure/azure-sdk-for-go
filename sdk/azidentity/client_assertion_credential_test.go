//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestClientAssertionCredential(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithBody(instanceDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(tenantDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(accessTokenRespSuccess))

	key := struct{}{}
	calls := 0
	getAssertion := func(c context.Context) (string, error) {
		if v := c.Value(key); v == nil || !v.(bool) {
			t.Fatal("unexpected context in getAssertion")
		}
		calls++
		return "assertion", nil
	}
	cred, err := NewClientAssertionCredential("tenant", "clientID", getAssertion, &ClientAssertionCredentialOptions{
		ClientOptions: azcore.ClientOptions{Transport: srv},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), key, true)
	_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	// silent authentication should now succeed
	_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestClientAssertionCredentialCallbackError(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithBody(instanceDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(tenantDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(accessTokenRespSuccess))

	expectedError := errors.New("it didn't work")
	getAssertion := func(c context.Context) (string, error) { return "", expectedError }
	cred, err := NewClientAssertionCredential("tenant", "clientID", getAssertion, &ClientAssertionCredentialOptions{
		ClientOptions: azcore.ClientOptions{Transport: srv},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil || !strings.Contains(err.Error(), expectedError.Error()) {
		t.Fatalf(`unexpected error: "%v"`, err)
	}
}

func TestClientAssertionCredential_Live(t *testing.T) {
	data, err := os.ReadFile(liveSP.pemPath)
	if err != nil {
		t.Fatalf(`failed to read cert: %v`, err)
	}
	certs, key, err := ParseCertificates(data, nil)
	if err != nil {
		t.Fatalf(`failed to parse cert: %v`, err)
	}
	for _, d := range []bool{true, false} {
		name := "default options"
		if d {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			o, stop := initRecording(t)
			defer stop()
			cred, err := NewClientAssertionCredential(liveSP.tenantID, liveSP.clientID,
				func(context.Context) (string, error) {
					return getAssertion(certs[0], key)
				},
				&ClientAssertionCredentialOptions{ClientOptions: o, DisableInstanceDiscovery: d},
			)
			if err != nil {
				t.Fatal(err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}
