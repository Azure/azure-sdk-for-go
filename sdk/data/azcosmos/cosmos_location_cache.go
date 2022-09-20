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
	return
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
