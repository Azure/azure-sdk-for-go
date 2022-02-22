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

type ConfigurationSettingResult struct {
	ConfigurationSetting
	ETag      *azcore.ETag
	SyncToken *string
}

type configurationSettingResponse struct {
	Result      ConfigurationSettingResult
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

type GetConfigurationSettingResult struct {
	ConfigurationSettingResult
	LastModified *time.Time
}

type GetConfigurationSettingResponse struct {
	Result      GetConfigurationSettingResult
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

type GetRevisionsResult struct {
	Items     []ConfigurationSetting
	SyncToken *string
}

type GetRevisionsPage struct {
	GetRevisionsResult
	RawResponse *http.Response
}

func getRevisionsPageFromGenerated(g generated.AzureAppConfigurationClientGetKeyValuesResponse) GetRevisionsPage {
	var css []ConfigurationSetting
	for _, cs := range g.Items {
		if cs != nil {
			css = append(css, configurationSettingFromGenerated(cs))
		}
	}

	return GetRevisionsPage{
		GetRevisionsResult: GetRevisionsResult{
			Items:     css,
			SyncToken: g.SyncToken,
		},
		RawResponse: g.RawResponse,
	}
}
