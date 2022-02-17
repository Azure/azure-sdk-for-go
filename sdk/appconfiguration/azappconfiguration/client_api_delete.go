//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"
	"sdk/appconfiguration/sdk/appconfiguration/azappconfiguration/internal/generated"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type DeleteConfigurationSettingResponse configurationSettingResponse

type DeleteConfigurationSettingOptions struct {
	OnlyIfUnchanged bool
}

func (cs ConfigurationSetting) toGeneratedDeleteOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientDeleteKeyValueOptions {
	return &generated.AzureAppConfigurationClientDeleteKeyValueOptions{
		IfMatch: (*string)(ifMatch),
		Label:   cs.Label,
	}
}

func (c *Client) DeleteConfigurationSetting(ctx context.Context, setting ConfigurationSetting, options *DeleteConfigurationSettingOptions) (DeleteConfigurationSettingResponse, error) {
	var ifMatch *azcore.ETag
	if options != nil && options.OnlyIfUnchanged {
		ifMatch = setting.etag
	}

	resp, err := c.appConfigClient.DeleteKeyValue(ctx, *setting.Key, setting.toGeneratedDeleteOptions(ifMatch))
	if err != nil {
		return DeleteConfigurationSettingResponse{}, err
	}

	return (DeleteConfigurationSettingResponse)(responseFromGeneratedDelete(resp)), nil
}
