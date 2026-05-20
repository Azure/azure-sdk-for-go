// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
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
	// throttled to refreshTimeInterval just like a successful one. This
	// prevents the failure-loop described in
	// https://github.com/Azure/azure-sdk-for-go/issues/25468 where every
	// caller after a failed Update would re-issue the HTTP call immediately.
	lastAttemptTime time.Time
	// lastUpdateErr is the error from the most recent refresh attempt. It is
	// returned to callers that hit the throttle window before the GEM has
	// ever been successfully populated, so a chronic bootstrap failure is
	// surfaced on every request rather than silently swallowed.
	lastUpdateErr error
	// everPopulated is set true on the first successful refresh and never
	// reset, even when invalidate() zeroes lastUpdateTime to force a fresh
	// refresh. populated() reads this flag rather than the timestamp, so a
	// transient post-invalidation refresh failure does not stall the data
	// plane: requests can still route through the existing locationCache
	// topology while the next refresh attempt is pending the throttle
	// window. See issue #25468 deep-review finding.
	everPopulated bool
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

func (gem *globalEndpointManager) MarkEndpointUnavailableForWrite(endpoint url.URL) error {
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
	wasAlreadyUnavailable, err := gem.locationCache.markEndpointUnavailableForWrite(endpoint)
	if err != nil {
		return err
	}
	if !wasAlreadyUnavailable {
		gem.invalidate()
	}
	return nil
}

func (gem *globalEndpointManager) MarkEndpointUnavailableForRead(endpoint url.URL) error {
	wasAlreadyUnavailable, err := gem.locationCache.markEndpointUnavailableForRead(endpoint)
	if err != nil {
		return err
	}
	if !wasAlreadyUnavailable {
		gem.invalidate()
	}
	return nil
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
	return gem.locationCache.canUseMultipleWriteLocs()
}

func (gem *globalEndpointManager) IsEndpointUnavailable(endpoint url.URL, ops requestedOperations) bool {
	return gem.locationCache.isEndpointUnavailable(endpoint, ops)
}

func (gem *globalEndpointManager) RefreshStaleEndpoints() {
	gem.locationCache.refreshStaleEndpoints()
}

// populated reports whether the GEM has ever been successfully populated.
// Unlike a check of !lastUpdateTime.IsZero(), this remains true after
// invalidate() zeroes the timestamps -- so a transient post-invalidation
// refresh failure does not stall the data plane: the policy can still
// route requests through the existing locationCache topology while the
// next refresh attempt waits for the throttle window. See issue #25468
// deep-review finding.
func (gem *globalEndpointManager) populated() bool {
	gem.gemMutex.Lock()
	defer gem.gemMutex.Unlock()
	return gem.everPopulated
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
	// immediately (issue #25468).
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
// refreshTimeInterval -- this prevents the failure-storm described in
// https://github.com/Azure/azure-sdk-for-go/issues/25468.
//
// If the GEM has never been successfully populated and the throttle is
// active, Update returns the cached error from the most recent failed
// attempt so a chronic bootstrap failure is surfaced on every request
// rather than silently swallowed. Once the GEM HAS been populated, a
// throttled Update returns nil even after a subsequent refresh failure --
// the data plane keeps using the cached topology until the throttle expires
// and a fresh refresh can run.
//
// Update respects ctx cancellation in both leader and waiter roles. The
// leader's HTTP call uses ctx directly (or a derived context inside
// GetAccountProperties). Waiters select between flight completion and
// ctx.Done() so a caller-side timeout cannot be exceeded by an unrelated
// stuck refresh.
func (gem *globalEndpointManager) Update(ctx context.Context, forceRefresh bool) error {
	gem.gemMutex.Lock()
	if !gem.shouldRefresh() && !forceRefresh {
		// Throttled. Surface the cached error only if we have NEVER
		// successfully populated the GEM -- otherwise the data plane has
		// a valid cached topology and should continue working until the
		// next refresh attempt succeeds.
		var cached error
		if !gem.everPopulated {
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
		flight := gem.inflight
		gem.gemMutex.Unlock()
		select {
		case <-flight.done:
			return flight.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	// We are the leader. Publish the inflight flight and snapshot the
	// invalidation generation, then release the lock while we perform the
	// HTTP call so ShouldRefresh and other non-Update paths don't block on
	// a network round-trip.
	flight := &updateFlight{done: make(chan struct{})}
	gem.inflight = flight
	genAtStart := gem.invalidationGen
	gem.gemMutex.Unlock()

	// Panic-safe cleanup: if refreshOnce (or anything it transitively calls
	// -- the pipeline, JSON unmarshal, locationCache.update) panics, we
	// MUST still clear gem.inflight and close flight.done, otherwise every
	// subsequent Update caller blocks forever on <-flight.done. We capture
	// any panic, record it as the flight error, and re-panic after cleanup.
	var err error
	func() {
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
					gem.everPopulated = true
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
		err = gem.refreshOnce(ctx)
	}()
	return err
}

// refreshOnce performs the actual GetAccountProperties HTTP call and
// propagates the result into the location cache. It must not hold gemMutex.
func (gem *globalEndpointManager) refreshOnce(ctx context.Context) error {
	accountProperties, err := gem.GetAccountProperties(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve account properties: %v", err)
	}
	if err := gem.locationCache.update(
		accountProperties.WriteRegions,
		accountProperties.ReadRegions,
		gem.preferredLocations,
		&accountProperties.EnableMultipleWriteLocations,
	); err != nil {
		return fmt.Errorf("failed to update location cache: %v", err)
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

	req.Raw().Header.Set(headerXmsDate, time.Now().UTC().Format(http.TimeFormat))
	req.Raw().Header.Set(headerXmsVersion, apiVersion)
	req.Raw().Header.Set(cosmosHeaderSDKSupportedCapabilities, supportedCapabilitiesHeaderValue)

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
