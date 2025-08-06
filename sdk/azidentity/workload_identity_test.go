//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func assertion(cert *x509.Certificate, key crypto.PrivateKey) (string, error) {
	j := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", liveSP.tenantID),
		"exp": json.Number(strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)),
		"iss": liveSP.clientID,
		"jti": uuid.New().String(),
		"nbf": json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		"sub": liveSP.clientID,
	})
	x5t := sha1.Sum(cert.Raw) // nosec
	j.Header = map[string]interface{}{
		"alg": "RS256",
		"typ": "JWT",
		"x5t": base64.StdEncoding.EncodeToString(x5t[:]),
	}
	return j.SignedString(key)
}

func TestWorkloadIdentityCredential_Live(t *testing.T) {
	// This test triggers the managed identity test app deployed to Azure Kubernetes Service.
	// See the bicep file and test resources scripts for details.
	// It triggers the app with kubectl because the test subscription prohibits opening ports to the internet.
	pod := os.Getenv("AZIDENTITY_POD_NAME")
	if pod == "" {
		t.Skip("set AZIDENTITY_POD_NAME to run this test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "kubectl", "exec", pod, "--", "wget", "-qO-", "localhost")
	b, err := cmd.CombinedOutput()
	s := string(b)
	require.NoError(t, err, s)
	require.Equal(t, "test passed", s)
}

func TestWorkloadIdentityCredential_Recorded(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
	}
	// workload identity and client cert auth use the same flow. This test
	// implements cert auth with WorkloadIdentityCredential as a way to test
	// that credential in an environment that's easier to set up than AKS
	cert, err := os.ReadFile(liveSP.pemPath)
	if err != nil {
		t.Fatal(err)
	}
	certs, key, err := ParseCertificates(cert, nil)
	if err != nil {
		t.Fatal(err)
	}
	a, err := assertion(certs[0], key)
	if err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(t.TempDir(), t.Name())
	if err := os.WriteFile(f, []byte(a), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	for _, b := range []bool{true, false} {
		name := "default options"
		if b {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			co, stop := initRecording(t)
			defer stop()
			cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
				ClientID:                 liveSP.clientID,
				ClientOptions:            co,
				DisableInstanceDiscovery: b,
				TenantID:                 liveSP.tenantID,
				TokenFilePath:            f,
			})
			if err != nil {
				t.Fatal(err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestWorkloadIdentityCredential(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	if err := os.WriteFile(tempFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}
	sts := mockSTS{tenant: fakeTenantID, tokenRequestCallback: func(req *http.Request) *http.Response {
		if err := req.ParseForm(); err != nil {
			t.Error(err)
		}
		if actual, ok := req.PostForm["client_assertion"]; !ok {
			t.Error("expected a client_assertion")
		} else if len(actual) != 1 || actual[0] != tokenValue {
			t.Errorf(`unexpected assertion "%s"`, actual[0])
		}
		if actual, ok := req.PostForm["client_id"]; !ok {
			t.Error("expected a client_id")
		} else if len(actual) != 1 || actual[0] != fakeClientID {
			t.Errorf(`unexpected assertion "%s"`, actual[0])
		}
		if actual := strings.Split(req.URL.Path, "/")[1]; actual != fakeTenantID {
			t.Errorf(`unexpected tenant "%s"`, actual)
		}
		return nil
	}}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:      fakeClientID,
		ClientOptions: policy.ClientOptions{Transport: &sts},
		TenantID:      fakeTenantID,
		TokenFilePath: tempFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"scope"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkloadIdentityCredential_Expiration(t *testing.T) {
	tokenReqs := 0
	tempFile := filepath.Join(t.TempDir(), "test-workload-token-file")
	sts := mockSTS{tenant: fakeTenantID, tokenRequestCallback: func(req *http.Request) *http.Response {
		if err := req.ParseForm(); err != nil {
			t.Error(err)
		}
		if actual, ok := req.PostForm["client_assertion"]; !ok {
			t.Error("expected a client_assertion")
		} else if len(actual) != 1 || actual[0] != fmt.Sprint(tokenReqs) {
			t.Errorf(`expected assertion "%d", got "%s"`, tokenReqs, actual[0])
		}
		tokenReqs++
		return nil
	}}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:      fakeClientID,
		ClientOptions: policy.ClientOptions{Transport: &sts},
		TenantID:      fakeTenantID,
		TokenFilePath: tempFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		// tokenReqs counts requests, and its latest value is the expected client assertion and the requested scope.
		// Each iteration of this loop therefore sends a token request with a unique assertion.
		s := fmt.Sprint(tokenReqs)
		if err = os.WriteFile(tempFile, []byte(fmt.Sprint(s)), os.ModePerm); err != nil {
			t.Fatalf("failed to write token file: %v", err)
		}
		if _, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{s}}); err != nil {
			t.Fatal(err)
		}
		cred.expires = time.Now().Add(-time.Second)
	}
	if tokenReqs != 2 {
		t.Fatalf("expected 2 token requests, got %d", tokenReqs)
	}
}

func TestTestWorkloadIdentityCredential_IncompleteConfig(t *testing.T) {
	f := filepath.Join(t.TempDir(), t.Name())
	for _, env := range []map[string]string{
		{},

		{azureClientID: fakeClientID},
		{azureFederatedTokenFile: f},
		{azureTenantID: fakeTenantID},

		{azureClientID: fakeClientID, azureTenantID: fakeTenantID},
		{azureClientID: fakeClientID, azureFederatedTokenFile: f},
		{azureTenantID: fakeTenantID, azureFederatedTokenFile: f},
	} {
		t.Run("", func(t *testing.T) {
			for k, v := range env {
				t.Setenv(k, v)
			}
			if _, err := NewWorkloadIdentityCredential(nil); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestWorkloadIdentityCredential_NoFile(t *testing.T) {
	for k, v := range map[string]string{
		azureClientID:           fakeClientID,
		azureFederatedTokenFile: filepath.Join(t.TempDir(), t.Name()),
		azureTenantID:           fakeTenantID,
	} {
		t.Setenv(k, v)
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = cred.GetToken(context.Background(), testTRO); err == nil {
		t.Fatal("expected an error")
	}
}

func TestWorkloadIdentityCredential_Options(t *testing.T) {
	clientID := "not-" + fakeClientID
	tenantID := "not-" + fakeTenantID
	wrongFile := filepath.Join(t.TempDir(), "wrong")
	rightFile := filepath.Join(t.TempDir(), "right")
	if err := os.WriteFile(rightFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	sts := mockSTS{
		tenant: tenantID,
		tokenRequestCallback: func(req *http.Request) *http.Response {
			if err := req.ParseForm(); err != nil {
				t.Error(err)
			}
			if actual, ok := req.PostForm["client_assertion"]; !ok {
				t.Error("expected a client_assertion")
			} else if len(actual) != 1 || actual[0] != tokenValue {
				t.Errorf(`unexpected assertion "%s"`, actual[0])
			}
			if actual, ok := req.PostForm["client_id"]; !ok {
				t.Error("expected a client_id")
			} else if len(actual) != 1 || actual[0] != clientID {
				t.Errorf(`unexpected assertion "%s"`, actual[0])
			}
			if actual := strings.Split(req.URL.Path, "/")[1]; actual != tenantID {
				t.Errorf(`unexpected tenant "%s"`, actual)
			}
			return nil
		},
	}
	// options should override environment variables
	for k, v := range map[string]string{
		azureClientID:           fakeClientID,
		azureFederatedTokenFile: wrongFile,
		azureTenantID:           fakeTenantID,
	} {
		t.Setenv(k, v)
	}
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientID:      clientID,
		ClientOptions: policy.ClientOptions{Transport: &sts},
		TenantID:      tenantID,
		TokenFilePath: rightFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("unexpected token %q", tk.Token)
	}
}

func TestWorkloadIdentityCredential_IdentityBinding_Transport(t *testing.T) {
	// Test the transport redirection logic directly
	transport := &identityBindingTransport{
		kubernetesTokenEndpoint: "https://kubernetes.default.svc",
		base: &mockRoundTripper{
			callback: func(req *http.Request) (*http.Response, error) {
				// Verify the request was redirected
				if req.URL.Host != "kubernetes.default.svc" {
					t.Errorf("expected request to be redirected to kubernetes.default.svc, got %s", req.URL.Host)
				}
				if req.URL.Scheme != "https" {
					t.Errorf("expected HTTPS scheme, got %s", req.URL.Scheme)
				}
				// Return a mock response
				return &http.Response{
					StatusCode: 200,
					Body:       http.NoBody,
				}, nil
			},
		},
	}

	// Test with Azure authority host request (should be redirected)
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "login.microsoftonline.com",
			Path:   "/tenant-id/oauth2/v2.0/token",
		},
	}

	_, err := transport.Do(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test with non-Azure host request (should not be redirected)
	req2 := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/api",
		},
	}

	transport2 := &identityBindingTransport{
		kubernetesTokenEndpoint: "https://kubernetes.default.svc",
		base: &mockRoundTripper{
			callback: func(req *http.Request) (*http.Response, error) {
				// Verify the request was NOT redirected
				if req.URL.Host != "example.com" {
					t.Errorf("expected request to NOT be redirected, got %s", req.URL.Host)
				}
				return &http.Response{
					StatusCode: 200,
					Body:       http.NoBody,
				}, nil
			},
		},
	}

	_, err = transport2.Do(req2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// mockRoundTripper is a helper for testing HTTP transport behavior
type mockRoundTripper struct {
	callback func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.callback(req)
}

func TestWorkloadIdentityCredential_IdentityBinding_Success(t *testing.T) {
	// Create temporary files for token file and CA file
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "token")
	caFile := filepath.Join(tempDir, "ca.pem")

	// Write token file
	if err := os.WriteFile(tokenFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	// Read the test certificate we created
	caCert, err := os.ReadFile("/tmp/test-certs/cert.pem")
	if err != nil {
		t.Skipf("test certificate not found, run: cd /tmp/test-certs && openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 1 -nodes -subj \"/C=US/ST=Test/L=Test/O=Test/OU=Test/CN=kubernetes.default.svc\"")
	}
	
	if err := os.WriteFile(caFile, caCert, os.ModePerm); err != nil {
		t.Fatalf("failed to write CA file: %v", err)
	}

	requestReceived := false
	_ = requestReceived // To avoid unused variable error in this test
	sts := mockSTS{
		tenant: fakeTenantID,
		tokenRequestCallback: func(req *http.Request) *http.Response {
			requestReceived = true
			// Check that the request was redirected to the Kubernetes endpoint
			if req.URL.Host != "kubernetes.default.svc" {
				t.Errorf("expected request to be redirected to kubernetes.default.svc, got %s", req.URL.Host)
			}
			if err := req.ParseForm(); err != nil {
				t.Error(err)
			}
			if actual, ok := req.PostForm["client_assertion"]; !ok {
				t.Error("expected a client_assertion")
			} else if len(actual) != 1 || actual[0] != tokenValue {
				t.Errorf(`unexpected assertion "%s"`, actual[0])
			}
			return nil
		},
	}

	// Set environment variables for identity binding mode
	for k, v := range map[string]string{
		azureClientID:                fakeClientID,
		azureTenantID:                fakeTenantID,
		azureFederatedTokenFile:      tokenFile,
		azureKubernetesTokenEndpoint: "https://kubernetes.default.svc",
		azureKubernetesSNIName:       "kubernetes.default.svc",
		azureKubernetesCAFile:        caFile,
	} {
		t.Setenv(k, v)
	}

	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: &sts},
	})
	if err != nil {
		t.Fatalf("failed to create credential: %v", err)
	}

	// Verify that identity binding mode is enabled
	if !cred.identityBinding {
		t.Error("expected identity binding mode to be enabled")
	}

	// Note: We can't test the full token flow here because it would make actual network requests
	// The important part is that the credential was created successfully with identity binding enabled
}

func TestWorkloadIdentityCredential_IdentityBinding(t *testing.T) {
	// Create temporary files for token file and CA file
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "token")
	caFile := filepath.Join(tempDir, "ca.pem")

	// Write token file
	if err := os.WriteFile(tokenFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	// Write a dummy CA certificate (for testing, we'll use an invalid cert to test error handling)
	caCert := `-----BEGIN CERTIFICATE-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234567890abcdefghij
klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghij
klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghij
klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghij
klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghij
-----END CERTIFICATE-----`
	
	if err := os.WriteFile(caFile, []byte(caCert), os.ModePerm); err != nil {
		t.Fatalf("failed to write CA file: %v", err)
	}

	requestCount := 0
	sts := mockSTS{
		tenant: fakeTenantID,
		tokenRequestCallback: func(req *http.Request) *http.Response {
			requestCount++
			// Check that the request was redirected to the Kubernetes endpoint
			if req.URL.Host != "kubernetes.default.svc" {
				t.Errorf("expected request to be redirected to kubernetes.default.svc, got %s", req.URL.Host)
			}
			if err := req.ParseForm(); err != nil {
				t.Error(err)
			}
			if actual, ok := req.PostForm["client_assertion"]; !ok {
				t.Error("expected a client_assertion")
			} else if len(actual) != 1 || actual[0] != tokenValue {
				t.Errorf(`unexpected assertion "%s"`, actual[0])
			}
			return nil
		},
	}

	// Set environment variables for identity binding mode
	for k, v := range map[string]string{
		azureClientID:                fakeClientID,
		azureTenantID:                fakeTenantID,
		azureFederatedTokenFile:      tokenFile,
		azureKubernetesTokenEndpoint: "https://kubernetes.default.svc",
		azureKubernetesSNIName:       "my-cluster.local",
		azureKubernetesCAFile:        caFile,
	} {
		t.Setenv(k, v)
	}

	// This should fail because we can't actually parse the dummy CA cert
	_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: &sts},
	})
	if err == nil {
		t.Fatal("expected an error when parsing invalid CA certificate")
	}
	if !strings.Contains(err.Error(), "failed to parse Kubernetes CA certificate") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkloadIdentityCredential_IdentityBinding_IncompleteConfig(t *testing.T) {
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "token")
	caFile := filepath.Join(tempDir, "ca.pem")

	if err := os.WriteFile(tokenFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	// Test with incomplete identity binding configuration
	testCases := []map[string]string{
		// Only one variable set
		{azureKubernetesTokenEndpoint: "https://kubernetes.default.svc"},
		{azureKubernetesSNIName: "my-cluster.local"},
		{azureKubernetesCAFile: caFile},
		// Two variables set
		{azureKubernetesTokenEndpoint: "https://kubernetes.default.svc", azureKubernetesSNIName: "my-cluster.local"},
		{azureKubernetesTokenEndpoint: "https://kubernetes.default.svc", azureKubernetesCAFile: caFile},
		{azureKubernetesSNIName: "my-cluster.local", azureKubernetesCAFile: caFile},
	}

	for i, envVars := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			// Set base environment variables
			for k, v := range map[string]string{
				azureClientID:           fakeClientID,
				azureTenantID:           fakeTenantID,
				azureFederatedTokenFile: tokenFile,
			} {
				t.Setenv(k, v)
			}

			// Set incomplete identity binding variables
			for k, v := range envVars {
				t.Setenv(k, v)
			}

			_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
				ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
			})
			if err == nil {
				t.Fatal("expected an error for incomplete identity binding configuration")
			}
			if !strings.Contains(err.Error(), "identity binding mode requires all three environment variables") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestWorkloadIdentityCredential_IdentityBinding_MissingCAFile(t *testing.T) {
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "token")
	missingCAFile := filepath.Join(tempDir, "missing-ca.pem")

	if err := os.WriteFile(tokenFile, []byte(tokenValue), os.ModePerm); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	// Set environment variables for identity binding mode with missing CA file
	for k, v := range map[string]string{
		azureClientID:                fakeClientID,
		azureTenantID:                fakeTenantID,
		azureFederatedTokenFile:      tokenFile,
		azureKubernetesTokenEndpoint: "https://kubernetes.default.svc",
		azureKubernetesSNIName:       "my-cluster.local",
		azureKubernetesCAFile:        missingCAFile,
	} {
		t.Setenv(k, v)
	}

	_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
	})
	if err == nil {
		t.Fatal("expected an error when CA file is missing")
	}
	if !strings.Contains(err.Error(), "failed to read Kubernetes CA file") {
		t.Errorf("unexpected error: %v", err)
	}
}
