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

type SetReadOnlyResponse configurationSettingResponse

type SetReadOnlyOptions struct {
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
