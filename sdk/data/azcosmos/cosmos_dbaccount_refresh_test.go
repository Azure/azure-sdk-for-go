// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Regression tests for excess GetDatabaseAccount calls observed with
// preferred regions configured.
//
// These tests cover:
//   F1: failed GEM Update is throttled to refreshTimeInterval (lastAttemptTime)
//   F2: concurrent Update callers are coalesced into a single HTTP call
//   F3: write-retry on 403/WriteForbidden force-refreshes the GEM per attempt
//   F4: a failed initial bootstrap Update is retried on the next request
//   F5: locationCache.readEndpoints does not deadlock on the stale+unavailable path
//   Soak: under sustained mixed load, total GEM calls respect refreshTimeInterval

package azcosmos

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

// countingTransport counts how many requests reach the test endpoint and
// optionally returns a canned status / error so we can simulate
// GetDatabaseAccount behaviour. body, when set, is returned as the response
// body.
type countingTransport struct {
	count    atomic.Int64
	status   int
	body     []byte
	respErr  error
	delay    time.Duration
	respFunc func() (int, []byte) // when non-nil, overrides status/body per call
}

func (c *countingTransport) Do(req *http.Request) (*http.Response, error) {
	c.count.Add(1)
	if c.delay > 0 {
		time.Sleep(c.delay)
	}
	if c.respErr != nil {
		return nil, c.respErr
	}
	status, body := c.status, c.body
	if c.respFunc != nil {
		status, body = c.respFunc()
	}
	resp := &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       jsonBody(body),
		Request:    req,
	}
	return resp, nil
}

func jsonBody(b []byte) io.ReadCloser {
	return &jsonReadCloser{b: b}
}

type jsonReadCloser struct {
	b []byte
	i int
}

func (r *jsonReadCloser) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
func (r *jsonReadCloser) Close() error { return nil }

func newGEMWithTransport(t *testing.T, preferred []string, transport policy.Transporter, refresh time.Duration) *globalEndpointManager {
	t.Helper()
	pl := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: transport})
	gem, err := newGlobalEndpointManager("https://fake.documents.azure.com:443/", pl, preferred, refresh, true)
	require.NoError(t, err)
	return gem
}

// ----------------------------------------------------------------------------
// F1: a failed Update must throttle subsequent refresh attempts to
// refreshTimeInterval. The bug had every subsequent caller re-issuing
// GetDatabaseAccount because lastUpdateTime was only set on success.
// We still surface the cached error to callers (see F3b) but the HTTP call
// must not repeat within the throttle window.
// ----------------------------------------------------------------------------
func TestFix1_FailedUpdateIsThrottled(t *testing.T) {
	transport := &countingTransport{status: http.StatusBadRequest}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// First Update: lastAttemptTime is zero -> shouldRefresh()==true -> HTTP call -> fails.
	err := gem.Update(context.Background(), false)
	require.Error(t, err, "first attempt against failing endpoint must surface the error")
	require.Equal(t, int64(1), transport.count.Load())

	// Next 50 Update calls within the refresh interval return the cached
	// error -- they do NOT re-issue the HTTP call.
	for i := 0; i < 50; i++ {
		err := gem.Update(context.Background(), false)
		require.Error(t, err, "throttled Update must still surface the cached error so callers know the GEM is not populated")
	}
	require.Equal(t, int64(1), transport.count.Load(),
		"a failed Update must be throttled exactly like a successful one")
}

// ----------------------------------------------------------------------------
// F2: concurrent Update callers are coalesced into a single in-flight HTTP
// call via the single-in-flight pattern in gem.Update.
// ----------------------------------------------------------------------------
func TestFix2_ConcurrentUpdateCallersCoalesce(t *testing.T) {
	// Slow transport so the first refresh is still in flight while the rest
	// of the goroutines arrive.
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 100 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	const concurrency = 200
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	errs := make([]error, concurrency)
	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			defer wg.Done()
			errs[idx] = gem.Update(context.Background(), false)
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		require.NoErrorf(t, err, "goroutine %d failed", i)
	}
	require.Equal(t, int64(1), transport.count.Load(),
		"concurrent Update callers must coalesce into a single HTTP call")
}

// ----------------------------------------------------------------------------
// F3: write-retry on 403/WriteForbidden force-refreshes the GEM on every
// retry attempt so the client picks up topology changes immediately.
// ----------------------------------------------------------------------------
func TestFix3_WriteRetryForceRefreshesGEM(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)

	westRegion := accountRegion{Name: "West US", Endpoint: defaultEndpoint.String()}
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{westRegion},
		WriteRegions: []accountRegion{westRegion},
	})
	transport := &countingTransport{status: http.StatusOK, body: body}
	gemPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0",
		azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: transport})

	gem := &globalEndpointManager{
		clientEndpoint:      defaultEndpoint.String(),
		pipeline:            gemPipeline,
		preferredLocations:  []string{"West US"},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Now(),
	}
	retry := &clientRetryPolicy{gem: gem}

	req, err := azruntime.NewRequest(context.Background(), http.MethodPost, defaultEndpoint.String())
	require.NoError(t, err)

	const writeRetries = 5
	rc := &retryContext{}
	for i := 0; i < writeRetries; i++ {
		shouldRetry, err := retry.attemptRetryOnEndpointFailure(req, true, rc)
		require.NoError(t, err)
		require.True(t, shouldRetry)
		rc.retryCount++
	}
	require.Equal(t, int64(writeRetries), transport.count.Load(),
		"write retries on 403/WriteForbidden must force-refresh the GEM on every attempt")
}

// ----------------------------------------------------------------------------
// F1b: an invalidate() that fires while a refresh is in flight must NOT be
// lost. Before the fix, the in-flight leader's "set lastUpdateTime" would
// overwrite the invalidation timestamps, causing the next caller to observe
// the refresh as completed and skip the refresh that was demanded by the
// invalidation. The leader now snapshots invalidationGen and refuses to
// commit timestamps if it changed during the flight.
// ----------------------------------------------------------------------------
func TestFix1b_InvalidateDuringInflightRefreshIsHonored(t *testing.T) {
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 150 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// Kick off a refresh; while it is in flight, invalidate.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = gem.Update(context.Background(), false)
	}()
	// Wait deterministically for the leader to claim the in-flight slot.
	// A naive time.Sleep is racy on loaded CI hosts; polling hasInflight()
	// confirms the leader has actually entered refreshOnce.
	require.Eventually(t, gem.hasInflight, time.Second, 2*time.Millisecond,
		"leader must claim the in-flight slot within 1s")
	gem.invalidate()
	wg.Wait()

	// First refresh returned (count=1) but invalidation should have
	// prevented timestamps from being committed. A subsequent Update must
	// therefore issue a fresh HTTP call (count=2).
	require.Equal(t, int64(1), transport.count.Load())
	require.NoError(t, gem.Update(context.Background(), false))
	require.Equal(t, int64(2), transport.count.Load(),
		"invalidation during an in-flight refresh must force the next Update to actually refresh")
}

// ----------------------------------------------------------------------------
// F3c: concurrent MarkEndpointUnavailable* calls for the same endpoint may
// each observe wasUnavailable==false (the check is not atomic with the
// mark). Each one may call invalidate(). The single-in-flight pattern in
// Update bounds the resulting HTTP calls so the user-visible blast radius
// is at most one extra refresh per concurrent burst -- not one per marker.
// This test documents and bounds that behaviour.
// ----------------------------------------------------------------------------
func TestFix3c_ConcurrentSameEndpointMarksAreBounded(t *testing.T) {
	defaultEndpoint, _ := url.Parse("https://fake.documents.azure.com:443/")
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: defaultEndpoint.String()}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: defaultEndpoint.String()}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 50 * time.Millisecond}
	gemPipeline := azruntime.NewPipeline("azcosmosgemtest", "v1.0.0",
		azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: transport})
	gem := &globalEndpointManager{
		clientEndpoint:      defaultEndpoint.String(),
		pipeline:            gemPipeline,
		preferredLocations:  []string{"West US"},
		locationCache:       CreateMockLC(*defaultEndpoint, false),
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Now(),
	}

	const concurrency = 50
	wg := sync.WaitGroup{}
	wg.Add(concurrency * 2)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_ = gem.MarkEndpointUnavailableForWrite(*defaultEndpoint)
		}()
		go func() {
			defer wg.Done()
			_ = gem.Update(context.Background(), false)
		}()
	}
	wg.Wait()
	// Trigger one more Update after the burst settles. This is required for
	// scheduler determinism: under some timings every concurrent Update can
	// complete (observing the still-fresh lastUpdateTime) before any marker
	// fires invalidate(), leaving the post-invalidation refresh un-triggered.
	// This explicit post-burst Update guarantees the invalidation is acted on
	// exactly once -- the singleflight then collapses any in-flight leader
	// and the new caller to a single HTTP round-trip.
	_ = gem.Update(context.Background(), false)
	// Give any spawned refresh time to drain.
	time.Sleep(200 * time.Millisecond)

	// With the atomic check-and-mark in markEndpointUnavailable, only the
	// FIRST goroutine to win the mapMutex Lock observes wasAlreadyUnavailable
	// == false and triggers invalidate(). All other markers see true and
	// skip the invalidation. Combined with the in-flight singleflight in
	// gem.Update, the upper bound is provably 1 HTTP call -- the leader's
	// refresh handles the single invalidation, and waiters share its result.
	// If invalidate were to land AFTER the leader started its flight
	// (impossible here because the test fixture starts with a non-stale
	// lastUpdateTime and the marker is the trigger for the first refresh),
	// invalidationGen would advance mid-flight and a second leader would
	// fire -- so 2 is the theoretical maximum across all timings. We
	// assert the tight bound to guard against regressions in atomicity.
	calls := transport.count.Load()
	require.Equal(t, int64(1), calls,
		"concurrent same-endpoint marks must collapse to exactly 1 GEM call (got %d for concurrency=%d)", calls, concurrency)
}

// to every request while the GEM has never been successfully populated, not
// just the very first request. The cached lastUpdateErr is returned from
// throttled Update calls when !populated.
// ----------------------------------------------------------------------------
func TestFix3b_BootstrapFailureIsSurfacedOnEveryRequestUntilThrottleExpires(t *testing.T) {
	transport := &countingTransport{status: http.StatusBadRequest}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)
	pol := &globalEndpointManagerPolicy{gem: gem}

	downstream := &countingTransport{status: http.StatusOK}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{pol}},
		&policy.ClientOptions{Transport: downstream})

	for i := 0; i < 5; i++ {
		r, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
		r.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
		_, err := pl.Do(r)
		require.Error(t, err, "request %d must surface the cached bootstrap error", i)
	}
	require.Equal(t, int64(1), transport.count.Load(),
		"only one actual HTTP call must be made; the rest return the cached error")
}

// resettableOnce only latches on success, so the next caller retries the
// bootstrap. Combined with F1, the retries are throttled to one per
// refreshTimeInterval -- they don't fan out per request.
// ----------------------------------------------------------------------------
func TestFix4_InitialBootstrapFailureIsRetriedAndThrottled(t *testing.T) {
	transport := &countingTransport{status: http.StatusBadRequest, delay: 5 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)
	pol := &globalEndpointManagerPolicy{gem: gem}

	downstream := &countingTransport{status: http.StatusOK}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{pol}},
		&policy.ClientOptions{Transport: downstream})

	// First request: synchronous bootstrap fires the GEM HTTP call, which
	// fails. The error must surface to the caller.
	r1, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
	r1.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
	_, err := pl.Do(r1)
	require.Error(t, err, "the failed bootstrap error must surface to the caller")

	first := transport.count.Load()
	require.Equal(t, int64(1), first)

	// Next 20 sequential requests must NOT each issue a fresh
	// GetDatabaseAccount call. Bootstrap is retried on the next request but
	// throttled by lastAttemptTime; subsequent requests inside the throttle
	// window see ShouldRefresh()==false and skip the async refresh too.
	// Capture the error from each follow-up to confirm the bootstrap path
	// is actually entered (returning the cached error) -- this prevents a
	// regression where populated() accidentally returned true and the
	// bootstrap was silently skipped.
	const followUp = 20
	for i := 0; i < followUp; i++ {
		r, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
		r.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
		_, err := pl.Do(r)
		require.Error(t, err,
			"follow-up request %d must surface the cached bootstrap error -- if it doesn't, populated() is incorrectly returning true and we're silently routing to an unpopulated GEM", i)
	}
	// Drain any goroutines that the policy may have spawned.
	time.Sleep(200 * time.Millisecond)

	total := transport.count.Load()
	require.Equal(t, int64(1), total,
		"failed bootstrap must be throttled, not retried on every request (got %d HTTP calls)", total)
}

// ----------------------------------------------------------------------------
// F5: locationCache.readEndpoints / writeEndpoints no longer self-deadlock on
// the stale+unavailable refresh path. Before the fix, this hung forever
// because RLock could not be upgraded to Lock inside refreshStaleEndpoints.
// ----------------------------------------------------------------------------
func TestFix5_ReadEndpointsDoesNotDeadlock(t *testing.T) {
	defaultEndpoint, _ := url.Parse("https://fake.documents.azure.com:443/")
	lc := CreateMockLC(*defaultEndpoint, false)
	lc.locationUnavailabilityInfoMap[*defaultEndpoint] = locationUnavailabilityInfo{
		lastCheckTime:  time.Now(),
		unavailableOps: read,
	}
	lc.lastUpdateTime = time.Now().Add(-10 * time.Minute)

	done := make(chan error, 1)
	go func() {
		_, err := lc.readEndpoints()
		done <- err
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("readEndpoints deadlocked on the stale+unavailable path -- F5 regression")
	}

	done2 := make(chan error, 1)
	go func() {
		_, err := lc.writeEndpoints()
		done2 <- err
	}()
	select {
	case err := <-done2:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("writeEndpoints deadlocked on the stale+unavailable path -- F5 regression")
	}
}

// ----------------------------------------------------------------------------
// Soak test: a high-concurrency burst against a healthy GEM with the default
// 5-min refresh interval should issue exactly one GetDatabaseAccount call.
// ----------------------------------------------------------------------------
func TestRegression_HealthyHighConcurrencyStaysAtOneGEMCall(t *testing.T) {
	if testing.Short() {
		t.Skip("500-goroutine soak; skipped under -short")
	}
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 20 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)
	pol := &globalEndpointManagerPolicy{gem: gem}

	downstream := &countingTransport{status: http.StatusOK}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{pol}},
		&policy.ClientOptions{Transport: downstream})

	const concurrency = 500
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			r, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
			r.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
			_, _ = pl.Do(r)
		}()
	}
	wg.Wait()
	time.Sleep(300 * time.Millisecond)

	calls := transport.count.Load()
	require.Equal(t, int64(1), calls,
		"500 concurrent requests on a healthy client must produce exactly one GetDatabaseAccount call (got %d)", calls)
}

// Soak test variant: the same burst against a FAILING GEM must also produce
// exactly one HTTP call, demonstrating that F1 + F2 together close the
// failure-storm path.
func TestRegression_FailingGEMHighConcurrencyStaysAtOneGEMCall(t *testing.T) {
	if testing.Short() {
		t.Skip("500-goroutine soak; skipped under -short")
	}
	transport := &countingTransport{status: http.StatusBadRequest, delay: 20 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)
	pol := &globalEndpointManagerPolicy{gem: gem}

	downstream := &countingTransport{status: http.StatusOK}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{pol}},
		&policy.ClientOptions{Transport: downstream})

	const concurrency = 500
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			r, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
			r.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
			_, _ = pl.Do(r)
		}()
	}
	wg.Wait()
	time.Sleep(300 * time.Millisecond)

	calls := transport.count.Load()
	require.Equal(t, int64(1), calls,
		"500 concurrent requests against a failing GEM must still produce only one GetDatabaseAccount call (got %d)", calls)
}

// ----------------------------------------------------------------------------
// Default-endpoint elimination.
// After the GEM is populated, data-plane requests must never resolve to the
// customer-supplied endpoint -- with the single exception of the degenerate
// "zero write regions" case for write requests.
// ----------------------------------------------------------------------------

// makeGEMWithRegions builds a populated GEM whose default endpoint host
// differs from every account-region host, so a routing fallback to the
// default endpoint is observable.
func makeGEMWithRegions(t *testing.T, isMultiMaster bool, preferred []string, enableCrossRegion bool) *globalEndpointManager {
	t.Helper()
	defaultEndpoint, _ := url.Parse("https://customer-endpoint.documents.azure.com:443/")
	lc := &locationCache{
		defaultEndpoint:                   *defaultEndpoint,
		locationUnavailabilityInfoMap:     map[url.URL]locationUnavailabilityInfo{},
		unavailableLocationExpirationTime: defaultExpirationTime,
		enableCrossRegionRetries:          enableCrossRegion,
		enableMultipleWriteLocations:      isMultiMaster,
	}
	writeRegions := []accountRegion{
		{Name: "East US", Endpoint: "https://east-us.documents.azure.com:443/"},
	}
	if isMultiMaster {
		writeRegions = append(writeRegions, accountRegion{Name: "Central US", Endpoint: "https://central-us.documents.azure.com:443/"})
	}
	readRegions := []accountRegion{
		{Name: "East US", Endpoint: "https://east-us.documents.azure.com:443/"},
		{Name: "Central US", Endpoint: "https://central-us.documents.azure.com:443/"},
		{Name: "West US", Endpoint: "https://west-us.documents.azure.com:443/"},
	}
	require.NoError(t, lc.update(writeRegions, readRegions, preferred, &isMultiMaster))
	gem := &globalEndpointManager{
		clientEndpoint:      defaultEndpoint.String(),
		preferredLocations:  preferred,
		locationCache:       lc,
		refreshTimeInterval: defaultExpirationTime,
		lastUpdateTime:      time.Now(),
	}
	return gem
}

func TestDefaultEndpointElim_DataPlaneNeverHitsDefaultWhenPopulated(t *testing.T) {
	defaultHost := "customer-endpoint.documents.azure.com:443"

	cases := []struct {
		name      string
		multi     bool
		preferred []string
	}{
		{"singleMaster_withPreferred", false, []string{"West US", "East US"}},
		{"singleMaster_noPreferred", false, []string{}},
		{"multiMaster_withPreferred", true, []string{"Central US", "East US"}},
		{"multiMaster_noPreferred", true, []string{}},
	}
	resourceTypes := []resourceType{resourceTypeDocument, resourceTypeCollection, resourceTypeDatabase}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gem := makeGEMWithRegions(t, tc.multi, tc.preferred, true)
			for _, rt := range resourceTypes {
				for _, isWrite := range []bool{false, true} {
					for _, useWrite := range []bool{false, true} {
						for idx := 0; idx < 25; idx++ {
							ep := gem.ResolveServiceEndpoint(idx, rt, isWrite, useWrite)
							require.NotEqual(t, defaultHost, ep.Host,
								"resourceType=%v isWrite=%v useWrite=%v idx=%d resolved to default endpoint", rt, isWrite, useWrite, idx)
						}
					}
				}
			}
		})
	}
}

func TestDefaultEndpointElim_SessionRetryOnSingleMasterRoutesToRegionalWrite(t *testing.T) {
	defaultHost := "customer-endpoint.documents.azure.com:443"
	gem := makeGEMWithRegions(t, false /*single master*/, []string{"West US"}, true)

	// Simulate the session-not-available retry: useWriteEndpoint=true on a
	// single-master read. The resolved endpoint must be the regional write
	// region, not the customer endpoint.
	ep := gem.ResolveServiceEndpoint(1, resourceTypeDocument, false /*isWrite*/, true /*useWrite*/)
	require.NotEqual(t, defaultHost, ep.Host)
	require.Contains(t, ep.Host, "east-us")
}

func TestDefaultEndpointElim_UnmatchedPreferredRegionRoutesToRegional(t *testing.T) {
	defaultHost := "customer-endpoint.documents.azure.com:443"
	// Customer set preferred regions that the account does NOT advertise.
	gem := makeGEMWithRegions(t, false, []string{"South Africa North", "France Central"}, true)

	for idx := 0; idx < 10; idx++ {
		readEP := gem.ResolveServiceEndpoint(idx, resourceTypeDocument, false, false)
		require.NotEqual(t, defaultHost, readEP.Host,
			"read with unmatched preferred region resolved to default endpoint at idx=%d", idx)
		writeEP := gem.ResolveServiceEndpoint(idx, resourceTypeDocument, true, false)
		require.NotEqual(t, defaultHost, writeEP.Host,
			"write with unmatched preferred region resolved to default endpoint at idx=%d", idx)
	}
}

func TestDefaultEndpointElim_ZeroWriteRegionsRetainsDefaultFallback(t *testing.T) {
	// Degenerate account metadata: no write regions advertised. Per the
	// project decision, this is the ONLY remaining data-plane code path
	// where a write request routes to the customer-supplied endpoint.
	defaultEndpoint, _ := url.Parse("https://customer-endpoint.documents.azure.com:443/")
	lc := &locationCache{
		defaultEndpoint:                   *defaultEndpoint,
		locationUnavailabilityInfoMap:     map[url.URL]locationUnavailabilityInfo{},
		unavailableLocationExpirationTime: defaultExpirationTime,
		enableCrossRegionRetries:          true,
		enableMultipleWriteLocations:      false,
	}
	readRegions := []accountRegion{
		{Name: "East US", Endpoint: "https://east-us.documents.azure.com:443/"},
	}
	noWrites := []accountRegion{} // explicit zero
	multi := false
	require.NoError(t, lc.update(noWrites, readRegions, []string{}, &multi))
	gem := &globalEndpointManager{
		clientEndpoint: defaultEndpoint.String(), locationCache: lc,
		refreshTimeInterval: defaultExpirationTime, lastUpdateTime: time.Now(),
	}

	ep := gem.ResolveServiceEndpoint(0, resourceTypeDocument, true /*isWrite*/, false)
	require.Equal(t, defaultEndpoint.Host, ep.Host,
		"zero-write-regions write must still route to the customer endpoint (documented degenerate case)")
}

func TestDefaultEndpointElim_CrossRegionRetriesDisabledUsesDefaultEndpoint(t *testing.T) {
	defaultHost := "customer-endpoint.documents.azure.com:443"
	gem := makeGEMWithRegions(t, false, []string{"East US"}, false /*enableCrossRegion=false*/)

	// With cross-region retries disabled, single-master writes route to the
	// default (customer-supplied) endpoint -- preserving the long-standing
	// behavior of this flag.
	for idx := 0; idx < 10; idx++ {
		ep := gem.ResolveServiceEndpoint(idx, resourceTypeDocument, true, false)
		require.Equal(t, defaultHost, ep.Host,
			"cross-region-retries=false must route writes to the default endpoint at idx=%d", idx)
	}
}

// TestDefaultEndpointElim_ZeroWriteRegionsReadGoesToReadRegion verifies the
// corollary of TestDefaultEndpointElim_ZeroWriteRegionsRetainsDefaultFallback:
// when the account advertises zero write regions, READS must still resolve to
// an advertised read region rather than the customer endpoint. Only the
// write request in this scenario is allowed to fall through to default.
func TestDefaultEndpointElim_ZeroWriteRegionsReadGoesToReadRegion(t *testing.T) {
	defaultEndpoint, _ := url.Parse("https://customer-endpoint.documents.azure.com:443/")
	lc := &locationCache{
		defaultEndpoint:                   *defaultEndpoint,
		locationUnavailabilityInfoMap:     map[url.URL]locationUnavailabilityInfo{},
		unavailableLocationExpirationTime: defaultExpirationTime,
		enableCrossRegionRetries:          true,
		enableMultipleWriteLocations:      false,
	}
	readRegions := []accountRegion{
		{Name: "East US", Endpoint: "https://east-us.documents.azure.com:443/"},
		{Name: "Central US", Endpoint: "https://central-us.documents.azure.com:443/"},
	}
	multi := false
	require.NoError(t, lc.update([]accountRegion{}, readRegions, []string{}, &multi))
	gem := &globalEndpointManager{
		clientEndpoint: defaultEndpoint.String(), locationCache: lc,
		refreshTimeInterval: defaultExpirationTime, lastUpdateTime: time.Now(),
	}

	for idx := 0; idx < 10; idx++ {
		ep := gem.ResolveServiceEndpoint(idx, resourceTypeDocument, false /*isWrite*/, false)
		require.NotEqual(t, defaultEndpoint.Host, ep.Host,
			"reads must route to a read region even when there are zero write regions; got default at idx=%d", idx)
	}
}

// ----------------------------------------------------------------------------
// F7: a transient post-invalidation refresh failure must NOT stall the data
// plane for the full refreshTimeInterval. Once the GEM has ever been
// populated, requests should keep routing through the cached topology even
// while a refresh attempt is throttled after a failure. This guards the
// regression flagged by the deep reviewer where invalidate() zeroed
// lastUpdateTime, which in turn made populated() return false and caused
// the policy's bootstrap path to surface the cached error on every request.
// ----------------------------------------------------------------------------
func TestFix7_InvalidateThenRefreshFailureDoesNotStallDataPlane(t *testing.T) {
	// Programmable transport: succeed once (bootstrap), then fail.
	responses := []int{http.StatusOK, http.StatusBadRequest, http.StatusBadRequest}
	idx := atomic.Int32{}
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://west-us.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://west-us.documents.azure.com:443/"}},
	})
	transport := &countingTransport{}
	transport.respFunc = func() (int, []byte) {
		i := idx.Add(1) - 1
		if int(i) < len(responses) {
			status := responses[i]
			if status == http.StatusOK {
				return status, body
			}
			return status, nil
		}
		return http.StatusBadRequest, nil
	}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)
	pol := &globalEndpointManagerPolicy{gem: gem}

	downstream := &countingTransport{status: http.StatusOK}
	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{pol}},
		&policy.ClientOptions{Transport: downstream})

	// 1. Bootstrap succeeds.
	r1, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
	r1.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
	_, err := pl.Do(r1)
	require.NoError(t, err, "bootstrap must succeed")
	require.True(t, gem.populated(), "GEM must be populated after a successful bootstrap")

	// 2. Simulate a regional 403 -> MarkEndpointUnavailableForWrite -> invalidate.
	west, _ := url.Parse("https://west-us.documents.azure.com:443/")
	require.NoError(t, gem.MarkEndpointUnavailableForWrite(*west))

	// 3. populated() must remain true even though lastUpdateTime is zero.
	require.True(t, gem.populated(),
		"populated() must remain true after invalidate() -- otherwise a transient refresh failure stalls the data plane")

	// 4. A subsequent data-plane request must succeed (routing through the cached
	// topology) even if the post-invalidation refresh attempt fails.
	r2, _ := azruntime.NewRequest(context.Background(), http.MethodGet, "https://fake.documents.azure.com/")
	r2.SetOperationValue(pipelineRequestOptions{resourceType: resourceTypeDocument})
	_, err = pl.Do(r2)
	require.NoError(t, err,
		"data-plane request must succeed via cached topology even when post-invalidate refresh fails")
}

// ----------------------------------------------------------------------------
// F8: a waiter on an in-flight refresh must respect its own ctx deadline,
// not block for the leader's full HTTP round-trip duration. This matters
// because clientRetryPolicy's GEM-refresh calls pass req.Raw().Context()
// (a cancellable user context), so a caller-side timeout must take effect
// promptly even when the user happens to be a waiter rather than the leader.
// ----------------------------------------------------------------------------
func TestFix8_UpdateWaiterRespectsContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("blocks on a 2s leader by design; skipped under -short")
	}
	// Slow leader: holds the in-flight slot for 2 seconds.
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 2 * time.Second}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// Leader starts a refresh in the background.
	leaderDone := make(chan struct{})
	go func() {
		defer close(leaderDone)
		_ = gem.Update(context.Background(), false)
	}()

	// Wait until the leader is actually in-flight via the hasInflight test
	// helper rather than peeking at internal mutex state -- keeps the test
	// resilient to changes in the GEM's internal synchronization primitives.
	require.Eventually(t, gem.hasInflight, 500*time.Millisecond, 5*time.Millisecond,
		"leader must claim the in-flight slot within 500ms")

	// Waiter has a 100 ms deadline. It must return promptly with the
	// deadline error rather than blocking for the leader's 2-second HTTP
	// call to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := gem.Update(ctx, false)
	waited := time.Since(start)

	require.ErrorIs(t, err, context.DeadlineExceeded, "waiter must surface ctx error")
	require.Less(t, waited, 500*time.Millisecond,
		"waiter must respect ctx deadline -- waited %v but ctx deadline was 100ms", waited)

	// Clean up: let the leader finish so the test doesn't leak the goroutine.
	<-leaderDone
}
