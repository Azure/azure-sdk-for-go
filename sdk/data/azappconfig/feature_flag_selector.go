// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
)

// FeatureFlagSelector is a set of options that allows selecting a filtered set of feature flags
// from the configuration store, and optionally allows indicating which fields of each flag to retrieve.
type FeatureFlagSelector struct {
	// Name filter that will be used to select a set of feature flags.
	NameFilter *string

	// Label filter that will be used to select a set of feature flags.
	LabelFilter *string

	// Tags filter that will be used to select a set of feature flags.
	// This is a list of tag filters in the format {tagName=tagValue}. For more information about filtering by tags, see:
	// https://aka.ms/azconfig/docs/keyvaluefiltering
	TagsFilter []string

	// Indicates the point in time in the revision history of the selected feature flags to retrieve.
	// If set, all properties of the feature flags in the returned group will be exactly what they were at this time.
	AcceptDateTime *time.Time

	// The fields of the feature flag to retrieve for each flag in the retrieved group.
	Fields []FeatureFlagFields
}

// AllFeatureFlagFields returns a collection of all feature flag fields to use in [FeatureFlagSelector].
func AllFeatureFlagFields() []FeatureFlagFields {
	return []FeatureFlagFields{
		FeatureFlagFieldsName,
		FeatureFlagFieldsLabel,
		FeatureFlagFieldsEnabled,
		FeatureFlagFieldsDescription,
		FeatureFlagFieldsConditions,
		FeatureFlagFieldsAllocation,
		FeatureFlagFieldsVariants,
		FeatureFlagFieldsTelemetry,
		FeatureFlagFieldsTags,
		FeatureFlagFieldsETag,
		FeatureFlagFieldsLastModified,
	}
}

func (fs FeatureFlagSelector) toGeneratedGetFeatureFlags() *generated.AzureAppConfigurationFeatureFlagClientGetFeatureFlagsOptions {
	var dt *string
	if fs.AcceptDateTime != nil {
		str := fs.AcceptDateTime.Format(timeFormat)
		dt = &str
	}
	return &generated.AzureAppConfigurationFeatureFlagClientGetFeatureFlagsOptions{
		AcceptDatetime: dt,
		Name:           fs.NameFilter,
		Label:          fs.LabelFilter,
		Select:         fs.Fields,
		Tags:           fs.TagsFilter,
	}
}

func (fs FeatureFlagSelector) toGeneratedGetFeatureFlagRevisions() *generated.AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsOptions {
	return &generated.AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsOptions{
		Name:   fs.NameFilter,
		Label:  fs.LabelFilter,
		Select: fs.Fields,
		Tags:   fs.TagsFilter,
	}
}
