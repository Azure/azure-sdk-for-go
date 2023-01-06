// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kubernetesconfiguration/armkubernetesconfiguration](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kubernetesconfiguration/armkubernetesconfiguration). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package kubernetesconfigurationapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/kubernetesconfiguration/mgmt/2022-04-02-preview/kubernetesconfiguration"
)

// ExtensionsClientAPI contains the set of methods on the ExtensionsClient type.
type ExtensionsClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, extensionName string, extension kubernetesconfiguration.Extension) (result kubernetesconfiguration.ExtensionsCreateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, extensionName string, forceDelete *bool) (result kubernetesconfiguration.ExtensionsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, extensionName string) (result kubernetesconfiguration.Extension, err error)
	List(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.ExtensionsListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.ExtensionsListIterator, err error)
	Update(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, extensionName string, patchExtension kubernetesconfiguration.PatchExtension) (result kubernetesconfiguration.ExtensionsUpdateFuture, err error)
}

var _ ExtensionsClientAPI = (*kubernetesconfiguration.ExtensionsClient)(nil)

// OperationStatusClientAPI contains the set of methods on the OperationStatusClient type.
type OperationStatusClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, extensionName string, operationID string) (result kubernetesconfiguration.OperationStatusResult, err error)
	List(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.OperationStatusListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.OperationStatusListIterator, err error)
}

var _ OperationStatusClientAPI = (*kubernetesconfiguration.OperationStatusClient)(nil)

// FluxConfigurationsClientAPI contains the set of methods on the FluxConfigurationsClient type.
type FluxConfigurationsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, fluxConfigurationName string, fluxConfiguration kubernetesconfiguration.FluxConfiguration) (result kubernetesconfiguration.FluxConfigurationsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, fluxConfigurationName string, forceDelete *bool) (result kubernetesconfiguration.FluxConfigurationsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, fluxConfigurationName string) (result kubernetesconfiguration.FluxConfiguration, err error)
	List(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.FluxConfigurationsListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.FluxConfigurationsListIterator, err error)
	Update(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, fluxConfigurationName string, fluxConfigurationPatch kubernetesconfiguration.FluxConfigurationPatch) (result kubernetesconfiguration.FluxConfigurationsUpdateFuture, err error)
}

var _ FluxConfigurationsClientAPI = (*kubernetesconfiguration.FluxConfigurationsClient)(nil)

// FluxConfigOperationStatusClientAPI contains the set of methods on the FluxConfigOperationStatusClient type.
type FluxConfigOperationStatusClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, fluxConfigurationName string, operationID string) (result kubernetesconfiguration.OperationStatusResult, err error)
}

var _ FluxConfigOperationStatusClientAPI = (*kubernetesconfiguration.FluxConfigOperationStatusClient)(nil)

// SourceControlConfigurationsClientAPI contains the set of methods on the SourceControlConfigurationsClient type.
type SourceControlConfigurationsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, sourceControlConfigurationName string, sourceControlConfiguration kubernetesconfiguration.SourceControlConfiguration) (result kubernetesconfiguration.SourceControlConfiguration, err error)
	Delete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, sourceControlConfigurationName string) (result kubernetesconfiguration.SourceControlConfigurationsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, sourceControlConfigurationName string) (result kubernetesconfiguration.SourceControlConfiguration, err error)
	List(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.SourceControlConfigurationListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string) (result kubernetesconfiguration.SourceControlConfigurationListIterator, err error)
}

var _ SourceControlConfigurationsClientAPI = (*kubernetesconfiguration.SourceControlConfigurationsClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result kubernetesconfiguration.ResourceProviderOperationListPage, err error)
	ListComplete(ctx context.Context) (result kubernetesconfiguration.ResourceProviderOperationListIterator, err error)
}

var _ OperationsClientAPI = (*kubernetesconfiguration.OperationsClient)(nil)

// PrivateLinkScopesClientAPI contains the set of methods on the PrivateLinkScopesClient type.
type PrivateLinkScopesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, scopeName string, parameters kubernetesconfiguration.PrivateLinkScope) (result kubernetesconfiguration.PrivateLinkScope, err error)
	Delete(ctx context.Context, resourceGroupName string, scopeName string) (result kubernetesconfiguration.PrivateLinkScopesDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, scopeName string) (result kubernetesconfiguration.PrivateLinkScope, err error)
	List(ctx context.Context) (result kubernetesconfiguration.PrivateLinkScopeListResultPage, err error)
	ListComplete(ctx context.Context) (result kubernetesconfiguration.PrivateLinkScopeListResultIterator, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result kubernetesconfiguration.PrivateLinkScopeListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result kubernetesconfiguration.PrivateLinkScopeListResultIterator, err error)
	UpdateTags(ctx context.Context, resourceGroupName string, scopeName string, privateLinkScopeTags kubernetesconfiguration.TagsResource) (result kubernetesconfiguration.PrivateLinkScope, err error)
}

var _ PrivateLinkScopesClientAPI = (*kubernetesconfiguration.PrivateLinkScopesClient)(nil)

// PrivateLinkResourcesClientAPI contains the set of methods on the PrivateLinkResourcesClient type.
type PrivateLinkResourcesClientAPI interface {
	Get(ctx context.Context, resourceGroupName string, scopeName string, groupName string) (result kubernetesconfiguration.PrivateLinkResource, err error)
	ListByPrivateLinkScope(ctx context.Context, resourceGroupName string, scopeName string) (result kubernetesconfiguration.PrivateLinkResourceListResult, err error)
}

var _ PrivateLinkResourcesClientAPI = (*kubernetesconfiguration.PrivateLinkResourcesClient)(nil)

// PrivateEndpointConnectionsClientAPI contains the set of methods on the PrivateEndpointConnectionsClient type.
type PrivateEndpointConnectionsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, scopeName string, privateEndpointConnectionName string, properties kubernetesconfiguration.PrivateEndpointConnection) (result kubernetesconfiguration.PrivateEndpointConnectionsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, scopeName string, privateEndpointConnectionName string) (result kubernetesconfiguration.PrivateEndpointConnectionsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, scopeName string, privateEndpointConnectionName string) (result kubernetesconfiguration.PrivateEndpointConnection, err error)
	ListByPrivateLinkScope(ctx context.Context, resourceGroupName string, scopeName string) (result kubernetesconfiguration.PrivateEndpointConnectionListResult, err error)
}

var _ PrivateEndpointConnectionsClientAPI = (*kubernetesconfiguration.PrivateEndpointConnectionsClient)(nil)
