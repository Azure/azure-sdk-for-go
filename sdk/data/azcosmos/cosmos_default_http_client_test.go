// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCosmosHTTPClient_Timeouts(t *testing.T) {
	require.Equal(t, 5*time.Second, defaultConnectTimeout)
	require.Equal(t, 65*time.Second, defaultHTTPRoundTripTimeout)

	c := newDefaultCosmosHTTPClient()
	require.NotNil(t, c)
	require.Equal(t, defaultHTTPRoundTripTimeout, c.Timeout)

	pool, ok := c.Transport.(*cosmosConnectionPool)
	require.True(t, ok, "expected *cosmosConnectionPool, got %T", c.Transport)
	transport := pool.transport
	require.NotNil(t, transport.DialContext)
}

func TestDefaultCosmosHTTPClient_TransportSettings(t *testing.T) {
	c := newDefaultCosmosHTTPClient()
	pool, ok := c.Transport.(*cosmosConnectionPool)
	require.True(t, ok)
	transport := pool.transport
	assert.Equal(t, 32, pool.limit)

	assert.NotNil(t, transport.Proxy)
	assert.True(t, transport.ForceAttemptHTTP2)
	assert.Equal(t, 1000, transport.MaxIdleConns)
	assert.Equal(t, 100, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 1000, transport.MaxConnsPerHost)
	assert.Equal(t, 90*time.Second, transport.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, transport.TLSHandshakeTimeout)
	assert.Equal(t, time.Second, transport.ExpectContinueTimeout)
	require.NotNil(t, transport.TLSClientConfig)
	assert.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
	assert.Equal(t, tls.RenegotiateFreelyAsClient, transport.TLSClientConfig.Renegotiation)
	require.NotNil(t, transport.HTTP2)
	assert.Equal(t, 2*time.Second, transport.HTTP2.SendPingTimeout)
	assert.Equal(t, time.Second, transport.HTTP2.PingTimeout)
}

func TestDefaultCosmosHTTPClient_Protocols(t *testing.T) {
	for _, tt := range []struct {
		name       string
		tls        bool
		http2      bool
		protoMajor int
	}{
		{name: "HTTP1", protoMajor: 1},
		{name: "HTTPS_HTTP1", tls: true, protoMajor: 1},
		{name: "HTTPS_HTTP2", tls: true, http2: true, protoMajor: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			server.EnableHTTP2 = tt.http2
			if tt.tls {
				server.StartTLS()
			} else {
				server.Start()
			}
			t.Cleanup(server.Close)

			c := newDefaultCosmosHTTPClient()
			t.Cleanup(c.CloseIdleConnections)
			pool, ok := c.Transport.(*cosmosConnectionPool)
			require.True(t, ok)
			transport := pool.transport
			if tt.tls {
				serverTransport, ok := server.Client().Transport.(*http.Transport)
				require.True(t, ok)
				transport.TLSClientConfig.RootCAs = serverTransport.TLSClientConfig.RootCAs
			}

			resp, err := c.Get(server.URL)
			require.NoError(t, err)
			defer func() { assert.NoError(t, resp.Body.Close()) }()
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
			assert.Equal(t, tt.protoMajor, resp.ProtoMajor)
			require.NotNil(t, resp.Request)
			assert.Equal(t, server.URL, resp.Request.URL.String())
		})
	}
}

func TestWithDefaultTransport_NilOptions(t *testing.T) {
	got := withDefaultTransport(nil)
	require.NotNil(t, got)
	require.Same(t, defaultCosmosHTTPClient, got.Transport)
}

func TestWithDefaultTransport_NoTransportSet(t *testing.T) {
	in := &ClientOptions{
		EnableContentResponseOnWrite: true,
	}
	got := withDefaultTransport(in)
	require.NotSame(t, in, got, "expected withDefaultTransport to clone options")
	require.Nil(t, in.Transport, "caller-supplied options must not be mutated")
	require.Same(t, defaultCosmosHTTPClient, got.Transport)
	require.True(t, got.EnableContentResponseOnWrite)
}

func TestWithDefaultTransport_PreservesCallerTransport(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	in := &ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: srv},
	}
	got := withDefaultTransport(in)
	require.Same(t, in, got, "expected withDefaultTransport to return same options when Transport is provided")
	require.Same(t, srv, got.Transport)
}
