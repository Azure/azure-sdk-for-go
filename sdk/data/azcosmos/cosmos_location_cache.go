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

func (lc *locationCache) update(writeLocations []accountRegion, readLocations []accountRegion, prefList []string, enableMultipleWriteLocations *bool) error {
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
	lc.refreshStaleEndpoints()
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

	nextLoc.writeEndpoints = lc.getPrefAvailableEndpoints(nextLoc.availWriteEndpointsByLocation, nextLoc.availWriteLocations, write, lc.defaultEndpoint)
	nextLoc.readEndpoints = lc.getPrefAvailableEndpoints(nextLoc.availReadEndpointsByLocation, nextLoc.availReadLocations, read, nextLoc.writeEndpoints[0])

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

func (lc *locationCache) readEndpoints() ([]url.URL, error) {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		err := lc.update(nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	return lc.locationInfo.readEndpoints, nil
}

func (lc *locationCache) writeEndpoints() ([]url.URL, error) {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		err := lc.update(nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	return lc.locationInfo.writeEndpoints, nil
}

func (lc *locationCache) getLocation(endpoint url.URL) regionId {
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

func (lc *locationCache) canUseMultipleWriteLocs() bool {
	return lc.enableMultipleWriteLocations
}

func (lc *locationCache) markEndpointUnavailableForRead(endpoint url.URL) error {
	return lc.markEndpointUnavailable(endpoint, read)
}

func (lc *locationCache) markEndpointUnavailableForWrite(endpoint url.URL) error {
	return lc.markEndpointUnavailable(endpoint, write)
}

func (lc *locationCache) markEndpointUnavailable(endpoint url.URL, op requestedOperations) error {
	now := time.Now()
	region := lc.getLocation(endpoint)

	lc.mapMutex.Lock()
	if info, ok := lc.locationUnavailabilityInfoMap[endpoint]; ok {
		info.lastCheckTime = now
		info.unavailableOps |= op
		lc.locationUnavailabilityInfoMap[endpoint] = info
	} else {
		info = locationUnavailabilityInfo{
			lastCheckTime:  now,
			unavailableOps: op,
		}
		lc.locationUnavailabilityInfoMap[endpoint] = info
	}
	lc.mapMutex.Unlock()

	log.Writef(EventEndpointManager,
		"Marked endpoint unavailable: endpoint=%s, region=%s, operation=%s",
		endpoint.Host, region, operationName(op))

	err := lc.update(nil, nil, nil, nil)
	return err
}

func (lc *locationCache) databaseAccountRead(dbAcct accountProperties) error {
	return lc.update(dbAcct.WriteRegions, dbAcct.ReadRegions, nil, &dbAcct.EnableMultipleWriteLocations)
}

func (lc *locationCache) refreshStaleEndpoints() {
	lc.mapMutex.Lock()
	defer lc.mapMutex.Unlock()
	for endpoint, info := range lc.locationUnavailabilityInfoMap {
		t := time.Since(info.lastCheckTime)
		if t > lc.unavailableLocationExpirationTime {
			region := lc.getLocation(endpoint)
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
		region := lc.getLocation(endpoint)
		log.Writef(EventEndpointManager,
			"Endpoint is now available: endpoint=%s, region=%s, unavailableFor=%v",
			endpoint.Host, region, t)
		delete(lc.locationUnavailabilityInfoMap, endpoint)
	}
}

func (lc *locationCache) isEndpointUnavailable(endpoint url.URL, ops requestedOperations) bool {
	lc.mapMutex.RLock()
	info, ok := lc.locationUnavailabilityInfoMap[endpoint]
	lc.mapMutex.RUnlock()
	if ops == none || !ok || ops&info.unavailableOps != ops {
		return false
	}
	return time.Since(info.lastCheckTime) < lc.unavailableLocationExpirationTime
}

func (lc *locationCache) getPrefAvailableEndpoints(endpointsByLoc map[regionId]url.URL, locs []regionId, availOps requestedOperations, fallbackEndpoint url.URL) []url.URL {
	endpoints := make([]url.URL, 0)
	if lc.enableCrossRegionRetries {
		if lc.canUseMultipleWriteLocs() || availOps&read != 0 {
			unavailEndpoints := make([]url.URL, 0)
			addedFallback := false
			for _, loc := range lc.locationInfo.prefLocations {
				if endpoint, ok := endpointsByLoc[loc]; ok {
					if lc.isEndpointUnavailable(endpoint, availOps) {
						unavailEndpoints = append(unavailEndpoints, endpoint)
					} else {
						endpoints = append(endpoints, endpoint)
					}
					// Remember that we added the fallback endpoint, so we don't duplicate it at the end
					if endpoint == fallbackEndpoint {
						addedFallback = true
					}
				}
			}
			endpoints = append(endpoints, unavailEndpoints...)
			if !addedFallback {
				// If we didn't put the fallback endpoint anywhere in the list, add it to the end now
				endpoints = append(endpoints, fallbackEndpoint)
			}
		} else {
			for _, loc := range locs {
				if endpoint, ok := endpointsByLoc[loc]; ok && loc != "" {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}
	if len(endpoints) == 0 {
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
			endpointsByLoc[loc.Name] = *endpoint
			parsedLocs = append(parsedLocs, loc.Name)
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
