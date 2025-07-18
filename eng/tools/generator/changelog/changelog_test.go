// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
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

	sortResult := SortFuncItem(s)
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
		"`ConstA`, `ConstB`, `ConstC` from enum `Const` has been removed",
		"Enum `RemovedTypeAlias` has been removed",
	}
	assert.Equal(t, expected, actual)
}

func TestCombineSimilarItem(t *testing.T) {
	r := &report.Package{
		AdditiveChanges: &report.AdditiveChanges{
			Added: &delta.Content{
				Content: exports.Content{
					Structs: map[string]exports.Struct{
						"Struct": {
							AnonymousFields: []string{"AnonymousA", "AnonymousB"},
							Fields: map[string]string{
								"FieldB": "",
								"FieldA": "",
							},
						},
					},
				},
			},
		},
		BreakingChanges: &report.BreakingChanges{
			Removed: &delta.Content{
				Content: exports.Content{
					Structs: map[string]exports.Struct{
						"RemovedStruct": {
							AnonymousFields: []string{"RemovedAnonymousA", "RemovedAnonymousB"},
							Fields: map[string]string{
								"RemovedFieldB": "",
								"RemovedFieldA": "",
							},
						},
					},
				},
			},
		},
	}
	actual := writeChangelogForPackage(r)
	expected := "### Breaking Changes\n\n- Field `RemovedAnonymousA`, `RemovedAnonymousB` of struct `RemovedStruct` has been removed\n- Field `RemovedFieldA`, `RemovedFieldB` of struct `RemovedStruct` has been removed\n\n### Features Added\n\n- New anonymous field `AnonymousA`, `AnonymousB` in struct `Struct`\n- New field `FieldA`, `FieldB` in struct `Struct`\n"
	assert.Equal(t, expected, actual)
}
