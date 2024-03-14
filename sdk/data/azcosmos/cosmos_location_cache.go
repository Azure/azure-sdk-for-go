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
	UserConsistencyPolicy        userConsistencyPolicy `json:"userConsistencyPolicy"`
	ContinuousBackupEnabled      bool                  `json:"continousBackupEnabled"`
}

func (accountProps accountProperties) String() string {
	return fmt.Sprintf("Read Regions: %v\nWrite Regions: %v\nUse Multiple Write Locations: %v\nDefault Consistency Level: %v\nContinous Backup Enabled: %v\n",
		accountProps.ReadRegions, accountProps.WriteRegions, accountProps.EnableMultipleWriteLocations, accountProps.UserConsistencyPolicy.DefaultConsistencyLevel, accountProps.ContinuousBackupEnabled)
}

type locationCache struct {
	locationInfo                      databaseAccountLocationsInfo
	defaultEndpoint                   url.URL
	enableEndpointDiscovery           bool
	useMultipleWriteLocations         bool
	locationUnavailabilityInfoMap     map[url.URL]locationUnavailabilityInfo
	mapMutex                          sync.RWMutex
	lastUpdateTime                    time.Time
	enableMultipleWriteLocations      bool
	unavailableLocationExpirationTime time.Duration
}

func newLocationCache(prefLocations []string, defaultEndpoint url.URL) *locationCache {
	return &locationCache{
		defaultEndpoint:                   defaultEndpoint,
		locationInfo:                      *newDatabaseAccountLocationsInfo(prefLocations, defaultEndpoint),
		locationUnavailabilityInfoMap:     make(map[url.URL]locationUnavailabilityInfo),
		unavailableLocationExpirationTime: defaultExpirationTime,
	}
}

func (lc *locationCache) update(writeLocations []accountRegion, readLocations []accountRegion, prefList []string, enableMultipleWriteLocations *bool) error {
	nextLoc := copyDatabaseAccountLocationsInfo(lc.locationInfo)
	if prefList != nil {
		nextLoc.prefLocations = prefList
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
	lc.lastUpdateTime = time.Now()
	lc.locationInfo = nextLoc
	// TODO: log
	return nil
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

func (lc *locationCache) canUseMultipleWriteLocs() bool {
	return lc.useMultipleWriteLocations && lc.enableMultipleWriteLocations
}

func (lc *locationCache) markEndpointUnavailableForRead(endpoint url.URL) error {
	return lc.markEndpointUnavailable(endpoint, read)
}

func (lc *locationCache) markEndpointUnavailableForWrite(endpoint url.URL) error {
	return lc.markEndpointUnavailable(endpoint, write)
}

func (lc *locationCache) markEndpointUnavailable(endpoint url.URL, op requestedOperations) error {
	now := time.Now()
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
			delete(lc.locationUnavailabilityInfoMap, endpoint)
		}
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

func (lc *locationCache) getPrefAvailableEndpoints(endpointsByLoc map[string]url.URL, locs []string, availOps requestedOperations, fallbackEndpoint url.URL) []url.URL {
	endpoints := make([]url.URL, 0)
	if lc.enableEndpointDiscovery {
		if lc.canUseMultipleWriteLocs() || availOps&read != 0 {
			unavailEndpoints := make([]url.URL, 0)
			unavailEndpoints = append(unavailEndpoints, fallbackEndpoint)
			for _, loc := range lc.locationInfo.prefLocations {
				if endpoint, ok := endpointsByLoc[loc]; ok && endpoint != fallbackEndpoint {
					if lc.isEndpointUnavailable(endpoint, availOps) {
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
