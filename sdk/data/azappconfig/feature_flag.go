// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
)

// FeatureFlag represents an Azure App Configuration feature flag.
//
// Feature flags are managed through the feature flag endpoint. Use [FeatureFlagClient]
// to interact with them.
type FeatureFlag struct {
	// The enabled state of the feature flag.
	Enabled *bool

	// The name of the feature flag. Required when passing a FeatureFlag to a
	// FeatureFlagClient operation.
	Name *string

	// The allocation of the feature flag.
	Allocation *FeatureFlagAllocation

	// The conditions that must be met for the feature flag to be enabled.
	Conditions *FeatureFlagConditions

	// A description of the feature flag.
	Description *string

	// A label used to group feature flags.
	Label *string

	// An ETag indicating the state of a feature flag within a configuration store.
	ETag *azcore.ETag

	// A dictionary of tags used to assign additional properties to a feature flag.
	Tags map[string]*string

	// The telemetry settings of the feature flag.
	Telemetry *FeatureFlagTelemetryConfiguration

	// The variants of the feature flag.
	Variants []FeatureFlagVariantDefinition

	// The last time a modifying operation was performed on the given feature flag.
	LastModified *time.Time
}

// FeatureFlagAllocation defines how to allocate variants based on context.
type FeatureFlagAllocation = generated.FeatureFlagAllocation

// FeatureFlagConditions describes the conditions that must be met for a feature flag to be enabled.
type FeatureFlagConditions = generated.FeatureFlagConditions

// FeatureFlagFilter is a filter that will conditionally enable or disable a feature flag.
type FeatureFlagFilter = generated.FeatureFlagFilter

// FeatureFlagTelemetryConfiguration describes the telemetry settings for a feature flag.
type FeatureFlagTelemetryConfiguration = generated.FeatureFlagTelemetryConfiguration

// FeatureFlagVariantDefinition describes a variant of a feature flag.
type FeatureFlagVariantDefinition = generated.FeatureFlagVariantDefinition

// GroupAllocation allocates groups to a feature flag variant.
type GroupAllocation = generated.GroupAllocation

// PercentileAllocation allocates a percentile range to a feature flag variant.
type PercentileAllocation = generated.PercentileAllocation

// UserAllocation allocates users to a feature flag variant.
type UserAllocation = generated.UserAllocation

func featureFlagFromGenerated(ff generated.FeatureFlag) FeatureFlag {
	return FeatureFlag{
		Enabled:      ff.Enabled,
		Name:         ff.Name,
		Allocation:   ff.Allocation,
		Conditions:   ff.Conditions,
		Description:  ff.Description,
		Label:        ff.Label,
		ETag:         (*azcore.ETag)(ff.Etag),
		Tags:         ff.Tags,
		Telemetry:    ff.Telemetry,
		Variants:     ff.Variants,
		LastModified: ff.LastModified,
	}
}

func (ff FeatureFlag) toGenerated() generated.FeatureFlag {
	return generated.FeatureFlag{
		Enabled:      ff.Enabled,
		Name:         ff.Name,
		Allocation:   ff.Allocation,
		Conditions:   ff.Conditions,
		Description:  ff.Description,
		Label:        ff.Label,
		Etag:         (*string)(ff.ETag),
		Tags:         ff.Tags,
		Telemetry:    ff.Telemetry,
		Variants:     ff.Variants,
		LastModified: ff.LastModified,
	}
}
