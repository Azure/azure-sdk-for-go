// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

// fakePolicy is a no-op policy used to satisfy newClient/newInternalPipeline
// signatures in tests that never dispatch a request through the pipeline.
type fakePolicy struct{}

func (fakePolicy) Do(req *policy.Request) (*http.Response, error) {
	return req.Next()
}

func TestDefaultCosmosHTTPClient_IsConfigured(t *testing.T) {
	require.NotNil(t, defaultCosmosHTTPClient, "package default HTTP client must be initialized")
	require.IsType(t, &http.Transport{}, defaultCosmosHTTPClient.Transport, "default HTTP client transport must be *http.Transport")
}

func TestNewDefaultCosmosHTTPClient_HTTP2PingValues(t *testing.T) {
	client, h2 := newDefaultCosmosHTTPClient()
	require.NotNil(t, client)
	require.NotNil(t, h2, "expected http2.ConfigureTransports to succeed on a freshly built transport")

	require.Equal(t, 1*time.Second, h2.ReadIdleTimeout, "Cosmos default ReadIdleTimeout must be 1s")
	require.Equal(t, 2*time.Second, h2.PingTimeout, "Cosmos default PingTimeout must be 2s")
}

func TestNewDefaultCosmosHTTPClient_TransportTuning(t *testing.T) {
	client, _ := newDefaultCosmosHTTPClient()
	tr, ok := client.Transport.(*http.Transport)
	require.True(t, ok, "default client transport must be *http.Transport")

	require.True(t, tr.ForceAttemptHTTP2, "ForceAttemptHTTP2 must be enabled so HTTP/2 ping settings take effect")
	require.NotNil(t, tr.TLSClientConfig, "TLSClientConfig must be set")
	require.Equal(t, uint16(tls.VersionTLS12), tr.TLSClientConfig.MinVersion, "MinVersion must be TLS 1.2")
}

func TestOptionsWithDefaultTransport_NilOptions_UsesDefault(t *testing.T) {
	out := optionsWithDefaultTransport(nil)
	require.Same(t, defaultCosmosHTTPClient, out.Transport, "nil options must yield the package default Transport")
}

func TestOptionsWithDefaultTransport_NilTransport_UsesDefault(t *testing.T) {
	in := &ClientOptions{}
	out := optionsWithDefaultTransport(in)
	require.Same(t, defaultCosmosHTTPClient, out.Transport, "nil Transport must be replaced by the package default")
	require.Nil(t, in.Transport, "caller-supplied options must not be mutated")
}

func TestOptionsWithDefaultTransport_PreservesCustomerTransport(t *testing.T) {
	custom := &http.Client{Transport: http.DefaultTransport}
	in := &ClientOptions{}
	in.Transport = custom

	out := optionsWithDefaultTransport(in)

	require.Same(t, custom, out.Transport, "customer-supplied Transport must be preserved")
	require.Same(t, custom, in.Transport, "caller-supplied options must not be mutated")
}

func TestNewClient_NilOptions_BuildsClient(t *testing.T) {
	gem := &globalEndpointManager{preferredLocations: []string{}}

	internal, err := newClient(fakePolicy{}, gem, nil)
	require.NoError(t, err)
	require.NotNil(t, internal)
}

func TestNewClient_NilTransportInOptions_DoesNotMutateCallerOptions(t *testing.T) {
	gem := &globalEndpointManager{preferredLocations: []string{}}
	opts := &ClientOptions{}

	internal, err := newClient(fakePolicy{}, gem, opts)
	require.NoError(t, err)
	require.NotNil(t, internal)
	require.Nil(t, opts.Transport, "newClient must not mutate the caller's ClientOptions when injecting the default Transport")
}

func TestNewClient_PreservesCustomerTransport(t *testing.T) {
	gem := &globalEndpointManager{preferredLocations: []string{}}
	custom := &http.Client{Transport: http.DefaultTransport}
	opts := &ClientOptions{}
	opts.Transport = custom

	internal, err := newClient(fakePolicy{}, gem, opts)
	require.NoError(t, err)
	require.NotNil(t, internal)
	require.Same(t, custom, opts.Transport, "customer-supplied Transport must be preserved through newClient")
}

func TestNewInternalPipeline_NilTransportInOptions_DoesNotMutateCallerOptions(t *testing.T) {
	opts := &ClientOptions{}
	_ = newInternalPipeline(fakePolicy{}, opts)
	require.Nil(t, opts.Transport, "newInternalPipeline must not mutate the caller's ClientOptions when injecting the default Transport")
}

func TestNewInternalPipeline_PreservesCustomerTransport(t *testing.T) {
	custom := &http.Client{Transport: http.DefaultTransport}
	opts := &ClientOptions{}
	opts.Transport = custom

	_ = newInternalPipeline(fakePolicy{}, opts)

	require.Same(t, custom, opts.Transport, "customer-supplied Transport must be preserved through newInternalPipeline")
}
