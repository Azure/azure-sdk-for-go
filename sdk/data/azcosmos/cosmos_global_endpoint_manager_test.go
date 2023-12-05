// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestGlobalEndpointManagerGetWriteEndpoints(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)

	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	writeEndpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err)

	serverEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	expectedWriteEndpoints := []url.URL{
		*serverEndpoint,
	}
	assert.Equal(t, expectedWriteEndpoints, writeEndpoints)
}

func TestGlobalEndpointManagerGetReadEndpoints(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	serverEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	expectedReadEndpoints := []url.URL{
		*serverEndpoint,
	}
	assert.Equal(t, expectedReadEndpoints, readEndpoints)
}

func TestGlobalEndpointManagerMarkEndpointUnavailableForRead(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)

	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	endpoint, err := url.Parse(client.endpoint)
	assert.NoError(t, err)

	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)
	print(readEndpoints)

	err = gem.MarkEndpointUnavailableForRead(*endpoint)
	assert.NoError(t, err)

	unavailable := gem.IsEndpointUnavailable(*endpoint, 1)
	assert.True(t, unavailable)
}

func TestGlobalEndpointManagerMarkEndpointUnavailableForWrite(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)

	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	endpoint, err := url.Parse(client.endpoint)
	assert.NoError(t, err)

	err = gem.MarkEndpointUnavailableForWrite(*endpoint)
	assert.NoError(t, err)

	unavailable := gem.IsEndpointUnavailable(*endpoint, 2)
	assert.True(t, unavailable)
}

func TestGlobalEndpointManagerGetEndpointLocation(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	westRegion := accountRegion{
		Name:     "West US",
		Endpoint: srv.URL(),
	}

	properties := accountProperties{
		ReadRegions:                  []accountRegion{westRegion},
		WriteRegions:                 []accountRegion{westRegion},
		EnableMultipleWriteLocations: false,
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}
	srv.SetResponse(mock.WithBody(jsonString))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)

	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	gem, err := newGlobalEndpointManager(client, []string{}, 5*time.Minute)
	assert.NoError(t, err)

	serverEndpoint, err := url.Parse(client.endpoint)
	assert.NoError(t, err)

	location := gem.GetEndpointLocation(*serverEndpoint)

	expectedLocation := "West US"
	assert.Equal(t, expectedLocation, location)
}

func TestGlobalEndpointManagerGetAccountProperties(t *testing.T) {
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

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	accountProps, err := gem.GetAccountProperties()
	assert.NoError(t, err)

	expectedAccountProps := accountProperties{
		ReadRegions:                  nil,
		WriteRegions:                 nil,
		EnableMultipleWriteLocations: false,
	}
	assert.Equal(t, expectedAccountProps, accountProps)
}

func TestGlobalEndpointManagerCanUseMultipleWriteLocations(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)

	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	preferredRegions := []string{"West US", "Central US"}

	serverEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	mockLc := newLocationCache(preferredRegions, *serverEndpoint)
	mockLc.enableMultipleWriteLocations = true
	mockLc.useMultipleWriteLocations = true

	mockGem := globalEndpointManager{
		client:              client,
		preferredLocations:  preferredRegions,
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	// Multiple locations should be false for default GEM
	canUseMultipleWriteLocs := gem.CanUseMultipleWriteLocations()
	assert.False(t, canUseMultipleWriteLocs)

	// Mock GEM with multiple write locations available should show true
	canUseMultipleWriteLocs = mockGem.CanUseMultipleWriteLocations()
	assert.True(t, canUseMultipleWriteLocs)
}

func TestGlobalEndpointManagerUpdate(t *testing.T) { //This test should be testing for the lock on the update/ refresh mechanism
	// Create a mock client
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	westRegion := accountRegion{
		Name:     "West US",
		Endpoint: srv.URL(),
	}

	properties := accountProperties{
		ReadRegions:                  []accountRegion{westRegion},
		WriteRegions:                 []accountRegion{westRegion},
		EnableMultipleWriteLocations: false,
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}
	srv.SetResponse(mock.WithBody(jsonString))

	verifier := headerPoliciesVerify{}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, &verifier}}, &policy.ClientOptions{Transport: srv})
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	assert.NoError(t, err)

	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
	})

	client := &Client{endpoint: srv.URL(), pipeline: pl}

	gem, err := newGlobalEndpointManager(client, []string{}, 5*time.Minute)
	assert.NoError(t, err)

	// Update the location cache and client's default endpoint
	err = gem.Update()
	assert.NoError(t, err)

	// Get the write endpoints after the update
	writeEndpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err)

	serverEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	// Assert the expected write endpoints after the update
	expectedWriteEndpoints := []url.URL{
		*serverEndpoint,
	}
	assert.Equal(t, expectedWriteEndpoints, writeEndpoints)

	// Get the read endpoints after the update
	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	// Assert the expected read endpoints
	expectedReadEndpoints := []url.URL{
		*serverEndpoint,
	}
	assert.Equal(t, expectedReadEndpoints, readEndpoints)
}

// Emulator Test

func TestGlobalEndpointManagerEmulator(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)
	emulatorRegionName := "South Central US"
	preferredRegions := []string{}
	emulatorRegion := accountRegion{Name: emulatorRegionName, Endpoint: "https://127.0.0.1:8081/"}

	gem, err := newGlobalEndpointManager(client, preferredRegions, 5*time.Minute)
	assert.NoError(t, err)

	accountProps, err := gem.GetAccountProperties()
	assert.NoError(t, err)

	// Verify the expected account properties
	expectedAccountProps := accountProperties{
		ReadRegions:                  []accountRegion{emulatorRegion},
		WriteRegions:                 []accountRegion{emulatorRegion},
		EnableMultipleWriteLocations: false,
	}
	assert.Equal(t, expectedAccountProps, accountProps)

	emulatorEndpoint, err := url.Parse("https://localhost:8081/")
	assert.NoError(t, err)

	// Verify the read endpoints
	readEndpoints, err := gem.GetReadEndpoints()
	assert.NoError(t, err)

	expectedEndpoints := []url.URL{
		*emulatorEndpoint,
	}
	assert.Equal(t, expectedEndpoints, readEndpoints)

	// Verify the write endpoints
	writeEndpoints, err := gem.GetWriteEndpoints()
	assert.NoError(t, err)

	assert.Equal(t, expectedEndpoints, writeEndpoints)

	// Assert location cache is not populated until update() is called
	locationInfo := gem.locationCache.locationInfo
	availableLocation := []string{}
	availableEndpointsByLocation := map[string]url.URL{}

	assert.Equal(t, locationInfo.availReadLocations, availableLocation)
	assert.Equal(t, locationInfo.availWriteLocations, availableLocation)
	assert.Equal(t, locationInfo.availReadEndpointsByLocation, availableEndpointsByLocation)
	assert.Equal(t, locationInfo.availWriteEndpointsByLocation, availableEndpointsByLocation)

	//update and assert available locations are now populated in location cache
	gem.Update()
	locationInfo = gem.locationCache.locationInfo

	assert.Equal(t, len(locationInfo.availReadLocations), len(availableLocation)+1)
	assert.Equal(t, len(locationInfo.availWriteLocations), len(availableLocation)+1)
	assert.Equal(t, locationInfo.availWriteLocations[0], emulatorRegionName)
	assert.Equal(t, locationInfo.availReadLocations[0], emulatorRegionName)
	assert.Equal(t, len(locationInfo.availReadEndpointsByLocation), len(availableEndpointsByLocation)+1)
	assert.Equal(t, len(locationInfo.availWriteEndpointsByLocation), len(availableEndpointsByLocation)+1)
}
