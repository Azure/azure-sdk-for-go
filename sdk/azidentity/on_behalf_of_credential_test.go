//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestOnBehalfOfCredential(t *testing.T) {
	realGetClient := getConfidentialClient
	t.Cleanup(func() { getConfidentialClient = realGetClient })
	expectedAssertion := "user-assertion"
	for _, test := range []struct {
		ctor    func(policy.Transporter) (*OnBehalfOfCredential, error)
		name    string
		sendX5C bool
	}{
		{
			ctor: func(tp policy.Transporter) (*OnBehalfOfCredential, error) {
				certs, key := allCertTests[0].certs, allCertTests[0].key
				o := OnBehalfOfCredentialOptions{ClientOptions: policy.ClientOptions{Transport: tp}}
				return NewOnBehalfOfCredentialWithCertificate(fakeTenantID, fakeClientID, expectedAssertion, certs, key, &o)
			},
			name: "certificate",
		},
		{
			ctor: func(tp policy.Transporter) (*OnBehalfOfCredential, error) {
				certs, key := allCertTests[0].certs, allCertTests[0].key
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
			srv := mockSTS{tokenRequestCallback: func(r *http.Request) {
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
			}}
			cred, err := test.ctor(&srv)
			if err != nil {
				t.Fatal(err)
			}
			tk, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
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
