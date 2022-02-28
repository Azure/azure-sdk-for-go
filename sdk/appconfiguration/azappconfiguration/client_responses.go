//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfiguration/internal/generated"
)

// ConfigurationSettingResult contains the configuration setting returned from Azure App Configuration client methods.
type ConfigurationSettingResult struct {
	ConfigurationSetting

	// ETag of the configuration setting.
	ETag *azcore.ETag

	// Sync token for the Azure App Configuration client, corresponding to the curent state of the client.
	SyncToken *string
}

type configurationSettingResponse struct {
	// Contains the configuration setting returned by Azure App Configuration client methods.
	Result ConfigurationSettingResult

	// Contains the raw HTTP response from Azure App Configuration client method invocation.
	RawResponse *http.Response
}

func responseFromGeneratedPut(g generated.AzureAppConfigurationClientPutKeyValueResponse) configurationSettingResponse {
	return configurationSettingResponse{
		Result: ConfigurationSettingResult{
			ConfigurationSetting: configurationSettingFromGenerated(g.KeyValue),
			ETag:                 (*azcore.ETag)(g.Etag),
			SyncToken:            g.SyncToken,
		},
		RawResponse: g.RawResponse,
	}
}

func responseFromGeneratedDelete(g generated.AzureAppConfigurationClientDeleteKeyValueResponse) configurationSettingResponse {
	return configurationSettingResponse{
		Result: ConfigurationSettingResult{
			ConfigurationSetting: configurationSettingFromGenerated(g.KeyValue),
			ETag:                 (*azcore.ETag)(g.Etag),
			SyncToken:            g.SyncToken,
		},
		RawResponse: g.RawResponse,
	}
}

func responseFromGeneratedPutLock(g generated.AzureAppConfigurationClientPutLockResponse) configurationSettingResponse {
	return configurationSettingResponse{
		Result: ConfigurationSettingResult{
			ConfigurationSetting: configurationSettingFromGenerated(g.KeyValue),
			ETag:                 (*azcore.ETag)(g.Etag),
			SyncToken:            g.SyncToken,
		},
		RawResponse: g.RawResponse,
	}
}

func responseFromGeneratedDeleteLock(g generated.AzureAppConfigurationClientDeleteLockResponse) configurationSettingResponse {
	return configurationSettingResponse{
		Result: ConfigurationSettingResult{
			ConfigurationSetting: configurationSettingFromGenerated(g.KeyValue),
			ETag:                 (*azcore.ETag)(g.Etag),
			SyncToken:            g.SyncToken,
		},
		RawResponse: g.RawResponse,
	}
}

// GetConfigurationSettingResult contains the configuration setting retrieved by GetConfigurationSetting method.
type GetConfigurationSettingResult struct {
	ConfigurationSettingResult

	// Contains the timestamp of when the configuration setting was last modified.
	LastModified *time.Time
}

// GetConfigurationSettingResponse contains the response from GetConfigurationSetting method.
type GetConfigurationSettingResponse struct {
	// Contains the configuration setting returned by the GetConfigurationSetting method.
	Result GetConfigurationSettingResult

	// Contains the raw HTTP response from Azure App Configuration client method invocation.
	RawResponse *http.Response
}

func responseFromGeneratedGet(g generated.AzureAppConfigurationClientGetKeyValueResponse) GetConfigurationSettingResponse {
	var t *time.Time
	if g.LastModified != nil {
		if tt, err := time.Parse(timeFormat, *g.LastModified); err == nil {
			t = &tt
		}
	}

	return GetConfigurationSettingResponse{
		Result: GetConfigurationSettingResult{
			ConfigurationSettingResult: ConfigurationSettingResult{
				ConfigurationSetting: configurationSettingFromGenerated(g.KeyValue),
				ETag:                 (*azcore.ETag)(g.Etag),
				SyncToken:            g.SyncToken,
			},
			LastModified: t,
		},
		RawResponse: g.RawResponse,
	}
}

// ConfigurationSettingResult contains the configuration setting returned from Azure App Configuration client methods.
type GetRevisionsResult struct {
	// Contains the configuration settings returned that match the setting selector provided.
	Items []ConfigurationSetting

	// Sync token for the Azure App Configuration client, corresponding to the curent state of the client.
	SyncToken *string
}

// GetRevisionsPage contains the response returned from the GetRevisions method call.
type GetRevisionsPage struct {
	// Contains the configuration settings returned by the GetRevisions method call.
	Result GetRevisionsResult

	// Contains the raw HTTP response from Azure App Configuration client method invocation.
	RawResponse *http.Response
}

func getRevisionsPageFromGenerated(g generated.AzureAppConfigurationClientGetRevisionsResponse) GetRevisionsPage {
	var css []ConfigurationSetting
	for _, cs := range g.Items {
		if cs != nil {
			css = append(css, configurationSettingFromGenerated(*cs))
		}
	}

	return GetRevisionsPage{
		Result: GetRevisionsResult{
			Items:     css,
			SyncToken: g.SyncToken,
		},
		RawResponse: g.RawResponse,
	}
}
