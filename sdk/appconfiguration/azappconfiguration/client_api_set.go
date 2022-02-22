//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type SetConfigurationSettingResponse configurationSettingResponse

type SetConfigurationSettingOptions struct {
	OnlyIfUnchanged bool
}

func (c *Client) SetConfigurationSetting(ctx context.Context, setting ConfigurationSetting, options *SetConfigurationSettingOptions) (SetConfigurationSettingResponse, error) {
	var ifMatch *azcore.ETag
	if options != nil && options.OnlyIfUnchanged {
		ifMatch = setting.etag
	}

	resp, err := c.appConfigClient.PutKeyValue(ctx, *setting.key, setting.toGeneratedPutOptions(ifMatch, nil))
	if err != nil {
		return SetConfigurationSettingResponse{}, err
	}

	return (SetConfigurationSettingResponse)(responseFromGeneratedPut(resp)), nil
}
