// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

func azurePowerShellTokenOutput(expiresOn int64) []byte {
	return []byte(fmt.Sprintf(`{
  "Token": %q,
  "ExpiresOn": %d
}`, tokenValue, expiresOn))
}

func mockAzurePowerShellFailure(_ context.Context, credName string, _ string) ([]byte, error) {
	if credName != credNameAzurePowerShell {
		return nil, errors.New("unexpected credential name: " + credName)
	}
	return nil, newAuthenticationFailedError(credNameAzurePowerShell, "Azure PowerShell error", nil)
}

func mockAzurePowerShellSuccess(_ context.Context, credName string, _ string) ([]byte, error) {
	if credName != credNameAzurePowerShell {
		return nil, errors.New("unexpected credential name: " + credName)
	}
	return azurePowerShellTokenOutput(638930167310000000), nil
}

func TestAzurePowerShellCredential_Claims(t *testing.T) {
	tro := policy.TokenRequestOptions{
		Scopes: []string{liveTestScope},
		Claims: `{"access_token":{"xms_cc":{"values":["cp1"]}}}`,
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(tro.Claims))
	exec := func(context.Context, string, string) ([]byte, error) {
		t.Fatal("GetToken shouldn't run Azure PowerShell when claims are specified")
		return nil, nil
	}

	cred, err := NewAzurePowerShellCredential(&AzurePowerShellCredentialOptions{exec: exec})
	require.NoError(t, err)
	_, err = cred.GetToken(ctx, tro)
	require.ErrorContains(t, err, fmt.Sprintf("Connect-AzAccount -ClaimsChallenge '%s'", encoded))

	t.Run("with tenant", func(t *testing.T) {
		expected := fmt.Sprintf("Connect-AzAccount -TenantId '%s' -ClaimsChallenge '%s'", fakeTenantID, encoded)

		cp := tro
		cp.TenantID = fakeTenantID
		_, err = cred.GetToken(ctx, cp)
		require.ErrorContains(t, err, expected)

		cred, err := NewAzurePowerShellCredential(&AzurePowerShellCredentialOptions{TenantID: fakeTenantID, exec: exec})
		require.NoError(t, err)
		_, err = cred.GetToken(ctx, tro)
		require.ErrorContains(t, err, expected)
	})
}

func TestAzurePowerShellCredential_DefaultChainError(t *testing.T) {
	cred, err := NewAzurePowerShellCredential(&AzurePowerShellCredentialOptions{
		inDefaultChain: true,
		exec:           mockAzurePowerShellFailure,
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
		exec: func(context.Context, string, string) ([]byte, error) {
			authNs++
			return nil, expected
		},
	}
	cred, err := NewAzurePowerShellCredential(&o)
	require.NoError(t, err)
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
			exec: func(context.Context, string, string) ([]byte, error) {
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
	options.exec = mockAzurePowerShellFailure
	cred, err := NewAzurePowerShellCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestAzurePowerShellCredential_TenantID(t *testing.T) {
	expected := "expected-tenant-id"
	called := false
	options := AzurePowerShellCredentialOptions{
		TenantID: expected,
		exec: func(ctx context.Context, credName, command string) ([]byte, error) {
			called = true
			splitCommand := strings.Split(command, " ")
			encodedScript := splitCommand[len(splitCommand)-1]
			decodedScript, err := base64DecodeUTF16LE(encodedScript)
			require.NoError(t, err)
			require.Contains(t, decodedScript, fmt.Sprintf("$params['TenantId'] = '%s'", expected))
			return mockAzurePowerShellSuccess(ctx, credName, command)
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
