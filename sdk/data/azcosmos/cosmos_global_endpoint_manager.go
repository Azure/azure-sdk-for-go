// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const defaultUnavailableLocationRefreshInterval = 5 * time.Minute

type globalEndpointManager struct {
	clientEndpoint      string
	pipeline            azruntime.Pipeline
	preferredLocations  []string
	locationCache       *locationCache
	refreshTimeInterval time.Duration
	gemMutex            sync.Mutex
	lastUpdateTime      time.Time
	// lastAttemptTime records the most recent Update attempt regardless of
	// outcome. shouldRefresh() honours it so a failed GetAccountProperties is
	// throttled to refreshTimeInterval just like a successful one, preventing
	// a failure-loop where every caller re-issues the HTTP call immediately.
	lastAttemptTime time.Time
	// lastUpdateErr is the error from the most recent refresh attempt. It is
	// returned to callers that hit the throttle window before the GEM has
	// ever been successfully populated, so a chronic bootstrap failure is
	// surfaced on every request rather than silently swallowed.
	lastUpdateErr error
	// everPopulated is set true on the first successful refresh and never
	// reset, even when invalidate() zeroes lastUpdateTime. populated() reads
	// this flag rather than the timestamp, so a transient post-invalidation
	// refresh failure does not stall the data plane: requests can still route
	// through the existing locationCache topology while the next refresh
	// attempt waits for the throttle window. Stored as atomic.Bool so the
	// hot-path policy check is lock-free.
	everPopulated atomic.Bool
	// inflight coalesces concurrent Update callers: only the first does the
	// HTTP call; the rest wait on the per-flight done channel and read
	// per-flight err. Each refresh has its own *updateFlight so late waiters
	// cannot accidentally read a subsequent flight's error.
	inflight *updateFlight
	// invalidationGen is bumped by invalidate(). The Update leader snapshots
	// it before the HTTP call and, if it changed by completion, declines to
	// advance lastUpdateTime/lastAttemptTime -- so an invalidation that
	// happens during an in-flight refresh is not lost.
	invalidationGen uint64
}

// updateFlight tracks a single in-flight refresh.
type updateFlight struct {
	done chan struct{}
	err  error
	// genAtStart is the invalidationGen the leader observed when it
	// claimed the flight. Waiters use it to decide whether the flight
	// pre-dates a later invalidate() (in which case forceRefresh callers
	// must loop and start a fresh flight after this one completes).
	genAtStart uint64
}

func newGlobalEndpointManager(clientEndpoint string, pipeline azruntime.Pipeline, preferredLocations []string, refreshTimeInterval time.Duration, enableCrossRegionRetries bool) (*globalEndpointManager, error) {
	endpoint, err := url.Parse(clientEndpoint)
	if err != nil {
		return &globalEndpointManager{}, err
	}

	if refreshTimeInterval == 0 {
		refreshTimeInterval = defaultUnavailableLocationRefreshInterval
	}

	gem := &globalEndpointManager{
		clientEndpoint:      clientEndpoint,
		pipeline:            pipeline,
		preferredLocations:  preferredLocations,
		locationCache:       newLocationCache(preferredLocations, *endpoint, enableCrossRegionRetries),
		refreshTimeInterval: refreshTimeInterval,
		lastUpdateTime:      time.Time{},
	}

	return gem, nil
}

func (gem *globalEndpointManager) GetWriteEndpoints() ([]url.URL, error) {
	return gem.locationCache.writeEndpoints()
}

func (gem *globalEndpointManager) GetReadEndpoints() ([]url.URL, error) {
	return gem.locationCache.readEndpoints()
}

func (gem *globalEndpointManager) MarkEndpointUnavailableForWrite(endpoint url.URL) (wasAlreadyUnavailable bool, err error) {
	// markEndpointUnavailableForWrite atomically reports whether the
	// endpoint was already unavailable for write in the same critical
	// section that performs the mark. This eliminates the check-then-act
	// race that would otherwise let concurrent markers all observe
	// "wasn't unavailable" and each invalidate the GEM. A new event
	// invalidates the GEM cache so the next Update(false) actually fires
	// -- the first 403/WriteForbidden or network error for an endpoint
	// may indicate a failover and we want to learn about new write
	// regions promptly. Subsequent retries within the unavailability
	// window do not invalidate.
	//
	// The caller receives wasAlreadyUnavailable so it can decide whether
	// to additionally force a fresh refresh; the GEM-internal invalidate
	// already arms the next non-force Update either way.
	wasAlreadyUnavailable, err = gem.locationCache.markEndpointUnavailableForWrite(endpoint)
	if err != nil {
		return false, err
	}
	if !wasAlreadyUnavailable {
		gem.invalidate()
	}
	return wasAlreadyUnavailable, nil
}

func (gem *globalEndpointManager) MarkEndpointUnavailableForRead(endpoint url.URL) (wasAlreadyUnavailable bool, err error) {
	wasAlreadyUnavailable, err = gem.locationCache.markEndpointUnavailableForRead(endpoint)
	if err != nil {
		return false, err
	}
	if !wasAlreadyUnavailable {
		gem.invalidate()
	}
	return wasAlreadyUnavailable, nil
}

// invalidate forces the next non-force Update to actually issue a refresh by
// clearing both lastUpdateTime and lastAttemptTime, and bumps
// invalidationGen so that a refresh currently in flight cannot mask the
// invalidation by writing the post-call timestamps. Used when we learn about
// a newly-unavailable endpoint and want to discover potential failover
// targets without waiting for the next refresh interval.
func (gem *globalEndpointManager) invalidate() {
	gem.gemMutex.Lock()
	defer gem.gemMutex.Unlock()
	gem.lastUpdateTime = time.Time{}
	gem.lastAttemptTime = time.Time{}
	gem.invalidationGen++
}

func (gem *globalEndpointManager) GetEndpointLocation(endpoint url.URL) string {
	return gem.locationCache.getLocation(endpoint)
}

func (gem *globalEndpointManager) CanUseMultipleWriteLocations() bool {
	return gem.locationCache.CanUseMultipleWriteLocs()
}

func (gem *globalEndpointManager) IsEndpointUnavailable(endpoint url.URL, ops requestedOperations) bool {
	return gem.locationCache.isEndpointUnavailable(endpoint, ops)
}

func (gem *globalEndpointManager) RefreshStaleEndpoints() {
	gem.locationCache.refreshStaleEndpoints()
}

// populated reports whether the GEM has ever been successfully populated.
// Unlike !lastUpdateTime.IsZero(), this remains true after invalidate()
// zeroes the timestamps, so a transient post-invalidation refresh failure
// does not stall the data plane. Lock-free atomic read keeps the policy
// hot path off gemMutex.
func (gem *globalEndpointManager) populated() bool {
	return gem.everPopulated.Load()
}

// hasInflight is a test-only accessor for the in-flight refresh slot.
// Tests use it to wait until a leader has claimed the slot before firing
// follow-up calls; keeping the field access inside the GEM means tests
// don't need to grab gemMutex directly.
func (gem *globalEndpointManager) hasInflight() bool {
	gem.gemMutex.Lock()
	defer gem.gemMutex.Unlock()
	return gem.inflight != nil
}

func (gem *globalEndpointManager) ShouldRefresh() bool {
	gem.gemMutex.Lock()
	defer gem.gemMutex.Unlock()
	return gem.shouldRefresh()
}

func (gem *globalEndpointManager) shouldRefresh() bool {
	// Honor whichever happened more recently: a successful update or an
	// attempt that failed. Failures must be throttled too, otherwise a
	// failing endpoint causes every caller to re-issue GetDatabaseAccount
	// immediately.
	last := gem.lastUpdateTime
	if gem.lastAttemptTime.After(last) {
		last = gem.lastAttemptTime
	}
	return time.Since(last) > gem.refreshTimeInterval
}

func (gem *globalEndpointManager) ResolveServiceEndpoint(locationIndex int, resourceType resourceType, isWriteOperation, useWriteEndpoint bool) url.URL {
	return gem.locationCache.resolveServiceEndpoint(locationIndex, resourceType, isWriteOperation, useWriteEndpoint)
}

// Update refreshes the GEM cache by calling GetAccountProperties when needed.
// Concurrent callers are coalesced via a single-in-flight pattern so at most
// one HTTP call is in flight per client at any time. Both successes and
// failures advance lastAttemptTime so the next refresh is throttled to
// refreshTimeInterval, preventing a failure-storm.
//
// If the GEM has never been successfully populated and the throttle is
// active, Update returns the cached error from the most recent failed
// attempt so a chronic bootstrap failure is surfaced on every request
// rather than silently swallowed. Once the GEM HAS been populated, a
// throttled Update returns nil even after a subsequent refresh failure --
// the data plane keeps using the cached topology until the throttle expires
// and a fresh refresh can run.
//
// Update respects ctx cancellation for waiters. The leader's HTTP call runs
// under context.WithoutCancel(ctx) so an unrelated caller-side cancellation
// does not poison the shared flight result for other coalesced waiters
// (GetAccountProperties applies its own 60s timeout). Waiters select between
// flight completion and their own ctx.Done() so a caller-side timeout cannot
// be exceeded by an unrelated stuck refresh.
// maxForceRefreshRetries bounds the total number of times a
// forceRefresh caller may re-enter the leadership-or-wait loop after
// observing that the last flight predates a later invalidation. The
// cap covers BOTH paths: a leader whose own flight saw a mid-flight
// invalidate(), and a waiter that joined a flight which turned out
// to predate a later invalidate(). Without this combined cap, a
// hostile invalidation+leadership-churn pattern could keep a single
// goroutine bouncing between waiter and leader roles indefinitely.
// After the cap is reached the call returns its last result; the
// retry-policy state machine then sees the outcome (Idle on success,
// Failed on error) and the NEXT request that needs a fresh refresh
// will spawn a new goroutine, naturally rate-limiting the topology
// poll while still making progress.
const maxForceRefreshRetries = 3

func (gem *globalEndpointManager) Update(ctx context.Context, forceRefresh bool) error {
	retries := 0
	for {
		var flight *updateFlight
		var genAtStart uint64
		// Acquire leadership or wait on an in-flight refresh.
	leader:
		for {
			gem.gemMutex.Lock()
			if !gem.shouldRefresh() && !forceRefresh {
				// Throttled. Surface the cached error only if we have NEVER
				// successfully populated the GEM -- otherwise the data plane has
				// a valid cached topology and should continue working until the
				// next refresh attempt succeeds. The cached error is shared across
				// force=true and force=false callers: both want to surface
				// "bootstrap is broken" and there's no caller-visible distinction.
				var cached error
				if !gem.everPopulated.Load() {
					cached = gem.lastUpdateErr
				}
				gem.gemMutex.Unlock()
				return cached
			}
			if gem.inflight != nil {
				// Another goroutine is performing a refresh. Wait for it and share
				// its result rather than spawning a duplicate HTTP call. The result
				// lives on the per-flight struct so subsequent flights cannot
				// overwrite it. Honour the waiter's ctx so a caller-side timeout
				// is not extended by the leader's HTTP call duration.
				waitFlight := gem.inflight
				gem.gemMutex.Unlock()
				select {
				case <-waitFlight.done:
					if forceRefresh {
						// Re-read invalidationGen UNDER THE LOCK after the
						// flight completes. Sampling it before the wait would
						// miss any invalidate() that fires while we are
						// blocked on <-waitFlight.done -- in which case the
						// flight we waited on would still pre-date the
						// latest invalidation but we would not know it, and
						// we would return the stale result instead of
						// looping to lead a fresh post-invalidation refresh.
						gem.gemMutex.Lock()
						latestGen := gem.invalidationGen
						gem.gemMutex.Unlock()
						// Consume one slot of the shared force-refresh
						// retry budget so a waiter that keeps joining
						// stale flights cannot loop forever.
						if waitFlight.genAtStart < latestGen && retries < maxForceRefreshRetries {
							retries++
							continue
						}
					}
					return waitFlight.err
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			// We are the leader. Publish the inflight flight and snapshot the
			// invalidation generation BEFORE releasing the lock so the HTTP
			// call below does not block ShouldRefresh and other non-Update
			// callers on a network round-trip.
			genAtStart = gem.invalidationGen
			flight = &updateFlight{done: make(chan struct{}), genAtStart: genAtStart}
			gem.inflight = flight
			gem.gemMutex.Unlock()
			break leader
		}

		err := gem.runLeaderFlight(ctx, flight, genAtStart)

		// If we are a forceRefresh caller and an invalidate() fired during
		// our flight, our flight's genAtStart predates the latest
		// invalidation -- the topology we just learned does not reflect
		// the most recent unavailability mark. Loop and either lead a
		// fresh flight that covers it or coalesce into one another
		// goroutine has already started. Without this, a 403 (or
		// connection error) whose MarkEndpointUnavailable* invalidate()
		// happened mid-flight would be silently coalesced into the
		// in-progress refresh, leave asyncRefreshState=Idle in the retry
		// policy, and not trigger another forced refresh.
		//
		// The same shared retries budget bounds this leader-side
		// re-entry so a sustained invalidation storm cannot keep one
		// goroutine spinning here forever.
		if forceRefresh {
			gem.gemMutex.Lock()
			latestGen := gem.invalidationGen
			gem.gemMutex.Unlock()
			if latestGen > genAtStart && retries < maxForceRefreshRetries {
				retries++
				continue
			}
		}
		return err
	}
}

// runLeaderFlight performs the actual GetAccountProperties HTTP call as
// the in-flight leader and runs the panic-safe defer that commits the
// flight result and timestamps. It MUST be called only after the caller
// has claimed leadership (set gem.inflight = flight). Returns the
// flight's error (nil on success).
func (gem *globalEndpointManager) runLeaderFlight(ctx context.Context, flight *updateFlight, genAtStart uint64) (err error) {
	// Panic-safe cleanup: if refreshOnce (or anything it transitively calls
	// -- the pipeline, JSON unmarshal, locationCache.update) panics, we
	// MUST still clear gem.inflight and close flight.done, otherwise every
	// subsequent Update caller blocks forever on <-flight.done. We capture
	// any panic, record it as the flight error, and re-panic after
	// cleanup so the original stack trace is preserved for the host's
	// panic handler.
	defer func() {
		r := recover()
		gem.gemMutex.Lock()
		if r != nil && err == nil {
			err = fmt.Errorf("panic in GEM refresh: %v", r)
		}
		flight.err = err
		gem.lastUpdateErr = err
		if gem.invalidationGen == genAtStart {
			// No invalidation occurred during the flight, so commit the
			// timestamps and let the throttle take effect.
			gem.lastAttemptTime = time.Now()
			if err == nil {
				gem.lastUpdateTime = gem.lastAttemptTime
				gem.everPopulated.Store(true)
			}
		}
		// If invalidationGen changed, leave the timestamps untouched so
		// the next caller observes shouldRefresh()==true and performs a
		// fresh refresh that reflects the post-invalidation state.
		gem.inflight = nil
		gem.gemMutex.Unlock()
		close(flight.done)
		if r != nil {
			panic(r)
		}
	}()
	err = gem.refreshOnce(context.WithoutCancel(ctx))
	return err
}

// refreshOnce performs the actual GetAccountProperties HTTP call and
// propagates the result into the location cache. It must not hold gemMutex.
func (gem *globalEndpointManager) refreshOnce(ctx context.Context) error {
	accountProperties, err := gem.GetAccountProperties(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve account properties: %w", err)
	}
	if err := gem.locationCache.update(
		accountProperties.WriteRegions,
		accountProperties.ReadRegions,
		gem.preferredLocations,
		&accountProperties.EnableMultipleWriteLocations,
	); err != nil {
		return fmt.Errorf("failed to update location cache: %w", err)
	}
	return nil
}

func (gem *globalEndpointManager) GetAccountProperties(ctx context.Context) (accountProperties, error) {
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabaseAccount,
		resourceAddress: "",
	}

	ctxt, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := azruntime.NewRequest(ctxt, http.MethodGet, gem.clientEndpoint)
	if err != nil {
		return accountProperties{}, err
	}

	addDefaultHeaders(req)

	req.SetOperationValue(operationContext)

	azResponse, err := gem.pipeline.Do(req)
	if err != nil {
		return accountProperties{}, err
	}

	successResponse := (azResponse.StatusCode >= 200 && azResponse.StatusCode < 300)
	if successResponse {
		properties, err := newAccountProperties(azResponse)
		if err != nil {
			return accountProperties{}, fmt.Errorf("failed to parse account properties: %v", err)
		}
		log.Write(azlog.EventResponse, "\n===== Database Account Information:\n"+properties.String()+"\n=====\n")
		return properties, nil
	}

	return accountProperties{}, azruntime.NewResponseErrorWithErrorCode(azResponse, azResponse.Status)
}

func newAccountProperties(azResponse *http.Response) (accountProperties, error) {
	properties := accountProperties{}
	err := azruntime.UnmarshalAsJSON(azResponse, &properties)
	if err != nil {
		return properties, err
	}

	return properties, nil
}
