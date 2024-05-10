//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func unsetEnvironmentVarsForTest(t *testing.T) {
	for _, k := range []string{
		azureClientCertificatePath, azureClientID, azureClientSecret, azurePassword, azureTenantID, azureUsername,
	} {
		if v, set := os.LookupEnv(k); set {
			require.NoError(t, os.Unsetenv(k))
			t.Cleanup(func() { require.NoError(t, os.Setenv(k, v)) })
		}
	}
}

func TestEnvironmentCredential(t *testing.T) {
	unsetEnvironmentVarsForTest(t)
	for _, test := range []struct {
		env  map[string]string
		cred azcore.TokenCredential
	}{
		{
			cred: &ClientCertificateCredential{},
			env: map[string]string{
				azureClientCertificatePath: "testdata/certificate.pem",
				azureClientID:              fakeClientID,
				azureTenantID:              fakeTenantID,
			},
		},
		{
			cred: &ClientSecretCredential{},
			env: map[string]string{
				azureClientID:     fakeClientID,
				azureClientSecret: fakeSecret,
				azureTenantID:     fakeTenantID,
			},
		},
		{
			cred: &UsernamePasswordCredential{},
			env: map[string]string{
				azureClientID: fakeClientID,
				azurePassword: "fake",
				azureTenantID: fakeTenantID,
				azureUsername: fakeUsername,
			},
		},
	} {
		t.Run(fmt.Sprintf("%T", test.cred), func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}
			cred, err := NewEnvironmentCredential(nil)
			require.NoError(t, err)
			require.IsType(t, test.cred, cred.cred)
			for k := range test.env {
				t.Run("missing "+k, func(t *testing.T) {
					before := os.Getenv(k)
					require.NoError(t, os.Unsetenv(k))
					defer os.Setenv(k, before)
					_, err := NewEnvironmentCredential(nil)
					require.Error(t, err)
				})
			}
		})
	}
}

func TestEnvironmentCredential_CertificateErrors(t *testing.T) {
	unsetEnvironmentVarsForTest(t)
	for _, test := range []struct {
		name, path string
	}{
		{"file doesn't exist", filepath.Join(t.TempDir(), t.Name())},
		{"invalid file", "testdata/certificate-wrong-key.pem"},
	} {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range map[string]string{
				azureClientID:              fakeClientID,
				azureClientCertificatePath: test.path,
				azureTenantID:              fakeTenantID,
			} {
				t.Setenv(k, v)
				_, err := NewEnvironmentCredential(nil)
				if err == nil {
					t.Fatal("expected an error")
				}
			}
		})
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

func TestEnvironmentCredential_SendCertificateChain(t *testing.T) {
	certData, err := os.ReadFile(liveSP.pfxPath)
	if err != nil {
		t.Fatal(err)
	}
	certs, _, err := ParseCertificates(certData, nil)
	if err != nil {
		t.Fatal(err)
	}
	unsetEnvironmentVarsForTest(t)
	sts := mockSTS{tokenRequestCallback: validateX5C(t, certs)}
	vars := map[string]string{
		azureClientID:              liveSP.clientID,
		azureClientCertificatePath: liveSP.pfxPath,
		azureTenantID:              liveSP.tenantID,
		envVarSendCertChain:        "true",
	}
	setEnvironmentVariables(t, vars)
	cred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{ClientOptions: azcore.ClientOptions{Transport: &sts}})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("unexpected token: %s", tk.Token)
	}
}

func TestEnvironmentCredential_ClientSecretLive(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
	}
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
	tk, err := cred.GetToken(context.Background(), testTRO)
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
	tk, err := cred.GetToken(context.Background(), testTRO)
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
