//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// SetConfigurationSettingResponse contains the response from SetConfigurationSetting method.
type SetConfigurationSettingResponse configurationSettingResponse

// SetConfigurationSettingOptions contains the optional parameters for the SetConfigurationSetting method.
type SetConfigurationSettingOptions struct {
	// If set to true and the configuration setting exists in the configuration store, overwrite the setting
	// if the passed-in configuration setting is the same version as the one in the configuration store.
	// The setting versions are the same if their ETag fields match.
	OnlyIfUnchanged bool
}

// SetConfigurationSetting creates a configuration setting if it doesn't exist or overwrites the existing setting in the configuration store.
func (c *Client) SetConfigurationSetting(ctx context.Context, setting Setting, options *SetConfigurationSettingOptions) (SetConfigurationSettingResponse, error) {
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
