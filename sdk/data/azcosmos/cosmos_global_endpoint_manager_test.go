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

	expectedLocation := newRegionId("West US")
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

	assert.True(t, strings.Contains(selectedEndpoint.Host, "centralus"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeDocument, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "eastus"))
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

	assert.True(t, strings.Contains(selectedEndpoint.Host, "centralus"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeDocument, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "centralus"))
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

	assert.True(t, strings.Contains(selectedEndpoint.Host, "centralus"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeCollection, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "eastus"))
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

	assert.True(t, strings.Contains(selectedEndpoint.Host, "centralus"))

	// Writes should go to primary endpoint
	writeOperation = true
	selectedEndpoint = mockGem.ResolveServiceEndpoint(0, resourceTypeCollection, writeOperation, false)

	assert.True(t, strings.Contains(selectedEndpoint.Host, "eastus"))
}

// A policy that captures all requests made.
type requestCollector struct {
	CapturedRequests []*policy.Request
}

func (p *requestCollector) Do(req *policy.Request) (*http.Response, error) {
	p.CapturedRequests = append(p.CapturedRequests, req)
	return req.Next()
}

func TestRequestToUpdateGEMPreservesIncomingContextWithoutCancellation(t *testing.T) {
	type contextKey string

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	// The GEM needs it's own pipeline that doesn't have the GEM policy in it to avoid deadlocking.
	capturePolicy := &requestCollector{}
	gemPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer, PerCallPolicies: []policy.Policy{capturePolicy}})
	mockGem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            gemPipeline,
		preferredLocations:  []string{"Central US"},
		locationCache:       &locationCache{},
		refreshTimeInterval: 5 * time.Minute,
	}

	gemPolicy := &globalEndpointManagerPolicy{
		gem: mockGem,
	}

	// For the "main" pipeline under test, we can insert the GEM policy, which will cause GEM updates to run (through the GEM pipeline).
	testPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer, PerCallPolicies: []policy.Policy{gemPolicy}})

	// Create a context so we can track that it flows through.
	// The context has a test value which SHOULD be preserved, and then we cancel it before even issuing the request.
	// This allows us to verify that the GEM update proceeds, even if the request is canceled.
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), contextKey("test"), "testValue"))
	cancel()

	// Issue a test request
	req, err := azruntime.NewRequest(ctx, http.MethodGet, gemServer.URL())
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = testPipeline.Do(req)

	// The _main_ request should correctly have been canceled.
	// If the GEM request had been cancelled, the error would be the "failed to retrieve account properties" error GEM returns.
	if err != context.Canceled {
		t.Fatalf("expected context to be canceled, got %v", err)
	}

	// Make sure we actually got a request to get account properties
	if len(capturePolicy.CapturedRequests) != 1 {
		t.Fatalf("expected to capture the request to the GEM, got %d requests", len(capturePolicy.CapturedRequests))
	}
	capturedReq := capturePolicy.CapturedRequests[0]
	if capturedReq.Raw().URL.String() != gemServer.URL() {
		t.Fatalf("expected the captured request to be to the account metadata endpoint, got %s", capturedReq.Raw().URL.String())
	}
	if capturedReq.Raw().Method != http.MethodGet {
		t.Fatalf("expected the captured request to be a GET, got %s", capturedReq.Raw().Method)
	}

	// Validate that the context of THAT request is non-canceled and has our test value.
	capturedContext := capturedReq.Raw().Context()
	if _, ok := capturedContext.Deadline(); !ok {
		t.Fatalf("expected the context to not have a deadline")
	}
	value := capturedContext.Value(contextKey("test"))
	if value != "testValue" {
		t.Fatalf("expected a captured context to contain test=testValue, got test=%v", value)
	}
}

func TestAddedAllowTentativeHeaderGEMPolicy(t *testing.T) {
	type contextKey string

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))
	serverEndpoint, _ := url.Parse("https://myaccount.documents.azure.com:443/")
	mocked_response := "{\"_self\":\"\",\"id\":\"my_account\",\"_rid\":\"my_account-westus.sql.cosmos.azure.com\",\"media\":\"//media/\",\"addresses\":\"//addresses/\",\"_dbs\":\"//dbs/\",\"writableLocations\":[{\"name\":\"West US\",\"databaseAccountEndpoint\":\"https://my_account-westus.documents.azure.com:443/\"},{\"name\":\"West US 3\",\"databaseAccountEndpoint\":\"https://my_account-westus3.documents.azure.com:443/\"}],\"readableLocations\":[{\"name\":\"West US\",\"databaseAccountEndpoint\":\"https://my_account-westus.documents.azure.com:443/\"},{\"name\":\"West US 3\",\"databaseAccountEndpoint\":\"https://my_account-westus3.documents.azure.com:443/\"}], \"enableMultipleWriteLocations\":true}"

	gemServer.SetResponse(mock.WithBody([]byte(mocked_response)))
	mockLc := createLocationCacheForGem(*serverEndpoint, true)

	// The GEM needs it's own pipeline that doesn't have the GEM policy in it to avoid deadlocking.
	capturePolicy := &requestCollector{}
	gemPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer, PerCallPolicies: []policy.Policy{capturePolicy}})
	mockGem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            gemPipeline,
		preferredLocations:  []string{"Central US"},
		locationCache:       mockLc,
		refreshTimeInterval: 5 * time.Minute,
	}

	gemPolicy := &globalEndpointManagerPolicy{
		gem: mockGem,
	}

	// For the "main" pipeline under test, we can insert the GEM policy, which will cause GEM updates to run (through the GEM pipeline).
	testPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer, PerCallPolicies: []policy.Policy{gemPolicy}})

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), contextKey("test"), "testValue"))
	defer cancel()

	// Issue a test request
	req, err := azruntime.NewRequest(ctx, http.MethodGet, gemServer.URL())
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, _ := testPipeline.Do(req)
	// tentative write header should be sent for multi write account
	if resp.Request.Header.Get(cosmosHeaderAllowTentativeWrites) == "" {
		t.Fatalf("expected %s header to be set", cosmosHeaderAllowTentativeWrites)
	}

	// tentative write header should not be sent if the account is not multi-write
	mocked_response = "{\"_self\":\"\",\"id\":\"my_account\",\"_rid\":\"my_account-westus.sql.cosmos.azure.com\",\"media\":\"//media/\",\"addresses\":\"//addresses/\",\"_dbs\":\"//dbs/\",\"writableLocations\":[{\"name\":\"West US\",\"databaseAccountEndpoint\":\"https://my_account-westus.documents.azure.com:443/\"},{\"name\":\"West US 3\",\"databaseAccountEndpoint\":\"https://my_account-westus3.documents.azure.com:443/\"}],\"readableLocations\":[{\"name\":\"West US\",\"databaseAccountEndpoint\":\"https://my_account-westus.documents.azure.com:443/\"},{\"name\":\"West US 3\",\"databaseAccountEndpoint\":\"https://my_account-westus3.documents.azure.com:443/\"}], \"enableMultipleWriteLocations\":false}"
	gemServer.SetResponse(mock.WithBody([]byte(mocked_response)))
	// change time to trigger another get account properties call
	mockGem.lastUpdateTime = time.Now().Add(-10 * time.Minute)

	// Issue another test request
	req, err = azruntime.NewRequest(ctx, http.MethodGet, gemServer.URL())
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Used to trigger another get account properties call in the background
	_, err = testPipeline.Do(req)
	if err != nil {
		t.Fatalf("testPipeline.Do failed: %v", err)
	}

	// Issue another test request that will use the updated account properties
	req, err = azruntime.NewRequest(ctx, http.MethodGet, gemServer.URL())
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, _ = testPipeline.Do(req)
	if resp.Request.Header.Get(cosmosHeaderAllowTentativeWrites) != "" {
		t.Fatalf("expected %s header not to be set", cosmosHeaderAllowTentativeWrites)
	}
}

func createLocationCacheForGem(defaultEndpoint url.URL, isMultiMaster bool) *locationCache {
	availableWriteLocs := []regionId{newRegionId("East US")}
	if isMultiMaster {
		availableWriteLocs = []regionId{newRegionId("East US"), newRegionId("Central US")}
	}
	availableReadLocs := []regionId{newRegionId("East US"), newRegionId("Central US"), newRegionId("East US 2")}
	availableWriteEndpointsByLoc := map[regionId]url.URL{}
	availableReadEndpointsByLoc := map[regionId]url.URL{}
	writeEndpoints := []url.URL{}
	readEndpoints := []url.URL{}

	for _, region := range availableWriteLocs {
		regionalEndpoint, _ := url.Parse(defaultEndpoint.Scheme + "://" + defaultEndpoint.Hostname() + "-" + region.String())
		availableWriteEndpointsByLoc[region] = *regionalEndpoint
		writeEndpoints = append(writeEndpoints, *regionalEndpoint)
	}

	for _, region := range availableReadLocs {
		regionalEndpoint, _ := url.Parse(defaultEndpoint.Scheme + "://" + defaultEndpoint.Hostname() + "-" + region.String())
		availableReadEndpointsByLoc[region] = *regionalEndpoint
		readEndpoints = append(readEndpoints, *regionalEndpoint)
	}

	dbAccountLocationInfo := &databaseAccountLocationsInfo{
		prefLocations:                 []regionId{newRegionId("Central US")},
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
	_ = cache.update(nil, nil, nil, nil)

	return &cache
}
