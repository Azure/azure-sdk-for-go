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

// SetReadOnlyResponse contains the response from SetReadOnly method.
type SetReadOnlyResponse configurationSettingResponse

// SetReadOnlyOptions contains the optional parameters for the SetReadOnly method.
type SetReadOnlyOptions struct {
	// If set to true and the configuration setting exists in the configuration store, update the setting
	// if the passed-in configuration setting is the same version as the one in the configuration store.
	// The setting versions are the same if their ETag fields match.
	OnlyIfUnchanged bool
}

func (cs ConfigurationSetting) toGeneratedPutLockOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientPutLockOptions {
	return &generated.AzureAppConfigurationClientPutLockOptions{
		IfMatch: (*string)(ifMatch),
		Label:   cs.label,
	}
}

func (cs ConfigurationSetting) toGeneratedDeleteLockOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientDeleteLockOptions {
	return &generated.AzureAppConfigurationClientDeleteLockOptions{
		IfMatch: (*string)(ifMatch),
		Label:   cs.label,
	}
}

// SetReadOnly sets an existing configuration setting to read only or read write state in the configuration store.
func (c *Client) SetReadOnly(ctx context.Context, setting ConfigurationSetting, isReadOnly bool, options *SetReadOnlyOptions) (SetReadOnlyResponse, error) {
	var ifMatch *azcore.ETag
	if options != nil && options.OnlyIfUnchanged {
		ifMatch = setting.etag
	}

	var err error
	if isReadOnly {
		var resp generated.AzureAppConfigurationClientPutLockResponse
		resp, err = c.appConfigClient.PutLock(ctx, *setting.key, setting.toGeneratedPutLockOptions(ifMatch))
		if err == nil {
			return (SetReadOnlyResponse)(responseFromGeneratedPutLock(resp)), nil
		}
	} else {
		var resp generated.AzureAppConfigurationClientDeleteLockResponse
		resp, err = c.appConfigClient.DeleteLock(ctx, *setting.key, setting.toGeneratedDeleteLockOptions(ifMatch))
		if err == nil {
			return (SetReadOnlyResponse)(responseFromGeneratedDeleteLock(resp)), nil
		}
	}

	return SetReadOnlyResponse{}, err
}
