// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
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
	expected := []string{
		"NewPrivateEndpointConnectionsClient",
		"*PrivateEndpointConnectionsClient.BeginCreate",
		"*PrivateEndpointConnectionsClient.BeginDelete",
		"*PrivateEndpointConnectionsClient.Get",
		"*PrivateEndpointConnectionsClient.NewListPager",
	}
	assert.Equal(t, expected, sortResult)
}

func TestRemovedConstAndTypeAlias(t *testing.T) {
	removedConst := delta.Content{
		Content: exports.Content{
			Consts: map[string]exports.Const{
				"ConstA": {
					Type: "Const",
				},
				"ConstB": {
					Type: "Const",
				},
				"ConstC": {
					Type: "Const",
				},
				"RemovedTypeAliasA": {
					Type: "RemovedTypeAlias",
				},
				"RemovedTypeAliasB": {
					Type: "RemovedTypeAlias",
				},
			},
			TypeAliases: map[string]exports.TypeAlias{
				"RemovedTypeAlias": {
					UnderlayingType: "string",
				},
			},
		},
	}

	actual := getRemovedContent(&removedConst)
	expected := []string{
		"Const `ConstA`, `ConstB`, `ConstC` from type alias `Const` has been removed",
		"Type alias `RemovedTypeAlias` has been removed",
	}
	assert.Equal(t, expected, actual)
}
