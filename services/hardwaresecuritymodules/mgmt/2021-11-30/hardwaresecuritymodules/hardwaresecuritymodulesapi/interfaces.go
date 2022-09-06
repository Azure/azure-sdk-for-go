package hardwaresecuritymodulesapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/hardwaresecuritymodules/mgmt/2021-11-30/hardwaresecuritymodules"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result hardwaresecuritymodules.DedicatedHsmOperationListResult, err error)
}

var _ OperationsClientAPI = (*hardwaresecuritymodules.OperationsClient)(nil)

// DedicatedHsmClientAPI contains the set of methods on the DedicatedHsmClient type.
type DedicatedHsmClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, name string, parameters hardwaresecuritymodules.DedicatedHsm) (result hardwaresecuritymodules.DedicatedHsmCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, name string) (result hardwaresecuritymodules.DedicatedHsmDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, name string) (result hardwaresecuritymodules.DedicatedHsm, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, top *int32) (result hardwaresecuritymodules.DedicatedHsmListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string, top *int32) (result hardwaresecuritymodules.DedicatedHsmListResultIterator, err error)
	ListBySubscription(ctx context.Context, top *int32) (result hardwaresecuritymodules.DedicatedHsmListResultPage, err error)
	ListBySubscriptionComplete(ctx context.Context, top *int32) (result hardwaresecuritymodules.DedicatedHsmListResultIterator, err error)
	ListOutboundNetworkDependenciesEndpoints(ctx context.Context, resourceGroupName string, name string) (result hardwaresecuritymodules.OutboundEnvironmentEndpointCollectionPage, err error)
	ListOutboundNetworkDependenciesEndpointsComplete(ctx context.Context, resourceGroupName string, name string) (result hardwaresecuritymodules.OutboundEnvironmentEndpointCollectionIterator, err error)
	Update(ctx context.Context, resourceGroupName string, name string, parameters hardwaresecuritymodules.DedicatedHsmPatchParameters) (result hardwaresecuritymodules.DedicatedHsmUpdateFuture, err error)
}

var _ DedicatedHsmClientAPI = (*hardwaresecuritymodules.DedicatedHsmClient)(nil)
