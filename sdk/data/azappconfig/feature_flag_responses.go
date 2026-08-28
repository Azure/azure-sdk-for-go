// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AddFeatureFlagResponse contains the response from the FeatureFlagClient.AddFeatureFlag method.
type AddFeatureFlagResponse struct {
	FeatureFlag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// SetFeatureFlagResponse contains the response from the FeatureFlagClient.SetFeatureFlag method.
type SetFeatureFlagResponse struct {
	FeatureFlag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// GetFeatureFlagResponse contains the response from the FeatureFlagClient.GetFeatureFlag method.
type GetFeatureFlagResponse struct {
	FeatureFlag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// DeleteFeatureFlagResponse contains the response from the FeatureFlagClient.DeleteFeatureFlag method.
type DeleteFeatureFlagResponse struct {
	FeatureFlag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// ListFeatureFlagsPageResponse contains the response from the FeatureFlagClient.NewListFeatureFlagsPager method.
type ListFeatureFlagsPageResponse struct {
	// Contains the feature flags that match the selector provided.
	FeatureFlags []FeatureFlag

	// An ETag indicating the state of a page of feature flags within a configuration store.
	ETag *azcore.ETag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// ListFeatureFlagRevisionsPageResponse contains the response from the FeatureFlagClient.NewListFeatureFlagRevisionsPager method.
type ListFeatureFlagRevisionsPageResponse struct {
	// Contains the feature flag revisions that match the selector provided.
	FeatureFlags []FeatureFlag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}
