//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

type fakeCredential struct{}

func (mc fakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestNewClient(t *testing.T) {
	client, err := NewClient("module", "v1.0.0", fakeCredential{}, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint, client.Endpoint())
	require.NotZero(t, client.Pipeline())
	require.Zero(t, client.Tracer())

	client, err = NewClient("module", "", fakeCredential{}, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureChina,
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cloud.AzureChina.Services[cloud.ResourceManager].Endpoint, client.Endpoint())
}

func TestNewClientError(t *testing.T) {
	client, err := NewClient("module", "malformed", fakeCredential{}, nil)
	require.Error(t, err)
	require.Nil(t, client)

	badCloud := cloud.Configuration{
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "fake/audience",
			},
		},
	}
	client, err = NewClient("module", "v1.0.0", fakeCredential{}, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: badCloud,
		},
	})
	require.Error(t, err)
	require.Nil(t, client)
}
