// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
)

// SettingFields are fields to retrieve from a configuration setting.
type SettingFields = generated.SettingFields

const (
	// The primary identifier of a configuration setting.
	SettingFieldsKey SettingFields = generated.SettingFieldsKey

	// A label used to group configuration settings.
	SettingFieldsLabel SettingFields = generated.SettingFieldsLabel

	// The value of the configuration setting.
	SettingFieldsValue SettingFields = generated.SettingFieldsValue

	// The content type of the configuration setting's value.
	SettingFieldsContentType SettingFields = generated.SettingFieldsContentType

	// The description of a configuration setting.
	SettingFieldsDescription SettingFields = generated.SettingFieldsDescription

	// An ETag indicating the version of a configuration setting within a configuration store.
	SettingFieldsETag SettingFields = generated.SettingFieldsEtag

	// The last time a modifying operation was performed on the given configuration setting.
	SettingFieldsLastModified SettingFields = generated.SettingFieldsLastModified

	// A value indicating whether the configuration setting is read-only.
	SettingFieldsIsReadOnly SettingFields = generated.SettingFieldsLocked

	// A list of tags that can help identify what a configuration setting may be applicable for.
	SettingFieldsTags SettingFields = generated.SettingFieldsTags
)

// SnapshotFields are fields to retrieve from a snapshot.
type SnapshotFields = generated.SnapshotFields

const (
	// The composition type of a snapshot.
	SnapshotFieldsCompositionType SnapshotFields = generated.SnapshotFieldsCompositionType

	// The time when the snapshot was created.
	SnapshotFieldsCreated SnapshotFields = generated.SnapshotFieldsCreated

	// The description of a snapshot.
	SnapshotFieldsDescription SnapshotFields = generated.SnapshotFieldsDescription

	// An ETag indicating the version of a snapshot.
	SnapshotFieldsETag SnapshotFields = generated.SnapshotFieldsETag

	// The time when the snapshot will expire once archived.
	SnapshotFieldsExpires SnapshotFields = generated.SnapshotFieldsExpires

	// A list of filters used to generate the snapshot.
	SnapshotFieldsFilters SnapshotFields = generated.SnapshotFieldsFilters

	// The number of items in the snapshot.
	SnapshotFieldsItemsCount SnapshotFields = generated.SnapshotFieldsItemsCount

	// The primary identifier of a snapshot.
	SnapshotFieldsName SnapshotFields = generated.SnapshotFieldsName

	// Retention period in seconds of the snapshot upon archiving.
	SnapshotFieldsRetentionPeriod SnapshotFields = generated.SnapshotFieldsRetentionPeriod

	// Size of the snapshot.
	SnapshotFieldsSize SnapshotFields = generated.SnapshotFieldsSize

	// Status of the snapshot.
	SnapshotFieldsStatus SnapshotFields = generated.SnapshotFieldsStatus

	// A list of tags on the snapshot.
	SnapshotFieldsTags SnapshotFields = generated.SnapshotFieldsTags
)

// SnapshotStatus contains the current status of the snapshot
type SnapshotStatus = generated.SnapshotStatus

const (
	// Snapshot is archived state.
	SnapshotStatusArchived SnapshotStatus = generated.SnapshotStatusArchived

	// Snapshot is in failing state.
	SnapshotStatusFailed SnapshotStatus = generated.SnapshotStatusFailed

	// Snapshot is in provisioning state.
	SnapshotStatusProvisioning SnapshotStatus = generated.SnapshotStatusProvisioning

	// Snapshot is in ready state.
	SnapshotStatusReady SnapshotStatus = generated.SnapshotStatusReady
)

// CompositionType is the composition of filters used to create a snapshot.
type CompositionType = generated.CompositionType

const (
	// Snapshot is composed with a Key filter
	CompositionTypeKey CompositionType = generated.CompositionTypeKey

	// Snapshot is composed with a Key and Label filter
	CompositionTypeKeyLabel CompositionType = generated.CompositionTypeKeyLabel
)

// FeatureFlagFields are fields to retrieve from a feature flag.
type FeatureFlagFields = generated.FeatureFlagFields

const (
	// The name of a feature flag.
	FeatureFlagFieldsName FeatureFlagFields = generated.FeatureFlagFieldsName

	// The label of a feature flag.
	FeatureFlagFieldsLabel FeatureFlagFields = generated.FeatureFlagFieldsLabel

	// The enabled state of a feature flag.
	FeatureFlagFieldsEnabled FeatureFlagFields = generated.FeatureFlagFieldsEnabled

	// The description of a feature flag.
	FeatureFlagFieldsDescription FeatureFlagFields = generated.FeatureFlagFieldsDescription

	// The conditions of a feature flag.
	FeatureFlagFieldsConditions FeatureFlagFields = generated.FeatureFlagFieldsConditions

	// The allocation of a feature flag.
	FeatureFlagFieldsAllocation FeatureFlagFields = generated.FeatureFlagFieldsAllocation

	// The variants of a feature flag.
	FeatureFlagFieldsVariants FeatureFlagFields = generated.FeatureFlagFieldsVariants

	// The telemetry configuration of a feature flag.
	FeatureFlagFieldsTelemetry FeatureFlagFields = generated.FeatureFlagFieldsTelemetry

	// The tags of a feature flag.
	FeatureFlagFieldsTags FeatureFlagFields = generated.FeatureFlagFieldsTags

	// An ETag indicating the version of a feature flag within a configuration store.
	FeatureFlagFieldsETag FeatureFlagFields = generated.FeatureFlagFieldsEtag

	// The last time a modifying operation was performed on the feature flag.
	FeatureFlagFieldsLastModified FeatureFlagFields = generated.FeatureFlagFieldsLastModified
)

// RequirementType describes how filters on a feature flag are combined.
type RequirementType = generated.RequirementType

const (
	// The feature flag is enabled when all filters evaluate to true.
	RequirementTypeAll RequirementType = generated.RequirementTypeAll

	// The feature flag is enabled when any filter evaluates to true.
	RequirementTypeAny RequirementType = generated.RequirementTypeAny
)

// StatusOverride overrides the enabled state of a feature flag when a variant is chosen.
type StatusOverride = generated.StatusOverride

const (
	// Do not override the enabled state of the feature flag.
	StatusOverrideNone StatusOverride = generated.StatusOverrideNone

	// Force the feature flag to be enabled.
	StatusOverrideEnabled StatusOverride = generated.StatusOverrideEnabled

	// Force the feature flag to be disabled.
	StatusOverrideDisabled StatusOverride = generated.StatusOverrideDisabled
)
