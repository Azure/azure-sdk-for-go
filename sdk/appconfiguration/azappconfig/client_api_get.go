//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfig/internal/generated"
)

// GetConfigurationSettingOptions contains the optional parameters for the GetConfigurationSetting method.
type GetConfigurationSettingOptions struct {
	// If set to true, only retrieve the setting from the configuration store if it has changed since the client last retrieved it.
	// It is determined to have changed if the ETag field on the passed-in configuration setting is different from the ETag
	// of the setting in the configuration store.
	OnlyIfChanged bool

	// The setting will be retrieved exactly as it existed at the provided time.
	AcceptDateTime *time.Time
}

func (cs Setting) toGeneratedGetOptions(ifNoneMatch *azcore.ETag, acceptDateTime *time.Time) *generated.AzureAppConfigurationClientGetKeyValueOptions {
	var dt *string
	if acceptDateTime != nil {
		str := acceptDateTime.Format(timeFormat)
		dt = &str
	}

	return &generated.AzureAppConfigurationClientGetKeyValueOptions{
		AcceptDatetime: dt,
		IfNoneMatch:    (*string)(ifNoneMatch),
		Label:          cs.Label,
	}
}

// GetConfigurationSetting retrieves an existing configuration setting from the configuration store.
func (c *Client) GetConfigurationSetting(ctx context.Context, setting Setting, options *GetConfigurationSettingOptions) (GetConfigurationSettingResponse, error) {
	var ifNoneMatch *azcore.ETag
	var acceptDateTime *time.Time
	if options != nil {
		if options.OnlyIfChanged {
			ifNoneMatch = setting.ETag
		}

		acceptDateTime = options.AcceptDateTime
	}

	resp, err := c.appConfigClient.GetKeyValue(ctx, *setting.Key, setting.toGeneratedGetOptions(ifNoneMatch, acceptDateTime))
	if err != nil {
		return GetConfigurationSettingResponse{}, err
	}

	return responseFromGeneratedGet(resp), nil
}
