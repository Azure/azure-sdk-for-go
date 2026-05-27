// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

const defaultExpirationTime time.Duration = time.Minute * 5

const (
	none requestedOperations = iota
	read
	write
	all
)

type requestedOperations int

type locationUnavailabilityInfo struct {
	lastCheckTime  time.Time
	unavailableOps requestedOperations
}

type databaseAccountLocationsInfo struct {
	prefLocations                 []string
	availWriteLocations           []string
	availReadLocations            []string
	availWriteEndpointsByLocation map[string]url.URL
	availReadEndpointsByLocation  map[string]url.URL
	writeEndpoints                []url.URL
	readEndpoints                 []url.URL
}

type accountRegion struct {
	Name     string `json:"name"`
	Endpoint string `json:"databaseAccountEndpoint"`
}

type userConsistencyPolicy struct {
	DefaultConsistencyLevel string `json:"defaultConsistencyLevel"`
}

type accountProperties struct {
	ReadRegions                  []accountRegion       `json:"readableLocations"`
	WriteRegions                 []accountRegion       `json:"writableLocations"`
	EnableMultipleWriteLocations bool                  `json:"enableMultipleWriteLocations"`
	AccountConsistency           userConsistencyPolicy `json:"userConsistencyPolicy"`
}

func (accountProps accountProperties) String() string {
	return fmt.Sprintf("Read regions: %v\nWrite regions: %v\nMulti-region writes: %v\nAccount consistency level: %v",
		accountProps.ReadRegions, accountProps.WriteRegions, accountProps.EnableMultipleWriteLocations, accountProps.AccountConsistency.DefaultConsistencyLevel)
}

type locationCache struct {
	locationInfo                      databaseAccountLocationsInfo
	defaultEndpoint                   url.URL
	enableCrossRegionRetries          bool
	locationUnavailabilityInfoMap     map[url.URL]locationUnavailabilityInfo
	mapMutex                          sync.RWMutex
	lastUpdateTime                    time.Time
	enableMultipleWriteLocations      bool
	unavailableLocationExpirationTime time.Duration
}

func newLocationCache(prefLocations []string, defaultEndpoint url.URL, enableCrossRegionRetries bool) *locationCache {
	return &locationCache{
		defaultEndpoint:                   defaultEndpoint,
		locationInfo:                      *newDatabaseAccountLocationsInfo(prefLocations, defaultEndpoint),
		locationUnavailabilityInfoMap:     make(map[url.URL]locationUnavailabilityInfo),
		unavailableLocationExpirationTime: defaultExpirationTime,
		enableCrossRegionRetries:          enableCrossRegionRetries,
	}
}

// update refreshes the location cache. It acquires mapMutex internally; do
// not call it while holding mapMutex (use updateLocked from inside such
// sections). Public callers go through update; internal callers that already
// hold the write lock call updateLocked.
func (lc *locationCache) update(writeLocations []accountRegion, readLocations []accountRegion, prefList []string, enableMultipleWriteLocations *bool) error {
	lc.mapMutex.Lock()
	defer lc.mapMutex.Unlock()
	return lc.updateLocked(writeLocations, readLocations, prefList, enableMultipleWriteLocations)
}

func (lc *locationCache) updateLocked(writeLocations []accountRegion, readLocations []accountRegion, prefList []string, enableMultipleWriteLocations *bool) error {
	nextLoc := copyDatabaseAccountLocationsInfo(lc.locationInfo)
	if prefList != nil {
		nextLoc.prefLocations = prefList
	}
	if enableMultipleWriteLocations != nil {
		lc.enableMultipleWriteLocations = *enableMultipleWriteLocations
	}
	lc.refreshStaleEndpointsLocked()
	if readLocations != nil {
		availReadEndpointsByLocation, availReadLocations, err := getEndpointsByLocation(readLocations)
		if err != nil {
			return err
		}
		nextLoc.availReadEndpointsByLocation = availReadEndpointsByLocation
		nextLoc.availReadLocations = availReadLocations
	}

	if writeLocations != nil {
		availWriteEndpointsByLocation, availWriteLocations, err := getEndpointsByLocation(writeLocations)
		if err != nil {
			return err
		}
		nextLoc.availWriteEndpointsByLocation = availWriteEndpointsByLocation
		nextLoc.availWriteLocations = availWriteLocations
	}

	// Choose regional fallbacks so the route lists never trail into the
	// customer-supplied default endpoint. The only data-plane code path
	// that may still hit the default endpoint is the degenerate "zero
	// write regions on a write" case.
	writeFallback := lc.defaultEndpoint
	if len(nextLoc.availWriteLocations) > 0 {
		if ep, ok := nextLoc.availWriteEndpointsByLocation[nextLoc.availWriteLocations[0]]; ok {
			writeFallback = ep
		}
	}
	nextLoc.writeEndpoints = lc.getPrefAvailableEndpointsLocked(nextLoc.availWriteEndpointsByLocation, nextLoc.availWriteLocations, nextLoc.prefLocations, write, writeFallback)
	// Prefer the first available read region for the read fallback. Only
	// fall back to the first write endpoint (or, transitively, the default
	// endpoint) when the account advertises zero read regions -- accounts
	// with valid read regions must never resolve a read to the default
	// endpoint.
	readFallback := nextLoc.writeEndpoints[0]
	if len(nextLoc.availReadLocations) > 0 {
		if ep, ok := nextLoc.availReadEndpointsByLocation[nextLoc.availReadLocations[0]]; ok {
			readFallback = ep
		}
	}
	nextLoc.readEndpoints = lc.getPrefAvailableEndpointsLocked(nextLoc.availReadEndpointsByLocation, nextLoc.availReadLocations, nextLoc.prefLocations, read, readFallback)
	lc.lastUpdateTime = time.Now()
	lc.locationInfo = nextLoc
	// TODO: log
	return nil
}

func (lc *locationCache) resolveServiceEndpoint(locationIndex int, resourceType resourceType, isWriteOperation, useWriteEndpoint bool) url.URL {
	// Take a read lock for the duration of endpoint resolution. The
	// fields read here (locationInfo, enableMultipleWriteLocations) are
	// rewritten atomically under mapMutex.Lock() by update/updateLocked,
	// and a concurrent forced refresh (e.g. from the retry policy's
	// asyncForceRefreshGEM or the GEM policy's background refresh) can
	// race with us without this lock.
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	if (isWriteOperation || useWriteEndpoint) && !lc.canUseMultipleWriteLocsToRoute(resourceType) {
		if lc.enableCrossRegionRetries && len(lc.locationInfo.availWriteLocations) > 0 {
			locationIndex = min(locationIndex%2, len(lc.locationInfo.availWriteLocations)-1)
			writeLocation := lc.locationInfo.availWriteLocations[locationIndex]
			return lc.locationInfo.availWriteEndpointsByLocation[writeLocation]
		}
		return lc.defaultEndpoint
	}

	endpoints := lc.locationInfo.readEndpoints
	if isWriteOperation {
		endpoints = lc.locationInfo.writeEndpoints
	}
	return endpoints[locationIndex%len(endpoints)]
}

func (lc *locationCache) canUseMultipleWriteLocsToRoute(resourceType resourceType) bool {
	return lc.canUseMultipleWriteLocs() && resourceType == resourceTypeDocument
}

// readEndpoints returns the cached preferred read endpoints, refreshing the
// stale-endpoint set if the cache hasn't been updated within the
// unavailableLocationExpirationTime AND at least one unavailability entry is
// recorded. The refresh path used to call lc.update while still holding
// mapMutex.RLock, which deadlocks with refreshStaleEndpoints's Lock()
// (sync.RWMutex cannot upgrade RLock to Lock). We now capture the staleness
// decision under RLock, release it, and let update acquire the write lock.
func (lc *locationCache) readEndpoints() ([]url.URL, error) {
	lc.mapMutex.RLock()
	stale := time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0
	lc.mapMutex.RUnlock()
	if stale {
		if err := lc.update(nil, nil, nil, nil); err != nil {
			return nil, err
		}
	}
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return lc.locationInfo.readEndpoints, nil
}

func (lc *locationCache) writeEndpoints() ([]url.URL, error) {
	lc.mapMutex.RLock()
	stale := time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0
	lc.mapMutex.RUnlock()
	if stale {
		if err := lc.update(nil, nil, nil, nil); err != nil {
			return nil, err
		}
	}
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return lc.locationInfo.writeEndpoints, nil
}

func (lc *locationCache) getLocation(endpoint url.URL) string {
	firstLoc := ""
	for location, uri := range lc.locationInfo.availWriteEndpointsByLocation {
		if uri == endpoint {
			return location
		}
		if firstLoc == "" {
			firstLoc = location
		}
	}

	for location, uri := range lc.locationInfo.availReadEndpointsByLocation {
		if uri == endpoint {
			return location
		}
	}

	if endpoint == lc.defaultEndpoint && !lc.canUseMultipleWriteLocs() {
		if len(lc.locationInfo.availWriteEndpointsByLocation) > 0 {
			return firstLoc
		}
	}
	return ""
}

// canUseMultipleWriteLocs returns whether the account supports multi-master
// writes. Callers that already hold lc.mapMutex use this; the public
// CanUseMultipleWriteLocations entrypoint locks first.
func (lc *locationCache) canUseMultipleWriteLocs() bool {
	return lc.enableMultipleWriteLocations
}

func (lc *locationCache) CanUseMultipleWriteLocs() bool {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return lc.enableMultipleWriteLocations
}

func (lc *locationCache) markEndpointUnavailableForRead(endpoint url.URL) (wasAlreadyUnavailable bool, err error) {
	return lc.markEndpointUnavailable(endpointKey(endpoint), read)
}

func (lc *locationCache) markEndpointUnavailableForWrite(endpoint url.URL) (wasAlreadyUnavailable bool, err error) {
	return lc.markEndpointUnavailable(endpointKey(endpoint), write)
}

// endpointKey normalizes a url.URL to the form used as a key in
// locationUnavailabilityInfoMap and stored in availReadEndpointsByLocation
// / availWriteEndpointsByLocation: scheme + host only. Callers of
// MarkEndpointUnavailable* commonly pass the full request URL (including
// path, query, fragment, RawPath, etc.), which would never struct-equal
// the base URLs the cache uses; without normalization the marks are
// recorded under keys nothing else ever looks up and the demote step in
// getPrefAvailableEndpointsLocked silently does nothing.
func endpointKey(u url.URL) url.URL {
	return url.URL{Scheme: u.Scheme, Host: u.Host}
}

// markEndpointUnavailable atomically samples whether the endpoint was already
// unavailable for `op` and records the unavailability. Returning the prior
// state from inside the same critical section that performs the mark
// eliminates the check-then-act race exploited by concurrent callers.
func (lc *locationCache) markEndpointUnavailable(endpoint url.URL, op requestedOperations) (wasAlreadyUnavailable bool, err error) {
	now := time.Now()
	lc.mapMutex.Lock()
	defer lc.mapMutex.Unlock()
	if info, ok := lc.locationUnavailabilityInfoMap[endpoint]; ok {
		wasAlreadyUnavailable = op&info.unavailableOps == op &&
			time.Since(info.lastCheckTime) < lc.unavailableLocationExpirationTime
		info.lastCheckTime = now
		info.unavailableOps |= op
		lc.locationUnavailabilityInfoMap[endpoint] = info
	} else {
		lc.locationUnavailabilityInfoMap[endpoint] = locationUnavailabilityInfo{
			lastCheckTime:  now,
			unavailableOps: op,
		}
	}
	// Fast path: if the endpoint was already unavailable for this op within
	// the unavailability window, the route lists cannot change -- skip the
	// full updateLocked recompute. The only state mutation is the bumped
	// lastCheckTime, which getPrefAvailableEndpointsLocked observes via
	// isEndpointUnavailableLocked and which yields the same result.
	if wasAlreadyUnavailable {
		return wasAlreadyUnavailable, nil
	}
	return wasAlreadyUnavailable, lc.updateLocked(nil, nil, nil, nil)
}

func (lc *locationCache) databaseAccountRead(dbAcct accountProperties) error {
	return lc.update(dbAcct.WriteRegions, dbAcct.ReadRegions, nil, &dbAcct.EnableMultipleWriteLocations)
}

func (lc *locationCache) refreshStaleEndpoints() {
	lc.mapMutex.Lock()
	defer lc.mapMutex.Unlock()
	lc.refreshStaleEndpointsLocked()
}

func (lc *locationCache) refreshStaleEndpointsLocked() {
	for endpoint, info := range lc.locationUnavailabilityInfoMap {
		if time.Since(info.lastCheckTime) > lc.unavailableLocationExpirationTime {
			delete(lc.locationUnavailabilityInfoMap, endpoint)
		}
	}
}

func (lc *locationCache) isEndpointUnavailable(endpoint url.URL, ops requestedOperations) bool {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return lc.isEndpointUnavailableLocked(endpoint, ops)
}

func (lc *locationCache) isEndpointUnavailableLocked(endpoint url.URL, ops requestedOperations) bool {
	info, ok := lc.locationUnavailabilityInfoMap[endpointKey(endpoint)]
	if ops == none || !ok || ops&info.unavailableOps != ops {
		return false
	}
	return time.Since(info.lastCheckTime) < lc.unavailableLocationExpirationTime
}

// getPrefAvailableEndpointsLocked returns the endpoints for the customer's
// preferred locations in priority order, with unavailable endpoints moved to
// the tail. Callers pass prefLocations explicitly so updateLocked can compute
// route lists from the in-progress nextLoc snapshot rather than the
// already-committed lc.locationInfo.
func (lc *locationCache) getPrefAvailableEndpointsLocked(endpointsByLoc map[string]url.URL, locs []string, prefLocations []string, availOps requestedOperations, fallbackEndpoint url.URL) []url.URL {
	endpoints := make([]url.URL, 0)
	if lc.enableCrossRegionRetries {
		if lc.canUseMultipleWriteLocs() || availOps&read != 0 {
			unavailEndpoints := make([]url.URL, 0)
			for _, loc := range prefLocations {
				if endpoint, ok := endpointsByLoc[loc]; ok {
					if lc.isEndpointUnavailableLocked(endpoint, availOps) {
						unavailEndpoints = append(unavailEndpoints, endpoint)
					} else {
						endpoints = append(endpoints, endpoint)
					}
				}
			}
			endpoints = append(endpoints, unavailEndpoints...)
		} else {
			for _, loc := range locs {
				if endpoint, ok := endpointsByLoc[loc]; ok && loc != "" {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}
	if len(endpoints) == 0 {
		// Last resort: none of the customer's preferred regions matched the
		// account's regions (or cross-region retries are off and the
		// non-multi-master read branch yielded nothing, or the account
		// itself advertises zero regions). The caller passes a regional
		// fallback whenever availWriteLocations is non-empty, so the only
		// time this is the customer-supplied default endpoint is the
		// degenerate "zero regions" case.
		endpoints = append(endpoints, fallbackEndpoint)
	}
	return endpoints
}

func getEndpointsByLocation(locs []accountRegion) (map[string]url.URL, []string, error) {
	endpointsByLoc := make(map[string]url.URL)
	parsedLocs := make([]string, 0)
	for _, loc := range locs {
		endpoint, err := url.Parse(loc.Endpoint)
		if err != nil {
			return nil, nil, err
		}
		if loc.Name != "" {
			endpointsByLoc[loc.Name] = *endpoint
			parsedLocs = append(parsedLocs, loc.Name)
		}
		// TODO else: log
	}
	return endpointsByLoc, parsedLocs, nil
}

func newDatabaseAccountLocationsInfo(prefLocations []string, defaultEndpoint url.URL) *databaseAccountLocationsInfo {
	availWriteLocs := make([]string, 0)
	availReadLocs := make([]string, 0)
	availWriteEndpointsByLocation := make(map[string]url.URL)
	availReadEndpointsByLocation := make(map[string]url.URL)
	// Pre-populated seed: the lists contain defaultEndpoint until the first
	// successful Update() replaces them with regional endpoints. This is
	// safe because the pipeline policy (globalEndpointManagerPolicy) blocks
	// data-plane requests on a synchronous bootstrap and surfaces the GEM
	// error if it fails, so resolveServiceEndpoint is never consulted for a
	// real data-plane request while these seeded values are still in effect.
	writeEndpoints := []url.URL{defaultEndpoint}
	readEndpoints := []url.URL{defaultEndpoint}
	return &databaseAccountLocationsInfo{
		prefLocations:                 prefLocations,
		availWriteLocations:           availWriteLocs,
		availReadLocations:            availReadLocs,
		availWriteEndpointsByLocation: availWriteEndpointsByLocation,
		availReadEndpointsByLocation:  availReadEndpointsByLocation,
		writeEndpoints:                writeEndpoints,
		readEndpoints:                 readEndpoints,
	}
}

func copyDatabaseAccountLocationsInfo(other databaseAccountLocationsInfo) databaseAccountLocationsInfo {
	return databaseAccountLocationsInfo{
		prefLocations:                 other.prefLocations,
		availWriteLocations:           other.availWriteLocations,
		availReadLocations:            other.availReadLocations,
		availWriteEndpointsByLocation: other.availWriteEndpointsByLocation,
		availReadEndpointsByLocation:  other.availReadEndpointsByLocation,
		writeEndpoints:                other.writeEndpoints,
		readEndpoints:                 other.readEndpoints,
	}
}
