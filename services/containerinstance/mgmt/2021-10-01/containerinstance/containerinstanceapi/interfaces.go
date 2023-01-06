// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package containerinstanceapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2021-10-01/containerinstance"
	"github.com/Azure/go-autorest/autorest"
)

// ContainerGroupsClientAPI contains the set of methods on the ContainerGroupsClient type.
type ContainerGroupsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, containerGroupName string, containerGroup containerinstance.ContainerGroup) (result containerinstance.ContainerGroupsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, containerGroupName string) (result containerinstance.ContainerGroupsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, containerGroupName string) (result containerinstance.ContainerGroup, err error)
	GetOutboundNetworkDependenciesEndpoints(ctx context.Context, resourceGroupName string, containerGroupName string) (result containerinstance.ListString, err error)
	List(ctx context.Context) (result containerinstance.ContainerGroupListResultPage, err error)
	ListComplete(ctx context.Context) (result containerinstance.ContainerGroupListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result containerinstance.ContainerGroupListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result containerinstance.ContainerGroupListResultIterator, err error)
	Restart(ctx context.Context, resourceGroupName string, containerGroupName string) (result containerinstance.ContainerGroupsRestartFuture, err error)
	Start(ctx context.Context, resourceGroupName string, containerGroupName string) (result containerinstance.ContainerGroupsStartFuture, err error)
	Stop(ctx context.Context, resourceGroupName string, containerGroupName string) (result autorest.Response, err error)
	Update(ctx context.Context, resourceGroupName string, containerGroupName string, resource containerinstance.Resource) (result containerinstance.ContainerGroup, err error)
}

var _ ContainerGroupsClientAPI = (*containerinstance.ContainerGroupsClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result containerinstance.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result containerinstance.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*containerinstance.OperationsClient)(nil)

// LocationClientAPI contains the set of methods on the LocationClient type.
type LocationClientAPI interface {
	ListCachedImages(ctx context.Context, location string) (result containerinstance.CachedImagesListResultPage, err error)
	ListCachedImagesComplete(ctx context.Context, location string) (result containerinstance.CachedImagesListResultIterator, err error)
	ListCapabilities(ctx context.Context, location string) (result containerinstance.CapabilitiesListResultPage, err error)
	ListCapabilitiesComplete(ctx context.Context, location string) (result containerinstance.CapabilitiesListResultIterator, err error)
	ListUsage(ctx context.Context, location string) (result containerinstance.UsageListResult, err error)
}

var _ LocationClientAPI = (*containerinstance.LocationClient)(nil)

// ContainersClientAPI contains the set of methods on the ContainersClient type.
type ContainersClientAPI interface {
	Attach(ctx context.Context, resourceGroupName string, containerGroupName string, containerName string) (result containerinstance.ContainerAttachResponse, err error)
	ExecuteCommand(ctx context.Context, resourceGroupName string, containerGroupName string, containerName string, containerExecRequest containerinstance.ContainerExecRequest) (result containerinstance.ContainerExecResponse, err error)
	ListLogs(ctx context.Context, resourceGroupName string, containerGroupName string, containerName string, tail *int32, timestamps *bool) (result containerinstance.Logs, err error)
}

var _ ContainersClientAPI = (*containerinstance.ContainersClient)(nil)

// SubnetServiceAssociationLinkClientAPI contains the set of methods on the SubnetServiceAssociationLinkClient type.
type SubnetServiceAssociationLinkClientAPI interface {
	Delete(ctx context.Context, resourceGroupName string, virtualNetworkName string, subnetName string) (result containerinstance.SubnetServiceAssociationLinkDeleteFuture, err error)
}

var _ SubnetServiceAssociationLinkClientAPI = (*containerinstance.SubnetServiceAssociationLinkClient)(nil)
