// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"sync"
	"time"
)

const defaultExpirationTime time.Duration = time.Minute * 5

const (
	none opType = iota
	read
	write
	all
)

type opType int

type locationUnavailabilityInfo struct {
	lastCheckTime  time.Time
	unavailableOps opType
}

type dbAcctLocationsInfo struct {
	prefLocations              []string
	availWriteLocations        []string
	availReadLocations         []string
	availWriteEndptsByLocation map[string]url.URL
	availReadEndptsByLocation  map[string]url.URL
	writeEndpts                []url.URL
	readEndpts                 []url.URL
}

// AcctRegion represents a Azure Cosmos DB Database Account in a specific region.
type AcctRegion struct {
	// The name of the database account location in the Azure Cosmos DB service.
	name string
	// The URL of the database account location in the Azure Cosmos DB service.
	endpoint string
	// props    map[string]interface{}
}

// AcctProperties represents a container for Azure Cosmos DB databases.
type AcctProperties struct {
	// The list of readable regions of this database account from the Azure Cosmos DB service.
	readRegions []AcctRegion
	// The list of writable regions of this database account from the Azure Cosmos DB service.
	writeRegions []AcctRegion
	// Indicates whether multiple write locations are enabled for this database account.
	enableMultipleWriteLocations bool
}

type locationCache struct {
	locationInfo              dbAcctLocationsInfo
	defaultEndpt              url.URL
	enableEndptDiscovery      bool
	useMultipleWriteLocations bool
	// connLimit                         int
	rwMutex                           sync.RWMutex
	locationUnavailabilityInfoMap     map[url.URL]locationUnavailabilityInfo
	mapMutex                          sync.RWMutex
	lastUpdateTime                    time.Time
	enableMultipleWriteLocations      bool
	unavailableLocationExpirationTime time.Duration
}

func newLocationCache(prefLocations []string, defaultEndpt url.URL) *locationCache {
	return &locationCache{
		defaultEndpt:                      defaultEndpt,
		locationInfo:                      *newDbAcctLocationsInfo(prefLocations, defaultEndpt),
		locationUnavailabilityInfoMap:     make(map[url.URL]locationUnavailabilityInfo),
		unavailableLocationExpirationTime: defaultExpirationTime,
	}
}

func (lc *locationCache) update(writeLocations []AcctRegion, readLocations []AcctRegion, prefList []string, enableMultipleWriteLocations bool) error {
	lc.rwMutex.RLock()
	nextLoc := copyDbAcctLocationsInfo(lc.locationInfo)
	if prefList != nil {
		nextLoc.prefLocations = prefList
	}
	lc.enableMultipleWriteLocations = enableMultipleWriteLocations
	lc.refreshStaleEndpts()
	if readLocations != nil {
		availReadEndptsByLocation, availReadLocations, err := getEndptsByLocation(readLocations)
		if err != nil {
			lc.rwMutex.RUnlock()
			return err
		}
		nextLoc.availReadEndptsByLocation = availReadEndptsByLocation
		nextLoc.availReadLocations = availReadLocations
	}

	if writeLocations != nil {
		availWriteEndptsByLocation, availWriteLocations, err := getEndptsByLocation(writeLocations)
		if err != nil {
			lc.rwMutex.RUnlock()
			return err
		}
		nextLoc.availWriteEndptsByLocation = availWriteEndptsByLocation
		nextLoc.availWriteLocations = availWriteLocations
	}

	nextLoc.writeEndpts = lc.getPrefAvailableEndpts(nextLoc.availWriteEndptsByLocation, nextLoc.availWriteLocations, write, lc.defaultEndpt)
	nextLoc.readEndpts = lc.getPrefAvailableEndpts(nextLoc.availReadEndptsByLocation, nextLoc.availReadLocations, read, nextLoc.writeEndpts[0])
	lc.lastUpdateTime = time.Now()
	lc.rwMutex.RUnlock()
	lc.rwMutex.Lock()
	lc.locationInfo = nextLoc
	lc.rwMutex.Unlock()
	return nil
}

func (lc *locationCache) readEndpts() ([]url.URL, error) {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		err := lc.update(nil, nil, nil, lc.enableMultipleWriteLocations)
		if err != nil {
			return nil, err
		}
	}
	return lc.locationInfo.readEndpts, nil
}

func (lc *locationCache) writeEndpts() ([]url.URL, error) {
	lc.mapMutex.RLock()
	defer lc.mapMutex.RUnlock()
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		err := lc.update(nil, nil, nil, lc.enableMultipleWriteLocations)
		if err != nil {
			return nil, err
		}
	}
	return lc.locationInfo.writeEndpts, nil
}

func (lc *locationCache) getLocation(endpoint url.URL) string {
	firstLoc := ""
	lc.rwMutex.RLock()
	defer lc.rwMutex.RUnlock()
	// TODO: Find workaround for firstLoc, maps are unordered
	for location, uri := range lc.locationInfo.availWriteEndptsByLocation {
		if uri == endpoint {
			return location
		}
		if firstLoc == "" {
			firstLoc = location
		}
	}

	for location, uri := range lc.locationInfo.availReadEndptsByLocation {
		if uri == endpoint {
			return location
		}
	}

	if endpoint == lc.defaultEndpt && !lc.canUseMultipleWriteLocs() {
		if len(lc.locationInfo.availWriteEndptsByLocation) > 0 {
			return firstLoc
		}
	}
	return ""
}

func (lc *locationCache) canUseMultipleWriteLocs() bool {
	lc.rwMutex.RLock()
	defer lc.rwMutex.RUnlock()
	return lc.useMultipleWriteLocations && lc.enableMultipleWriteLocations
}

func (lc *locationCache) markEndptUnavailableForRead(endpoint url.URL) error {
	return lc.markEndptUnavailable(endpoint, read)
}

func (lc *locationCache) markEndptUnavailableForWrite(endpoint url.URL) error {
	return lc.markEndptUnavailable(endpoint, write)
}

func (lc *locationCache) markEndptUnavailable(endpoint url.URL, op opType) error {
	currTime := time.Now()
	lc.mapMutex.Lock()
	if info, ok := lc.locationUnavailabilityInfoMap[endpoint]; ok {
		info.lastCheckTime = currTime
		info.unavailableOps |= op
		lc.locationUnavailabilityInfoMap[endpoint] = info
	} else {
		info = locationUnavailabilityInfo{
			lastCheckTime:  currTime,
			unavailableOps: op,
		}
		lc.locationUnavailabilityInfoMap[endpoint] = info
	}
	lc.mapMutex.Unlock()
	err := lc.update(nil, nil, nil, lc.enableMultipleWriteLocations)
	return err
}

func (lc *locationCache) dbAcctRead(dbAcct AcctProperties) error {
	return lc.update(dbAcct.writeRegions, dbAcct.readRegions, nil, dbAcct.enableMultipleWriteLocations)
}

func (lc *locationCache) refreshStaleEndpts() {
	lc.mapMutex.Lock()
	for endpoint, info := range lc.locationUnavailabilityInfoMap {
		t := time.Since(info.lastCheckTime)
		if t > lc.unavailableLocationExpirationTime {
			delete(lc.locationUnavailabilityInfoMap, endpoint)
		}
	}
	lc.mapMutex.Unlock()
}

func (lc *locationCache) isEndptUnavailable(endpoint url.URL, ops opType) bool {
	lc.mapMutex.RLock()
	info, ok := lc.locationUnavailabilityInfoMap[endpoint]
	lc.mapMutex.RUnlock()
	if ops == none || !ok || ops&info.unavailableOps != ops {
		return false
	}
	lc.rwMutex.RLock()
	defer lc.rwMutex.RUnlock()
	return time.Since(info.lastCheckTime) < lc.unavailableLocationExpirationTime
}

func (lc *locationCache) getPrefAvailableEndpts(endptsByLoc map[string]url.URL, locs []string, availOps opType, fallbackEndpt url.URL) []url.URL {
	endpts := make([]url.URL, 0)
	lc.rwMutex.RLock()
	if lc.enableEndptDiscovery {
		if lc.canUseMultipleWriteLocs() || availOps&read != 0 {
			unavailEndpts := make([]url.URL, 0)
			unavailEndpts = append(unavailEndpts, fallbackEndpt)
			for _, loc := range lc.locationInfo.prefLocations {
				if endpt, ok := endptsByLoc[loc]; ok && endpt != fallbackEndpt {
					if lc.isEndptUnavailable(endpt, availOps) {
						unavailEndpts = append(unavailEndpts, endpt)
					} else {
						endpts = append(endpts, endpt)
					}
				}
			}
			endpts = append(endpts, unavailEndpts...)
		} else {
			for _, loc := range locs {
				if endpt, ok := endptsByLoc[loc]; ok && loc != "" {
					endpts = append(endpts, endpt)
				}
			}
		}
	}
	lc.rwMutex.RUnlock()
	if len(endpts) == 0 {
		endpts = append(endpts, fallbackEndpt)
	}
	return endpts
}

func getEndptsByLocation(locs []AcctRegion) (map[string]url.URL, []string, error) {
	endptsByLoc := make(map[string]url.URL)
	parsedLocs := make([]string, 0)
	for _, loc := range locs {
		endpt, err := url.Parse(loc.endpoint)
		if err != nil {
			return nil, nil, err
		}
		if loc.name != "" {
			endptsByLoc[loc.name] = *endpt
			parsedLocs = append(parsedLocs, loc.name)
			// set service pt needed?
		}
	}
	return endptsByLoc, parsedLocs, nil
}

func newDbAcctLocationsInfo(prefLocations []string, defaultEndpt url.URL) *dbAcctLocationsInfo {
	availWriteLocs := make([]string, 0)
	availReadLocs := make([]string, 0)
	availWriteEndptsByLocation := make(map[string]url.URL)
	availReadEndptsByLocation := make(map[string]url.URL)
	writeEndpts := []url.URL{defaultEndpt}
	readEndpts := []url.URL{defaultEndpt}
	return &dbAcctLocationsInfo{
		prefLocations:              prefLocations,
		availWriteLocations:        availWriteLocs,
		availReadLocations:         availReadLocs,
		availWriteEndptsByLocation: availWriteEndptsByLocation,
		availReadEndptsByLocation:  availReadEndptsByLocation,
		writeEndpts:                writeEndpts,
		readEndpts:                 readEndpts,
	}
}

func copyDbAcctLocationsInfo(other dbAcctLocationsInfo) dbAcctLocationsInfo {
	return dbAcctLocationsInfo{
		prefLocations:              other.prefLocations,
		availWriteLocations:        other.availWriteLocations,
		availReadLocations:         other.availReadLocations,
		availWriteEndptsByLocation: other.availWriteEndptsByLocation,
		availReadEndptsByLocation:  other.availReadEndptsByLocation,
		writeEndpts:                other.writeEndpts,
		readEndpts:                 other.readEndpts,
	}
}
