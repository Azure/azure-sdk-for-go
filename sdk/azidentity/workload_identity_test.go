//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
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
	"io"
	"math/big"
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

// errReadCloser is a helper that forces an error when reading the request body
type errReadCloser struct{}

func (e errReadCloser) Read(p []byte) (int, error) { return 0, fmt.Errorf("forced read error") }
func (e errReadCloser) Close() error               { return nil }

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

func TestParseAndValidateCustomTokenEndpoint(t *testing.T) {
	cases := []struct {
		name      string
		endpoint  string
		expectErr bool
		check     func(t testing.TB, u *url.URL, err error)
	}{
		{
			name:     "valid https endpoint without path",
			endpoint: "https://example.com",
			check: func(t testing.TB, u *url.URL, err error) {
				require.NoError(t, err)
				require.Equal(t, "https", u.Scheme)
				require.Equal(t, "example.com", u.Host)
				require.Equal(t, "", u.RawQuery)
				require.Equal(t, "", u.Fragment)
			},
		},
		{
			name:     "valid https endpoint with path",
			endpoint: "https://example.com/token/path",
			check: func(t testing.TB, u *url.URL, err error) {
				require.NoError(t, err)
				require.Equal(t, "/token/path", u.Path)
			},
		},
		{
			name:      "reject non-https scheme",
			endpoint:  "http://example.com",
			expectErr: true,
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "https scheme")
			},
		},
		{
			name:      "reject user info",
			endpoint:  "https://user:pass@example.com/token",
			expectErr: true,
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "must not contain user info")
			},
		},
		{
			name:      "reject query params",
			endpoint:  "https://example.com/token?foo=bar",
			expectErr: true,
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "must not contain a query")
			},
		},
		{
			name:      "reject fragment",
			endpoint:  "https://example.com/token#frag",
			expectErr: true,
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "must not contain a fragment")
			},
		},
		{
			name:      "reject unparseable URL",
			endpoint:  "https://example.com/%zz",
			expectErr: true,
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to parse custom token endpoint URL")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u, err := parseAndValidateCustomTokenEndpoint(c.endpoint)
			if c.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, u)
			}
			if c.check != nil {
				c.check(t, u, err)
			}
		})
	}
}

func TestConfigureCustomTokenEndpoint(t *testing.T) {
	cases := []struct {
		name          string
		envs          map[string]string
		clientOptions policy.ClientOptions

		expectErr                bool
		checkErr                 func(t testing.TB, err error) // optional check on error
		expectOverridesTransport bool
	}{
		{
			name:                     "no custom endpoint",
			expectErr:                false,
			expectOverridesTransport: false,
		},
		{
			name:      "custom endpoint enabled with minimal settings",
			expectErr: false,
			envs: map[string]string{
				azureKubernetesTokenEndpoint: "https://custom-endpoint.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint enabled with CA file + SNI",
			expectErr: false,
			envs: map[string]string{
				azureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				azureKubernetesCAFile:        "/path/to/custom-ca-file",
				azureKubernetesSNIName:       "custom-sni.example.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint enabled with CA data + SNI",
			expectErr: false,
			envs: map[string]string{
				azureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				azureKubernetesCAData:        "custom-ca-data",
				azureKubernetesSNIName:       "custom-sni.example.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint enabled with SNI",
			expectErr: false,
			envs: map[string]string{
				azureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				azureKubernetesSNIName:       "custom-sni.example.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint disabled with extra environment variables",
			expectErr: true,
			envs: map[string]string{
				azureKubernetesSNIName: "custom-sni.example.com",
			},
			checkErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointEnvSetWithoutTokenEndpoint)
			},
		},
		{
			name:      "custom endpoint enabled with both CAData and CAFile",
			expectErr: true,
			envs: map[string]string{
				azureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				azureKubernetesCAData:        "custom-ca-data",
				azureKubernetesCAFile:        "/path/to/custom-ca-file",
			},
			checkErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointMultipleCASourcesSet)
			},
		},
		{
			name:      "custom endpoint enabled with invalid endpoint",
			expectErr: true,
			envs: map[string]string{
				// http endpoint is not allowed
				azureKubernetesTokenEndpoint: "http://custom-endpoint.com",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.envs) > 0 {
				for k, v := range c.envs {
					t.Setenv(k, v)
				}
			}
			err := configureCustomTokenEndpoint(&c.clientOptions)

			if c.expectErr {
				require.Error(t, err)
				if c.checkErr != nil {
					c.checkErr(t, err)
				}
				return
			}

			require.NoError(t, err)
			if c.expectOverridesTransport {
				require.NotNil(t, c.clientOptions.Transport)
				require.IsType(t, &customTokenEndpointTransport{}, c.clientOptions.Transport)
			}
		})
	}
}

// createTestCA creates a valid CA as bytes
func createTestCA(t testing.TB) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

// createTestCAFile creates a valid CA file in a temporary directory.
// It returns the path to the CA file.
func createTestCAFile(t testing.TB) string {
	t.Helper()
	caData := createTestCA(t)
	tempDir := t.TempDir()
	caFile := filepath.Join(tempDir, "test-ca.pem")
	if err := os.WriteFile(caFile, caData, 0600); err != nil {
		t.Fatalf("failed to write CA file: %v", err)
	}
	return caFile
}

func TestCustomTokenEndpointTransport_loadCAPool(t *testing.T) {
	cases := []struct {
		name string
		tr   *customTokenEndpointTransport

		expectErr            bool
		expectCustomCertPool bool
	}{
		{
			name:                 "fallback to system pool",
			tr:                   &customTokenEndpointTransport{},
			expectErr:            false,
			expectCustomCertPool: false,
		},
		{
			name:      "load CA from CA file - invalid CA file",
			tr:        &customTokenEndpointTransport{caFile: "/path/to/invalid/ca/file"},
			expectErr: true,
		},
		{
			name: "load CA from CA file - invalid CA file content",
			tr: &customTokenEndpointTransport{
				caFile: func() string {
					d := t.TempDir()
					caFile := filepath.Join(d, "test-ca.pem")
					err := os.WriteFile(caFile, []byte("invalid-ca-content"), 0600)
					require.NoError(t, err)

					return caFile
				}(),
			},
			expectErr: true,
		},
		{
			name: "load CA from CA file - empty CA file content",
			tr: &customTokenEndpointTransport{
				caFile: func() string {
					d := t.TempDir()
					caFile := filepath.Join(d, "test-ca.pem")
					err := os.WriteFile(caFile, []byte(""), 0600)
					require.NoError(t, err)

					return caFile
				}(),
			},
			expectErr: true,
		},
		{
			name:                 "load CA from CA file - valid CA",
			tr:                   &customTokenEndpointTransport{caFile: createTestCAFile(t)},
			expectErr:            false,
			expectCustomCertPool: true,
		},
		{
			name:      "load CA from CA data - invalid CA data",
			tr:        &customTokenEndpointTransport{caData: "invalid-ca"},
			expectErr: true,
		},
		{
			name:                 "load CA from CA data - valid CA",
			tr:                   &customTokenEndpointTransport{caData: string(createTestCA(t))},
			expectErr:            false,
			expectCustomCertPool: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caPool, err := c.tr.loadCAPool()
			if c.expectErr {
				require.Error(t, err)
				require.Nil(t, caPool)
				return
			}

			require.NoError(t, err)
			if c.expectCustomCertPool {
				require.NotNil(t, caPool)
			}
		})
	}
}

func TestCustomTokenEndpointTransport_getTokenTransporter(t *testing.T) {
	cases := []struct {
		name string
		tr   *customTokenEndpointTransport

		expectErr         bool
		validateTransport func(t testing.TB, httpTr *http.Transport)
	}{
		{
			name: "no overrides",
			tr: &customTokenEndpointTransport{
				baseTransport: &http.Transport{},
			},
			expectErr: false,
		},
		{
			name: "with custom CA",
			tr: &customTokenEndpointTransport{
				baseTransport: &http.Transport{},
				caFile:        createTestCAFile(t),
			},
			expectErr: false,
			validateTransport: func(t testing.TB, httpTr *http.Transport) {
				require.NotNil(t, httpTr.TLSClientConfig)
				require.NotNil(t, httpTr.TLSClientConfig.RootCAs)
			},
		},
		{
			name: "invalid CA",
			tr: &customTokenEndpointTransport{
				baseTransport: &http.Transport{},
				caData:        "invalid-ca-data",
			},
			expectErr: true,
		},
		{
			name: "with SNI",
			tr: &customTokenEndpointTransport{
				baseTransport: &http.Transport{},
				sniName:       "example.com",
			},
			expectErr: false,
			validateTransport: func(t testing.TB, httpTr *http.Transport) {
				require.NotNil(t, httpTr.TLSClientConfig)
				require.NotEmpty(t, httpTr.TLSClientConfig.ServerName)
				require.Equal(t, "example.com", httpTr.TLSClientConfig.ServerName)
			},
		},
		{
			name: "with CA + SNI",
			tr: &customTokenEndpointTransport{
				baseTransport: &http.Transport{},
				sniName:       "example.com",
				caFile:        createTestCAFile(t),
			},
			expectErr: false,
			validateTransport: func(t testing.TB, httpTr *http.Transport) {
				require.NotNil(t, httpTr.TLSClientConfig)
				require.NotNil(t, httpTr.TLSClientConfig.RootCAs)
				require.NotEmpty(t, httpTr.TLSClientConfig.ServerName)
				require.Equal(t, "example.com", httpTr.TLSClientConfig.ServerName)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			transport, err := c.tr.getTokenTransporter()
			if c.expectErr {
				require.Error(t, err)
				require.Nil(t, transport)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, transport)
			if c.validateTransport != nil {
				c.validateTransport(t, transport)
			}
		})
	}
}

func TestCustomTokenEndpointTransport_isTokenRequest(t *testing.T) {
	// helper to build a request with given parameters
	newReq := func(method, contentType string, body io.Reader) *http.Request {
		t.Helper()

		if method == "" {
			method = http.MethodPost
		}
		if body == nil {
			// http.NewRequest with nil body sets ContentLength to 0
			req, err := http.NewRequest(method, "https://example.com/token", nil)
			require.NoError(t, err)
			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}
			return req
		}

		req, err := http.NewRequest(method, "https://example.com/token", body)
		require.NoError(t, err)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		// ensure ContentLength is positive when we know the body size
		switch r := body.(type) {
		case *bytes.Reader:
			req.ContentLength = int64(r.Len())
		case *strings.Reader:
			req.ContentLength = int64(r.Len())
		case *bytes.Buffer:
			req.ContentLength = int64(r.Len())
		}
		return req
	}

	newPostFormReq := func(body string) *http.Request {
		return newReq(
			http.MethodPost,
			"application/x-www-form-urlencoded",
			strings.NewReader(body),
		)
	}

	cases := []struct {
		name string
		req  *http.Request

		expectErr         bool
		expected          bool
		checkInspectedReq func(t testing.TB, req *http.Request)
	}{
		{
			name:     "non-POST method",
			req:      newReq(http.MethodGet, "application/x-www-form-urlencoded", strings.NewReader("a=b")),
			expected: false,
		},
		{
			name:     "missing Content-Type",
			req:      newReq(http.MethodPost, "", strings.NewReader("a=b")),
			expected: false,
		},
		{
			name:     "non form Content-Type",
			req:      newReq(http.MethodPost, "application/json", strings.NewReader("{}")),
			expected: false,
		},
		{
			name:     "nil body",
			req:      newReq(http.MethodPost, "application/x-www-form-urlencoded", nil),
			expected: false,
		},
		{
			name: "http.NoBody with non-empty content length",
			req: func() *http.Request {
				r := newPostFormReq("a=b")
				r.Body = http.NoBody
				r.ContentLength = 10
				return r
			}(),
			expected: false,
		},
		{
			name: "unreadable body returns error",
			req: func() *http.Request {
				r := newPostFormReq("a=b")
				r.Body = errReadCloser{}
				r.ContentLength = 10
				return r
			}(),
			expectErr: true,
		},
		{
			name:     "invalid form body (unparseable)",
			req:      newPostFormReq("%zz"),
			expected: false,
		},
		{
			name:     "missing client_assertion_type",
			req:      newPostFormReq("client_assertion=abc"),
			expected: false,
		},
		{
			name:     "missing client_assertion",
			req:      newPostFormReq("client_assertion_type=type"),
			expected: false,
		},
		{
			name:     "valid token request",
			req:      newPostFormReq("client_assertion=abc&client_assertion_type=type"),
			expected: true,
			checkInspectedReq: func(t testing.TB, req *http.Request) {
				// Body should be reset and readable
				b1, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				require.Equal(t, "client_assertion=abc&client_assertion_type=type", string(b1))
				// GetBody should be set and return the same bytes
				rc, err := req.GetBody()
				require.NoError(t, err)
				b2, err := io.ReadAll(rc)
				require.NoError(t, err)
				require.Equal(t, string(b1), string(b2))
			},
		},
		{
			name: "valid token request with charset",
			req: newReq(
				http.MethodPost,
				"application/x-www-form-urlencoded; charset=utf-8",
				strings.NewReader("client_assertion=abc&client_assertion_type=type"),
			),
			expected: true,
		},
		{
			name: "lowercase method post",
			req: newReq(
				"post",
				"application/x-www-form-urlencoded",
				strings.NewReader("client_assertion=abc&client_assertion_type=type"),
			),
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tr := &customTokenEndpointTransport{}
			v, err := tr.isTokenRequest(c.req)

			if c.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.expected, v)
			if c.checkInspectedReq != nil {
				c.checkInspectedReq(t, c.req)
			}
		})
	}
}
