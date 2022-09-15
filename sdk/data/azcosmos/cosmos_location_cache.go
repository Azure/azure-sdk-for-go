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

func (lc *LocationCache) update(writeLocations []string, readLocations []string, prefList []string, enableMultipleWriteLocations bool) {

}

func (lc *LocationCache) ReadEndpoints() []url.URL {
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		// lc.Update()
	}
	return lc.locationInfo.readEndpts
}

func (lc *LocationCache) WriteEndpoints() []url.URL {
	if time.Since(lc.lastUpdateTime) > lc.unavailableLocationExpirationTime && len(lc.locationUnavailabilityInfoMap) > 0 {
		// lc.update()
	}
	return lc.locationInfo.writeEndpts
}

func (lc *LocationCache) GetLocation(uri url.URL) string {
	// loc := lc.locationInfo.availWriteEndptsByLocation
	// return loc
}
