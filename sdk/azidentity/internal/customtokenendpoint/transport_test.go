// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenendpoint

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

// errReadCloser is a helper that forces an error when reading the request body
type errReadCloser struct{}

func (e errReadCloser) Read(p []byte) (int, error) { return 0, fmt.Errorf("forced read error") }
func (e errReadCloser) Close() error               { return nil }

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

func TestConfigure(t *testing.T) {
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
				AzureKubernetesTokenEndpoint: "https://custom-endpoint.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint enabled with CA file + SNI",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				AzureKubernetesCAFile:        "/path/to/custom-ca-file",
				AzureKubernetesSNIName:       "custom-sni.example.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint enabled with CA data + SNI",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				AzureKubernetesCAData:        "custom-ca-data",
				AzureKubernetesSNIName:       "custom-sni.example.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint enabled with SNI",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				AzureKubernetesSNIName:       "custom-sni.example.com",
			},
			expectOverridesTransport: true,
		},
		{
			name:      "custom endpoint disabled with extra environment variables",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesSNIName: "custom-sni.example.com",
			},
			checkErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointEnvSetWithoutTokenEndpoint)
			},
		},
		{
			name:      "custom endpoint enabled with both CAData and CAFile",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesTokenEndpoint: "https://custom-endpoint.com",
				AzureKubernetesCAData:        "custom-ca-data",
				AzureKubernetesCAFile:        "/path/to/custom-ca-file",
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
				AzureKubernetesTokenEndpoint: "http://custom-endpoint.com",
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
			err := Configure(&c.clientOptions, http.DefaultClient)

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
