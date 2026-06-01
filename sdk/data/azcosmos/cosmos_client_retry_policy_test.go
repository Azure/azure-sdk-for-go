// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore azcosmosgemtest azcosmostest retriable

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	// Setting up responses for consistent failures
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should fail since 404/1002 retries once for non-multi master accounts
	assert.Error(t, err)
	assert.True(t, verifier.requests[0].retryContext.sessionRetryCount == 1)

	// Setting up responses for single failure
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithStatusCode(200))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should succeed since 404/1002 retries once for non-multi master accounts
	assert.NoError(t, err)
	assert.True(t, verifier.requests[0].retryContext.sessionRetryCount == 1)

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
	assert.True(t, verifier.requests[0].retryContext.sessionRetryCount == 1)

	// Setting up responses for single failure
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1002"),
		mock.WithStatusCode(404))
	srv.AppendResponse(
		mock.WithStatusCode(200))
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	// Request should succeed since 404/1002 retries once for non-multi master accounts
	assert.NoError(t, err)
	assert.True(t, verifier.requests[0].retryContext.sessionRetryCount == 1)
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

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

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should fail since 404/1002 retries once per available region multi master accounts (3 read regions)
	assert.Error(t, err)
	assert.True(t, verifier.requests[0].retryContext.sessionRetryCount == 3)

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
	assert.True(t, verifier.requests[1].retryContext.sessionRetryCount == 3)

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
	assert.True(t, verifier.requests[2].retryContext.sessionRetryCount == 2)

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
	assert.True(t, verifier.requests[3].retryContext.sessionRetryCount == 2)
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	// Setting up responses for retrying twice
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1008"),
		mock.WithStatusCode(403))
	srv.AppendResponse(
		mock.WithHeader("x-ms-substatus", "1008"),
		mock.WithStatusCode(403))
	srv.AppendResponse(
		mock.WithStatusCode(200))

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.NoError(t, err)
	assert.True(t, verifier.requests[0].retryContext.retryCount == 2)
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
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
	assert.True(t, verifier.requests[0].retryContext.retryCount == 2)
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
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
	assert.True(t, verifier.requests[0].retryContext.retryCount == 2)

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
	assert.True(t, verifier.requests[0].retryContext.retryCount == 2)

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
	assert.True(t, verifier.requests[1].retryContext.retryCount == 2)
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
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
	assert.True(t, verifier.requests[0].retryContext.retryCount == 0)

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
	assert.True(t, verifier.requests[1].retryContext.retryCount == 2)

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
	assert.True(t, verifier.requests[2].retryContext.retryCount == 2)
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	// Setting up responses for retrying and succeeding, we still have one 503 saved in server responses
	DNSerr := &net.DNSError{}
	srv.AppendError(DNSerr)
	srv.AppendError(DNSerr)
	srv.AppendResponse(
		mock.WithStatusCode(200))

	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	// Request should retry twice on the same region and then succeed; no
	// cross-region failover, so retryCount stays at 0 and
	// sameRegionRetryCount reflects the two retries.
	assert.NoError(t, err)
	assert.Equal(t, 0, verifier.requests[0].retryContext.retryCount)
	assert.Equal(t, 2, verifier.requests[0].retryContext.sameRegionRetryCount)

}

// setupRetryPolicyTestClient creates a Client wired to a single mock cosmos
// server with the client retry policy under test. Tests can append responses
// (or errors) to `srv` to drive the retry behavior.
func setupRetryPolicyTestClient(t *testing.T) (*Client, *mock.Server, *clientRetryPolicyVerifier, func()) {
	return setupRetryPolicyTestClientOpts(t, true /*multiMaster*/, true /*enableCrossRegion*/)
}

func setupRetryPolicyTestClientOpts(t *testing.T, multiMaster, enableCrossRegion bool) (*Client, *mock.Server, *clientRetryPolicyVerifier, func()) {
	t.Helper()
	srv, srvClose := mock.NewTLSServer()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	lc := CreateMockLC(*defaultEndpoint, multiMaster)
	lc.enableCrossRegionRetries = enableCrossRegion

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{},
		locationCache:       lc,
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}
	verifier := &clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	cleanup := func() {
		srvClose()
		gemClose()
	}
	return client, srv, verifier, cleanup
}

func TestConnectionErrorReadFailsOverAfterThreeSameRegionAttempts(t *testing.T) {
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	dnsErr := &net.DNSError{}
	for i := 0; i < 4; i++ {
		srv.AppendError(dnsErr)
	}
	srv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.NoError(t, err)
	rc := verifier.requests[0].retryContext
	// 3 same-region attempts exhausted, then a cross-region failover that
	// succeeded. After failover sameRegionRetryCount is reset and
	// retryCount is incremented to pick a different endpoint.
	assert.Equal(t, 0, rc.sameRegionRetryCount)
	assert.Equal(t, 0, rc.retryCount) // post-fix: retryCount not incremented on connection-error failover; demote-in-cache handles routing
}

// TestConnectionErrorReadFailsOverWhenGlobalEndpointIsUnreachable simulates
// a regional gateway outage where the global account endpoint also resolves
// to the same regional FE pool that has been blocked (the common case for
// single-region writes — global FQDN points at the write region's FE).
//
// Before the fix, attemptRetryOnNetworkError had three interlocking
// problems that prevented the cross-region failover from ever taking effect:
//  1. It forced a synchronous gem.Update(ctx, true) after
//     MarkEndpointUnavailable*. With the global endpoint unreachable, the
//     refresh failed (in production typically a connect timeout when the
//     global FQDN resolves to a blocked regional FE pool; the test injects
//     a net.DNSError as a deterministic stand-in for any gem.Update
//     failure) — causing the policy to surface the original connection
//     failure without ever attempting the cross-region retry.
//  2. It incremented retryCount after the mark. MarkEndpointUnavailable*
//     demotes the bad endpoint to the TAIL of readEndpoints rather than
//     removing it, so readEndpoints becomes [good, bad]. With retryCount
//     bumped to 1, ResolveServiceEndpoint(1 % 2) returns the still-bad
//     endpoint at the tail — the failover attempt would hit the same
//     dead region again.
//  3. MarkEndpointUnavailable* was called with the full request URL
//     (path, query, etc. included) but the unavailability map and the
//     cache's per-region endpoint lookup were keyed by base URLs
//     (scheme+host). The marks were therefore written under keys nothing
//     else looked up, so isEndpointUnavailableLocked always returned false
//     and the demote silently did nothing.
//
// The fix drops the forced refresh, leaves retryCount at 0 so the next
// ResolveServiceEndpoint returns readEndpoints[0] (the just-promoted
// preferred region), and normalizes URLs to scheme+host on both write
// and read sides of the unavailability map.
//
// To actually exercise the routing this test wires up TWO distinct mock
// servers (badSrv = original/unhealthy region, goodSrv = failover region)
// and points the location cache's read endpoints at both. badSrv only
// serves DNS errors; goodSrv serves the 200 the request needs. If the
// resolver returns badSrv after failover (because of any of the three
// pre-fix conditions) the test fails.
func TestConnectionErrorReadFailsOverWhenGlobalEndpointIsUnreachable(t *testing.T) {
	badSrv, badClose := mock.NewTLSServer()
	defer badClose()
	goodSrv, goodClose := mock.NewTLSServer()
	defer goodClose()

	badURL, err := url.Parse(badSrv.URL())
	require.NoError(t, err)
	goodURL, err := url.Parse(goodSrv.URL())
	require.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	// Simulate the global endpoint being unreachable for the duration of
	// the regional outage. In production this typically manifests as a
	// connect timeout (global FQDN resolves to a blocked regional FE
	// pool); a net.DNSError gives us the same gem.Update(ctx,true)
	// failure deterministically and without test-time sleeps.
	gemServer.SetError(&net.DNSError{})

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	// Build a location cache with TWO distinct regional endpoints so the
	// routing decision after failover is observable. "East US" (badSrv)
	// is the user's application region (index 0); "Central US" (goodSrv)
	// is the next preferred.
	lc := newLocationCache([]string{"East US", "Central US"}, *badURL, true /*enableCrossRegionRetries*/)
	require.NoError(t, lc.update(
		[]accountRegion{{Name: "East US", Endpoint: badSrv.URL()}},
		[]accountRegion{
			{Name: "East US", Endpoint: badSrv.URL()},
			{Name: "Central US", Endpoint: goodSrv.URL()},
		},
		[]string{"East US", "Central US"},
		nil,
	))

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{"East US", "Central US"},
		locationCache:       lc,
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	// azcore needs to dispatch to whichever URL the policy resolves to,
	// not a fixed Transport. routingMockTransport keys by host so a single
	// client sees distinct backing servers per region.
	routingTransport := routingMockTransport{
		byHost: map[string]*mock.Server{
			badURL.Host:  badSrv,
			goodURL.Host: goodSrv,
		},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}
	verifier := &clientRetryPolicyVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{verifier, retryPolicy}}, &policy.ClientOptions{Transport: &routingTransport})
	client := &Client{endpoint: badSrv.URL(), endpointUrl: badURL, internal: internalClient, gem: gem}

	dnsErr := &net.DNSError{}
	// 1 initial + 3 same-region retries on the bad region.
	for i := 0; i < 4; i++ {
		badSrv.AppendError(dnsErr)
	}
	// Cross-region failover should hit the good region.
	goodSrv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	require.NoError(t, err, "cross-region failover should reach the good region")
	rc := verifier.requests[0].retryContext
	assert.True(t, rc.crossRegionFailoverDone, "expected one cross-region failover")
	// retryCount stays at 0: MarkEndpointUnavailable* demoted the bad
	// endpoint so ResolveServiceEndpoint(0) now returns the good region.
	// Bumping retryCount to 1 would index back to the demoted-tail slot.
	assert.Equal(t, 0, rc.retryCount)
	assert.Equal(t, 0, rc.sameRegionRetryCount)
	// 1 initial + 3 same-region retries against badSrv.
	assert.Equal(t, 4, badSrv.Requests())
	// Exactly one request against the good region (the failover).
	assert.Equal(t, 1, goodSrv.Requests())
}

// routingMockTransport routes each request to the mock server matching
// the request URL's host. This lets a single client see distinct backing
// servers per region without azcore short-circuiting to a fixed mock.
type routingMockTransport struct {
	byHost map[string]*mock.Server
}

func (r *routingMockTransport) Do(req *http.Request) (*http.Response, error) {
	srv, ok := r.byHost[req.URL.Host]
	if !ok {
		return nil, fmt.Errorf("no mock server registered for host %q", req.URL.Host)
	}
	return srv.Do(req)
}

func TestNotSentConnectionErrorWriteFailsOver(t *testing.T) {
	// Multi-master account: writes can fail over to another write region
	// when the failure is classified as not-sent (DNS in this case).
	client, srv, verifier, cleanup := setupRetryPolicyTestClientOpts(t, true /*multiMaster*/, true)
	defer cleanup()

	dnsErr := &net.DNSError{}
	for i := 0; i < 4; i++ {
		srv.AppendError(dnsErr)
	}
	srv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, mErr := json.Marshal(item)
	require.NoError(t, mErr)
	_, err := container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	assert.NoError(t, err)
	rc := verifier.requests[0].retryContext
	assert.Equal(t, 0, rc.sameRegionRetryCount)
	assert.Equal(t, 0, rc.retryCount) // post-fix: retryCount not incremented on connection-error failover; demote-in-cache handles routing
}

// fakeAmbiguousNetError satisfies net.Error and wraps syscall.ECONNRESET
// so classifyNetworkError categorizes it as connectionErrorAmbiguous. It
// does NOT implement *net.OpError; tests that want to exercise the
// real OpError code path should use &net.OpError{Op: "read", ...} directly.
type fakeAmbiguousNetError struct{ msg string }

func (e *fakeAmbiguousNetError) Error() string   { return e.msg }
func (e *fakeAmbiguousNetError) Timeout() bool   { return false }
func (e *fakeAmbiguousNetError) Temporary() bool { return false }
func (e *fakeAmbiguousNetError) Unwrap() error   { return syscall.ECONNRESET }

func TestAmbiguousConnectionErrorWriteDoesNotFailOver(t *testing.T) {
	// Ambiguous transport errors on a write must NOT be retried at all by
	// the client retry policy (neither same-region nor cross-region) —
	// the request body may already have been applied server-side and
	// retrying could produce duplicate mutations (e.g. PatchItem with
	// Increment, TransactionalBatch with non-idempotent ops).
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	ambErr := &fakeAmbiguousNetError{msg: "connection reset by peer"}
	// More errors than we should ever consume.
	for i := 0; i < 6; i++ {
		srv.AppendError(ambErr)
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, mErr := json.Marshal(item)
	require.NoError(t, mErr)
	_, err := container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	assert.Error(t, err)
	rc := verifier.requests[0].retryContext
	// No same-region retries, no cross-region failover; exactly 1 HTTP
	// request reached the server.
	assert.Equal(t, 0, rc.sameRegionRetryCount)
	assert.Equal(t, 0, rc.retryCount)
	assert.False(t, rc.crossRegionFailoverDone)
	assert.Equal(t, 1, srv.Requests())
}

func TestAmbiguousConnectionErrorReadFailsOver(t *testing.T) {
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	ambErr := &fakeAmbiguousNetError{msg: "connection reset by peer"}
	for i := 0; i < 4; i++ {
		srv.AppendError(ambErr)
	}
	srv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.NoError(t, err)
	rc := verifier.requests[0].retryContext
	assert.Equal(t, 0, rc.sameRegionRetryCount)
	assert.Equal(t, 0, rc.retryCount) // post-fix: retryCount not incremented on connection-error failover; demote-in-cache handles routing
}

func TestCallerDeadlineExceededDoesNotRetry(t *testing.T) {
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	srv.AppendError(&net.DNSError{})

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(ctx, NewPartitionKeyString("1"), "doc1", nil)

	assert.Error(t, err)
	// No retries should have been attempted.
	rc := verifier.requests[0].retryContext
	assert.Equal(t, 0, rc.sameRegionRetryCount)
	assert.Equal(t, 0, rc.retryCount)
	assert.Equal(t, 1, len(verifier.requests))
}

func TestNotSentConnectionErrorMultiMasterWriteFailsOver(t *testing.T) {
	// Multi-master account: writes can fail over to a different write region.
	// 3 same-region DNS errors + 1 cross-region failover that succeeds.
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	dnsErr := &net.DNSError{}
	for i := 0; i < 4; i++ {
		srv.AppendError(dnsErr)
	}
	srv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, mErr := json.Marshal(item)
	require.NoError(t, mErr)
	_, err := container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	assert.NoError(t, err)
	rc := verifier.requests[0].retryContext
	assert.True(t, rc.crossRegionFailoverDone, "expected one cross-region failover")
	assert.Equal(t, 0, rc.retryCount) // post-fix: retryCount not incremented on connection-error failover; demote-in-cache handles routing
	assert.Equal(t, 0, rc.sameRegionRetryCount)
}

func TestConnectionErrorGivesUpAfterSingleCrossRegionFailover(t *testing.T) {
	// 3 same-region errors + 1 cross-region failover that ALSO errors →
	// the policy must NOT chain a second failover; it returns the error.
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	dnsErr := &net.DNSError{}
	// More errors than we should ever consume (3 same-region + 1 failover = 4).
	for i := 0; i < 8; i++ {
		srv.AppendError(dnsErr)
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.Error(t, err)
	rc := verifier.requests[0].retryContext
	// One cross-region failover happened and then we gave up.
	assert.True(t, rc.crossRegionFailoverDone)
	assert.Equal(t, 0, rc.retryCount) // post-fix: retryCount not incremented on connection-error failover; demote-in-cache handles routing
	// Mock server should have served exactly 5 requests:
	// 1 initial + 3 same-region retries + 1 cross-region failover.
	assert.Equal(t, 5, srv.Requests())
}

func TestRequestTimeoutReadRetriesCrossRegion(t *testing.T) {
	// One 408 on a read → single cross-region retry that succeeds.
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.NoError(t, err)
	rc := verifier.requests[0].retryContext
	assert.True(t, rc.requestTimeoutRetryDone)
	assert.Equal(t, 1, rc.retryCount)
	assert.Equal(t, 2, srv.Requests())
}

func TestRequestTimeoutReadGivesUpAfterOneCrossRegionRetry(t *testing.T) {
	// Two consecutive 408s on a read → exactly one cross-region retry,
	// then the policy returns the 408 as non-retriable.
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	for i := 0; i < 4; i++ {
		srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	assert.Error(t, err)
	rc := verifier.requests[0].retryContext
	assert.True(t, rc.requestTimeoutRetryDone)
	assert.Equal(t, 1, rc.retryCount)
	// 1 initial + 1 cross-region retry = 2 requests.
	assert.Equal(t, 2, srv.Requests())
}

func TestRequestTimeoutWriteDoesNotRetry(t *testing.T) {
	// A 408 on a write should NOT be retried — writes on 408 are
	// ambiguous and could lead to duplicates.
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(200))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, mErr := json.Marshal(item)
	require.NoError(t, mErr)
	_, err := container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	assert.Error(t, err)
	rc := verifier.requests[0].retryContext
	assert.False(t, rc.requestTimeoutRetryDone)
	assert.Equal(t, 0, rc.retryCount)
	assert.Equal(t, 1, srv.Requests())
}

func TestSingleMasterWriteDoesNotFailoverOnConnectionError(t *testing.T) {
	// Single-master: writes have only one possible write region, so a
	// "cross-region failover" would just route back to the same region.
	// The policy should give up after 3 same-region attempts without
	// performing a wasted cross-region attempt and without marking the
	// only write endpoint unavailable for write.
	client, srv, verifier, cleanup := setupRetryPolicyTestClientOpts(t, false /*multiMaster*/, true)
	defer cleanup()

	dnsErr := &net.DNSError{}
	for i := 0; i < 6; i++ {
		srv.AppendError(dnsErr)
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, err := json.Marshal(item)
	require.NoError(t, err)
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	require.Error(t, err)
	rc := verifier.requests[0].retryContext
	// 1 initial + 3 same-region retries = 4 total; no cross-region attempt.
	assert.Equal(t, 4, srv.Requests())
	assert.False(t, rc.crossRegionFailoverDone)
	// No endpoint should be marked unavailable for write on a single-master
	// account — there is nowhere else to send writes.
	for _, info := range client.gem.locationCache.locationUnavailabilityInfoMap {
		assert.NotEqual(t, write, info.unavailableOps, "single-master write endpoint should not be marked write-unavailable")
		assert.NotEqual(t, all, info.unavailableOps, "single-master write endpoint should not be marked all-unavailable")
	}
}

func TestAmbiguousWriteMarksEndpointUnavailableForRead(t *testing.T) {
	// Multi-master write that gives up on an ambiguous transport error
	// should still mark the endpoint unavailable for read so concurrent
	// requests learn about the regional outage.
	client, srv, verifier, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	ambErr := &fakeAmbiguousNetError{msg: "connection reset by peer"}
	for i := 0; i < 6; i++ {
		srv.AppendError(ambErr)
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, err := json.Marshal(item)
	require.NoError(t, err)
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)

	// At least one endpoint must have been marked unavailable for write
	// (single-master would NOT do this; we use multi-master here).
	require.Error(t, err)
	rc := verifier.requests[0].retryContext
	assert.False(t, rc.crossRegionFailoverDone)
	// Marked unavailable for read for at least one endpoint.
	var markedForRead bool
	for _, info := range client.gem.locationCache.locationUnavailabilityInfoMap {
		if info.unavailableOps == read || info.unavailableOps == all {
			markedForRead = true
			break
		}
	}
	assert.True(t, markedForRead, "expected at least one endpoint marked unavailable for read")
}

func TestConnectionErrorWithCrossRegionRetriesDisabledFailsFast(t *testing.T) {
	// With enableCrossRegionRetries=false the policy must not retry at
	// all — neither same-region nor cross-region — preserving the
	// pre-existing "fail fast" semantics.
	client, srv, _, cleanup := setupRetryPolicyTestClientOpts(t, true, false /*disable cross region*/)
	defer cleanup()

	for i := 0; i < 4; i++ {
		srv.AppendError(&net.DNSError{})
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)

	require.Error(t, err)
	assert.Equal(t, 1, srv.Requests(), "no retries should be performed when cross-region retries are disabled")
}

func TestCallerDeadlineDuringBackoffShortCircuits(t *testing.T) {
	// While the policy is sleeping between same-region retries, if the
	// caller's context expires the sleep must return early and the
	// policy must give up.
	client, srv, _, cleanup := setupRetryPolicyTestClient(t)
	defer cleanup()

	for i := 0; i < 6; i++ {
		srv.AppendError(&net.DNSError{})
	}

	// Deadline shorter than defaultBackoff so the first backoff is
	// guaranteed to be interrupted by ctx cancellation.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	start := time.Now()
	_, err := container.ReadItem(ctx, NewPartitionKeyString("1"), "doc1", nil)
	elapsed := time.Since(start)

	require.Error(t, err)
	// The cancellation cause must be preserved so callers can
	// errors.Is(err, context.DeadlineExceeded) regardless of whether
	// the deadline fired before req.Next() or during backoff.
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	// We should give up before the full 3-second same-region budget.
	assert.Less(t, elapsed, defaultBackoff*time.Second, "policy should not sleep through caller deadline")
	// At most 2 requests: the initial attempt and possibly one more
	// before the deadline fires during backoff.
	assert.LessOrEqual(t, srv.Requests(), 2)
}

func TestClassifyNetworkError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want connectionErrorKind
	}{
		{"nil", nil, connectionErrorNone},
		{"dns", &net.DNSError{}, connectionErrorNotSent},
		{"dial op", &net.OpError{Op: "dial", Err: errors.New("boom")}, connectionErrorNotSent},
		{"connection refused", syscall.ECONNREFUSED, connectionErrorNotSent},
		{"host unreachable", syscall.EHOSTUNREACH, connectionErrorNotSent},
		{"etimedout", syscall.ETIMEDOUT, connectionErrorNotSent},
		{"real opError read wrapping econnreset", &net.OpError{Op: "read", Err: syscall.ECONNRESET}, connectionErrorAmbiguous},
		{"eof", io.EOF, connectionErrorAmbiguous},
		{"unexpected eof", io.ErrUnexpectedEOF, connectionErrorAmbiguous},
		{"connection reset", syscall.ECONNRESET, connectionErrorAmbiguous},
		{"deadline exceeded", context.DeadlineExceeded, connectionErrorAmbiguous},
		{"read op", &net.OpError{Op: "read", Err: errors.New("boom")}, connectionErrorAmbiguous},
		{"plain error", errors.New("nope"), connectionErrorNone},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, classifyNetworkError(tc.err))
		})
	}
}

// TestWriteForbiddenFailsOverToHealthyRegion is the routing-level
// regression for the 403/WriteForbidden path. It mirrors
// TestConnectionErrorReadFailsOverWhenGlobalEndpointIsUnreachable for
// the network-error path: two distinct backend mock servers wired
// through a host-routing transport, the first returns
// 403/WriteForbidden, and the failover must reach the second.
//
// Before this PR also fixed the 403 path, MarkEndpointUnavailable*
// demoted the bad write endpoint to the tail of writeEndpoints, then
// the outer Do() loop bumped retryContext.retryCount += 1, and the
// next ResolveServiceEndpoint(1 % 2) routed right back to the demoted
// bad endpoint. The fix sets retryContext.resolveFromHead = true so
// the next resolve uses locationIndex 0.
func TestWriteForbiddenFailsOverToHealthyRegion(t *testing.T) {
	badSrv, badClose := mock.NewTLSServer()
	defer badClose()
	goodSrv, goodClose := mock.NewTLSServer()
	defer goodClose()

	badURL, err := url.Parse(badSrv.URL())
	require.NoError(t, err)
	goodURL, err := url.Parse(goodSrv.URL())
	require.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetError(&net.DNSError{})
	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	lc := newLocationCache([]string{"East US", "Central US"}, *badURL, true /*enableCrossRegionRetries*/)
	require.NoError(t, lc.update(
		[]accountRegion{
			{Name: "East US", Endpoint: badSrv.URL()},
			{Name: "Central US", Endpoint: goodSrv.URL()},
		},
		[]accountRegion{
			{Name: "East US", Endpoint: badSrv.URL()},
			{Name: "Central US", Endpoint: goodSrv.URL()},
		},
		[]string{"East US", "Central US"},
		boolPtr(true), // enable multi-master so writes can fail over
	))

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{"East US", "Central US"},
		locationCache:       lc,
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	routingTransport := routingMockTransport{
		byHost: map[string]*mock.Server{
			badURL.Host:  badSrv,
			goodURL.Host: goodSrv,
		},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}
	verifier := &clientRetryPolicyVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{verifier, retryPolicy}}, &policy.ClientOptions{Transport: &routingTransport})
	client := &Client{endpoint: badSrv.URL(), endpointUrl: badURL, internal: internalClient, gem: gem}

	// 1 initial 403/WriteForbidden on the bad region.
	badSrv.AppendResponse(
		mock.WithHeader("x-ms-substatus", subStatusWriteForbidden),
		mock.WithStatusCode(http.StatusForbidden))
	// Cross-region failover should hit the good region.
	goodSrv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	item := map[string]interface{}{"id": "1", "value": "2"}
	marshalled, mErr := json.Marshal(item)
	require.NoError(t, mErr)
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
	require.NoError(t, err, "403/WriteForbidden must fail over to the healthy region")

	// Exactly one request on each: 1 initial 403 against badSrv, 1
	// failover success against goodSrv. A regression that re-routed to
	// the demoted endpoint would show 2 requests on badSrv.
	assert.Equal(t, 1, badSrv.Requests(), "no further requests should hit the demoted write endpoint")
	assert.Equal(t, 1, goodSrv.Requests(), "the failover request must reach the healthy write endpoint")
}

// TestConnectionErrorFailoverResetsNonZeroRetryCount covers the mixed
// failure sequence: a prior HTTP-status retry (e.g. 408) bumps
// retryCount, then a connection error triggers the cross-region
// failover. The failover must still land on the healthy region; if the
// connection-error path merely "does not increment" retryCount instead
// of forcing the next resolve to head, the inherited non-zero index
// indexes back to the demoted-tail bad endpoint.
func TestConnectionErrorFailoverResetsNonZeroRetryCount(t *testing.T) {
	badSrv, badClose := mock.NewTLSServer()
	defer badClose()
	goodSrv, goodClose := mock.NewTLSServer()
	defer goodClose()

	badURL, err := url.Parse(badSrv.URL())
	require.NoError(t, err)
	goodURL, err := url.Parse(goodSrv.URL())
	require.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetError(&net.DNSError{})
	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	lc := newLocationCache([]string{"East US", "Central US"}, *badURL, true)
	require.NoError(t, lc.update(
		[]accountRegion{{Name: "East US", Endpoint: badSrv.URL()}},
		[]accountRegion{
			{Name: "East US", Endpoint: badSrv.URL()},
			{Name: "Central US", Endpoint: goodSrv.URL()},
		},
		[]string{"East US", "Central US"},
		nil,
	))

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{"East US", "Central US"},
		locationCache:       lc,
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	routingTransport := routingMockTransport{
		byHost: map[string]*mock.Server{
			badURL.Host:  badSrv,
			goodURL.Host: goodSrv,
		},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}
	verifier := &clientRetryPolicyVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{verifier, retryPolicy}}, &policy.ClientOptions{Transport: &routingTransport})
	client := &Client{endpoint: badSrv.URL(), endpointUrl: badURL, internal: internalClient, gem: gem}

	// Sequence on badSrv:
	//   1) 408 (read) -> outer loop bumps retryCount to 1, picks
	//      readEndpoints[1] = Central US for the next attempt.
	// Sequence on goodSrv (now selected):
	//   2) initial attempt: 4x DNSError (initial + 3 same-region) ->
	//      triggers cross-region failover via attemptRetryOnNetworkError.
	// After the failover the inherited retryCount is 1 (or higher).
	// resolveFromHead must force the next resolve to index 0, which
	// after demote-to-tail of Central US is East US (badSrv).
	//   3) failover hits badSrv again: serve a 200.
	badSrv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	badSrv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	dnsErr := &net.DNSError{}
	for i := 0; i < 4; i++ {
		goodSrv.AppendError(dnsErr)
	}

	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	require.NoError(t, err, "mixed 408+connection-error sequence must still fail over to a healthy host")

	rc := verifier.requests[0].retryContext
	assert.True(t, rc.requestTimeoutRetryDone, "the 408 retry should have run")
	assert.True(t, rc.crossRegionFailoverDone, "the connection-error failover should have run")
	// 1 initial 408 + 1 final success against badSrv.
	assert.Equal(t, 2, badSrv.Requests(), "expected initial 408 + post-failover success on the head-of-list endpoint")
	// 1 initial + 3 same-region retries against goodSrv (the 408 routed us here, then DNS killed it).
	assert.Equal(t, 4, goodSrv.Requests(), "expected initial + 3 same-region attempts before failover gave up on the bad region")
}

func boolPtr(b bool) *bool { return &b }

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

type clientRetryPolicyVerifier struct {
	requests []clientRetryPolicyVerifierRequest
}

type clientRetryPolicyVerifierRequest struct {
	retryContext *retryContext
}

func (p *clientRetryPolicyVerifier) Do(req *policy.Request) (*http.Response, error) {
	resp, err := req.Next()
	pr := clientRetryPolicyVerifierRequest{}
	o := retryContext{}
	req.OperationValue(&o)
	pr.retryContext = &o
	p.requests = append(p.requests, pr)
	return resp, err
}

// multiReadEndpointLC builds a mock location cache whose readEndpoints slice has
// more than one entry, all pointing at the single mock server (defaultEndpoint).
// The server-error cross-region retry path requires len(readEndpoints) > 1; using
// identical URLs keeps every retry routed back to the same mock server so the
// queued responses are consumed in order.
func multiReadEndpointLC(defaultEndpoint url.URL, isMultiMaster bool) *locationCache {
	lc := CreateMockLC(defaultEndpoint, isMultiMaster)
	lc.locationInfo.readEndpoints = []url.URL{defaultEndpoint, defaultEndpoint, defaultEndpoint}
	return lc
}

// hostRoutingTransport dispatches each request to a backing transport selected
// by the request's URL host, and records the host of every attempt in order.
// It lets a test prove that a cross-region retry actually changes the request
// target endpoint (instead of merely advancing retry counters).
type hostRoutingTransport struct {
	routes    map[string]policy.Transporter
	seenHosts []string
}

func (t *hostRoutingTransport) Do(req *http.Request) (*http.Response, error) {
	t.seenHosts = append(t.seenHosts, req.URL.Host)
	backing, ok := t.routes[req.URL.Host]
	if !ok {
		return nil, fmt.Errorf("no backing transport registered for host %q", req.URL.Host)
	}
	return backing.Do(req)
}

// TestReadServerError_CrossRegionRoutesToDifferentEndpoint proves that the
// cross-region 5xx retry targets a genuinely different endpoint than the
// in-region attempts. The in-region read endpoint and the next preferred
// region are backed by two distinct mock servers, and a host-routing transport
// records which endpoint each attempt actually hit. The in-region server fails
// both the initial request and the in-region retry; the request only succeeds
// once the cross-region retry routes to the second server.
func TestReadServerError_CrossRegionRoutesToDifferentEndpoint(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			// In-region endpoint: fails the initial request and the in-region retry.
			srvA, closeA := mock.NewTLSServer()
			defer closeA()
			srvA.AppendResponse(mock.WithStatusCode(statusCode))
			srvA.AppendResponse(mock.WithStatusCode(statusCode))

			// Cross-region endpoint: succeeds once the failover routes here.
			srvB, closeB := mock.NewTLSServer()
			defer closeB()
			srvB.AppendResponse(mock.WithStatusCode(200))

			endpointA, err := url.Parse(srvA.URL())
			assert.NoError(t, err)
			endpointB, err := url.Parse(srvB.URL())
			assert.NoError(t, err)
			// Sanity: the two regions must be genuinely distinct endpoints,
			// otherwise the routing assertion below would be meaningless.
			assert.NotEqual(t, endpointA.Host, endpointB.Host)

			gemServer, gemClose := mock.NewTLSServer()
			defer gemClose()
			gemServer.SetResponse(mock.WithStatusCode(200))

			internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

			// Resolve read index 0 -> endpointA (in-region) and index 1 ->
			// endpointB (next preferred region).
			lc := CreateMockLC(*endpointA, false)
			lc.locationInfo.readEndpoints = []url.URL{*endpointA, *endpointB}

			gem := &globalEndpointManager{
				clientEndpoint:      gemServer.URL(),
				pipeline:            internalPipeline,
				preferredLocations:  []string{"East US", "Central US"},
				locationCache:       lc,
				refreshTimeInterval: defaultExpirationTime,
				lastUpdateTime:      time.Time{},
			}

			retryPolicy := &clientRetryPolicy{gem: gem}
			verifier := clientRetryPolicyVerifier{}

			transport := &hostRoutingTransport{
				routes: map[string]policy.Transporter{
					endpointA.Host: srvA,
					endpointB.Host: srvB,
				},
			}

			internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: transport})

			client := &Client{endpoint: srvA.URL(), endpointUrl: endpointA, internal: internalClient, gem: gem}
			db, _ := client.NewDatabase("database_id")
			container, _ := db.NewContainer("container_id")

			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			assert.NoError(t, err)

			// Counter behavior: one in-region retry then one cross-region retry.
			assert.Equal(t, 2, verifier.requests[0].retryContext.serverErrorRetryCount)
			assert.Equal(t, 1, verifier.requests[0].retryContext.retryCount)
			assert.Equal(t, 1, verifier.requests[0].retryContext.preferredLocationIndex)

			// Routing behavior: the first two attempts hit the in-region
			// endpoint and only the third (cross-region) attempt hit the
			// second, distinct endpoint. This proves the failover changes the
			// request target rather than only mutating counters.
			assert.Equal(t, []string{endpointA.Host, endpointA.Host, endpointB.Host}, transport.seenHosts)
		})
	}
}

func TestReadServerError(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
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
				locationCache:       multiReadEndpointLC(*defaultEndpoint, false),
				refreshTimeInterval: defaultExpirationTime,
				lastUpdateTime:      time.Time{},
			}

			retryPolicy := &clientRetryPolicy{gem: gem}
			verifier := clientRetryPolicyVerifier{}

			internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

			client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
			db, _ := client.NewDatabase("database_id")
			container, _ := db.NewContainer("container_id")

			// Setting up responses for in-region retry + cross-region retry then succeeding.
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(200))
			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			// Request should retry in-region then cross-region and succeed on the third attempt.
			assert.NoError(t, err)
			assert.True(t, verifier.requests[0].retryContext.serverErrorRetryCount == 2)
			// retryCount advances only on the cross-region retry, not the in-region one.
			assert.True(t, verifier.requests[0].retryContext.retryCount == 1)
			assert.True(t, verifier.requests[0].retryContext.preferredLocationIndex == 1)

			// Setting up responses for both retries failing.
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			// Request should retry once in-region and once cross-region and then fail.
			assert.Error(t, err)
			assert.True(t, verifier.requests[1].retryContext.serverErrorRetryCount == 2)
			assert.True(t, verifier.requests[1].retryContext.retryCount == 1)

			// Without preferred locations, only the in-region retry should occur.
			gem.preferredLocations = []string{}
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			assert.Error(t, err)
			assert.True(t, verifier.requests[2].retryContext.serverErrorRetryCount == 1)
			assert.True(t, verifier.requests[2].retryContext.retryCount == 0)
		})
	}
}

func TestWriteServerErrorNotRetried(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
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
				locationCache:       CreateMockLC(*defaultEndpoint, true),
				refreshTimeInterval: defaultExpirationTime,
				lastUpdateTime:      time.Time{},
			}

			retryPolicy := &clientRetryPolicy{gem: gem}
			verifier := clientRetryPolicyVerifier{}

			internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

			client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
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

			// Even with preferred locations and multi-master configured, writes must not retry 5xx.
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(200))
			_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("1"), marshalled, nil)
			assert.Error(t, err)
			assert.True(t, verifier.requests[0].retryContext.retryCount == 0)
			assert.True(t, verifier.requests[0].retryContext.serverErrorRetryCount == 0)
		})
	}
}

func TestReadServerError_InRegionRetrySucceeds(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
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
			verifier := clientRetryPolicyVerifier{}

			internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

			client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
			db, _ := client.NewDatabase("database_id")
			container, _ := db.NewContainer("container_id")

			// In-region retry should be enough; no cross-region failover required.
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(200))
			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			assert.NoError(t, err)
			assert.True(t, verifier.requests[0].retryContext.serverErrorRetryCount == 1)
			assert.True(t, verifier.requests[0].retryContext.retryCount == 0)
			assert.True(t, verifier.requests[0].retryContext.preferredLocationIndex == 0)
		})
	}
}

func TestReadServerError_CrossRegionRetriesDisabled(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			srv, closeFunc := mock.NewTLSServer()
			defer closeFunc()

			defaultEndpoint, err := url.Parse(srv.URL())
			assert.NoError(t, err)

			gemServer, gemClose := mock.NewTLSServer()
			defer gemClose()
			gemServer.SetResponse(mock.WithStatusCode(200))

			internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

			lc := CreateMockLC(*defaultEndpoint, false)
			lc.enableCrossRegionRetries = false

			gem := &globalEndpointManager{
				clientEndpoint:      gemServer.URL(),
				pipeline:            internalPipeline,
				preferredLocations:  []string{"East US", "Central US"},
				locationCache:       lc,
				refreshTimeInterval: defaultExpirationTime,
				lastUpdateTime:      time.Time{},
			}

			retryPolicy := &clientRetryPolicy{gem: gem}
			verifier := clientRetryPolicyVerifier{}

			internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

			client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
			db, _ := client.NewDatabase("database_id")
			container, _ := db.NewContainer("container_id")

			// With cross-region disabled only the in-region retry should be attempted.
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(200))
			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			assert.Error(t, err)
			assert.True(t, verifier.requests[0].retryContext.serverErrorRetryCount == 1)
			assert.True(t, verifier.requests[0].retryContext.retryCount == 0)
			assert.True(t, verifier.requests[0].retryContext.preferredLocationIndex == 0)
		})
	}
}

func TestReadServerError_ErrorIsResponseError(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
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
				locationCache:       multiReadEndpointLC(*defaultEndpoint, false),
				refreshTimeInterval: defaultExpirationTime,
				lastUpdateTime:      time.Time{},
			}

			retryPolicy := &clientRetryPolicy{gem: gem}
			verifier := clientRetryPolicyVerifier{}

			internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

			client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
			db, _ := client.NewDatabase("database_id")
			container, _ := db.NewContainer("container_id")

			// Exhaust both retries and verify the surfaced error preserves the response.
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			srv.AppendResponse(mock.WithStatusCode(statusCode))
			_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
			assert.Error(t, err)
			var responseErr *azcore.ResponseError
			assert.True(t, errors.As(err, &responseErr))
			assert.Equal(t, statusCode, responseErr.StatusCode)
		})
	}
}

func TestReadServerError_NonRetriable5xxNotRetried(t *testing.T) {
	// 501 Not Implemented is a 5xx status that should NOT be retried.
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
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	srv.AppendResponse(mock.WithStatusCode(http.StatusNotImplemented))
	srv.AppendResponse(mock.WithStatusCode(200))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	assert.Error(t, err)
	assert.True(t, verifier.requests[0].retryContext.serverErrorRetryCount == 0)
	assert.True(t, verifier.requests[0].retryContext.retryCount == 0)
}

func TestReadServerError_MixedWith503(t *testing.T) {
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
		locationCache:       multiReadEndpointLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	// Sequence: 500 (in-region retry) -> 503 (cross-region via preferredLocationIndex) ->
	// 500 (cross-region 5xx retry; consumes the last preferred location) -> 200.
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusServiceUnavailable))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(200))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	assert.NoError(t, err)
	assert.True(t, verifier.requests[0].retryContext.serverErrorRetryCount == 2)
	// 503 advances preferredLocationIndex by 1 and retryCount by 1;
	// the cross-region 5xx retry advances both by 1 again.
	assert.True(t, verifier.requests[0].retryContext.preferredLocationIndex == 2)
	assert.True(t, verifier.requests[0].retryContext.retryCount == 2)
}

// TestReadServerError_SingleReadEndpoint verifies that when the location cache has
// only one resolved read endpoint, the cross-region retry is skipped because failing
// over would just hit the same endpoint as the in-region retry. This covers
// single-region accounts and the case where preferred locations resolve to only one
// available read endpoint.
func TestReadServerError_SingleReadEndpoint(t *testing.T) {
	srv, closeFunc := mock.NewTLSServer()
	defer closeFunc()

	defaultEndpoint, err := url.Parse(srv.URL())
	assert.NoError(t, err)

	gemServer, gemClose := mock.NewTLSServer()
	defer gemClose()
	gemServer.SetResponse(mock.WithStatusCode(200))

	internalPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: gemServer})

	lc := CreateMockLC(*defaultEndpoint, false)
	// Simulate a single-region account: only one resolved read endpoint, even though
	// the caller provided multiple preferred locations.
	lc.locationInfo.readEndpoints = []url.URL{*defaultEndpoint}

	gem := &globalEndpointManager{
		clientEndpoint:      gemServer.URL(),
		pipeline:            internalPipeline,
		preferredLocations:  []string{"East US", "Central US"},
		locationCache:       lc,
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Time{},
	}

	retryPolicy := &clientRetryPolicy{gem: gem}
	verifier := clientRetryPolicyVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerRetry: []policy.Policy{&verifier, retryPolicy}}, &policy.ClientOptions{Transport: srv})

	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	db, _ := client.NewDatabase("database_id")
	container, _ := db.NewContainer("container_id")

	// Two 5xx responses queued: the in-region retry consumes the second one and
	// fails. The cross-region retry must NOT fire because readEndpoints has length 1.
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	_, err = container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	assert.Error(t, err)
	assert.Equal(t, 1, verifier.requests[0].retryContext.serverErrorRetryCount)
	assert.Equal(t, 0, verifier.requests[0].retryContext.retryCount)
	assert.Equal(t, 0, verifier.requests[0].retryContext.preferredLocationIndex)
}
