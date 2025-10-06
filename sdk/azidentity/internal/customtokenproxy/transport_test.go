// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenproxy

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

func TestParseAndValidateCustomTokenProxy(t *testing.T) {
	cases := []struct {
		name     string
		endpoint string
		check    func(t testing.TB, u *url.URL, err error)
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
				require.Equal(t, "/", u.Path, "should set path to '/' if not present")
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
			name:     "reject non-https scheme",
			endpoint: "http://example.com",
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "https scheme")
			},
		},
		{
			name:     "reject user info",
			endpoint: "https://user:pass@example.com/token",
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "must not contain user info")
			},
		},
		{
			name:     "reject query params",
			endpoint: "https://example.com/token?foo=bar",
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "must not contain a query")
			},
		},
		{
			name:     "reject fragment",
			endpoint: "https://example.com/token#frag",
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "must not contain a fragment")
			},
		},
		{
			name:     "reject unparseable URL",
			endpoint: "https://example.com/%zz",
			check: func(t testing.TB, _ *url.URL, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to parse custom token proxy URL")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u, err := parseAndValidateCustomTokenProxy(c.endpoint)
			c.check(t, u, err)
		})
	}
}

func TestConfigure(t *testing.T) {
	var (
		testCAData = string(createTestCA(t))
		testCAFile = createTestCAFile(t)
	)

	cases := []struct {
		name          string
		envs          map[string]string
		clientOptions policy.ClientOptions

		expectErr       bool
		checkErr        func(t testing.TB, err error) // optional check on error
		expectTransport bool
	}{
		{
			name:            "no custom endpoint",
			expectErr:       false,
			expectTransport: false,
		},
		{
			name:      "custom endpoint enabled with minimal settings",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
			},
			expectTransport: true,
		},
		{
			name:      "custom endpoint enabled with CA file + SNI",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile:     testCAFile,
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			expectTransport: true,
		},
		{
			name:      "custom endpoint enabled with invalid CA file",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile:     "/non/existent/path/to/custom-ca-file.pem",
			},
			expectTransport: false,
		},
		{
			name:      "custom endpoint enabled with CA file contains invalid CA data",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile: func() string {
					t.Helper()

					tempDir := t.TempDir()
					caFile := filepath.Join(tempDir, "invalid-ca-file.pem")
					require.NoError(t, os.WriteFile(caFile, []byte("invalid-ca-cert"), 0600))
					return caFile
				}(),
			},
			expectTransport: false,
		},
		{
			name:      "custom endpoint enabled with CA data + SNI",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAData:     testCAData,
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			expectTransport: true,
		},
		{
			name:      "custom endpoint enabled with invalid CA data",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAData:     string("invalid-ca-cert"),
			},
			expectTransport: false,
		},
		{
			name:      "custom endpoint enabled with SNI",
			expectErr: false,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			expectTransport: true,
		},
		{
			name:      "custom endpoint disabled with extra environment variables",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesSNIName: "custom-sni.example.com",
			},
			checkErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointEnvSetWithoutTokenProxy)
			},
		},
		{
			name:      "custom endpoint enabled with both CAData and CAFile",
			expectErr: true,
			envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAData:     testCAData,
				AzureKubernetesCAFile:     testCAFile,
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
				AzureKubernetesTokenProxy: "http://custom-endpoint.com",
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
			err := Configure(&c.clientOptions)

			if c.expectErr {
				require.Error(t, err)
				if c.checkErr != nil {
					c.checkErr(t, err)
				}
				return
			}

			require.NoError(t, err)
			if c.expectTransport {
				require.NotNil(t, c.clientOptions.Transport)
				require.IsType(t, &customTokenProxyTransport{}, c.clientOptions.Transport)
			} else {
				require.Nil(t, c.clientOptions.Transport)
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

func TestCustomTokenProxyTransport_getTokenTransporter(t *testing.T) {
	cases := []struct {
		name string
		tr   *customTokenProxyTransport

		expectErr         bool
		validateTransport func(t testing.TB, httpTr *http.Transport)
	}{
		{
			name:      "no overrides",
			tr:        &customTokenProxyTransport{},
			expectErr: false,
		},
		{
			name: "with custom CA",
			tr: &customTokenProxyTransport{
				caFile: createTestCAFile(t),
			},
			expectErr: false,
			validateTransport: func(t testing.TB, httpTr *http.Transport) {
				require.NotNil(t, httpTr.TLSClientConfig)
				require.NotNil(t, httpTr.TLSClientConfig.RootCAs)
			},
		},
		{
			name: "invalid CA",
			tr: &customTokenProxyTransport{
				caData: []byte("invalid-ca-data"),
			},
			expectErr: true,
		},
		{
			name: "with SNI",
			tr: &customTokenProxyTransport{
				sniName: "example.com",
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
			tr: &customTokenProxyTransport{
				sniName: "example.com",
				caFile:  createTestCAFile(t),
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
			require.NotNil(t, c.tr.transport)
			require.Equal(t, c.tr.transport, transport, "should set the same transport to policy")
			if c.validateTransport != nil {
				c.validateTransport(t, transport)
			}
		})
	}
}

func TestCustomTokenProxyTransport_getTokenTransporter_reentry(t *testing.T) {
	t.Run("no CA overrides", func(t *testing.T) {
		tr := &customTokenProxyTransport{}
		transport, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport)

		transport2, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport2)
		require.Equal(t, transport, transport2, "should return the same transport on re-entry")
	})

	t.Run("with CAData overrides", func(t *testing.T) {
		tr := customTokenProxyTransport{
			caData: createTestCA(t),
		}
		transport, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport)

		transport2, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport2)
		require.Equal(t, transport, transport2, "should return the same transport on re-entry")
	})

	t.Run("with CAFile overrides", func(t *testing.T) {
		caFile := createTestCAFile(t)
		tr := customTokenProxyTransport{
			caFile: caFile,
		}
		transport, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport)

		transport2, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport2)
		require.Equal(t, transport, transport2, "should return the same transport on re-entry if ca file doesn't change")

		require.NoError(t, os.Truncate(caFile, 0))
		transport3, err := tr.getTokenTransporter()
		require.NoError(t, err, "empty CA file with existing transporter should not return error")
		require.NotNil(t, transport3)
		require.NotEmpty(t, tr.caData, "previous loaded CA data should be retained")
		require.NotNil(t, tr.transport, "previous transport should be retained")
		require.Equal(t, transport, transport3, "should return the same transport on re-entry if ca file is empty")

		newCAData := createTestCA(t)
		require.NoError(t, os.WriteFile(caFile, newCAData, 0600))
		transport4, err := tr.getTokenTransporter()
		require.NoError(t, err)
		require.NotNil(t, transport4)
		require.NotEqual(t, transport, transport4, "should return new transport on re-entry if ca file content is updated")
	})

	t.Run("with CAFile overrides and empty CA file on first call", func(t *testing.T) {
		caFile := filepath.Join(t.TempDir(), "empty-ca-file.pem")
		require.NoError(t, os.WriteFile(caFile, []byte{}, 0600))

		tr := customTokenProxyTransport{
			caFile: caFile,
		}
		transport, err := tr.getTokenTransporter()
		require.Error(t, err, "empty CA file on first call should return error")
		require.Nil(t, transport)
	})
}

// this provides a minimal behavior test on the transport.
// The full coverage can be found in workload identity credential tests.
func TestCustomTokenProxyTransport_Do(t *testing.T) {
	mux := http.NewServeMux()
	testServer := httptest.NewTLSServer(mux)

	ca := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: testServer.Certificate().Raw})
	require.NotEmpty(t, ca)

	const testSNIName = "test-sni-name.example.com"

	tokenProxyURL, err := url.Parse(testServer.URL + "/extra/root/path")
	require.NoError(t, err)

	transport := customTokenProxyTransport{
		caData:     ca,
		sniName:    testSNIName,
		tokenProxy: tokenProxyURL,
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://original-request.com/client-path?query=1",
		nil,
	)
	require.NoError(t, err)

	mux.HandleFunc("/extra/root/path/client-path", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, testSNIName, r.TLS.ServerName)
		require.Equal(t, "1", r.URL.Query().Get("query"))

		w.WriteHeader(http.StatusOK)
	})

	resp, err := transport.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRewriteProxyRequestURL(t *testing.T) {
	tests := []struct {
		name            string
		proxyURL        *url.URL
		reqURL          *url.URL
		wantScheme      string
		wantHost        string
		wantPath        string
		wantEscapedPath string
		wantRawQuery    string
	}{
		{
			name: "proxy url with / path; request path has no leading slash",
			proxyURL: &url.URL{
				Scheme: "https",
				Host:   "proxy.example.com",
				Path:   "/",
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "login", // no leading slash
				RawPath:  "",
				RawQuery: "a=1&b=2",
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/login",
			wantEscapedPath: "/login",
			wantRawQuery:    "a=1&b=2",
		},
		{
			name: "proxy url with / path; request path has no path",
			proxyURL: &url.URL{
				Scheme: "https",
				Host:   "proxy.example.com",
				Path:   "/",
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "",
				RawPath:  "",
				RawQuery: "a=1&b=2",
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/",
			wantEscapedPath: "/",
			wantRawQuery:    "a=1&b=2",
		},
		{
			name: "no RawPath on either; add slash between",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/base", // no trailing slash
				RawPath: "",      // explicitly empty
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "login", // no leading slash
				RawPath:  "",
				RawQuery: "a=1&b=2",
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/base/login",
			wantEscapedPath: "/base/login",
			wantRawQuery:    "a=1&b=2",
		},
		{
			name: "no RawPath; collapse double slash",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/v1/", // trailing slash
				RawPath: "",
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "/oauth2/token", // leading slash
				RawPath:  "",
				RawQuery: "x=1",
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/v1/oauth2/token",
			wantEscapedPath: "/v1/oauth2/token",
			wantRawQuery:    "x=1",
		},
		{
			name: "with RawPath; maintain escaped segments and collapse slash",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/base/",
				RawPath: "/base/",
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "/a b",   // space in segment
				RawPath:  "/a%20b", // encoded form
				RawQuery: "q=1",
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/base/a b",
			wantEscapedPath: "/base/a%20b",
			wantRawQuery:    "q=1",
		},
		{
			name: "with RawPath both sides no slashes; insert slash",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/api", // no trailing slash
				RawPath: "/api", // no trailing slash
			},
			reqURL: &url.URL{
				Scheme:  "https",
				Host:    "orig.example.com",
				Path:    "v1", // no leading slash
				RawPath: "v1", // no leading slash
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/api/v1",
			wantEscapedPath: "/api/v1",
			wantRawQuery:    "",
		},
		{
			name: "with RawPath on proxy only; preserve encoded path",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/p a",
				RawPath: "/p%20a",
			},
			reqURL: &url.URL{
				Scheme:  "https",
				Host:    "orig.example.com",
				Path:    "/b",
				RawPath: "",
			},
			wantScheme:      "https",
			wantHost:        "proxy.example.com",
			wantPath:        "/p a/b",
			wantEscapedPath: "/p%20a/b",
			wantRawQuery:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{URL: tc.reqURL}
			rewriteProxyRequestURL(req, tc.proxyURL)

			require.Equal(t, tc.wantScheme, req.URL.Scheme, "scheme mismatch")
			require.Equal(t, tc.wantHost, req.URL.Host, "host mismatch")
			require.Equal(t, tc.wantPath, req.URL.Path, "path mismatch")
			require.Equal(t, tc.wantEscapedPath, req.URL.EscapedPath(), "escaped path mismatch")
			require.Equal(t, tc.wantRawQuery, req.URL.RawQuery, "query mismatch")
		})
	}
}
