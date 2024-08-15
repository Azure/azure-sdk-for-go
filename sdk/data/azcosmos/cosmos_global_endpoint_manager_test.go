// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

type countPolicy struct {
	callCount int
}

func (p *countPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.callCount += 1
	return req.Next()
}

func TestGlobalEndpointManagerGetWriteEndpoints(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{"West US", "Central US"}, 5*time.Minute, true)
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
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{"West US", "Central US"}, 5*time.Minute, true)
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
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	endpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{"West US", "Central US"}, 5*time.Minute, true)
	assert.NoError(t, err)

	err = gem.MarkEndpointUnavailableForRead(*endpoint)
	assert.NoError(t, err)

	unavailable := gem.IsEndpointUnavailable(*endpoint, 1)
	assert.True(t, unavailable)
}

func TestGlobalEndpointManagerMarkEndpointUnavailableForWrite(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	endpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{"West US", "Central US"}, 5*time.Minute, true)
	assert.NoError(t, err)

	err = gem.MarkEndpointUnavailableForWrite(*endpoint)
	assert.NoError(t, err)

	unavailable := gem.IsEndpointUnavailable(*endpoint, 2)
	assert.True(t, unavailable)
}

func TestGlobalEndpointManagerGetEndpointLocation(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()

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
	assert.NoError(t, err)

	srv.SetResponse(mock.WithStatusCode(200))
	srv.SetResponse(mock.WithBody(jsonString))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	serverEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{}, 5*time.Minute, true)
	assert.NoError(t, err)

	err = gem.Update(context.Background(), false)
	assert.NoError(t, err)

	location := gem.GetEndpointLocation(*serverEndpoint)

	expectedLocation := "West US"
	assert.Equal(t, expectedLocation, location)
}

func TestGlobalEndpointManagerGetAccountProperties(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{"West US", "Central US"}, 5*time.Minute, true)
	assert.NoError(t, err)

	accountProps, err := gem.GetAccountProperties(context.Background())
	assert.NoError(t, err)

	expectedAccountProps := accountProperties{
		ReadRegions:                  nil,
		WriteRegions:                 nil,
		EnableMultipleWriteLocations: false,
	}
	assert.Equal(t, expectedAccountProps, accountProps)
}

func TestGlobalEndpointManagerCanUseMultipleWriteLocations(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), internal: internalClient}

	preferredRegions := []string{"West US", "Central US"}

	serverEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	mockLc := newLocationCache(preferredRegions, *serverEndpoint, true)
	mockLc.enableMultipleWriteLocations = true

	mockGem := globalEndpointManager{
		clientEndpoint:      client.endpoint,
		preferredLocations:  preferredRegions,
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	gem, err := newGlobalEndpointManager(srv.URL(), internalClient.Pipeline(), []string{}, 5*time.Minute, true)
	assert.NoError(t, err)

	// Multiple locations should be false for default GEM
	canUseMultipleWriteLocs := gem.CanUseMultipleWriteLocations()
	assert.False(t, canUseMultipleWriteLocs)

	// Mock GEM with multiple write locations available should show true
	canUseMultipleWriteLocs = mockGem.CanUseMultipleWriteLocations()
	assert.True(t, canUseMultipleWriteLocs)
}

func TestGlobalEndpointManagerConcurrentUpdate(t *testing.T) {
	countPolicy := &countPolicy{}
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()
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

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{countPolicy}}, &policy.ClientOptions{Transport: srv})

	gem, err := newGlobalEndpointManager(srv.URL(), pl, []string{}, 5*time.Second, true)
	assert.NoError(t, err)

	// Call update concurrently and see how many times the policy gets called
	concurrency := 5
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// Call the function in each goroutine
			err := gem.Update(context.Background(), false)
			assert.NoError(t, err)
		}(wg)
	}

	wg.Wait()

	// Check that the function was called the right number of times
	callCount := countPolicy.callCount
	assert.Equal(t, callCount, 1)

	err = gem.Update(context.Background(), false)
	assert.NoError(t, err)
	callCount = countPolicy.callCount
	assert.Equal(t, callCount, 1)

	time.Sleep(5 * time.Second)

	err = gem.Update(context.Background(), false)
	assert.NoError(t, err)
	callCount = countPolicy.callCount
	assert.Equal(t, callCount, 2)
}

func TestGlobalEndpointManagerResolveEndpointSingleMasterDocumentOperation(t *testing.T) {
	serverEndpoint, _ := url.Parse("https://myaccount.documents.azure.com:443/")

	mockLc := createLocationCacheForGem(*serverEndpoint, false)

	mockGem := globalEndpointManager{
		clientEndpoint:      "https://localhost",
		preferredLocations:  []string{"Central US"},
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	// Reads should follow preferred locations
	writeOperation := false
	selectedEndpoint := mockGem.ResolveServiceEndpoint(0, resourceTypeDocument, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "central-us"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeDocument, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "east-us"))
}

func TestGlobalEndpointManagerResolveEndpointMultiMasterDocumentOperation(t *testing.T) {
	serverEndpoint, _ := url.Parse("https://myaccount.documents.azure.com:443/")

	mockLc := createLocationCacheForGem(*serverEndpoint, true)

	mockGem := globalEndpointManager{
		clientEndpoint:      "https://localhost",
		preferredLocations:  []string{"Central US"},
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	// Reads and Writes should follow preferred locations
	writeOperation := false
	selectedEndpoint := mockGem.ResolveServiceEndpoint(0, resourceTypeDocument, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "central-us"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeDocument, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "central-us"))
}

func TestGlobalEndpointManagerResolveEndpointSingleMasterMetadataOperation(t *testing.T) {
	serverEndpoint, _ := url.Parse("https://myaccount.documents.azure.com:443/")

	mockLc := createLocationCacheForGem(*serverEndpoint, false)

	mockGem := globalEndpointManager{
		clientEndpoint:      "https://localhost",
		preferredLocations:  []string{"Central US"},
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	// Reads should follow preferred locations
	writeOperation := false
	selectedEndpoint := mockGem.ResolveServiceEndpoint(0, resourceTypeCollection, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "central-us"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeCollection, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "east-us"))
}

func TestGlobalEndpointManagerResolveEndpointMultiMasterMetadataOperation(t *testing.T) {
	serverEndpoint, _ := url.Parse("https://myaccount.documents.azure.com:443/")

	mockLc := createLocationCacheForGem(*serverEndpoint, true)

	mockGem := globalEndpointManager{
		clientEndpoint:      "https://localhost",
		preferredLocations:  []string{"Central US"},
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	// Reads should follow preferred locations
	writeOperation := false
	selectedEndpoint := mockGem.ResolveServiceEndpoint(0, resourceTypeCollection, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "central-us"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeCollection, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "east-us"))
}

func createLocationCacheForGem(defaultEndpoint url.URL, isMultiMaster bool) *locationCache {
	availableWriteLocs := []string{"East US"}
	if isMultiMaster {
		availableWriteLocs = []string{"East US", "Central US"}
	}
	availableReadLocs := []string{"East US", "Central US", "East US 2"}
	availableWriteEndpointsByLoc := map[string]url.URL{}
	availableReadEndpointsByLoc := map[string]url.URL{}
	writeEndpoints := []url.URL{}
	readEndpoints := []url.URL{}

	for _, value := range availableWriteLocs {
		regionalEndpoint, _ := url.Parse(defaultEndpoint.Scheme + "://" + defaultEndpoint.Hostname() + "-" + strings.ToLower(strings.ReplaceAll(value, " ", "-")))
		availableWriteEndpointsByLoc[value] = *regionalEndpoint
		writeEndpoints = append(writeEndpoints, *regionalEndpoint)
	}

	for _, value := range availableReadLocs {
		regionalEndpoint, _ := url.Parse(defaultEndpoint.Scheme + "://" + defaultEndpoint.Hostname() + "-" + strings.ToLower(strings.ReplaceAll(value, " ", "-")))
		availableReadEndpointsByLoc[value] = *regionalEndpoint
		readEndpoints = append(readEndpoints, *regionalEndpoint)
	}

	dbAccountLocationInfo := &databaseAccountLocationsInfo{
		prefLocations:                 []string{"Central US"},
		availWriteLocations:           availableWriteLocs,
		availReadLocations:            availableReadLocs,
		availWriteEndpointsByLocation: availableWriteEndpointsByLoc,
		availReadEndpointsByLocation:  availableReadEndpointsByLoc,
		writeEndpoints:                writeEndpoints,
		readEndpoints:                 readEndpoints,
	}

	cache := locationCache{
		defaultEndpoint:                   defaultEndpoint,
		locationInfo:                      *dbAccountLocationInfo,
		locationUnavailabilityInfoMap:     make(map[url.URL]locationUnavailabilityInfo),
		unavailableLocationExpirationTime: defaultExpirationTime,
		enableCrossRegionRetries:          true,
		enableMultipleWriteLocations:      isMultiMaster,
	}

	// Order by preference
	cache.update(nil, nil, nil, nil)

	return &cache
}
