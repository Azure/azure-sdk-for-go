//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// azTokenOutput returns JSON output similar to az account get-access-token.
// All versions of az return expiresOn, a local timestamp. v2.54.0+
// additionally return expires_on, a Unix timestamp. If the expires_on
// argument to this function is 0, the returned JSON omits expires_on.
func azTokenOutput(expiresOn string, expires_on int64) []byte {
	e_o := ""
	if expires_on != 0 {
		e_o = fmt.Sprintf(`
		"expires_on": %d,
`, expires_on)
	}
	return []byte(fmt.Sprintf(`{
  "accessToken": %q,
  "expiresOn": %q,%s
  "subscription": "fake-subscription",
  "tenant": %q,
  "tokenType": "Bearer"
}`, tokenValue, expiresOn, e_o, fakeTenantID))
}

func mockAzTokenProviderFailure(context.Context, []string, string, string) ([]byte, error) {
	return nil, newAuthenticationFailedError(credNameAzureCLI, "mock provider error", nil, nil)
}

func mockAzTokenProviderSuccess(context.Context, []string, string, string) ([]byte, error) {
	return azTokenOutput("2001-02-03 04:05:06.000007", 0), nil
}

func TestAzureCLICredential_DefaultChainError(t *testing.T) {
	cred, err := NewAzureCLICredential(&AzureCLICredentialOptions{
		inDefaultChain: true,
		tokenProvider:  mockAzTokenProviderFailure,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	var cu credentialUnavailable
	if !errors.As(err, &cu) {
		t.Fatalf("expected %T, got %T: %q", cu, err, err)
	}
}

func TestAzureCLICredential_Error(t *testing.T) {
	// GetToken shouldn't invoke the CLI a second time after a failure
	authNs := 0
	expected := newCredentialUnavailableError(credNameAzureCLI, "it didn't work")
	o := AzureCLICredentialOptions{
		tokenProvider: func(context.Context, []string, string, string) ([]byte, error) {
			authNs++
			return nil, expected
		},
	}
	cred, err := NewAzureCLICredential(&o)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatal("expected an error")
	}
	if err != expected {
		t.Fatalf("expected %v, got %v", expected, err)
	}
	if authNs != 1 {
		t.Fatalf("expected 1 authN, got %d", authNs)
	}
}

func TestAzureCLICredential_GetTokenSuccess(t *testing.T) {
	expectedExpiresOn := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	for _, withExpires_on := range []bool{false, true} {
		name := "without expires_on"
		if withExpires_on {
			name = "with expires_on"
		}
		t.Run(name, func(t *testing.T) {
			ExpiresOn := expectedExpiresOn.Local().Format("2006-01-02 15:04:05.999999999")
			expires_on := int64(0)
			if withExpires_on {
				// set the wrong time for ExpiresOn so this test fails if the credential uses it
				ExpiresOn = "2001-01-01 01:01:01.000000"
				expires_on = expectedExpiresOn.Unix()
			}
			cred, err := NewAzureCLICredential(&AzureCLICredentialOptions{
				tokenProvider: func(context.Context, []string, string, string) ([]byte, error) {
					output := azTokenOutput(ExpiresOn, expires_on)
					return output, nil
				},
			})
			require.NoError(t, err)

			actual, err := cred.GetToken(context.Background(), testTRO)
			require.NoError(t, err)
			require.True(t, actual.ExpiresOn.Equal(expectedExpiresOn))
			require.Equal(t, time.UTC, actual.ExpiresOn.Location())
			require.Equal(t, tokenValue, actual.Token)
		})
	}
}

func TestAzureCLICredential_GetTokenInvalidToken(t *testing.T) {
	options := AzureCLICredentialOptions{}
	options.tokenProvider = mockAzTokenProviderFailure
	cred, err := NewAzureCLICredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestAzureCLICredential_Subscription(t *testing.T) {
	called := false
	for _, want := range []string{"", "expected-subscription"} {
		t.Run(fmt.Sprintf("subscription=%q", want), func(t *testing.T) {
			options := AzureCLICredentialOptions{
				Subscription: want,
				tokenProvider: func(ctx context.Context, scopes []string, tenant, subscription string) ([]byte, error) {
					called = true
					if subscription != want {
						t.Fatalf("wanted subscription %q, got %q", want, subscription)
					}
					return mockAzTokenProviderSuccess(ctx, scopes, tenant, subscription)
				},
			}
			cred, err := NewAzureCLICredential(&options)
			if err != nil {
				t.Fatal(err)
			}
			_, err = cred.GetToken(context.Background(), testTRO)
			if err != nil {
				t.Fatal(err)
			}
			if !called {
				t.Fatal("token provider wasn't called")
			}
		})
	}
}

func TestAzureCLICredential_TenantID(t *testing.T) {
	expected := "expected-tenant-id"
	called := false
	options := AzureCLICredentialOptions{
		TenantID: expected,
		tokenProvider: func(ctx context.Context, scopes []string, tenantID, subscription string) ([]byte, error) {
			called = true
			if tenantID != expected {
				t.Fatal("Unexpected tenant ID: " + tenantID)
			}
			return mockAzTokenProviderSuccess(ctx, scopes, tenantID, subscription)
		},
	}
	cred, err := NewAzureCLICredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !called {
		t.Fatal("token provider wasn't called")
	}
}
