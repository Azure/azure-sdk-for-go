//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import "time"

// AddSettingResponse contains the response from AddSetting method.
type AddSettingResponse struct {
	Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

// DeleteSettingResponse contains the response from DeleteSetting method.
type DeleteSettingResponse struct {
	Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

// GetSettingResponse contains the configuration setting retrieved by GetSetting method.
type GetSettingResponse struct {
	Setting

	// Contains the timestamp of when the configuration setting was last modified.
	LastModified *time.Time

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

// ListRevisionsPageResponse contains the configuration settings returned by ListRevisionsPager.
type ListRevisionsPageResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

// ListSettingsPageResponse contains the configuration settings returned by ListRevisionsPager.
type ListSettingsPageResponse struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

// SetReadOnlyResponse contains the response from SetReadOnly method.
type SetReadOnlyResponse struct {
	Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

// SetSettingResponse contains the response from SetSetting method.
type SetSettingResponse struct {
	Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}
