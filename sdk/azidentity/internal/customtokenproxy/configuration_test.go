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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal/exported"
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
		Options         exported.CustomTokenProxyOptions
		ClientOptions   policy.ClientOptions
		ExpectErr       bool
		AssertErr       func(t testing.TB, err error)
		ExpectTransport bool
	}{
		{
			Name:            "no custom endpoint",
			Envs:            map[string]string{},
			Options:         exported.CustomTokenProxyOptions{},
			ExpectErr:       false,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with minimal settings",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint enabled with CA file + SNI",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAFile:     testCAFile,
				SNIName:    "custom-sni.example.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint enabled with invalid CA file",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAFile:     "/non/existent/path/to/custom-ca-file.pem",
			},
			ExpectErr:       true,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with CA file contains invalid CA data",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAFile: func() string {
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
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAData:     testCAData,
				SNIName:    "custom-sni.example.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint enabled with invalid CA data",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAData:     string("invalid-ca-cert"),
			},
			ExpectErr:       true,
			ExpectTransport: false,
		},
		{
			Name: "custom endpoint enabled with SNI",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				SNIName:    "custom-sni.example.com",
			},
			ExpectErr:       false,
			ExpectTransport: true,
		},
		{
			Name: "custom endpoint disabled with extra environment variables",
			Options: exported.CustomTokenProxyOptions{
				SNIName: "custom-sni.example.com",
			},
			ExpectErr: true,
			AssertErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointSetWithoutTokenProxy)
			},
		},
		{
			Name: "custom endpoint enabled with both CAData and CAFile",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAData:     testCAData,
				CAFile:     testCAFile,
			},
			ExpectErr: true,
			AssertErr: func(t testing.TB, err error) {
				require.ErrorIs(t, err, errCustomEndpointMultipleCASourcesSet)
			},
		},
		{
			Name: "custom endpoint enabled with invalid endpoint",
			Options: exported.CustomTokenProxyOptions{
				// http endpoint is not allowed
				TokenProxy: "http://custom-endpoint.com",
			},
			ExpectErr: true,
		},
		{
			Name: "set by environment variables",
			Envs: map[string]string{
				EnvAzureKubernetesTokenProxy: "https://custom-endpoint.com",
				EnvAzureKubernetesCAFile:     testCAFile,
				EnvAzureKubernetesSNIName:    "custom-sni.example.com",
			},
			Options:         exported.CustomTokenProxyOptions{},
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

			mutateClientOptions, err := GetClientOptionsConfigurer(&tt.Options)
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

func TestBackfillOptionsFromEnv(t *testing.T) {
	tests := []struct {
		Name     string
		Options  exported.CustomTokenProxyOptions
		Envs     map[string]string
		Expected exported.CustomTokenProxyOptions
	}{
		{
			Name:     "empty",
			Options:  exported.CustomTokenProxyOptions{},
			Envs:     map[string]string{},
			Expected: exported.CustomTokenProxyOptions{},
		},
		{
			Name: "options field is not nil",
			Options: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAData:     "testCAData",
				CAFile:     "testCAFile",
				SNIName:    "custom-sni.example.com",
			},
			Envs: map[string]string{
				EnvAzureKubernetesTokenProxy: "https://endpoint-from-env.com",
				EnvAzureKubernetesCAData:     "ca-data-from-env",
				EnvAzureKubernetesCAFile:     "ca-file-from-env",
				EnvAzureKubernetesSNIName:    "sni-name-from-env",
			},
			Expected: exported.CustomTokenProxyOptions{
				TokenProxy: "https://custom-endpoint.com",
				CAData:     "testCAData",
				CAFile:     "testCAFile",
				SNIName:    "custom-sni.example.com",
			},
		},
		{
			Name:    "options field is nil",
			Options: exported.CustomTokenProxyOptions{},
			Envs: map[string]string{
				EnvAzureKubernetesTokenProxy: "https://endpoint-from-env.com",
				EnvAzureKubernetesCAData:     "ca-data-from-env",
				EnvAzureKubernetesCAFile:     "ca-file-from-env",
				EnvAzureKubernetesSNIName:    "sni-name-from-env",
			},
			Expected: exported.CustomTokenProxyOptions{
				TokenProxy: "https://endpoint-from-env.com",
				CAData:     "ca-data-from-env",
				CAFile:     "ca-file-from-env",
				SNIName:    "sni-name-from-env",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			for k, v := range tt.Envs {
				t.Setenv(k, v)
			}

			backfillOptionsFromEnv(&tt.Options)
			require.Equal(t, tt.Expected, tt.Options)
		})
	}
}
