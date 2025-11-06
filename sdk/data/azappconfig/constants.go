//go:build go1.18
// +build go1.18

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
