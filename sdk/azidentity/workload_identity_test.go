//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

// createTestCA generates a self-signed CA certificate for testing
func createTestCA(t *testing.T) ([]byte, string) {
	t.Helper()

	// Generate private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test CA"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test City"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	// Encode certificate to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Write to temporary file
	caFile := filepath.Join(t.TempDir(), "ca.crt")
	err = os.WriteFile(caFile, certPEM, 0600)
	require.NoError(t, err)

	return certPEM, caFile
}

// customTokenEndpointTransportRecorder records HTTP requests made by the transport
type customTokenEndpointTransportRecorder struct {
	requests []*http.Request
	response *http.Response
	err      error
	mtx      sync.Mutex
}

func (r *customTokenEndpointTransportRecorder) RoundTrip(req *http.Request) (*http.Response, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.requests = append(r.requests, req)
	if r.err != nil {
		return nil, r.err
	}
	if r.response != nil {
		return r.response, nil
	}
	return &http.Response{
		StatusCode: 200,
		Body:       http.NoBody,
	}, nil
}

func (r *customTokenEndpointTransportRecorder) Do(req *http.Request) (*http.Response, error) {
	return r.RoundTrip(req)
}

func (r *customTokenEndpointTransportRecorder) getRequests() []*http.Request {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	// Return copy to avoid race conditions
	result := make([]*http.Request, len(r.requests))
	copy(result, r.requests)
	return result
}

func TestWorkloadIdentityCredential_IdentityBinding_Detection(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "all variables present - should configure identity binding",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "https://kubernetes.default.svc",
				azureKubernetesSNIName:       "test.cluster.local",
				azureKubernetesCAFile:        "valid-ca-file",
			},
			wantErr: false,
		},
		{
			name:    "no variables present - should use standard mode",
			envVars: map[string]string{},
			wantErr: false,
		},
		{
			name: "only token endpoint present - should error on missing CA file",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "https://kubernetes.default.svc",
			},
			wantErr:    true,
			wantErrMsg: "read CA file",
		},
		{
			name: "only SNI name present - should use standard mode",
			envVars: map[string]string{
				azureKubernetesSNIName: "test.cluster.local",
			},
			wantErr: false,
		},
		{
			name: "only CA file present - should use standard mode",
			envVars: map[string]string{
				azureKubernetesCAFile: "test-ca.crt",
			},
			wantErr: false,
		},
		{
			name: "two of three variables present - should error on missing CA file",
			envVars: map[string]string{
				azureKubernetesTokenEndpoint: "https://kubernetes.default.svc",
				azureKubernetesSNIName:       "test.cluster.local",
			},
			wantErr:    true,
			wantErrMsg: "read CA file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test CA file if needed
			var caFile string
			if tt.envVars[azureKubernetesCAFile] != "" {
				_, caFile = createTestCA(t)
				tt.envVars[azureKubernetesCAFile] = caFile
			}

			// Set up environment
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			// Set up required workload identity variables
			tempTokenFile := filepath.Join(t.TempDir(), "token")
			err := os.WriteFile(tempTokenFile, []byte("test-token"), 0600)
			require.NoError(t, err)

			t.Setenv(azureClientID, fakeClientID)
			t.Setenv(azureTenantID, fakeTenantID)
			t.Setenv(azureFederatedTokenFile, tempTokenFile)

			// Create credential
			cred, err := NewWorkloadIdentityCredential(nil)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg)
				require.Nil(t, cred)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cred)
			}
		})
	}
}

func TestWorkloadIdentityCredential_IdentityBinding_InvalidURL(t *testing.T) {
	_, caFile := createTestCA(t)

	t.Setenv(azureKubernetesTokenEndpoint, "invalid://url with spaces")
	t.Setenv(azureKubernetesSNIName, "test.cluster.local")
	t.Setenv(azureKubernetesCAFile, caFile)
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, filepath.Join(t.TempDir(), "token"))

	_, err := NewWorkloadIdentityCredential(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse token endpoint URL")
}

func TestWorkloadIdentityCredential_IdentityBinding_NonHTTPS(t *testing.T) {
	_, caFile := createTestCA(t)

	t.Setenv(azureKubernetesTokenEndpoint, "http://kubernetes.default.svc")
	t.Setenv(azureKubernetesSNIName, "test.cluster.local")
	t.Setenv(azureKubernetesCAFile, caFile)
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, filepath.Join(t.TempDir(), "token"))

	_, err := NewWorkloadIdentityCredential(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token endpoint must use https scheme")
}

func TestWorkloadIdentityCredential_IdentityBinding_InvalidCAFile(t *testing.T) {
	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc")
	t.Setenv(azureKubernetesSNIName, "test.cluster.local")
	t.Setenv(azureKubernetesCAFile, "/nonexistent/ca.crt")
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, filepath.Join(t.TempDir(), "token"))

	_, err := NewWorkloadIdentityCredential(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "read CA file")
}

func TestWorkloadIdentityCredential_IdentityBinding_InvalidCAContent(t *testing.T) {
	// Create invalid CA file
	caFile := filepath.Join(t.TempDir(), "invalid-ca.crt")
	err := os.WriteFile(caFile, []byte("not a valid certificate"), 0600)
	require.NoError(t, err)

	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc")
	t.Setenv(azureKubernetesSNIName, "test.cluster.local")
	t.Setenv(azureKubernetesCAFile, caFile)
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, filepath.Join(t.TempDir(), "token"))

	_, err = NewWorkloadIdentityCredential(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no valid certificates found")
}

func TestWorkloadIdentityCredential_IdentityBinding_TransportRedirection(t *testing.T) {
	_, caFile := createTestCA(t)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport directly
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc:443", recorder)
	require.NoError(t, err)

	// Test token request (should be redirected)
	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/tenant-id/oauth2/v2.0/token", nil)
	require.NoError(t, err)
	_, _ = transport.Do(req)

	// Since token requests go through the custom transport's internal logic and don't hit the fallback transport,
	// we can't easily verify the redirection through the recorder. Instead, we verify that the fallback transport
	// was NOT used (recorder should have no requests)
	requests := recorder.getRequests()
	require.Len(t, requests, 0, "Token requests should not go through fallback transport")

	// Test non-token request (should go through fallback transport)
	req2, err := http.NewRequest("GET", "https://login.microsoftonline.com/common/discovery/instance", nil)
	require.NoError(t, err)
	_, _ = transport.Do(req2)

	// Verify fallback transport was used for non-token request
	requests = recorder.getRequests()
	require.Len(t, requests, 1, "Non-token requests should go through fallback transport")
	require.Equal(t, "login.microsoftonline.com", requests[0].URL.Host)
}

func TestWorkloadIdentityCredential_IdentityBinding_NonTokenRequestPassthrough(t *testing.T) {
	_, caFile := createTestCA(t)

	t.Setenv(azureKubernetesTokenEndpoint, "https://kubernetes.default.svc")
	t.Setenv(azureKubernetesSNIName, "test.cluster.local")
	t.Setenv(azureKubernetesCAFile, caFile)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err)

	// Test non-token request (should pass through to fallback transport)
	req, _ := http.NewRequest("GET", "https://example.com/api/data", nil)
	_, _ = transport.Do(req)

	// Verify fallback transport was used
	requests := recorder.getRequests()
	require.Len(t, requests, 1)
	require.Equal(t, "example.com", requests[0].URL.Host)
}

func TestWorkloadIdentityCredential_IdentityBinding_TokenRequestRedirection(t *testing.T) {
	_, caFile := createTestCA(t)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	ibTransport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc:443", recorder)
	require.NoError(t, err)

	// Test token request (should be redirected)
	req, _ := http.NewRequest("POST", "https://login.microsoftonline.com/tenant/oauth2/v2.0/token", nil)
	_, _ = ibTransport.Do(req)

	// Verify request was redirected to Kubernetes endpoint
	// The request should go through the custom transport (internal), not the fallback transport
	// So we shouldn't see any requests in the recorder (fallback transport)
	requests := recorder.getRequests()
	require.Len(t, requests, 0) // No requests should go to fallback transport
}

func TestWorkloadIdentityCredential_IdentityBinding_CAReloading(t *testing.T) {
	tempDir := t.TempDir()
	caFile := filepath.Join(tempDir, "ca.crt")

	// Create initial CA
	initialCA, _ := createTestCA(t)
	err := os.WriteFile(caFile, initialCA, 0600)
	require.NoError(t, err)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err)

	// Verify initial CA was loaded
	originalTransport, err := transport.getTokenTransporter()
	require.NoError(t, err)
	require.NotNil(t, originalTransport)

	// Force reload by setting nextRead to past
	transport.mtx.Lock()
	transport.nextRead = time.Now().Add(-time.Hour)
	transport.mtx.Unlock()

	// Create new CA and update file
	newCA, _ := createTestCA(t)
	err = os.WriteFile(caFile, newCA, 0600)
	require.NoError(t, err)

	// Get transport again to trigger reload
	newTransport, err := transport.getTokenTransporter()
	require.NoError(t, err)
	require.NotNil(t, newTransport)

	// Verify that CA was reloaded (content should be different)
	transport.mtx.RLock()
	currentCA := transport.currentCA
	transport.mtx.RUnlock()
	require.Equal(t, newCA, currentCA)
}

func TestWorkloadIdentityCredential_IdentityBinding_EmptyCAFile(t *testing.T) {
	tempDir := t.TempDir()
	caFile := filepath.Join(tempDir, "empty-ca.crt")

	// Create empty CA file
	err := os.WriteFile(caFile, []byte(""), 0600)
	require.NoError(t, err)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport - empty files are handled gracefully during construction
	// The implementation allows empty files to handle CA file rotation scenarios
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err) // Empty file should not error during construction
	require.NotNil(t, transport)
}

func TestWorkloadIdentityCredential_IdentityBinding_CAFileRotation(t *testing.T) {
	tempDir := t.TempDir()
	caFile := filepath.Join(tempDir, "rotating-ca.crt")

	// Create initial CA
	initialCA, _ := createTestCA(t)
	err := os.WriteFile(caFile, initialCA, 0600)
	require.NoError(t, err)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err)

	// Simulate file rotation - write empty file first (common during rotation)
	err = os.WriteFile(caFile, []byte(""), 0600)
	require.NoError(t, err)

	// Force reload with empty file
	transport.mtx.Lock()
	transport.nextRead = time.Now().Add(-time.Hour)
	err = transport.reloadCA() // Should not error with empty file
	transport.mtx.Unlock()
	require.NoError(t, err) // Empty file should be handled gracefully (no-op)

	// Write new CA content
	newCA, _ := createTestCA(t)
	err = os.WriteFile(caFile, newCA, 0600)
	require.NoError(t, err)

	// Force another reload
	transport.mtx.Lock()
	transport.nextRead = time.Now().Add(-time.Hour)
	err = transport.reloadCA()
	transport.mtx.Unlock()
	require.NoError(t, err)

	// Verify new CA was loaded
	transport.mtx.RLock()
	currentCA := transport.currentCA
	transport.mtx.RUnlock()
	require.Equal(t, newCA, currentCA)
}

func TestWorkloadIdentityCredential_IdentityBinding_ConcurrentAccess(t *testing.T) {
	_, caFile := createTestCA(t)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err)

	// Test concurrent access to transport
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = transport.getTokenTransporter() // Should not panic or deadlock
		}()
	}

	wg.Wait() // Should complete without hanging
}

func TestWorkloadIdentityCredential_IdentityBinding_ConcurrentAccessWithCARotation(t *testing.T) {
	tempDir := t.TempDir()
	caFile := filepath.Join(tempDir, "ca.crt")

	// Create initial CA
	initialCA, _ := createTestCA(t)
	err := os.WriteFile(caFile, initialCA, 0600)
	require.NoError(t, err)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	transport, err := newIdentityBindingTransport(caFile, "test.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err)

	// Test concurrent access while rotating CA
	var wg sync.WaitGroup

	// Start concurrent readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, _ = transport.getTokenTransporter() // Should not panic or deadlock
				time.Sleep(time.Millisecond * 10)      // Small delay to create overlap
			}
		}()
	}

	// Start CA rotation goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for k := 0; k < 3; k++ {
			time.Sleep(time.Millisecond * 50)
			// Generate new CA and update file
			newCA, _ := createTestCA(t)
			err := os.WriteFile(caFile, newCA, 0600)
			if err != nil {
				return // File might be locked during concurrent access
			}
			// Force reload by updating nextRead time
			transport.mtx.Lock()
			transport.nextRead = time.Now().Add(-time.Hour)
			transport.mtx.Unlock()
		}
	}()

	wg.Wait() // Should complete without hanging or deadlock
}

func TestWorkloadIdentityCredential_IdentityBinding_TLSConfiguration(t *testing.T) {
	_, caFile := createTestCA(t)

	// Create transport recorder
	recorder := &customTokenEndpointTransportRecorder{}

	// Create identity binding transport
	transport, err := newIdentityBindingTransport(caFile, "custom-sni.cluster.local", "https://kubernetes.default.svc", recorder)
	require.NoError(t, err)

	// Get the configured transport
	tokenTransport, err := transport.getTokenTransporter()
	require.NoError(t, err)

	// Verify TLS configuration
	require.NotNil(t, tokenTransport.TLSClientConfig)
	require.Equal(t, "custom-sni.cluster.local", tokenTransport.TLSClientConfig.ServerName)
	require.NotNil(t, tokenTransport.TLSClientConfig.RootCAs)
}

func TestWorkloadIdentityCredential_IdentityBinding_BackwardCompatibility(t *testing.T) {
	// Test that standard workload identity still works when no identity binding variables are set
	tempTokenFile := filepath.Join(t.TempDir(), "token")
	err := os.WriteFile(tempTokenFile, []byte(tokenValue), 0600)
	require.NoError(t, err)

	// Set up standard workload identity environment (no binding variables)
	t.Setenv(azureClientID, fakeClientID)
	t.Setenv(azureTenantID, fakeTenantID)
	t.Setenv(azureFederatedTokenFile, tempTokenFile)

	// Create mock STS for standard requests
	sts := mockSTS{tenant: fakeTenantID, tokenRequestCallback: func(req *http.Request) *http.Response {
		// Should receive standard MSAL token request
		require.Contains(t, req.URL.Host, "login.microsoftonline.com")
		return nil
	}}

	// Create credential - should work in standard mode
	cred, err := NewWorkloadIdentityCredential(&WorkloadIdentityCredentialOptions{
		ClientOptions: policy.ClientOptions{Transport: &sts},
	})
	require.NoError(t, err)

	// Should be able to get token using standard flow
	testGetTokenSuccess(t, cred)
}
