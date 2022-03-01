//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfiguration/internal/generated"
)

// DeleteConfigurationSettingResponse contains the response from DeleteConfigurationSetting method.
type DeleteConfigurationSettingResponse configurationSettingResponse

// DeleteConfigurationSettingOptions contains the optional parameters for the DeleteConfigurationSetting method.
type DeleteConfigurationSettingOptions struct {
	// If set to true and the configuration setting exists in the configuration store,
	// delete the setting if the passed-in configuration setting is the same version as the one in the configuration store.
	// The setting versions are the same if their ETag fields match.
	OnlyIfUnchanged bool
}

func (cs Setting) toGeneratedDeleteOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientDeleteKeyValueOptions {
	return &generated.AzureAppConfigurationClientDeleteKeyValueOptions{
		IfMatch: (*string)(ifMatch),
		Label:   cs.Label,
	}
}

// DeleteConfigurationSetting deletes a configuration setting from the configuration store.
func (c *Client) DeleteConfigurationSetting(ctx context.Context, setting Setting, options *DeleteConfigurationSettingOptions) (DeleteConfigurationSettingResponse, error) {
	var ifMatch *azcore.ETag
	if options != nil && options.OnlyIfUnchanged {
		ifMatch = setting.ETag
	}

	resp, err := c.appConfigClient.DeleteKeyValue(ctx, *setting.Key, setting.toGeneratedDeleteOptions(ifMatch))
	if err != nil {
		return DeleteConfigurationSettingResponse{}, err
	}

	return (DeleteConfigurationSettingResponse)(responseFromGeneratedDelete(resp)), nil
}
