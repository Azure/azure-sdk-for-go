// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const defaultExpirationTime time.Duration = time.Minute * 5

const (
	none  requestedOperations = 0
	read                      = 0b01
	write                     = 0b10
	all                       = read | write
)

type requestedOperations int

type locationUnavailabilityInfo struct {
	lastCheckTime  time.Time
	unavailableOps requestedOperations
}

type databaseAccountLocationsInfo struct {
	prefLocations                 []regionId
	availWriteLocations           []regionId
	availReadLocations            []regionId
	availWriteEndpointsByLocation map[regionId]url.URL
	availReadEndpointsByLocation  map[regionId]url.URL
	writeEndpoints                []url.URL
	readEndpoints                 []url.URL
}

type accountRegion struct {
	Name     regionId `json:"name"`
	Endpoint string   `json:"databaseAccountEndpoint"`
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
	prefRegions := make([]regionId, len(prefLocations))
	for i, loc := range prefLocations {
		prefRegions[i] = newRegionId(loc)
	}
	return &locationCache{
		defaultEndpoint:                   defaultEndpoint,
		locationInfo:                      *newDatabaseAccountLocationsInfo(prefRegions, defaultEndpoint),
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
		prefRegions := make([]regionId, len(prefList))
		for i, loc := range prefList {
			prefRegions[i] = newRegionId(loc)
		}
		nextLoc.prefLocations = prefRegions
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

	// Only compare and log if the event is enabled
	if log.Should(EventEndpointManager) {
		writeEndpointsChanged := !urlSlicesEqual(lc.locationInfo.writeEndpoints, nextLoc.writeEndpoints)
		readEndpointsChanged := !urlSlicesEqual(lc.locationInfo.readEndpoints, nextLoc.readEndpoints)

		if writeEndpointsChanged || readEndpointsChanged {
			log.Writef(EventEndpointManager,
				"\n===== Endpoint Priority Recomputed =====\n"+
					"Preferred regions: %s\n"+
					"Write endpoint priority: %s\n"+
					"Read endpoint priority: %s\n"+
					"Multi-write locations enabled: %v\n"+
					"=========================================\n",
				formatRegionList(nextLoc.prefLocations),
				formatEndpointList(nextLoc.writeEndpoints),
				formatEndpointList(nextLoc.readEndpoints),
				lc.enableMultipleWriteLocations)
		}
	}

	lc.lastUpdateTime = time.Now()
	lc.locationInfo = nextLoc
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

func (lc *locationCache) getLocation(endpoint url.URL) regionId {
	// Take a read lock for the duration of the lookup. The reads of
	// locationInfo.availWriteEndpointsByLocation /
	// availReadEndpointsByLocation and enableMultipleWriteLocations race
	// the writes in update / updateLocked, especially now that the retry
	// policy's asyncForceRefreshGEM can trigger a refresh concurrently
	// with the data-plane lookup that calls into here.
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return lc.getLocationLocked(endpoint)
}

// getLocationLocked is the non-locking variant of getLocation for callers
// that already hold lc.mapMutex (read or write).
func (lc *locationCache) getLocationLocked(endpoint url.URL) regionId {
	var firstLoc regionId
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

// sessionRetrySnapshot returns a coherent snapshot of the fields the
// session-unavailable retry path needs to make a routing decision:
// (canUseMultipleWriteLocs, availReadLocationCount, availWriteLocationCount).
// Taking these reads under a single RLock prevents a concurrent
// locationCache.update (e.g. from an async GEM refresh) from rewriting
// enableMultipleWriteLocations between the multi-write branch decision
// and the slice-length sampling that follows it.
func (lc *locationCache) sessionRetrySnapshot() (multiWrite bool, readN, writeN int) {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return lc.enableMultipleWriteLocations,
		len(lc.locationInfo.availReadLocations),
		len(lc.locationInfo.availWriteLocations)
}

// readEndpointCount returns the number of resolved preferred read endpoints
// under RLock. The server-error retry path uses it to decide whether a
// cross-region failover would actually target a different endpoint. Reading
// the slice length under the lock prevents a torn read against a concurrent
// locationCache.update (e.g. from an async GEM refresh).
func (lc *locationCache) readEndpointCount() int {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	return len(lc.locationInfo.readEndpoints)
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
	region := lc.getLocation(endpoint)

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
	log.Writef(EventEndpointManager,
		"Marked endpoint unavailable: endpoint=%s, region=%s, operation=%s",
		endpoint.Host, region, operationName(op))
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
		t := time.Since(info.lastCheckTime)
		if t > lc.unavailableLocationExpirationTime {
			region := lc.getLocationLocked(endpoint)
			log.Writef(EventEndpointManager,
				"Endpoint is now available: endpoint=%s, region=%s, unavailableFor=%v",
				endpoint.Host, region, t)
			delete(lc.locationUnavailabilityInfoMap, endpoint)
		}
	}
}

// forceRefreshStaleEndpoints forces all stale endpoints to be cleared immediately (for testing)
func (lc *locationCache) forceRefreshStaleEndpoints() {
	lc.mapMutex.Lock()
	defer lc.mapMutex.Unlock()
	for endpoint, info := range lc.locationUnavailabilityInfoMap {
		t := time.Since(info.lastCheckTime)
		region := lc.getLocationLocked(endpoint)
		log.Writef(EventEndpointManager,
			"Endpoint is now available: endpoint=%s, region=%s, unavailableFor=%v",
			endpoint.Host, region, t)
		delete(lc.locationUnavailabilityInfoMap, endpoint)
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
func (lc *locationCache) getPrefAvailableEndpointsLocked(endpointsByLoc map[regionId]url.URL, locs []regionId, prefLocations []regionId, availOps requestedOperations, fallbackEndpoint url.URL) []url.URL {
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

func getEndpointsByLocation(locs []accountRegion) (map[regionId]url.URL, []regionId, error) {
	endpointsByLoc := make(map[regionId]url.URL)
	parsedLocs := make([]regionId, 0)
	for _, loc := range locs {
		endpoint, err := url.Parse(loc.Endpoint)
		if err != nil {
			return nil, nil, err
		}
		if loc.Name != "" {
			// Always canonicalize so that names set directly in Go code
			// (e.g., accountRegion{Name: "East US"}) match the canonicalized
			// prefLocations produced by newRegionId / UnmarshalJSON.
			canonical := newRegionId(string(loc.Name))
			endpointsByLoc[canonical] = *endpoint
			parsedLocs = append(parsedLocs, canonical)
		}
		// TODO else: log
	}
	return endpointsByLoc, parsedLocs, nil
}

func newDatabaseAccountLocationsInfo(prefLocations []regionId, defaultEndpoint url.URL) *databaseAccountLocationsInfo {
	availWriteLocs := make([]regionId, 0)
	availReadLocs := make([]regionId, 0)
	availWriteEndpointsByLocation := make(map[regionId]url.URL)
	availReadEndpointsByLocation := make(map[regionId]url.URL)
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

// Helper function to format region lists for logging
func formatRegionList(regions []regionId) string {
	if len(regions) == 0 {
		return "[]"
	}
	strs := make([]string, len(regions))
	for i, r := range regions {
		strs[i] = r.String()
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

// Helper function to format endpoint lists for logging
func formatEndpointList(endpoints []url.URL) string {
	if len(endpoints) == 0 {
		return "[]"
	}
	strs := make([]string, len(endpoints))
	for i, e := range endpoints {
		strs[i] = e.Host
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

// Helper to get operation name string
func operationName(op requestedOperations) string {
	switch op {
	case read:
		return "read"
	case write:
		return "write"
	case all:
		return "all"
	default:
		return "none"
	}
}

// Helper to compare URL slices
func urlSlicesEqual(a, b []url.URL) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Host != b[i].Host {
			return false
		}
	}
	return true
}
