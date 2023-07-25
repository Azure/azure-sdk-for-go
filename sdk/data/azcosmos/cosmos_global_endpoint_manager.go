package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const defaultBackgroundRefreshInterval = 5 * time.Minute

type globalEndpointManager struct {
	client           *Client
	preferredRegions []string
	locationCache    *locationCache
}

// newGlobalEndpointManager creates a new global endpoint manager.
// It takes a client, preferredRegions, and refreshInterval as input parameters.
// If the refreshInterval is zero, it uses the default value (5 minutes).
// It initializes a new globalEndpointManager, creates a location cache, and starts a background refresh process.
// The background refresh process periodically updates the location cache based on the specified refreshInterval.
func newGlobalEndpointManager(client *Client, preferredRegions []string, refreshInterval time.Duration) (*globalEndpointManager, error) {
	endpoint, err := url.Parse(client.endpoint)
	if err != nil {
		return nil, err
	}

	// If refreshInterval is zero, use the default value (5 minutes)
	if refreshInterval == 0 {
		refreshInterval = defaultBackgroundRefreshInterval
	}

	gem, err := &globalEndpointManager{
		client:           client,
		preferredRegions: preferredRegions,
		locationCache:    newLocationCache(preferredRegions, *endpoint),
	}, nil
	if err != nil {
		return nil, err
	}

	// Start the background refresh process with the specified refreshInterval.
	gem.startBackgroundRefresh(refreshInterval)

	return gem, nil
}

// GetWriteEndpoints returns the write endpoints from the location cache.
func (gem *globalEndpointManager) GetWriteEndpoints() ([]url.URL, error) {
	return gem.locationCache.writeEndpoints()
}

// GetReadEndpoints returns the read endpoints from the location cache.
func (gem *globalEndpointManager) GetReadEndpoints() ([]url.URL, error) {
	return gem.locationCache.readEndpoints()
}

// GetAccountProperties retrieves account properties from the Cosmos DB instance.
func (gem *globalEndpointManager) GetAccountProperties() (accountProperties, error) {
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabaseAccount,
		resourceAddress: "",
	}

	path, err := generatePathForNameBased(resourceTypeDatabaseAccount, "", false)
	if err != nil {
		return accountProperties{}, fmt.Errorf("failed to generate path for name-based request: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	azResponse, err := gem.client.sendGetRequest(path, ctx, operationContext, nil, nil)
	cancel()
	if err != nil {
		return accountProperties{}, fmt.Errorf("failed to send GET request: %v", err)
	}

	properties, err := newAccountProperties(azResponse)
	if err != nil {
		return accountProperties{}, fmt.Errorf("failed to parse account properties: %v", err)
	}

	return properties, nil
}

// newAccountProperties parses the account properties from the HTTP response.
func newAccountProperties(azResponse *http.Response) (accountProperties, error) {
	properties := accountProperties{}
	err := azruntime.UnmarshalAsJSON(azResponse, &properties)
	if err != nil {
		return properties, err
	}
	return properties, nil
}

// GetLocation returns the location for the given endpoint from the location cache.
func (gem *globalEndpointManager) GetLocation(endpoint url.URL) string {
	return gem.locationCache.getLocation(endpoint)
}

// MarkEndpointUnavailableForRead marks an endpoint as unavailable for read operations.
func (gem *globalEndpointManager) MarkEndpointUnavailableForRead(endpoint url.URL) error {
	return gem.locationCache.markEndpointUnavailableForRead(endpoint)
}

// MarkEndpointUnavailableForWrite marks an endpoint as unavailable for write operations.
func (gem *globalEndpointManager) MarkEndpointUnavailableForWrite(endpoint url.URL) error {
	return gem.locationCache.markEndpointUnavailableForWrite(endpoint)
}

// Update updates the location cache and the client's default endpoint based on the provided account properties.
func (gem *globalEndpointManager) Update() error {
	accountProperties, err := gem.GetAccountProperties()
	if err != nil {
		return fmt.Errorf("failed to retrieve account properties: %v", err)
	}

	if err := gem.locationCache.update(accountProperties.WriteRegions, accountProperties.ReadRegions, gem.preferredRegions, &accountProperties.EnableMultipleWriteLocations); err != nil {
		return fmt.Errorf("failed to update location cache: %v", err)
	}
	return nil
}

// startBackgroundRefresh starts a background refresh process that periodically updates the location cache
// based on the given refreshInterval.
// It creates a context and a cancellation function to gracefully stop the background refresh process.
// The function runs a goroutine that listens for two events:
//  1. If the context is canceled (e.g., due to stopBackgroundRefresh being called), it performs an update
//     to refresh the location cache and cancels the goroutine.
//  2. After the specified refreshInterval, it performs an update to refresh the location cache and cancels the goroutine.
//
// The method is designed to be used as a goroutine and will continue running until explicitly stopped.
func (gem *globalEndpointManager) startBackgroundRefresh(refreshInterval time.Duration) {
	// Create an error channel to receive possible errors from the goroutine.
	errChan := make(chan error)

	// Create a new context and a cancellation function to stop the background refresh.
	ctx, cancel := context.WithCancel(context.Background())

	// Start a new goroutine to handle the background refresh process.
	go func() {
		for {
			select {
			// If the context is canceled (e.g., due to stopBackgroundRefresh being called), perform an update
			// to refresh the location cache and cancel the goroutine.
			case <-ctx.Done():
				_ = gem.Update() // Perform the update to refresh the location cache.
				cancel()         // Cancel the context to stop the goroutine.
				return           // Return from the goroutine.
			// After the specified refreshInterval, perform an update to refresh the location cache and cancel the goroutine.
			case <-time.After(refreshInterval):
				_ = gem.Update() // Perform the update to refresh the location cache.
				return           // Return from the goroutine.
			}
		}
	}()

	// Wait for an error to be received from the error channel (this will not happen in this case).
	// If there's an error, panic (this should never happen).
	if err := <-errChan; err != nil {
		panic(err)
	}
}

// RefreshStaleEndpoints triggers a refresh of stale endpoints in the location cache.
func (gem *globalEndpointManager) RefreshStaleEndpoints() {
	gem.locationCache.refreshStaleEndpoints()
}

// IsEndpointUnavailable checks if an endpoint is marked as unavailable for the given requested operations.
func (gem *globalEndpointManager) IsEndpointUnavailable(endpoint url.URL, ops requestedOperations) bool {
	return gem.locationCache.isEndpointUnavailable(endpoint, ops)
}

// CanUseMultipleWriteLocations returns whether multiple write locations can be used.
func (gem *globalEndpointManager) CanUseMultipleWriteLocations() bool {
	return gem.locationCache.canUseMultipleWriteLocs()
}
