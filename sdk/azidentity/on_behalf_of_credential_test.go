//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestOnBehalfOfCredential(t *testing.T) {
	expectedAssertion := "user-assertion"
	certs, key := allCertTests[0].certs, allCertTests[0].key
	for _, test := range []struct {
		ctor    func(policy.Transporter) (*OnBehalfOfCredential, error)
		name    string
		sendX5C bool
	}{
		{
			ctor: func(tp policy.Transporter) (*OnBehalfOfCredential, error) {
				o := OnBehalfOfCredentialOptions{ClientOptions: policy.ClientOptions{Transport: tp}}
				return NewOnBehalfOfCredentialWithCertificate(fakeTenantID, fakeClientID, expectedAssertion, certs, key, &o)
			},
			name: "certificate",
		},
		{
			ctor: func(tp policy.Transporter) (*OnBehalfOfCredential, error) {
				o := OnBehalfOfCredentialOptions{ClientOptions: policy.ClientOptions{Transport: tp}, SendCertificateChain: true}
				return NewOnBehalfOfCredentialWithCertificate(fakeTenantID, fakeClientID, expectedAssertion, certs, key, &o)
			},
			name:    "SNI",
			sendX5C: true,
		},
		{
			ctor: func(tp policy.Transporter) (*OnBehalfOfCredential, error) {
				o := OnBehalfOfCredentialOptions{ClientOptions: policy.ClientOptions{Transport: tp}}
				return NewOnBehalfOfCredentialWithSecret(fakeTenantID, fakeClientID, expectedAssertion, "secret", &o)
			},
			name: "secret",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			key := struct{}{}
			ctx := context.WithValue(context.Background(), key, true)
			srv := mockSTS{tokenRequestCallback: func(r *http.Request) *http.Response {
				if c := r.Context(); c == nil {
					t.Fatal("AcquireTokenOnBehalfOf received no Context")
				} else if v := c.Value(key); v == nil || !v.(bool) {
					t.Fatal("AcquireTokenOnBehalfOf received unexpected Context")
				}
				if err := r.ParseForm(); err != nil {
					t.Fatal(err)
				}
				if scope := r.FormValue("scope"); !strings.Contains(scope, liveTestScope) {
					t.Errorf(`unexpected scopes "%v"`, scope)
				}
				if assertion := r.FormValue("assertion"); assertion != expectedAssertion {
					t.Errorf(`unexpected assertion "%s"`, assertion)
				}
				if test.sendX5C {
					validateX5C(t, certs)(r)
				}
				return nil
			}}
			cred, err := test.ctor(&srv)
			if err != nil {
				t.Fatal(err)
			}
			tk, err := cred.GetToken(ctx, testTRO)
			if err != nil {
				t.Fatal(err)
			}
			if tk.Token != tokenValue {
				t.Errorf(`unexpected token "%s"`, tk.Token)
			}
			if tk.ExpiresOn.Before(time.Now()) {
				t.Error("GetToken returned an invalid expiration time")
			}
			if tk.ExpiresOn.Location() != time.UTC {
				t.Error("ExpiresOn isn't UTC")
			}
		})
	}
}

func TestOnBehalfOfCredential_Error(t *testing.T) {
	// GetToken shouldn't send a second token request after the first fails
	tokenReqs := 0
	cred, err := NewOnBehalfOfCredentialWithSecret("tenant", "clientID", "assertion", "secret", &OnBehalfOfCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Transport: &mockSTS{
				tokenRequestCallback: func(*http.Request) *http.Response {
					tokenReqs++
					return &http.Response{Body: io.NopCloser(strings.NewReader("")), StatusCode: 400}
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatal("expected an error")
	}
	if tokenReqs != 1 {
		t.Fatalf("expected 1 token request, got %d", tokenReqs)
	}
}
