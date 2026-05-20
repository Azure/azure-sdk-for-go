// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestDefaultCosmosHTTPClient_Timeouts(t *testing.T) {
	require.Equal(t, 5*time.Second, defaultConnectTimeout)
	require.Equal(t, 65*time.Second, defaultRequestTimeout)

	c := newDefaultCosmosHTTPClient()
	require.NotNil(t, c)
	require.Equal(t, defaultRequestTimeout, c.Timeout)

	transport, ok := c.Transport.(*http.Transport)
	require.True(t, ok, "expected *http.Transport, got %T", c.Transport)
	require.NotNil(t, transport.DialContext)
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
