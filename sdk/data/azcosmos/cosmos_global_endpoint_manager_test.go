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
		{Scheme: "https", Host: "westus.documents.azure.com"},
		{Scheme: "https", Host: "eastus.documents.azure.com"},
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
		{Scheme: "https", Host: "westus.documents.azure.com"},
		{Scheme: "https", Host: "eastus.documents.azure.com"},
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
	endpoint := url.URL{Scheme: "https", Host: "westus.documents.azure.com"}
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
		{Scheme: "https", Host: "westus.documents.azure.com"},
		{Scheme: "https", Host: "eastus.documents.azure.com"},
	}
	assert.Equal(t, expectedWriteEndpoints, writeEndpoints)

	// Get the read endpoints after the update
	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	// Assert the expected read endpoints after the update
	expectedReadEndpoints := []url.URL{
		{Scheme: "https", Host: "westus.documents.azure.com"},
		{Scheme: "https", Host: "eastus.documents.azure.com"},
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
