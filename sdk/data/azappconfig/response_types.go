// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/exported"
)

// SyncToken contains data used in the Sync-Token header.
// See [Azure App Configuration documentation] for more information on sync tokens.
//
// [Azure App Configuration documentation]: https://learn.microsoft.com/azure/azure-app-configuration/rest-api-consistency
type SyncToken = exported.SyncToken

// AddSettingResponse contains the response from AddSetting method.
type AddSettingResponse struct {
	Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// DeleteSettingResponse contains the response from DeleteSetting method.
type DeleteSettingResponse struct {
	Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// GetSettingResponse contains the configuration setting retrieved by GetSetting method.
type GetSettingResponse struct {
	Setting

	// Contains the timestamp of when the configuration setting was last modified.
	LastModified *time.Time

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// ListRevisionsPageResponse contains the response from the NewListRevisionsPager method.
type ListRevisionsPageResponse struct {
	// Contains the configuration setting revisions that match the setting selector provided.
	Settings []Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// ListSettingsPageResponse contains the response from the NewListSettingsPager method.
type ListSettingsPageResponse struct {
	// Contains the configuration settings that match the setting selector provided.
	Settings []Setting

	// An ETag indicating the state of a page of configuration settings within a configuration store.
	ETag *azcore.ETag

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// SetReadOnlyResponse contains the response from SetReadOnly method.
type SetReadOnlyResponse struct {
	Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// SetSettingResponse contains the response from SetSetting method.
type SetSettingResponse struct {
	Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// ArchiveSnapshotResponse contains the response from the ArchiveSnapshot method.
type ArchiveSnapshotResponse struct {
	Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken SyncToken
}

// ListSnapshotsResponse contains the response from the NewGetSnapshotsPager method.
type ListSnapshotsResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Snapshots []Snapshot

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// CreateSnapshotResponse contains the response from the BeginCreateSnapshot method.
type CreateSnapshotResponse struct {
	// Read-Only information about the snapshot retrieved from a Create Snapshot operation.
	Snapshot
}

// ListSettingsForSnapshotResponse contains the response from the ListConfigurationSettingsForSnapshot method.
type ListSettingsForSnapshotResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// GetSnapshotResponse contains the response from the GetSnapshot method.
type GetSnapshotResponse struct {
	// Snapshot object in GetSnapshot Response
	Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken SyncToken
}

// RecoverSnapshotResponse contains the response from the RecoverSnapshot method.
type RecoverSnapshotResponse struct {
	Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken SyncToken
}

// updateSnapshotStatusResponse contains the response from the UpdateSnapshotStatus method.
type updateSnapshotStatusResponse struct {
	Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken SyncToken
}
