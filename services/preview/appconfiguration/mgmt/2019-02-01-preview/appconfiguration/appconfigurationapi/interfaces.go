package appconfigurationapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/appconfiguration/mgmt/2019-02-01-preview/appconfiguration"
)

// ConfigurationStoresClientAPI contains the set of methods on the ConfigurationStoresClient type.
type ConfigurationStoresClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, configStoreName string, configStoreCreationParameters appconfiguration.ConfigurationStore) (result appconfiguration.ConfigurationStoresCreateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, configStoreName string) (result appconfiguration.ConfigurationStoresDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, configStoreName string) (result appconfiguration.ConfigurationStore, err error)
	List(ctx context.Context, skipToken string) (result appconfiguration.ConfigurationStoreListResultPage, err error)
	ListComplete(ctx context.Context, skipToken string) (result appconfiguration.ConfigurationStoreListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, skipToken string) (result appconfiguration.ConfigurationStoreListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string, skipToken string) (result appconfiguration.ConfigurationStoreListResultIterator, err error)
	ListKeys(ctx context.Context, resourceGroupName string, configStoreName string, skipToken string) (result appconfiguration.APIKeyListResultPage, err error)
	ListKeysComplete(ctx context.Context, resourceGroupName string, configStoreName string, skipToken string) (result appconfiguration.APIKeyListResultIterator, err error)
	ListKeyValue(ctx context.Context, resourceGroupName string, configStoreName string, listKeyValueParameters appconfiguration.ListKeyValueParameters) (result appconfiguration.KeyValue, err error)
	RegenerateKey(ctx context.Context, resourceGroupName string, configStoreName string, regenerateKeyParameters appconfiguration.RegenerateKeyParameters) (result appconfiguration.APIKey, err error)
	Update(ctx context.Context, resourceGroupName string, configStoreName string, configStoreUpdateParameters appconfiguration.ConfigurationStoreUpdateParameters) (result appconfiguration.ConfigurationStoresUpdateFuture, err error)
}

var _ ConfigurationStoresClientAPI = (*appconfiguration.ConfigurationStoresClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	CheckNameAvailability(ctx context.Context, checkNameAvailabilityParameters appconfiguration.CheckNameAvailabilityParameters) (result appconfiguration.NameAvailabilityStatus, err error)
	List(ctx context.Context, skipToken string) (result appconfiguration.OperationDefinitionListResultPage, err error)
	ListComplete(ctx context.Context, skipToken string) (result appconfiguration.OperationDefinitionListResultIterator, err error)
}

var _ OperationsClientAPI = (*appconfiguration.OperationsClient)(nil)
