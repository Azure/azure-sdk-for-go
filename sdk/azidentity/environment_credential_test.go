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
	clearEnvVars("AZURE_TENANT_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET", "AZURE_CLIENT_CERTIFICATE_PATH", "AZURE_USERNAME", "AZURE_PASSWORD")
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

func TestEnvironmentCredential_ClientSecretSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_SECRET", secret)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Did not expect an error. Received: %v", err)
	}
	if _, ok := cred.cred.(*ClientSecretCredential); !ok {
		t.Fatalf("Did not receive the right credential type. Expected *azidentity.ClientSecretCredential, Received: %t", cred)
	}
}

func TestEnvironmentCredential_ClientCertificatePathSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", certificatePath)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Did not expect an error. Received: %v", err)
	}
	if _, ok := cred.cred.(*ClientCertificateCredential); !ok {
		t.Fatalf("Did not receive the right credential type. Expected *azidentity.ClientCertificateCredential, Received: %t", cred)
	}
}

func TestEnvironmentCredential_UsernameOnlySet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_USERNAME", "username")
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

func TestEnvironmentCredential_UsernamePasswordSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_USERNAME", "username")
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv("AZURE_PASSWORD", "password")
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Did not expect an error. Received: %v", err)
	}
	if _, ok := cred.cred.(*UsernamePasswordCredential); !ok {
		t.Fatalf("Did not receive the right credential type. Expected *azidentity.UsernamePasswordCredential, Received: %t", cred)
	}
}
