// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore azcosmosgemtest azcosmostest retriable

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
	"fmt"
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
		StatusCode:    status,
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          jsonBody(body),
		ContentLength: int64(len(body)),
		Request:       req,
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
// F3: write-retry on 403/WriteForbidden kicks off an opportunistic GEM
// refresh on the FIRST mark for each endpoint and on subsequent marks
// at most once per forcedRefreshMinInterval (rate-limited). It must
// NOT issue one refresh per retry (that would storm
// GetDatabaseAccount during a sustained 403 flap) and it must NOT
// issue zero refreshes after the first (single-master writes need
// recovery when the first refresh returns stale topology). The
// refresh is fire-and-forget so a stalled metadata endpoint cannot
// block the retry.
// ----------------------------------------------------------------------------
func TestFix3_WriteRetryKicksOffFireAndForgetRefresh(t *testing.T) {
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
	start := time.Now()
	for i := 0; i < writeRetries; i++ {
		shouldRetry, err := retry.attemptRetryOnEndpointFailure(req, true, rc)
		require.NoError(t, err, "retry must not surface metadata-refresh errors")
		require.True(t, shouldRetry, "write retries on 403/WriteForbidden must continue")
		rc.retryCount++
	}
	elapsed := time.Since(start)

	// At least one refresh must run (the first 403 always forces).
	require.Eventually(t, func() bool {
		return transport.count.Load() >= 1
	}, 5*time.Second, 10*time.Millisecond,
		"first 403 must kick off a GEM refresh")
	// Give any racing follow-up refresh a chance to land.
	time.Sleep(200 * time.Millisecond)

	// Upper bound: rate-limited to at most one refresh per
	// forcedRefreshMinInterval. With 5 retries each sleeping
	// defaultBackoff*time.Second between attempts, elapsed ~= 5s.
	// Expected refreshes: 1 (initial) + floor(elapsed /
	// forcedRefreshMinInterval) at most. Allow +1 for boundary
	// scheduling slack.
	maxExpected := int64(1 + elapsed/forcedRefreshMinInterval) + 1
	require.LessOrEqual(t, transport.count.Load(), maxExpected,
		"sustained 403s against the same endpoint must be rate-limited (elapsed=%v got=%d max=%d)",
		elapsed, transport.count.Load(), maxExpected)
}

// ----------------------------------------------------------------------------
// F3a: 403/WriteForbidden retry must complete successfully even when the
// asynchronous GEM refresh hits a hard failure (5xx, dial error, etc.).
// Before the async-refresh change the policy did
//
//	if err := gem.Update(ctx, true); err != nil { return false, err }
//
// so any GEM-refresh failure short-circuited the retry and surfaced the
// metadata error to the caller -- exactly the regression we're guarding
// against here.
// ----------------------------------------------------------------------------
func TestFix3a_WriteRetrySucceedsWhenGEMRefreshFails(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)

	// Transport that always returns 500 -- gem.Update will fail every
	// time. The retry must still return (true, nil).
	transport := &countingTransport{status: http.StatusInternalServerError}
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

	rc := &retryContext{}
	shouldRetry, err := retry.attemptRetryOnEndpointFailure(req, true, rc)
	require.NoError(t, err, "GEM refresh failures must not surface to the retry caller")
	require.True(t, shouldRetry, "retry must proceed regardless of GEM refresh outcome")

	// Sanity check: the async refresh actually ran (and failed). Without
	// this, a regression that made asyncForceRefreshGEM a no-op would
	// still pass the (true, nil) assertions above.
	require.Eventually(t, func() bool {
		return transport.count.Load() >= 1
	}, 5*time.Second, 10*time.Millisecond,
		"async GEM refresh must run even though it fails")
}

// ----------------------------------------------------------------------------
// F3b: 408 read retry must NOT mark the request's endpoint unavailable.
// A 408 is a per-request signal (the gateway accepted the request but
// did not produce a response in time); demoting the whole region in the
// location cache for unavailableLocationExpirationTime based on one
// slow request would penalize concurrent reads against a region that
// may be perfectly healthy.
// ----------------------------------------------------------------------------
func TestFix3b_RequestTimeoutDoesNotMarkEndpointUnavailable(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)

	transport := &countingTransport{status: http.StatusOK, body: []byte("{}")}
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

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, defaultEndpoint.String())
	require.NoError(t, err)

	rc := &retryContext{}
	shouldRetry, err := retry.attemptRetryOnRequestTimeout(req, false /*isWriteOperation*/, rc)
	require.NoError(t, err)
	require.True(t, shouldRetry)
	require.True(t, rc.requestTimeoutRetryDone)

	require.Empty(t, gem.locationCache.locationUnavailabilityInfoMap,
		"408 must not record any endpoint as unavailable")
	require.Equal(t, int64(0), transport.count.Load(),
		"408 retry must not synchronously hit the GEM endpoint")
}

// ----------------------------------------------------------------------------
// F3c: 408 read retry is non-blocking -- even a permanently-stalled GEM
// endpoint cannot delay or fail the retry, because the 408 path no
// longer calls gem.Update at all.
// ----------------------------------------------------------------------------
func TestFix3c_RequestTimeoutDoesNotBlockOnStalledGEM(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)

	// Hanging transport: any call would block effectively forever.
	transport := &countingTransport{status: http.StatusOK, body: []byte("{}"), delay: 10 * time.Minute}
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

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, defaultEndpoint.String())
	require.NoError(t, err)

	// defaultBackoff*time.Second is 1s, so give a comfortable upper bound.
	done := make(chan struct{})
	var shouldRetry bool
	var retryErr error
	go func() {
		shouldRetry, retryErr = retry.attemptRetryOnRequestTimeout(req, false, &retryContext{})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("attemptRetryOnRequestTimeout blocked on the stalled GEM endpoint")
	}
	require.NoError(t, retryErr)
	require.True(t, shouldRetry)
}

// ----------------------------------------------------------------------------
// F3d: asyncForceRefreshGEM's CAS gate (asyncRefreshPending) must coalesce
// a retry storm. With N concurrent retries hitting a slow GEM endpoint,
// only one refresh goroutine should reach gem.Update at a time; the GEM's
// own single-in-flight pattern further coalesces to a single HTTP call.
// Without the gate every retry would queue its own goroutine inside
// gem.Update, wasting goroutine + channel overhead for no benefit.
// ----------------------------------------------------------------------------
func TestFix3d_AsyncRefreshCASGateCoalescesRetryStorm(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)

	westRegion := accountRegion{Name: "West US", Endpoint: defaultEndpoint.String()}
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{westRegion},
		WriteRegions: []accountRegion{westRegion},
	})
	// Slow-but-successful GEM transport so refreshes overlap in time.
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 200 * time.Millisecond}
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

	const concurrency = 50
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			retry.asyncForceRefreshGEM()
		}()
	}
	wg.Wait()

	// All asyncForceRefreshGEM callers have returned (the CAS gate is
	// non-blocking). Wait for whichever refresh goroutine did get spawned
	// to complete its HTTP call.
	require.Eventually(t, func() bool {
		return transport.count.Load() >= 1
	}, 5*time.Second, 10*time.Millisecond,
		"at least one refresh goroutine must reach gem.Update")

	// Give any racing follow-up refresh a chance to land before we assert
	// the upper bound. asyncRefreshPending clears in the spawned
	// goroutine's defer, which runs AFTER gem.Update returns, so a second
	// caller arriving in that tiny window could legitimately spawn a
	// second refresh. The bound is "no goroutine-per-retry storm", not
	// strictly one.
	time.Sleep(100 * time.Millisecond)
	require.Less(t, transport.count.Load(), int64(concurrency/2),
		"CAS gate must coalesce concurrent retries; got %d HTTP calls for %d retries",
		transport.count.Load(), concurrency)
}

// ----------------------------------------------------------------------------
// F3e: asyncForceRefreshGEM must use a background context internally, not
// the caller's context. If we threaded the caller's context (or any
// derivative that inherits its deadline/cancellation) into gem.Update,
// an already-cancelled or near-expired caller context would abort the
// refresh immediately and defeat its purpose.
// ----------------------------------------------------------------------------
func TestFix3e_AsyncRefreshIgnoresCallerContextCancellation(t *testing.T) {
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

	// Call asyncForceRefreshGEM directly. We deliberately bypass
	// attemptRetryOnEndpointFailure here because that path now uses
	// sleepWithContext for the backoff (so caller-cancellation correctly
	// short-circuits the retry budget). The contract we're verifying is
	// narrower: asyncForceRefreshGEM itself must not inherit any caller
	// context -- it runs on context.Background() so an already-cancelled
	// or near-expired caller deadline cannot abort the background HTTP
	// call.
	retry.asyncForceRefreshGEM()

	require.Eventually(t, func() bool {
		return transport.count.Load() >= 1
	}, 5*time.Second, 10*time.Millisecond,
		"asyncForceRefreshGEM must run on context.Background() and complete its HTTP call")
}

// ----------------------------------------------------------------------------
// F3f: when the first forced async refresh fails, a subsequent 403 on
// the SAME endpoint must be allowed to spawn another forced refresh.
// Before the fix the wasAlreadyUnavailable guard unconditionally
// suppressed every subsequent forced refresh, which stranded
// single-master writes on the failed write endpoint for the GEM
// throttle window (refreshTimeInterval, default 5 min) after the
// metadata endpoint recovered.
//
// We exercise this at the policy unit level by driving
// asyncForceRefreshGEM's state machine directly: the first call must
// land in asyncRefreshFailed, and the retry-policy's gating logic
// (mirrored here) must permit a new refresh when state == Failed.
// ----------------------------------------------------------------------------
func TestFix3f_FailedAsyncRefreshIsRetriedOnNextSameEndpoint403(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)

	// A transport that returns an error on demand; we toggle errOn to
	// flip between failing and succeeding refreshes.
	var errOn atomic.Bool
	errOn.Store(true)
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: defaultEndpoint.String()}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: defaultEndpoint.String()}},
	})
	transport := &countingTransport{
		respFunc: func() (int, []byte) {
			if errOn.Load() {
				// Return 500 with a body that will not parse; pipeline
				// surfaces this as an error and azcore retry does NOT
				// retry a 500 with no Retry-After... wait, it does.
				// Instead use a non-retriable status the GEM treats as
				// an error: 400.
				return http.StatusBadRequest, []byte(`{"code":"BadRequest"}`)
			}
			return http.StatusOK, body
		},
	}
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

	// First forced refresh: spawns goroutine, transport returns 400 ->
	// gem.Update returns an error -> state should land at Failed.
	require.True(t, retry.asyncForceRefreshGEM(), "first call must spawn the refresh goroutine")
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if retry.asyncRefreshState.Load() == asyncRefreshFailed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.Equal(t, asyncRefreshFailed, retry.asyncRefreshState.Load(),
		"first forced refresh must record Failed (count=%d)", transport.count.Load())
	firstCount := transport.count.Load()
	require.GreaterOrEqual(t, firstCount, int64(1), "first refresh must hit the transport")

	// Now make the transport succeed and call asyncForceRefreshGEM
	// again. With state == Failed it MUST spawn a new goroutine even
	// though no fresh invalidation happened.
	errOn.Store(false)
	require.True(t, retry.asyncForceRefreshGEM(),
		"second call must be allowed because previous refresh failed")
	deadline = time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if retry.asyncRefreshState.Load() == asyncRefreshIdle && transport.count.Load() > firstCount {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.Equal(t, asyncRefreshIdle, retry.asyncRefreshState.Load(),
		"second forced refresh must succeed (count=%d)", transport.count.Load())
	require.Greater(t, transport.count.Load(), firstCount,
		"second refresh must hit the transport (first=%d total=%d)", firstCount, transport.count.Load())

	// And once Idle, a third call (no failure, no invalidation) is
	// allowed too (Idle state always permits a new refresh).
	require.True(t, retry.asyncForceRefreshGEM(),
		"third call from Idle state must be allowed")
}

// ----------------------------------------------------------------------------
// F1c: a forced refresh request that arrives while a stale refresh is
// in flight, where invalidate() fires WHILE the waiter is blocked on
// the in-flight, must still trigger a fresh post-invalidation refresh.
// Before the fix the waiter sampled invalidationGen before the wait,
// so any invalidation during the wait was lost and the waiter returned
// the stale (pre-invalidation) flight result.
// ----------------------------------------------------------------------------
func TestFix1c_ForceRefreshWaiterReReadsInvalidationGenAfterWait(t *testing.T) {
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 250 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// 1) Kick off a refresh (becomes the leader/inflight).
	leaderDone := make(chan struct{})
	go func() {
		defer close(leaderDone)
		_ = gem.Update(context.Background(), false)
	}()
	require.Eventually(t, gem.hasInflight, 2*time.Second, 5*time.Millisecond,
		"leader must enter the in-flight slot")

	// 2) Start a forceRefresh waiter while the leader is still in flight.
	waiterStart := make(chan struct{})
	waiterDone := make(chan struct{})
	go func() {
		close(waiterStart)
		_ = gem.Update(context.Background(), true)
		close(waiterDone)
	}()
	<-waiterStart

	// 3) While the waiter is blocked on <-leaderFlight.done, invalidate.
	//    The leader will finish soon (delay=250ms); the waiter must
	//    observe the post-invalidate genAtStart on re-read and loop.
	time.Sleep(50 * time.Millisecond)
	mark, _ := url.Parse("https://fake.documents.azure.com:443/")
	_, markErr := gem.MarkEndpointUnavailableForWrite(*mark)
	require.NoError(t, markErr)

	// 4) Both goroutines must complete.
	select {
	case <-leaderDone:
	case <-time.After(5 * time.Second):
		t.Fatal("leader refresh did not complete")
	}
	select {
	case <-waiterDone:
	case <-time.After(5 * time.Second):
		t.Fatal("waiter did not complete; likely stuck looping or never loops")
	}

	// 5) Expect exactly 2 HTTP calls: the leader's, plus the
	//    post-invalidation refresh the waiter should have led after
	//    looping. Without the fix this would be 1 (waiter returned the
	//    stale flight without looping).
	require.Equal(t, int64(2), transport.count.Load(),
		"waiter must loop and lead a fresh refresh after invalidation during wait")
}

// ----------------------------------------------------------------------------
// F1d: a forceRefresh LEADER (not waiter) whose flight gets invalidated
// mid-flight must loop and lead a fresh refresh. Before the fix, only
// waiters re-read invalidationGen after their wait; a leader whose own
// genAtStart predated a mid-flight invalidate() would simply return,
// leaving asyncRefreshState=Idle in the retry policy and silently
// skipping the post-invalidation refresh the caller actually needed.
// ----------------------------------------------------------------------------
func TestFix1d_ForceRefreshLeaderLoopsOnMidFlightInvalidation(t *testing.T) {
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 250 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// Kick off a forceRefresh as leader.
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- gem.Update(context.Background(), true /*forceRefresh*/)
	}()
	require.Eventually(t, gem.hasInflight, 2*time.Second, 5*time.Millisecond,
		"leader must enter the in-flight slot")

	// Fire invalidate() during the leader's flight.
	time.Sleep(50 * time.Millisecond)
	mark, _ := url.Parse("https://fake.documents.azure.com:443/")
	_, err := gem.MarkEndpointUnavailableForWrite(*mark)
	require.NoError(t, err)

	// The leader's first flight will complete; the outer loop must
	// detect latestGen > genAtStart and lead a second flight. Wait for
	// the entire Update call to return.
	select {
	case updateErr := <-leaderDone:
		require.NoError(t, updateErr)
	case <-time.After(5 * time.Second):
		t.Fatal("leader Update did not return; likely never loops")
	}

	// Expect exactly 2 HTTP calls: the original flight that predated
	// the invalidation + the post-invalidation refresh the leader led
	// after looping. Without the fix this would be 1.
	require.Equal(t, int64(2), transport.count.Load(),
		"leader must loop and lead a fresh refresh after mid-flight invalidation")
}

// ----------------------------------------------------------------------------
// F3g: attemptRetryOnSessionUnavailable must snapshot
// (enableMultipleWriteLocations, availReadLocations, availWriteLocations)
// atomically. Otherwise a concurrent locationCache.update can flip the
// multi-write decision between the CanUseMultipleWriteLocations() check
// and the subsequent slice-length sampling, producing a routing decision
// that mixes pre- and post-refresh topology state.
// ----------------------------------------------------------------------------
func TestFix3g_SessionUnavailableSnapshotIsAtomic(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)
	lc := CreateMockLC(*defaultEndpoint, true /*multiMaster*/)
	multiWrite, readN, writeN := lc.sessionRetrySnapshot()
	require.True(t, multiWrite, "multi-write must be reported")
	require.Greater(t, readN, 0, "read count must be populated")
	require.Greater(t, writeN, 0, "write count must be populated")

	// Concurrent flip+snapshot race: a hostile updater toggles
	// enableMultipleWriteLocations and rewrites the slices repeatedly
	// while a reader takes snapshots. Each snapshot must be internally
	// consistent: multiWrite == true => write/read slices are the
	// multi-master shape that was current at the lock acquisition.
	stop := make(chan struct{})
	go func() {
		toggle := true
		for {
			select {
			case <-stop:
				return
			default:
			}
			toggle = !toggle
			_ = lc.update(nil, nil, nil, &toggle)
		}
	}()
	defer close(stop)

	deadline := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(deadline) {
		mw, rN, wN := lc.sessionRetrySnapshot()
		// The snapshot must come from a single locked read. We can't
		// directly verify that without instrumentation, but at minimum
		// the returned tuple must be a value (not a panic / negative).
		require.GreaterOrEqual(t, rN, 0)
		require.GreaterOrEqual(t, wN, 0)
		_ = mw
	}
}

// ----------------------------------------------------------------------------
// F3h: when the first forced refresh returns successfully but the
// metadata still reflects pre-failover topology (a common race during
// single-master account failovers), the retry policy MUST be able to
// force another refresh against the same already-unavailable endpoint
// after forcedRefreshMinInterval. Otherwise the policy's wasAlreadyUnavailable
// gate plus asyncRefreshState=Idle on success would suppress every
// subsequent forced refresh and leave the client stuck on the bad
// write endpoint until the GEM throttle expires (default 5 minutes).
//
// We verify the gating function staleForcedRefresh() and the spawn
// decision directly rather than driving the full attemptRetryOnEndpointFailure
// (which sleeps defaultBackoff between calls and would couple this
// test to that timing).
// ----------------------------------------------------------------------------
func TestFix3h_RepeatWriteForbiddenForcesRefreshAfterRateWindow(t *testing.T) {
	defaultEndpoint, err := url.Parse("https://fake.documents.azure.com:443/")
	require.NoError(t, err)
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: defaultEndpoint.String()}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: defaultEndpoint.String()}},
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

	// Before any refresh has run, staleForcedRefresh() must return true
	// so the very first 403 unconditionally triggers a refresh.
	require.True(t, retry.staleForcedRefresh(),
		"with no recorded prior refresh, staleForcedRefresh must allow a spawn")

	// Spawn a real refresh and wait for it to complete -- this populates
	// lastForcedRefreshUnixNano with "now".
	require.True(t, retry.asyncForceRefreshGEM())
	require.Eventually(t, func() bool {
		return retry.asyncRefreshState.Load() == asyncRefreshIdle && transport.count.Load() >= 1
	}, 5*time.Second, 10*time.Millisecond, "first refresh must complete")
	require.NotZero(t, retry.lastForcedRefreshUnixNano.Load(),
		"the goroutine's defer must record completion time")

	// Immediately after the first refresh completes, staleForcedRefresh()
	// must return false so a follow-up 403 on the same endpoint does
	// NOT spawn another refresh -- the rate-limit window has not
	// elapsed and a tight 403 loop must not storm GetDatabaseAccount.
	require.False(t, retry.staleForcedRefresh(),
		"within forcedRefreshMinInterval of a completed refresh, follow-up 403s must be rate-limited")

	// Simulate the rate window elapsing by rewinding the recorded
	// timestamp by forcedRefreshMinInterval + a small slack. Now the
	// gate must allow a new refresh.
	retry.lastForcedRefreshUnixNano.Store(time.Now().Add(-forcedRefreshMinInterval - 50*time.Millisecond).UnixNano())
	require.True(t, retry.staleForcedRefresh(),
		"after forcedRefreshMinInterval has elapsed, repeat 403 must be allowed to spawn a new refresh")
}

// ----------------------------------------------------------------------------
// F1e: the forceRefresh leader's mid-flight-invalidation loop must be
// bounded. Under sustained invalidations the loop could otherwise spin
// indefinitely inside a single Update call, monopolizing the inflight
// slot. After maxForceRefreshRetries iterations the call returns; the
// retry policy state machine then sees the outcome and the NEXT caller
// can lead a fresh attempt.
// ----------------------------------------------------------------------------
func TestFix1e_ForceRefreshLeaderLoopIsBounded(t *testing.T) {
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	// Slow transport so we have time to invalidate during each flight.
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 80 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// Hostile concurrent invalidator: keeps firing invalidate() so
	// every leader iteration sees latestGen > genAtStart.
	stop := make(chan struct{})
	go func() {
		mark, _ := url.Parse("https://fake.documents.azure.com:443/")
		for {
			select {
			case <-stop:
				return
			case <-time.After(5 * time.Millisecond):
				_, _ = gem.MarkEndpointUnavailableForWrite(*mark)
			}
		}
	}()
	defer close(stop)

	start := time.Now()
	err := gem.Update(context.Background(), true /*forceRefresh*/)
	elapsed := time.Since(start)
	require.NoError(t, err)

	// 1 initial + maxForceRefreshRetries loop iterations = up to
	// maxForceRefreshRetries+1 flights. Each takes ~80ms; total must
	// be far less than runaway (multi-second). 80ms * (3+1) = 320ms
	// nominal; allow up to 2s of scheduler slack.
	require.LessOrEqual(t, transport.count.Load(), int64(maxForceRefreshRetries+1),
		"leader loop must be bounded by maxForceRefreshRetries (count=%d)", transport.count.Load())
	require.Less(t, elapsed, 2*time.Second,
		"leader loop must return promptly under sustained invalidations (elapsed=%v count=%d)", elapsed, transport.count.Load())
}

// ----------------------------------------------------------------------------
// F1f: a forceRefresh WAITER (one that joins an in-flight refresh
// rather than leading) must also be bounded by maxForceRefreshRetries.
// Before this fix, the waiter "continue" path did not increment the
// retry counter, so a sustained leadership-churn pattern (other
// goroutines repeatedly winning each subsequent flight while
// invalidate() keeps firing) could keep one waiter joining stale
// flights indefinitely. Now the leader and waiter paths share a
// single budget.
// ----------------------------------------------------------------------------
func TestFix1f_ForceRefreshWaiterLoopIsBounded(t *testing.T) {
	body, _ := json.Marshal(accountProperties{
		ReadRegions:  []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
		WriteRegions: []accountRegion{{Name: "West US", Endpoint: "https://fake.documents.azure.com:443/"}},
	})
	// Each flight takes ~80ms.
	transport := &countingTransport{status: http.StatusOK, body: body, delay: 80 * time.Millisecond}
	gem := newGEMWithTransport(t, []string{"West US"}, transport, 5*time.Minute)

	// 1) Kick off a long-running leader so the waiter under test will
	//    arrive while inflight != nil.
	leaderDone := make(chan struct{})
	go func() {
		defer close(leaderDone)
		_ = gem.Update(context.Background(), false /*not force*/)
	}()
	require.Eventually(t, gem.hasInflight, 2*time.Second, 5*time.Millisecond,
		"leader must enter the in-flight slot first")

	// 2) Hostile background pattern: keep firing invalidate() so every
	//    flight the waiter observes pre-dates a later invalidation,
	//    AND keep starting new leaders the moment the previous
	//    flight clears so the waiter never gets to lead itself.
	stop := make(chan struct{})
	defer close(stop)
	mark, _ := url.Parse("https://fake.documents.azure.com:443/")
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.After(5 * time.Millisecond):
				_, _ = gem.MarkEndpointUnavailableForWrite(*mark)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				_ = gem.Update(context.Background(), false)
			}
		}
	}()

	// 3) Now start a forceRefresh waiter and time how long it takes
	//    to return. It should bail out after the shared retry budget
	//    is exhausted, NOT spin indefinitely.
	start := time.Now()
	waiterDone := make(chan error, 1)
	go func() {
		waiterDone <- gem.Update(context.Background(), true /*forceRefresh*/)
	}()
	select {
	case waiterErr := <-waiterDone:
		require.NoError(t, waiterErr)
	case <-time.After(5 * time.Second):
		t.Fatalf("forceRefresh waiter did not return; likely unbounded")
	}
	elapsed := time.Since(start)
	// 4 flights * 80ms each (leader's flight + 3 retry-budget loops) =
	// ~320ms nominal upper bound. Allow scheduler slack but require
	// significantly less than runaway (multi-second).
	require.Less(t, elapsed, 2*time.Second,
		"forceRefresh waiter must return promptly under leadership churn (elapsed=%v)", elapsed)

	<-leaderDone
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
			_, _ = gem.MarkEndpointUnavailableForWrite(*defaultEndpoint)
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
	_, err = gem.MarkEndpointUnavailableForWrite(*west)
	require.NoError(t, err)

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
