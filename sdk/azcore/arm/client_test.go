// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
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

func TestClientTracingSanitizer(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	endpoint := srv.URL() + "?sig=secret&visible=value"
	expectedErr := &url.Error{Op: "Get", URL: endpoint, Err: errors.New("failed")}
	srv.SetError(expectedErr)
	var description string
	client, err := NewClient("module", "v1.0.0", fakeCredential{}, &ClientOptions{
		DisableRPRegistration: true,
		ClientOptions: azcore.ClientOptions{
			InsecureAllowCredentialWithHTTP: true,
			Transport:                       srv,
			Retry:                           policy.RetryOptions{MaxRetries: -1},
			Logging:                         policy.LogOptions{AllowedQueryParams: []string{"visible"}},
			TracingProvider: tracing.NewProvider(func(name, version string) tracing.Tracer {
				return tracing.NewTracer(func(ctx context.Context, name string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
					return ctx, tracing.NewSpan(tracing.SpanImpl{
						SetStatus: func(_ tracing.SpanStatus, desc string) {
							if name == "Client.Operation" {
								description = desc
							}
						},
					})
				}, nil)
			}, nil),
		},
	})
	require.NoError(t, err)
	ctx, end := runtime.StartSpan(context.Background(), "Client.Operation", client.Tracer(), nil)
	req, err := runtime.NewRequest(ctx, http.MethodGet, endpoint)
	require.NoError(t, err)
	_, err = client.Pipeline().Do(req)
	end(err)
	require.Same(t, expectedErr, err)
	require.Contains(t, description, "?sig=REDACTED&visible=value")
	require.NotContains(t, description, "secret")
}
