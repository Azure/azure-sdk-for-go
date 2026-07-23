// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AddFeatureFlagOptions contains the optional parameters for the FeatureFlagClient.AddFeatureFlag method.
type AddFeatureFlagOptions struct {
	// placeholder for future options
}

// SetFeatureFlagOptions contains the optional parameters for the FeatureFlagClient.SetFeatureFlag method.
type SetFeatureFlagOptions struct {
	// If set, and the feature flag exists in the configuration store, overwrite the flag
	// only if the passed-in ETag is the same as the flag's ETag in the configuration store.
	//
	// This has IfMatch semantics.
	OnlyIfUnchanged *azcore.ETag
}

// GetFeatureFlagOptions contains the optional parameters for the FeatureFlagClient.GetFeatureFlag method.
type GetFeatureFlagOptions struct {
	// The feature flag will be retrieved exactly as it existed at the provided time.
	AcceptDateTime *time.Time

	// Feature flag label.
	Label *string

	// If set, only retrieve the feature flag from the configuration store if the flag has changed
	// since the client last retrieved it with the ETag provided.
	//
	// This has IfNoneMatch semantics.
	OnlyIfChanged *azcore.ETag
}

// DeleteFeatureFlagOptions contains the optional parameters for the FeatureFlagClient.DeleteFeatureFlag method.
type DeleteFeatureFlagOptions struct {
	// Feature flag label.
	Label *string

	// If set, and the feature flag exists in the configuration store,
	// delete the flag only if the passed-in ETag is the same as the flag's ETag in the configuration store.
	//
	// This has IfMatch semantics.
	OnlyIfUnchanged *azcore.ETag
}

// ListFeatureFlagsOptions contains the optional parameters for the FeatureFlagClient.NewListFeatureFlagsPager method.
type ListFeatureFlagsOptions struct {
	// The match conditions used when making the request.
	// Conditions are applied to pages one by one in the order specified.
	MatchConditions []azcore.MatchConditions
}

// ListFeatureFlagRevisionsOptions contains the optional parameters for the FeatureFlagClient.NewListFeatureFlagRevisionsPager method.
type ListFeatureFlagRevisionsOptions struct {
	// placeholder for future options
}
