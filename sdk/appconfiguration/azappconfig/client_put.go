//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfig/internal/generated"
)

func (cs Setting) toGeneratedPutOptions(ifMatch *azcore.ETag, ifNoneMatch *azcore.ETag) *generated.AzureAppConfigurationClientPutKeyValueOptions {
	return &generated.AzureAppConfigurationClientPutKeyValueOptions{
		Entity:      cs.toGenerated(),
		IfMatch:     (*string)(ifMatch),
		IfNoneMatch: (*string)(ifNoneMatch),
		Label:       cs.Label,
	}
}
