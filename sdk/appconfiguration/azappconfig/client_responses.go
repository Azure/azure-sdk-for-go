//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfig/internal/generated"
)

type configurationSettingResult struct {
	Setting

	// Sync token for the Azure App Configuration client, corresponding to the curent state of the client.
	SyncToken *string
}

// AddConfigurationSettingResult contains the response from AddConfigurationSetting method.
type AddConfigurationSettingResult configurationSettingResult

// SetConfigurationSettingResult contains the response from SetConfigurationSetting method.
type SetConfigurationSettingResult configurationSettingResult

// DeleteConfigurationSettingResult contains the response from DeleteConfigurationSetting method.
type DeleteConfigurationSettingResult configurationSettingResult

// SetReadOnlyResult contains the response from SetReadOnly method.
type SetReadOnlyResult configurationSettingResult

func fromGeneratedPut(g generated.AzureAppConfigurationClientPutKeyValueResponse) configurationSettingResult {
	return configurationSettingResult{
		Setting:   configurationSettingFromGenerated(g.KeyValue),
		SyncToken: g.SyncToken,
	}
}

func fromGeneratedDelete(g generated.AzureAppConfigurationClientDeleteKeyValueResponse) DeleteConfigurationSettingResult {
	return DeleteConfigurationSettingResult{
		Setting:   configurationSettingFromGenerated(g.KeyValue),
		SyncToken: g.SyncToken,
	}
}

func fromGeneratedPutLock(g generated.AzureAppConfigurationClientPutLockResponse) SetReadOnlyResult {
	return SetReadOnlyResult{
		Setting:   configurationSettingFromGenerated(g.KeyValue),
		SyncToken: g.SyncToken,
	}
}

func fromGeneratedDeleteLock(g generated.AzureAppConfigurationClientDeleteLockResponse) SetReadOnlyResult {
	return SetReadOnlyResult{
		Setting:   configurationSettingFromGenerated(g.KeyValue),
		SyncToken: g.SyncToken,
	}
}

// GetConfigurationSettingResult contains the configuration setting retrieved by GetConfigurationSetting method.
type GetConfigurationSettingResult struct {
	configurationSettingResult

	// Contains the timestamp of when the configuration setting was last modified.
	LastModified *time.Time
}

func fromGeneratedGet(g generated.AzureAppConfigurationClientGetKeyValueResponse) GetConfigurationSettingResult {
	var t *time.Time
	if g.LastModified != nil {
		if tt, err := time.Parse(timeFormat, *g.LastModified); err == nil {
			t = &tt
		}
	}

	return GetConfigurationSettingResult{
		configurationSettingResult: configurationSettingResult{
			Setting:   configurationSettingFromGenerated(g.KeyValue),
			SyncToken: g.SyncToken,
		},
		LastModified: t,
	}
}

// ListRevisionsPage contains the configuration settings returned by ListRevisionsPager.
type ListRevisionsPage struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Settings []Setting

	// Sync token for the Azure App Configuration client, corresponding to the current state of the client.
	SyncToken *string
}

func fromGeneratedGetRevisionsPage(g generated.AzureAppConfigurationClientGetRevisionsResponse) ListRevisionsPage {
	var css []Setting
	for _, cs := range g.Items {
		if cs != nil {
			css = append(css, configurationSettingFromGenerated(*cs))
		}
	}

	return ListRevisionsPage{
		Settings:  css,
		SyncToken: g.SyncToken,
	}
}
