// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

type fakeTokenCredential struct{}

func (fakeTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-bearer-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// bearerStub stands in for the bearer token policy so tests can tell whether it was wrapped.
type bearerStub struct{}

func (bearerStub) Do(req *policy.Request) (*http.Response, error) { return req.Next() }

const (
	accountBlobURL      = "https://account.blob.core.windows.net/container/blob"
	customDomainBlobURL = "https://files.contoso.com/container/blob?sv=2026-12-06&sig=secret"
)

func sessionOptions(enabled *bool, accountName string) *ClientOptions {
	return &ClientOptions{Session: exported.SessionOptions{Enabled: enabled, AccountName: accountName}}
}

func TestSessionAuthPolicyResolvesEnabled(t *testing.T) {
	bearer := bearerStub{}
	for _, tt := range []struct {
		name        string
		enabled     *bool
		wantSession bool
	}{
		{"NilResolvesToDisabled", nil, false},
		{"False", to.Ptr(false), false},
		{"True", to.Ptr(true), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p, err := newSessionAuthPolicy(accountBlobURL, fakeTokenCredential{}, bearer, sessionOptions(tt.enabled, ""))
			require.NoError(t, err)
			if tt.wantSession {
				require.NotEqual(t, policy.Policy(bearer), p, "the bearer policy is wrapped by the session policy")
			} else {
				require.Equal(t, policy.Policy(bearer), p, "the bearer policy is used as-is")
			}
		})
	}
}

func TestSessionAuthPolicyAccountName(t *testing.T) {
	bearer := bearerStub{}

	t.Run("DerivedFromURL", func(t *testing.T) {
		p, err := newSessionAuthPolicy(accountBlobURL, fakeTokenCredential{}, bearer, sessionOptions(to.Ptr(true), ""))
		require.NoError(t, err)
		require.NotEqual(t, policy.Policy(bearer), p)
	})
	t.Run("DerivedFromPathStyleURL", func(t *testing.T) {
		p, err := newSessionAuthPolicy("http://127.0.0.1:10000/devstoreaccount1/container/blob", fakeTokenCredential{}, bearer, sessionOptions(to.Ptr(true), ""))
		require.NoError(t, err)
		require.NotEqual(t, policy.Policy(bearer), p)
	})
	t.Run("TrueCustomDomainFails", func(t *testing.T) {
		_, err := newSessionAuthPolicy(customDomainBlobURL, fakeTokenCredential{}, bearer, sessionOptions(to.Ptr(true), ""))
		require.ErrorContains(t, err, "account name could not be determined")
		require.ErrorContains(t, err, "https://files.contoso.com/container/blob")
		require.NotContains(t, err.Error(), "sig=", "the error must not disclose the SAS")
		require.NotContains(t, err.Error(), "secret")
	})
	t.Run("NilCustomDomainUsesBearer", func(t *testing.T) {
		p, err := newSessionAuthPolicy(customDomainBlobURL, fakeTokenCredential{}, bearer, sessionOptions(nil, ""))
		require.NoError(t, err, "only an explicit true fails construction")
		require.Equal(t, policy.Policy(bearer), p)
	})
	t.Run("FalseCustomDomainUsesBearer", func(t *testing.T) {
		p, err := newSessionAuthPolicy(customDomainBlobURL, fakeTokenCredential{}, bearer, sessionOptions(to.Ptr(false), ""))
		require.NoError(t, err, "only an explicit true fails construction")
		require.Equal(t, policy.Policy(bearer), p)
	})
	t.Run("ExplicitAccountNameForCustomDomain", func(t *testing.T) {
		p, err := newSessionAuthPolicy(customDomainBlobURL, fakeTokenCredential{}, bearer, sessionOptions(to.Ptr(true), "account"))
		require.NoError(t, err)
		require.NotEqual(t, policy.Policy(bearer), p)
	})
}

func TestGetAzClientSessionsRequireTokenCredential(t *testing.T) {
	sharedKey, err := exported.NewSharedKeyCredential("account", "dGVzdC1rZXk=")
	require.NoError(t, err)

	for name, enabled := range map[string]*bool{"Nil": nil, "False": to.Ptr(false), "True": to.Ptr(true)} {
		t.Run("SharedKey"+name, func(t *testing.T) {
			_, err := GetAzClient(accountBlobURL, nil, sharedKey, sessionOptions(enabled, ""))
			require.NoError(t, err, "shared key clients ignore session options")
		})
		t.Run("SASOrAnonymous"+name, func(t *testing.T) {
			_, err := GetAzClient(accountBlobURL+"?sv=2026-12-06&sig=x", nil, nil, sessionOptions(enabled, ""))
			require.NoError(t, err, "SAS and anonymous clients ignore session options")
		})
	}
	t.Run("SharedKeyEnabledCustomDomain", func(t *testing.T) {
		// the account name check belongs to the token path only
		_, err := GetAzClient(customDomainBlobURL, nil, sharedKey, sessionOptions(to.Ptr(true), ""))
		require.NoError(t, err)
	})
	t.Run("TokenEnabledCustomDomain", func(t *testing.T) {
		_, err := GetAzClient(customDomainBlobURL, fakeTokenCredential{}, nil, sessionOptions(to.Ptr(true), ""))
		require.Error(t, err)
	})
}

func TestNewContainerSessionProvider(t *testing.T) {
	_, err := NewContainerSessionProvider(accountBlobURL, nil, nil)
	require.ErrorContains(t, err, "token credential is required")

	_, err = NewContainerSessionProvider("://not a url", fakeTokenCredential{}, nil)
	require.Error(t, err)

	p, err := NewContainerSessionProvider(accountBlobURL+"?sv=2026-12-06&sig=x", fakeTokenCredential{}, nil)
	require.NoError(t, err)
	require.NotNil(t, p)

	// the caller's options aren't modified by forcing sessions off for the provider's own pipeline
	opts := sessionOptions(to.Ptr(true), "account")
	_, err = NewContainerSessionProvider(accountBlobURL, fakeTokenCredential{}, opts)
	require.NoError(t, err)
	require.Equal(t, to.Ptr(true), opts.Session.Enabled)
	require.Equal(t, "account", opts.Session.AccountName)
}

func TestRedactedURL(t *testing.T) {
	require.Equal(t, "https://files.contoso.com/c/b", redactedURL("https://user:pass@files.contoso.com/c/b?sig=secret#frag"))
	require.Equal(t, "<unparseable URL>", redactedURL("://bad"))
}
