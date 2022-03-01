//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AddConfigurationSettingResponse contains the response from AddConfigurationSetting method.
type AddConfigurationSettingResponse configurationSettingResponse

// AddConfigurationSettingOptions contains the optional parameters for the AddConfigurationSetting method.
type AddConfigurationSettingOptions struct {
}

// AddConfigurationSetting creates a configuration setting only if the setting does not already exist in the configuration store.
func (c *Client) AddConfigurationSetting(ctx context.Context, setting Setting, options *AddConfigurationSettingOptions) (AddConfigurationSettingResponse, error) {
	etagAny := azcore.ETagAny
	resp, err := c.appConfigClient.PutKeyValue(ctx, *setting.key, setting.toGeneratedPutOptions(nil, &etagAny))
	if err != nil {
		return AddConfigurationSettingResponse{}, err
	}

	return (AddConfigurationSettingResponse)(responseFromGeneratedPut(resp)), nil
}
