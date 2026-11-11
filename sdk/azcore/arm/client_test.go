// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// tests assume Azure Public Cloud unless they specify otherwise
	if err := os.Unsetenv(shared.EnvAzureEnvironment); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

type fakeCredential struct{}

func (mc fakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestNewClient(t *testing.T) {
	t.Setenv(shared.EnvAzureEnvironment, "")
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

// scopeCredential records the scopes of each token request
type scopeCredential struct {
	scopes *[]string
}

func (sc scopeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	*sc.scopes = append(*sc.scopes, options.Scopes...)
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// requireClientCloud asserts that a client created with opts uses expected's endpoint and token audience
func requireClientCloud(t *testing.T, expected cloud.Configuration, opts *ClientOptions) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse()
	if opts == nil {
		opts = &ClientOptions{}
	}
	opts.DisableRPRegistration = true
	opts.Transport = srv
	before := opts.Cloud
	scopes := []string{}
	client, err := NewClient("module", "v1.0.0", scopeCredential{scopes: &scopes}, opts)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, before, opts.Cloud, "NewClient shouldn't modify the caller's options")
	require.Equal(t, expected.Services[cloud.ResourceManager].Endpoint, client.Endpoint())

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	_, err = client.Pipeline().Do(req)
	require.NoError(t, err)
	require.Equal(t, []string{expected.Services[cloud.ResourceManager].Audience + "/.default"}, scopes)
}

func TestNewClientEnvironment(t *testing.T) {
	for name, expected := range map[string]cloud.Configuration{
		"":                       cloud.AzurePublic,
		"AzureChina":             cloud.AzureChina,
		"AzureChinaCloud":        cloud.AzureChina,
		"AzureCloud":             cloud.AzurePublic,
		"AzureGovernment":        cloud.AzureGovernment,
		"AzurePublic":            cloud.AzurePublic,
		"AzurePublicCloud":       cloud.AzurePublic,
		"AzureUSGovernment":      cloud.AzureGovernment,
		"AzureUSGovernmentCloud": cloud.AzureGovernment,
		"azurechina":             cloud.AzureChina,
		"AZUREUSGOVERNMENT":      cloud.AzureGovernment,
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(shared.EnvAzureEnvironment, name)
			requireClientCloud(t, expected, nil)

			// a cloud specified in code takes precedence
			for _, specified := range []cloud.Configuration{cloud.AzureChina, cloud.AzureGovernment, cloud.AzurePublic} {
				requireClientCloud(t, specified, &ClientOptions{
					ClientOptions: azcore.ClientOptions{
						Cloud: specified,
					},
				})
			}
		})
	}
}

// hostRecorder records the host of each request before sending it to a mock server
type hostRecorder struct {
	srv   *mock.Server
	hosts []string
}

func (h *hostRecorder) Do(req *http.Request) (*http.Response, error) {
	h.hosts = append(h.hosts, req.URL.Host)
	return h.srv.Do(req)
}

func TestNewClientEnvironmentRPRegistration(t *testing.T) {
	const (
		resource     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account"
		unregistered = `{"error":{"code":"MissingSubscriptionRegistration"}}`
		registered   = `{"id":"/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage","namespace":"Microsoft.Storage","registrationState":"Registered"}`
	)
	for name, expected := range map[string]cloud.Configuration{
		"AzureChina":        cloud.AzureChina,
		"AzureUSGovernment": cloud.AzureGovernment,
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(shared.EnvAzureEnvironment, name)
			srv, close := mock.NewTLSServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			// the resource provider is unregistered
			srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(unregistered)))
			// responses to the policy's registration and status requests
			srv.RepeatResponse(2, mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(registered)))
			// response to the retried original request
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

			scopes := []string{}
			recorder := &hostRecorder{srv: srv}
			client, err := NewClient("module", "v1.0.0", scopeCredential{scopes: &scopes}, &ClientOptions{
				ClientOptions: azcore.ClientOptions{Transport: recorder},
			})
			require.NoError(t, err)
			req, err := runtime.NewRequest(context.Background(), http.MethodGet, client.Endpoint()+resource)
			require.NoError(t, err)
			resp, err := client.Pipeline().Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			// the original request, registration, status check, and retry should all go to the named cloud
			u, err := url.Parse(expected.Services[cloud.ResourceManager].Endpoint)
			require.NoError(t, err)
			require.Equal(t, []string{u.Host, u.Host, u.Host, u.Host}, recorder.hosts)
			require.NotEmpty(t, scopes)
			for _, scope := range scopes {
				require.Equal(t, expected.Services[cloud.ResourceManager].Audience+"/.default", scope)
			}
		})
	}
}

func TestResolveCloud(t *testing.T) {
	custom := cloud.Configuration{ActiveDirectoryAuthorityHost: "https://localhost"}
	servicesOnly := cloud.Configuration{Services: map[cloud.ServiceName]cloud.ServiceConfiguration{}}
	for _, test := range []struct {
		name, value string
		specified   cloud.Configuration
		expected    cloud.Configuration
	}{
		{name: "default", expected: cloud.AzurePublic},
		{name: "named cloud", value: "AzureChina", expected: cloud.AzureChina},
		{name: "surrounding whitespace", value: " AzureChina\n", expected: cloud.AzureChina},
		{name: "specified cloud", specified: custom, expected: custom},
		{name: "specified cloud overrides variable", value: "AzureChina", specified: custom, expected: custom},
		{name: "specified cloud with invalid variable", value: "invalid", specified: custom, expected: custom},
		{name: "specified services only", value: "AzureChina", specified: servicesOnly, expected: servicesOnly},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(shared.EnvAzureEnvironment, test.value)
			actual, err := resolveCloud(test.specified)
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
	for _, value := range []string{"AzureGermanCloud", "Azure China", "invalid"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv(shared.EnvAzureEnvironment, value)
			actual, err := resolveCloud(cloud.Configuration{})
			require.ErrorContains(t, err, shared.EnvAzureEnvironment)
			require.ErrorContains(t, err, value)
			require.Zero(t, actual)
		})
	}
}

func TestNewClientInvalidEnvironment(t *testing.T) {
	t.Setenv(shared.EnvAzureEnvironment, "invalid")
	client, err := NewClient("module", "v1.0.0", fakeCredential{}, nil)
	require.ErrorContains(t, err, shared.EnvAzureEnvironment)
	require.Nil(t, client)

	// the variable is ignored when a cloud is specified in code
	client, err = NewClient("module", "v1.0.0", fakeCredential{}, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureGovernment,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cloud.AzureGovernment.Services[cloud.ResourceManager].Endpoint, client.Endpoint())
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

	// the environment variable doesn't complete or replace a partial configuration
	t.Setenv(shared.EnvAzureEnvironment, "AzureChina")
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
