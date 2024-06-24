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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

func TestClientAssertionCredential(t *testing.T) {
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
		ClientOptions: azcore.ClientOptions{Transport: &mockSTS{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), key, true)
	_, err = cred.GetToken(ctx, testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	// silent authentication should now succeed
	_, err = cred.GetToken(ctx, testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestClientAssertionCredentialCallbackError(t *testing.T) {
	expectedError := errors.New("it didn't work")
	getAssertion := func(c context.Context) (string, error) { return "", expectedError }
	cred, err := NewClientAssertionCredential("tenant", "clientID", getAssertion, &ClientAssertionCredentialOptions{
		ClientOptions: azcore.ClientOptions{Transport: &mockSTS{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil || !strings.Contains(err.Error(), expectedError.Error()) {
		t.Fatalf(`unexpected error: "%v"`, err)
	}
}

func TestClientAssertionCredentialNilCallback(t *testing.T) {
	_, err := NewClientAssertionCredential(fakeTenantID, fakeClientID, nil, nil)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientAssertionCredential_Live(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
	}
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
					return assertion(certs[0], key)
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
