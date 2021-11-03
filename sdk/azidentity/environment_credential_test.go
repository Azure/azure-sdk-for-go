// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
	err = os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", "testdata/certificate.pem")
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

func TestEnvironmentCredential_ClientSecretLive(t *testing.T) {
	vars := map[string]string{
		"AZURE_CLIENT_ID":     liveSP.clientID,
		"AZURE_CLIENT_SECRET": liveSP.secret,
		"AZURE_TENANT_ID":     liveSP.tenantID,
	}
	for _, v := range vars {
		if v == "" {
			t.Skip("missing live service principal configuration")
		}
	}
	setEnvironmentVariables(t, vars)
	cred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if tk.Token == "" {
		t.Fatalf("GetToken returned an invalid token")
	}
	if !tk.ExpiresOn.After(time.Now().UTC()) {
		t.Fatalf("GetToken returned an invalid expiration time")
	}
}

func TestEnvironmentCredential_UserPasswordLive(t *testing.T) {
	vars := map[string]string{
		"AZURE_CLIENT_ID": developerSignOnClientID,
		"AZURE_TENANT_ID": os.Getenv("AZURE_IDENTITY_TEST_TENANTID"),
		"AZURE_USERNAME":  os.Getenv("AZURE_IDENTITY_TEST_USERNAME"),
		"AZURE_PASSWORD":  os.Getenv("AZURE_IDENTITY_TEST_PASSWORD"),
	}
	for _, v := range vars {
		if v == "" {
			t.Skip("missing live user configuration")
		}
	}
	setEnvironmentVariables(t, vars)
	cred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if tk.Token == "" {
		t.Fatalf("GetToken returned an invalid token")
	}
	if !tk.ExpiresOn.After(time.Now().UTC()) {
		t.Fatalf("GetToken returned an invalid expiration time")
	}
}
