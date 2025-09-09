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

func azurePowerShellTokenOutput(expiresOn int64) []byte {
	return []byte(fmt.Sprintf(`{
  "Token": %q,
  "ExpiresOn": %d
}`, tokenValue, expiresOn))
}

func mockazurePowerShellTokenProviderFailure(context.Context, []string, string, string) ([]byte, error) {
	return nil, newAuthenticationFailedError(credNameAzurePowerShell, "mock provider error", nil)
}

func mockazurePowerShellTokenProviderSuccess(context.Context, []string, string, string) ([]byte, error) {
	return azurePowerShellTokenOutput(638930167310000000), nil
}

func TestAzurePowerShellCredential_DefaultChainError(t *testing.T) {
	cred, err := NewAzurePowerShellCredential(&AzurePowerShellCredentialOptions{
		inDefaultChain: true,
		tokenProvider:  mockazurePowerShellTokenProviderFailure,
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

func TestAzurePowerShellCredential_Error(t *testing.T) {
	// GetToken shouldn't invoke Azure PowerShell a second time after a failure
	authNs := 0
	expected := newCredentialUnavailableError(credNameAzurePowerShell, "it didn't work")
	o := AzurePowerShellCredentialOptions{
		tokenProvider: func(context.Context, []string, string, string) ([]byte, error) {
			authNs++
			return nil, expected
		},
	}
	cred, err := NewAzurePowerShellCredential(&o)
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

func TestAzurePowerShellCredential_GetTokenSuccess(t *testing.T) {
	expectedExpiresOn := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	t.Run("fetches token with correct expiration", func(t *testing.T) {
		ExpiresOn := epochTicks + expectedExpiresOn.UTC().UnixNano()/100
		cred, err := NewAzurePowerShellCredential(&AzurePowerShellCredentialOptions{
			tokenProvider: func(context.Context, []string, string, string) ([]byte, error) {
				output := azurePowerShellTokenOutput(ExpiresOn)
				return output, nil
			},
		})
		require.NoError(t, err)

		actual, err := cred.GetToken(context.Background(), testTRO)
		require.NoError(t, err)
		require.NotEmpty(t, actual.Token, "Token should not be empty")
		require.True(t, actual.ExpiresOn.Equal(expectedExpiresOn))
		require.Equal(t, time.UTC, actual.ExpiresOn.Location())
	})
}

func TestAzurePowerShellCredential_GetTokenInvalidToken(t *testing.T) {
	options := AzurePowerShellCredentialOptions{}
	options.tokenProvider = mockazurePowerShellTokenProviderFailure
	cred, err := NewAzurePowerShellCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestAzurePowerShellCredential_Subscription(t *testing.T) {
	called := false
	for _, want := range []string{"", "expected-subscription"} {
		t.Run(fmt.Sprintf("subscription=%q", want), func(t *testing.T) {
			options := AzurePowerShellCredentialOptions{
				Subscription: want,
				tokenProvider: func(ctx context.Context, scopes []string, tenant, subscription string) ([]byte, error) {
					called = true
					if subscription != want {
						t.Fatalf("wanted subscription %q, got %q", want, subscription)
					}
					return mockazurePowerShellTokenProviderSuccess(ctx, scopes, tenant, subscription)
				},
			}
			cred, err := NewAzurePowerShellCredential(&options)
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

func TestAzurePowerShellCredential_TenantID(t *testing.T) {
	expected := "expected-tenant-id"
	called := false
	options := AzurePowerShellCredentialOptions{
		TenantID: expected,
		tokenProvider: func(ctx context.Context, scopes []string, tenantID, subscription string) ([]byte, error) {
			called = true
			if tenantID != expected {
				t.Fatal("Unexpected tenant ID: " + tenantID)
			}
			return mockazurePowerShellTokenProviderSuccess(ctx, scopes, tenantID, subscription)
		},
	}
	cred, err := NewAzurePowerShellCredential(&options)
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
