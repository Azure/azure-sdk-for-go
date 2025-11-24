// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenproxy

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestParseTokenProxyURL(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			u, err := parseTokenProxyURL(c.endpoint)
			c.check(t, u, err)
		})
	}
}

func TestOptions_Configure(t *testing.T) {
	var (
		testCAData = string(createTestCA(t))
		testCAFile = createTestCAFile(t)
	)

	tests := []struct {
		Name            string
		Envs            map[string]string
		Options         Options
		ClientOptions   policy.ClientOptions
		ExpectErr       bool
		AssertErr       func(t testing.TB, err error)
		ExpectTransport bool
	}{
		{
			Name:            "no custom endpoint",
			Envs:            map[string]string{},
			Options:         Options{},
			ExpectErr:       false,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with minimal settings",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint enabled with CA file + SNI",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile:     testCAFile,
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint enabled with invalid CA file",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile:     "/non/existent/path/to/custom-ca-file.pem",
			},
			ExpectErr:       true,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with CA file contains invalid CA data",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile: func() string {
					t.Helper()

					tempDir := t.TempDir()
					caFile := filepath.Join(tempDir, "invalid-ca-file.pem")
					require.NoError(t, os.WriteFile(caFile, []byte("invalid-ca-cert"), 0600))
					return caFile
				}(),
			},
			ExpectErr:       true,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with CA data + SNI",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAData:     testCAData,
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint enabled with invalid CA data",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAData:     string("invalid-ca-cert"),
			},
			ExpectErr:       true,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with SNI",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint disabled with extra environment variables",
			Options: Options{
				AzureKubernetesSNIName: "custom-sni.example.com",
			},
			ExpectErr: true,
			AssertErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointSetWithoutTokenProxy)
			},
		},
		{
			Name: "custom endpoint enabled with both CAData and CAFile",
			Options: Options{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAData:     testCAData,
				AzureKubernetesCAFile:     testCAFile,
			},
			ExpectErr: true,
			AssertErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointMultipleCASourcesSet)
			},
		},
		{
			Name: "custom endpoint enabled with invalid endpoint",
			Options: Options{
				// http endpoint is not allowed
				AzureKubernetesTokenProxy: "http://custom-endpoint.com",
			},
			ExpectErr: true,
		},
		{
			Name: "set by environment variables",
			Envs: map[string]string{
				AzureKubernetesTokenProxy: "https://custom-endpoint.com",
				AzureKubernetesCAFile:     testCAFile,
				AzureKubernetesSNIName:    "custom-sni.example.com",
			},
			Options:         Options{},
			ExpectErr:       false,
			ExpectTransport: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if len(tt.Envs) > 0 {
				for k, v := range tt.Envs {
					t.Setenv(k, v)
				}
			}

			mutateClientOptions, err := Apply(&tt.Options)
			if tt.ExpectErr {
				require.Error(t, err)
				if tt.AssertErr != nil {
					tt.AssertErr(t, err)
				}
				return
			}

			require.NoError(t, err)

			mutateClientOptions(&tt.ClientOptions)
			if tt.ExpectTransport {
				require.NotNil(t, tt.ClientOptions.Transport)
				require.IsType(t, &transport{}, tt.ClientOptions.Transport)
			} else {
				require.Nil(t, tt.ClientOptions.Transport)
			}
		})
	}
}

func TestConfiguration_defaults(t *testing.T) {
	t.Run("fills from env when empty", func(t *testing.T) {
		expected := map[string]string{
			AzureKubernetesTokenProxy: "https://custom-endpoint.com",
			AzureKubernetesSNIName:    "sni.example.com",
			AzureKubernetesCAFile:     "/path/to/ca.pem",
			AzureKubernetesCAData:     "pem-data",
		}
		for k, v := range expected {
			t.Setenv(k, v)
		}

		opts := Options{}
		opts.defaults()

		require.Equal(t, expected[AzureKubernetesTokenProxy], opts.AzureKubernetesTokenProxy)
		require.Equal(t, expected[AzureKubernetesSNIName], opts.AzureKubernetesSNIName)
		require.Equal(t, expected[AzureKubernetesCAFile], opts.AzureKubernetesCAFile)
		require.Equal(t, expected[AzureKubernetesCAData], opts.AzureKubernetesCAData)
	})

	t.Run("preserves explicit values", func(t *testing.T) {
		t.Setenv(AzureKubernetesTokenProxy, "https://env-value.com")
		opts := Options{
			AzureKubernetesTokenProxy: "https://explicit.com",
			AzureKubernetesSNIName:    "explicit-sni",
			AzureKubernetesCAFile:     "/explicit/ca.pem",
			AzureKubernetesCAData:     "explicit-ca-data",
		}

		opts.defaults()

		require.Equal(t, "https://explicit.com", opts.AzureKubernetesTokenProxy)
		require.Equal(t, "explicit-sni", opts.AzureKubernetesSNIName)
		require.Equal(t, "/explicit/ca.pem", opts.AzureKubernetesCAFile)
		require.Equal(t, "explicit-ca-data", opts.AzureKubernetesCAData)
	})
}
