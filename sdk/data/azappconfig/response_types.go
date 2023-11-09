//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/internal/generated"
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

// ListRevisionsPageResponse contains the configuration settings returned by ListRevisionsPager.
type ListRevisionsPageResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

// ListSettingsPageResponse contains the configuration settings returned by ListRevisionsPager.
type ListSettingsPageResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

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
type ArchiveSnapshotResponse struct {
	generated.Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken *string
}

type ListSnapshotsPagerResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Snapshots []generated.Snapshot

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}
type BeginCreateSnapshotResponse struct {
}

type ListConfigurationSettingsForSnapshotResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

	// SyncToken contains the value returned in the Sync-Token header.
	SyncToken SyncToken
}

type ListSnapshotResponse struct {
	//TODO: Snapshot to constants?
	generated.Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken *string
}

type RecoverSnapshotResponse struct {
	generated.Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken *string
}

type updateSnapshotStatusResponse struct {
	generated.Snapshot

	// Link contains the information returned from the Link header response.
	Link *string

	// SyncToken contains the information returned from the Sync-Token header response.
	SyncToken *string
}
