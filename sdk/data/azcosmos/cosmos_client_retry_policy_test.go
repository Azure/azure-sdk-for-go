// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
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
	assert.Equal(t, 1, rc.retryCount)
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
	assert.Equal(t, 1, rc.retryCount)
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
	// 3 same-region attempts, no cross-region failover for ambiguous writes.
	assert.Equal(t, 3, rc.sameRegionRetryCount)
	assert.Equal(t, 0, rc.retryCount)
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
	assert.Equal(t, 1, rc.retryCount)
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
	assert.Equal(t, 1, rc.retryCount)
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
	assert.Equal(t, 1, rc.retryCount)
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
