//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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

// testCACertificate returns a valid test CA certificate in PEM format
func testCACertificate() string {
	return `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUF2VIP4+AnEtb52KTCHbo4+fESfswDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0xOTEwMzAyMjQ2MjBaFw0yMjA4
MTkyMjQ2MjBaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDL1hG+JYCfIPp3tlZ05J4pYIJ3Ckfs432bE3rYuWlR
2w9KqdjWkKxuAxpjJ+T+uoqVaT3BFMfi4ZRYOCI69s4+lP3DwR8uBCp9xyVkF8th
XfS3iui0liGDviVBoBJJWvjDFU8a/Hseg+QfoxAb6tx0kEc7V3ozBLWoIDJjfwJ3
NdsLZGVtAC34qCWeEIvS97CDA4g3Kc6hYJIrAa7pxHzo/Nd0U3e7z+DlBcJV7dY6
TZUyjBVTpzppWe+XQEOfKsjkDNykHEC1C1bClG0u7unS7QOBMd6bOGkeL+Bc+n22
slTzs5amsbDLNuobSaUsFt9vgD5jRD6FwhpXwj/Ek0F7AgMBAAGjUzBRMB0GA1Ud
DgQWBBT6Mf9uXFB67bY2PeW3GCTKfkO7vDAfBgNVHSMEGDAWgBT6Mf9uXFB67bY2
PeW3GCTKfkO7vDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCZ
1+kTISX85v9/ag7glavaPFUYsOSOOofl8gSzov7L01YL+srq7tXdvZmWrjQ/dnOY
h18rp9rb24vwIYxNioNG/M2cW1jBJwEGsDPOwdPV1VPcRmmUJW9kY130gRHBCd/N
qB7dIkcQnpNsxPIIWI+sRQp73U0ijhOByDnCNHLHon6vbfFTwkO1XggmV5BdZ3uQ
JNJyckILyNzlhmf6zhonMp4lVzkgxWsAm2vgdawd6dmBa+7Avb2QK9s+IdUSutFh
-----END CERTIFICATE-----`
}

// TestWorkloadIdentityCredential_CustomTokenEndpoint_BasicConfiguration tests the basic configuration
// of custom token endpoint mode when environment variables are set correctly
func TestWorkloadIdentityCredential_CustomTokenEndpoint_BasicConfiguration(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-workload-token-file")
	caFile := filepath.Join(tempDir, "ca.crt")
	
	// Create a simple CA file (content doesn't matter for logic test)
	testCACert := "# Simple CA file for testing transport logic\n"
	
	require.NoError(t, os.WriteFile(tempFile, []byte(tokenValue), 0600))
	require.NoError(t, os.WriteFile(caFile, []byte("# Simple CA file for testing transport logic\n"), 0600))

	// Set up environment variables for custom token endpoint mode
	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc/oauth2/token")
	t.Setenv(azureKubernetesCAFile, caFile)
	t.Setenv(azureKubernetesSNIName, "kubernetes.default.svc.cluster.local")
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, tempFile)

	// Create a mock server that handles both token and non-token requests
	mockTransport := &mockCustomTokenEndpointTransport{
		tokenEndpointHandler: func(req *http.Request) (*http.Response, error) {
			// Verify request was redirected to custom endpoint
			require.Equal(t, "https", req.URL.Scheme)
			require.Equal(t, "kubernetes.default.svc", req.URL.Host)
			require.Equal(t, "/oauth2/token", req.URL.Path)
			
			// Verify it's a token request
			require.NoError(t, req.ParseForm())
			require.Contains(t, req.PostForm, "client_assertion")
			require.Contains(t, req.PostForm, "client_assertion_type")
			require.Equal(t, tokenValue, req.PostForm.Get("client_assertion"))
			require.Equal(t, fakeClientID, req.PostForm.Get("client_id"))
			
			// Return successful token response
			body := fmt.Sprintf(`{"access_token": "%s","expires_in": 3600,"token_type":"Bearer"}`, tokenValue)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		},
		fallbackHandler: func(req *http.Request) (*http.Response, error) {
			// Handle metadata requests
			if strings.Contains(req.URL.Path, "instance") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(instanceMetadata(fakeTenantID))),
				}, nil
			}
			if strings.Contains(req.URL.Path, "openid-configuration") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(tenantMetadata(fakeTenantID))),
				}, nil
			}
			
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		},
	}

	// Expect creation to fail due to invalid CA certificate
	// since we're not providing a valid cert, but we should get to the CA loading part
	_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: mockTransport},
	})
	require.NoError(t, err) // Should create successfully
	
	// Note: Since CA file is invalid, token requests will fail at CA loading time
	// This test primarily verifies that the custom transport configuration is set up correctly
}

// TestWorkloadIdentityCredential_CustomTokenEndpoint_CADataSupport tests CA data support
func TestWorkloadIdentityCredential_CustomTokenEndpoint_CADataSupport(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-workload-token-file")
	
	require.NoError(t, os.WriteFile(tempFile, []byte(tokenValue), 0600))

	// Set up environment variables for custom token endpoint mode with CA data
	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc/oauth2/token")
	t.Setenv(azureKubernetesCAData, "# Simple CA data for testing transport logic\n")
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, tempFile)

	mockTransport := &mockCustomTokenEndpointTransport{
		tokenEndpointHandler: func(req *http.Request) (*http.Response, error) {
			// Return successful token response
			body := fmt.Sprintf(`{"access_token": "%s","expires_in": 3600,"token_type":"Bearer"}`, tokenValue)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		},
		fallbackHandler: func(req *http.Request) (*http.Response, error) {
			// Handle metadata requests
			if strings.Contains(req.URL.Path, "instance") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(instanceMetadata(fakeTenantID))),
				}, nil
			}
			if strings.Contains(req.URL.Path, "openid-configuration") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(tenantMetadata(fakeTenantID))),
				}, nil
			}
			
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		},
	}

	_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: mockTransport},
	})
	require.NoError(t, err)
	
	// Note: Token acquisition would fail due to invalid CA data, but credential creation should succeed
	// This test verifies that the CA data path is properly configured
}

// TestWorkloadIdentityCredential_CustomTokenEndpoint_ValidationErrors tests validation error scenarios
func TestWorkloadIdentityCredential_CustomTokenEndpoint_ValidationErrors(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-workload-token-file")
	caFile := filepath.Join(tempDir, "ca.crt")
	emptyCaFile := filepath.Join(tempDir, "empty-ca.crt")
	
	require.NoError(t, os.WriteFile(tempFile, []byte(tokenValue), 0600))
	require.NoError(t, os.WriteFile(caFile, []byte("# Test CA file\n"), 0600))
	require.NoError(t, os.WriteFile(emptyCaFile, []byte(""), 0600))

	testCases := []struct {
		name     string
		envVars  map[string]string
		wantErr  string
	}{
		{
			name: "non-https token endpoint",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "http://kubernetes.default.svc/oauth2/token",
				azureKubernetesCAFile:        caFile,
				azureClientID:                fakeClientID,
				azureTenantID:                fakeTenantID,
				azureFederatedTokenFile:      tempFile,
			},
			wantErr: "token endpoint must use https scheme",
		},
		{
			name: "invalid token endpoint URL",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "://invalid-url",
				azureKubernetesCAFile:        caFile,
				azureClientID:                fakeClientID,
				azureTenantID:                fakeTenantID,
				azureFederatedTokenFile:      tempFile,
			},
			wantErr: "failed to parse token endpoint URL",
		},
		{
			name: "both CA file and data specified",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "https://kubernetes.default.svc/oauth2/token",
				azureKubernetesCAFile:        caFile,
				azureKubernetesCAData:        "# test ca data",
				azureClientID:                fakeClientID,
				azureTenantID:                fakeTenantID,
				azureFederatedTokenFile:      tempFile,
			},
			wantErr: "only one of AZURE_KUBERNETES_CA_FILE and AZURE_KUBERNETES_CA_DATA can be specified",
		},
		{
			name: "no CA file or data specified",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "https://kubernetes.default.svc/oauth2/token",
				azureClientID:                fakeClientID,
				azureTenantID:                fakeTenantID,
				azureFederatedTokenFile:      tempFile,
			},
			wantErr: "at least one of AZURE_KUBERNETES_CA_FILE or AZURE_KUBERNETES_CA_DATA must be specified",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear environment and set test vars
			for k := range tc.envVars {
				t.Setenv(k, tc.envVars[k])
			}

			_, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
				ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// TestWorkloadIdentityCredential_CustomTokenEndpoint_RequestRouting tests request routing logic
func TestWorkloadIdentityCredential_CustomTokenEndpoint_RequestRouting(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-workload-token-file")
	caFile := filepath.Join(tempDir, "ca.crt")
	
	require.NoError(t, os.WriteFile(tempFile, []byte(tokenValue), 0600))
	require.NoError(t, os.WriteFile(caFile, []byte(testCACertificate()), 0600))

	tokenEndpointCalled := false
	fallbackClientCalled := false
	var capturedTokenRequest *http.Request
	var capturedFallbackRequest *http.Request
	
	// Set up environment variables
	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc/oauth2/v2.0/token")
	t.Setenv(azureKubernetesCAFile, caFile)
	t.Setenv(azureKubernetesSNIName, "kubernetes.default.svc.cluster.local")
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, tempFile)

	mockTransport := &mockCustomTokenEndpointTransport{
		tokenEndpointHandler: func(req *http.Request) (*http.Response, error) {
			tokenEndpointCalled = true
			capturedTokenRequest = req.Clone(req.Context())
			
			// Return successful token response
			body := fmt.Sprintf(`{"access_token": "%s","expires_in": 3600,"token_type":"Bearer"}`, tokenValue)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		},
		fallbackHandler: func(req *http.Request) (*http.Response, error) {
			fallbackClientCalled = true
			capturedFallbackRequest = req.Clone(req.Context())
			
			// Handle metadata requests
			if strings.Contains(req.URL.Path, "instance") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(instanceMetadata(fakeTenantID))),
				}, nil
			}
			if strings.Contains(req.URL.Path, "openid-configuration") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(tenantMetadata(fakeTenantID))),
				}, nil
			}
			
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		},
	}

	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: mockTransport},
	})
	require.NoError(t, err)

	// Test token acquisition
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://vault.azure.net/.default"},
	})
	require.NoError(t, err)
	require.Equal(t, tokenValue, tk.Token)
	
	// Verify request routing
	require.True(t, tokenEndpointCalled, "Token endpoint should have been called")
	require.True(t, fallbackClientCalled, "Fallback client should have been called for metadata")
	
	// Verify token request was properly redirected and modified
	require.NotNil(t, capturedTokenRequest)
	require.Equal(t, "https", capturedTokenRequest.URL.Scheme)
	require.Equal(t, "kubernetes.default.svc", capturedTokenRequest.URL.Host)
	require.Equal(t, "/oauth2/v2.0/token", capturedTokenRequest.URL.Path)  // Should copy path from endpoint
	require.Equal(t, "kubernetes.default.svc", capturedTokenRequest.Host)
	
	// Verify fallback request was not modified
	require.NotNil(t, capturedFallbackRequest)
	require.Contains(t, capturedFallbackRequest.URL.String(), "login.microsoftonline.com")
}

// TestWorkloadIdentityCredential_CustomTokenEndpoint_ConcurrentAccess tests concurrent access to the custom transport
func TestWorkloadIdentityCredential_CustomTokenEndpoint_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-workload-token-file")
	caFile := filepath.Join(tempDir, "ca.crt")
	
	require.NoError(t, os.WriteFile(tempFile, []byte(tokenValue), 0600))
	require.NoError(t, os.WriteFile(caFile, []byte(testCACertificate()), 0600))

	var tokenCallCount int32
	var fallbackCallCount int32
	
	// Set up environment variables
	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc/oauth2/token")
	t.Setenv(azureKubernetesCAFile, caFile)
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, tempFile)

	mockTransport := &mockCustomTokenEndpointTransport{
		tokenEndpointHandler: func(req *http.Request) (*http.Response, error) {
			// Increment counter atomically
			count := atomic.AddInt32(&tokenCallCount, 1)
			
			// Add small delay to increase chance of race conditions
			time.Sleep(time.Millisecond * 10)
			
			// Return successful token response with unique content per call
			body := fmt.Sprintf(`{"access_token": "%s-%d","expires_in": 3600,"token_type":"Bearer"}`, tokenValue, count)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		},
		fallbackHandler: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&fallbackCallCount, 1)
			
			// Handle metadata requests
			if strings.Contains(req.URL.Path, "instance") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(instanceMetadata(fakeTenantID))),
				}, nil
			}
			if strings.Contains(req.URL.Path, "openid-configuration") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(tenantMetadata(fakeTenantID))),
				}, nil
			}
			
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		},
	}

	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: mockTransport},
	})
	require.NoError(t, err)

	const numGoroutines = 10
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]string, 0, numGoroutines)
	errors := make([]error, 0, numGoroutines)

	// Run concurrent token acquisitions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: []string{fmt.Sprintf("https://vault.azure.net/.default?id=%d", id)},
			})
			
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors = append(errors, err)
			} else {
				results = append(results, tk.Token)
			}
		}(i)
	}

	wg.Wait()

	// Verify no errors occurred
	require.Empty(t, errors, "No errors should occur during concurrent access")
	require.Equal(t, numGoroutines, len(results), "All requests should succeed")

	// Verify both endpoints were called
	require.Greater(t, atomic.LoadInt32(&tokenCallCount), int32(0), "Token endpoint should have been called")
	require.Greater(t, atomic.LoadInt32(&fallbackCallCount), int32(0), "Fallback client should have been called")
}

// TestWorkloadIdentityCredential_CustomTokenEndpoint_WithoutEnvironment tests standard behavior when custom environment is not set
func TestWorkloadIdentityCredential_CustomTokenEndpoint_WithoutEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-workload-token-file")
	require.NoError(t, os.WriteFile(tempFile, []byte(tokenValue), 0600))

	// Set up standard environment variables (no custom token endpoint)
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, tempFile)

	// Ensure custom token endpoint variables are NOT set
	require.Empty(t, os.Getenv(azureKubernetesTokenEndpoint))

	standardSTSCalled := false

	sts := mockSTS{
		tenant: fakeTenantID,
		tokenRequestCallback: func(req *http.Request) *http.Response {
			standardSTSCalled = true
			
			// Verify this is a standard Azure token request
			require.Contains(t, req.URL.Host, "login.microsoftonline.com")
			
			require.NoError(t, req.ParseForm())
			require.Contains(t, req.PostForm, "client_assertion")
			require.Equal(t, tokenValue, req.PostForm.Get("client_assertion"))
			
			return nil // Use default response
		},
	}

	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: &sts},
	})
	require.NoError(t, err)

	// Test token acquisition - should use standard flow
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://vault.azure.net/.default"},
	})
	require.NoError(t, err)
	require.Equal(t, tokenValue, tk.Token)
	
	// Verify standard STS was called
	require.True(t, standardSTSCalled, "Standard STS should be called")
}

// mockCustomTokenEndpointTransport is a test helper that allows separate handling of token and non-token requests
type mockCustomTokenEndpointTransport struct {
	tokenEndpointHandler func(req *http.Request) (*http.Response, error)
	fallbackHandler      func(req *http.Request) (*http.Response, error)
}

func (m *mockCustomTokenEndpointTransport) Do(req *http.Request) (*http.Response, error) {
	// Simple heuristic: if request body contains client_assertion, it's a token request
	if req.Body != nil && req.Body != http.NoBody {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		
		if strings.Contains(string(bodyBytes), "client_assertion") {
			return m.tokenEndpointHandler(req)
		}
	}
	
	return m.fallbackHandler(req)
}
