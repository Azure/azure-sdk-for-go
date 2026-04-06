// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestWithLayoutEndpoint(t *testing.T) {
	ctx := context.Background()
	endpoint := "layout.blob.core.windows.net"

	ctxWithEndpoint := WithLayoutEndpoint(ctx, endpoint)

	value := ctxWithEndpoint.Value(ctxLayoutEndpointKey{})
	require.NotNil(t, value)
	require.Equal(t, endpoint, value.(string))
}

func TestWithLayoutEndpointEmptyString(t *testing.T) {
	ctx := context.Background()

	ctxWithEndpoint := WithLayoutEndpoint(ctx, "")

	value := ctxWithEndpoint.Value(ctxLayoutEndpointKey{})
	require.Nil(t, value)
}

func TestLayoutPolicyWithLayoutEndpoint(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithStatusCode(200))

	layoutHost := "layout.blob.core.windows.net:443"
	originalHost := "original.blob.core.windows.net"

	p := NewLayoutPolicy()
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerCall: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)

	ctx := WithLayoutEndpoint(context.Background(), "https://"+layoutHost)
	req, err := runtime.NewRequest(ctx, http.MethodGet, "https://"+originalHost+"/container/blob")
	require.NoError(t, err)

	// Verify the Host header is set to original and URL host is changed
	_, err = pl.Do(req)
	require.NoError(t, err)

	// After policy execution, the Host header should be set to original host
	require.Equal(t, originalHost, req.Raw().Host)
	// The URL host should be changed to layout endpoint
	require.Equal(t, layoutHost, req.Raw().URL.Host)
}

func TestLayoutPolicyWithLayoutEndpointEmpty(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithStatusCode(200))

	originalHost := "original.blob.core.windows.net"

	p := NewLayoutPolicy()
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerCall: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)

	// Manually set empty string in context (bypassing WithLayoutEndpoint)
	ctx := context.WithValue(context.Background(), ctxLayoutEndpointKey{}, "")
	req, err := runtime.NewRequest(ctx, http.MethodGet, "https://"+originalHost+"/container/blob")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.NoError(t, err)

	// URL host should remain unchanged when endpoint is empty
	require.Equal(t, originalHost, req.Raw().URL.Host)
}

func TestLayoutPolicyWithLayoutEndpointInvalid(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()

	p := NewLayoutPolicy()
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerCall: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)

	// Use an invalid URL that will fail parsing
	ctx := context.WithValue(context.Background(), ctxLayoutEndpointKey{}, "://invalid-url")
	req, err := runtime.NewRequest(ctx, http.MethodGet, "https://original.blob.core.windows.net/container/blob")
	require.NoError(t, err)

	_, err = pl.Do(req)
	require.Error(t, err)
}

func TestLayoutPolicyWithoutLayoutEndpoint(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithStatusCode(200))

	originalHost := "original.blob.core.windows.net"

	p := NewLayoutPolicy()
	pl := runtime.NewPipeline("", "",
		runtime.PipelineOptions{PerCall: []policy.Policy{p}},
		&policy.ClientOptions{Transport: srv},
	)

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://"+originalHost+"/container/blob")
	require.NoError(t, err)

	originalURLHost := req.Raw().URL.Host
	originalReqHost := req.Raw().Host

	_, err = pl.Do(req)
	require.NoError(t, err)

	// Without layout endpoint, Host header and URL host should remain unchanged
	require.Equal(t, originalURLHost, req.Raw().URL.Host)
	require.Equal(t, originalReqHost, req.Raw().Host)
}

func TestNewLayoutPolicy(t *testing.T) {
	p := NewLayoutPolicy()
	require.NotNil(t, p)
	_, ok := p.(layoutPolicy)
	require.True(t, ok)
}
