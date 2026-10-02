// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

type tokenCredFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (t tokenCredFunc) GetToken(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if l := len(tro.Scopes); l != 1 {
		return azcore.AccessToken{}, fmt.Errorf("unexpected scopes len %d", l)
	}
	return t(ctx, tro)
}

type fakeTransport struct{}

func (fakeTransport) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Request:    req,
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
		Header:     http.Header{},
	}, nil
}

func TestNewClient_AutoDetectAudience(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		scope    string
		cloud    cloud.Configuration
	}{
		{
			name:     "AzurePublicAzConfig",
			endpoint: "https://example.azconfig.io",
			scope:    "https://azconfig.io/.default",
		},
		{
			name:     "AzurePublicAppConfig",
			endpoint: "https://example.appconfig.azure.com",
			scope:    "https://appconfig.azure.com/.default",
		},
		{
			name:     "AzureChina",
			endpoint: "https://example.azconfig.azure.cn",
			scope:    "https://azconfig.azure.cn/.default",
		},
		{
			name:     "AzureGovernment",
			endpoint: "https://example.appconfig.azure.us",
			scope:    "https://appconfig.azure.us/.default",
		},
		{
			name:     "CustomSovereignCloud",
			endpoint: "https://example.appconfig.sovcloud-api.fr",
			scope:    "https://appconfig.sovcloud-api.fr/.default",
		},
		{
			name:     "CustomSovereignCloudWithRegion",
			endpoint: "https://example.eastus.APPconfig.sovereign.cloud:8443",
			scope:    "https://appconfig.sovereign.cloud/.default",
		},
		{
			name:     "Staging",
			endpoint: "https://example.appconfig-staging.azure.com",
			scope:    "https://appconfig-staging.azure.com/.default",
		},
		{
			name:     "HyphenatedMarker",
			endpoint: "https://example.appconfig-test.azure.com",
			scope:    "https://appconfig.azure.com/.default",
		},
		{
			name:     "LookalikeMarker",
			endpoint: "https://example.fooappconfig.azure.us",
			scope:    "https://appconfig.azure.com/.default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.endpoint, tokenCredFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				require.Equal(t, tt.scope, tro.Scopes[0])
				return azcore.AccessToken{}, nil
			}), &ClientOptions{ClientOptions: policy.ClientOptions{Cloud: tt.cloud, Transport: &fakeTransport{}}})
			require.NoError(t, err)

			_, _ = client.GetSetting(context.Background(), "fake-key", nil)
		})
	}
}

func TestNewClient_SovereignClouds(t *testing.T) {
	azureBleu := cloud.Configuration{
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			ServiceName: {
				Audience: "https://appconfig.sovcloud-api.fr",
			},
		},
	}

	tests := []struct {
		label    string
		endpoint string
		scope    string
		cfg      cloud.Configuration
	}{
		{
			label:    "AzureChina",
			endpoint: "https://example.azconfig.azure.cn",
			scope:    "https://appconfig.azure.cn/.default",
			cfg:      cloud.AzureChina,
		},
		{
			label:    "AzureGovernment",
			endpoint: "https://example.azconfig.azure.us",
			scope:    "https://appconfig.azure.us/.default",
			cfg:      cloud.AzureGovernment,
		},
		{
			label:    "AzurePublic",
			endpoint: "https://example.azconfig.io",
			scope:    "https://appconfig.azure.com/.default",
			cfg:      cloud.AzurePublic,
		},
		{
			label:    "AzureBleu",
			endpoint: "https://example.azconfig.sovcloud-api.fr",
			scope:    "https://appconfig.sovcloud-api.fr/.default",
			cfg:      azureBleu,
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			client, err := NewClient(tt.endpoint, tokenCredFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				require.Equal(t, tt.scope, tro.Scopes[0])
				return azcore.AccessToken{}, nil
			}), &ClientOptions{
				ClientOptions: policy.ClientOptions{
					Cloud:     tt.cfg,
					Transport: &fakeTransport{},
				},
			})
			require.NoError(t, err)

			// Call an API to trigger the pipeline which will call GetToken on our fake cred
			_, _ = client.GetSetting(context.Background(), "fake-key", nil)
		})
	}
}
