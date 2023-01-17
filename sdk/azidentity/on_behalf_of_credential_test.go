//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

func TestOnBehalfOfCredential(t *testing.T) {
	realGetClient := getConfidentialClient
	t.Cleanup(func() { getConfidentialClient = realGetClient })
	expectedAssertion := "user-assertion"
	for _, test := range []struct {
		ctor    func() (*OnBehalfOfCredential, error)
		name    string
		sendX5C bool
	}{
		{
			ctor: func() (*OnBehalfOfCredential, error) {
				certs, key := allCertTests[0].certs, allCertTests[0].key
				return NewOnBehalfOfCredentialFromCertificate(fakeTenantID, fakeClientID, expectedAssertion, certs, key, nil)
			},
			name: "certificate",
		},
		{
			ctor: func() (*OnBehalfOfCredential, error) {
				certs, key := allCertTests[0].certs, allCertTests[0].key
				return NewOnBehalfOfCredentialFromCertificate(fakeTenantID, fakeClientID, expectedAssertion, certs, key, &OnBehalfOfCredentialOptions{SendCertificateChain: true})
			},
			name:    "certificate_SNI",
			sendX5C: true,
		},
		{
			ctor: func() (*OnBehalfOfCredential, error) {
				return NewOnBehalfOfCredentialFromSecret(fakeTenantID, fakeClientID, expectedAssertion, "secret", nil)
			},
			name: "secret",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			called := false
			key := struct{}{}
			ctx := context.WithValue(context.Background(), key, true)
			fake := fakeConfidentialClient{
				ar: confidential.AuthResult{AccessToken: tokenValue, ExpiresOn: time.Now().Add(time.Hour)},
				oboCallback: func(c context.Context, assertion string, scopes []string) {
					called = true
					if v := c.Value(key); v == nil || !v.(bool) {
						t.Error("AcquireTokenOnBehalfOf received unexpected Context")
					}
					if len(scopes) != 1 || scopes[0] != liveTestScope {
						t.Errorf(`unexpected scopes "%v"`, scopes)
					}
					if assertion != expectedAssertion {
						t.Errorf(`unexpected assertion "%s"`, assertion)
					}
				},
			}
			getConfidentialClient = func(clientID, tenantID string, cred confidential.Credential, co *azcore.ClientOptions, opts ...confidential.Option) (confidentialClient, error) {
				if clientID != fakeClientID {
					t.Errorf(`unexpected clientID "%s"`, clientID)
				}
				if tenantID != fakeTenantID {
					t.Errorf(`unexpected tenantID "%s"`, tenantID)
				}
				msalOpts := confidential.Options{}
				for _, o := range opts {
					o(&msalOpts)
				}
				if test.sendX5C != msalOpts.SendX5C {
					t.Fatal("incorrect value for SendX5C")
				}
				return fake, nil
			}
			cred, err := test.ctor()
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
			if !called {
				t.Fatal("validation function wasn't called")
			}
		})
	}
}
