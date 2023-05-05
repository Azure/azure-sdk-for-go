//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

func resetEnvironmentVarsForTest() {
	clearEnvVars(azureTenantID, azureClientID, azureClientSecret, azureClientCertificatePath, azureUsername, azurePassword)
}

func TestEnvironmentCredential_TenantIDNotSet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv(azureClientID, fakeClientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientSecret, fakeSecret)
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
	err := os.Setenv(azureTenantID, fakeTenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientSecret, fakeSecret)
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
	err := os.Setenv(azureTenantID, fakeTenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientID, fakeClientID)
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
	err := os.Setenv(azureTenantID, fakeTenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientID, fakeClientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientSecret, fakeSecret)
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
	err := os.Setenv(azureTenantID, fakeTenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientID, fakeClientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientCertificatePath, "testdata/certificate.pem")
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

func TestEnvironmentCredential_ClientCertificatePassword(t *testing.T) {
	for key, value := range map[string]string{
		azureTenantID:              fakeTenantID,
		azureClientID:              fakeClientID,
		azureClientCertificatePath: "testdata/certificate_encrypted_key.pfx",
	} {
		t.Setenv(key, value)
	}
	for _, correctPassword := range []bool{true, false} {
		t.Run(fmt.Sprintf("%v", correctPassword), func(t *testing.T) {
			password := "wrong password"
			if correctPassword {
				password = "password"
			}
			t.Setenv(azureClientCertificatePassword, password)
			cred, err := NewEnvironmentCredential(nil)
			if correctPassword {
				if err != nil {
					t.Fatal(err)
				}
				if _, ok := cred.cred.(*ClientCertificateCredential); !ok {
					t.Fatalf("expected *azidentity.ClientCertificateCredential, got %t", cred)
				}
			} else if err == nil || !strings.Contains(err.Error(), "password") {
				t.Fatal("expected an error about the password")
			}
		})
	}
}

func TestEnvironmentCredential_UsernameOnlySet(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv(azureTenantID, fakeTenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientID, fakeClientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureUsername, "username")
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
	err := os.Setenv(azureTenantID, fakeTenantID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureClientID, fakeClientID)
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azureUsername, "username")
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	err = os.Setenv(azurePassword, "password")
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

func TestEnvironmentCredential_SendCertificateChain(t *testing.T) {
	certData, err := os.ReadFile(liveSP.pfxPath)
	if err != nil {
		t.Fatal(err)
	}
	certs, _, err := ParseCertificates(certData, nil)
	if err != nil {
		t.Fatal(err)
	}
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithBody(instanceDiscoveryResponse))
	srv.AppendResponse(mock.WithBody(tenantDiscoveryResponse))
	srv.AppendResponse(mock.WithPredicate(validateX5C(t, certs)), mock.WithBody(accessTokenRespSuccess))
	srv.AppendResponse()

	vars := map[string]string{
		azureClientID:              liveSP.clientID,
		azureClientCertificatePath: liveSP.pfxPath,
		azureTenantID:              liveSP.tenantID,
		envVarSendCertChain:        "true",
	}
	setEnvironmentVariables(t, vars)
	cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: srv}})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("unexpected token: %s", tk.Token)
	}
}

func TestEnvironmentCredential_ClientSecretLive(t *testing.T) {
	vars := map[string]string{
		azureClientID:     liveSP.clientID,
		azureClientSecret: liveSP.secret,
		azureTenantID:     liveSP.tenantID,
	}
	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			setEnvironmentVariables(t, vars)
			opts, stop := initRecording(t)
			defer stop()
			cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
				ClientOptions:            opts,
				DisableInstanceDiscovery: disabledID,
			})
			if err != nil {
				t.Fatalf("failed to construct credential: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestEnvironmentCredentialADFS_ClientSecretLive(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if adfsLiveSP.clientID == "" || adfsLiveSP.secret == "" {
			t.Skip("set ADFS_SP_* environment variables to run this test live")
		}
	}
	vars := map[string]string{
		azureClientID:      adfsLiveSP.clientID,
		azureClientSecret:  adfsLiveSP.secret,
		azureTenantID:      "adfs",
		azureAuthorityHost: adfsAuthority,
	}
	setEnvironmentVariables(t, vars)
	opts, stop := initRecording(t)
	defer stop()
	cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
		ClientOptions:            opts,
		DisableInstanceDiscovery: true,
	})
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}

func TestEnvironmentCredential_InvalidClientSecretLive(t *testing.T) {
	vars := map[string]string{
		azureClientID:     liveSP.clientID,
		azureClientSecret: "invalid secret",
		azureTenantID:     liveSP.tenantID,
	}
	setEnvironmentVariables(t, vars)
	opts, stop := initRecording(t)
	defer stop()
	cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: opts})
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if !reflect.ValueOf(tk).IsZero() {
		t.Fatal("expected a zero value AccessToken")
	}
	if e, ok := err.(*AuthenticationFailedError); ok {
		if e.RawResponse == nil {
			t.Fatal("expected a non-nil RawResponse")
		}
	} else {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	if !strings.HasPrefix(err.Error(), credNameSecret) {
		t.Fatal("missing credential type prefix")
	}
}

func TestEnvironmentCredential_UserPasswordLive(t *testing.T) {
	vars := map[string]string{
		azureClientID: developerSignOnClientID,
		azureTenantID: liveUser.tenantID,
		azureUsername: liveUser.username,
		azurePassword: liveUser.password,
	}
	setEnvironmentVariables(t, vars)
	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			opts, stop := initRecording(t)
			defer stop()
			cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
				ClientOptions:            opts,
				DisableInstanceDiscovery: disabledID,
			})
			if err != nil {
				t.Fatalf("failed to construct credential: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestEnvironmentCredentialADFS_UserPasswordLive(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if adfsLiveUser.clientID == "" || adfsLiveUser.username == "" || adfsLiveUser.password == "" {
			t.Skip("set ADFS_IDENTITY_TEST_* environment variables to run this test live")
		}
	}
	vars := map[string]string{
		azureClientID:      adfsLiveUser.clientID,
		azureTenantID:      "adfs",
		azureUsername:      adfsLiveUser.username,
		azurePassword:      adfsLiveUser.password,
		azureAuthorityHost: adfsAuthority,
	}
	setEnvironmentVariables(t, vars)
	opts, stop := initRecording(t)
	defer stop()
	cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
		ClientOptions:            opts,
		DisableInstanceDiscovery: true,
	})
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}

func TestEnvironmentCredential_InvalidPasswordLive(t *testing.T) {
	vars := map[string]string{
		azureClientID: developerSignOnClientID,
		azureTenantID: liveUser.tenantID,
		azureUsername: liveUser.username,
		azurePassword: "invalid password",
	}
	setEnvironmentVariables(t, vars)
	opts, stop := initRecording(t)
	defer stop()
	cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: opts})
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if !reflect.ValueOf(tk).IsZero() {
		t.Fatal("expected a zero value AccessToken")
	}
	if e, ok := err.(*AuthenticationFailedError); ok {
		if e.RawResponse == nil {
			t.Fatal("expected a non-nil RawResponse")
		}
	} else {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	if !strings.HasPrefix(err.Error(), credNameUserPassword) {
		t.Fatal("missing credential type prefix")
	}
}
