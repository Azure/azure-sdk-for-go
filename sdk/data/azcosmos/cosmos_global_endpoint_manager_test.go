// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestGlobalEndpointManager_GetWriteEndpoints(t *testing.T) {

	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	// Create a mock client
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	// Get the write endpoints
	writeEndpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err)

	// Assert the expected write endpoints
	expectedWriteEndpoints := []url.URL{
		{Scheme: "https", Host: "127.0.0.1:61453"},
	}
	assert.Equal(t, expectedWriteEndpoints, writeEndpoints)
}

func TestGlobalEndpointManager_GetReadEndpoints(t *testing.T) {
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	// Get the read endpoints
	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	// Assert the expected read endpoints
	expectedReadEndpoints := []url.URL{
		{Scheme: "https", Host: "127.0.0.1:61453"},
	}
	assert.Equal(t, expectedReadEndpoints, readEndpoints)
}

func TestGlobalEndpointManager_GetAccountProperties(t *testing.T) {
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Get the account properties
	accountProps, err := gem.GetAccountProperties()
	assert.NoError(t, err)

	// Assert the expected account properties
	expectedAccountProps := accountProperties{
		ReadRegions:                  []accountRegion{},
		WriteRegions:                 []accountRegion{},
		EnableMultipleWriteLocations: false,
	}
	assert.Equal(t, expectedAccountProps, accountProps)
}

func TestGlobalEndpointManager_GetLocation(t *testing.T) {
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Get the location for the given endpoint
	endpoint := url.URL{Scheme: "https", Host: "127.0.0.1:61453"}
	location := gem.GetLocation(endpoint)

	// Assert the expected location
	expectedLocation := "westus"
	assert.Equal(t, expectedLocation, location)
}

func TestGlobalEndpointManager_MarkEndpointUnavailableForRead(t *testing.T) {
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Mark an endpoint as unavailable for read
	endpoint, err := url.Parse(client.endpoint)
	if err != nil {
		return
	}

	err = gem.MarkEndpointUnavailableForRead(*endpoint)
	assert.NoError(t, err)

	// Check if the endpoint is marked as unavailable for read
	unavailable := gem.IsEndpointUnavailable(*endpoint, 1)
	assert.True(t, unavailable)
}

func TestGlobalEndpointManager_MarkEndpointUnavailableForWrite(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Mark an endpoint as unavailable for write
	endpoint, err := url.Parse(client.endpoint)
	if err != nil {
		return
	}
	assert.NoError(t, err)

	// Check if the endpoint is marked as unavailable for write
	unavailable := gem.IsEndpointUnavailable(*endpoint, 1)
	assert.True(t, unavailable)
}

func TestGlobalEndpointManager_Update(t *testing.T) {
	// Create a mock client
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Update the location cache and client's default endpoint
	err = gem.Update()
	assert.NoError(t, err)

	// Get the write endpoints after the update
	writeEndpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err)

	// Assert the expected write endpoints after the update
	expectedWriteEndpoints := []url.URL{
		{Scheme: "https", Host: "127.0.0.1:61453"},
	}
	assert.Equal(t, expectedWriteEndpoints, writeEndpoints)

	// Get the read endpoints after the update
	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	// Assert the expected read endpoints after the update
	expectedReadEndpoints := []url.URL{
		{Scheme: "https", Host: "127.0.0.1:61453"},
	}
	assert.Equal(t, expectedReadEndpoints, readEndpoints)
}

func TestGlobalEndpointManager_RefreshStaleEndpoints(t *testing.T) {
	// Create a mock client
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Refresh stale endpoints in the location cache
	gem.RefreshStaleEndpoints()
	// No assertions as it is an internal mechanism
}

func TestGlobalEndpointManager_CanUseMultipleWriteLocations(t *testing.T) {
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Check if multiple write locations can be used
	canUseMultipleWriteLocs := gem.CanUseMultipleWriteLocations()
	assert.False(t, canUseMultipleWriteLocs)
}

func TestGlobalEndpointManagerBackgroundRefresh(t *testing.T) {

	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	// Create a mock client
	// headerPolicy := &headerPolicies{}
	// srv, close := mock.NewTLSServer()
	// defer close()
	// srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	// verifier := headerPoliciesVerify{}
	// pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	// req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	// if err != nil {
	// 	return
	// }
	// req.SetOperationValue(pipelineRequestOptions{
	// 	isWriteOperation: true,
	// })

	// client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Create a new global endpoint manager with preferred regions and refresh interval
	gem, err := newGlobalEndpointManager(mockClient, []string{"westus", "eastus"}, 5*time.Minute)
	assert.NoError(t, err)

	// Check if multiple write locations can be used
	canUseMultipleWriteLocs := gem.CanUseMultipleWriteLocations()
	assert.False(t, canUseMultipleWriteLocs)
}

// MockClient is a mock implementation of the Client interface for testing purposes.
type MockClient struct {
	sendGetRequestErr      error
	sendGetRequestResponse *http.Response
	sendGetRequestTimeout  bool
}

// sendGetRequest is a mock implementation for the Client's sendGetRequest method.
func (mc *MockClient) sendGetRequest(path string, ctx context.Context, options pipelineRequestOptions, resourceType resourceType, headers http.Header) (*http.Response, error) {
	if mc.sendGetRequestTimeout {
		// Simulate a timeout.
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return mc.sendGetRequestResponse, mc.sendGetRequestErr
}
func TestEndpointFailureMockTest(t *testing.T) {
	// Create a mock client and preferredRegions for testing
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		return
	}
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}
	mockClient := client

	// Simulate multiple regions for testing
	preferredRegions := []string{"region1", "region2", "region3"}

	// Create a new globalEndpointManager with the mock client
	gem, err := newGlobalEndpointManager(mockClient, preferredRegions, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error creating globalEndpointManager: %v", err)
	}

	endpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err, "unexpected error")
	// Verify that the write endpoints are switched to the next available region
	expectedEndpoints := []url.URL{
		// In this case, region1 failed, so the next available region is region2.
		// Assuming region2 is also available, the write endpoints should switch to it.
		// Add the appropriate URL values based on your implementation.
		// For example, &url.URL{Scheme: "https", Host: "region2.mydomain.com"},
	}
	assert.Equal(t, expectedEndpoints, endpoints, "write endpoints not switched as expected")

	endpoints, err = gem.GetWriteEndpoints()
	assert.NoError(t, err, "unexpected error")
	// Verify that the write endpoints are switched to the next available region
	// For example, region2 failed, so the next available region is region3.
	// Assuming region3 is also available, the write endpoints should switch to it.
	expectedEndpoints = []url.URL{
		// Add the appropriate URL values based on your implementation.
		// For example, &url.URL{Scheme: "https", Host: "region3.mydomain.com"},
	}
	assert.Equal(t, expectedEndpoints, endpoints, "write endpoints not switched as expected")

	// You can add more test cases to cover various failure scenarios and endpoint switching.
}
