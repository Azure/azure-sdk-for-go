// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type DatabaseAccountLocation struct {
	Name                    string `json:"name"`
	DatabaseAccountEndpoint string `json:"databaseAccountEndpoint"`
}

type DatabaseAccount struct {
	ID                           string                    `json:"id"`
	ResourceID                   string                    `json:"_rid"`
	WritableLocations            []DatabaseAccountLocation `json:"writableLocations"`
	ReadableLocations            []DatabaseAccountLocation `json:"readableLocations"`
	EnableMultipleWriteLocations bool                      `json:"enableMultipleWriteLocations"`
}

func NewDatabaseAccountFromJSON(jsonStr string) (*DatabaseAccount, error) {
	var account DatabaseAccount
	if err := json.Unmarshal([]byte(jsonStr), &account); err != nil {
		return nil, err
	}
	return &account, nil
}

type DatabaseAccountManagerInternal interface {
	GetDatabaseAccountFromEndpoint(ctx context.Context, endpoint *url.URL) (*DatabaseAccount, error)
	GetServiceEndpoint() *url.URL
}

type ConnectionPolicy struct {
	EnableEndpointDiscovery     bool
	PreferredLocations          []string
	UsingMultipleWriteLocations bool
}

type DatabaseAccountLocationsInfo struct {
	PreferredLocations               []string
	AvailableWriteEndpointByLocation map[string]*url.URL
	AvailableReadEndpointByLocation  map[string]*url.URL
	AvailableWriteLocations          []string
	AvailableReadLocations           []string
	WriteEndpoints                   []*url.URL
	ReadEndpoints                    []*url.URL
}

func newDatabaseAccountLocationsInfo(preferredLocations []string) *DatabaseAccountLocationsInfo {
	return &DatabaseAccountLocationsInfo{
		PreferredLocations:               preferredLocations,
		AvailableWriteEndpointByLocation: make(map[string]*url.URL),
		AvailableReadEndpointByLocation:  make(map[string]*url.URL),
		AvailableWriteLocations:          []string{},
		AvailableReadLocations:           []string{},
		WriteEndpoints:                   []*url.URL{},
		ReadEndpoints:                    []*url.URL{},
	}
}

type LocationCache struct {
	mu                                   sync.RWMutex
	locationInfo                         *DatabaseAccountLocationsInfo
	connectionPolicy                     *ConnectionPolicy
	defaultEndpoint                      *url.URL
	locationUnavailabilityInfoByEndpoint map[string]time.Time
	enableMultipleWriteLocations         bool
}

func NewLocationCache(connectionPolicy *ConnectionPolicy, defaultEndpoint *url.URL) *LocationCache {
	lc := &LocationCache{
		locationInfo:                         newDatabaseAccountLocationsInfo(connectionPolicy.PreferredLocations),
		connectionPolicy:                     connectionPolicy,
		defaultEndpoint:                      defaultEndpoint,
		locationUnavailabilityInfoByEndpoint: make(map[string]time.Time),
	}
	return lc
}

func (lc *LocationCache) GetReadEndpoints() []*url.URL {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	if len(lc.locationInfo.ReadEndpoints) == 0 {
		return []*url.URL{lc.defaultEndpoint}
	}
	return lc.locationInfo.ReadEndpoints
}

func (lc *LocationCache) GetWriteEndpoints() []*url.URL {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	if len(lc.locationInfo.WriteEndpoints) == 0 {
		return []*url.URL{lc.defaultEndpoint}
	}
	return lc.locationInfo.WriteEndpoints
}

func (lc *LocationCache) GetAvailableReadEndpointByLocation() map[string]*url.URL {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	result := make(map[string]*url.URL)
	for k, v := range lc.locationInfo.AvailableReadEndpointByLocation {
		result[k] = v
	}
	return result
}

func (lc *LocationCache) GetAvailableWriteEndpointByLocation() map[string]*url.URL {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	result := make(map[string]*url.URL)
	for k, v := range lc.locationInfo.AvailableWriteEndpointByLocation {
		result[k] = v
	}
	return result
}

func (lc *LocationCache) MarkEndpointUnavailableForRead(endpoint *url.URL) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.locationUnavailabilityInfoByEndpoint[endpoint.String()] = time.Now()
	lc.updateReadEndpoints()
}

func (lc *LocationCache) MarkEndpointUnavailableForWrite(endpoint *url.URL) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.locationUnavailabilityInfoByEndpoint[endpoint.String()] = time.Now()
	lc.updateWriteEndpoints()
}

func (lc *LocationCache) OnDatabaseAccountRead(databaseAccount *DatabaseAccount) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.enableMultipleWriteLocations = databaseAccount.EnableMultipleWriteLocations

	newLocationInfo := newDatabaseAccountLocationsInfo(lc.connectionPolicy.PreferredLocations)

	for _, loc := range databaseAccount.WritableLocations {
		endpoint, err := url.Parse(loc.DatabaseAccountEndpoint)
		if err == nil {
			newLocationInfo.AvailableWriteEndpointByLocation[loc.Name] = endpoint
			newLocationInfo.AvailableWriteLocations = append(newLocationInfo.AvailableWriteLocations, loc.Name)
		}
	}

	for _, loc := range databaseAccount.ReadableLocations {
		endpoint, err := url.Parse(loc.DatabaseAccountEndpoint)
		if err == nil {
			newLocationInfo.AvailableReadEndpointByLocation[loc.Name] = endpoint
			newLocationInfo.AvailableReadLocations = append(newLocationInfo.AvailableReadLocations, loc.Name)
		}
	}

	lc.locationInfo = newLocationInfo
	lc.updateReadEndpoints()
	lc.updateWriteEndpoints()
}

func (lc *LocationCache) updateReadEndpoints() {
	var endpoints []*url.URL
	var unavailablePreferredEndpoints []*url.URL

	unavailableEndpointMap := make(map[string]bool)
	for endpoint := range lc.locationUnavailabilityInfoByEndpoint {
		unavailableEndpointMap[endpoint] = true
	}

	for _, location := range lc.locationInfo.PreferredLocations {
		if endpoint, ok := lc.locationInfo.AvailableReadEndpointByLocation[location]; ok {
			if unavailableEndpointMap[endpoint.String()] {
				unavailablePreferredEndpoints = append(unavailablePreferredEndpoints, endpoint)
			} else {
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	if len(endpoints) == 0 {
		// Use first write endpoint as fallback (matches Java behavior)
		// This ensures that when write endpoint is unavailable for read,
		// shouldRefreshEndpoints returns canRefreshInBackground=false
		fallback := lc.defaultEndpoint
		if len(lc.locationInfo.WriteEndpoints) > 0 {
			fallback = lc.locationInfo.WriteEndpoints[0]
		}
		endpoints = append(endpoints, fallback)
	}

	endpoints = append(endpoints, unavailablePreferredEndpoints...)

	lc.locationInfo.ReadEndpoints = endpoints
}

func (lc *LocationCache) updateWriteEndpoints() {
	var endpoints []*url.URL
	var unavailablePreferredEndpoints []*url.URL

	unavailableEndpointMap := make(map[string]bool)
	for endpoint := range lc.locationUnavailabilityInfoByEndpoint {
		unavailableEndpointMap[endpoint] = true
	}

	if lc.CanUseMultipleWriteLocations() {
		for _, location := range lc.locationInfo.PreferredLocations {
			if endpoint, ok := lc.locationInfo.AvailableWriteEndpointByLocation[location]; ok {
				if unavailableEndpointMap[endpoint.String()] {
					unavailablePreferredEndpoints = append(unavailablePreferredEndpoints, endpoint)
				} else {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	} else {
		for _, location := range lc.locationInfo.AvailableWriteLocations {
			if endpoint, ok := lc.locationInfo.AvailableWriteEndpointByLocation[location]; ok {
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	if len(endpoints) == 0 {
		endpoints = append(endpoints, lc.defaultEndpoint)
	}

	endpoints = append(endpoints, unavailablePreferredEndpoints...)

	lc.locationInfo.WriteEndpoints = endpoints
}

func (lc *LocationCache) isEndpointUnavailable(endpoint *url.URL, forRead bool) bool {
	_, exists := lc.locationUnavailabilityInfoByEndpoint[endpoint.String()]
	return exists
}

func (lc *LocationCache) anyEndpointsAvailable(endpoints []*url.URL) bool {
	for _, endpoint := range endpoints {
		if _, unavailable := lc.locationUnavailabilityInfoByEndpoint[endpoint.String()]; !unavailable {
			return true
		}
	}
	return false
}

func (lc *LocationCache) shouldRefreshEndpoints() (shouldRefresh bool, canRefreshInBackground bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	canRefreshInBackground = true

	if !lc.connectionPolicy.EnableEndpointDiscovery {
		return false, true
	}

	readEndpoints := lc.locationInfo.ReadEndpoints
	if len(readEndpoints) == 0 {
		return true, false
	}

	if lc.isEndpointUnavailable(readEndpoints[0], true) {
		canRefreshInBackground = lc.anyEndpointsAvailable(readEndpoints)
		return true, canRefreshInBackground
	}

	var mostPreferredLocation string
	if len(lc.locationInfo.PreferredLocations) > 0 {
		mostPreferredLocation = lc.locationInfo.PreferredLocations[0]
	}

	if mostPreferredLocation != "" {
		if mostPreferredEndpoint, ok := lc.locationInfo.AvailableReadEndpointByLocation[mostPreferredLocation]; ok {
			if mostPreferredEndpoint.String() != readEndpoints[0].String() {
				return true, true
			}
		} else {
			return true, true
		}
	}

	writeEndpoints := lc.locationInfo.WriteEndpoints
	if len(writeEndpoints) == 0 {
		return true, false
	}

	if !lc.CanUseMultipleWriteLocations() {
		if lc.isEndpointUnavailable(writeEndpoints[0], false) {
			canRefreshInBackground = lc.anyEndpointsAvailable(writeEndpoints)
			return true, canRefreshInBackground
		}
		shouldRefresh = lc.connectionPolicy.UsingMultipleWriteLocations && !lc.enableMultipleWriteLocations
		return shouldRefresh, true
	}

	if mostPreferredLocation != "" {
		if mostPreferredEndpoint, ok := lc.locationInfo.AvailableWriteEndpointByLocation[mostPreferredLocation]; ok {
			if mostPreferredEndpoint.String() != writeEndpoints[0].String() {
				return true, true
			}
		} else {
			return true, true
		}
	}

	shouldRefresh = lc.connectionPolicy.UsingMultipleWriteLocations && !lc.enableMultipleWriteLocations
	return shouldRefresh, true
}

func (lc *LocationCache) CanUseMultipleWriteLocations() bool {
	return lc.connectionPolicy.UsingMultipleWriteLocations && lc.enableMultipleWriteLocations
}

type GlobalEndpointManager struct {
	databaseAccountManager                    DatabaseAccountManagerInternal
	connectionPolicy                          *ConnectionPolicy
	locationCache                             *LocationCache
	isRefreshing                              atomic.Bool
	refreshInBackground                       atomic.Bool
	backgroundRefreshLocationTimeIntervalInMS int
	lastCacheUpdateTime                       time.Time
	mu                                        sync.Mutex
	stopChan                                  chan struct{}
}

func NewGlobalEndpointManager(
	databaseAccountManager DatabaseAccountManagerInternal,
	connectionPolicy *ConnectionPolicy,
) *GlobalEndpointManager {
	serviceEndpoint := databaseAccountManager.GetServiceEndpoint()

	gem := &GlobalEndpointManager{
		databaseAccountManager:                    databaseAccountManager,
		connectionPolicy:                          connectionPolicy,
		locationCache:                             NewLocationCache(connectionPolicy, serviceEndpoint),
		backgroundRefreshLocationTimeIntervalInMS: 5 * 60 * 1000,
		stopChan: make(chan struct{}),
	}

	return gem
}

func (gem *GlobalEndpointManager) Init(ctx context.Context) error {
	return gem.refreshLocations(ctx)
}

func (gem *GlobalEndpointManager) GetLocationCache() *LocationCache {
	return gem.locationCache
}

func (gem *GlobalEndpointManager) IsRefreshing() bool {
	return gem.isRefreshing.Load()
}

func (gem *GlobalEndpointManager) RefreshInBackground() bool {
	return gem.refreshInBackground.Load()
}

func (gem *GlobalEndpointManager) SetBackgroundRefreshLocationTimeIntervalInMS(ms int) {
	gem.mu.Lock()
	defer gem.mu.Unlock()
	gem.backgroundRefreshLocationTimeIntervalInMS = ms
}

func (gem *GlobalEndpointManager) MarkEndpointUnavailableForRead(endpoint *url.URL) {
	gem.locationCache.MarkEndpointUnavailableForRead(endpoint)
}

func (gem *GlobalEndpointManager) MarkEndpointUnavailableForWrite(endpoint *url.URL) {
	gem.locationCache.MarkEndpointUnavailableForWrite(endpoint)
}

func (gem *GlobalEndpointManager) RefreshLocationAsync(ctx context.Context, forceRefresh bool) error {
	if !gem.connectionPolicy.EnableEndpointDiscovery {
		return nil
	}

	if forceRefresh {
		err := gem.refreshLocations(ctx)
		if err != nil {
			return err
		}
		gem.startBackgroundRefreshIfNeeded()
		return nil
	}

	if !gem.isRefreshing.CompareAndSwap(false, true) {
		return nil
	}
	defer gem.isRefreshing.Store(false)

	shouldRefresh, canRefreshInBackground := gem.locationCache.shouldRefreshEndpoints()

	if shouldRefresh {
		if !canRefreshInBackground {
			err := gem.refreshLocations(ctx)
			if err != nil {
				return err
			}
		}
		gem.startBackgroundRefreshIfNeeded()
	}

	return nil
}

func (gem *GlobalEndpointManager) refreshLocations(ctx context.Context) error {
	serviceEndpoint := gem.databaseAccountManager.GetServiceEndpoint()

	databaseAccount, err := gem.databaseAccountManager.GetDatabaseAccountFromEndpoint(ctx, serviceEndpoint)
	if err != nil {
		return err
	}

	gem.locationCache.OnDatabaseAccountRead(databaseAccount)
	gem.lastCacheUpdateTime = time.Now()

	return nil
}

func (gem *GlobalEndpointManager) startBackgroundRefreshIfNeeded() {
	if gem.locationCache.CanUseMultipleWriteLocations() {
		gem.refreshInBackground.Store(false)
		return
	}

	gem.refreshInBackground.Store(true)
}

func (gem *GlobalEndpointManager) StartRefreshLocationTimerAsync(ctx context.Context) {
	go func() {
		gem.mu.Lock()
		interval := gem.backgroundRefreshLocationTimeIntervalInMS
		gem.mu.Unlock()

		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-gem.stopChan:
				return
			case <-ticker.C:
				gem.backgroundRefresh(ctx)
			}
		}
	}()
}

func (gem *GlobalEndpointManager) backgroundRefresh(ctx context.Context) {
	gem.refreshLocations(ctx)
	gem.refreshInBackground.Store(false)

	shouldRefresh, canRefreshInBackground := gem.locationCache.shouldRefreshEndpoints()
	if shouldRefresh && !canRefreshInBackground {
		gem.refreshLocations(ctx)
	}
}

func (gem *GlobalEndpointManager) Close() {
	close(gem.stopChan)
}
