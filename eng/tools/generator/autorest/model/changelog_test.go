// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/stretchr/testify/assert"
)

func TestSortFuncItem(t *testing.T) {
	get := "PrivateEndpointConnectionsClientGetResponse, error"
	beginDelete := "*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error"
	beginCreate := "*runtime.Poller[PrivateEndpointConnectionsClientCreateResponse], error"
	newListPager := "*runtime.Pager[PrivateEndpointConnectionsClientListResponse]"
	newClient := "*PrivateEndpointConnectionsClient, error"

	s := map[string]exports.Func{
		"*PrivateEndpointConnectionsClient.Get": {
			Returns: &get,
		},
		"*PrivateEndpointConnectionsClient.BeginDelete": {
			Returns: &beginDelete,
		},
		"*PrivateEndpointConnectionsClient.BeginCreate": {
			Returns: &beginCreate,
		},
		"*PrivateEndpointConnectionsClient.NewListPager": {
			Returns: &newListPager,
		},
		"NewPrivateEndpointConnectionsClient": {
			Returns: &newClient,
		},
	}

	sortResult := sortFuncItem(s)
	expcted := []string{
		"NewPrivateEndpointConnectionsClient",
		"*PrivateEndpointConnectionsClient.BeginCreate",
		"*PrivateEndpointConnectionsClient.BeginDelete",
		"*PrivateEndpointConnectionsClient.Get",
		"*PrivateEndpointConnectionsClient.NewListPager",
	}
	assert.Equal(t, expcted, sortResult)
}
