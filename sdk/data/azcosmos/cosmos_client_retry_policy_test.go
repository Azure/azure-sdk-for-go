// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestSessionNotAvailableSingleMaster(t *testing.T) {
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	// Setting up responses for consistent failures
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should fail since 404/1002 retries once for non-multi master accounts
	assert.Error(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 1)

	// Setting up responses for single failure
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithStatusCode(200))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should succeed since 404/1002 retries once for non-multi master accounts
	assert.NoError(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 1)

	// Testing write requests
	item := map[string]interface{}{
		"id":    "1",
		"value": "2",
	}
	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	// Setting up responses for consistent failures
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should fail since 404/1002 retries once for non-multi master accounts
	assert.Error(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 1)

	// Setting up responses for single failure
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithStatusCode(200))
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should succeed since 404/1002 retries once for non-multi master accounts
	assert.NoError(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 1)
}

func TestSessionNotAvailableMultiMaster(t *testing.T) {
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{},
		locationCache:       CreateMockLC(*defaultEndpoint, true),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	// Setting up responses for using all retries and failing
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should fail since 404/1002 retries once per available region multi master accounts (3 read regions)
	assert.Error(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 3)

	// Setting up responses for using all retries and succeeding
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithStatusCode(200))

	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should succeed since 404/1002 retries once per available region multi master accounts (3 read regions)
	assert.NoError(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 3)

	// Testing write requests
	item := map[string]interface{}{
		"id":    "1",
		"value": "2",
	}
	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	// Setting up responses for using all retries and failing
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))

	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should fail since 404/1002 retries once per available region multi master accounts (2 write regions)
	assert.Error(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 2)

	// Setting up responses for using all retries and succeeding
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithStatusCode(200))

	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should succeed since 404/1002 retries once per available region multi master accounts (2 write regions)
	assert.NoError(t, err)
	assert.True(t, retryPolicy.sessionRetryCount == 2)
}

func TestReadEndpointFailure(t *testing.T) {
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	// Setting up responses for retrying twice
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1008"),
		mock.WithStatusCode(403))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1008"),
		mock.WithStatusCode(403))
	srv.AppendResponse(
		mock.WithStatusCode(200))

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.NoError(t, err)
	assert.True(t, retryPolicy.retryCount == 2)
	// Verify region is marked as read unavailable
	assert.True(t, len(gem.locationCache.locationUnavailabilityInfoMap) == 1)
	locationKeys := []url.URL{}
	for k := range gem.locationCache.locationUnavailabilityInfoMap {
		locationKeys = append(locationKeys, k)
	}
	assert.True(t, gem.locationCache.locationUnavailabilityInfoMap[locationKeys[0]].unavailableOps == 1)
}

func TestWriteEndpointFailure(t *testing.T) {
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	item := map[string]interface{}{
		"id":    "1",
		"value": "2",
	}
	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	// Setting up responses for retrying twice
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "3"),
		mock.WithStatusCode(403))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "3"),
		mock.WithStatusCode(403))
	srv.AppendResponse(
		mock.WithStatusCode(200))

	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	assert.NoError(t, err)
	assert.True(t, retryPolicy.retryCount == 2)
	// Verify region is marked as write unavailable
	locationKeys := []url.URL{}
	for k := range gem.locationCache.locationUnavailabilityInfoMap {
		locationKeys = append(locationKeys, k)
	}
	assert.True(t, gem.locationCache.locationUnavailabilityInfoMap[locationKeys[0]].unavailableOps == 2)
}

func TestReadServiceUnavailable(t *testing.T) {
	// depends on length of preferred locations, if its write request has to be multi master
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{"East US", "Central US"},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	// Setting up responses for retrying and succeeding
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(200))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should retry twice and then succeed (2 preferred regions)
	assert.NoError(t, err)
	assert.True(t, retryPolicy.retryCount == 2)
	fmt.Println("we here 1")

	// Setting up responses for retrying and failing
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should retry twice and then fail (2 preferred regions)
	assert.Error(t, err)
	assert.True(t, retryPolicy.retryCount == 2)

	// Setting up multi master location cache to test same behavior
	client.gem.locationCache = CreateMockLC(*defaultEndpoint, true)

	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should retry twice and then fail (2 preferred regions)
	assert.Error(t, err)
	assert.True(t, retryPolicy.retryCount == 2)
}

func TestWriteServiceUnavailable(t *testing.T) {
	// depends on length of preferred locations, if its write request has to be multi master
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{"East US", "Central US"},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	item := map[string]interface{}{
		"id":    "1",
		"value": "2",
	}
	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	// Setting up responses for single master write failure
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))

	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Assert we do not retry the request since we are not multi master
	assert.Error(t, err)
	assert.True(t, retryPolicy.retryCount == 0)

	// Setting up multi master location cache to test same behavior
	client.gem.locationCache = CreateMockLC(*defaultEndpoint, true)

	// Setting up responses for retrying and succeeding, we still have one 503 saved in server responses
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(200))

	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should retry twice and then succeed (2 preferred regions)
	assert.NoError(t, err)
	assert.True(t, retryPolicy.retryCount == 2)

	// Setting up responses for retrying and failing
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))
	srv.AppendResponse(
		mock.WithStatusCode(503))

	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should retry twice and then fail (2 preferred regions)
	assert.Error(t, err)
	assert.True(t, retryPolicy.retryCount == 2)
}

func TestDnsErrorRetry(t *testing.T) {
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), pipeline: pl, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	// Setting up responses for retrying and succeeding, we still have one 503 saved in server responses
	DNSerr := &net.DNSError{}
	srv.AppendError(DNSerr)
	srv.AppendError(DNSerr)
	srv.AppendResponse(
		mock.WithStatusCode(200))

	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should retry twice and then succeed
	assert.NoError(t, err)
	assert.True(t, retryPolicy.retryCount == 2)

}

func CreateMockLC(defaultEndpoint url.URL, isMultiMaster bool) *locationCache {
	availableWriteLocs := []string{"East US"}
	if isMultiMaster {
		availableWriteLocs = []string{"East US", "Central US"}
	}
	availableReadLocs := []string{"East US", "Central US", "East US 2"}
	availableWriteEndpointsByLoc := map[string]url.URL{}
	availableReadEndpointsByLoc := map[string]url.URL{}
	dereferencedEndpoint := defaultEndpoint

	for _, value := range availableWriteLocs {
		availableWriteEndpointsByLoc[value] = defaultEndpoint
	}

	for _, value := range availableReadLocs {
		availableReadEndpointsByLoc[value] = defaultEndpoint
	}

	dbAccountLocationInfo := &databaseAccountLocationsInfo{
		prefLocations:                 []string{},
		availWriteLocations:           availableWriteLocs,
		availReadLocations:            availableReadLocs,
		availWriteEndpointsByLocation: availableWriteEndpointsByLoc,
		availReadEndpointsByLocation:  availableReadEndpointsByLoc,
		writeEndpoints:                []url.URL{dereferencedEndpoint},
		readEndpoints:                 []url.URL{dereferencedEndpoint},
	}

	return &locationCache{
		defaultEndpoint:                   defaultEndpoint,
		locationInfo:                      *dbAccountLocationInfo,
		locationUnavailabilityInfoMap:     make(map[url.URL]locationUnavailabilityInfo),
		unavailableLocationExpirationTime: defaultExpirationTime,
		enableCrossRegionRetries:          true,
		enableMultipleWriteLocations:      isMultiMaster,
	}

}
