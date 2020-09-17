// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"os"
	"testing"
)

func initEnvironmentVarsForTest() error {
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		return err
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		return err
	}
	err = os.Setenv("AZURE_CLIENT_SECRET", secret)
	if err != nil {
		return err
	}
	return nil
}

func resetEnvironmentVarsForTest() {
	clearEnvVars("AZURE_TENANT_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET")
}

func TestEnvironmentCredential_TenantIDNotSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_SECRET", secret)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	_, err = NewEnvironmentCredential(nil)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
	var credentialUnavailable *CredentialUnavailableError
	if !errors.As(err, &credentialUnavailable) {
		t.Fatalf("Expected a credential unavailable error, instead received: %T", err)
	}
}

func TestEnvironmentCredential_ClientIDNotSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_SECRET", secret)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	_, err = NewEnvironmentCredential(nil)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
	var credentialUnavailable *CredentialUnavailableError
	if !errors.As(err, &credentialUnavailable) {
		t.Fatalf("Expected a credential unavailable error, instead received: %T", err)
	}
}

func TestEnvironmentCredential_ClientSecretNotSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	_, err = NewEnvironmentCredential(nil)
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
	var credentialUnavailable *CredentialUnavailableError
	if !errors.As(err, &credentialUnavailable) {
		t.Fatalf("Expected a credential unavailable error, instead received: %T", err)
	}
}
