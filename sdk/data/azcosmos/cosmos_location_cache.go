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

type locationCache struct {
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

