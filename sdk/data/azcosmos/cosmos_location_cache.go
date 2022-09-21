// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"sync"
	"time"
)

const DefaultExpirationTime time.Duration = time.Minute * 5

const (
	none opType = iota
	read
	write
)

type opType int

type locationUnavailabilityInfo struct {
	lastCheckTime time.Time
	unavailableOp opType
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

type AcctRegion struct {
	name     string
	endpoint string // make url?
	props    map[string]interface{}
}
type AcctProperties struct {
	readRegions                  []AcctRegion
	writeRegions                 []AcctRegion
	enableMultipleWriteLocations bool
}

type LocationCache struct {
	locationInfo                      dbAcctLocationsInfo
	defaultEndpt                      url.URL
	enableEndptDiscovery              bool
	useMultipleWriteLocations         bool
	connLimit                         int
	mu                                sync.Mutex // possibly try RWMutex
	locationUnavailabilityInfoMap     map[url.URL]locationUnavailabilityInfo
	lastUpdateTime                    time.Time
	enableMultipleWriteLocations      bool
	unavailableLocationExpirationTime time.Duration
}

func (lc *LocationCache) updateCore() {
}
func (lc *LocationCache) updatePrefLocations(prefLocs []string) {
}
func (lc *LocationCache) update(writeLocations []AcctRegion, readLocations []AcctRegion, prefList []string, enableMultipleWriteLocations bool) {
}

func (lc *LocationCache) ReadEndpoints() []url.URL {
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		lc.updateCore()
	}
	return lc.locationInfo.readEndpts
}

func (lc *LocationCache) WriteEndpoints() []url.URL {
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		lc.updateCore()
	}
	return lc.locationInfo.writeEndpts
}

func (lc *LocationCache) GetLocation(endpoint url.URL) string {
	firstLoc := ""
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

	if endpoint == lc.defaultEndpt && !lc.CanUseMultipleWriteLocs() {
		if len(lc.locationInfo.availWriteEndptsByLocation) > 0 {
			return firstLoc
		}
	}
	return ""

}

func (lc *LocationCache) CanUseMultipleWriteLocs() bool {
	return lc.useMultipleWriteLocations && lc.enableMultipleWriteLocations
}

func (lc *LocationCache) MarkEndptUnavailableForRead(endpoint url.URL) {
	lc.MarkEndptUnavailable(endpoint, read)
}

func (lc *LocationCache) MarkEndptUnavailableForWrite(endpoint url.URL) {
	lc.MarkEndptUnavailable(endpoint, write)
}

func (lc *LocationCache) MarkEndptUnavailable(endpoint url.URL, op opType) {
	currTime := time.Now()
	lc.mu.Lock()
	if info, ok := lc.locationUnavailabilityInfoMap[endpoint]; ok {
		info.lastCheckTime = currTime
		info.unavailableOp |= op
		lc.locationUnavailabilityInfoMap[endpoint] = info
	} else {
		info = locationUnavailabilityInfo{
			lastCheckTime: currTime,
			unavailableOp: op,
		}
		lc.locationUnavailabilityInfoMap[endpoint] = info
	}
	lc.mu.Unlock()
	lc.updateCore()
}

func (lc *LocationCache) DbAcctRead(dbAcct AcctProperties) {
	lc.update(dbAcct.writeRegions, dbAcct.readRegions, nil, dbAcct.enableMultipleWriteLocations)
}

func (lc *LocationCache) RefreshStaleEndpts() {
	lc.mu.Lock()
	for endpoint, info := range lc.locationUnavailabilityInfoMap {
		if time.Since(info.lastCheckTime) > lc.unavailableLocationExpirationTime {
			delete(lc.locationUnavailabilityInfoMap, endpoint)
		}
	}
}
