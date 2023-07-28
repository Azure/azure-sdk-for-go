// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// globalEndpointManagerPolicies contains policies related to caching and the global endpoint manager.
type globalEndpointManagerPolicies struct {
	endpointMgr *globalEndpointManager // endpointMgr is a pointer to the global endpoint manager instance.
}

// NewGlobalEndpointManagerPolicies creates a new instance of globalEndpointManagerPolicies.
func NewGlobalEndpointManagerPolicies(endpointMgr *globalEndpointManager) *globalEndpointManagerPolicies {
	return &globalEndpointManagerPolicies{
		endpointMgr: endpointMgr,
	}
}

// Do implements the policy.RequestHandler interface and allows the custom policy to be executed.
func (p *globalEndpointManagerPolicies) Do(req *policy.Request) (*http.Response, error) {
	// Trigger Update method to refresh the location cache in the global endpoint manager.
	if err := p.updateLocationCache(); err != nil {
		// Handle error if there's an issue updating the location cache.
	}

	// Continue processing the request by calling req.Next().
	// This allows other policies in the chain to execute.
	return req.Next()
}

// updateLocationCache updates the location cache using the global endpoint manager.
// This method should be called when the cache needs to be refreshed.
func (p *globalEndpointManagerPolicies) updateLocationCache() error {
	// Retrieve the latest account properties from the global endpoint manager.
	accountProps, err := p.endpointMgr.GetAccountProperties()
	if err != nil {
		return fmt.Errorf("failed to get account properties: %w", err)
	}

	// Update the location cache using the latest account properties.
	if err := p.endpointMgr.locationCache.update(accountProps.WriteRegions, accountProps.ReadRegions, p.endpointMgr.preferredRegions, &accountProps.EnableMultipleWriteLocations); err != nil {
		return fmt.Errorf("failed to update location cache: %w", err)
	}

	return nil
}
