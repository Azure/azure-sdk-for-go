// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package audience

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

// mockErrorTransporter is a custom transporter that returns a specific error
type mockErrorTransporter struct {
	err error
}

func (m *mockErrorTransporter) Do(req *http.Request) (*http.Response, error) {
	return nil, m.err
}

func TestAudienceErrorHandlingPolicy_NoError(t *testing.T) {
	// Test that the policy doesn't interfere when there's no error
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	errPolicy := NewAudienceErrorHandlingPolicy(false)
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{errPolicy},
	}, &policy.ClientOptions{
		Transport: srv,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAudienceErrorHandlingPolicy_NonAudienceError(t *testing.T) {
	// Test that non-audience errors are passed through unchanged
	originalErr := errors.New("some other error")
	transport := &mockErrorTransporter{err: originalErr}

	errPolicy := NewAudienceErrorHandlingPolicy(false)
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{errPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.com")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, originalErr, err)
}

func TestAudienceErrorHandlingPolicy_AudienceErrorNotConfigured(t *testing.T) {
	// Test that when audience is not configured, we get the "no audience" error message
	audienceErr := errors.New("authentication failed: " + AadAudienceErrorCode + " invalid audience")
	transport := &mockErrorTransporter{err: audienceErr}

	errPolicy := NewAudienceErrorHandlingPolicy(false)
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{errPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.com")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "No authentication token audience was provided")
	require.Contains(t, err.Error(), "https://aka.ms/appconfig/client-token-audience")
}

func TestAudienceErrorHandlingPolicy_AudienceErrorConfigured(t *testing.T) {
	// Test that when audience is configured, we get the "wrong audience" error message
	audienceErr := errors.New("authentication failed: " + AadAudienceErrorCode + " invalid audience")
	transport := &mockErrorTransporter{err: audienceErr}

	errPolicy := NewAudienceErrorHandlingPolicy(true)
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{errPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.com")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "An incorrect token audience was provided")
	require.Contains(t, err.Error(), "https://aka.ms/appconfig/client-token-audience")
}

func TestAudienceErrorHandlingPolicy_WrappedError(t *testing.T) {
	// Test that wrapped errors containing the error code are also handled
	innerErr := errors.New(AadAudienceErrorCode + ": invalid audience")
	wrappedErr := errors.New("wrapped: " + innerErr.Error())
	transport := &mockErrorTransporter{err: wrappedErr}

	errPolicy := NewAudienceErrorHandlingPolicy(false)
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{errPolicy},
	}, &policy.ClientOptions{
		Transport: transport,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.com")
	require.NoError(t, err)

	resp, err := pl.Do(req)
	require.Error(t, err)
	require.Nil(t, resp)
	require.True(t, strings.Contains(err.Error(), "No authentication token audience was provided") ||
		strings.Contains(err.Error(), AadAudienceErrorCode))
}
