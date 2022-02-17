//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"
	"sdk/appconfiguration/sdk/appconfiguration/azappconfiguration/internal/generated"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type GetConfigurationSettingOptions struct {
	OnlyIfChanged  bool
	AcceptDateTime *time.Time
}

func (cs ConfigurationSetting) toGeneratedGetOptions(ifNoneMatch *azcore.ETag, acceptDateTime *time.Time) *generated.AzureAppConfigurationClientGetKeyValueOptions {
	dt := acceptDateTime.Format(timeFormat)
	return &generated.AzureAppConfigurationClientGetKeyValueOptions{
		AcceptDatetime: &dt,
		IfNoneMatch:    (*string)(ifNoneMatch),
		Label:          cs.Label,
	}
}

func (c *Client) GetConfigurationSetting(ctx context.Context, setting ConfigurationSetting, options *GetConfigurationSettingOptions) (GetConfigurationSettingResponse, error) {
	var ifNoneMatch *azcore.ETag
	var acceptDateTime *time.Time
	if options != nil {
		if options.OnlyIfChanged {
			ifNoneMatch = setting.etag
		}

		acceptDateTime = options.AcceptDateTime
	}

	resp, err := c.appConfigClient.GetKeyValue(ctx, *setting.Key, setting.toGeneratedGetOptions(ifNoneMatch, acceptDateTime))
	if err != nil {
		return GetConfigurationSettingResponse{}, err
	}

	return responseFromGeneratedGet(resp), nil
}
