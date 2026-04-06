// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AddSettingOptions contains the optional parameters for the AddSetting method.
type AddSettingOptions struct {
	// Configuration setting content type.
	ContentType *string

	// Configuration setting label.
	Label *string

	// A dictionary of tags used to assign additional properties to a configuration setting.
	Tags map[string]*string
}

// DeleteSettingOptions contains the optional parameters for the DeleteSetting method.
type DeleteSettingOptions struct {
	// Configuration setting label.
	Label *string

	// If set, and the configuration setting exists in the configuration store,
	// delete the setting if the passed-in ETag is the same as the setting's ETag in the configuration store.
	//
	// This has IfMatch semantics.
	OnlyIfUnchanged *azcore.ETag
}

// GetSettingOptions contains the optional parameters for the GetSetting method.
type GetSettingOptions struct {
	// The setting will be retrieved exactly as it existed at the provided time.
	AcceptDateTime *time.Time

	// Configuration setting label.
	Label *string

	// If set, only retrieve the setting from the configuration store if setting has changed
	// since the client last retrieved it with the ETag provided.
	//
	// This has IfNoneMatch semantics.
	OnlyIfChanged *azcore.ETag
}

// ListRevisionsOptions contains the optional parameters for the NewListRevisionsPager method.
type ListRevisionsOptions struct {
	// placeholder for future options
}

// ListSettingsOptions contains the optional parameters for the NewListSettingsPager method.
type ListSettingsOptions struct {
	// The match conditions used when making the request.
	// Conditions are applied to pages one by one in the order specified.
	MatchConditions []azcore.MatchConditions
}

// SetReadOnlyOptions contains the optional parameters for the SetReadOnly method.
type SetReadOnlyOptions struct {
	// Configuration setting label.
	Label *string

	// If set, and the configuration setting exists in the configuration store, update the setting
	// if the passed-in configuration setting ETag is the same version as the one in the configuration store.
	//
	// This has IfMatch semantics.
	OnlyIfUnchanged *azcore.ETag
}

// SetSettingOptions contains the optional parameters for the SetSetting method.
type SetSettingOptions struct {
	// Configuration setting content type.
	ContentType *string

	// Configuration setting label.
	Label *string

	// A dictionary of tags used to assign additional properties to a configuration setting.
	// These can be used to indicate how a configuration setting may be applied.
	Tags map[string]*string

	// If set, and the configuration setting exists in the configuration store, overwrite the setting
	// if the passed-in ETag is the same version as the one in the configuration store.
	//
	// This has IfMatch semantics.
	OnlyIfUnchanged *azcore.ETag
}

// BeginCreateSnapshotOptions contains the optional parameters for the BeginCreateSnapshot method.
type BeginCreateSnapshotOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string

	// The composition type describes how the key-values within the snapshot are composed. The 'key' composition type ensures
	// there are no two key-values containing the same key. The 'key_label' composition
	// type ensures there are no two key-values containing the same key and label.
	CompositionType *CompositionType

	// The amount of time, in seconds, that a snapshot will remain in the archived state before expiring. This property is only
	// writable during the creation of a snapshot. If not specified, the default
	// lifetime of key-value revisions will be used.
	RetentionPeriod *int64

	// The tags of the snapshot.
	Tags map[string]*string
}

// ArchiveSnapshotOptions contains the optional parameters for the ArchiveSnapshot method.
type ArchiveSnapshotOptions struct {
	// Used to perform an operation only if the targeted resource's etag matches the value provided.
	IfMatch *azcore.ETag

	// Used to perform an operation only if the targeted resource's etag does not match the value provided.
	IfNoneMatch *azcore.ETag
}

// RestoreSnapshotOptions contains the optional parameters for the RestoreSnapshot method.
type RestoreSnapshotOptions struct {
	// Used to perform an operation only if the targeted resource's etag matches the value provided.
	IfMatch *azcore.ETag

	// Used to perform an operation only if the targeted resource's etag does not match the value provided.
	IfNoneMatch *azcore.ETag
}

// ListSnapshotsOptions contains the optional parameters for the ListSnapshotsPager method.
type ListSnapshotsOptions struct {
	// Instructs the server to return elements that appear after the element referred to by the specified token.
	After *string

	// A filter for the name of the returned snapshots.
	Name *string

	// Used to select what fields are present in the returned resource(s).
	Select []SnapshotFields

	// Used to filter returned snapshots by their status property.
	Status []SnapshotStatus
}

// ListSettingsForSnapshotOptions contains the optional parameters for the NewListSettingsForSnapshotPager method.
type ListSettingsForSnapshotOptions struct {
	// Requests the server to respond with the state of the resource at the specified time.
	AcceptDatetime *string

	// Instructs the server to return elements that appear after the element referred to by the specified token.
	After *string

	// Used to perform an operation only if the targeted resource's etag matches the value provided.
	IfMatch *azcore.ETag

	// Used to perform an operation only if the targeted resource's etag does not match the value provided.
	IfNoneMatch *azcore.ETag

	// Used to select what fields are present in the returned resource(s).
	Select []SettingFields

	// A filter used to match Keys
	Key string

	// A filter used to match Labels
	Label string
}

// GetSnapshotOptions contains the optional parameters for the GetSnapshot method.
type GetSnapshotOptions struct {
	// Used to perform an operation only if the targeted resource's etag matches the value provided.
	IfMatch *azcore.ETag

	// Used to perform an operation only if the targeted resource's etag does not match the value provided.
	IfNoneMatch *azcore.ETag

	// Used to select what fields are present in the returned resource(s).
	Select []SnapshotFields
}

// RecoverSnapshotOptions contains the optional parameters for the RecoverSnapshot method.
type RecoverSnapshotOptions struct {
	// Used to perform an operation only if the targeted resource's etag matches the value provided.
	IfMatch *azcore.ETag

	// Used to perform an operation only if the targeted resource's etag does not match the value provided.
	IfNoneMatch *azcore.ETag
}

// UpdateSnapshotStatusOptions contains the optional parameters for the UpdateSnapshotStatus method.
type updateSnapshotStatusOptions struct {
	// Used to perform an operation only if the targeted resource's etag matches the value provided.
	IfMatch *azcore.ETag

	// Used to perform an operation only if the targeted resource's etag does not match the value provided.
	IfNoneMatch *azcore.ETag
}
