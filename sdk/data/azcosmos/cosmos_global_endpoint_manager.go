// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const defaultUnavailableLocationRefreshInterval = 5 * time.Minute

type globalEndpointManager struct {
	clientEndpoint      string
	pipeline            azruntime.Pipeline
	preferredLocations  []string
	locationCache       *locationCache
	refreshTimeInterval time.Duration
	gemMutex            sync.Mutex
	lastUpdateTime      time.Time
}

func newGlobalEndpointManager(clientEndpoint string, pipeline azruntime.Pipeline, preferredLocations []string, refreshTimeInterval time.Duration, enableEndpointDiscovery bool) (*globalEndpointManager, error) {
	endpoint, err := url.Parse(clientEndpoint)
	if err != nil {
		return &globalEndpointManager{}, err
	}

	if refreshTimeInterval == 0 {
		refreshTimeInterval = defaultUnavailableLocationRefreshInterval
	}

	gem := &globalEndpointManager{
		clientEndpoint:      clientEndpoint,
		pipeline:            pipeline,
		preferredLocations:  preferredLocations,
		locationCache:       newLocationCache(preferredLocations, *endpoint, enableEndpointDiscovery),
		refreshTimeInterval: refreshTimeInterval,
		lastUpdateTime:      time.Time{},
	}

	return gem, nil
}

func (gem *globalEndpointManager) GetWriteEndpoints() ([]url.URL, error) {
	return gem.locationCache.writeEndpoints()
}

func (gem *globalEndpointManager) GetReadEndpoints() ([]url.URL, error) {
	return gem.locationCache.readEndpoints()
}

func (gem *globalEndpointManager) MarkEndpointUnavailableForWrite(endpoint url.URL) error {
	return gem.locationCache.markEndpointUnavailableForWrite(endpoint)
}

func (gem *globalEndpointManager) MarkEndpointUnavailableForRead(endpoint url.URL) error {
	return gem.locationCache.markEndpointUnavailableForRead(endpoint)
}

func (gem *globalEndpointManager) GetEndpointLocation(endpoint url.URL) string {
	return gem.locationCache.getLocation(endpoint)
}

func (gem *globalEndpointManager) CanUseMultipleWriteLocations() bool {
	return gem.locationCache.canUseMultipleWriteLocs()
}

func (gem *globalEndpointManager) IsEndpointUnavailable(endpoint url.URL, ops requestedOperations) bool {
	return gem.locationCache.isEndpointUnavailable(endpoint, ops)
}

func (gem *globalEndpointManager) RefreshStaleEndpoints() {
	gem.locationCache.refreshStaleEndpoints()
}

func (gem *globalEndpointManager) ShouldRefresh() bool {
	return time.Since(gem.lastUpdateTime) > gem.refreshTimeInterval
}

func (gem *globalEndpointManager) ResolveServiceEndpoint(locationIndex int, isWriteOperation bool) url.URL {
	return gem.locationCache.resolveServiceEndpoint(locationIndex, isWriteOperation)
}

func (gem *globalEndpointManager) GetPreferredLocationEndpoint(preferredLocationIndex int, currentUrl url.URL) url.URL {
	endpointString := currentUrl.String()
	location := gem.preferredLocations[preferredLocationIndex]
	endpointParts := strings.Split(endpointString, ".")
	if len(endpointParts) > 0 {
		databaseAccountName := endpointParts[0]
		locationalDatabaseAccountName := databaseAccountName + "-" + strings.ToLower(strings.ReplaceAll(location, " ", ""))
		endpointParts[0] = locationalDatabaseAccountName
		locationalString := strings.Join(endpointParts, ".")
		locationalURL, err := url.Parse(locationalString)
		if err != nil {
			// error parsing the new url
			return currentUrl
		}
		return *locationalURL
	}
	return currentUrl
}

func (gem *globalEndpointManager) Update(ctx context.Context) error {
	gem.gemMutex.Lock()
	defer gem.gemMutex.Unlock()
	if !gem.ShouldRefresh() {
		return nil
	}
	accountProperties, err := gem.GetAccountProperties(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve account properties: %v", err)
	}
	err = gem.locationCache.update(
		accountProperties.WriteRegions,
		accountProperties.ReadRegions,
		gem.preferredLocations,
		&accountProperties.EnableMultipleWriteLocations)
	if err != nil {
		return fmt.Errorf("failed to update location cache: %v", err)
	}
	gem.lastUpdateTime = time.Now()
	return nil
}

func (gem *globalEndpointManager) GetAccountProperties(ctx context.Context) (accountProperties, error) {
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabaseAccount,
		resourceAddress: "",
	}

	ctxt, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := azruntime.NewRequest(ctxt, http.MethodGet, gem.clientEndpoint)
	if err != nil {
		return accountProperties{}, err
	}

	req.Raw().Header.Set(headerXmsDate, time.Now().UTC().Format(http.TimeFormat))
	req.Raw().Header.Set(headerXmsVersion, apiVersion)
	req.Raw().Header.Set(cosmosHeaderSDKSupportedCapabilities, supportedCapabilitiesHeaderValue)

	req.SetOperationValue(operationContext)

	azResponse, err := gem.pipeline.Do(req)
	if err != nil {
		return accountProperties{}, err
	}

	successResponse := (azResponse.StatusCode >= 200 && azResponse.StatusCode < 300)
	if successResponse {
		properties, err := newAccountProperties(azResponse)
		if err != nil {
			return accountProperties{}, fmt.Errorf("failed to parse account properties: %v", err)
		}
		return properties, nil
	}

	return accountProperties{}, newCosmosError(azResponse)
}

func newAccountProperties(azResponse *http.Response) (accountProperties, error) {
	properties := accountProperties{}
	err := azruntime.UnmarshalAsJSON(azResponse, &properties)
	if err != nil {
		return properties, err
	}

	return properties, nil
}
