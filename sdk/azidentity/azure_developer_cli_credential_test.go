// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

var (
	mockAzdSuccess = func(_ context.Context, credName string, _ string) ([]byte, error) {
		if credName != credNameAzureDeveloperCLI {
			return nil, errors.New("unexpected credential name: " + credName)
		}
		return []byte(`{
  "token": "mocktoken",
  "expiresOn": "2001-02-03T04:05:06Z"
}
`), nil
	}
	mockAzdFailure = func(_ context.Context, credName string, _ string) ([]byte, error) {
		if credName != credNameAzureDeveloperCLI {
			return nil, errors.New("unexpected credential name: " + credName)
		}
		return nil, newAuthenticationFailedError(credNameAzureDeveloperCLI, "azd error", nil)
	}
)

func TestAzureDeveloperCLICredential_Claims(t *testing.T) {
	tro := policy.TokenRequestOptions{
		Scopes: []string{liveTestScope},
		Claims: `{"access_token":{"xms_cc":{"values":["cp1"]}}}`,
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(tro.Claims))
	t.Run("old azd", func(t *testing.T) {
		cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
			exec: func(_ context.Context, _, command string) ([]byte, error) {
				require.Contains(t, command, "--claims "+encoded)
				return nil, newAuthenticationFailedError(credNameAzureDeveloperCLI, "unknown flag: --claims", nil)
			},
		})
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, tro)
		require.ErrorContains(t, err, mfaRequired)
		require.ErrorContains(t, err, "Upgrade")
	})

	t.Run("recent azd", func(t *testing.T) {
		cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
			exec: func(_ context.Context, _, command string) ([]byte, error) {
				require.Contains(t, command, "--claims "+encoded)
				return nil, newAuthenticationFailedError(credNameAzureDeveloperCLI, `{"type":"consoleMessage","timestamp":"...","data":{"message":"\nERROR: fetching token: AADSTS50079: Due to a configuration change made by your administrator, or because you moved to a new location, you must enroll in multi-factor authentication'. Trace ID: ... Correlation ID: ... Timestamp: ...\n"}}
{"type":"consoleMessage","timestamp":"...","data":{"message":"Suggestion: reauthentication required, run azd auth login --scope ... to acquire a new token.\n"}}`, nil)
			},
		})
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, tro)
		require.ErrorContains(t, err, mfaRequired)
		require.ErrorAs(t, err, new(*AuthenticationFailedError))
	})
}

func TestAzureDeveloperCLICredential_DefaultChainError(t *testing.T) {
	cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
		inDefaultChain: true,
		exec:           mockAzdFailure,
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

func TestAzureDeveloperCLICredential_Error(t *testing.T) {
	// GetToken shouldn't invoke the CLI a second time after a failure
	authNs := 0
	expected := newCredentialUnavailableError(credNameAzureDeveloperCLI, "it didn't work")
	o := AzureDeveloperCLICredentialOptions{
		exec: func(context.Context, string, string) ([]byte, error) {
			authNs++
			return nil, expected
		},
	}
	cred, err := NewAzureDeveloperCLICredential(&o)
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

func TestAzureDeveloperCLICredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
		exec: func(ctx context.Context, credName, command string) ([]byte, error) {
			require.Equal(t, command, "azd auth token -o json --no-prompt --scope "+liveTestScope)
			return mockAzdSuccess(ctx, credName, command)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	at, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if at.Token != "mocktoken" {
		t.Fatalf("unexpected access token %q", at.Token)
	}
	expected := time.Date(2001, 2, 3, 4, 5, 6, 000, time.UTC)
	if actual := at.ExpiresOn; !actual.Equal(expected) || actual.Location() != time.UTC {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestAzureDeveloperCLICredential_TenantID(t *testing.T) {
	expected := "expected-tenant-id"
	called := false
	options := AzureDeveloperCLICredentialOptions{
		TenantID: expected,
		exec: func(ctx context.Context, credName, command string) ([]byte, error) {
			called = true
			require.Contains(t, command, " --tenant-id "+expected)
			return mockAzdSuccess(ctx, credName, command)
		},
	}
	cred, err := NewAzureDeveloperCLICredential(&options)
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

func TestAzureDeveloperCLICredential_JSONErrorParsing(t *testing.T) {
	// These tests verify that when shellExec parses JSON errors from azd stderr,
	// the credential properly handles the cleaned error messages.
	// The exec mock simulates what shellExec would return after parsing.

	t.Run("parsed JSON message", func(t *testing.T) {
		cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
			exec: func(_ context.Context, _, _ string) ([]byte, error) {
				// Simulates shellExec having parsed: {"data":{"message":"\nERROR: fetching token: authentication failed\n"}}
				return nil, newAuthenticationFailedError(credNameAzureDeveloperCLI, "ERROR: fetching token: authentication failed", nil)
			},
		})
		require.NoError(t, err)
		_, err = cred.GetToken(context.Background(), testTRO)
		require.Error(t, err)
		require.ErrorContains(t, err, "ERROR: fetching token: authentication failed")
	})

	t.Run("plain text error", func(t *testing.T) {
		cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
			exec: func(_ context.Context, _, _ string) ([]byte, error) {
				// Simulates shellExec having received non-JSON stderr
				return nil, newAuthenticationFailedError(credNameAzureDeveloperCLI, "ERROR: plain text error message", nil)
			},
		})
		require.NoError(t, err)
		_, err = cred.GetToken(context.Background(), testTRO)
		require.Error(t, err)
		require.ErrorContains(t, err, "ERROR: plain text error message")
	})
}
